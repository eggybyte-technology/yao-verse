package interfaces

import (
	"context"
	"time"

	oracleinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
	sharedinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
)

// DataPersister defines the interface for persisting blockchain data
type DataPersister interface {
	// PersistExecutionResult persists the execution result of a block
	PersistExecutionResult(ctx context.Context, result *types.ExecutionResult) error

	// PersistBlock persists a block to storage
	PersistBlock(ctx context.Context, block *types.Block) error

	// PersistTransactions persists transactions to storage
	PersistTransactions(ctx context.Context, txs []*types.Transaction) error

	// PersistReceipts persists transaction receipts to storage
	PersistReceipts(ctx context.Context, receipts []*types.Receipt) error

	// PersistStateDiff persists state changes to storage
	PersistStateDiff(ctx context.Context, stateDiff *types.StateDiff) error

	// PersistLogs persists event logs to storage
	PersistLogs(ctx context.Context, logs []*types.Log) error

	// GetPersistenceMetrics returns persistence metrics
	GetPersistenceMetrics() *PersistenceMetrics
}

// StateCommitter defines the interface for committing state changes
type StateCommitter interface {
	// CommitState commits state changes in a transaction
	CommitState(ctx context.Context, result *types.ExecutionResult) error

	// RollbackState rolls back uncommitted state changes
	RollbackState(ctx context.Context, blockNumber uint64) error

	// GetStateRoot returns the state root for a given block
	GetStateRoot(ctx context.Context, blockNumber uint64) (common.Hash, error)

	// ValidateStateConsistency validates state consistency
	ValidateStateConsistency(ctx context.Context, fromBlock, toBlock uint64) (*StateConsistencyReport, error)

	// GetCommitMetrics returns commit metrics
	GetCommitMetrics() *CommitMetrics
}

// CacheInvalidator defines the interface for cache invalidation
type CacheInvalidator interface {
	// InvalidateCaches invalidates relevant caches after state commit
	InvalidateCaches(ctx context.Context, stateDiff *types.StateDiff) error

	// GenerateInvalidationKeys generates cache keys that need to be invalidated
	GenerateInvalidationKeys(ctx context.Context, stateDiff *types.StateDiff) ([]*types.CacheKey, error)

	// BroadcastInvalidation broadcasts cache invalidation message
	BroadcastInvalidation(ctx context.Context, message *types.CacheInvalidationMessage) error

	// GetInvalidationMetrics returns invalidation metrics
	GetInvalidationMetrics() *InvalidationMetrics
}

// BackupManager defines the interface for backup and recovery operations
type BackupManager interface {
	// CreateBackup creates a full or incremental backup
	CreateBackup(ctx context.Context, config *BackupConfig) (*BackupInfo, error)

	// RestoreBackup restores data from a backup
	RestoreBackup(ctx context.Context, backupID string, config *RestoreConfig) error

	// ListBackups lists all available backups
	ListBackups(ctx context.Context, filter *BackupFilter) ([]*BackupInfo, error)

	// DeleteBackup deletes a backup
	DeleteBackup(ctx context.Context, backupID string) error

	// ValidateBackup validates backup integrity
	ValidateBackup(ctx context.Context, backupID string) (*BackupValidation, error)

	// GetBackupMetrics returns backup metrics
	GetBackupMetrics() *BackupMetrics
}

// IndexManager defines the interface for database index management
type IndexManager interface {
	// CreateIndex creates a database index
	CreateIndex(ctx context.Context, indexDef *IndexDefinition) error

	// DropIndex drops a database index
	DropIndex(ctx context.Context, indexName string) error

	// RebuildIndex rebuilds an existing index
	RebuildIndex(ctx context.Context, indexName string) error

	// ListIndexes lists all database indexes
	ListIndexes(ctx context.Context, tableName string) ([]*IndexInfo, error)

	// GetIndexStatistics returns index usage statistics
	GetIndexStatistics(ctx context.Context, indexName string) (*IndexStatistics, error)

	// OptimizeIndexes analyzes and optimizes database indexes
	OptimizeIndexes(ctx context.Context) (*IndexOptimization, error)

	// GetIndexMetrics returns index-related metrics
	GetIndexMetrics() *IndexMetrics
}

// DataValidator defines the interface for data validation and integrity checks
type DataValidator interface {
	// ValidateBlock validates block data integrity
	ValidateBlock(ctx context.Context, blockNumber uint64) (*BlockValidation, error)

	// ValidateTransaction validates transaction data integrity
	ValidateTransaction(ctx context.Context, txHash common.Hash) (*TransactionValidation, error)

	// ValidateStateData validates state data integrity
	ValidateStateData(ctx context.Context, fromBlock, toBlock uint64) (*StateValidation, error)

	// PerformIntegrityCheck performs comprehensive data integrity check
	PerformIntegrityCheck(ctx context.Context, config *IntegrityCheckConfig) (*IntegrityCheckReport, error)

	// GetValidationMetrics returns validation metrics
	GetValidationMetrics() *ValidationMetrics
}

