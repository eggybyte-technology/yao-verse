package interfaces

import (
	"context"
	"math/big"
	"time"

	oracleinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
	sharedinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
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

	// GetStateDBForHash returns the state database for a specific block hash
	GetStateDBForHash(ctx context.Context, blockHash common.Hash) (oracleinterfaces.StateDB, error)

	// CreateSnapshot creates a state snapshot
	CreateSnapshot(ctx context.Context) (int, error)

	// RevertToSnapshot reverts state to a snapshot
	RevertToSnapshot(ctx context.Context, snapshotID int) error

	// CalculateStateDiff calculates state difference between two snapshots
	CalculateStateDiff(ctx context.Context, beforeSnapshot, afterSnapshot int) (*types.StateDiff, error)

	// CommitState commits state changes
	CommitState(ctx context.Context, stateDiff *types.StateDiff) error

	// GetMetrics returns state manager metrics
	GetMetrics() *StateManagerMetrics

	// HealthCheck performs state manager health check
	HealthCheck(ctx context.Context) error
}

// GasCalculator defines the interface for gas calculations
type GasCalculator interface {
	// EstimateGas estimates gas required for a transaction
	EstimateGas(ctx context.Context, tx *types.Transaction, blockNumber uint64) (uint64, error)

	// CalculateGasPrice calculates optimal gas price
	CalculateGasPrice(ctx context.Context) (*big.Int, error)

	// ValidateGasLimit validates if gas limit is sufficient
	ValidateGasLimit(ctx context.Context, tx *types.Transaction) error

	// GetBaseFee returns the base fee for EIP-1559
	GetBaseFee(ctx context.Context, blockNumber uint64) (*big.Int, error)

	// GetGasLimit returns the current gas limit
	GetGasLimit() uint64

	// SetGasLimit sets the gas limit
	SetGasLimit(gasLimit uint64)

	// GetGasMetrics returns gas calculation metrics
	GetGasMetrics() *GasMetrics
}

// ParallelExecutionScheduler defines the interface for parallel block execution
type ParallelExecutionScheduler interface {
	// ScheduleBlockExecution schedules a block for parallel execution
	ScheduleBlockExecution(ctx context.Context, block *types.Block, priority ExecutionPriority) (*ExecutionJob, error)

	// GetExecutionQueue returns the current execution queue status
	GetExecutionQueue(ctx context.Context) ([]*ExecutionJob, error)

	// CancelExecution cancels a scheduled execution
	CancelExecution(ctx context.Context, jobID string) error

	// GetWorkerStatus returns the status of all execution workers
	GetWorkerStatus(ctx context.Context) ([]*WorkerStatus, error)

	// SetWorkerCount sets the number of execution workers
	SetWorkerCount(ctx context.Context, count int) error

	// GetSchedulerMetrics returns scheduler metrics
	GetSchedulerMetrics() *SchedulerMetrics
}

// CacheInvalidationHandler defines the interface for handling cache invalidation
type CacheInvalidationHandler interface {
	// HandleInvalidationMessage handles a cache invalidation message
	HandleInvalidationMessage(ctx context.Context, message *types.CacheInvalidationMessage) error

	// InvalidateKeys invalidates specific cache keys
	InvalidateKeys(ctx context.Context, keys []string) error

	// GetInvalidationMetrics returns cache invalidation metrics
	GetInvalidationMetrics() *CacheInvalidationMetrics

	// RegisterInvalidationCallback registers a callback for cache invalidation
	RegisterInvalidationCallback(callback CacheInvalidationCallback) error
}

// TransactionTracer defines the interface for transaction tracing and debugging
type TransactionTracer interface {
	// TraceTransaction traces a transaction execution
	TraceTransaction(ctx context.Context, tx *types.Transaction, config *TraceConfig) (*TransactionTrace, error)

	// TraceBlock traces all transactions in a block
	TraceBlock(ctx context.Context, block *types.Block, config *TraceConfig) ([]*TransactionTrace, error)

	// TraceCall traces a call execution
	TraceCall(ctx context.Context, args CallArgs, config *TraceConfig) (*CallTrace, error)

	// GetTraceHistory returns trace history for analysis
	GetTraceHistory(ctx context.Context, filter *TraceFilter) ([]*TransactionTrace, error)

	// GetTracingMetrics returns tracing metrics
	GetTracingMetrics() *TracingMetrics
}

