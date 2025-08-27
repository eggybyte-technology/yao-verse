package interfaces

import (
	"context"
	"time"

	sharedinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
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
	OrderTransactions(txs []*types.Transaction) ([]*types.Transaction, error)

	// ValidateTransaction validates a transaction for ordering
	ValidateTransaction(tx *types.Transaction) error

	// GetStrategy returns the current ordering strategy
	GetStrategy() OrderingStrategy

	// SetStrategy sets the ordering strategy
	SetStrategy(strategy OrderingStrategy) error

	// GetMetrics returns ordering metrics
	GetMetrics() *OrderingMetrics
}

// OrderingStrategy defines transaction ordering strategies
type OrderingStrategy string

const (
	OrderingStrategyFIFO     OrderingStrategy = "fifo"
	OrderingStrategyGasPrice OrderingStrategy = "gas_price"
	OrderingStrategyNonce    OrderingStrategy = "nonce"
	OrderingStrategyHybrid   OrderingStrategy = "hybrid"
	OrderingStrategyWeighted OrderingStrategy = "weighted"
)

// LeaderElector defines the interface for leader election
type LeaderElector interface {
	// Start starts the leader election process
	Start(ctx context.Context) error

	// Stop stops the leader election process
	Stop(ctx context.Context) error

	// IsLeader returns true if this node is the current leader
	IsLeader() bool

	// GetLeader returns information about the current leader
	GetLeader(ctx context.Context) (*LeaderInfo, error)

	// StepDown causes the current leader to step down
	StepDown(ctx context.Context) error

	// RegisterCallback registers callback for leadership changes
	RegisterCallback(callback LeadershipCallback) error

	// GetElectionStatus returns the current election status
	GetElectionStatus() *ElectionStatus

	// GetMetrics returns election metrics
	GetMetrics() *ElectionMetrics

	// HealthCheck performs leader election health check
	HealthCheck(ctx context.Context) error
}

// LeadershipCallback defines the callback for leadership changes
type LeadershipCallback func(isLeader bool, leaderInfo *LeaderInfo)

// ElectionStatus represents the status of leader election
type ElectionStatus struct {
	IsLeader     bool      `json:"isLeader"`
	LeaderID     string    `json:"leaderId,omitempty"`
	Term         int64     `json:"term"`
	LastElection time.Time `json:"lastElection"`
	Participants []string  `json:"participants"`
	Status       string    `json:"status"`
}

// LeaderInfo represents information about the cluster leader
type LeaderInfo struct {
	NodeID      string    `json:"nodeId"`
	NodeName    string    `json:"nodeName"`
	Address     string    `json:"address"`
	ElectedAt   time.Time `json:"electedAt"`
	Term        int64     `json:"term"`
	LastContact time.Time `json:"lastContact"`
}

// ChronosServer defines the main YaoChronos server interface
type ChronosServer interface {
	// Start starts the Chronos server
	Start(ctx context.Context) error

	// Stop stops the Chronos server gracefully
	Stop(ctx context.Context) error

	// GetSequencer returns the sequencer instance
	GetSequencer() Sequencer

	// GetBlockBuilder returns the block builder instance
	GetBlockBuilder() BlockBuilder

	// GetTransactionOrdering returns the transaction ordering instance
	GetTransactionOrdering() TransactionOrdering

	// GetLeaderElector returns the leader elector instance
	GetLeaderElector() LeaderElector

	// GetMessageProducer returns the message producer for publishing blocks
	GetMessageProducer() sharedinterfaces.MessageProducer

	// GetMessageConsumer returns the message consumer for receiving transactions
	GetMessageConsumer() sharedinterfaces.MessageConsumer

	// PauseSequencing pauses transaction sequencing
	PauseSequencing(ctx context.Context) error

	// ResumeSequencing resumes transaction sequencing
	ResumeSequencing(ctx context.Context) error

	// GetStatus returns the current server status
	GetStatus(ctx context.Context) *ChronosStatus

	// GetMetrics returns comprehensive server metrics
	GetMetrics(ctx context.Context) *ChronosMetrics

	// HealthCheck performs comprehensive health check
	HealthCheck(ctx context.Context) error
}