// TransactionManager defines the interface for database transaction management
type TransactionManager interface {
	// BeginTransaction starts a new database transaction
	BeginTransaction(ctx context.Context, options *TransactionOptions) (Transaction, error)

	// GetActiveTransactions returns all active transactions
	GetActiveTransactions(ctx context.Context) ([]*TransactionInfo, error)

	// KillTransaction kills a long-running transaction
	KillTransaction(ctx context.Context, transactionID string) error

	// GetTransactionMetrics returns transaction metrics
	GetTransactionMetrics() *TransactionMetrics
}

// PerformanceOptimizer defines the interface for database performance optimization
type PerformanceOptimizer interface {
	// AnalyzePerformance analyzes database performance
	AnalyzePerformance(ctx context.Context, duration time.Duration) (*PerformanceAnalysis, error)

	// OptimizeQueries analyzes and optimizes slow queries
	OptimizeQueries(ctx context.Context) (*QueryOptimization, error)

	// UpdateTableStatistics updates table statistics for better query planning
	UpdateTableStatistics(ctx context.Context, tableName string) error

	// CompactTables compacts database tables to reclaim space
	CompactTables(ctx context.Context, tables []string) (*CompactionResult, error)

	// GetPerformanceMetrics returns performance metrics
	GetPerformanceMetrics() *PerformanceMetrics
}

// ArchiveServer defines the main YaoArchive server interface
type ArchiveServer interface {
	// Start starts the archive server
	Start(ctx context.Context) error

	// Stop stops the archive server gracefully
	Stop(ctx context.Context) error

	// GetDataPersister returns the data persister instance
	GetDataPersister() DataPersister

	// GetStateCommitter returns the state committer instance
	GetStateCommitter() StateCommitter

	// GetCacheInvalidator returns the cache invalidator instance
	GetCacheInvalidator() CacheInvalidator

	// GetBackupManager returns the backup manager instance
	GetBackupManager() BackupManager

	// GetIndexManager returns the index manager instance
	GetIndexManager() IndexManager

	// GetDataValidator returns the data validator instance
	GetDataValidator() DataValidator

	// GetTransactionManager returns the transaction manager instance
	GetTransactionManager() TransactionManager

	// GetPerformanceOptimizer returns the performance optimizer instance
	GetPerformanceOptimizer() PerformanceOptimizer

	// GetStorage returns the L3 storage interface
	GetStorage() oracleinterfaces.L3Storage

	// GetMessageProducer returns the message producer for publishing invalidations
	GetMessageProducer() sharedinterfaces.MessageProducer

	// GetMessageConsumer returns the message consumer for receiving execution results
	GetMessageConsumer() sharedinterfaces.MessageConsumer

	// GetStatus returns the current server status
	GetStatus(ctx context.Context) *ArchiveStatus

	// GetMetrics returns comprehensive server metrics
	GetMetrics(ctx context.Context) *ArchiveMetrics

	// HealthCheck performs a comprehensive health check
	HealthCheck(ctx context.Context) error
}

// Supporting types for new and existing interfaces

// Transaction represents a database transaction
type Transaction interface {
	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error

	// GetID returns the transaction ID
	GetID() string

	// IsActive returns true if transaction is active
	IsActive() bool
}

