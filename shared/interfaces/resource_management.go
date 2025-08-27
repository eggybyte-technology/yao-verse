package interfaces

import (
	"context"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/types"
)

// ResourceManager defines the interface for comprehensive resource management
type ResourceManager interface {
	// Memory management
	GetMemoryManager() MemoryManager

	// Connection pool management
	GetConnectionManager() ConnectionManager

	// Thread pool management
	GetThreadPoolManager() ThreadPoolManager

	// Resource limiting
	GetResourceLimiter() ResourceLimiter

	// Resource monitoring
	GetResourceStats() *ResourceStats

	// Resource cleanup
	StartResourceCleanup(ctx context.Context) error
	StopResourceCleanup() error

	// Resource health
	HealthCheck(ctx context.Context) error
}

// MemoryManager defines the interface for memory management
type MemoryManager interface {
	// Memory allocation
	Allocate(size int64, category string) (*MemoryBlock, error)
	Deallocate(block *MemoryBlock) error
	Reallocate(block *MemoryBlock, newSize int64) (*MemoryBlock, error)

	// Memory pools
	CreateMemoryPool(name string, config *MemoryPoolConfig) (MemoryPool, error)
	GetMemoryPool(name string) (MemoryPool, error)
	DestroyMemoryPool(name string) error
	ListMemoryPools() ([]string, error)

	// Memory monitoring
	GetMemoryStats() *MemoryStats
	GetMemoryUsage() *types.MemoryUsage
	GetMemoryPressure() MemoryPressureLevel

	// Memory optimization
	TriggerGC() error
	CompactMemory() error
	GetGCStats() *GCStats

	// Memory limits
	SetMemoryLimit(limit int64) error
	GetMemoryLimit() int64
	GetAvailableMemory() int64

	// Memory alerts
	SetMemoryThreshold(threshold MemoryThreshold) error
	GetMemoryAlerts() ([]*MemoryAlert, error)

	// Memory profiling
	StartMemoryProfiling(duration time.Duration) (*MemoryProfile, error)
	GetMemoryProfile() (*MemoryProfile, error)
}

// MemoryPool defines the interface for memory pools
type MemoryPool interface {
	// Pool operations
	Get() (*MemoryBlock, error)
	Put(block *MemoryBlock) error

	// Pool management
	Resize(newCapacity int) error
	Clear() error
	Close() error

	// Pool statistics
	GetStats() *MemoryPoolStats
	GetCapacity() int
	GetUsage() int
	GetAvailable() int
}

// ConnectionManager defines the interface for connection management
type ConnectionManager interface {
	// Connection pools
	CreateConnectionPool(name string, config *ConnectionPoolConfig) (ConnectionPool, error)
	GetConnectionPool(name string) (ConnectionPool, error)
	DestroyConnectionPool(name string) error
	ListConnectionPools() ([]string, error)

	// Connection monitoring
	GetConnectionStats() *ConnectionStats
	GetActiveConnections() ([]*Connection, error)

	// Connection limits
	SetConnectionLimit(limit int) error
	GetConnectionLimit() int
	GetAvailableConnections() int

	// Connection health
	CheckConnectionHealth(connID string) error
	GetUnhealthyConnections() ([]*Connection, error)

	// Connection cleanup
	CloseIdleConnections(maxIdleTime time.Duration) error
	CloseExpiredConnections() error
}

// ConnectionPool defines the interface for connection pools
type ConnectionPool interface {
	// Connection operations
	GetConnection(ctx context.Context) (Connection, error)
	ReturnConnection(conn Connection) error

	// Pool management
	Resize(newSize int) error
	Close() error

	// Pool statistics
	GetStats() *ConnectionPoolStats
	GetActiveCount() int
	GetIdleCount() int
	GetTotalCount() int

	// Pool health
	HealthCheck(ctx context.Context) error
	ValidateConnections(ctx context.Context) error
}

// Connection represents a managed connection
type Connection interface {
	// Connection info
	GetID() string
	GetType() string
	GetRemoteAddr() string
	GetLocalAddr() string

	// Connection state
	IsActive() bool
	IsIdle() bool
	GetCreatedAt() time.Time
	GetLastUsed() time.Time

	// Connection operations
	Reset() error
	Close() error
	HealthCheck(ctx context.Context) error

	// Connection metadata
	GetMetadata() map[string]interface{}
	SetMetadata(key string, value interface{})

	// Connection metrics
	GetMetrics() *ConnectionMetrics
}