// ContractManager defines the interface for contract management
type ContractManager interface {
	// DeployContract deploys a new contract
	DeployContract(ctx context.Context, deployment *ContractDeployment) (*ContractInfo, error)

	// GetContract returns contract information
	GetContract(ctx context.Context, address common.Address) (*ContractInfo, error)

	// GetContractCode returns contract bytecode
	GetContractCode(ctx context.Context, address common.Address) ([]byte, error)

	// ValidateContract validates contract bytecode
	ValidateContract(ctx context.Context, bytecode []byte) (*ContractValidation, error)

	// GetContractMetrics returns contract-related metrics
	GetContractMetrics() *ContractMetrics
}

// VulcanServer defines the main YaoVulcan server interface
type VulcanServer interface {
	// Start starts the Vulcan server
	Start(ctx context.Context) error

	// Stop stops the Vulcan server gracefully
	Stop(ctx context.Context) error

	// GetEVMExecutor returns the EVM executor instance
	GetEVMExecutor() EVMExecutor

	// GetBlockProcessor returns the block processor instance
	GetBlockProcessor() BlockProcessor

	// GetStateManager returns the state manager instance
	GetStateManager() StateManager

	// GetGasCalculator returns the gas calculator instance
	GetGasCalculator() GasCalculator

	// GetParallelScheduler returns the parallel execution scheduler
	GetParallelScheduler() ParallelExecutionScheduler

	// GetCacheInvalidationHandler returns the cache invalidation handler
	GetCacheInvalidationHandler() CacheInvalidationHandler

	// GetTransactionTracer returns the transaction tracer
	GetTransactionTracer() TransactionTracer

	// GetContractManager returns the contract manager
	GetContractManager() ContractManager

	// GetYaoOracle returns the embedded YaoOracle instance
	GetYaoOracle() oracleinterfaces.YaoOracle

	// GetMessageProducer returns the message producer for publishing results
	GetMessageProducer() sharedinterfaces.MessageProducer

	// GetMessageConsumer returns the message consumer for receiving blocks
	GetMessageConsumer() sharedinterfaces.MessageConsumer

	// PauseExecution pauses block execution
	PauseExecution(ctx context.Context) error

	// ResumeExecution resumes block execution
	ResumeExecution(ctx context.Context) error

	// GetStatus returns the current server status
	GetStatus(ctx context.Context) *VulcanStatus

	// GetMetrics returns comprehensive server metrics
	GetMetrics(ctx context.Context) *VulcanMetrics

	// HealthCheck performs comprehensive health check
	HealthCheck(ctx context.Context) error
}

// Supporting types for new and existing interfaces

// ExecutionPriority defines execution priority levels
type ExecutionPriority string

const (
	ExecutionPriorityLow      ExecutionPriority = "low"
	ExecutionPriorityNormal   ExecutionPriority = "normal"
	ExecutionPriorityHigh     ExecutionPriority = "high"
	ExecutionPriorityCritical ExecutionPriority = "critical"
)

// ExecutionJob represents a block execution job
type ExecutionJob struct {
	JobID          string                 `json:"jobId"`
	Block          *types.Block           `json:"block"`
	Priority       ExecutionPriority      `json:"priority"`
	Status         JobStatus              `json:"status"`
	AssignedWorker string                 `json:"assignedWorker,omitempty"`
	ScheduledAt    time.Time              `json:"scheduledAt"`
	StartedAt      *time.Time             `json:"startedAt,omitempty"`
	CompletedAt    *time.Time             `json:"completedAt,omitempty"`
	Result         *types.ExecutionResult `json:"result,omitempty"`
	ErrorMessage   string                 `json:"errorMessage,omitempty"`
	Metadata       map[string]string      `json:"metadata,omitempty"`
}

