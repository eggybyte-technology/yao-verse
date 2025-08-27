package interfaces

import (
	"context"
	"math/big"
	"time"

	sharedinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
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
	From      *common.Address `json:"from"`
	To        *common.Address `json:"to"`
	Gas       *uint64         `json:"gas"`
	GasPrice  *big.Int        `json:"gasPrice"`
	Value     *big.Int        `json:"value"`
	Data      hexutil.Bytes   `json:"data"`
	GasFeeCap *big.Int        `json:"maxFeePerGas"`
	GasTipCap *big.Int        `json:"maxPriorityFeePerGas"`
}

// FilterQuery contains options for contract log filtering
type FilterQuery struct {
	BlockHash *common.Hash     `json:"blockHash"`
	FromBlock *big.Int         `json:"fromBlock"`
	ToBlock   *big.Int         `json:"toBlock"`
	Addresses []common.Address `json:"addresses"`
	Topics    [][]common.Hash  `json:"topics"`
}

// AccountResult represents the account proof result
type AccountResult struct {
	Address      common.Address  `json:"address"`
	Balance      *big.Int        `json:"balance"`
	CodeHash     common.Hash     `json:"codeHash"`
	Nonce        uint64          `json:"nonce"`
	StorageHash  common.Hash     `json:"storageHash"`
	AccountProof []string        `json:"accountProof"`
	StorageProof []StorageResult `json:"storageProof"`
}

// StorageResult represents the storage proof result
type StorageResult struct {
	Key   string   `json:"key"`
	Value *big.Int `json:"value"`
	Proof []string `json:"proof"`
}

// TxPoolService defines the interface for transaction pool management
type TxPoolService interface {
	// AddTransaction adds a transaction to the pool
	AddTransaction(ctx context.Context, tx *types.Transaction) error

	// RemoveTransaction removes a transaction from the pool
	RemoveTransaction(ctx context.Context, hash common.Hash) error

	// GetTransaction returns a transaction from the pool
	GetTransaction(ctx context.Context, hash common.Hash) (*types.Transaction, error)

	// GetPendingTransactions returns all pending transactions
	GetPendingTransactions(ctx context.Context) ([]*types.Transaction, error)

	// GetQueuedTransactions returns all queued transactions
	GetQueuedTransactions(ctx context.Context) ([]*types.Transaction, error)

	// GetTransactionStatus returns the status of a transaction
	GetTransactionStatus(ctx context.Context, hash common.Hash) (*TxStatus, error)

	// ValidateTransaction validates a transaction
	ValidateTransaction(ctx context.Context, tx *types.Transaction) error

	// GetPoolStatus returns the current pool status
	GetPoolStatus(ctx context.Context) (*TxPoolStatus, error)

	// SetGasPrice sets the minimum gas price
	SetGasPrice(gasPrice *big.Int) error

	// GetGasPrice returns the minimum gas price
	GetGasPrice() *big.Int

	// Clear clears all transactions from the pool
	Clear(ctx context.Context) error

	// GetMetrics returns transaction pool metrics
	GetMetrics() *TxPoolMetrics
}