// BackupConfig defines configuration for backup operations
type BackupConfig struct {
	BackupType      BackupType        `json:"backupType"`
	Compression     CompressionType   `json:"compression"`
	IncludeIndexes  bool              `json:"includeIndexes"`
	IncludeLogs     bool              `json:"includeLogs"`
	FromBlock       *uint64           `json:"fromBlock,omitempty"`
	ToBlock         *uint64           `json:"toBlock,omitempty"`
	StorageLocation string            `json:"storageLocation"`
	Encryption      *EncryptionConfig `json:"encryption,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// BackupType defines types of backups
type BackupType string

const (
	BackupTypeFull         BackupType = "full"
	BackupTypeIncremental  BackupType = "incremental"
	BackupTypeDifferential BackupType = "differential"
	BackupTypeSnapshot     BackupType = "snapshot"
)

// CompressionType defines compression types
type CompressionType string

const (
	CompressionTypeNone CompressionType = "none"
	CompressionTypeGzip CompressionType = "gzip"
	CompressionTypeLZ4  CompressionType = "lz4"
	CompressionTypeZstd CompressionType = "zstd"
)

// EncryptionConfig defines encryption configuration
type EncryptionConfig struct {
	Enabled   bool   `json:"enabled"`
	Algorithm string `json:"algorithm"`
	KeyID     string `json:"keyId"`
}

// BackupInfo represents information about a backup
type BackupInfo struct {
	BackupID        string            `json:"backupId"`
	BackupType      BackupType        `json:"backupType"`
	CreatedAt       time.Time         `json:"createdAt"`
	CompletedAt     time.Time         `json:"completedAt"`
	Size            int64             `json:"size"`
	CompressedSize  int64             `json:"compressedSize"`
	FromBlock       uint64            `json:"fromBlock"`
	ToBlock         uint64            `json:"toBlock"`
	StorageLocation string            `json:"storageLocation"`
	IsEncrypted     bool              `json:"isEncrypted"`
	Checksum        string            `json:"checksum"`
	Status          BackupStatus      `json:"status"`
	ErrorMessage    string            `json:"errorMessage,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// BackupStatus defines backup status
type BackupStatus string

const (
	BackupStatusPending   BackupStatus = "pending"
	BackupStatusRunning   BackupStatus = "running"
	BackupStatusCompleted BackupStatus = "completed"
	BackupStatusFailed    BackupStatus = "failed"
	BackupStatusCorrupted BackupStatus = "corrupted"
)

// RestoreConfig defines configuration for restore operations
type RestoreConfig struct {
	TargetDatabase  string            `json:"targetDatabase"`
	ReplaceExisting bool              `json:"replaceExisting"`
	SkipIndexes     bool              `json:"skipIndexes"`
	SkipLogs        bool              `json:"skipLogs"`
	RestoreToBlock  *uint64           `json:"restoreToBlock,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// BackupFilter defines filter criteria for backup queries
type BackupFilter struct {
	BackupType *BackupType   `json:"backupType,omitempty"`
	FromDate   *time.Time    `json:"fromDate,omitempty"`
	ToDate     *time.Time    `json:"toDate,omitempty"`
	Status     *BackupStatus `json:"status,omitempty"`
	Limit      int           `json:"limit,omitempty"`
	Offset     int           `json:"offset,omitempty"`
}

// BackupValidation represents backup validation results
type BackupValidation struct {
	BackupID       string            `json:"backupId"`
	IsValid        bool              `json:"isValid"`
	ValidationTime time.Duration     `json:"validationTime"`
	ChecksumMatch  bool              `json:"checksumMatch"`
	Issues         []ValidationIssue `json:"issues,omitempty"`
	ValidatedAt    time.Time         `json:"validatedAt"`
}

// ValidationIssue represents a validation issue
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
	IssueTypeChecksum      IssueType = "checksum"
)

// IndexDefinition defines a database index
type IndexDefinition struct {
	IndexName   string                 `json:"indexName"`
	TableName   string                 `json:"tableName"`
	Columns     []string               `json:"columns"`
	IndexType   IndexType              `json:"indexType"`
	IsUnique    bool                   `json:"isUnique"`
	WhereClause string                 `json:"whereClause,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// IndexType defines types of database indexes
type IndexType string

const (
	IndexTypeBTree  IndexType = "btree"
	IndexTypeHash   IndexType = "hash"
	IndexTypeGIN    IndexType = "gin"
	IndexTypeGiST   IndexType = "gist"
	IndexTypeSpGiST IndexType = "spgist"
	IndexTypeBRIN   IndexType = "brin"
)

// IndexInfo represents information about a database index
type IndexInfo struct {
	IndexName  string     `json:"indexName"`
	TableName  string     `json:"tableName"`
	Columns    []string   `json:"columns"`
	IndexType  IndexType  `json:"indexType"`
	IsUnique   bool       `json:"isUnique"`
	Size       int64      `json:"size"`
	CreatedAt  time.Time  `json:"createdAt"`
	LastUsed   *time.Time `json:"lastUsed,omitempty"`
	UsageCount uint64     `json:"usageCount"`
	IsValid    bool       `json:"isValid"`
	Definition string     `json:"definition"`
}

// IndexStatistics represents index usage statistics
type IndexStatistics struct {
	IndexName         string        `json:"indexName"`
	TableName         string        `json:"tableName"`
	TotalScans        uint64        `json:"totalScans"`
	TupleReads        uint64        `json:"tupleReads"`
	TupleFetches      uint64        `json:"tupleFetches"`
	LastScan          *time.Time    `json:"lastScan,omitempty"`
	AverageScanTime   time.Duration `json:"averageScanTime"`
	SelectivityRatio  float64       `json:"selectivityRatio"`
	MaintenanceCost   float64       `json:"maintenanceCost"`
	RecommendedAction string        `json:"recommendedAction,omitempty"`
}

// IndexOptimization represents index optimization results
type IndexOptimization struct {
	AnalyzedAt          time.Time              `json:"analyzedAt"`
	TotalIndexes        int                    `json:"totalIndexes"`
	UnusedIndexes       []string               `json:"unusedIndexes"`
	DuplicateIndexes    [][]string             `json:"duplicateIndexes"`
	MissingIndexes      []*IndexRecommendation `json:"missingIndexes"`
	OptimizationSummary string                 `json:"optimizationSummary"`
	EstimatedSavings    *OptimizationSavings   `json:"estimatedSavings"`
}

// IndexRecommendation represents a recommended index
type IndexRecommendation struct {
	TableName       string    `json:"tableName"`
	Columns         []string  `json:"columns"`
	IndexType       IndexType `json:"indexType"`
	Reason          string    `json:"reason"`
	EstimatedImpact string    `json:"estimatedImpact"`
	Priority        Priority  `json:"priority"`
}

// Priority defines recommendation priority levels
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// OptimizationSavings represents estimated optimization savings
type OptimizationSavings struct {
	StorageSavings  int64   `json:"storageSavings"`  // bytes
	QuerySpeedup    float64 `json:"querySpeedup"`    // percentage
	MaintenanceCost float64 `json:"maintenanceCost"` // percentage reduction
}

// BlockValidation represents block validation results
type BlockValidation struct {
	BlockNumber    uint64            `json:"blockNumber"`
	BlockHash      common.Hash       `json:"blockHash"`
	IsValid        bool              `json:"isValid"`
	ValidationTime time.Duration     `json:"validationTime"`
	Issues         []ValidationIssue `json:"issues,omitempty"`
	ValidatedAt    time.Time         `json:"validatedAt"`
}

// TransactionValidation represents transaction validation results
type TransactionValidation struct {
	TxHash         common.Hash       `json:"txHash"`
	BlockNumber    uint64            `json:"blockNumber"`
	IsValid        bool              `json:"isValid"`
	ValidationTime time.Duration     `json:"validationTime"`
	Issues         []ValidationIssue `json:"issues,omitempty"`
	ValidatedAt    time.Time         `json:"validatedAt"`
}

// StateValidation represents state validation results
type StateValidation struct {
	FromBlock      uint64            `json:"fromBlock"`
	ToBlock        uint64            `json:"toBlock"`
	IsValid        bool              `json:"isValid"`
	ValidationTime time.Duration     `json:"validationTime"`
	Issues         []ValidationIssue `json:"issues,omitempty"`
	ValidatedAt    time.Time         `json:"validatedAt"`
}

// IntegrityCheckConfig defines configuration for integrity checks
type IntegrityCheckConfig struct {
	CheckType       IntegrityCheckType `json:"checkType"`
	FromBlock       *uint64            `json:"fromBlock,omitempty"`
	ToBlock         *uint64            `json:"toBlock,omitempty"`
	IncludeLogs     bool               `json:"includeLogs"`
	IncludeReceipts bool               `json:"includeReceipts"`
	ParallelChecks  int                `json:"parallelChecks"`
	Timeout         time.Duration      `json:"timeout"`
}

// IntegrityCheckType defines types of integrity checks
type IntegrityCheckType string

const (
	IntegrityCheckTypeFull     IntegrityCheckType = "full"
	IntegrityCheckTypeBlocks   IntegrityCheckType = "blocks"
	IntegrityCheckTypeState    IntegrityCheckType = "state"
	IntegrityCheckTypeLogs     IntegrityCheckType = "logs"
	IntegrityCheckTypeReceipts IntegrityCheckType = "receipts"
)

// IntegrityCheckReport represents comprehensive integrity check results
type IntegrityCheckReport struct {
	CheckID         string                      `json:"checkId"`
	CheckType       IntegrityCheckType          `json:"checkType"`
	StartedAt       time.Time                   `json:"startedAt"`
	CompletedAt     time.Time                   `json:"completedAt"`
	Duration        time.Duration               `json:"duration"`
	BlocksChecked   uint64                      `json:"blocksChecked"`
	IssuesFound     int                         `json:"issuesFound"`
	OverallStatus   IntegrityStatus             `json:"overallStatus"`
	DetailedResults []*ComponentIntegrityResult `json:"detailedResults"`
	Summary         string                      `json:"summary"`
}

// IntegrityStatus defines integrity status
type IntegrityStatus string

const (
	IntegrityStatusHealthy IntegrityStatus = "healthy"
	IntegrityStatusWarning IntegrityStatus = "warning"
	IntegrityStatusCorrupt IntegrityStatus = "corrupt"
	IntegrityStatusFailed  IntegrityStatus = "failed"
)

// ComponentIntegrityResult represents integrity results for a specific component
type ComponentIntegrityResult struct {
	Component   string                 `json:"component"`
	Status      IntegrityStatus        `json:"status"`
	ChecksRun   int                    `json:"checksRun"`
	IssuesFound int                    `json:"issuesFound"`
	Issues      []ValidationIssue      `json:"issues,omitempty"`
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
}

// TransactionOptions defines options for database transactions
type TransactionOptions struct {
	IsolationLevel IsolationLevel      `json:"isolationLevel"`
	ReadOnly       bool                `json:"readOnly"`
	Timeout        time.Duration       `json:"timeout"`
	Priority       TransactionPriority `json:"priority"`
	Metadata       map[string]string   `json:"metadata,omitempty"`
}

// IsolationLevel defines transaction isolation levels
type IsolationLevel string

const (
	IsolationLevelReadUncommitted IsolationLevel = "read_uncommitted"
	IsolationLevelReadCommitted   IsolationLevel = "read_committed"
	IsolationLevelRepeatableRead  IsolationLevel = "repeatable_read"
	IsolationLevelSerializable    IsolationLevel = "serializable"
)

// TransactionPriority defines transaction priority levels
type TransactionPriority string

const (
	TransactionPriorityLow    TransactionPriority = "low"
	TransactionPriorityNormal TransactionPriority = "normal"
	TransactionPriorityHigh   TransactionPriority = "high"
)

// TransactionInfo represents information about a database transaction
type TransactionInfo struct {
	TransactionID  string              `json:"transactionId"`
	Status         TransactionStatus   `json:"status"`
	StartTime      time.Time           `json:"startTime"`
	IsolationLevel IsolationLevel      `json:"isolationLevel"`
	Priority       TransactionPriority `json:"priority"`
	Query          string              `json:"query,omitempty"`
	BlockedBy      []string            `json:"blockedBy,omitempty"`
	Blocking       []string            `json:"blocking,omitempty"`
	LocksHeld      int                 `json:"locksHeld"`
	RowsAffected   uint64              `json:"rowsAffected"`
}

// TransactionStatus defines transaction status
type TransactionStatus string

const (
	TransactionStatusActive     TransactionStatus = "active"
	TransactionStatusIdle       TransactionStatus = "idle"
	TransactionStatusBlocked    TransactionStatus = "blocked"
	TransactionStatusCommitted  TransactionStatus = "committed"
	TransactionStatusRolledBack TransactionStatus = "rolled_back"
)

// PerformanceAnalysis represents database performance analysis results
type PerformanceAnalysis struct {
	AnalysisID      string                       `json:"analysisId"`
	AnalyzedAt      time.Time                    `json:"analyzedAt"`
	Duration        time.Duration                `json:"duration"`
	SlowQueries     []*SlowQuery                 `json:"slowQueries"`
	ResourceUsage   *ResourceUsageAnalysis       `json:"resourceUsage"`
	Bottlenecks     []*PerformanceBottleneck     `json:"bottlenecks"`
	Recommendations []*PerformanceRecommendation `json:"recommendations"`
	OverallScore    float64                      `json:"overallScore"` // 0-100
}

// SlowQuery represents a slow database query
type SlowQuery struct {
	Query          string        `json:"query"`
	AverageTime    time.Duration `json:"averageTime"`
	TotalTime      time.Duration `json:"totalTime"`
	ExecutionCount uint64        `json:"executionCount"`
	RowsExamined   uint64        `json:"rowsExamined"`
	RowsReturned   uint64        `json:"rowsReturned"`
	FirstSeen      time.Time     `json:"firstSeen"`
	LastSeen       time.Time     `json:"lastSeen"`
	Recommendation string        `json:"recommendation,omitempty"`
}

// ResourceUsageAnalysis represents resource usage analysis
type ResourceUsageAnalysis struct {
	CPUUsage        float64 `json:"cpuUsage"`      // percentage
	MemoryUsage     float64 `json:"memoryUsage"`   // percentage
	DiskIORate      float64 `json:"diskIoRate"`    // MB/s
	NetworkIORate   float64 `json:"networkIoRate"` // MB/s
	ConnectionCount int     `json:"connectionCount"`
	ActiveQueries   int     `json:"activeQueries"`
	QueuedQueries   int     `json:"queuedQueries"`
}

// PerformanceBottleneck represents a performance bottleneck
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
	BottleneckTypeDisk    BottleneckType = "disk"
	BottleneckTypeNetwork BottleneckType = "network"
	BottleneckTypeLock    BottleneckType = "lock"
	BottleneckTypeIndex   BottleneckType = "index"
	BottleneckTypeQuery   BottleneckType = "query"
)

// PerformanceRecommendation represents a performance optimization recommendation
type PerformanceRecommendation struct {
	RecommendationID     string   `json:"recommendationId"`
	Title                string   `json:"title"`
	Description          string   `json:"description"`
	Priority             Priority `json:"priority"`
	EstimatedImpact      float64  `json:"estimatedImpact"` // 0-100
	ImplementationEffort string   `json:"implementationEffort"`
	Actions              []string `json:"actions"`
}

// QueryOptimization represents query optimization results
type QueryOptimization struct {
	OptimizedAt         time.Time         `json:"optimizedAt"`
	QueriesAnalyzed     int               `json:"queriesAnalyzed"`
	QueriesOptimized    int               `json:"queriesOptimized"`
	AverageSpeedupPct   float64           `json:"averageSpeedupPct"`
	OptimizedQueries    []*OptimizedQuery `json:"optimizedQueries"`
	OptimizationSummary string            `json:"optimizationSummary"`
}

// OptimizedQuery represents an optimized query
type OptimizedQuery struct {
	OriginalQuery    string        `json:"originalQuery"`
	OptimizedQuery   string        `json:"optimizedQuery"`
	OriginalTime     time.Duration `json:"originalTime"`
	OptimizedTime    time.Duration `json:"optimizedTime"`
	SpeedupPercent   float64       `json:"speedupPercent"`
	OptimizationType string        `json:"optimizationType"`
}

// CompactionResult represents table compaction results
type CompactionResult struct {
	CompactedAt     time.Time                `json:"compactedAt"`
	TablesCompacted []string                 `json:"tablesCompacted"`
	SpaceReclaimed  int64                    `json:"spaceReclaimed"` // bytes
	CompactionTime  time.Duration            `json:"compactionTime"`
	Results         []*TableCompactionResult `json:"results"`
}

// TableCompactionResult represents compaction results for a single table
type TableCompactionResult struct {
	TableName      string        `json:"tableName"`
	OriginalSize   int64         `json:"originalSize"`
	CompactedSize  int64         `json:"compactedSize"`
	SpaceReclaimed int64         `json:"spaceReclaimed"`
	CompactionTime time.Duration `json:"compactionTime"`
}

// StateConsistencyReport represents state consistency validation results
type StateConsistencyReport struct {
	FromBlock         uint64                   `json:"fromBlock"`
	ToBlock           uint64                   `json:"toBlock"`
	CheckedAt         time.Time                `json:"checkedAt"`
	IsConsistent      bool                     `json:"isConsistent"`
	IssuesFound       int                      `json:"issuesFound"`
	BlocksChecked     uint64                   `json:"blocksChecked"`
	ConsistencyIssues []*StateConsistencyIssue `json:"consistencyIssues,omitempty"`
}

// StateConsistencyIssue represents a state consistency issue
type StateConsistencyIssue struct {
	BlockNumber   uint64        `json:"blockNumber"`
	IssueType     IssueType     `json:"issueType"`
	Severity      IssueSeverity `json:"severity"`
	Description   string        `json:"description"`
	ExpectedValue interface{}   `json:"expectedValue,omitempty"`
	ActualValue   interface{}   `json:"actualValue,omitempty"`
}

// Status and metrics types

// ArchiveStatus represents the status of YaoArchive server
type ArchiveStatus struct {
	ServerStatus       ServerStatus   `json:"serverStatus"`
	LastCommittedBlock uint64         `json:"lastCommittedBlock"`
	PendingCommits     int            `json:"pendingCommits"`
	CommitRate         float64        `json:"commitRate"` // commits per second
	DatabaseHealth     DatabaseHealth `json:"databaseHealth"`
	BackupStatus       BackupHealth   `json:"backupStatus"`
	Uptime             time.Duration  `json:"uptime"`
	StartTime          time.Time      `json:"startTime"`
	Version            string         `json:"version"`
}

// ServerStatus defines server status
type ServerStatus string

const (
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusRunning  ServerStatus = "running"
	ServerStatusStopping ServerStatus = "stopping"
	ServerStatusStopped  ServerStatus = "stopped"
	ServerStatusError    ServerStatus = "error"
)

// DatabaseHealth represents database health status
type DatabaseHealth struct {
	Status          HealthStatus `json:"status"`
	ConnectionCount int          `json:"connectionCount"`
	ActiveQueries   int          `json:"activeQueries"`
	LockedTables    int          `json:"lockedTables"`
	DiskUsage       float64      `json:"diskUsage"` // percentage
	LastCheckup     time.Time    `json:"lastCheckup"`
}

// HealthStatus defines health status levels
type HealthStatus string

const (
	HealthStatusHealthy  HealthStatus = "healthy"
	HealthStatusWarning  HealthStatus = "warning"
	HealthStatusCritical HealthStatus = "critical"
	HealthStatusUnknown  HealthStatus = "unknown"
)

// BackupHealth represents backup system health
type BackupHealth struct {
	Status        HealthStatus `json:"status"`
	LastBackup    *time.Time   `json:"lastBackup,omitempty"`
	BackupsCount  int          `json:"backupsCount"`
	FailedBackups int          `json:"failedBackups"`
	StorageUsage  int64        `json:"storageUsage"`
}

// Comprehensive metrics types

// PersistenceMetrics provides data persistence metrics
type PersistenceMetrics struct {
	TotalOperations       uint64        `json:"totalOperations"`
	SuccessfulOperations  uint64        `json:"successfulOperations"`
	FailedOperations      uint64        `json:"failedOperations"`
	AverageOperationTime  time.Duration `json:"averageOperationTime"`
	BlocksPersisted       uint64        `json:"blocksPersisted"`
	TransactionsPersisted uint64        `json:"transactionsPersisted"`
	ReceiptsPersisted     uint64        `json:"receiptsPersisted"`
	LogsPersisted         uint64        `json:"logsPersisted"`
	Timestamp             time.Time     `json:"timestamp"`
}

// CommitMetrics provides state commit metrics
type CommitMetrics struct {
	TotalCommits      uint64        `json:"totalCommits"`
	SuccessfulCommits uint64        `json:"successfulCommits"`
	FailedCommits     uint64        `json:"failedCommits"`
	RolledBackCommits uint64        `json:"rolledBackCommits"`
	AverageCommitTime time.Duration `json:"averageCommitTime"`
	P99CommitTime     time.Duration `json:"p99CommitTime"`
	CommitsPerSecond  float64       `json:"commitsPerSecond"`
	Timestamp         time.Time     `json:"timestamp"`
}

// InvalidationMetrics provides cache invalidation metrics
type InvalidationMetrics struct {
	TotalInvalidations      uint64        `json:"totalInvalidations"`
	SuccessfulInvalidations uint64        `json:"successfulInvalidations"`
	FailedInvalidations     uint64        `json:"failedInvalidations"`
	KeysInvalidated         uint64        `json:"keysInvalidated"`
	AverageInvalidationTime time.Duration `json:"averageInvalidationTime"`
	InvalidationsPerSecond  float64       `json:"invalidationsPerSecond"`
	Timestamp               time.Time     `json:"timestamp"`
}

// BackupMetrics provides backup-related metrics
type BackupMetrics struct {
	TotalBackups       uint64        `json:"totalBackups"`
	SuccessfulBackups  uint64        `json:"successfulBackups"`
	FailedBackups      uint64        `json:"failedBackups"`
	AverageBackupTime  time.Duration `json:"averageBackupTime"`
	AverageBackupSize  int64         `json:"averageBackupSize"`
	TotalBackupStorage int64         `json:"totalBackupStorage"`
	LastBackupTime     *time.Time    `json:"lastBackupTime,omitempty"`
	RestoreOperations  uint64        `json:"restoreOperations"`
	SuccessfulRestores uint64        `json:"successfulRestores"`
	Timestamp          time.Time     `json:"timestamp"`
}

// IndexMetrics provides index-related metrics
type IndexMetrics struct {
	TotalIndexes     int           `json:"totalIndexes"`
	UsedIndexes      int           `json:"usedIndexes"`
	UnusedIndexes    int           `json:"unusedIndexes"`
	TotalIndexSize   int64         `json:"totalIndexSize"`
	IndexScans       uint64        `json:"indexScans"`
	IndexCreations   uint64        `json:"indexCreations"`
	IndexDrops       uint64        `json:"indexDrops"`
	IndexRebuildTime time.Duration `json:"indexRebuildTime"`
	LastOptimization *time.Time    `json:"lastOptimization,omitempty"`
	Timestamp        time.Time     `json:"timestamp"`
}

// ValidationMetrics provides data validation metrics
type ValidationMetrics struct {
	TotalValidations      uint64        `json:"totalValidations"`
	SuccessfulValidations uint64        `json:"successfulValidations"`
	FailedValidations     uint64        `json:"failedValidations"`
	IssuesFound           uint64        `json:"issuesFound"`
	AverageValidationTime time.Duration `json:"averageValidationTime"`
	LastValidation        *time.Time    `json:"lastValidation,omitempty"`
	BlocksValidated       uint64        `json:"blocksValidated"`
	TransactionsValidated uint64        `json:"transactionsValidated"`
	Timestamp             time.Time     `json:"timestamp"`
}

// TransactionMetrics provides database transaction metrics
type TransactionMetrics struct {
	ActiveTransactions     int           `json:"activeTransactions"`
	TotalTransactions      uint64        `json:"totalTransactions"`
	CommittedTransactions  uint64        `json:"committedTransactions"`
	RolledBackTransactions uint64        `json:"rolledBackTransactions"`
	AverageTransactionTime time.Duration `json:"averageTransactionTime"`
	LongestTransaction     time.Duration `json:"longestTransaction"`
	DeadlocksDetected      uint64        `json:"deadlocksDetected"`
	Timestamp              time.Time     `json:"timestamp"`
}

// PerformanceMetrics provides database performance metrics
type PerformanceMetrics struct {
	QueriesPerSecond      float64       `json:"queriesPerSecond"`
	AverageQueryTime      time.Duration `json:"averageQueryTime"`
	SlowQueries           uint64        `json:"slowQueries"`
	CacheHitRatio         float64       `json:"cacheHitRatio"`
	ConnectionUtilization float64       `json:"connectionUtilization"`
	DiskIOPS              float64       `json:"diskIops"`
	NetworkThroughput     float64       `json:"networkThroughput"`
	CPUUtilization        float64       `json:"cpuUtilization"`
	MemoryUtilization     float64       `json:"memoryUtilization"`
	Timestamp             time.Time     `json:"timestamp"`
}

// ArchiveMetrics defines comprehensive metrics for YaoArchive
type ArchiveMetrics struct {
	ServerMetrics       *ServerMetrics       `json:"serverMetrics"`
	PersistenceMetrics  *PersistenceMetrics  `json:"persistenceMetrics"`
	CommitMetrics       *CommitMetrics       `json:"commitMetrics"`
	InvalidationMetrics *InvalidationMetrics `json:"invalidationMetrics"`
	BackupMetrics       *BackupMetrics       `json:"backupMetrics"`
	IndexMetrics        *IndexMetrics        `json:"indexMetrics"`
	ValidationMetrics   *ValidationMetrics   `json:"validationMetrics"`
	TransactionMetrics  *TransactionMetrics  `json:"transactionMetrics"`
	PerformanceMetrics  *PerformanceMetrics  `json:"performanceMetrics"`
	MessageMetrics      *MessageMetrics      `json:"messageMetrics"`
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

// ArchiveConfig defines comprehensive configuration for YaoArchive
type ArchiveConfig struct {
	// Node configuration
	NodeID   string `json:"nodeId"`
	NodeName string `json:"nodeName"`

	// Database configuration (TiDB)
	DatabaseDSN         string `json:"databaseDsn"`
	DatabaseMaxOpenConn int    `json:"databaseMaxOpenConn"`
	DatabaseMaxIdleConn int    `json:"databaseMaxIdleConn"`
	DatabaseMaxLifetime int    `json:"databaseMaxLifetime"`

	// Message queue configuration
	RocketMQEndpoints []string `json:"rocketmqEndpoints"`
	RocketMQAccessKey string   `json:"rocketmqAccessKey"`
	RocketMQSecretKey string   `json:"rocketmqSecretKey"`
	ConsumerGroup     string   `json:"consumerGroup"`
	ProducerGroup     string   `json:"producerGroup"`

	// Processing configuration
	WorkerCount         int `json:"workerCount"`
	ProcessingTimeout   int `json:"processingTimeout"`
	MaxConcurrentBlocks int `json:"maxConcurrentBlocks"`
	BatchSize           int `json:"batchSize"`

	// Cache invalidation configuration
	EnableCacheInvalidation    bool `json:"enableCacheInvalidation"`
	CacheInvalidationBatchSize int  `json:"cacheInvalidationBatchSize"`
	CacheInvalidationTimeout   int  `json:"cacheInvalidationTimeout"`

	// Backup configuration
	EnableBackups     bool              `json:"enableBackups"`
	BackupInterval    time.Duration     `json:"backupInterval"`
	BackupRetention   int               `json:"backupRetention"`
	BackupStoragePath string            `json:"backupStoragePath"`
	BackupCompression CompressionType   `json:"backupCompression"`
	BackupEncryption  *EncryptionConfig `json:"backupEncryption,omitempty"`

	// Index management configuration
	EnableIndexOptimization   bool          `json:"enableIndexOptimization"`
	IndexOptimizationInterval time.Duration `json:"indexOptimizationInterval"`
	AutoDropUnusedIndexes     bool          `json:"autoDropUnusedIndexes"`

	// Data validation configuration
	EnableDataValidation   bool          `json:"enableDataValidation"`
	ValidationInterval     time.Duration `json:"validationInterval"`
	IntegrityCheckInterval time.Duration `json:"integrityCheckInterval"`

	// Performance optimization configuration
	EnablePerformanceOptimization bool          `json:"enablePerformanceOptimization"`
	PerformanceAnalysisInterval   time.Duration `json:"performanceAnalysisInterval"`
	QueryOptimizationThreshold    time.Duration `json:"queryOptimizationThreshold"`
	AutoTableCompaction           bool          `json:"autoTableCompaction"`

	// Monitoring configuration
	EnableMetrics       bool          `json:"enableMetrics"`
	MetricsInterval     time.Duration `json:"metricsInterval"`
	EnableHealthCheck   bool          `json:"enableHealthCheck"`
	HealthCheckInterval time.Duration `json:"healthCheckInterval"`

	// Logging configuration
	LogLevel  string `json:"logLevel"`
	LogFormat string `json:"logFormat"`
	LogFile   string `json:"logFile"`
	LogSQL    bool   `json:"logSql"`
}
