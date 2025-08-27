package interfaces

import (
	"context"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/types"
)

// ErrorHandler defines the interface for error handling
type ErrorHandler interface {
	// HandleError handles an error with context
	HandleError(ctx context.Context, err error, operation string) error

	// ShouldRetry determines if an error should trigger a retry
	ShouldRetry(err error) bool

	// GetErrorSeverity returns the severity level of an error
	GetErrorSeverity(err error) ErrorSeverity

	// ClassifyError classifies an error into a category
	ClassifyError(err error) ErrorCategory

	// RegisterErrorMapping registers a custom error mapping
	RegisterErrorMapping(errorPattern string, category ErrorCategory, severity ErrorSeverity) error

	// GetErrorStats returns error statistics
	GetErrorStats() *ErrorStats

	// LogError logs an error with appropriate severity
	LogError(ctx context.Context, err error, operation string, metadata map[string]interface{})

	// NotifyError sends error notifications if configured
	NotifyError(ctx context.Context, err error, operation string) error
}

// RetryManager defines the interface for retry management
type RetryManager interface {
	// Execute executes an operation with retry logic
	Execute(ctx context.Context, operation RetryableOperation) error

	// ExecuteWithPolicy executes an operation with a specific retry policy
	ExecuteWithPolicy(ctx context.Context, operation RetryableOperation, policy *RetryPolicy) error

	// GetRetryPolicy returns the default retry policy
	GetRetryPolicy() *RetryPolicy

	// SetRetryPolicy sets the default retry policy
	SetRetryPolicy(policy *RetryPolicy) error

	// RegisterCustomPolicy registers a named retry policy
	RegisterCustomPolicy(name string, policy *RetryPolicy) error

	// GetCustomPolicy returns a named retry policy
	GetCustomPolicy(name string) (*RetryPolicy, error)

	// GetRetryStats returns retry statistics
	GetRetryStats() *RetryStats

	// IsOperationRetryable checks if an operation can be retried
	IsOperationRetryable(operationType string) bool
}

// CircuitBreaker defines the interface for circuit breaker pattern
type CircuitBreaker interface {
	// Execute executes an operation through the circuit breaker
	Execute(ctx context.Context, operation CircuitBreakerOperation) error

	// GetState returns the current state of the circuit breaker
	GetState() CircuitBreakerState

	// ForceOpen forces the circuit breaker to open state
	ForceOpen() error

	// ForceClose forces the circuit breaker to close state
	ForceClose() error

	// Reset resets the circuit breaker to initial state
	Reset() error

	// GetMetrics returns circuit breaker metrics
	GetMetrics() *CircuitBreakerMetrics

	// RegisterStateChangeCallback registers callback for state changes
	RegisterStateChangeCallback(callback StateChangeCallback) error
}

// Timeout defines the interface for timeout management
type Timeout interface {
	// WithTimeout wraps an operation with timeout
	WithTimeout(ctx context.Context, operation TimeoutOperation, timeout time.Duration) error

	// WithDeadline wraps an operation with deadline
	WithDeadline(ctx context.Context, operation TimeoutOperation, deadline time.Time) error

	// GetDefaultTimeout returns the default timeout for an operation type
	GetDefaultTimeout(operationType string) time.Duration

	// SetDefaultTimeout sets the default timeout for an operation type
	SetDefaultTimeout(operationType string, timeout time.Duration) error

	// GetTimeoutStats returns timeout statistics
	GetTimeoutStats() *TimeoutStats
}

// BulkheadIsolation defines the interface for bulkhead isolation pattern
type BulkheadIsolation interface {
	// Execute executes an operation in an isolated thread pool
	Execute(ctx context.Context, poolName string, operation BulkheadOperation) error

	// CreatePool creates a new isolation pool
	CreatePool(name string, config *PoolConfig) error

	// GetPoolStats returns statistics for a pool
	GetPoolStats(poolName string) (*PoolStats, error)

	// GetAllPoolStats returns statistics for all pools
	GetAllPoolStats() (map[string]*PoolStats, error)

	// ResizePool resizes an existing pool
	ResizePool(poolName string, newSize int) error

	// ShutdownPool shuts down a pool gracefully
	ShutdownPool(poolName string, timeout time.Duration) error
}

// FaultInjector defines the interface for fault injection (testing purposes)
type FaultInjector interface {
	// InjectError injects an error for testing
	InjectError(operationType string, errorType ErrorType, probability float64) error

	// InjectLatency injects latency for testing
	InjectLatency(operationType string, latency time.Duration, probability float64) error

	// ClearFaults clears all injected faults
	ClearFaults() error

	// GetInjectedFaults returns all currently injected faults
	GetInjectedFaults() []FaultInjection

	// EnableFaultInjection enables or disables fault injection
	EnableFaultInjection(enabled bool) error

	// IsEnabled returns true if fault injection is enabled
	IsEnabled() bool
}

