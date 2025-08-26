package interfaces

import (
	"context"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
)

// Sequencer defines the interface for the transaction sequencer
type Sequencer interface {
	// Start starts the sequencer
	Start(ctx context.Context) error

	// Stop stops the sequencer gracefully
	Stop(ctx context.Context) error

	// ProcessTransactions processes a batch of transactions
	ProcessTransactions(ctx context.Context, txs []*types.Transaction) (*types.Block, error)

	// GetNextBlockNumber returns the next block number to be created
	GetNextBlockNumber() uint64

	// GetCurrentBlock returns the current block being processed
	GetCurrentBlock() *types.Block

	// GetBlockByNumber returns a block by its number
	GetBlockByNumber(ctx context.Context, number uint64) (*types.Block, error)

	// GetBlockByHash returns a block by its hash
	GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)

	// GetMetrics returns sequencer metrics
	GetMetrics() *SequencerMetrics

	// HealthCheck performs a health check
	HealthCheck(ctx context.Context) error
}

// BlockBuilder defines the interface for building blocks
type BlockBuilder interface {
	// NewBlock creates a new block with the given transactions
	NewBlock(parentHash common.Hash, number uint64, txs []*types.Transaction) (*types.Block, error)

	// ValidateBlock validates a block structure
	ValidateBlock(block *types.Block) error

	// CalculateBlockHash calculates the hash of a block
	CalculateBlockHash(block *types.Block) common.Hash

	// SetGasLimit sets the gas limit for new blocks
	SetGasLimit(gasLimit uint64)

	// GetGasLimit returns the current gas limit
	GetGasLimit() uint64
}

// TransactionOrdering defines the interface for transaction ordering strategies
type TransactionOrdering interface {
	// OrderTransactions orders a batch of transactions according to the strategy
	OrderTransactions(txs []*types.Transaction) []*types.Transaction

	// ValidateTransaction validates a transaction for ordering
	ValidateTransaction(tx *types.Transaction) error

	// GetOrderingStrategy returns the current ordering strategy
	GetOrderingStrategy() OrderingStrategy
}

// OrderingStrategy defines different transaction ordering strategies
type OrderingStrategy string

const (
	// OrderingStrategyFIFO orders transactions in first-in-first-out order
	OrderingStrategyFIFO OrderingStrategy = "fifo"

	// OrderingStrategyGasPrice orders transactions by gas price (highest first)
	OrderingStrategyGasPrice OrderingStrategy = "gas_price"

	// OrderingStrategyNonce orders transactions by account nonce
	OrderingStrategyNonce OrderingStrategy = "nonce"

	// OrderingStrategyHybrid combines multiple strategies
	OrderingStrategyHybrid OrderingStrategy = "hybrid"
)

// LeaderElector defines the interface for leader election
type LeaderElector interface {
	// Start starts the leader election process
	Start(ctx context.Context) error

	// Stop stops the leader election process
	Stop(ctx context.Context) error

	// IsLeader returns true if this node is the current leader
	IsLeader() bool

	// GetLeaderID returns the ID of the current leader
	GetLeaderID() string

	// GetLeaderEndpoint returns the endpoint of the current leader
	GetLeaderEndpoint() string

	// RegisterLeadershipCallback registers a callback for leadership changes
	RegisterLeadershipCallback(callback LeadershipCallback) error

	// GetElectionStatus returns the current election status
	GetElectionStatus() *ElectionStatus
}

// LeadershipCallback defines the callback for leadership changes
type LeadershipCallback func(isLeader bool, leaderID string)

// ElectionStatus represents the status of leader election
type ElectionStatus struct {
	IsLeader     bool      `json:"isLeader"`
	LeaderID     string    `json:"leaderId"`
	LeaderTerm   uint64    `json:"leaderTerm"`
	LastElection time.Time `json:"lastElection"`
	Participants []string  `json:"participants"`
}

// ChronosServer defines the main YaoChronos server interface
type ChronosServer interface {
	// Start starts the chronos server
	Start(ctx context.Context) error

	// Stop stops the chronos server gracefully
	Stop(ctx context.Context) error

	// GetSequencer returns the sequencer instance
	GetSequencer() Sequencer

	// GetBlockBuilder returns the block builder instance
	GetBlockBuilder() BlockBuilder

	// GetTransactionOrdering returns the transaction ordering instance
	GetTransactionOrdering() TransactionOrdering

	// GetLeaderElector returns the leader elector instance
	GetLeaderElector() LeaderElector

	// GetMetrics returns server metrics
	GetMetrics() *ChronosMetrics

	// HealthCheck performs a comprehensive health check
	HealthCheck(ctx context.Context) error

	// GetStatus returns the current server status
	GetStatus() *ChronosStatus
}

