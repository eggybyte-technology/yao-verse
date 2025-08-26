package interfaces

import (
	"context"
	"math/big"

	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
)

// StateReader defines the interface for reading state data
type StateReader interface {
	// GetAccount retrieves an account by address
	GetAccount(ctx context.Context, address common.Address) (*types.Account, error)

	// GetBalance retrieves the balance of an account
	GetBalance(ctx context.Context, address common.Address) (*big.Int, error)

	// GetNonce retrieves the nonce of an account
	GetNonce(ctx context.Context, address common.Address) (uint64, error)

	// GetCode retrieves the code of a contract
	GetCode(ctx context.Context, address common.Address) ([]byte, error)

	// GetCodeHash retrieves the code hash of a contract
	GetCodeHash(ctx context.Context, address common.Address) (common.Hash, error)

	// GetStorageAt retrieves a storage value at a specific key
	GetStorageAt(ctx context.Context, address common.Address, key common.Hash) (common.Hash, error)

	// HasAccount checks if an account exists
	HasAccount(ctx context.Context, address common.Address) (bool, error)

	// IsContract checks if an address is a contract
	IsContract(ctx context.Context, address common.Address) (bool, error)
}

// StateWriter defines the interface for writing state data
type StateWriter interface {
	// SetAccount sets account data
	SetAccount(ctx context.Context, address common.Address, account *types.Account) error

	// SetBalance sets the balance of an account
	SetBalance(ctx context.Context, address common.Address, balance *big.Int) error

	// SetNonce sets the nonce of an account
	SetNonce(ctx context.Context, address common.Address, nonce uint64) error

	// SetCode sets the code of a contract
	SetCode(ctx context.Context, address common.Address, code []byte) error

	// SetStorageAt sets a storage value at a specific key
	SetStorageAt(ctx context.Context, address common.Address, key, value common.Hash) error

	// DeleteAccount deletes an account
	DeleteAccount(ctx context.Context, address common.Address) error

	// CreateAccount creates a new account
	CreateAccount(ctx context.Context, address common.Address) error
}

// StateDB defines the combined state database interface
type StateDB interface {
	StateReader
	StateWriter

	// Clone creates a copy of the state database
	Clone() StateDB

	// Copy creates a deep copy of the state database
	Copy() StateDB

	// Snapshot creates a snapshot of the current state
	Snapshot() int

	// RevertToSnapshot reverts the state to a previous snapshot
	RevertToSnapshot(revid int)

	// GetLogs gets all logs generated during transaction execution
	GetLogs(txHash common.Hash, blockHash common.Hash, blockNumber uint64) []*types.Log

	// AddLog adds a log entry
	AddLog(log *types.Log)

	// GetRefund returns the current value of the refund counter
	GetRefund() uint64

	// Finalise finalizes the state and returns the state root
	Finalise(deleteEmptyObjects bool) (common.Hash, error)

	// Commit commits the state changes to the underlying database
	Commit(deleteEmptyObjects bool) (common.Hash, error)
}

// CacheManager defines the interface for managing multi-level cache
type CacheManager interface {
	// L1Cache operations (local in-memory cache)
	GetFromL1(key string) (interface{}, bool)
	SetToL1(key string, value interface{}) error
	DeleteFromL1(key string) error
	ClearL1() error
	GetL1Stats() *CacheStats

	// L2Cache operations (Redis cluster)
	GetFromL2(ctx context.Context, key string) (interface{}, error)
	SetToL2(ctx context.Context, key string, value interface{}) error
	DeleteFromL2(ctx context.Context, key string) error
	ClearL2(ctx context.Context) error
	GetL2Stats() *CacheStats

	// Cache invalidation
	InvalidateKeys(ctx context.Context, keys []*types.CacheKey) error

	// Cache warming
	WarmupCache(ctx context.Context) error

	// Cache health check
	HealthCheck(ctx context.Context) error
}

// CacheStats represents cache statistics
type CacheStats struct {
	Hits     uint64  `json:"hits"`
	Misses   uint64  `json:"misses"`
	HitRatio float64 `json:"hitRatio"`
	Size     uint64  `json:"size"`
	Capacity uint64  `json:"capacity"`
}

// L3Storage defines the interface for persistent storage (TiDB)
type L3Storage interface {
	// Account operations
	GetAccount(ctx context.Context, address common.Address) (*types.Account, error)
	SetAccount(ctx context.Context, address common.Address, account *types.Account) error
	DeleteAccount(ctx context.Context, address common.Address) error
	HasAccount(ctx context.Context, address common.Address) (bool, error)

	// Storage operations
	GetStorageAt(ctx context.Context, address common.Address, key common.Hash) (common.Hash, error)
	SetStorageAt(ctx context.Context, address common.Address, key, value common.Hash) error

	// Code operations
	GetCode(ctx context.Context, address common.Address) ([]byte, error)
	SetCode(ctx context.Context, address common.Address, code []byte) error
	GetCodeHash(ctx context.Context, address common.Address) (common.Hash, error)

	// Block operations
	GetBlock(ctx context.Context, blockNumber uint64) (*types.Block, error)
	GetBlockByHash(ctx context.Context, blockHash common.Hash) (*types.Block, error)
	SetBlock(ctx context.Context, block *types.Block) error

	// Transaction operations
	GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, error)
	GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)

	// State root operations
	GetStateRoot(ctx context.Context, blockNumber uint64) (common.Hash, error)
	SetStateRoot(ctx context.Context, blockNumber uint64, stateRoot common.Hash) error

	// Batch operations for performance
	BeginTransaction(ctx context.Context) (Transaction, error)
}

