// Package types contains common types and structures used across YaoVerse components
package types

import (
	"context"
	"time"
)

// TimeRange represents a time interval with start and end times
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// IsValid checks if the time range is valid
func (tr TimeRange) IsValid() bool {
	return tr.Start.Before(tr.End) || tr.Start.Equal(tr.End)
}

// Duration returns the duration of the time range
func (tr TimeRange) Duration() time.Duration {
	if !tr.IsValid() {
		return 0
	}
	return tr.End.Sub(tr.Start)
}

// Contains checks if a given time is within the range
func (tr TimeRange) Contains(t time.Time) bool {
	return (t.Equal(tr.Start) || t.After(tr.Start)) && (t.Equal(tr.End) || t.Before(tr.End))
}

// HealthStatus represents the health status of a component
type HealthStatus string

const (
	HealthStatusHealthy   HealthStatus = "healthy"
	HealthStatusUnhealthy HealthStatus = "unhealthy"
	HealthStatusDegraded  HealthStatus = "degraded"
	HealthStatusUnknown   HealthStatus = "unknown"
)

// IsHealthy checks if the status is healthy
func (hs HealthStatus) IsHealthy() bool {
	return hs == HealthStatusHealthy
}

// String returns the string representation of health status
func (hs HealthStatus) String() string {
	return string(hs)
}

// HealthRecord represents a health check record
type HealthRecord struct {
	Timestamp     time.Time     `json:"timestamp"`
	Status        HealthStatus  `json:"status"`
	Message       string        `json:"message,omitempty"`
	Details       interface{}   `json:"details,omitempty"`
	Error         string        `json:"error,omitempty"`
	CheckedBy     string        `json:"checkedBy,omitempty"`
	NodeID        string        `json:"nodeId,omitempty"`
	ComponentName string        `json:"componentName,omitempty"`
	Duration      time.Duration `json:"duration,omitempty"`
	ErrorCount    int           `json:"errorCount,omitempty"`
}

// HealthCheckFunc defines the signature for health check functions
type HealthCheckFunc func(ctx context.Context) (*HealthRecord, error)

// MemoryUsage represents memory usage statistics
type MemoryUsage struct {
	Total       uint64  `json:"total"`       // Total memory
	Used        uint64  `json:"used"`        // Used memory
	Available   uint64  `json:"available"`   // Available memory
	UsedPercent float64 `json:"usedPercent"` // Usage percentage
	Buffers     uint64  `json:"buffers"`     // Buffer memory
	Cached      uint64  `json:"cached"`      // Cached memory
	SwapTotal   uint64  `json:"swapTotal"`   // Total swap memory
	SwapUsed    uint64  `json:"swapUsed"`    // Used swap memory
}

// IsHigh checks if memory usage is high (>80%)
func (mu MemoryUsage) IsHigh() bool {
	return mu.UsedPercent > 80.0
}

// IsCritical checks if memory usage is critical (>95%)
func (mu MemoryUsage) IsCritical() bool {
	return mu.UsedPercent > 95.0
}

// NotificationChannel represents a notification delivery channel
type NotificationChannel struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name,omitempty"`
	Type     string                 `json:"type"` // email, webhook, slack, etc.
	Enabled  bool                   `json:"enabled"`
	Config   map[string]interface{} `json:"config"`             // General config/settings
	Settings map[string]interface{} `json:"settings,omitempty"` // Alternative name for config
	Filters  []string               `json:"filters,omitempty"`  // event type filters
	Labels   map[string]string      `json:"labels,omitempty"`   // Additional labels
	Created  time.Time              `json:"created,omitempty"`
	Updated  time.Time              `json:"updated,omitempty"`
}

// IsActive checks if the notification channel is active
func (nc NotificationChannel) IsActive() bool {
	return nc.Enabled
}

// SupportsEventType checks if the channel supports a specific event type
func (nc NotificationChannel) SupportsEventType(eventType string) bool {
	if len(nc.Filters) == 0 {
		return true // No filters means all events are supported
	}

	for _, filter := range nc.Filters {
		if filter == eventType {
			return true
		}
	}
	return false
}

// LogLevel represents the severity level of a log entry
type LogLevel int

const (
	LogLevelTrace LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

// String returns the string representation of log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelTrace:
		return "trace"
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	case LogLevelFatal:
		return "fatal"
	default:
		return "unknown"
	}
}

