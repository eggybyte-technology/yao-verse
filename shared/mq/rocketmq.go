package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/eggybyte-technology/yao-verse-shared/types"
)

// RocketMQProducer implements MessageProducer interface using RocketMQ
type RocketMQProducer struct {
	producer rocketmq.Producer
	config   *RocketMQConfig
}

// RocketMQConsumer implements MessageConsumer interface using RocketMQ
type RocketMQConsumer struct {
	consumers map[string]rocketmq.PushConsumer // topic -> consumer
	config    *RocketMQConfig
	handlers  map[string]types.MessageHandler // topic -> handler
	mutex     sync.RWMutex
}

// RocketMQConfig defines configuration for RocketMQ
type RocketMQConfig struct {
	Endpoints     []string `json:"endpoints"`
	AccessKey     string   `json:"accessKey"`
	SecretKey     string   `json:"secretKey"`
	ProducerGroup string   `json:"producerGroup"`
	ConsumerGroup string   `json:"consumerGroup"`
	Namespace     string   `json:"namespace"`
	RetryTimes    int      `json:"retryTimes"`
}

// NewRocketMQProducer creates a new RocketMQ producer
func NewRocketMQProducer(config *RocketMQConfig) (*RocketMQProducer, error) {
	opts := make([]producer.Option, 0)
	opts = append(opts, producer.WithNameServer(config.Endpoints))
	opts = append(opts, producer.WithRetry(config.RetryTimes))
	opts = append(opts, producer.WithGroupName(config.ProducerGroup))

	if config.Namespace != "" {
		opts = append(opts, producer.WithNamespace(config.Namespace))
	}

	if config.AccessKey != "" && config.SecretKey != "" {
		opts = append(opts, producer.WithCredentials(primitive.Credentials{
			AccessKey: config.AccessKey,
			SecretKey: config.SecretKey,
		}))
	}

	prod, err := rocketmq.NewProducer(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create RocketMQ producer: %w", err)
	}

	err = prod.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start RocketMQ producer: %w", err)
	}

	return &RocketMQProducer{
		producer: prod,
		config:   config,
	}, nil
}

// NewRocketMQConsumer creates a new RocketMQ consumer
func NewRocketMQConsumer(config *RocketMQConfig) (*RocketMQConsumer, error) {
	return &RocketMQConsumer{
		consumers: make(map[string]rocketmq.PushConsumer),
		config:    config,
		handlers:  make(map[string]types.MessageHandler),
	}, nil
}

// SendTxSubmission sends a transaction submission message
func (p *RocketMQProducer) SendTxSubmission(msg *types.TxSubmissionMessage) error {
	return p.sendMessage(types.TopicTxIngress, msg.MessageID, msg)
}

// SendBlockExecution sends a block execution message
func (p *RocketMQProducer) SendBlockExecution(msg *types.BlockExecutionMessage) error {
	return p.sendMessage(types.TopicBlockExecute, msg.MessageID, msg)
}

// SendStateCommit sends a state commit message
func (p *RocketMQProducer) SendStateCommit(msg *types.StateCommitMessage) error {
	return p.sendMessage(types.TopicStateCommit, msg.MessageID, msg)
}

// SendCacheInvalidation broadcasts a cache invalidation message
func (p *RocketMQProducer) SendCacheInvalidation(msg *types.CacheInvalidationMessage) error {
	// Generate a unique message ID based on block number and timestamp
	messageID := fmt.Sprintf("cache-invalidation-%d-%d", msg.BlockNumber, msg.Timestamp.Unix())
	return p.sendMessage(types.TopicCacheInvalidation, messageID, msg)
}

// Close closes the producer
func (p *RocketMQProducer) Close() error {
	return p.producer.Shutdown()
}

// sendMessage sends a message to the specified topic
func (p *RocketMQProducer) sendMessage(topic, messageID string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}

	// Set message key for routing and tracing
	if messageID != "" {
		msg.WithKeys([]string{messageID})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := p.producer.SendSync(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send message to topic %s: %w", topic, err)
	}

	log.Printf("Message sent successfully to topic %s, msgId: %s", topic, result.MsgID)
	return nil
}

// Subscribe subscribes to a topic with a handler (cluster mode)
func (c *RocketMQConsumer) Subscribe(topic string, consumerGroup string, handler types.MessageHandler) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if already subscribed to this topic
	if _, exists := c.consumers[topic]; exists {
		return fmt.Errorf("already subscribed to topic: %s", topic)
	}

	opts := make([]consumer.Option, 0)
	opts = append(opts, consumer.WithNameServer(c.config.Endpoints))
	opts = append(opts, consumer.WithGroupName(consumerGroup))
	opts = append(opts, consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset))
	opts = append(opts, consumer.WithConsumerModel(consumer.Clustering)) // Cluster mode for load balancing

	if c.config.Namespace != "" {
		opts = append(opts, consumer.WithNamespace(c.config.Namespace))
	}

	if c.config.AccessKey != "" && c.config.SecretKey != "" {
		opts = append(opts, consumer.WithCredentials(primitive.Credentials{
			AccessKey: c.config.AccessKey,
			SecretKey: c.config.SecretKey,
		}))
	}

	pushConsumer, err := rocketmq.NewPushConsumer(opts...)
	if err != nil {
		return fmt.Errorf("failed to create consumer for topic %s: %w", topic, err)
	}

	// Subscribe to topic with message handler
	err = pushConsumer.Subscribe(topic, consumer.MessageSelector{}, c.createMessageHandler(handler))
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	c.consumers[topic] = pushConsumer
	c.handlers[topic] = handler

	log.Printf("Subscribed to topic %s with consumer group %s", topic, consumerGroup)
	return nil
}

