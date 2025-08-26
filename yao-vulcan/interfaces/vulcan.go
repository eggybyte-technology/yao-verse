package interfaces

import (
	"context"
	"math/big"
	"time"

	oracleinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

// EVMExecutor defines the interface for EVM execution
type EVMExecutor interface {
	// ExecuteBlock executes all transactions in a block
	ExecuteBlock(ctx context.Context, block *types.Block) (*types.ExecutionResult, error)

	// ExecuteTransaction executes a single transaction
	ExecuteTransaction(ctx context.Context, tx *types.Transaction, blockNumber uint64, blockHash common.Hash) (*types.Receipt, error)

	// ExecuteCall executes a call without state changes
	ExecuteCall(ctx context.Context, args CallArgs, blockNumber *big.Int) ([]byte, error)

	// EstimateGas estimates gas usage for a transaction
	EstimateGas(ctx context.Context, args CallArgs, blockNumber *big.Int) (uint64, error)

	// CreateEVM creates a new EVM instance
	CreateEVM(ctx context.Context, blockNumber uint64, blockHash common.Hash, timestamp uint64) (*vm.EVM, error)

	// GetVMConfig returns the EVM configuration
	GetVMConfig() *vm.Config

	// SetVMConfig sets the EVM configuration
	SetVMConfig(config *vm.Config)

	// GetExecutionMetrics returns execution metrics
	GetExecutionMetrics() *ExecutionMetrics
}

// CallArgs represents the arguments for EVM calls
type CallArgs struct {
	From      *common.Address `json:"from"`
	To        *common.Address `json:"to"`
	Gas       *uint64         `json:"gas"`
	GasPrice  *big.Int        `json:"gasPrice"`
	Value     *big.Int        `json:"value"`
	Data      []byte          `json:"data"`
	GasFeeCap *big.Int        `json:"maxFeePerGas"`
	GasTipCap *big.Int        `json:"maxPriorityFeePerGas"`
}

// BlockProcessor defines the interface for processing blocks
type BlockProcessor interface {
	// ProcessBlock processes a block and returns execution results
	ProcessBlock(ctx context.Context, block *types.Block) (*types.ExecutionResult, error)

	// ValidateBlock validates a block before processing
	ValidateBlock(ctx context.Context, block *types.Block) error

	// GetProcessingStatus returns the current processing status
	GetProcessingStatus() *ProcessingStatus

	// GetLastProcessedBlock returns the last processed block number
	GetLastProcessedBlock() uint64

	// SetBlockProcessedCallback sets callback for when a block is processed
	SetBlockProcessedCallback(callback BlockProcessedCallback) error
}

// BlockProcessedCallback defines the callback for when a block is processed
type BlockProcessedCallback func(result *types.ExecutionResult)

// ProcessingStatus represents the current processing status
type ProcessingStatus struct {
	CurrentBlock   uint64        `json:"currentBlock"`
	ProcessingTime time.Duration `json:"processingTime"`
	Status         ProcessStatus `json:"status"`
	ErrorMessage   string        `json:"errorMessage,omitempty"`
	LastUpdate     time.Time     `json:"lastUpdate"`
}

// ProcessStatus defines the processing status
type ProcessStatus string

const (
	ProcessStatusIdle       ProcessStatus = "idle"
	ProcessStatusProcessing ProcessStatus = "processing"
	ProcessStatusError      ProcessStatus = "error"
	ProcessStatusCompleted  ProcessStatus = "completed"
)

// StateManager defines the interface for managing blockchain state
type StateManager interface {
	// GetStateDB returns the state database for a specific block
	GetStateDB(ctx context.Context, blockNumber uint64) (oracleinterfaces.StateDB, error)

	// CreateSnapshot creates a state snapshot
	CreateSnapshot(ctx context.Context) (int, error)

	// RevertToSnapshot reverts state to a specific snapshot
	RevertToSnapshot(ctx context.Context, snapshot int) error

	// ApplyStateDiff applies state changes from execution
	ApplyStateDiff(ctx context.Context, stateDiff *types.StateDiff) error

	// GetStateDiff calculates state difference between two points
	GetStateDiff(ctx context.Context, fromSnapshot, toSnapshot int) (*types.StateDiff, error)

	// CommitState commits state changes to storage
	CommitState(ctx context.Context, blockNumber uint64) (common.Hash, error)

	// GetStateRoot returns the state root for a block
	GetStateRoot(ctx context.Context, blockNumber uint64) (common.Hash, error)
}

// GasCalculator defines the interface for gas calculations
type GasCalculator interface {
	// CalculateGasPrice calculates the gas price for transactions
	CalculateGasPrice(ctx context.Context) (*big.Int, error)

	// EstimateExecutionGas estimates gas for transaction execution
	EstimateExecutionGas(ctx context.Context, tx *types.Transaction) (uint64, error)

	// GetGasLimit returns the current gas limit
	GetGasLimit() uint64

	// SetGasLimit sets the gas limit
	SetGasLimit(gasLimit uint64)

	// GetBaseFee returns the base fee for EIP-1559
	GetBaseFee(ctx context.Context, blockNumber uint64) (*big.Int, error)
}

// VulcanServer defines the main YaoVulcan server interface
type VulcanServer interface {
	// Start starts the vulcan server
	Start(ctx context.Context) error

	// Stop stops the vulcan server gracefully
	Stop(ctx context.Context) error

	// GetEVMExecutor returns the EVM executor instance
	GetEVMExecutor() EVMExecutor

	// GetBlockProcessor returns the block processor instance
	GetBlockProcessor() BlockProcessor

	// GetStateManager returns the state manager instance
	GetStateManager() StateManager

	// GetGasCalculator returns the gas calculator instance
	GetGasCalculator() GasCalculator

	// GetOracle returns the embedded YaoOracle instance
	GetOracle() oracleinterfaces.YaoOracle

	// GetMetrics returns server metrics
	GetMetrics() *VulcanMetrics

	// HealthCheck performs a comprehensive health check
	HealthCheck(ctx context.Context) error

	// GetStatus returns the current server status
	GetStatus() *VulcanStatus
}

// ExecutionMetrics defines metrics for EVM execution
type ExecutionMetrics struct {
	// Transaction metrics
	TransactionsExecuted  uint64  `json:"transactionsExecuted"`
	TransactionsPerSecond float64 `json:"transactionsPerSecond"`
	FailedTransactions    uint64  `json:"failedTransactions"`

	// Gas metrics
	TotalGasUsed      uint64  `json:"totalGasUsed"`
	AverageGasUsed    uint64  `json:"averageGasUsed"`
	GasUsagePerSecond float64 `json:"gasUsagePerSecond"`

	// Execution time metrics
	TotalExecutionTime   int64 `json:"totalExecutionTime"`   // in nanoseconds
	AverageExecutionTime int64 `json:"averageExecutionTime"` // in nanoseconds
	MaxExecutionTime     int64 `json:"maxExecutionTime"`     // in nanoseconds
	MinExecutionTime     int64 `json:"minExecutionTime"`     // in nanoseconds

	// Contract metrics
	ContractCreations uint64 `json:"contractCreations"`
	ContractCalls     uint64 `json:"contractCalls"`
	ContractErrors    uint64 `json:"contractErrors"`

	// State access metrics
	StateReads   uint64  `json:"stateReads"`
	StateWrites  uint64  `json:"stateWrites"`
	CacheHitRate float64 `json:"cacheHitRate"`
}

// VulcanMetrics defines comprehensive metrics for the vulcan server
type VulcanMetrics struct {
	// Execution metrics
	ExecutionMetrics *ExecutionMetrics `json:"executionMetrics"`

	// Block processing metrics
	BlocksProcessed  uint64  `json:"blocksProcessed"`
	BlocksPerSecond  float64 `json:"blocksPerSecond"`
	ProcessingErrors uint64  `json:"processingErrors"`
	AverageBlockTime int64   `json:"averageBlockTime"` // in milliseconds

	// Oracle metrics (embedded YaoOracle)
	OracleMetrics *oracleinterfaces.OracleMetrics `json:"oracleMetrics"`

	// System metrics
	MemoryUsage    uint64  `json:"memoryUsage"`
	CPUUsage       float64 `json:"cpuUsage"`
	GoroutineCount int     `json:"goroutineCount"`
	HeapSize       uint64  `json:"heapSize"`

	// Message queue metrics
	MessagesConsumed uint64 `json:"messagesConsumed"`
	MessagesProduced uint64 `json:"messagesProduced"`
	MessageLatency   int64  `json:"messageLatency"` // in milliseconds
	ConsumerLag      uint64 `json:"consumerLag"`

	// Network metrics
	NetworkConnections uint64 `json:"networkConnections"`
	BytesReceived      uint64 `json:"bytesReceived"`
	BytesSent          uint64 `json:"bytesSent"`

	// Error metrics
	TotalErrors     uint64 `json:"totalErrors"`
	EVMErrors       uint64 `json:"evmErrors"`
	StateErrors     uint64 `json:"stateErrors"`
	ConsensusErrors uint64 `json:"consensusErrors"`
}

// VulcanStatus represents the current status of the vulcan server
type VulcanStatus struct {
	NodeID              string        `json:"nodeId"`
	Status              ServerStatus  `json:"status"`
	LastProcessedBlock  uint64        `json:"lastProcessedBlock"`
	CurrentProcessing   *uint64       `json:"currentProcessing,omitempty"`
	ProcessingStartTime *time.Time    `json:"processingStartTime,omitempty"`
	StartTime           time.Time     `json:"startTime"`
	Uptime              time.Duration `json:"uptime"`
	ConfigHash          string        `json:"configHash"`
	Version             string        `json:"version"`

	// Oracle status
	OracleStatus string `json:"oracleStatus"`
	CacheStatus  string `json:"cacheStatus"`

	// EVM status
	EVMVersion  string `json:"evmVersion"`
	ChainConfig string `json:"chainConfig"`
}

// ServerStatus defines the server status (reusing from chronos)
type ServerStatus string

const (
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusRunning  ServerStatus = "running"
	ServerStatusStopping ServerStatus = "stopping"
	ServerStatusStopped  ServerStatus = "stopped"
	ServerStatusError    ServerStatus = "error"
)

// VulcanConfig defines the configuration for YaoVulcan
type VulcanConfig struct {
	// Node configuration
	NodeID   string `json:"nodeId"`
	NodeName string `json:"nodeName"`

	// EVM configuration
	ChainID         int64  `json:"chainId"`
	NetworkID       int64  `json:"networkId"`
	EVMVersion      string `json:"evmVersion"`
	EnableTracing   bool   `json:"enableTracing"`
	EnableProfiling bool   `json:"enableProfiling"`
	JumpTable       string `json:"jumpTable"` // istanbul, berlin, london, etc.

	// Gas configuration
	GasLimit             uint64   `json:"gasLimit"`
	GasPrice             *big.Int `json:"gasPrice"`
	BaseFeeChangeDenom   uint64   `json:"baseFeeChangeDenom"`   // EIP-1559
	ElasticityMultiplier uint64   `json:"elasticityMultiplier"` // EIP-1559

	// Processing configuration
	WorkerCount         int `json:"workerCount"`
	MaxConcurrentBlocks int `json:"maxConcurrentBlocks"`
	ProcessingTimeout   int `json:"processingTimeout"`  // in milliseconds
	StateFlushInterval  int `json:"stateFlushInterval"` // in milliseconds

	// Message queue configuration
	RocketMQEndpoints []string `json:"rocketmqEndpoints"`
	RocketMQAccessKey string   `json:"rocketmqAccessKey"`
	RocketMQSecretKey string   `json:"rocketmqSecretKey"`
	ConsumerGroup     string   `json:"consumerGroup"`
	ProducerGroup     string   `json:"producerGroup"`
	BatchSize         int      `json:"batchSize"`
	BatchTimeout      int      `json:"batchTimeout"` // in milliseconds

	// YaoOracle configuration (embedded state library)
	OracleConfig *oracleinterfaces.OracleConfig `json:"oracleConfig"`

	// Performance tuning
	EnableJIT           bool   `json:"enableJit"`           // Enable JIT compilation
	EnableOptimizations bool   `json:"enableOptimizations"` // Enable various optimizations
	MaxTransactionGas   uint64 `json:"maxTransactionGas"`   // Maximum gas per transaction
	MaxCallDepth        int    `json:"maxCallDepth"`        // Maximum call depth
	EnableMemoryLimit   bool   `json:"enableMemoryLimit"`   // Enable memory limit enforcement
	MemoryLimit         int    `json:"memoryLimit"`         // Memory limit in MB

	// Metrics configuration
	EnableMetrics         bool `json:"enableMetrics"`
	MetricsInterval       int  `json:"metricsInterval"` // in seconds
	MetricsPort           int  `json:"metricsPort"`
	EnableDetailedMetrics bool `json:"enableDetailedMetrics"`

	// Debugging configuration
	EnableDebugAPI  bool `json:"enableDebugApi"`
	DebugTraceLimit int  `json:"debugTraceLimit"`
	EnableGCStats   bool `json:"enableGcStats"`
	ProfilerPort    int  `json:"profilerPort"`

	// Logging configuration
	LogLevel            string `json:"logLevel"`
	LogFormat           string `json:"logFormat"`
	LogFile             string `json:"logFile"`
	EnableStructuredLog bool   `json:"enableStructuredLog"`

	// Health check configuration
	HealthCheckInterval int `json:"healthCheckInterval"` // in seconds
	HealthCheckTimeout  int `json:"healthCheckTimeout"`  // in seconds

	// Cache invalidation configuration
	EnableCacheInvalidation    bool `json:"enableCacheInvalidation"`
	CacheInvalidationBatchSize int  `json:"cacheInvalidationBatchSize"`
	CacheInvalidationTimeout   int  `json:"cacheInvalidationTimeout"` // in milliseconds
}