// Data structures

// ErrorSeverity defines error severity levels
type ErrorSeverity string

const (
	ErrorSeverityTrace    ErrorSeverity = "trace"
	ErrorSeverityDebug    ErrorSeverity = "debug"
	ErrorSeverityInfo     ErrorSeverity = "info"
	ErrorSeverityWarning  ErrorSeverity = "warning"
	ErrorSeverityError    ErrorSeverity = "error"
	ErrorSeverityCritical ErrorSeverity = "critical"
	ErrorSeverityFatal    ErrorSeverity = "fatal"
)

// ErrorCategory defines error categories
type ErrorCategory string

const (
	ErrorCategoryNetwork       ErrorCategory = "network"
	ErrorCategoryDatabase      ErrorCategory = "database"
	ErrorCategoryTimeout       ErrorCategory = "timeout"
	ErrorCategoryValidation    ErrorCategory = "validation"
	ErrorCategoryAuthorization ErrorCategory = "authorization"
	ErrorCategoryConfiguration ErrorCategory = "configuration"
	ErrorCategoryResource      ErrorCategory = "resource"
	ErrorCategoryBusiness      ErrorCategory = "business"
	ErrorCategorySystem        ErrorCategory = "system"
	ErrorCategoryExternal      ErrorCategory = "external"
	ErrorCategoryUnknown       ErrorCategory = "unknown"
)

// ErrorStats represents error statistics
type ErrorStats struct {
	TotalErrors         uint64                    `json:"totalErrors"`
	ErrorsByCategory    map[ErrorCategory]uint64  `json:"errorsByCategory"`
	ErrorsBySeverity    map[ErrorSeverity]uint64  `json:"errorsBySeverity"`
	ErrorsByOperation   map[string]uint64         `json:"errorsByOperation"`
	RecentErrors        []ErrorEvent              `json:"recentErrors"`
	ErrorRate           float64                   `json:"errorRate"` // errors per second
	ErrorRateByCategory map[ErrorCategory]float64 `json:"errorRateByCategory"`
	LastReset           time.Time                 `json:"lastReset"`
	Timestamp           time.Time                 `json:"timestamp"`
}

// ErrorEvent represents a single error event
type ErrorEvent struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	Error      string                 `json:"error"`
	Operation  string                 `json:"operation"`
	Category   ErrorCategory          `json:"category"`
	Severity   ErrorSeverity          `json:"severity"`
	Context    map[string]interface{} `json:"context,omitempty"`
	StackTrace string                 `json:"stackTrace,omitempty"`
	Retryable  bool                   `json:"retryable"`
}

// RetryPolicy defines retry behavior
type RetryPolicy struct {
	MaxAttempts        int                   `json:"maxAttempts"`
	InitialDelay       time.Duration         `json:"initialDelay"`
	MaxDelay           time.Duration         `json:"maxDelay"`
	BackoffMultiplier  float64               `json:"backoffMultiplier"`
	BackoffStrategy    BackoffStrategy       `json:"backoffStrategy"`
	RetryableErrors    []ErrorCategory       `json:"retryableErrors"`
	NonRetryableErrors []ErrorCategory       `json:"nonRetryableErrors"`
	Jitter             bool                  `json:"jitter"`
	JitterRange        float64               `json:"jitterRange"`
	RetryCondition     func(error, int) bool `json:"-"` // custom retry condition
	OnRetry            func(error, int)      `json:"-"` // callback on retry attempt
}

// BackoffStrategy defines backoff strategies
type BackoffStrategy string

const (
	BackoffStrategyFixed       BackoffStrategy = "fixed"
	BackoffStrategyLinear      BackoffStrategy = "linear"
	BackoffStrategyExponential BackoffStrategy = "exponential"
	BackoffStrategyCustom      BackoffStrategy = "custom"
)

// RetryStats represents retry statistics
type RetryStats struct {
	TotalAttempts        uint64                   `json:"totalAttempts"`
	SuccessfulRetries    uint64                   `json:"successfulRetries"`
	FailedRetries        uint64                   `json:"failedRetries"`
	RetriesByOperation   map[string]uint64        `json:"retriesByOperation"`
	RetriesByErrorType   map[ErrorCategory]uint64 `json:"retriesByErrorType"`
	AverageRetryCount    float64                  `json:"averageRetryCount"`
	MaxRetryCount        int                      `json:"maxRetryCount"`
	RetrySuccessRate     float64                  `json:"retrySuccessRate"`
	AverageRetryDuration time.Duration            `json:"averageRetryDuration"`
	Timestamp            time.Time                `json:"timestamp"`
}

