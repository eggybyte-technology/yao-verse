package interfaces

import (
	"context"

	oracleinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
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
}

// StateCommitter defines the interface for committing state changes
type StateCommitter interface {
	// CommitState commits state changes in a transaction
	CommitState(ctx context.Context, result *types.ExecutionResult) error

	// RollbackState rolls back uncommitted state changes
	RollbackState(ctx context.Context, blockNumber uint64) error

	// GetStateRoot returns the state root for a given block
	GetStateRoot(ctx context.Context, blockNumber uint64) (common.Hash, error)
}

// CacheInvalidator defines the interface for cache invalidation
type CacheInvalidator interface {
	// InvalidateCaches invalidates relevant caches after state commit
	InvalidateCaches(ctx context.Context, stateDiff *types.StateDiff) error

	// GenerateInvalidationKeys generates cache keys that need to be invalidated
	GenerateInvalidationKeys(ctx context.Context, stateDiff *types.StateDiff) ([]*types.CacheKey, error)

	// BroadcastInvalidation broadcasts cache invalidation message
	BroadcastInvalidation(ctx context.Context, message *types.CacheInvalidationMessage) error
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

	// GetStorage returns the L3 storage interface
	GetStorage() oracleinterfaces.L3Storage

	// HealthCheck performs a comprehensive health check
	HealthCheck(ctx context.Context) error
}

// ArchiveConfig defines the configuration for YaoArchive
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

	// Cache invalidation configuration
	EnableCacheInvalidation    bool `json:"enableCacheInvalidation"`
	CacheInvalidationBatchSize int  `json:"cacheInvalidationBatchSize"`
	CacheInvalidationTimeout   int  `json:"cacheInvalidationTimeout"`

	// Logging configuration
	LogLevel  string `json:"logLevel"`
	LogFormat string `json:"logFormat"`
	LogFile   string `json:"logFile"`
}
