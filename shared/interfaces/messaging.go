package interfaces

import (
	"context"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/types"
)

// MessageProducer defines the interface for message producers
type MessageProducer interface {
	// Start starts the producer
	Start(ctx context.Context) error

	// Stop stops the producer gracefully
	Stop(ctx context.Context) error

	// Close closes the producer and releases resources
	Close() error

	// SendTxSubmission sends a transaction submission message
	SendTxSubmission(ctx context.Context, msg *types.TxSubmissionMessage) error

	// SendBlockExecution sends a block execution message
	SendBlockExecution(ctx context.Context, msg *types.BlockExecutionMessage) error

	// SendStateCommit sends a state commit message
	SendStateCommit(ctx context.Context, msg *types.StateCommitMessage) error

	// SendCacheInvalidation broadcasts a cache invalidation message
	SendCacheInvalidation(ctx context.Context, msg *types.CacheInvalidationMessage) error

	// SendWithRetry sends a message with retry logic
	SendWithRetry(ctx context.Context, topic string, message interface{}, maxRetries int, backoff time.Duration) error

	// GetMetrics returns producer metrics
	GetMetrics() *ProducerMetrics

	// HealthCheck performs producer health check
	HealthCheck(ctx context.Context) error
}

// MessageConsumer defines the interface for message consumers
type MessageConsumer interface {
	// Start starts the consumer
	Start(ctx context.Context) error

	// Stop stops the consumer gracefully
	Stop(ctx context.Context) error

	// Close closes the consumer and releases resources
	Close() error

	// Subscribe subscribes to a topic with clustering mode
	Subscribe(ctx context.Context, topic string, consumerGroup string, handler types.MessageHandler) error

	// SubscribeBroadcast subscribes to a topic with broadcast mode
	SubscribeBroadcast(ctx context.Context, topic string, handler types.MessageHandler) error

	// Unsubscribe unsubscribes from a topic
	Unsubscribe(ctx context.Context, topic string) error

	// Pause pauses message consumption
	Pause(ctx context.Context) error

	// Resume resumes message consumption
	Resume(ctx context.Context) error

	// GetMetrics returns consumer metrics
	GetMetrics() *ConsumerMetrics

	// HealthCheck performs consumer health check
	HealthCheck(ctx context.Context) error

	// GetConsumerLag returns the consumer lag for topics
	GetConsumerLag(ctx context.Context) (map[string]int64, error)
}

// MessageBroker defines the interface for message broker management
type MessageBroker interface {
	// CreateTopic creates a new topic
	CreateTopic(ctx context.Context, topic string, config *TopicConfig) error

	// DeleteTopic deletes a topic
	DeleteTopic(ctx context.Context, topic string) error

	// ListTopics lists all available topics
	ListTopics(ctx context.Context) ([]string, error)

	// GetTopicInfo returns information about a topic
	GetTopicInfo(ctx context.Context, topic string) (*TopicInfo, error)

	// CreateConsumerGroup creates a consumer group
	CreateConsumerGroup(ctx context.Context, group string, topics []string) error

	// DeleteConsumerGroup deletes a consumer group
	DeleteConsumerGroup(ctx context.Context, group string) error

	// ListConsumerGroups lists all consumer groups
	ListConsumerGroups(ctx context.Context) ([]string, error)

	// ResetConsumerGroupOffset resets consumer group offset
	ResetConsumerGroupOffset(ctx context.Context, group string, topic string, offset int64) error

	// GetBrokerMetrics returns broker metrics
	GetBrokerMetrics(ctx context.Context) (*BrokerMetrics, error)

	// HealthCheck performs broker health check
	HealthCheck(ctx context.Context) error
}

// MessageTransactionCoordinator defines interface for transactional messaging
type MessageTransactionCoordinator interface {
	// BeginTransaction starts a message transaction
	BeginTransaction(ctx context.Context) (MessageTransaction, error)

	// PrepareTransaction prepares a transaction for commit
	PrepareTransaction(ctx context.Context, txID string) error

	// CommitTransaction commits a prepared transaction
	CommitTransaction(ctx context.Context, txID string) error

	// RollbackTransaction rolls back a transaction
	RollbackTransaction(ctx context.Context, txID string) error

	// GetTransactionStatus returns transaction status
	GetTransactionStatus(ctx context.Context, txID string) (TransactionStatus, error)

	// CleanupExpiredTransactions cleans up expired transactions
	CleanupExpiredTransactions(ctx context.Context, maxAge time.Duration) error
}

// MessageTransaction represents a transactional message session
type MessageTransaction interface {
	// GetID returns the transaction ID
	GetID() string

	// Send sends a message within the transaction
	Send(ctx context.Context, topic string, message interface{}) error

	// Prepare prepares the transaction
	Prepare(ctx context.Context) error

	// Commit commits the transaction
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction
	Rollback(ctx context.Context) error

	// GetStatus returns transaction status
	GetStatus() TransactionStatus

	// GetMessages returns messages in this transaction
	GetMessages() []TransactionalMessage
}

