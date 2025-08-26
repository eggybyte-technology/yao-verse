package interfaces

import (
	"context"
	"math/big"

	oracleinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// JSONRPCService defines the Ethereum JSON-RPC compatible interface
type JSONRPCService interface {
	// Transaction related methods
	SendRawTransaction(ctx context.Context, data hexutil.Bytes) (common.Hash, error)
	GetTransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, error)
	GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error)
	GetTransactionCount(ctx context.Context, address common.Address, blockNumber *big.Int) (uint64, error)

	// Account related methods
	GetBalance(ctx context.Context, address common.Address, blockNumber *big.Int) (*big.Int, error)
	GetCode(ctx context.Context, address common.Address, blockNumber *big.Int) (hexutil.Bytes, error)
	GetStorageAt(ctx context.Context, address common.Address, key common.Hash, blockNumber *big.Int) (common.Hash, error)

	// Block related methods
	GetBlockByNumber(ctx context.Context, number *big.Int, fullTx bool) (*types.Block, error)
	GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (*types.Block, error)
	BlockNumber(ctx context.Context) (uint64, error)

	// Call methods
	Call(ctx context.Context, args CallArgs, blockNumber *big.Int) (hexutil.Bytes, error)
	EstimateGas(ctx context.Context, args CallArgs, blockNumber *big.Int) (uint64, error)

	// Network methods
	ChainId(ctx context.Context) (*big.Int, error)
	GasPrice(ctx context.Context) (*big.Int, error)

	// Log methods
	GetLogs(ctx context.Context, crit FilterQuery) ([]*types.Log, error)

	// Utility methods
	GetProof(ctx context.Context, address common.Address, keys []common.Hash, blockNumber *big.Int) (*AccountResult, error)
}