// TxStatus represents the status of a transaction
type TxStatus struct {
	Hash        common.Hash  `json:"hash"`
	Status      TxStatusType `json:"status"`
	BlockNumber *uint64      `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash `json:"blockHash,omitempty"`
	TxIndex     *uint64      `json:"transactionIndex,omitempty"`
	GasUsed     *uint64      `json:"gasUsed,omitempty"`
	Error       string       `json:"error,omitempty"`
}

// TxStatusType defines transaction status types
type TxStatusType string

const (
	TxStatusPending   TxStatusType = "pending"
	TxStatusQueued    TxStatusType = "queued"
	TxStatusConfirmed TxStatusType = "confirmed"
	TxStatusFailed    TxStatusType = "failed"
	TxStatusRejected  TxStatusType = "rejected"
	TxStatusDropped   TxStatusType = "dropped"
)

// TxPoolStatus represents the transaction pool status
type TxPoolStatus struct {
	PendingCount  int      `json:"pendingCount"`
	QueuedCount   int      `json:"queuedCount"`
	Capacity      int      `json:"capacity"`
	UsagePercent  float64  `json:"usagePercent"`
	MinGasPrice   *big.Int `json:"minGasPrice"`
	MaxTxLifetime int64    `json:"maxTxLifetime"` // seconds
}

// TxPoolMetrics defines metrics for transaction pool
type TxPoolMetrics struct {
	// Transaction counts
	TotalTransactions    uint64 `json:"totalTransactions"`
	PendingTransactions  int    `json:"pendingTransactions"`
	QueuedTransactions   int    `json:"queuedTransactions"`
	DroppedTransactions  uint64 `json:"droppedTransactions"`
	RejectedTransactions uint64 `json:"rejectedTransactions"`

	// Performance metrics
	TransactionsPerSecond float64 `json:"transactionsPerSecond"`
	AverageProcessingTime int64   `json:"averageProcessingTime"` // nanoseconds
	PoolUtilization       float64 `json:"poolUtilization"`       // percentage

	// Gas metrics
	AverageGasPrice *big.Int `json:"averageGasPrice"`
	MinGasPrice     *big.Int `json:"minGasPrice"`
	MaxGasPrice     *big.Int `json:"maxGasPrice"`

	// Error metrics
	ValidationErrors        uint64 `json:"validationErrors"`
	InsufficientFundsErrors uint64 `json:"insufficientFundsErrors"`
	NonceErrors             uint64 `json:"nonceErrors"`
	GasPriceErrors          uint64 `json:"gasPriceErrors"`
}

// PortalServer defines the main YaoPortal server interface
type PortalServer interface {
	// Start starts the Portal server
	Start(ctx context.Context) error

	// Stop stops the Portal server gracefully
	Stop(ctx context.Context) error

	// GetJSONRPCService returns the JSON-RPC service instance
	GetJSONRPCService() JSONRPCService

	// GetTxPoolService returns the transaction pool service instance
	GetTxPoolService() TxPoolService

	// GetMessageProducer returns the message producer for publishing transactions
	GetMessageProducer() sharedinterfaces.MessageProducer

	// GetStateReader returns the state reader for queries
	GetStateReader() sharedinterfaces.StateReader

	// RegisterWebSocketClient registers a WebSocket client for subscriptions
	RegisterWebSocketClient(ctx context.Context, clientID string, subscriptions []string) error

	// UnregisterWebSocketClient unregisters a WebSocket client
	UnregisterWebSocketClient(ctx context.Context, clientID string) error

	// BroadcastToClients broadcasts a message to all subscribed clients
	BroadcastToClients(ctx context.Context, subscription string, data interface{}) error

	// GetStatus returns the current server status
	GetStatus(ctx context.Context) *PortalStatus

	// GetMetrics returns comprehensive server metrics
	GetMetrics(ctx context.Context) *PortalMetrics

	// HealthCheck performs comprehensive health check
	HealthCheck(ctx context.Context) error
}

// PortalStatus represents the status of YaoPortal server
type PortalStatus struct {
	ServerStatus        ServerStatus  `json:"serverStatus"`
	ConnectedClients    int           `json:"connectedClients"`
	ActiveSubscriptions int           `json:"activeSubscriptions"`
	RequestsPerSecond   float64       `json:"requestsPerSecond"`
	TxPoolStatus        *TxPoolStatus `json:"txPoolStatus"`
	Uptime              time.Duration `json:"uptime"`
	StartTime           time.Time     `json:"startTime"`
	Version             string        `json:"version"`
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

// PortalMetrics defines comprehensive metrics for YaoPortal
type PortalMetrics struct {
	ServerMetrics    *ServerMetrics    `json:"serverMetrics"`
	JSONRPCMetrics   *JSONRPCMetrics   `json:"jsonrpcMetrics"`
	TxPoolMetrics    *TxPoolMetrics    `json:"txPoolMetrics"`
	WebSocketMetrics *WebSocketMetrics `json:"websocketMetrics"`
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

// JSONRPCMetrics defines metrics for JSON-RPC service
type JSONRPCMetrics struct {
	// Request metrics
	TotalRequests     uint64            `json:"totalRequests"`
	RequestsPerSecond float64           `json:"requestsPerSecond"`
	RequestsByMethod  map[string]uint64 `json:"requestsByMethod"`

	// Response metrics
	SuccessfulResponses uint64            `json:"successfulResponses"`
	ErrorResponses      uint64            `json:"errorResponses"`
	ErrorsByType        map[string]uint64 `json:"errorsByType"`

	// Performance metrics
	AverageResponseTime time.Duration `json:"averageResponseTime"`
	P99ResponseTime     time.Duration `json:"p99ResponseTime"`

	// Connection metrics
	ActiveConnections  int    `json:"activeConnections"`
	TotalConnections   uint64 `json:"totalConnections"`
	ConnectionsDropped uint64 `json:"connectionsDropped"`

	Timestamp time.Time `json:"timestamp"`
}

// WebSocketMetrics defines metrics for WebSocket connections
type WebSocketMetrics struct {
	// Connection metrics
	ActiveConnections  int    `json:"activeConnections"`
	TotalConnections   uint64 `json:"totalConnections"`
	ConnectionsDropped uint64 `json:"connectionsDropped"`

	// Subscription metrics
	ActiveSubscriptions int            `json:"activeSubscriptions"`
	SubscriptionsByType map[string]int `json:"subscriptionsByType"`

	// Message metrics
	MessagesSent     uint64 `json:"messagesSent"`
	MessagesReceived uint64 `json:"messagesReceived"`
	BroadcastsSent   uint64 `json:"broadcastsSent"`

	// Performance metrics
	AverageLatency   time.Duration `json:"averageLatency"`
	ConnectionUptime time.Duration `json:"connectionUptime"`

	Timestamp time.Time `json:"timestamp"`
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

// PortalConfig defines comprehensive configuration for YaoPortal
type PortalConfig struct {
	// Node configuration
	NodeID   string `json:"nodeId"`
	NodeName string `json:"nodeName"`

	// Server configuration
	HTTPBindAddress string `json:"httpBindAddress"`
	HTTPPort        int    `json:"httpPort"`
	WSBindAddress   string `json:"wsBindAddress"`
	WSPort          int    `json:"wsPort"`
	EnableTLS       bool   `json:"enableTls"`
	TLSCert         string `json:"tlsCert,omitempty"`
	TLSKey          string `json:"tlsKey,omitempty"`

	// JSON-RPC configuration
	EnableHTTP     bool          `json:"enableHttp"`
	EnableWS       bool          `json:"enableWs"`
	CORSOrigins    []string      `json:"corsOrigins"`
	MaxRequestSize int64         `json:"maxRequestSize"`
	RequestTimeout time.Duration `json:"requestTimeout"`
	BatchLimit     int           `json:"batchLimit"`

	// Transaction pool configuration
	TxPoolCapacity int           `json:"txPoolCapacity"`
	TxPoolTimeout  time.Duration `json:"txPoolTimeout"`
	MinGasPrice    *big.Int      `json:"minGasPrice,omitempty"`
	MaxTxLifetime  time.Duration `json:"maxTxLifetime"`

	// WebSocket configuration
	WSMaxConnections int           `json:"wsMaxConnections"`
	WSPingInterval   time.Duration `json:"wsPingInterval"`
	WSWriteTimeout   time.Duration `json:"wsWriteTimeout"`
	WSReadTimeout    time.Duration `json:"wsReadTimeout"`

	// State reading configuration
	StateReadTimeout time.Duration `json:"stateReadTimeout"`
	MaxBlockRange    uint64        `json:"maxBlockRange"`
	EnableTracing    bool          `json:"enableTracing"`

	// Message queue configuration
	RocketMQEndpoints []string `json:"rocketmqEndpoints"`
	RocketMQAccessKey string   `json:"rocketmqAccessKey"`
	RocketMQSecretKey string   `json:"rocketmqSecretKey"`
	ProducerGroup     string   `json:"producerGroup"`

	// Rate limiting configuration
	EnableRateLimit   bool    `json:"enableRateLimit"`
	RequestsPerSecond float64 `json:"requestsPerSecond"`
	BurstSize         int     `json:"burstSize"`

	// Authentication configuration
	EnableAuth     bool     `json:"enableAuth"`
	AllowedOrigins []string `json:"allowedOrigins"`
	APIKeys        []string `json:"apiKeys,omitempty"`

	// Monitoring configuration
	EnableMetrics       bool          `json:"enableMetrics"`
	MetricsInterval     time.Duration `json:"metricsInterval"`
	EnableHealthCheck   bool          `json:"enableHealthCheck"`
	HealthCheckInterval time.Duration `json:"healthCheckInterval"`

	// Logging configuration
	LogLevel     string `json:"logLevel"`
	LogFormat    string `json:"logFormat"`
	LogFile      string `json:"logFile"`
	LogRequests  bool   `json:"logRequests"`
	LogResponses bool   `json:"logResponses"`

	// Backend services configuration
	BackendServices *BackendServicesConfig `json:"backendServices"`

	// Ethereum configuration
	Ethereum *EthereumConfig `json:"ethereum"`

	// Additional configurations
	ConsumerGroup       string        `json:"consumerGroup,omitempty"`
	RateLimitWindowSize time.Duration `json:"rateLimitWindowSize,omitempty"`
	MetricsPort         int           `json:"metricsPort,omitempty"`
	EnableAccessLog     bool          `json:"enableAccessLog,omitempty"`
}

// BackendServicesConfig defines configuration for backend services
type BackendServicesConfig struct {
	ConsensusService *BackendServiceConfig `json:"consensusService"`
	ExecutionService *BackendServiceConfig `json:"executionService"`
	StateService     *BackendServiceConfig `json:"stateService"`
	ArchiveService   *BackendServiceConfig `json:"archiveService"`
}

// BackendServiceConfig defines configuration for a backend service
type BackendServiceConfig struct {
	Name            string   `json:"name"`
	Namespace       string   `json:"namespace"`
	Endpoints       []string `json:"endpoints"`
	HealthCheckPath string   `json:"healthCheckPath"`
}

// EthereumConfig defines Ethereum compatibility configuration
type EthereumConfig struct {
	NetworkName     string        `json:"networkName"`
	ChainId         int64         `json:"chainId"`
	DefaultGasPrice string        `json:"defaultGasPrice"`
	DefaultGasLimit uint64        `json:"defaultGasLimit"`
	BlockTime       time.Duration `json:"blockTime"`
	MaxBlockSize    uint64        `json:"maxBlockSize"`
}

// RateLimiter defines the interface for request rate limiting
type RateLimiter interface {
	// AllowRequest checks if a request should be allowed based on rate limits
	AllowRequest(ctx context.Context, clientIP string, endpoint string) (bool, error)

	// GetRateLimit returns the current rate limit settings for a client
	GetRateLimit(ctx context.Context, clientIP string) (*RateLimitStatus, error)

	// SetRateLimit sets a custom rate limit for a client
	SetRateLimit(ctx context.Context, clientIP string, limit *RateLimitConfig) error

	// ResetRateLimit resets rate limit counters for a client
	ResetRateLimit(ctx context.Context, clientIP string) error

	// GetGlobalRateStats returns global rate limiting statistics
	GetGlobalRateStats(ctx context.Context) (*RateLimitStats, error)
}

// LoadBalancer defines the interface for load balancing across backend services
type LoadBalancer interface {
	// SelectBackend selects an appropriate backend service for a request
	SelectBackend(ctx context.Context, request *BackendRequest) (*BackendInfo, error)

	// UpdateBackendHealth updates the health status of a backend
	UpdateBackendHealth(ctx context.Context, backendID string, isHealthy bool) error

	// GetBackendStatus returns the status of all backends
	GetBackendStatus(ctx context.Context) ([]*BackendStatus, error)

	// AddBackend adds a new backend to the load balancer
	AddBackend(ctx context.Context, backend *BackendConfig) error

	// RemoveBackend removes a backend from the load balancer
	RemoveBackend(ctx context.Context, backendID string) error

	// GetLoadBalancingStats returns load balancing statistics
	GetLoadBalancingStats(ctx context.Context) (*LoadBalancingStats, error)
}

// AuthenticationService defines the interface for request authentication
type AuthenticationService interface {
	// ValidateAPIKey validates an API key
	ValidateAPIKey(ctx context.Context, apiKey string) (*APIKeyInfo, error)

	// ValidateJWTToken validates a JWT token
	ValidateJWTToken(ctx context.Context, token string) (*TokenClaims, error)

	// ValidateOrigin validates the request origin
	ValidateOrigin(ctx context.Context, origin string) (bool, error)

	// CreateAPIKey creates a new API key
	CreateAPIKey(ctx context.Context, request *APIKeyRequest) (*APIKeyResponse, error)

	// RevokeAPIKey revokes an existing API key
	RevokeAPIKey(ctx context.Context, keyID string) error

	// GetAuthMetrics returns authentication metrics
	GetAuthMetrics(ctx context.Context) (*AuthMetrics, error)
}

// ConnectionManager defines the interface for managing client connections
type ConnectionManager interface {
	// RegisterConnection registers a new client connection
	RegisterConnection(ctx context.Context, conn *ClientConnection) error

	// UnregisterConnection unregisters a client connection
	UnregisterConnection(ctx context.Context, connectionID string) error

	// BroadcastMessage broadcasts a message to all connected clients
	BroadcastMessage(ctx context.Context, message *BroadcastMessage) error

	// SendToConnection sends a message to a specific connection
	SendToConnection(ctx context.Context, connectionID string, message interface{}) error

	// GetActiveConnections returns all active connections
	GetActiveConnections(ctx context.Context) ([]*ClientConnection, error)

	// CleanupStaleConnections removes stale/inactive connections
	CleanupStaleConnections(ctx context.Context) (int, error)

	// GetConnectionMetrics returns connection metrics
	GetConnectionMetrics(ctx context.Context) (*ConnectionMetrics, error)
}

// RequestValidator defines the interface for validating incoming requests
type RequestValidator interface {
	// ValidateJSONRPCRequest validates a JSON-RPC request structure
	ValidateJSONRPCRequest(ctx context.Context, request *JSONRPCRequest) error

	// ValidateTransactionRequest validates transaction parameters
	ValidateTransactionRequest(ctx context.Context, tx *types.Transaction) (*ValidationResult, error)

	// ValidateCallRequest validates eth_call parameters
	ValidateCallRequest(ctx context.Context, args *CallArgs) error

	// ValidateBlockRequest validates block query parameters
	ValidateBlockRequest(ctx context.Context, blockID string) error

	// GetValidationStats returns validation statistics
	GetValidationStats(ctx context.Context) (*ValidationStats, error)
}

// BackendServiceManager defines the interface for managing backend services
type BackendServiceManager interface {
	// GetService returns the endpoint for a specific service type
	GetService(ctx context.Context, serviceType string) (string, error)

	// GetHealthyService returns a healthy endpoint for a service type
	GetHealthyService(ctx context.Context, serviceType string) (string, error)

	// CheckServiceHealth checks the health of a service
	CheckServiceHealth(ctx context.Context, serviceType string, endpoint string) (bool, error)

	// GetAllServiceEndpoints returns all endpoints for a service type
	GetAllServiceEndpoints(ctx context.Context, serviceType string) ([]string, error)

	// UpdateServiceHealth updates the health status of a service endpoint
	UpdateServiceHealth(ctx context.Context, serviceType string, endpoint string, isHealthy bool) error

	// GetServiceMetrics returns metrics for backend services
	GetServiceMetrics(ctx context.Context) (*BackendServiceMetrics, error)
}

// BackendServiceMetrics provides metrics for backend services
type BackendServiceMetrics struct {
	ServiceHealth        map[string]bool          `json:"serviceHealth"`
	ServiceRequestCounts map[string]uint64        `json:"serviceRequestCounts"`
	ServiceResponseTimes map[string]time.Duration `json:"serviceResponseTimes"`
	ServiceErrorRates    map[string]float64       `json:"serviceErrorRates"`
	LastHealthCheckTime  map[string]time.Time     `json:"lastHealthCheckTime"`
	Timestamp            time.Time                `json:"timestamp"`
}

// WebSocketService defines the interface for WebSocket functionality
type WebSocketService interface {
	// HandleWebSocketConnection handles a new WebSocket connection
	HandleWebSocketConnection(ctx context.Context, conn interface{}) error

	// Subscribe subscribes a connection to a topic
	Subscribe(ctx context.Context, connectionID string, subscription string) error

	// Unsubscribe unsubscribes a connection from a topic
	Unsubscribe(ctx context.Context, connectionID string, subscription string) error

	// Broadcast broadcasts a message to all subscribers of a topic
	Broadcast(ctx context.Context, topic string, data interface{}) error

	// GetSubscriptions returns all subscriptions for a connection
	GetSubscriptions(ctx context.Context, connectionID string) ([]string, error)

	// GetWebSocketMetrics returns WebSocket metrics
	GetWebSocketMetrics(ctx context.Context) (*WebSocketMetrics, error)
}

// Supporting types for new interfaces

// RateLimitStatus represents the current rate limit status for a client
type RateLimitStatus struct {
	ClientIP       string        `json:"clientIP"`
	RequestCount   int           `json:"requestCount"`
	WindowStart    time.Time     `json:"windowStart"`
	WindowDuration time.Duration `json:"windowDuration"`
	Limit          int           `json:"limit"`
	Remaining      int           `json:"remaining"`
	ResetTime      time.Time     `json:"resetTime"`
	IsBlocked      bool          `json:"isBlocked"`
	BlockedUntil   *time.Time    `json:"blockedUntil,omitempty"`
}

// RateLimitConfig defines rate limit configuration for a client
type RateLimitConfig struct {
	RequestsPerSecond float64       `json:"requestsPerSecond"`
	BurstSize         int           `json:"burstSize"`
	WindowDuration    time.Duration `json:"windowDuration"`
	BlockDuration     time.Duration `json:"blockDuration,omitempty"`
}

// RateLimitStats provides global rate limiting statistics
type RateLimitStats struct {
	TotalRequests      uint64            `json:"totalRequests"`
	AllowedRequests    uint64            `json:"allowedRequests"`
	RejectedRequests   uint64            `json:"rejectedRequests"`
	ActiveClients      int               `json:"activeClients"`
	BlockedClients     int               `json:"blockedClients"`
	RequestsByEndpoint map[string]uint64 `json:"requestsByEndpoint"`
	AverageLatency     time.Duration     `json:"averageLatency"`
	Timestamp          time.Time         `json:"timestamp"`
}

// BackendRequest represents a request to be load balanced
type BackendRequest struct {
	Method     string                 `json:"method"`
	Parameters map[string]interface{} `json:"parameters"`
	ClientIP   string                 `json:"clientIP"`
	RequestID  string                 `json:"requestId"`
	Priority   RequestPriority        `json:"priority"`
}

// RequestPriority defines request priority levels
type RequestPriority string

const (
	PriorityLow      RequestPriority = "low"
	PriorityNormal   RequestPriority = "normal"
	PriorityHigh     RequestPriority = "high"
	PriorityCritical RequestPriority = "critical"
)

// BackendInfo represents information about a selected backend
type BackendInfo struct {
	BackendID string            `json:"backendId"`
	Address   string            `json:"address"`
	Port      int               `json:"port"`
	Weight    int               `json:"weight"`
	IsHealthy bool              `json:"isHealthy"`
	LoadScore float64           `json:"loadScore"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// BackendStatus represents the status of a backend service
type BackendStatus struct {
	BackendID       string        `json:"backendId"`
	Address         string        `json:"address"`
	IsHealthy       bool          `json:"isHealthy"`
	ResponseTime    time.Duration `json:"responseTime"`
	RequestCount    uint64        `json:"requestCount"`
	ErrorCount      uint64        `json:"errorCount"`
	LastHealthCheck time.Time     `json:"lastHealthCheck"`
	LoadScore       float64       `json:"loadScore"`
	Weight          int           `json:"weight"`
}

// BackendConfig defines configuration for a backend service
type BackendConfig struct {
	BackendID           string            `json:"backendId"`
	Address             string            `json:"address"`
	Port                int               `json:"port"`
	Weight              int               `json:"weight"`
	HealthCheckURL      string            `json:"healthCheckUrl"`
	HealthCheckInterval time.Duration     `json:"healthCheckInterval"`
	MaxConnections      int               `json:"maxConnections"`
	Timeout             time.Duration     `json:"timeout"`
	Metadata            map[string]string `json:"metadata,omitempty"`
}

// LoadBalancingStats provides load balancing statistics
type LoadBalancingStats struct {
	TotalRequests       uint64                   `json:"totalRequests"`
	RequestsByBackend   map[string]uint64        `json:"requestsByBackend"`
	BackendHealthStatus map[string]bool          `json:"backendHealthStatus"`
	AverageResponseTime map[string]time.Duration `json:"averageResponseTime"`
	BackendLoadScores   map[string]float64       `json:"backendLoadScores"`
	FailoverEvents      uint64                   `json:"failoverEvents"`
	Timestamp           time.Time                `json:"timestamp"`
}

// APIKeyInfo represents information about an API key
type APIKeyInfo struct {
	KeyID       string            `json:"keyId"`
	ClientName  string            `json:"clientName"`
	Permissions []string          `json:"permissions"`
	RateLimit   *RateLimitConfig  `json:"rateLimit,omitempty"`
	ExpiresAt   *time.Time        `json:"expiresAt,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	LastUsed    *time.Time        `json:"lastUsed,omitempty"`
	IsActive    bool              `json:"isActive"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	Subject     string            `json:"sub"`
	Issuer      string            `json:"iss"`
	Audience    string            `json:"aud"`
	ExpiresAt   time.Time         `json:"exp"`
	IssuedAt    time.Time         `json:"iat"`
	Permissions []string          `json:"permissions"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// APIKeyRequest represents a request to create an API key
type APIKeyRequest struct {
	ClientName  string            `json:"clientName"`
	Permissions []string          `json:"permissions"`
	RateLimit   *RateLimitConfig  `json:"rateLimit,omitempty"`
	ExpiresIn   *time.Duration    `json:"expiresIn,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// APIKeyResponse represents the response when creating an API key
type APIKeyResponse struct {
	KeyID     string     `json:"keyId"`
	APIKey    string     `json:"apiKey"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

// AuthMetrics provides authentication metrics
type AuthMetrics struct {
	TotalAuthRequests     uint64            `json:"totalAuthRequests"`
	SuccessfulAuth        uint64            `json:"successfulAuth"`
	FailedAuth            uint64            `json:"failedAuth"`
	InvalidAPIKeys        uint64            `json:"invalidApiKeys"`
	ExpiredTokens         uint64            `json:"expiredTokens"`
	ActiveAPIKeys         int               `json:"activeApiKeys"`
	AuthenticationMethods map[string]uint64 `json:"authenticationMethods"`
	AverageAuthTime       time.Duration     `json:"averageAuthTime"`
	Timestamp             time.Time         `json:"timestamp"`
}

// ClientConnection represents a client connection
type ClientConnection struct {
	ConnectionID   string            `json:"connectionId"`
	ClientIP       string            `json:"clientIP"`
	UserAgent      string            `json:"userAgent"`
	ConnectedAt    time.Time         `json:"connectedAt"`
	LastActivity   time.Time         `json:"lastActivity"`
	ConnectionType ConnectionType    `json:"connectionType"`
	Subscriptions  []string          `json:"subscriptions"`
	IsActive       bool              `json:"isActive"`
	BytesSent      uint64            `json:"bytesSent"`
	BytesReceived  uint64            `json:"bytesReceived"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

// ConnectionType defines the type of connection
type ConnectionType string

const (
	ConnectionTypeHTTP      ConnectionType = "http"
	ConnectionTypeWebSocket ConnectionType = "websocket"
	ConnectionTypeGRPC      ConnectionType = "grpc"
)

// BroadcastMessage represents a message to be broadcast
type BroadcastMessage struct {
	Type           string         `json:"type"`
	Topic          string         `json:"topic"`
	Data           interface{}    `json:"data"`
	TargetClients  []string       `json:"targetClients,omitempty"`
	ExcludeClients []string       `json:"excludeClients,omitempty"`
	TTL            *time.Duration `json:"ttl,omitempty"`
}

// ConnectionMetrics provides connection metrics
type ConnectionMetrics struct {
	TotalConnections    uint64                 `json:"totalConnections"`
	ActiveConnections   int                    `json:"activeConnections"`
	ConnectionsByType   map[ConnectionType]int `json:"connectionsByType"`
	AverageSessionTime  time.Duration          `json:"averageSessionTime"`
	ConnectionsDropped  uint64                 `json:"connectionsDropped"`
	ReconnectionCount   uint64                 `json:"reconnectionCount"`
	BandwidthUsage      *BandwidthUsage        `json:"bandwidthUsage"`
	ConnectionsByRegion map[string]int         `json:"connectionsByRegion,omitempty"`
	Timestamp           time.Time              `json:"timestamp"`
}

// BandwidthUsage represents bandwidth usage statistics
type BandwidthUsage struct {
	BytesSentPerSecond     float64 `json:"bytesSentPerSecond"`
	BytesReceivedPerSecond float64 `json:"bytesReceivedPerSecond"`
	TotalBytesSent         uint64  `json:"totalBytesSent"`
	TotalBytesReceived     uint64  `json:"totalBytesReceived"`
}

// JSONRPCRequest represents a JSON-RPC request structure
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id"`
}

// ValidationResult represents the result of request validation
type ValidationResult struct {
	IsValid        bool                   `json:"isValid"`
	Errors         []string               `json:"errors,omitempty"`
	Warnings       []string               `json:"warnings,omitempty"`
	GasEstimate    *uint64                `json:"gasEstimate,omitempty"`
	Suggestions    []string               `json:"suggestions,omitempty"`
	ValidationTime time.Duration          `json:"validationTime"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// ValidationStats provides validation statistics
type ValidationStats struct {
	TotalValidations      uint64            `json:"totalValidations"`
	SuccessfulValidations uint64            `json:"successfulValidations"`
	FailedValidations     uint64            `json:"failedValidations"`
	ValidationsByType     map[string]uint64 `json:"validationsByType"`
	CommonErrors          map[string]uint64 `json:"commonErrors"`
	AverageValidationTime time.Duration     `json:"averageValidationTime"`
	Timestamp             time.Time         `json:"timestamp"`
}