// TopicConfig defines configuration for a topic
type TopicConfig struct {
	Name              string            `json:"name"`
	Partitions        int               `json:"partitions"`
	ReplicationFactor int               `json:"replicationFactor"`
	RetentionTime     time.Duration     `json:"retentionTime"`
	MaxMessageSize    int64             `json:"maxMessageSize"`
	CompressionType   string            `json:"compressionType"`
	Properties        map[string]string `json:"properties"`
}

// TopicInfo contains information about a topic
type TopicInfo struct {
	Name              string            `json:"name"`
	Partitions        []PartitionInfo   `json:"partitions"`
	ReplicationFactor int               `json:"replicationFactor"`
	RetentionTime     time.Duration     `json:"retentionTime"`
	TotalMessages     int64             `json:"totalMessages"`
	TotalSize         int64             `json:"totalSize"`
	Properties        map[string]string `json:"properties"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

// PartitionInfo contains information about a topic partition
type PartitionInfo struct {
	ID             int   `json:"id"`
	Leader         int   `json:"leader"`
	Replicas       []int `json:"replicas"`
	ISR            []int `json:"isr"`
	MessageCount   int64 `json:"messageCount"`
	Size           int64 `json:"size"`
	EarliestOffset int64 `json:"earliestOffset"`
	LatestOffset   int64 `json:"latestOffset"`
}

// ProducerMetrics defines metrics for message producers
type ProducerMetrics struct {
	// Message metrics
	TotalMessages     uint64            `json:"totalMessages"`
	MessagesPerSecond float64           `json:"messagesPerSecond"`
	MessagesByTopic   map[string]uint64 `json:"messagesByTopic"`

	// Performance metrics
	AverageLatency int64   `json:"averageLatency"` // in milliseconds
	P99Latency     int64   `json:"p99Latency"`     // in milliseconds
	ThroughputMBps float64 `json:"throughputMbps"` // MB per second

	// Error metrics
	TotalErrors   uint64            `json:"totalErrors"`
	ErrorsByType  map[string]uint64 `json:"errorsByType"`
	RetryAttempts uint64            `json:"retryAttempts"`

	// Connection metrics
	ActiveConnections int       `json:"activeConnections"`
	ConnectionRetries uint64    `json:"connectionRetries"`
	LastConnected     time.Time `json:"lastConnected"`

	// Buffer metrics
	BufferUsage     float64 `json:"bufferUsage"` // percentage
	BufferSize      int64   `json:"bufferSize"`  // in bytes
	PendingMessages int64   `json:"pendingMessages"`
}

// ConsumerMetrics defines metrics for message consumers
type ConsumerMetrics struct {
	// Message metrics
	TotalMessages     uint64            `json:"totalMessages"`
	MessagesPerSecond float64           `json:"messagesPerSecond"`
	MessagesByTopic   map[string]uint64 `json:"messagesByTopic"`

	// Processing metrics
	AverageProcessTime int64   `json:"averageProcessTime"` // in milliseconds
	P99ProcessTime     int64   `json:"p99ProcessTime"`     // in milliseconds
	ProcessingRate     float64 `json:"processingRate"`     // messages per second

	// Error metrics
	TotalErrors           uint64            `json:"totalErrors"`
	ProcessingErrors      uint64            `json:"processingErrors"`
	DeserializationErrors uint64            `json:"deserializationErrors"`
	ErrorsByType          map[string]uint64 `json:"errorsByType"`

	// Consumer lag metrics
	ConsumerLag    map[string]int64 `json:"consumerLag"` // by topic
	MaxConsumerLag int64            `json:"maxConsumerLag"`

	// Connection metrics
	ActiveConnections int       `json:"activeConnections"`
	RebalanceCount    uint64    `json:"rebalanceCount"`
	LastRebalance     time.Time `json:"lastRebalance"`
}

// BrokerMetrics defines metrics for message brokers
type BrokerMetrics struct {
	// Broker info
	BrokerID int           `json:"brokerId"`
	Version  string        `json:"version"`
	Uptime   time.Duration `json:"uptime"`

	// Topic metrics
	TopicCount     int `json:"topicCount"`
	PartitionCount int `json:"partitionCount"`
	LeaderCount    int `json:"leaderCount"`

	// Message metrics
	TotalMessages     int64   `json:"totalMessages"`
	MessagesInPerSec  float64 `json:"messagesInPerSec"`
	MessagesOutPerSec float64 `json:"messagesOutPerSec"`
	BytesInPerSec     float64 `json:"bytesInPerSec"`
	BytesOutPerSec    float64 `json:"bytesOutPerSec"`

	// Storage metrics
	TotalLogSize     int64         `json:"totalLogSize"`     // in bytes
	LogRetentionSize int64         `json:"logRetentionSize"` // in bytes
	LogRetentionTime time.Duration `json:"logRetentionTime"`

	// Performance metrics
	RequestHandlerAvgIdlePercent   float64 `json:"requestHandlerAvgIdlePercent"`
	NetworkProcessorAvgIdlePercent float64 `json:"networkProcessorAvgIdlePercent"`

	// Connection metrics
	ActiveConnections int `json:"activeConnections"`
	ProducerCount     int `json:"producerCount"`
	ConsumerCount     int `json:"consumerCount"`

	// Error metrics
	TotalErrors           uint64 `json:"totalErrors"`
	FailedFetchRequests   uint64 `json:"failedFetchRequests"`
	FailedProduceRequests uint64 `json:"failedProduceRequests"`
}

// TransactionStatus represents the status of a message transaction
type TransactionStatus string

const (
	TransactionStatusBegin      TransactionStatus = "begin"
	TransactionStatusPrepared   TransactionStatus = "prepared"
	TransactionStatusCommitted  TransactionStatus = "committed"
	TransactionStatusRolledBack TransactionStatus = "rolledback"
	TransactionStatusExpired    TransactionStatus = "expired"
	TransactionStatusUnknown    TransactionStatus = "unknown"
)

// TransactionalMessage represents a message within a transaction
type TransactionalMessage struct {
	Topic     string      `json:"topic"`
	Message   interface{} `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
	MessageID string      `json:"messageId"`
}