// CircuitBreakerState defines circuit breaker states
type CircuitBreakerState string

const (
	CircuitBreakerStateClosed   CircuitBreakerState = "closed"
	CircuitBreakerStateOpen     CircuitBreakerState = "open"
	CircuitBreakerStateHalfOpen CircuitBreakerState = "half_open"
)

// CircuitBreakerConfig defines circuit breaker configuration
type CircuitBreakerConfig struct {
	FailureThreshold    int           `json:"failureThreshold"`
	RecoveryTimeout     time.Duration `json:"recoveryTimeout"`
	SuccessThreshold    int           `json:"successThreshold"`
	Timeout             time.Duration `json:"timeout"`
	MaxConcurrentCalls  int           `json:"maxConcurrentCalls"`
	SlidingWindowSize   int           `json:"slidingWindowSize"`
	SlidingWindowType   WindowType    `json:"slidingWindowType"`
	MinimumRequestCount int           `json:"minimumRequestCount"`
}

// WindowType defines sliding window types
type WindowType string

const (
	WindowTypeCount WindowType = "count"
	WindowTypeTime  WindowType = "time"
)

// CircuitBreakerMetrics represents circuit breaker metrics
type CircuitBreakerMetrics struct {
	State               CircuitBreakerState `json:"state"`
	TotalCalls          uint64              `json:"totalCalls"`
	SuccessfulCalls     uint64              `json:"successfulCalls"`
	FailedCalls         uint64              `json:"failedCalls"`
	RejectedCalls       uint64              `json:"rejectedCalls"`
	LastFailureTime     *time.Time          `json:"lastFailureTime,omitempty"`
	LastSuccessTime     *time.Time          `json:"lastSuccessTime,omitempty"`
	StateTransitions    uint64              `json:"stateTransitions"`
	SuccessRate         float64             `json:"successRate"`
	FailureRate         float64             `json:"failureRate"`
	AverageResponseTime time.Duration       `json:"averageResponseTime"`
	Timestamp           time.Time           `json:"timestamp"`
}

// TimeoutStats represents timeout statistics
type TimeoutStats struct {
	TotalTimeouts       uint64                   `json:"totalTimeouts"`
	TimeoutsByOperation map[string]uint64        `json:"timeoutsByOperation"`
	AverageTimeout      time.Duration            `json:"averageTimeout"`
	TimeoutRate         float64                  `json:"timeoutRate"`
	OperationTimeouts   map[string]time.Duration `json:"operationTimeouts"`
	RecentTimeouts      []TimeoutEvent           `json:"recentTimeouts"`
	Timestamp           time.Time                `json:"timestamp"`
}

// TimeoutEvent represents a timeout event
type TimeoutEvent struct {
	ID           string                 `json:"id"`
	Timestamp    time.Time              `json:"timestamp"`
	Operation    string                 `json:"operation"`
	TimeoutValue time.Duration          `json:"timeoutValue"`
	ActualTime   time.Duration          `json:"actualTime"`
	Context      map[string]interface{} `json:"context,omitempty"`
}

// PoolConfig defines configuration for bulkhead isolation pool
type PoolConfig struct {
	CoreSize        int                   `json:"coreSize"`
	MaxSize         int                   `json:"maxSize"`
	QueueCapacity   int                   `json:"queueCapacity"`
	KeepAliveTime   time.Duration         `json:"keepAliveTime"`
	RejectionPolicy types.RejectionPolicy `json:"rejectionPolicy"`
}

// RejectionPolicy defined in types/common.go

// PoolStats represents thread pool statistics
type PoolStats struct {
	PoolName         string        `json:"poolName"`
	CoreSize         int           `json:"coreSize"`
	MaxSize          int           `json:"maxSize"`
	ActiveThreads    int           `json:"activeThreads"`
	QueuedTasks      int           `json:"queuedTasks"`
	CompletedTasks   uint64        `json:"completedTasks"`
	RejectedTasks    uint64        `json:"rejectedTasks"`
	AverageTaskTime  time.Duration `json:"averageTaskTime"`
	Utilization      float64       `json:"utilization"`
	QueueUtilization float64       `json:"queueUtilization"`
	Timestamp        time.Time     `json:"timestamp"`
}

// FaultInjection represents an injected fault
type FaultInjection struct {
	ID            string      `json:"id"`
	OperationType string      `json:"operationType"`
	FaultType     FaultType   `json:"faultType"`
	Probability   float64     `json:"probability"`
	Parameters    interface{} `json:"parameters,omitempty"`
	CreatedAt     time.Time   `json:"createdAt"`
	HitCount      uint64      `json:"hitCount"`
	Enabled       bool        `json:"enabled"`
}