// SequencerMetrics defines metrics for the sequencer
type SequencerMetrics struct {
	// Block metrics
	BlocksCreated    uint64  `json:"blocksCreated"`
	BlocksPerSecond  float64 `json:"blocksPerSecond"`
	AverageBlockTime int64   `json:"averageBlockTime"` // in milliseconds
	CurrentBlockSize uint64  `json:"currentBlockSize"`
	AverageBlockSize uint64  `json:"averageBlockSize"`

	// Transaction metrics
	TransactionsProcessed uint64  `json:"transactionsProcessed"`
	TransactionsPerSecond float64 `json:"transactionsPerSecond"`
	PendingTransactions   uint64  `json:"pendingTransactions"`

	// Performance metrics
	ProcessingTime int64 `json:"processingTime"` // in nanoseconds
	QueueWaitTime  int64 `json:"queueWaitTime"`  // in nanoseconds

	// Error metrics
	TotalErrors    uint64 `json:"totalErrors"`
	InvalidTxs     uint64 `json:"invalidTxs"`
	OrderingErrors uint64 `json:"orderingErrors"`
}

// ChronosMetrics defines comprehensive metrics for the chronos server
type ChronosMetrics struct {
	// Sequencer metrics
	SequencerMetrics *SequencerMetrics `json:"sequencerMetrics"`

	// Leader election metrics
	ElectionCount   uint64 `json:"electionCount"`
	LeadershipTime  int64  `json:"leadershipTime"` // in seconds
	IsCurrentLeader bool   `json:"isCurrentLeader"`

	// System metrics
	MemoryUsage    uint64  `json:"memoryUsage"`
	CPUUsage       float64 `json:"cpuUsage"`
	GoroutineCount int     `json:"goroutineCount"`

	// Message queue metrics
	MessagesConsumed uint64 `json:"messagesConsumed"`
	MessagesProduced uint64 `json:"messagesProduced"`
	MessageLatency   int64  `json:"messageLatency"` // in milliseconds

	// Network metrics
	NetworkConnections uint64 `json:"networkConnections"`
	BytesReceived      uint64 `json:"bytesReceived"`
	BytesSent          uint64 `json:"bytesSent"`
}

// ChronosStatus represents the current status of the chronos server
type ChronosStatus struct {
	NodeID        string        `json:"nodeId"`
	IsLeader      bool          `json:"isLeader"`
	LeaderID      string        `json:"leaderId"`
	CurrentBlock  uint64        `json:"currentBlock"`
	LastBlockTime time.Time     `json:"lastBlockTime"`
	Status        ServerStatus  `json:"status"`
	StartTime     time.Time     `json:"startTime"`
	Uptime        time.Duration `json:"uptime"`
	ConfigHash    string        `json:"configHash"`
	Version       string        `json:"version"`
}

// ServerStatus defines the server status
type ServerStatus string

const (
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusRunning  ServerStatus = "running"
	ServerStatusStopping ServerStatus = "stopping"
	ServerStatusStopped  ServerStatus = "stopped"
	ServerStatusError    ServerStatus = "error"
)

// ChronosConfig defines the configuration for YaoChronos
type ChronosConfig struct {
	// Node configuration
	NodeID   string `json:"nodeId"`
	NodeName string `json:"nodeName"`

	// Leader election configuration
	EnableLeaderElection bool   `json:"enableLeaderElection"`
	LeaseNamespace       string `json:"leaseNamespace"`
	LeaseName            string `json:"leaseName"`
	LeaseDuration        int    `json:"leaseDuration"` // in seconds
	RenewDeadline        int    `json:"renewDeadline"` // in seconds
	RetryPeriod          int    `json:"retryPeriod"`   // in seconds

	// Block building configuration
	BlockGasLimit           uint64 `json:"blockGasLimit"`
	MaxBlockSize            int    `json:"maxBlockSize"` // in bytes
	BlockTime               int    `json:"blockTime"`    // target block time in milliseconds
	MaxTransactionsPerBlock int    `json:"maxTransactionsPerBlock"`

	// Transaction ordering configuration
	OrderingStrategy    OrderingStrategy `json:"orderingStrategy"`
	MaxPendingTxs       int              `json:"maxPendingTxs"`
	TxValidationTimeout int              `json:"txValidationTimeout"` // in milliseconds

	// Message queue configuration
	RocketMQEndpoints []string `json:"rocketmqEndpoints"`
	RocketMQAccessKey string   `json:"rocketmqAccessKey"`
	RocketMQSecretKey string   `json:"rocketmqSecretKey"`
	ConsumerGroup     string   `json:"consumerGroup"`
	ProducerGroup     string   `json:"producerGroup"`
	BatchSize         int      `json:"batchSize"`
	BatchTimeout      int      `json:"batchTimeout"` // in milliseconds

	// Performance tuning
	WorkerCount       int  `json:"workerCount"`
	ProcessingTimeout int  `json:"processingTimeout"` // in milliseconds
	EnableMetrics     bool `json:"enableMetrics"`
	MetricsInterval   int  `json:"metricsInterval"` // in seconds
	EnableProfiling   bool `json:"enableProfiling"`

	// Kubernetes configuration (for leader election)
	KubeConfig    string `json:"kubeConfig"`
	KubeNamespace string `json:"kubeNamespace"`

	// Logging configuration
	LogLevel  string `json:"logLevel"`
	LogFormat string `json:"logFormat"`
	LogFile   string `json:"logFile"`

	// Health check configuration
	HealthCheckInterval int `json:"healthCheckInterval"` // in seconds
	HealthCheckTimeout  int `json:"healthCheckTimeout"`  // in seconds
}