// LogEntry represents a single log entry
type LogEntry struct {
	ID        string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Component string                 `json:"component"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	TraceID   string                 `json:"traceId,omitempty"`
	SpanID    string                 `json:"spanId,omitempty"`
	Source    string                 `json:"source,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// LogQuery represents a log search query
type LogQuery struct {
	TimeRange  TimeRange              `json:"timeRange"`
	Levels     []LogLevel             `json:"levels,omitempty"`
	Components []string               `json:"components,omitempty"`
	Keywords   []string               `json:"keywords,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	TraceID    string                 `json:"traceId,omitempty"`
	Limit      int                    `json:"limit,omitempty"`
	Offset     int                    `json:"offset,omitempty"`
	SortBy     string                 `json:"sortBy,omitempty"`
	SortOrder  string                 `json:"sortOrder,omitempty"` // asc, desc
}

// LogQueryResult represents the result of a log query
type LogQueryResult struct {
	Entries    []*LogEntry   `json:"entries"`
	TotalCount int           `json:"totalCount"`
	HasMore    bool          `json:"hasMore"`
	QueryTime  time.Duration `json:"queryTime"`
	ExecutedAt time.Time     `json:"executedAt"`
}

// LogAnomaly represents an detected anomaly in logs
type LogAnomaly struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`     // spike, drop, pattern, etc.
	Severity    string                 `json:"severity"` // low, medium, high, critical
	Component   string                 `json:"component,omitempty"`
	Description string                 `json:"description"`
	TimeRange   TimeRange              `json:"timeRange"`
	Confidence  float64                `json:"confidence"` // 0.0 to 1.0
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
	DetectedAt  time.Time              `json:"detectedAt"`
}

// LogPattern represents a detected pattern in logs
type LogPattern struct {
	ID         string      `json:"id"`
	Pattern    string      `json:"pattern"`
	Frequency  int         `json:"frequency"`
	Components []string    `json:"components"`
	FirstSeen  time.Time   `json:"firstSeen"`
	LastSeen   time.Time   `json:"lastSeen"`
	Confidence float64     `json:"confidence"`
	Examples   []*LogEntry `json:"examples,omitempty"`
}

// ErrorTrends represents error trend analysis results
type ErrorTrends struct {
	TimeRange         TimeRange        `json:"timeRange"`
	TotalErrors       int              `json:"totalErrors"`
	ErrorsByLevel     map[LogLevel]int `json:"errorsByLevel"`
	ErrorsByComponent map[string]int   `json:"errorsByComponent"`
	ErrorsByHour      map[int]int      `json:"errorsByHour"` // hour of day -> count
	TopErrors         []*ErrorSummary  `json:"topErrors"`
	TrendDirection    string           `json:"trendDirection"` // increasing, decreasing, stable
	AnalyzedAt        time.Time        `json:"analyzedAt"`
}

// ErrorSummary represents a summary of a specific error
type ErrorSummary struct {
	Message    string    `json:"message"`
	Count      int       `json:"count"`
	FirstSeen  time.Time `json:"firstSeen"`
	LastSeen   time.Time `json:"lastSeen"`
	Components []string  `json:"components"`
}

// LogMetrics represents log-based metrics
type LogMetrics struct {
	TimeRange       TimeRange        `json:"timeRange"`
	TotalLogs       int              `json:"totalLogs"`
	LogsByLevel     map[LogLevel]int `json:"logsByLevel"`
	LogsByComponent map[string]int   `json:"logsByComponent"`
	ErrorRate       float64          `json:"errorRate"`
	WarnRate        float64          `json:"warnRate"`
	LogVelocity     float64          `json:"logVelocity"` // logs per second
	GeneratedAt     time.Time        `json:"generatedAt"`
}

// LogStats represents overall log statistics
type LogStats struct {
	TotalLogs       int64              `json:"totalLogs"`
	LogsByLevel     map[LogLevel]int64 `json:"logsByLevel"`
	LogsByComponent map[string]int64   `json:"logsByComponent"`
	OldestLog       time.Time          `json:"oldestLog"`
	NewestLog       time.Time          `json:"newestLog"`
	IndexSize       int64              `json:"indexSize"`
	StorageUsed     int64              `json:"storageUsed"`
	UpdatedAt       time.Time          `json:"updatedAt"`
}

// ValidationResult represents the result of a validation operation
type ValidationResult struct {
	Valid       bool                   `json:"valid"`
	Errors      []ValidationError      `json:"errors,omitempty"`
	Warnings    []ValidationWarning    `json:"warnings,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	ValidatedAt time.Time              `json:"validatedAt,omitempty"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Rule    string      `json:"rule,omitempty"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Path    string      `json:"path,omitempty"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Path    string      `json:"path,omitempty"`
}

// Priority represents task priority levels
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// String returns the string representation of priority
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// ConnectionPoolStats defined in interfaces/resource_management.go

// RejectionPolicy defines how to handle rejected tasks
type RejectionPolicy string

const (
	RejectionPolicyAbort         RejectionPolicy = "abort"
	RejectionPolicyCallerRuns    RejectionPolicy = "caller_runs"
	RejectionPolicyDiscardOldest RejectionPolicy = "discard_oldest"
	RejectionPolicyDiscard       RejectionPolicy = "discard"
)

// String returns the string representation of rejection policy
func (rp RejectionPolicy) String() string {
	return string(rp)
}

// Severity defines severity levels
type Severity string

const (
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// String returns the string representation of severity
func (s Severity) String() string {
	return string(s)
}
