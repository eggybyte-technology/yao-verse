package types

import "time"

// MessageType defines the type of messages in the system
type MessageType string

const (
	// Transaction ingress topic messages
	MessageTypeTxSubmission MessageType = "tx_submission"

	// Block execution topic messages
	MessageTypeBlockExecution MessageType = "block_execution"

	// State commit topic messages
	MessageTypeStateCommit MessageType = "state_commit"

	// Cache invalidation topic messages
	MessageTypeCacheInvalidation MessageType = "cache_invalidation"
)

// TxSubmissionMessage represents a transaction submission message
type TxSubmissionMessage struct {
	Type        MessageType  `json:"type"`
	Transaction *Transaction `json:"transaction"`
	SubmittedAt time.Time    `json:"submittedAt"`
	PortalNode  string       `json:"portalNode"`
	MessageID   string       `json:"messageId"`
}

// BlockExecutionMessage represents a block execution message
type BlockExecutionMessage struct {
	Type      MessageType `json:"type"`
	Block     *Block      `json:"block"`
	CreatedAt time.Time   `json:"createdAt"`
	Creator   string      `json:"creator"`
	MessageID string      `json:"messageId"`
}

// StateCommitMessage represents a state commit message
type StateCommitMessage struct {
	Type            MessageType      `json:"type"`
	ExecutionResult *ExecutionResult `json:"executionResult"`
	CommitAt        time.Time        `json:"commitAt"`
	MessageID       string           `json:"messageId"`
}

// Topic names used in RocketMQ
const (
	// TopicTxIngress is the topic for incoming transactions
	TopicTxIngress = "tx-ingress-topic"

	// TopicBlockExecute is the topic for block execution
	TopicBlockExecute = "block-execute-topic"

	// TopicStateCommit is the topic for state commit
	TopicStateCommit = "state-commit-topic"

	// TopicCacheInvalidation is the topic for cache invalidation (broadcast mode)
	TopicCacheInvalidation = "cache-invalidation-topic"
)

// Consumer group names
const (
	// ConsumerGroupChronos is the consumer group for YaoChronos
	ConsumerGroupChronos = "yao-chronos-group"

	// ConsumerGroupVulcan is the consumer group for YaoVulcan
	ConsumerGroupVulcan = "yao-vulcan-group"

	// ConsumerGroupArchive is the consumer group for YaoArchive
	ConsumerGroupArchive = "yao-archive-group"
)

// MessageHandler defines the interface for handling different types of messages
type MessageHandler interface {
	// HandleTxSubmission handles transaction submission messages
	HandleTxSubmission(msg *TxSubmissionMessage) error

	// HandleBlockExecution handles block execution messages
	HandleBlockExecution(msg *BlockExecutionMessage) error

	// HandleStateCommit handles state commit messages
	HandleStateCommit(msg *StateCommitMessage) error

	// HandleCacheInvalidation handles cache invalidation messages
	HandleCacheInvalidation(msg *CacheInvalidationMessage) error
}

// MessageProducer defines the interface for message producers
type MessageProducer interface {
	// SendTxSubmission sends a transaction submission message
	SendTxSubmission(msg *TxSubmissionMessage) error

	// SendBlockExecution sends a block execution message
	SendBlockExecution(msg *BlockExecutionMessage) error

	// SendStateCommit sends a state commit message
	SendStateCommit(msg *StateCommitMessage) error

	// SendCacheInvalidation broadcasts a cache invalidation message
	SendCacheInvalidation(msg *CacheInvalidationMessage) error

	// Close closes the producer
	Close() error
}

// MessageConsumer defines the interface for message consumers
type MessageConsumer interface {
	// Subscribe subscribes to a topic with a handler
	Subscribe(topic string, consumerGroup string, handler MessageHandler) error

	// SubscribeBroadcast subscribes to a broadcast topic with a handler
	SubscribeBroadcast(topic string, handler MessageHandler) error

	// Start starts consuming messages
	Start() error

	// Stop stops consuming messages
	Stop() error

	// Close closes the consumer
	Close() error
}