// JobStatus defines execution job status
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// WorkerStatus represents the status of an execution worker
type WorkerStatus struct {
	WorkerID       string        `json:"workerId"`
	Status         WorkerState   `json:"status"`
	CurrentJob     *string       `json:"currentJob,omitempty"`
	JobsCompleted  uint64        `json:"jobsCompleted"`
	JobsFailed     uint64        `json:"jobsFailed"`
	AverageJobTime time.Duration `json:"averageJobTime"`
	LastJobTime    *time.Time    `json:"lastJobTime,omitempty"`
	CPUUsage       float64       `json:"cpuUsage"`
	MemoryUsage    int64         `json:"memoryUsage"`
	ErrorMessage   string        `json:"errorMessage,omitempty"`
}

// WorkerState defines execution worker states
type WorkerState string

const (
	WorkerStateIdle    WorkerState = "idle"
	WorkerStateRunning WorkerState = "running"
	WorkerStateError   WorkerState = "error"
	WorkerStateStopped WorkerState = "stopped"
)

// SchedulerMetrics provides parallel execution scheduler metrics
type SchedulerMetrics struct {
	TotalJobs           uint64        `json:"totalJobs"`
	CompletedJobs       uint64        `json:"completedJobs"`
	FailedJobs          uint64        `json:"failedJobs"`
	CancelledJobs       uint64        `json:"cancelledJobs"`
	QueuedJobs          int           `json:"queuedJobs"`
	ActiveWorkers       int           `json:"activeWorkers"`
	IdleWorkers         int           `json:"idleWorkers"`
	AverageJobTime      time.Duration `json:"averageJobTime"`
	ThroughputPerSecond float64       `json:"throughputPerSecond"`
	QueueWaitTime       time.Duration `json:"queueWaitTime"`
	WorkerUtilization   float64       `json:"workerUtilization"` // percentage
	Timestamp           time.Time     `json:"timestamp"`
}

// CacheInvalidationCallback defines callback for cache invalidation events
type CacheInvalidationCallback func(keys []string)

// CacheInvalidationMetrics provides cache invalidation metrics
type CacheInvalidationMetrics struct {
	TotalInvalidations      uint64        `json:"totalInvalidations"`
	KeysInvalidated         uint64        `json:"keysInvalidated"`
	InvalidationErrors      uint64        `json:"invalidationErrors"`
	AverageInvalidationTime time.Duration `json:"averageInvalidationTime"`
	InvalidationRate        float64       `json:"invalidationRate"` // keys per second
	LastInvalidation        time.Time     `json:"lastInvalidation"`
	Timestamp               time.Time     `json:"timestamp"`
}

// TraceConfig defines configuration for transaction tracing
type TraceConfig struct {
	EnableMemory     bool                   `json:"enableMemory"`
	EnableStack      bool                   `json:"enableStack"`
	EnableStorage    bool                   `json:"enableStorage"`
	EnableReturnData bool                   `json:"enableReturnData"`
	Timeout          time.Duration          `json:"timeout"`
	TraceType        TraceType              `json:"traceType"`
	CustomConfig     map[string]interface{} `json:"customConfig,omitempty"`
}

// TraceType defines types of tracing
type TraceType string

const (
	TraceTypeCall      TraceType = "call"
	TraceTypeDebug     TraceType = "debug"
	TraceTypeVMTrace   TraceType = "vmTrace"
	TraceTypeStateDiff TraceType = "stateDiff"
	TraceTypePrestate  TraceType = "prestate"
)

// TransactionTrace represents a transaction execution trace
type TransactionTrace struct {
	TxHash      common.Hash      `json:"txHash"`
	BlockNumber uint64           `json:"blockNumber"`
	BlockHash   common.Hash      `json:"blockHash"`
	TxIndex     uint64           `json:"txIndex"`
	TraceType   TraceType        `json:"traceType"`
	From        common.Address   `json:"from"`
	To          common.Address   `json:"to"`
	Input       []byte           `json:"input"`
	Output      []byte           `json:"output"`
	Value       *big.Int         `json:"value"`
	Gas         uint64           `json:"gas"`
	GasUsed     uint64           `json:"gasUsed"`
	Error       string           `json:"error,omitempty"`
	Calls       []*CallTrace     `json:"calls,omitempty"`
	Logs        []*types.Log     `json:"logs,omitempty"`
	StateDiff   *types.StateDiff `json:"stateDiff,omitempty"`
	VMTrace     *VMTrace         `json:"vmTrace,omitempty"`
	Duration    time.Duration    `json:"duration"`
	GeneratedAt time.Time        `json:"generatedAt"`
}