// Transaction defines the interface for database transactions
type Transaction interface {
	// Commit commits the transaction
	Commit() error

	// Rollback rolls back the transaction
	Rollback() error

	// SetAccount sets account data in the transaction
	SetAccount(address common.Address, account *types.Account) error

	// SetStorageAt sets storage data in the transaction
	SetStorageAt(address common.Address, key, value common.Hash) error

	// SetCode sets contract code in the transaction
	SetCode(address common.Address, code []byte) error

	// SetBlock sets block data in the transaction
	SetBlock(block *types.Block) error

	// SetReceipts sets transaction receipts in the transaction
	SetReceipts(receipts []*types.Receipt) error

	// SetStateRoot sets state root in the transaction
	SetStateRoot(blockNumber uint64, stateRoot common.Hash) error
}

// YaoOracle defines the main interface for the embedded state library
type YaoOracle interface {
	StateDB
	CacheManager

	// Initialize initializes the oracle with configuration
	Initialize(ctx context.Context, config *OracleConfig) error

	// Start starts background services (cache invalidation listener, etc.)
	Start(ctx context.Context) error

	// Stop stops background services
	Stop(ctx context.Context) error

	// Close closes the oracle and releases resources
	Close() error

	// GetStorage returns the L3 storage interface
	GetStorage() L3Storage

	// ExecuteInTransaction executes a function within a database transaction
	ExecuteInTransaction(ctx context.Context, fn func(StateDB) error) error

	// GetStateDiff calculates state difference between snapshots
	GetStateDiff(beforeSnapshot, afterSnapshot int) *types.StateDiff

	// ApplyStateDiff applies state changes from a state diff
	ApplyStateDiff(ctx context.Context, stateDiff *types.StateDiff) error

	// GetCacheInvalidationHandler returns the cache invalidation message handler
	GetCacheInvalidationHandler() types.MessageHandler

	// GetMetrics returns oracle performance metrics
	GetMetrics() *OracleMetrics
}

// OracleConfig defines the configuration for YaoOracle
type OracleConfig struct {
	// L1 Cache configuration
	L1CacheSize int `json:"l1CacheSize"`
	L1CacheTTL  int `json:"l1CacheTtl"`

	// L2 Cache (Redis) configuration
	RedisEndpoints []string `json:"redisEndpoints"`
	RedisPassword  string   `json:"redisPassword"`
	RedisDB        int      `json:"redisDb"`
	RedisPoolSize  int      `json:"redisPoolSize"`
	L2CacheTTL     int      `json:"l2CacheTtl"`

	// L3 Storage (TiDB) configuration
	TiDBDSN         string `json:"tidbDsn"`
	TiDBMaxOpenConn int    `json:"tidbMaxOpenConn"`
	TiDBMaxIdleConn int    `json:"tidbMaxIdleConn"`
	TiDBMaxLifetime int    `json:"tidbMaxLifetime"`

	// Cache warming configuration
	EnableCacheWarming bool     `json:"enableCacheWarming"`
	WarmupAccounts     []string `json:"warmupAccounts"`

	// Message queue configuration for cache invalidation
	RocketMQEndpoints []string `json:"rocketmqEndpoints"`
	RocketMQAccessKey string   `json:"rocketmqAccessKey"`
	RocketMQSecretKey string   `json:"rocketmqSecretKey"`

	// Performance tuning
	BatchSize       int  `json:"batchSize"`
	EnableMetrics   bool `json:"enableMetrics"`
	MetricsInterval int  `json:"metricsInterval"`
}

// OracleMetrics defines metrics for YaoOracle performance monitoring
type OracleMetrics struct {
	// Cache metrics
	L1Stats *CacheStats `json:"l1Stats"`
	L2Stats *CacheStats `json:"l2Stats"`

	// Performance metrics
	TotalRequests       uint64 `json:"totalRequests"`
	AverageResponseTime int64  `json:"averageResponseTime"` // in nanoseconds
	L3Queries           uint64 `json:"l3Queries"`

	// Error metrics
	TotalErrors   uint64 `json:"totalErrors"`
	CacheErrors   uint64 `json:"cacheErrors"`
	StorageErrors uint64 `json:"storageErrors"`

	// System metrics
	MemoryUsage    uint64 `json:"memoryUsage"`
	GoroutineCount int    `json:"goroutineCount"`

	// Invalidation metrics
	InvalidationMessages uint64 `json:"invalidationMessages"`
	InvalidatedKeys      uint64 `json:"invalidatedKeys"`
}