// ThreadPoolManager defines the interface for thread pool management
type ThreadPoolManager interface {
	// Thread pools
	CreateThreadPool(name string, config *ThreadPoolConfig) (ThreadPool, error)
	GetThreadPool(name string) (ThreadPool, error)
	DestroyThreadPool(name string) error
	ListThreadPools() ([]string, error)

	// Task execution
	Submit(poolName string, task Task) error
	SubmitWithTimeout(poolName string, task Task, timeout time.Duration) error
	SubmitWithPriority(poolName string, task Task, priority types.Priority) error

	// Thread pool monitoring
	GetThreadPoolStats() *ThreadPoolStats
	GetActiveThreads() int
	GetQueuedTasks() int

	// Thread pool management
	PauseThreadPool(name string) error
	ResumeThreadPool(name string) error
	ShutdownThreadPool(name string, timeout time.Duration) error
}

// ThreadPool defines the interface for individual thread pools
type ThreadPool interface {
	// Task submission
	Submit(task Task) error
	SubmitWithTimeout(task Task, timeout time.Duration) error
	SubmitWithPriority(task Task, priority Priority) error

	// Pool control
	Pause() error
	Resume() error
	Shutdown(timeout time.Duration) error

	// Pool configuration
	Resize(coreSize, maxSize int) error
	SetQueueCapacity(capacity int) error
	SetKeepAliveTime(keepAlive time.Duration) error

	// Pool statistics
	GetStats() *ThreadPoolStats
	GetCoreSize() int
	GetMaxSize() int
	GetActiveCount() int
	GetQueueSize() int

	// Pool health
	HealthCheck(ctx context.Context) error
}

// Task represents a task to be executed
type Task interface {
	Execute(ctx context.Context) error
	GetID() string
	GetPriority() Priority
	GetCreatedAt() time.Time
	GetTimeout() time.Duration
	GetMetadata() map[string]interface{}
}

// ResourceLimiter defines the interface for resource limiting
type ResourceLimiter interface {
	// Rate limiting
	Allow(resource string, tokens int) bool
	AllowN(resource string, tokens int, n int) bool
	Reserve(resource string, tokens int) Reservation
	Wait(ctx context.Context, resource string, tokens int) error

	// Quota management
	SetQuota(resource string, quota Quota) error
	GetQuota(resource string) (*Quota, error)
	GetQuotaUsage(resource string) (*QuotaUsage, error)
	ResetQuota(resource string) error

	// Resource limits
	SetResourceLimit(resource string, limit ResourceLimit) error
	GetResourceLimit(resource string) (*ResourceLimit, error)
	CheckResourceLimit(resource string, usage int64) error

	// Limiting statistics
	GetLimitingStats() *LimitingStats
	GetRateLimitStats(resource string) (*RateLimitStats, error)
}

// Reservation represents a rate limit reservation
type Reservation interface {
	OK() bool
	Delay() time.Duration
	Cancel()
	CancelAt(time.Time)
}

// CacheManager defines the interface for cache resource management
type CacheManager interface {
	// Cache creation and management
	CreateCache(name string, config *CacheConfig) (Cache, error)
	GetCache(name string) (Cache, error)
	DestroyCache(name string) error
	ListCaches() ([]string, error)

	// Global cache operations
	ClearAllCaches() error
	CompactAllCaches() error

	// Cache monitoring
	GetCacheStats() *CacheStats
	GetGlobalCacheUsage() *CacheUsage

	// Cache policies
	SetGlobalCachePolicy(policy CachePolicy) error
	GetGlobalCachePolicy() *CachePolicy
}