// CallTrace represents a call trace within a transaction
type CallTrace struct {
	Type    string         `json:"type"`
	From    common.Address `json:"from"`
	To      common.Address `json:"to"`
	Input   []byte         `json:"input"`
	Output  []byte         `json:"output"`
	Value   *big.Int       `json:"value"`
	Gas     uint64         `json:"gas"`
	GasUsed uint64         `json:"gasUsed"`
	Error   string         `json:"error,omitempty"`
	Calls   []*CallTrace   `json:"calls,omitempty"`
	Depth   int            `json:"depth"`
}

// VMTrace represents detailed VM execution trace
type VMTrace struct {
	Code         []byte            `json:"code"`
	Operations   []*VMOperation    `json:"operations"`
	MemoryTrace  [][]byte          `json:"memoryTrace,omitempty"`
	StackTrace   [][]string        `json:"stackTrace,omitempty"`
	StorageTrace map[string]string `json:"storageTrace,omitempty"`
}

// VMOperation represents a single VM operation
type VMOperation struct {
	PC      uint64            `json:"pc"`
	Op      string            `json:"op"`
	Gas     uint64            `json:"gas"`
	GasCost uint64            `json:"gasCost"`
	Memory  []byte            `json:"memory,omitempty"`
	Stack   []string          `json:"stack,omitempty"`
	Storage map[string]string `json:"storage,omitempty"`
	Depth   int               `json:"depth"`
}

// TraceFilter defines filter criteria for trace queries
type TraceFilter struct {
	FromBlock   *uint64          `json:"fromBlock,omitempty"`
	ToBlock     *uint64          `json:"toBlock,omitempty"`
	FromAddress []common.Address `json:"fromAddress,omitempty"`
	ToAddress   []common.Address `json:"toAddress,omitempty"`
	TraceType   *TraceType       `json:"traceType,omitempty"`
	Limit       int              `json:"limit,omitempty"`
}

// TracingMetrics provides transaction tracing metrics
type TracingMetrics struct {
	TotalTraces      uint64               `json:"totalTraces"`
	TracesPerSecond  float64              `json:"tracesPerSecond"`
	AverageTraceTime time.Duration        `json:"averageTraceTime"`
	TracesByType     map[TraceType]uint64 `json:"tracesByType"`
	TracingErrors    uint64               `json:"tracingErrors"`
	StorageUsage     int64                `json:"storageUsage"`
	LastTrace        time.Time            `json:"lastTrace"`
	Timestamp        time.Time            `json:"timestamp"`
}