// CallArgs represents the arguments for a contract call
type CallArgs struct {
	From     *common.Address `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Data     *hexutil.Bytes  `json:"data"`
}

// FilterQuery represents a log filter query
type FilterQuery struct {
	BlockHash *common.Hash     `json:"blockHash"`
	FromBlock *big.Int         `json:"fromBlock"`
	ToBlock   *big.Int         `json:"toBlock"`
	Addresses []common.Address `json:"addresses"`
	Topics    [][]common.Hash  `json:"topics"`
}

// AccountResult represents the result of GetProof
type AccountResult struct {
	Address      common.Address  `json:"address"`
	AccountProof []string        `json:"accountProof"`
	Balance      *big.Int        `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        uint64          `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	StorageProof []StorageResult `json:"storageProof"`
}

// StorageResult represents the storage proof result
type StorageResult struct {
	Key   common.Hash `json:"key"`
	Value *big.Int    `json:"value"`
	Proof []string    `json:"proof"`
}

// TxPoolService defines the transaction pool interface
type TxPoolService interface {
	// SubmitTransaction submits a transaction to the pool
	SubmitTransaction(ctx context.Context, tx *types.Transaction) error

	// GetPendingTransactions returns pending transactions
	GetPendingTransactions(ctx context.Context) ([]*types.Transaction, error)

	// GetTransactionStatus returns the status of a transaction
	GetTransactionStatus(ctx context.Context, hash common.Hash) (*TxStatus, error)

	// ValidateTransaction validates a transaction before adding to pool
	ValidateTransaction(ctx context.Context, tx *types.Transaction) error

	// GetNonce returns the next nonce for an address
	GetNonce(ctx context.Context, address common.Address) (uint64, error)

	// RemoveTransaction removes a transaction from the pool
	RemoveTransaction(ctx context.Context, hash common.Hash) error
}

// TxStatus represents the status of a transaction
type TxStatus struct {
	Status      TxStatusType `json:"status"`
	BlockHash   *common.Hash `json:"blockHash,omitempty"`
	BlockNumber *uint64      `json:"blockNumber,omitempty"`
	Error       string       `json:"error,omitempty"`
}

// TxStatusType defines transaction status types
type TxStatusType string

const (
	TxStatusPending   TxStatusType = "pending"
	TxStatusConfirmed TxStatusType = "confirmed"
	TxStatusFailed    TxStatusType = "failed"
	TxStatusDropped   TxStatusType = "dropped"
)

// PortalServer defines the main YaoPortal server interface
type PortalServer interface {
	// Start starts the portal server
	Start(ctx context.Context) error

	// Stop stops the portal server gracefully
	Stop(ctx context.Context) error

	// RegisterJSONRPCService registers the JSON-RPC service
	RegisterJSONRPCService(service JSONRPCService) error

	// GetTxPoolService returns the transaction pool service
	GetTxPoolService() TxPoolService

	// GetMetrics returns server metrics
	GetMetrics() *PortalMetrics

	// HealthCheck performs a health check
	HealthCheck(ctx context.Context) error
}

// PortalMetrics defines metrics for the portal server
type PortalMetrics struct {
	// Request metrics
	TotalRequests       uint64            `json:"totalRequests"`
	RequestsPerSecond   float64           `json:"requestsPerSecond"`
	AverageResponseTime int64             `json:"averageResponseTime"` // in milliseconds
	RequestsByMethod    map[string]uint64 `json:"requestsByMethod"`

	// Transaction pool metrics
	PendingTransactions uint64 `json:"pendingTransactions"`
	SubmittedTxs        uint64 `json:"submittedTxs"`
	DroppedTxs          uint64 `json:"droppedTxs"`

	// Error metrics
	TotalErrors  uint64            `json:"totalErrors"`
	ErrorsByType map[string]uint64 `json:"errorsByType"`
	InvalidTxs   uint64            `json:"invalidTxs"`

	// Connection metrics
	ActiveConnections uint64 `json:"activeConnections"`
	TotalConnections  uint64 `json:"totalConnections"`

	// System metrics
	MemoryUsage    uint64  `json:"memoryUsage"`
	CPUUsage       float64 `json:"cpuUsage"`
	GoroutineCount int     `json:"goroutineCount"`
}

// PortalConfig defines the configuration for YaoPortal
type PortalConfig struct {
	// Server configuration
	Host string `json:"host"`
	Port int    `json:"port"`

	// TLS configuration
	EnableTLS bool   `json:"enableTls"`
	CertFile  string `json:"certFile"`
	KeyFile   string `json:"keyFile"`

	// CORS configuration
	EnableCORS     bool     `json:"enableCors"`
	AllowedOrigins []string `json:"allowedOrigins"`
	AllowedMethods []string `json:"allowedMethods"`
	AllowedHeaders []string `json:"allowedHeaders"`

	// Rate limiting
	EnableRateLimit bool `json:"enableRateLimit"`
	RateLimit       int  `json:"rateLimit"` // requests per second
	BurstLimit      int  `json:"burstLimit"`

	// Transaction pool configuration
	TxPoolSize    int `json:"txPoolSize"`
	TxPoolTimeout int `json:"txPoolTimeout"` // in seconds
	MaxTxDataSize int `json:"maxTxDataSize"` // in bytes

	// Message queue configuration
	RocketMQEndpoints []string `json:"rocketmqEndpoints"`
	RocketMQAccessKey string   `json:"rocketmqAccessKey"`
	RocketMQSecretKey string   `json:"rocketmqSecretKey"`
	ProducerGroup     string   `json:"producerGroup"`

	// YaoOracle configuration for read operations
	OracleConfig *oracleinterfaces.OracleConfig `json:"oracleConfig"`

	// Chain configuration
	ChainID     int64  `json:"chainId"`
	NetworkID   int64  `json:"networkId"`
	GenesisHash string `json:"genesisHash"`

	// Performance tuning
	MaxConcurrentRequests int  `json:"maxConcurrentRequests"`
	RequestTimeout        int  `json:"requestTimeout"` // in seconds
	EnableMetrics         bool `json:"enableMetrics"`
	MetricsInterval       int  `json:"metricsInterval"` // in seconds

	// Logging configuration
	LogLevel  string `json:"logLevel"`
	LogFormat string `json:"logFormat"`
	LogFile   string `json:"logFile"`
}