// Cache defines the interface for individual caches
type Cache interface {
	// Basic operations
	Get(key string) (interface{}, bool)
	Put(key string, value interface{}) error
	PutWithTTL(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Clear() error

	// Batch operations
	GetMulti(keys []string) (map[string]interface{}, error)
	PutMulti(items map[string]interface{}) error
	DeleteMulti(keys []string) error

	// Cache management
	Compact() error
	GetStats() *CacheStats
	GetSize() int
	GetCapacity() int

	// Cache policies
	SetEvictionPolicy(policy EvictionPolicy) error
	SetTTL(ttl time.Duration) error
}

// Data structures

// MemoryBlock represents an allocated memory block
type MemoryBlock struct {
	ID          string    `json:"id"`
	Size        int64     `json:"size"`
	Category    string    `json:"category"`
	Address     uintptr   `json:"address"`
	AllocatedAt time.Time `json:"allocatedAt"`
}

// MemoryPoolConfig defines configuration for memory pools
type MemoryPoolConfig struct {
	BlockSize    int64         `json:"blockSize"`
	InitialSize  int           `json:"initialSize"`
	MaxSize      int           `json:"maxSize"`
	GrowthFactor float64       `json:"growthFactor"`
	ShrinkFactor float64       `json:"shrinkFactor"`
	MaxIdleTime  time.Duration `json:"maxIdleTime"`
	EnableStats  bool          `json:"enableStats"`
}

// MemoryStats represents memory statistics
type MemoryStats struct {
	TotalAllocated    int64                       `json:"totalAllocated"`
	TotalDeallocated  int64                       `json:"totalDeallocated"`
	CurrentUsage      int64                       `json:"currentUsage"`
	PeakUsage         int64                       `json:"peakUsage"`
	AllocationCount   uint64                      `json:"allocationCount"`
	DeallocationCount uint64                      `json:"deallocationCount"`
	UsageByCategory   map[string]int64            `json:"usageByCategory"`
	PoolStats         map[string]*MemoryPoolStats `json:"poolStats"`
	GCStats           *GCStats                    `json:"gcStats"`
	Timestamp         time.Time                   `json:"timestamp"`
}

// MemoryUsage represents current memory usage
type MemoryUsage struct {
	Used       int64      `json:"used"`
	Available  int64      `json:"available"`
	Total      int64      `json:"total"`
	Percentage float64    `json:"percentage"`
	Swap       *SwapUsage `json:"swap,omitempty"`
}

// SwapUsage represents swap memory usage
type SwapUsage struct {
	Used       int64   `json:"used"`
	Available  int64   `json:"available"`
	Total      int64   `json:"total"`
	Percentage float64 `json:"percentage"`
}

// MemoryPressureLevel defines memory pressure levels
type MemoryPressureLevel string

const (
	MemoryPressureLow      MemoryPressureLevel = "low"
	MemoryPressureMedium   MemoryPressureLevel = "medium"
	MemoryPressureHigh     MemoryPressureLevel = "high"
	MemoryPressureCritical MemoryPressureLevel = "critical"
)

// GCStats represents garbage collection statistics
type GCStats struct {
	NumGC        uint32        `json:"numGC"`
	PauseTotal   time.Duration `json:"pauseTotal"`
	PauseEnd     []time.Time   `json:"pauseEnd"`
	PauseNs      []uint64      `json:"pauseNs"`
	LastGC       time.Time     `json:"lastGC"`
	NextGC       uint64        `json:"nextGC"`
	GCCPUPercent float64       `json:"gcCpuPercent"`
}

// MemoryThreshold defines memory thresholds for alerts
type MemoryThreshold struct {
	WarningPercent  float64 `json:"warningPercent"`
	CriticalPercent float64 `json:"criticalPercent"`
	MaxUsage        int64   `json:"maxUsage"`
	EnableAlerts    bool    `json:"enableAlerts"`
}

// MemoryAlert represents a memory-related alert
type MemoryAlert struct {
	ID             string          `json:"id"`
	Type           MemoryAlertType `json:"type"`
	Severity       Severity        `json:"severity"`
	Message        string          `json:"message"`
	CurrentValue   float64         `json:"currentValue"`
	ThresholdValue float64         `json:"thresholdValue"`
	Timestamp      time.Time       `json:"timestamp"`
	Resolved       bool            `json:"resolved"`
	ResolvedAt     *time.Time      `json:"resolvedAt,omitempty"`
}

// MemoryAlertType defines types of memory alerts
type MemoryAlertType string

const (
	MemoryAlertTypeUsage    MemoryAlertType = "usage"
	MemoryAlertTypeLeak     MemoryAlertType = "leak"
	MemoryAlertTypePressure MemoryAlertType = "pressure"
	MemoryAlertTypeGC       MemoryAlertType = "gc"
)

// MemoryProfile represents memory profiling results
type MemoryProfile struct {
	Duration    time.Duration          `json:"duration"`
	Samples     []*MemoryProfileSample `json:"samples"`
	Allocations []*AllocationProfile   `json:"allocations"`
	Hotspots    []*MemoryHotspot       `json:"hotspots"`
	Summary     *MemoryProfileSummary  `json:"summary"`
	GeneratedAt time.Time              `json:"generatedAt"`
}

// MemoryProfileSample represents a memory profile sample
type MemoryProfileSample struct {
	Timestamp   time.Time `json:"timestamp"`
	Usage       int64     `json:"usage"`
	Allocations int64     `json:"allocations"`
	Frees       int64     `json:"frees"`
	GCCount     uint32    `json:"gcCount"`
}

// AllocationProfile represents allocation profiling data
type AllocationProfile struct {
	Function   string  `json:"function"`
	File       string  `json:"file"`
	Line       int     `json:"line"`
	Count      int64   `json:"count"`
	Size       int64   `json:"size"`
	Percentage float64 `json:"percentage"`
}

// MemoryHotspot represents a memory hotspot
type MemoryHotspot struct {
	Location   string  `json:"location"`
	Size       int64   `json:"size"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
	Type       string  `json:"type"`
}

// MemoryProfileSummary represents memory profile summary
type MemoryProfileSummary struct {
	TotalAllocations int64   `json:"totalAllocations"`
	TotalSize        int64   `json:"totalSize"`
	PeakUsage        int64   `json:"peakUsage"`
	AverageUsage     int64   `json:"averageUsage"`
	GCOverhead       float64 `json:"gcOverhead"`
}

// MemoryPoolStats represents memory pool statistics
type MemoryPoolStats struct {
	PoolName      string    `json:"poolName"`
	BlockSize     int64     `json:"blockSize"`
	TotalBlocks   int       `json:"totalBlocks"`
	UsedBlocks    int       `json:"usedBlocks"`
	FreeBlocks    int       `json:"freeBlocks"`
	Utilization   float64   `json:"utilization"`
	HitRatio      float64   `json:"hitRatio"`
	Allocations   uint64    `json:"allocations"`
	Deallocations uint64    `json:"deallocations"`
	Timestamp     time.Time `json:"timestamp"`
}

// ConnectionPoolConfig defines configuration for connection pools
type ConnectionPoolConfig struct {
	MinConnections      int           `json:"minConnections"`
	MaxConnections      int           `json:"maxConnections"`
	ConnectionTimeout   time.Duration `json:"connectionTimeout"`
	IdleTimeout         time.Duration `json:"idleTimeout"`
	MaxLifetime         time.Duration `json:"maxLifetime"`
	ValidationQuery     string        `json:"validationQuery,omitempty"`
	ValidationTimeout   time.Duration `json:"validationTimeout"`
	RetryAttempts       int           `json:"retryAttempts"`
	RetryDelay          time.Duration `json:"retryDelay"`
	EnableHealthCheck   bool          `json:"enableHealthCheck"`
	HealthCheckInterval time.Duration `json:"healthCheckInterval"`
}

// ConnectionStats represents connection statistics
type ConnectionStats struct {
	TotalConnections   int                             `json:"totalConnections"`
	ActiveConnections  int                             `json:"activeConnections"`
	IdleConnections    int                             `json:"idleConnections"`
	FailedConnections  int                             `json:"failedConnections"`
	PoolStats          map[string]*ConnectionPoolStats `json:"poolStats"`
	ConnectionsByType  map[string]int                  `json:"connectionsByType"`
	ConnectionsByState map[ConnectionState]int         `json:"connectionsByState"`
	Timestamp          time.Time                       `json:"timestamp"`
}

// ConnectionState defines connection states
type ConnectionState string

const (
	ConnectionStateNew        ConnectionState = "new"
	ConnectionStateActive     ConnectionState = "active"
	ConnectionStateIdle       ConnectionState = "idle"
	ConnectionStateValidating ConnectionState = "validating"
	ConnectionStateClosed     ConnectionState = "closed"
	ConnectionStateFailed     ConnectionState = "failed"
)

// ConnectionPoolStats represents connection pool statistics
type ConnectionPoolStats struct {
	PoolName          string        `json:"poolName"`
	MinSize           int           `json:"minSize"`
	MaxSize           int           `json:"maxSize"`
	CurrentSize       int           `json:"currentSize"`
	ActiveConnections int           `json:"activeConnections"`
	IdleConnections   int           `json:"idleConnections"`
	PendingRequests   int           `json:"pendingRequests"`
	Utilization       float64       `json:"utilization"`
	WaitTime          time.Duration `json:"waitTime"`
	CreationRate      float64       `json:"creationRate"`
	DestructionRate   float64       `json:"destructionRate"`
	FailureRate       float64       `json:"failureRate"`
	Timestamp         time.Time     `json:"timestamp"`
}

// ConnectionMetrics represents individual connection metrics
type ConnectionMetrics struct {
	BytesRead    uint64        `json:"bytesRead"`
	BytesWritten uint64        `json:"bytesWritten"`
	RequestCount uint64        `json:"requestCount"`
	ErrorCount   uint64        `json:"errorCount"`
	LatencyP50   time.Duration `json:"latencyP50"`
	LatencyP95   time.Duration `json:"latencyP95"`
	LatencyP99   time.Duration `json:"latencyP99"`
	LastError    string        `json:"lastError,omitempty"`
	LastErrorAt  *time.Time    `json:"lastErrorAt,omitempty"`
}

// ThreadPoolConfig defines configuration for thread pools
type ThreadPoolConfig struct {
	CoreSize         int                   `json:"coreSize"`
	MaxSize          int                   `json:"maxSize"`
	QueueCapacity    int                   `json:"queueCapacity"`
	KeepAliveTime    time.Duration         `json:"keepAliveTime"`
	AllowCoreTimeout bool                  `json:"allowCoreTimeout"`
	RejectionPolicy  types.RejectionPolicy `json:"rejectionPolicy"`
	ThreadNamePrefix string                `json:"threadNamePrefix"`
	EnableMetrics    bool                  `json:"enableMetrics"`
}

// RejectionPolicy defined in types/common.go

// ThreadPoolStats represents thread pool statistics
type ThreadPoolStats struct {
	PoolName            string        `json:"poolName"`
	CoreSize            int           `json:"coreSize"`
	MaxSize             int           `json:"maxSize"`
	CurrentSize         int           `json:"currentSize"`
	ActiveCount         int           `json:"activeCount"`
	QueueSize           int           `json:"queueSize"`
	CompletedTasks      uint64        `json:"completedTasks"`
	TotalTasks          uint64        `json:"totalTasks"`
	RejectedTasks       uint64        `json:"rejectedTasks"`
	AverageTaskTime     time.Duration `json:"averageTaskTime"`
	AverageWaitTime     time.Duration `json:"averageWaitTime"`
	Utilization         float64       `json:"utilization"`
	QueueUtilization    float64       `json:"queueUtilization"`
	ThroughputPerSecond float64       `json:"throughputPerSecond"`
	Timestamp           time.Time     `json:"timestamp"`
}

// Priority defined in types/common.go

// Quota represents a resource quota
type Quota struct {
	Limit  int64         `json:"limit"`
	Period time.Duration `json:"period"`
	Scope  string        `json:"scope,omitempty"`
	Reset  QuotaReset    `json:"reset"`
	Burst  int64         `json:"burst,omitempty"`
}

// QuotaReset defines when quotas reset
type QuotaReset string

const (
	QuotaResetDaily   QuotaReset = "daily"
	QuotaResetHourly  QuotaReset = "hourly"
	QuotaResetMonthly QuotaReset = "monthly"
	QuotaResetNever   QuotaReset = "never"
)

// QuotaUsage represents current quota usage
type QuotaUsage struct {
	Used       int64     `json:"used"`
	Remaining  int64     `json:"remaining"`
	Limit      int64     `json:"limit"`
	ResetAt    time.Time `json:"resetAt"`
	Percentage float64   `json:"percentage"`
}

// ResourceLimit represents a resource limit
type ResourceLimit struct {
	SoftLimit int64         `json:"softLimit"`
	HardLimit int64         `json:"hardLimit"`
	Unit      string        `json:"unit"`
	Period    time.Duration `json:"period,omitempty"`
}

// LimitingStats represents resource limiting statistics
type LimitingStats struct {
	TotalRequests    uint64                     `json:"totalRequests"`
	AllowedRequests  uint64                     `json:"allowedRequests"`
	RejectedRequests uint64                     `json:"rejectedRequests"`
	RateLimitStats   map[string]*RateLimitStats `json:"rateLimitStats"`
	QuotaStats       map[string]*QuotaUsage     `json:"quotaStats"`
	Timestamp        time.Time                  `json:"timestamp"`
}

// RateLimitStats represents rate limiting statistics for a resource
type RateLimitStats struct {
	Resource         string    `json:"resource"`
	CurrentRate      float64   `json:"currentRate"`
	MaxRate          float64   `json:"maxRate"`
	BurstCapacity    int       `json:"burstCapacity"`
	TokensAvailable  int       `json:"tokensAvailable"`
	RequestsAllowed  uint64    `json:"requestsAllowed"`
	RequestsRejected uint64    `json:"requestsRejected"`
	LastRefill       time.Time `json:"lastRefill"`
}

// CacheConfig defines configuration for caches
type CacheConfig struct {
	MaxSize          int              `json:"maxSize"`
	TTL              time.Duration    `json:"ttl"`
	MaxIdleTime      time.Duration    `json:"maxIdleTime"`
	EvictionPolicy   EvictionPolicy   `json:"evictionPolicy"`
	EnableStats      bool             `json:"enableStats"`
	EnableMetrics    bool             `json:"enableMetrics"`
	Loader           CacheLoader      `json:"-"`
	Writer           CacheWriter      `json:"-"`
	EvictionListener EvictionListener `json:"-"`
}

// EvictionPolicy defines cache eviction policies
type EvictionPolicy string

const (
	EvictionPolicyLRU    EvictionPolicy = "lru"
	EvictionPolicyLFU    EvictionPolicy = "lfu"
	EvictionPolicyFIFO   EvictionPolicy = "fifo"
	EvictionPolicyRandom EvictionPolicy = "random"
	EvictionPolicyTTL    EvictionPolicy = "ttl"
)

// CacheStats represents cache statistics
type CacheStats struct {
	Name            string        `json:"name"`
	Size            int           `json:"size"`
	Capacity        int           `json:"capacity"`
	HitCount        uint64        `json:"hitCount"`
	MissCount       uint64        `json:"missCount"`
	LoadCount       uint64        `json:"loadCount"`
	EvictionCount   uint64        `json:"evictionCount"`
	HitRatio        float64       `json:"hitRatio"`
	MissRatio       float64       `json:"missRatio"`
	Utilization     float64       `json:"utilization"`
	AverageLoadTime time.Duration `json:"averageLoadTime"`
	Timestamp       time.Time     `json:"timestamp"`
}

// CacheUsage represents overall cache usage
type CacheUsage struct {
	TotalSize       int64     `json:"totalSize"`
	UsedSize        int64     `json:"usedSize"`
	AvailableSize   int64     `json:"availableSize"`
	Utilization     float64   `json:"utilization"`
	CacheCount      int       `json:"cacheCount"`
	TotalHits       uint64    `json:"totalHits"`
	TotalMisses     uint64    `json:"totalMisses"`
	OverallHitRatio float64   `json:"overallHitRatio"`
	Timestamp       time.Time `json:"timestamp"`
}

// CachePolicy defines global cache policies
type CachePolicy struct {
	DefaultTTL        time.Duration  `json:"defaultTtl"`
	MaxCacheSize      int64          `json:"maxCacheSize"`
	EvictionPolicy    EvictionPolicy `json:"evictionPolicy"`
	EnableCompression bool           `json:"enableCompression"`
	EnableEncryption  bool           `json:"enableEncryption"`
	StatsRetention    time.Duration  `json:"statsRetention"`
}

// ResourceStats represents overall resource statistics
type ResourceStats struct {
	Memory      *MemoryStats     `json:"memory"`
	Connections *ConnectionStats `json:"connections"`
	ThreadPools *ThreadPoolStats `json:"threadPools"`
	Caches      *CacheUsage      `json:"caches"`
	Limits      *LimitingStats   `json:"limits"`
	System      *SystemResources `json:"system"`
	Timestamp   time.Time        `json:"timestamp"`
}

// SystemResources represents system-level resource information
type SystemResources struct {
	CPU         *CPUResources     `json:"cpu"`
	Memory      *SystemMemory     `json:"memory"`
	Disk        *DiskResources    `json:"disk"`
	Network     *NetworkResources `json:"network"`
	FileHandles *FileHandles      `json:"fileHandles"`
}

// CPUResources represents CPU resource information
type CPUResources struct {
	Cores      int     `json:"cores"`
	Usage      float64 `json:"usage"`
	LoadAvg1   float64 `json:"loadAvg1"`
	LoadAvg5   float64 `json:"loadAvg5"`
	LoadAvg15  float64 `json:"loadAvg15"`
	UserTime   float64 `json:"userTime"`
	SystemTime float64 `json:"systemTime"`
	IdleTime   float64 `json:"idleTime"`
}

// SystemMemory represents system memory information
type SystemMemory struct {
	Total     uint64  `json:"total"`
	Available uint64  `json:"available"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	Cached    uint64  `json:"cached"`
	Buffers   uint64  `json:"buffers"`
	Usage     float64 `json:"usage"`
	SwapTotal uint64  `json:"swapTotal"`
	SwapUsed  uint64  `json:"swapUsed"`
	SwapFree  uint64  `json:"swapFree"`
}

// DiskResources represents disk resource information
type DiskResources struct {
	Partitions []DiskPartition `json:"partitions"`
	IOStats    *DiskIOStats    `json:"ioStats"`
}

// DiskPartition represents a disk partition
type DiskPartition struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	FileSystem  string  `json:"fileSystem"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	Usage       float64 `json:"usage"`
	InodesTotal uint64  `json:"inodesTotal"`
	InodesUsed  uint64  `json:"inodesUsed"`
	InodesFree  uint64  `json:"inodesFree"`
}

// DiskIOStats represents disk I/O statistics
type DiskIOStats struct {
	ReadBytes   uint64  `json:"readBytes"`
	WriteBytes  uint64  `json:"writeBytes"`
	ReadOps     uint64  `json:"readOps"`
	WriteOps    uint64  `json:"writeOps"`
	ReadTime    uint64  `json:"readTime"`
	WriteTime   uint64  `json:"writeTime"`
	IOTime      uint64  `json:"ioTime"`
	Utilization float64 `json:"utilization"`
}

// NetworkResources represents network resource information
type NetworkResources struct {
	Interfaces []NetworkInterface `json:"interfaces"`
	Stats      *NetworkStats      `json:"stats"`
}

// NetworkInterface represents a network interface
type NetworkInterface struct {
	Name            string `json:"name"`
	MTU             int    `json:"mtu"`
	Flags           string `json:"flags"`
	BytesReceived   uint64 `json:"bytesReceived"`
	BytesSent       uint64 `json:"bytesSent"`
	PacketsReceived uint64 `json:"packetsReceived"`
	PacketsSent     uint64 `json:"packetsSent"`
	ErrorsReceived  uint64 `json:"errorsReceived"`
	ErrorsSent      uint64 `json:"errorsSent"`
	DropsReceived   uint64 `json:"dropsReceived"`
	DropsSent       uint64 `json:"dropsSent"`
}

// NetworkStats represents overall network statistics
type NetworkStats struct {
	TotalBytesReceived   uint64 `json:"totalBytesReceived"`
	TotalBytesSent       uint64 `json:"totalBytesSent"`
	TotalPacketsReceived uint64 `json:"totalPacketsReceived"`
	TotalPacketsSent     uint64 `json:"totalPacketsSent"`
	TotalErrors          uint64 `json:"totalErrors"`
	TotalDrops           uint64 `json:"totalDrops"`
}

// FileHandles represents file handle information
type FileHandles struct {
	Current int     `json:"current"`
	Max     int     `json:"max"`
	Usage   float64 `json:"usage"`
}

// Callback interfaces
type CacheLoader func(key string) (interface{}, error)
type CacheWriter func(key string, value interface{}) error
type EvictionListener func(key string, value interface{}, reason EvictionReason)

// EvictionReason defines reasons for cache eviction
type EvictionReason string

const (
	EvictionReasonExpired  EvictionReason = "expired"
	EvictionReasonSize     EvictionReason = "size"
	EvictionReasonExplicit EvictionReason = "explicit"
	EvictionReasonReplaced EvictionReason = "replaced"
)

// Severity level (reused from other interfaces if needed)
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)