// ContractDeployment represents a contract deployment request
type ContractDeployment struct {
	Bytecode    []byte            `json:"bytecode"`
	Constructor []byte            `json:"constructor,omitempty"`
	Value       *big.Int          `json:"value"`
	Gas         uint64            `json:"gas"`
	GasPrice    *big.Int          `json:"gasPrice"`
	From        common.Address    `json:"from"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ContractInfo represents information about a deployed contract
type ContractInfo struct {
	Address     common.Address    `json:"address"`
	Creator     common.Address    `json:"creator"`
	CreationTx  common.Hash       `json:"creationTx"`
	CreatedAt   time.Time         `json:"createdAt"`
	BlockNumber uint64            `json:"blockNumber"`
	Bytecode    []byte            `json:"bytecode"`
	CodeHash    common.Hash       `json:"codeHash"`
	Balance     *big.Int          `json:"balance"`
	IsVerified  bool              `json:"isVerified"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ContractValidation represents contract validation results
type ContractValidation struct {
	IsValid                 bool              `json:"isValid"`
	ValidationTime          time.Duration     `json:"validationTime"`
	Issues                  []ValidationIssue `json:"issues,omitempty"`
	GasEstimate             uint64            `json:"gasEstimate"`
	SecurityScore           float64           `json:"securityScore"` // 0-100
	OptimizationSuggestions []string          `json:"optimizationSuggestions,omitempty"`
}

// ValidationIssue represents a contract validation issue
type ValidationIssue struct {
	Severity    IssueSeverity `json:"severity"`
	Type        IssueType     `json:"type"`
	Description string        `json:"description"`
	Location    *CodeLocation `json:"location,omitempty"`
	Resolution  string        `json:"resolution,omitempty"`
}

// IssueSeverity defines validation issue severity levels
type IssueSeverity string

const (
	IssueSeverityInfo     IssueSeverity = "info"
	IssueSeverityWarning  IssueSeverity = "warning"
	IssueSeverityError    IssueSeverity = "error"
	IssueSeverityCritical IssueSeverity = "critical"
)

// IssueType defines types of validation issues
type IssueType string

const (
	IssueTypeSecurity     IssueType = "security"
	IssueTypeOptimization IssueType = "optimization"
	IssueTypeSyntax       IssueType = "syntax"
	IssueTypeLogic        IssueType = "logic"
)

// CodeLocation represents a location in contract code
type CodeLocation struct {
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	Line     int    `json:"line,omitempty"`
	Column   int    `json:"column,omitempty"`
	Function string `json:"function,omitempty"`
}

// ContractMetrics provides contract-related metrics
type ContractMetrics struct {
	TotalContracts       uint64            `json:"totalContracts"`
	ContractDeployments  uint64            `json:"contractDeployments"`
	ContractCalls        uint64            `json:"contractCalls"`
	ContractErrors       uint64            `json:"contractErrors"`
	VerifiedContracts    uint64            `json:"verifiedContracts"`
	AverageDeploymentGas uint64            `json:"averageDeploymentGas"`
	ContractsByType      map[string]uint64 `json:"contractsByType,omitempty"`
	Timestamp            time.Time         `json:"timestamp"`
}

// Status and metrics types

// VulcanStatus represents the status of YaoVulcan server
type VulcanStatus struct {
	ServerStatus           ServerStatus  `json:"serverStatus"`
	LastProcessedBlock     uint64        `json:"lastProcessedBlock"`
	CurrentProcessingBlock uint64        `json:"currentProcessingBlock,omitempty"`
	QueuedBlocks           int           `json:"queuedBlocks"`
	ProcessingRate         float64       `json:"processingRate"` // blocks per second
	ExecutorCount          int           `json:"executorCount"`
	OracleStatus           string        `json:"oracleStatus"`
	CacheHitRate           float64       `json:"cacheHitRate"`
	Uptime                 time.Duration `json:"uptime"`
	StartTime              time.Time     `json:"startTime"`
	Version                string        `json:"version"`
}

// ServerStatus defines server status
type ServerStatus string

const (
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusRunning  ServerStatus = "running"
	ServerStatusPaused   ServerStatus = "paused"
	ServerStatusStopping ServerStatus = "stopping"
	ServerStatusStopped  ServerStatus = "stopped"
	ServerStatusError    ServerStatus = "error"
)

// ExecutionMetrics defines metrics for EVM execution
type ExecutionMetrics struct {
	// Execution metrics
	TotalExecutions      uint64        `json:"totalExecutions"`
	ExecutionsPerSecond  float64       `json:"executionsPerSecond"`
	AverageExecutionTime time.Duration `json:"averageExecutionTime"`
	P99ExecutionTime     time.Duration `json:"p99ExecutionTime"`

	// Transaction metrics
	TransactionsExecuted  uint64  `json:"transactionsExecuted"`
	TransactionsPerSecond float64 `json:"transactionsPerSecond"`
	SuccessfulTxs         uint64  `json:"successfulTxs"`
	FailedTxs             uint64  `json:"failedTxs"`

	// Gas metrics
	TotalGasUsed     uint64  `json:"totalGasUsed"`
	AverageGasUsed   uint64  `json:"averageGasUsed"`
	GasUsedPerSecond float64 `json:"gasUsedPerSecond"`

	// Contract metrics
	ContractCalls     uint64 `json:"contractCalls"`
	ContractCreations uint64 `json:"contractCreations"`
	ContractErrors    uint64 `json:"contractErrors"`

	// Error metrics
	ExecutionErrors uint64 `json:"executionErrors"`
	VMErrors        uint64 `json:"vmErrors"`
	OutOfGasErrors  uint64 `json:"outOfGasErrors"`

	Timestamp time.Time `json:"timestamp"`
}

// StateManagerMetrics defines metrics for state management
type StateManagerMetrics struct {
	// State operations
	StateReads               uint64        `json:"stateReads"`
	StateWrites              uint64        `json:"stateWrites"`
	StateOperationsPerSecond float64       `json:"stateOperationsPerSecond"`
	AverageStateOpTime       time.Duration `json:"averageStateOpTime"`

	// Snapshot metrics
	SnapshotsCreated  uint64 `json:"snapshotsCreated"`
	SnapshotsReverted uint64 `json:"snapshotsReverted"`
	ActiveSnapshots   int    `json:"activeSnapshots"`

	// State diff metrics
	StateDiffsCalculated uint64 `json:"stateDiffsCalculated"`
	StateCommits         uint64 `json:"stateCommits"`
	StateRollbacks       uint64 `json:"stateRollbacks"`

	// Error metrics
	StateErrors    uint64 `json:"stateErrors"`
	CommitErrors   uint64 `json:"commitErrors"`
	SnapshotErrors uint64 `json:"snapshotErrors"`

	Timestamp time.Time `json:"timestamp"`
}

// GasMetrics defines metrics for gas calculations
type GasMetrics struct {
	// Gas estimation metrics
	GasEstimations        uint64        `json:"gasEstimations"`
	AverageEstimationTime time.Duration `json:"averageEstimationTime"`
	EstimationErrors      uint64        `json:"estimationErrors"`

	// Gas price metrics
	GasPriceCalculations uint64   `json:"gasPriceCalculations"`
	CurrentGasPrice      *big.Int `json:"currentGasPrice"`
	AverageGasPrice      *big.Int `json:"averageGasPrice"`

	// Gas validation metrics
	GasValidations        uint64 `json:"gasValidations"`
	GasValidationFailures uint64 `json:"gasValidationFailures"`

	Timestamp time.Time `json:"timestamp"`
}

// VulcanMetrics defines comprehensive metrics for YaoVulcan
type VulcanMetrics struct {
	ServerMetrics            *ServerMetrics                  `json:"serverMetrics"`
	ExecutionMetrics         *ExecutionMetrics               `json:"executionMetrics"`
	StateManagerMetrics      *StateManagerMetrics            `json:"stateManagerMetrics"`
	GasMetrics               *GasMetrics                     `json:"gasMetrics"`
	SchedulerMetrics         *SchedulerMetrics               `json:"schedulerMetrics"`
	CacheInvalidationMetrics *CacheInvalidationMetrics       `json:"cacheInvalidationMetrics"`
	TracingMetrics           *TracingMetrics                 `json:"tracingMetrics"`
	ContractMetrics          *ContractMetrics                `json:"contractMetrics"`
	OracleMetrics            *oracleinterfaces.OracleMetrics `json:"oracleMetrics"`
	MessageMetrics           *MessageMetrics                 `json:"messageMetrics"`
}

// ServerMetrics defines general server metrics
type ServerMetrics struct {
	Status          ServerStatus  `json:"status"`
	Uptime          time.Duration `json:"uptime"`
	CPUUsage        float64       `json:"cpuUsage"`
	MemoryUsage     int64         `json:"memoryUsage"`
	GoroutineCount  int           `json:"goroutineCount"`
	GCPauseTime     time.Duration `json:"gcPauseTime"`
	RequestsHandled uint64        `json:"requestsHandled"`
	ErrorsCount     uint64        `json:"errorsCount"`
	Timestamp       time.Time     `json:"timestamp"`
}

// MessageMetrics defines messaging metrics
type MessageMetrics struct {
	MessagesProduced  uint64        `json:"messagesProduced"`
	MessagesConsumed  uint64        `json:"messagesConsumed"`
	ProducerErrors    uint64        `json:"producerErrors"`
	ConsumerErrors    uint64        `json:"consumerErrors"`
	AverageLatency    time.Duration `json:"averageLatency"`
	MessageQueueDepth int           `json:"messageQueueDepth"`
	ConnectionsActive int           `json:"connectionsActive"`
	ReconnectionCount uint64        `json:"reconnectionCount"`
	Timestamp         time.Time     `json:"timestamp"`
}

// VulcanConfig defines comprehensive configuration for YaoVulcan
type VulcanConfig struct {
	// Node configuration
	NodeID   string `json:"nodeId"`
	NodeName string `json:"nodeName"`

	// Server configuration
	BindAddress string `json:"bindAddress"`
	Port        int    `json:"port"`
	EnableTLS   bool   `json:"enableTls"`
	TLSCert     string `json:"tlsCert,omitempty"`
	TLSKey      string `json:"tlsKey,omitempty"`

	// EVM configuration
	EnableEVM         bool          `json:"enableEvm"`
	EVMChainID        uint64        `json:"evmChainId"`
	VMTimeout         time.Duration `json:"vmTimeout"`
	MaxCallStackDepth int           `json:"maxCallStackDepth"`
	EnableJIT         bool          `json:"enableJit"`
	EnableDebug       bool          `json:"enableDebug"`

	// Gas configuration
	GasLimit         uint64   `json:"gasLimit"`
	MaxGasPrice      *big.Int `json:"maxGasPrice,omitempty"`
	MinGasPrice      *big.Int `json:"minGasPrice,omitempty"`
	GasEstimationCap uint64   `json:"gasEstimationCap"`

	// Oracle configuration (embedded YaoOracle)
	OracleConfig oracleinterfaces.OracleConfig `json:"oracleConfig"`

	// Processing configuration
	WorkerCount         int           `json:"workerCount"`
	ProcessingTimeout   time.Duration `json:"processingTimeout"`
	MaxConcurrentBlocks int           `json:"maxConcurrentBlocks"`
	QueueCapacity       int           `json:"queueCapacity"`

	// Message queue configuration
	RocketMQEndpoints []string `json:"rocketmqEndpoints"`
	RocketMQAccessKey string   `json:"rocketmqAccessKey"`
	RocketMQSecretKey string   `json:"rocketmqSecretKey"`
	ConsumerGroup     string   `json:"consumerGroup"`
	ProducerGroup     string   `json:"producerGroup"`

	// State management configuration
	EnableSnapshots  bool `json:"enableSnapshots"`
	MaxSnapshots     int  `json:"maxSnapshots"`
	SnapshotInterval int  `json:"snapshotInterval"`
	EnableStateDiff  bool `json:"enableStateDiff"`

	// Performance tuning
	EnableOptimizations bool `json:"enableOptimizations"`
	CachePreloading     bool `json:"cachePreloading"`
	ParallelExecution   bool `json:"parallelExecution"`

	// Monitoring configuration
	EnableMetrics       bool          `json:"enableMetrics"`
	MetricsInterval     time.Duration `json:"metricsInterval"`
	EnableHealthCheck   bool          `json:"enableHealthCheck"`
	HealthCheckInterval time.Duration `json:"healthCheckInterval"`

	// Tracing configuration
	EnableTracing     bool          `json:"enableTracing"`
	TracingBufferSize int           `json:"tracingBufferSize"`
	TracingTimeout    time.Duration `json:"tracingTimeout"`

	// Cache invalidation configuration
	EnableCacheInvalidation    bool `json:"enableCacheInvalidation"`
	CacheInvalidationBatchSize int  `json:"cacheInvalidationBatchSize"`
	CacheInvalidationTimeout   int  `json:"cacheInvalidationTimeout"`

	// Logging configuration
	LogLevel    string `json:"logLevel"`
	LogFormat   string `json:"logFormat"`
	LogFile     string `json:"logFile"`
	LogEVMTrace bool   `json:"logEvmTrace"`
}