// FaultType defines types of faults that can be injected
type FaultType string

const (
	FaultTypeError   FaultType = "error"
	FaultTypeLatency FaultType = "latency"
	FaultTypeTimeout FaultType = "timeout"
	FaultTypeReject  FaultType = "reject"
)

// ErrorType defines types of errors for fault injection
type ErrorType string

const (
	ErrorTypeNetwork    ErrorType = "network_error"
	ErrorTypeDatabase   ErrorType = "database_error"
	ErrorTypeTimeout    ErrorType = "timeout_error"
	ErrorTypeResource   ErrorType = "resource_error"
	ErrorTypeValidation ErrorType = "validation_error"
	ErrorTypeUnexpected ErrorType = "unexpected_error"
)

// Operation function types
type RetryableOperation func(ctx context.Context) error
type CircuitBreakerOperation func(ctx context.Context) (interface{}, error)
type TimeoutOperation func(ctx context.Context) error
type BulkheadOperation func(ctx context.Context) (interface{}, error)

// Callback function types
type StateChangeCallback func(oldState, newState CircuitBreakerState)

// ErrorHandlerConfig defines configuration for error handling
type ErrorHandlerConfig struct {
	// Error classification
	EnableErrorClassification bool                        `json:"enableErrorClassification"`
	CustomErrorMappings       map[string]ErrorMappingRule `json:"customErrorMappings"`

	// Error logging
	LogErrors         bool          `json:"logErrors"`
	LogLevel          ErrorSeverity `json:"logLevel"`
	IncludeStackTrace bool          `json:"includeStackTrace"`

	// Error notification
	EnableNotifications    bool                  `json:"enableNotifications"`
	NotificationChannels   []NotificationChannel `json:"notificationChannels"`
	NotificationThresholds map[ErrorSeverity]int `json:"notificationThresholds"`

	// Statistics
	EnableStatistics        bool          `json:"enableStatistics"`
	StatisticsRetentionTime time.Duration `json:"statisticsRetentionTime"`
	RecentErrorsLimit       int           `json:"recentErrorsLimit"`
}

// ErrorMappingRule defines how to map an error pattern to category/severity
type ErrorMappingRule struct {
	Pattern   string        `json:"pattern"` // regex pattern or error type
	Category  ErrorCategory `json:"category"`
	Severity  ErrorSeverity `json:"severity"`
	Retryable bool          `json:"retryable"`
}

// NotificationChannel defines how to send error notifications
type NotificationChannel struct {
	Type    NotificationType       `json:"type"`
	Config  map[string]interface{} `json:"config"`
	Enabled bool                   `json:"enabled"`
	Filters []NotificationFilter   `json:"filters,omitempty"`
}

// NotificationType defines types of notification channels
type NotificationType string

const (
	NotificationTypeEmail     NotificationType = "email"
	NotificationTypeSlack     NotificationType = "slack"
	NotificationTypeWebhook   NotificationType = "webhook"
	NotificationTypeSMS       NotificationType = "sms"
	NotificationTypePagerDuty NotificationType = "pagerduty"
)

// NotificationFilter defines filters for notifications
type NotificationFilter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// ResilenceManager defines a comprehensive resilience management interface
type ResilienceManager interface {
	// Get individual components
	GetErrorHandler() ErrorHandler
	GetRetryManager() RetryManager
	GetCircuitBreaker(name string) CircuitBreaker
	GetTimeout() Timeout
	GetBulkheadIsolation() BulkheadIsolation
	GetFaultInjector() FaultInjector

	// Create circuit breakers
	CreateCircuitBreaker(name string, config *CircuitBreakerConfig) (CircuitBreaker, error)

	// Execute operations with full resilience
	ExecuteResilient(ctx context.Context, operationName string, operation func(ctx context.Context) error) error

	// Get comprehensive metrics
	GetResilienceMetrics() *ResilienceMetrics

	// Health check
	HealthCheck(ctx context.Context) error
}

// ResilienceMetrics represents comprehensive resilience metrics
type ResilienceMetrics struct {
	ErrorStats          *ErrorStats                       `json:"errorStats"`
	RetryStats          *RetryStats                       `json:"retryStats"`
	CircuitBreakerStats map[string]*CircuitBreakerMetrics `json:"circuitBreakerStats"`
	TimeoutStats        *TimeoutStats                     `json:"timeoutStats"`
	BulkheadStats       map[string]*PoolStats             `json:"bulkheadStats"`
	OverallSuccessRate  float64                           `json:"overallSuccessRate"`
	OverallAvailability float64                           `json:"overallAvailability"`
	Timestamp           time.Time                         `json:"timestamp"`
}