// SubscribeBroadcast subscribes to a broadcast topic with a handler (broadcast mode)
func (c *RocketMQConsumer) SubscribeBroadcast(topic string, handler types.MessageHandler) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if already subscribed to this topic
	if _, exists := c.consumers[topic]; exists {
		return fmt.Errorf("already subscribed to topic: %s", topic)
	}

	opts := make([]consumer.Option, 0)
	opts = append(opts, consumer.WithNameServer(c.config.Endpoints))
	opts = append(opts, consumer.WithGroupName(c.config.ConsumerGroup))
	opts = append(opts, consumer.WithConsumeFromWhere(consumer.ConsumeFromLastOffset))
	opts = append(opts, consumer.WithConsumerModel(consumer.BroadCasting)) // Broadcast mode

	if c.config.Namespace != "" {
		opts = append(opts, consumer.WithNamespace(c.config.Namespace))
	}

	if c.config.AccessKey != "" && c.config.SecretKey != "" {
		opts = append(opts, consumer.WithCredentials(primitive.Credentials{
			AccessKey: c.config.AccessKey,
			SecretKey: c.config.SecretKey,
		}))
	}

	pushConsumer, err := rocketmq.NewPushConsumer(opts...)
	if err != nil {
		return fmt.Errorf("failed to create broadcast consumer for topic %s: %w", topic, err)
	}

	// Subscribe to topic with message handler
	err = pushConsumer.Subscribe(topic, consumer.MessageSelector{}, c.createMessageHandler(handler))
	if err != nil {
		return fmt.Errorf("failed to subscribe to broadcast topic %s: %w", topic, err)
	}

	c.consumers[topic] = pushConsumer
	c.handlers[topic] = handler

	log.Printf("Subscribed to broadcast topic %s", topic)
	return nil
}

// Start starts consuming messages
func (c *RocketMQConsumer) Start() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for topic, consumer := range c.consumers {
		err := consumer.Start()
		if err != nil {
			return fmt.Errorf("failed to start consumer for topic %s: %w", topic, err)
		}
		log.Printf("Started consumer for topic %s", topic)
	}

	return nil
}

// Stop stops consuming messages
func (c *RocketMQConsumer) Stop() error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for topic, consumer := range c.consumers {
		err := consumer.Shutdown()
		if err != nil {
			log.Printf("Error shutting down consumer for topic %s: %v", topic, err)
		} else {
			log.Printf("Stopped consumer for topic %s", topic)
		}
	}

	return nil
}

// Close closes the consumer
func (c *RocketMQConsumer) Close() error {
	return c.Stop()
}

// createMessageHandler creates a RocketMQ message handler that wraps our MessageHandler
func (c *RocketMQConsumer) createMessageHandler(handler types.MessageHandler) func(context.Context, ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	return func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			err := c.processMessage(msg, handler)
			if err != nil {
				log.Printf("Error processing message from topic %s: %v", msg.Topic, err)
				return consumer.ConsumeRetryLater, err
			}
		}
		return consumer.ConsumeSuccess, nil
	}
}

// processMessage processes a single message based on its topic
func (c *RocketMQConsumer) processMessage(msg *primitive.MessageExt, handler types.MessageHandler) error {
	switch msg.Topic {
	case types.TopicTxIngress:
		var txMsg types.TxSubmissionMessage
		err := json.Unmarshal(msg.Body, &txMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal tx submission message: %w", err)
		}
		return handler.HandleTxSubmission(&txMsg)

	case types.TopicBlockExecute:
		var blockMsg types.BlockExecutionMessage
		err := json.Unmarshal(msg.Body, &blockMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal block execution message: %w", err)
		}
		return handler.HandleBlockExecution(&blockMsg)

	case types.TopicStateCommit:
		var stateMsg types.StateCommitMessage
		err := json.Unmarshal(msg.Body, &stateMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal state commit message: %w", err)
		}
		return handler.HandleStateCommit(&stateMsg)

	case types.TopicCacheInvalidation:
		var invalidationMsg types.CacheInvalidationMessage
		err := json.Unmarshal(msg.Body, &invalidationMsg)
		if err != nil {
			return fmt.Errorf("failed to unmarshal cache invalidation message: %w", err)
		}
		return handler.HandleCacheInvalidation(&invalidationMsg)

	default:
		return fmt.Errorf("unknown topic: %s", msg.Topic)
	}
}