// ChronosStatus represents the status of YaoChronos server
type ChronosStatus struct {
	ServerStatus    ServerStatus  `json:"serverStatus"`
	IsLeader        bool          `json:"isLeader"`
	CurrentTerm     int64         `json:"currentTerm"`
	LastBlockNumber uint64        `json:"lastBlockNumber"`
	QueuedTxCount   int           `json:"queuedTxCount"`
	ProcessingRate  float64       `json:"processingRate"` // blocks per second
	Uptime          time.Duration `json:"uptime"`
	StartTime       time.Time     `json:"startTime"`
	Version         string        `json:"version"`
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

// SequencerMetrics defines metrics for the sequencer
type SequencerMetrics struct {
	// Block metrics
	TotalBlocksCreated uint64        `json:"totalBlocksCreated"`
	BlocksPerSecond    float64       `json:"blocksPerSecond"`
	AverageBlockSize   int64         `json:"averageBlockSize"`
	AverageBlockTime   time.Duration `json:"averageBlockTime"`
	LastBlockNumber    uint64        `json:"lastBlockNumber"`

	// Transaction metrics
	TotalTransactions     uint64  `json:"totalTransactions"`
	TransactionsPerSecond float64 `json:"transactionsPerSecond"`
	QueuedTransactions    int     `json:"queuedTransactions"`
	DroppedTransactions   uint64  `json:"droppedTransactions"`

	// Performance metrics
	AverageProcessingTime time.Duration `json:"averageProcessingTime"`
	P99ProcessingTime     time.Duration `json:"p99ProcessingTime"`

	// Error metrics
	ProcessingErrors uint64 `json:"processingErrors"`
	ValidationErrors uint64 `json:"validationErrors"`

	Timestamp time.Time `json:"timestamp"`
}

// ElectionMetrics defines metrics for leader election
type ElectionMetrics struct {
	ElectionCount      uint64        `json:"electionCount"`
	LeadershipTime     time.Duration `json:"leadershipTime"`
	LastElectionTime   time.Time     `json:"lastElectionTime"`
	ElectionDuration   time.Duration `json:"electionDuration"`
	SplitBrainEvents   uint64        `json:"splitBrainEvents"`
	IsCurrentLeader    bool          `json:"isCurrentLeader"`
	CandidateCount     int           `json:"candidateCount"`
	HeartbeatsSent     uint64        `json:"heartbeatsSent"`
	HeartbeatsReceived uint64        `json:"heartbeatsReceived"`
	Timestamp          time.Time     `json:"timestamp"`
}

// OrderingMetrics defines metrics for transaction ordering
type OrderingMetrics struct {
	TotalOrderingOperations uint64                      `json:"totalOrderingOperations"`
	AverageOrderingTime     time.Duration               `json:"averageOrderingTime"`
	OrderingsByStrategy     map[OrderingStrategy]uint64 `json:"orderingsByStrategy"`
	CurrentStrategy         OrderingStrategy            `json:"currentStrategy"`
	StrategyChanges         uint64                      `json:"strategyChanges"`
	OrderingErrors          uint64                      `json:"orderingErrors"`
	Timestamp               time.Time                   `json:"timestamp"`
}

// ChronosMetrics defines comprehensive metrics for YaoChronos
type ChronosMetrics struct {
	ServerMetrics    *ServerMetrics    `json:"serverMetrics"`
	SequencerMetrics *SequencerMetrics `json:"sequencerMetrics"`
	ElectionMetrics  *ElectionMetrics  `json:"electionMetrics"`
	OrderingMetrics  *OrderingMetrics  `json:"orderingMetrics"`
	MessageMetrics   *MessageMetrics   `json:"messageMetrics"`
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

// ChronosConfig defines comprehensive configuration for YaoChronos
type ChronosConfig struct {
	// Node configuration
	NodeID   string `json:"nodeId"`
	NodeName string `json:"nodeName"`

	// Server configuration
	BindAddress string `json:"bindAddress"`
	Port        int    `json:"port"`
	EnableTLS   bool   `json:"enableTls"`
	TLSCert     string `json:"tlsCert,omitempty"`
	TLSKey      string `json:"tlsKey,omitempty"`

	// Sequencer configuration
	BlockInterval time.Duration `json:"blockInterval"`
	MaxTxPerBlock int           `json:"maxTxPerBlock"`
	MaxBlockSize  int64         `json:"maxBlockSize"`
	TxTimeout     time.Duration `json:"txTimeout"`
	QueueCapacity int           `json:"queueCapacity"`

	// Ordering configuration
	DefaultOrderingStrategy OrderingStrategy `json:"defaultOrderingStrategy"`
	AllowStrategyChange     bool             `json:"allowStrategyChange"`

	// Election configuration
	EnableElection    bool          `json:"enableElection"`
	ElectionTimeout   time.Duration `json:"electionTimeout"`
	HeartbeatInterval time.Duration `json:"heartbeatInterval"`
	MaxCandidates     int           `json:"maxCandidates"`

	// Message queue configuration
	RocketMQEndpoints []string `json:"rocketmqEndpoints"`
	RocketMQAccessKey string   `json:"rocketmqAccessKey"`
	RocketMQSecretKey string   `json:"rocketmqSecretKey"`
	ConsumerGroup     string   `json:"consumerGroup"`
	ProducerGroup     string   `json:"producerGroup"`

	// Performance configuration
	WorkerCount       int           `json:"workerCount"`
	ProcessingTimeout time.Duration `json:"processingTimeout"`
	BatchSize         int           `json:"batchSize"`

	// Monitoring configuration
	EnableMetrics       bool          `json:"enableMetrics"`
	MetricsInterval     time.Duration `json:"metricsInterval"`
	EnableHealthCheck   bool          `json:"enableHealthCheck"`
	HealthCheckInterval time.Duration `json:"healthCheckInterval"`

	// Logging configuration
	LogLevel  string `json:"logLevel"`
	LogFormat string `json:"logFormat"`
	LogFile   string `json:"logFile"`
}

// BatchProcessor defines the interface for batch transaction processing
type BatchProcessor interface {
	// CreateBatch creates a new transaction batch
	CreateBatch(ctx context.Context, maxSize int, timeout time.Duration) (*TransactionBatch, error)

	// AddTransactionToBatch adds a transaction to an existing batch
	AddTransactionToBatch(ctx context.Context, batchID string, tx *types.Transaction) error

	// ProcessBatch processes a complete batch of transactions
	ProcessBatch(ctx context.Context, batch *TransactionBatch) (*BatchResult, error)

	// GetBatchStatus returns the status of a batch
	GetBatchStatus(ctx context.Context, batchID string) (*BatchStatus, error)

	// GetActiveBatches returns all currently active batches
	GetActiveBatches(ctx context.Context) ([]*TransactionBatch, error)

	// CancelBatch cancels a pending batch
	CancelBatch(ctx context.Context, batchID string) error

	// GetBatchMetrics returns batch processing metrics
	GetBatchMetrics() *BatchMetrics
}

// FailoverManager defines the interface for managing leader failover
type FailoverManager interface {
	// TriggerFailover triggers a manual failover
	TriggerFailover(ctx context.Context, reason string) error

	// GetFailoverHistory returns the history of failover events
	GetFailoverHistory(ctx context.Context) ([]*FailoverEvent, error)

	// SetFailoverPolicy sets the failover policy
	SetFailoverPolicy(ctx context.Context, policy *FailoverPolicy) error

	// GetFailoverPolicy returns the current failover policy
	GetFailoverPolicy(ctx context.Context) (*FailoverPolicy, error)

	// CheckClusterHealth checks the health of all cluster nodes
	CheckClusterHealth(ctx context.Context) (*ClusterHealth, error)

	// RegisterFailoverCallback registers a callback for failover events
	RegisterFailoverCallback(callback FailoverCallback) error

	// GetFailoverMetrics returns failover metrics
	GetFailoverMetrics() *FailoverMetrics
}

// CheckpointManager defines the interface for managing sequencer checkpoints
type CheckpointManager interface {
	// CreateCheckpoint creates a new checkpoint at current state
	CreateCheckpoint(ctx context.Context, reason string) (*Checkpoint, error)

	// RestoreCheckpoint restores state from a checkpoint
	RestoreCheckpoint(ctx context.Context, checkpointID string) error

	// GetCheckpoint returns checkpoint information
	GetCheckpoint(ctx context.Context, checkpointID string) (*Checkpoint, error)

	// ListCheckpoints returns all available checkpoints
	ListCheckpoints(ctx context.Context, limit int, offset int) ([]*Checkpoint, error)

	// DeleteCheckpoint deletes a checkpoint
	DeleteCheckpoint(ctx context.Context, checkpointID string) error

	// ValidateCheckpoint validates checkpoint integrity
	ValidateCheckpoint(ctx context.Context, checkpointID string) (*CheckpointValidation, error)

	// GetCheckpointMetrics returns checkpoint metrics
	GetCheckpointMetrics() *CheckpointMetrics
}

// PerformanceOptimizer defines the interface for performance optimization
type PerformanceOptimizer interface {
	// OptimizeTransactionOrdering optimizes transaction ordering based on analysis
	OptimizeTransactionOrdering(ctx context.Context, transactions []*types.Transaction) ([]*types.Transaction, error)

	// AnalyzePerformanceBottlenecks analyzes system performance bottlenecks
	AnalyzePerformanceBottlenecks(ctx context.Context) (*PerformanceAnalysis, error)

	// ApplyOptimization applies a performance optimization
	ApplyOptimization(ctx context.Context, optimization *OptimizationConfig) error

	// GetOptimizationRecommendations returns optimization recommendations
	GetOptimizationRecommendations(ctx context.Context) ([]*OptimizationRecommendation, error)

	// GetPerformanceProfile returns detailed performance profiling data
	GetPerformanceProfile(ctx context.Context, duration time.Duration) (*PerformanceProfile, error)

	// EnablePerformanceTuning enables automatic performance tuning
	EnablePerformanceTuning(ctx context.Context, config *TuningConfig) error

	// GetOptimizationMetrics returns optimization metrics
	GetOptimizationMetrics() *OptimizationMetrics
}

// StateTracker defines the interface for tracking sequencer state
type StateTracker interface {
	// GetCurrentState returns the current sequencer state
	GetCurrentState(ctx context.Context) (*SequencerState, error)

	// GetStateHistory returns sequencer state history
	GetStateHistory(ctx context.Context, from, to time.Time) ([]*SequencerStateSnapshot, error)

	// ValidateState validates the current sequencer state
	ValidateState(ctx context.Context) (*StateValidation, error)

	// GetStateMetrics returns state tracking metrics
	GetStateMetrics() *StateMetrics
}

// Supporting types for new interfaces

// TransactionBatch represents a batch of transactions for processing
type TransactionBatch struct {
	BatchID      string               `json:"batchId"`
	Transactions []*types.Transaction `json:"transactions"`
	CreatedAt    time.Time            `json:"createdAt"`
	MaxSize      int                  `json:"maxSize"`
	Timeout      time.Duration        `json:"timeout"`
	Status       BatchStatus          `json:"status"`
	Priority     BatchPriority        `json:"priority"`
	Metadata     map[string]string    `json:"metadata,omitempty"`
}

// BatchStatus defines batch processing status
type BatchStatus string

const (
	BatchStatusPending    BatchStatus = "pending"
	BatchStatusProcessing BatchStatus = "processing"
	BatchStatusCompleted  BatchStatus = "completed"
	BatchStatusFailed     BatchStatus = "failed"
	BatchStatusCancelled  BatchStatus = "cancelled"
	BatchStatusTimeout    BatchStatus = "timeout"
)

// BatchPriority defines batch priority levels
type BatchPriority string

const (
	BatchPriorityLow      BatchPriority = "low"
	BatchPriorityNormal   BatchPriority = "normal"
	BatchPriorityHigh     BatchPriority = "high"
	BatchPriorityCritical BatchPriority = "critical"
)

// BatchResult represents the result of batch processing
type BatchResult struct {
	BatchID           string        `json:"batchId"`
	ProcessedBlock    *types.Block  `json:"processedBlock"`
	ProcessingTime    time.Duration `json:"processingTime"`
	SuccessfulTxCount int           `json:"successfulTxCount"`
	FailedTxCount     int           `json:"failedTxCount"`
	TotalGasUsed      uint64        `json:"totalGasUsed"`
	Status            BatchStatus   `json:"status"`
	ErrorMessage      string        `json:"errorMessage,omitempty"`
	CompletedAt       time.Time     `json:"completedAt"`
}

// BatchMetrics provides batch processing metrics
type BatchMetrics struct {
	TotalBatches          uint64        `json:"totalBatches"`
	CompletedBatches      uint64        `json:"completedBatches"`
	FailedBatches         uint64        `json:"failedBatches"`
	CancelledBatches      uint64        `json:"cancelledBatches"`
	AverageBatchSize      int           `json:"averageBatchSize"`
	AverageProcessingTime time.Duration `json:"averageProcessingTime"`
	BatchesPerSecond      float64       `json:"batchesPerSecond"`
	ActiveBatches         int           `json:"activeBatches"`
	BatchUtilization      float64       `json:"batchUtilization"` // percentage
	Timestamp             time.Time     `json:"timestamp"`
}

// FailoverEvent represents a failover event
type FailoverEvent struct {
	EventID      string          `json:"eventId"`
	OccurredAt   time.Time       `json:"occurredAt"`
	Reason       string          `json:"reason"`
	OldLeader    *LeaderInfo     `json:"oldLeader"`
	NewLeader    *LeaderInfo     `json:"newLeader"`
	Duration     time.Duration   `json:"duration"`
	EventType    FailoverType    `json:"eventType"`
	Success      bool            `json:"success"`
	ErrorMessage string          `json:"errorMessage,omitempty"`
	Impact       *FailoverImpact `json:"impact"`
}

// FailoverType defines types of failover events
type FailoverType string

const (
	FailoverTypeAutomatic FailoverType = "automatic"
	FailoverTypeManual    FailoverType = "manual"
	FailoverTypePlanned   FailoverType = "planned"
	FailoverTypeForced    FailoverType = "forced"
)

// FailoverImpact describes the impact of a failover event
type FailoverImpact struct {
	DowntimeDuration time.Duration `json:"downtimeDuration"`
	LostTransactions int           `json:"lostTransactions"`
	DelayedBlocks    int           `json:"delayedBlocks"`
	AffectedClients  int           `json:"affectedClients"`
	RecoveryTime     time.Duration `json:"recoveryTime"`
}

// FailoverPolicy defines the policy for automatic failover
type FailoverPolicy struct {
	EnableAutoFailover  bool          `json:"enableAutoFailover"`
	HealthCheckInterval time.Duration `json:"healthCheckInterval"`
	FailureThreshold    int           `json:"failureThreshold"`
	RecoveryThreshold   int           `json:"recoveryThreshold"`
	MaxFailoverAttempts int           `json:"maxFailoverAttempts"`
	FailoverCooldown    time.Duration `json:"failoverCooldown"`
	NotificationEnabled bool          `json:"notificationEnabled"`
	PreferredSuccessors []string      `json:"preferredSuccessors"`
}

// FailoverCallback defines the callback for failover events
type FailoverCallback func(event *FailoverEvent)

// ClusterHealth represents the health status of the cluster
type ClusterHealth struct {
	OverallHealth  HealthStatus          `json:"overallHealth"`
	NodeHealth     map[string]NodeHealth `json:"nodeHealth"`
	CurrentLeader  string                `json:"currentLeader"`
	HealthyNodes   int                   `json:"healthyNodes"`
	UnhealthyNodes int                   `json:"unhealthyNodes"`
	TotalNodes     int                   `json:"totalNodes"`
	LastUpdate     time.Time             `json:"lastUpdate"`
}

// HealthStatus defines health status levels
type HealthStatus string

const (
	HealthStatusHealthy  HealthStatus = "healthy"
	HealthStatusWarning  HealthStatus = "warning"
	HealthStatusCritical HealthStatus = "critical"
	HealthStatusUnknown  HealthStatus = "unknown"
)

// NodeHealth represents the health of a specific node
type NodeHealth struct {
	NodeID              string                 `json:"nodeId"`
	Status              HealthStatus           `json:"status"`
	LastHeartbeat       time.Time              `json:"lastHeartbeat"`
	ResponseTime        time.Duration          `json:"responseTime"`
	IsLeader            bool                   `json:"isLeader"`
	ConsecutiveFailures int                    `json:"consecutiveFailures"`
	ErrorMessage        string                 `json:"errorMessage,omitempty"`
	Metrics             map[string]interface{} `json:"metrics,omitempty"`
}

// FailoverMetrics provides failover-related metrics
type FailoverMetrics struct {
	TotalFailovers      uint64        `json:"totalFailovers"`
	SuccessfulFailovers uint64        `json:"successfulFailovers"`
	FailedFailovers     uint64        `json:"failedFailovers"`
	AverageFailoverTime time.Duration `json:"averageFailoverTime"`
	MTTR                time.Duration `json:"mttr"` // Mean Time To Recovery
	MTBF                time.Duration `json:"mtbf"` // Mean Time Between Failures
	CurrentUptime       time.Duration `json:"currentUptime"`
	TotalDowntime       time.Duration `json:"totalDowntime"`
	AvailabilityPercent float64       `json:"availabilityPercent"`
	Timestamp           time.Time     `json:"timestamp"`
}

// Checkpoint represents a sequencer checkpoint
type Checkpoint struct {
	CheckpointID string            `json:"checkpointId"`
	CreatedAt    time.Time         `json:"createdAt"`
	BlockNumber  uint64            `json:"blockNumber"`
	BlockHash    common.Hash       `json:"blockHash"`
	StateRoot    common.Hash       `json:"stateRoot"`
	Reason       string            `json:"reason"`
	Size         int64             `json:"size"`
	Verified     bool              `json:"verified"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CreatedBy    string            `json:"createdBy"`
}

// CheckpointValidation represents checkpoint validation result
type CheckpointValidation struct {
	CheckpointID   string            `json:"checkpointId"`
	IsValid        bool              `json:"isValid"`
	ValidationTime time.Duration     `json:"validationTime"`
	Issues         []ValidationIssue `json:"issues,omitempty"`
	ValidatedAt    time.Time         `json:"validatedAt"`
	ValidatedBy    string            `json:"validatedBy"`
}

// ValidationIssue represents an issue found during validation
type ValidationIssue struct {
	Severity    IssueSeverity `json:"severity"`
	Type        IssueType     `json:"type"`
	Description string        `json:"description"`
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
	IssueTypeCorruption    IssueType = "corruption"
	IssueTypeInconsistency IssueType = "inconsistency"
	IssueTypeMissing       IssueType = "missing"
	IssueTypeIntegrity     IssueType = "integrity"
)

// CheckpointMetrics provides checkpoint-related metrics
type CheckpointMetrics struct {
	TotalCheckpoints      uint64        `json:"totalCheckpoints"`
	ValidCheckpoints      uint64        `json:"validCheckpoints"`
	CorruptedCheckpoints  uint64        `json:"corruptedCheckpoints"`
	AverageCheckpointSize int64         `json:"averageCheckpointSize"`
	AverageCreationTime   time.Duration `json:"averageCreationTime"`
	StorageUsage          int64         `json:"storageUsage"`
	LastCheckpointTime    time.Time     `json:"lastCheckpointTime"`
	Timestamp             time.Time     `json:"timestamp"`
}

// PerformanceAnalysis represents system performance analysis results
type PerformanceAnalysis struct {
	AnalysisID      string                        `json:"analysisId"`
	PerformedAt     time.Time                     `json:"performedAt"`
	Duration        time.Duration                 `json:"duration"`
	Bottlenecks     []*PerformanceBottleneck      `json:"bottlenecks"`
	Recommendations []*OptimizationRecommendation `json:"recommendations"`
	OverallScore    float64                       `json:"overallScore"` // 0-100
	Summary         string                        `json:"summary"`
}

// PerformanceBottleneck represents a detected performance bottleneck
type PerformanceBottleneck struct {
	Type        BottleneckType         `json:"type"`
	Severity    IssueSeverity          `json:"severity"`
	Component   string                 `json:"component"`
	Description string                 `json:"description"`
	Impact      float64                `json:"impact"` // 0-100
	Metrics     map[string]interface{} `json:"metrics"`
}

// BottleneckType defines types of performance bottlenecks
type BottleneckType string

const (
	BottleneckTypeCPU     BottleneckType = "cpu"
	BottleneckTypeMemory  BottleneckType = "memory"
	BottleneckTypeNetwork BottleneckType = "network"
	BottleneckTypeStorage BottleneckType = "storage"
	BottleneckTypeQueue   BottleneckType = "queue"
	BottleneckTypeLock    BottleneckType = "lock"
)

// OptimizationRecommendation represents a performance optimization recommendation
type OptimizationRecommendation struct {
	RecommendationID     string               `json:"recommendationId"`
	Title                string               `json:"title"`
	Description          string               `json:"description"`
	Priority             OptimizationPriority `json:"priority"`
	EstimatedImpact      float64              `json:"estimatedImpact"` // 0-100
	ImplementationEffort string               `json:"implementationEffort"`
	Configuration        *OptimizationConfig  `json:"configuration,omitempty"`
}

// OptimizationPriority defines optimization priority levels
type OptimizationPriority string

const (
	OptimizationPriorityLow      OptimizationPriority = "low"
	OptimizationPriorityMedium   OptimizationPriority = "medium"
	OptimizationPriorityHigh     OptimizationPriority = "high"
	OptimizationPriorityCritical OptimizationPriority = "critical"
)

// OptimizationConfig represents configuration for an optimization
type OptimizationConfig struct {
	OptimizationType string                 `json:"optimizationType"`
	Parameters       map[string]interface{} `json:"parameters"`
	ApplyImmediately bool                   `json:"applyImmediately"`
	RollbackTime     *time.Time             `json:"rollbackTime,omitempty"`
}

// PerformanceProfile represents detailed performance profiling data
type PerformanceProfile struct {
	ProfileID        string                 `json:"profileId"`
	GeneratedAt      time.Time              `json:"generatedAt"`
	Duration         time.Duration          `json:"duration"`
	CPUProfile       *CPUProfileData        `json:"cpuProfile"`
	MemoryProfile    *MemoryProfileData     `json:"memoryProfile"`
	GoroutineProfile *GoroutineProfileData  `json:"goroutineProfile"`
	BlockProfile     *BlockProfileData      `json:"blockProfile"`
	CustomMetrics    map[string]interface{} `json:"customMetrics"`
}

// TuningConfig represents configuration for automatic performance tuning
type TuningConfig struct {
	EnableAutoTuning    bool          `json:"enableAutoTuning"`
	TuningInterval      time.Duration `json:"tuningInterval"`
	AggressivenessLevel int           `json:"aggressivenessLevel"` // 1-10
	SafetyMargin        float64       `json:"safetyMargin"`        // 0-1
	MaxChangesPerHour   int           `json:"maxChangesPerHour"`
}

// Profile data types
type CPUProfileData struct {
	TotalSamples   uint64             `json:"totalSamples"`
	SampleDuration time.Duration      `json:"sampleDuration"`
	TopFunctions   []*FunctionProfile `json:"topFunctions"`
	CPUUtilization float64            `json:"cpuUtilization"`
}

type MemoryProfileData struct {
	AllocObjects uint64        `json:"allocObjects"`
	AllocBytes   uint64        `json:"allocBytes"`
	InUseObjects uint64        `json:"inUseObjects"`
	InUseBytes   uint64        `json:"inUseBytes"`
	HeapSize     uint64        `json:"heapSize"`
	GCCount      uint64        `json:"gcCount"`
	TotalGCPause time.Duration `json:"totalGcPause"`
}

type GoroutineProfileData struct {
	TotalGoroutines   int `json:"totalGoroutines"`
	RunningGoroutines int `json:"runningGoroutines"`
	BlockedGoroutines int `json:"blockedGoroutines"`
	WaitingGoroutines int `json:"waitingGoroutines"`
}

type BlockProfileData struct {
	TotalBlocks      uint64        `json:"totalBlocks"`
	TotalBlockTime   time.Duration `json:"totalBlockTime"`
	AverageBlockTime time.Duration `json:"averageBlockTime"`
}

type FunctionProfile struct {
	Function    string        `json:"function"`
	SampleCount uint64        `json:"sampleCount"`
	Percentage  float64       `json:"percentage"`
	CumTime     time.Duration `json:"cumTime"`
}

// OptimizationMetrics provides optimization-related metrics
type OptimizationMetrics struct {
	TotalOptimizations      uint64             `json:"totalOptimizations"`
	SuccessfulOptimizations uint64             `json:"successfulOptimizations"`
	FailedOptimizations     uint64             `json:"failedOptimizations"`
	AverageImpact           float64            `json:"averageImpact"`
	LastOptimization        time.Time          `json:"lastOptimization"`
	PerformanceGain         float64            `json:"performanceGain"` // percentage
	ResourceSavings         map[string]float64 `json:"resourceSavings"`
	Timestamp               time.Time          `json:"timestamp"`
}

// SequencerState represents the current state of the sequencer
type SequencerState struct {
	NodeID              string                 `json:"nodeId"`
	IsLeader            bool                   `json:"isLeader"`
	CurrentBlockNumber  uint64                 `json:"currentBlockNumber"`
	LastBlockHash       common.Hash            `json:"lastBlockHash"`
	PendingTransactions int                    `json:"pendingTransactions"`
	QueuedBatches       int                    `json:"queuedBatches"`
	ProcessingRate      float64                `json:"processingRate"`
	State               SequencerStateType     `json:"state"`
	LastUpdate          time.Time              `json:"lastUpdate"`
	Metrics             map[string]interface{} `json:"metrics"`
}

// SequencerStateType defines sequencer state types
type SequencerStateType string

const (
	SequencerStateIdle       SequencerStateType = "idle"
	SequencerStateProcessing SequencerStateType = "processing"
	SequencerStateSyncing    SequencerStateType = "syncing"
	SequencerStateError      SequencerStateType = "error"
	SequencerStateStopped    SequencerStateType = "stopped"
)

// SequencerStateSnapshot represents a snapshot of sequencer state at a point in time
type SequencerStateSnapshot struct {
	SnapshotID string          `json:"snapshotId"`
	Timestamp  time.Time       `json:"timestamp"`
	State      *SequencerState `json:"state"`
}

// StateValidation represents the result of state validation
type StateValidation struct {
	IsValid        bool              `json:"isValid"`
	ValidationTime time.Duration     `json:"validationTime"`
	Issues         []ValidationIssue `json:"issues,omitempty"`
	ValidatedAt    time.Time         `json:"validatedAt"`
}

// StateMetrics provides state tracking metrics
type StateMetrics struct {
	StateValidations      uint64        `json:"stateValidations"`
	ValidStates           uint64        `json:"validStates"`
	InvalidStates         uint64        `json:"invalidStates"`
	StateSnapshots        uint64        `json:"stateSnapshots"`
	AverageValidationTime time.Duration `json:"averageValidationTime"`
	LastValidation        time.Time     `json:"lastValidation"`
	Timestamp             time.Time     `json:"timestamp"`
}