// MessageDeliveryGuarantee defines message delivery guarantees
type MessageDeliveryGuarantee string

const (
	DeliveryAtMostOnce  MessageDeliveryGuarantee = "at_most_once"
	DeliveryAtLeastOnce MessageDeliveryGuarantee = "at_least_once"
	DeliveryExactlyOnce MessageDeliveryGuarantee = "exactly_once"
)

// MessageOrderingGuarantee defines message ordering guarantees
type MessageOrderingGuarantee string

const (
	OrderingNone      MessageOrderingGuarantee = "none"
	OrderingPartition MessageOrderingGuarantee = "partition"
	OrderingGlobal    MessageOrderingGuarantee = "global"
)

// MessageQueueConfig defines configuration for message queue
type MessageQueueConfig struct {
	// Broker configuration
	Brokers   []string `json:"brokers"`
	AccessKey string   `json:"accessKey"`
	SecretKey string   `json:"secretKey"`

	// Connection configuration
	ConnectionTimeout time.Duration `json:"connectionTimeout"`
	ReadTimeout       time.Duration `json:"readTimeout"`
	WriteTimeout      time.Duration `json:"writeTimeout"`
	KeepAlive         time.Duration `json:"keepAlive"`
	MaxConnections    int           `json:"maxConnections"`

	// Producer configuration
	ProducerBatchSize    int           `json:"producerBatchSize"`
	ProducerFlushTimeout time.Duration `json:"producerFlushTimeout"`
	ProducerRetries      int           `json:"producerRetries"`
	ProducerRetryBackoff time.Duration `json:"producerRetryBackoff"`

	// Consumer configuration
	ConsumerFetchSize    int           `json:"consumerFetchSize"`
	ConsumerFetchTimeout time.Duration `json:"consumerFetchTimeout"`
	ConsumerRetries      int           `json:"consumerRetries"`
	ConsumerRetryBackoff time.Duration `json:"consumerRetryBackoff"`

	// Delivery guarantees
	DeliveryGuarantee MessageDeliveryGuarantee `json:"deliveryGuarantee"`
	OrderingGuarantee MessageOrderingGuarantee `json:"orderingGuarantee"`

	// Security configuration
	EnableSSL         bool   `json:"enableSsl"`
	SSLCertFile       string `json:"sslCertFile"`
	SSLKeyFile        string `json:"sslKeyFile"`
	SSLCAFile         string `json:"sslCaFile"`
	SSLVerifyHostname bool   `json:"sslVerifyHostname"`

	// Monitoring configuration
	EnableMetrics   bool          `json:"enableMetrics"`
	MetricsInterval time.Duration `json:"metricsInterval"`
}

// MessageFilter defines interface for message filtering
type MessageFilter interface {
	// Filter filters messages based on criteria
	Filter(ctx context.Context, message interface{}) (bool, error)

	// GetCriteria returns filter criteria
	GetCriteria() map[string]interface{}

	// SetCriteria sets filter criteria
	SetCriteria(criteria map[string]interface{}) error
}

// MessageTransformer defines interface for message transformation
type MessageTransformer interface {
	// Transform transforms a message
	Transform(ctx context.Context, message interface{}) (interface{}, error)

	// GetTransformationType returns transformation type
	GetTransformationType() string
}

// MessageRouter defines interface for message routing
type MessageRouter interface {
	// Route determines the destination topic for a message
	Route(ctx context.Context, message interface{}) (string, error)

	// GetRoutingRules returns current routing rules
	GetRoutingRules() []RoutingRule

	// AddRoutingRule adds a new routing rule
	AddRoutingRule(rule RoutingRule) error

	// RemoveRoutingRule removes a routing rule
	RemoveRoutingRule(ruleID string) error
}

// RoutingRule defines a message routing rule
type RoutingRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Condition   string                 `json:"condition"`   // expression to evaluate
	Destination string                 `json:"destination"` // target topic
	Priority    int                    `json:"priority"`
	Enabled     bool                   `json:"enabled"`
	Properties  map[string]interface{} `json:"properties"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}
