package interfaces

import (
	"context"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/types"
)

// MetricsCollector defines the interface for collecting metrics
type MetricsCollector interface {
	// Counter metrics
	IncrementCounter(name string, labels map[string]string) error
	IncrementCounterBy(name string, value float64, labels map[string]string) error
	GetCounter(name string, labels map[string]string) (float64, error)

	// Gauge metrics
	SetGauge(name string, value float64, labels map[string]string) error
	IncrementGauge(name string, labels map[string]string) error
	DecrementGauge(name string, labels map[string]string) error
	AddToGauge(name string, value float64, labels map[string]string) error
	GetGauge(name string, labels map[string]string) (float64, error)

	// Histogram metrics
	RecordHistogram(name string, value float64, labels map[string]string) error
	GetHistogram(name string, labels map[string]string) (*HistogramData, error)

	// Summary metrics
	RecordSummary(name string, value float64, labels map[string]string) error
	GetSummary(name string, labels map[string]string) (*SummaryData, error)

	// Time series data
	RecordTimeSeries(name string, value float64, timestamp time.Time, labels map[string]string) error
	GetTimeSeries(name string, start, end time.Time, labels map[string]string) ([]DataPoint, error)

	// Metric management
	RegisterMetric(metric *MetricDefinition) error
	UnregisterMetric(name string) error
	ListMetrics() ([]*MetricDefinition, error)
	GetMetricMetadata(name string) (*MetricMetadata, error)

	// Batch operations
	RecordBatch(batch []*MetricRecord) error

	// Export and query
	ExportMetrics(format ExportFormat) ([]byte, error)
	QueryMetrics(query string) (*QueryResult, error)
}

// PerformanceMonitor defines the interface for performance monitoring
type PerformanceMonitor interface {
	// Operation timing
	StartTimer(operationName string, labels map[string]string) Timer
	RecordDuration(operationName string, duration time.Duration, labels map[string]string) error
	RecordLatency(operationName string, latency time.Duration, labels map[string]string) error

	// Resource monitoring
	RecordCPUUsage(usage float64) error
	RecordMemoryUsage(usage int64) error
	RecordDiskUsage(path string, usage DiskUsage) error
	RecordNetworkUsage(interfaceName string, usage NetworkUsage) error

	// Custom performance metrics
	RecordThroughput(operationName string, count int64, duration time.Duration) error
	RecordQueueLength(queueName string, length int) error
	RecordConnectionPool(poolName string, stats ConnectionPoolStats) error

	// Performance analysis
	GetPerformanceReport(timeRange types.TimeRange) (*PerformanceReport, error)
	GetBottleneckAnalysis(timeRange types.TimeRange) (*BottleneckAnalysis, error)
	GetTrendAnalysis(metric string, timeRange types.TimeRange) (*TrendAnalysis, error)

	// Alerting on performance
	SetPerformanceThreshold(metric string, threshold Threshold) error
	GetPerformanceAlerts() ([]*Alert, error)
}

// AlertManager defines the interface for alert management
type AlertManager interface {
	// Alert rules
	CreateAlertRule(rule *AlertRule) error
	UpdateAlertRule(ruleID string, rule *AlertRule) error
	DeleteAlertRule(ruleID string) error
	GetAlertRule(ruleID string) (*AlertRule, error)
	ListAlertRules() ([]*AlertRule, error)

	// Alert evaluation
	EvaluateRules(ctx context.Context) error
	EvaluateRule(ctx context.Context, ruleID string) error
	GetRuleEvaluationHistory(ruleID string, limit int) ([]*RuleEvaluation, error)

	// Active alerts
	GetActiveAlerts() ([]*Alert, error)
	GetAlert(alertID string) (*Alert, error)
	ResolveAlert(alertID string, resolvedBy string, reason string) error
	AcknowledgeAlert(alertID string, acknowledgedBy string, reason string) error

	// Alert notification
	SendAlert(alert *Alert, channels []string) error
	GetNotificationChannels() ([]*NotificationChannel, error)
	CreateNotificationChannel(channel *NotificationChannel) error
	UpdateNotificationChannel(channelID string, channel *NotificationChannel) error
	DeleteNotificationChannel(channelID string) error

	// Alert statistics
	GetAlertStats(timeRange types.TimeRange) (*AlertStats, error)
	GetAlertHistory(timeRange types.TimeRange, limit int) ([]*Alert, error)
}

// HealthChecker defines the interface for health checking
type HealthChecker interface {
	// Component health
	RegisterHealthCheck(name string, checker HealthCheckFunc) error
	UnregisterHealthCheck(name string) error
	GetHealthStatus(name string) (*types.HealthStatus, error)
	GetAllHealthStatuses() (map[string]*types.HealthStatus, error)

	// System health
	GetSystemHealth() (*SystemHealth, error)
	GetServiceHealth(serviceName string) (*ServiceHealth, error)
	GetDependencyHealth() (map[string]*types.HealthStatus, error)

	// Health monitoring
	StartHealthMonitoring(ctx context.Context, interval time.Duration) error
	StopHealthMonitoring() error
	GetHealthHistory(componentName string, timeRange types.TimeRange) ([]*types.HealthRecord, error)

	// Health alerting
	SetHealthThreshold(componentName string, threshold HealthThreshold) error
	GetHealthAlerts() ([]*HealthAlert, error)
}

// LogAnalyzer defines the interface for log analysis
type LogAnalyzer interface {
	// Log ingestion
	IngestLog(entry *types.LogEntry) error
	IngestLogBatch(entries []*types.LogEntry) error

	// Log querying
	QueryLogs(query *types.LogQuery) (*types.LogQueryResult, error)
	GetLogsByLevel(level types.LogLevel, timeRange types.TimeRange) ([]*types.LogEntry, error)
	GetLogsByComponent(component string, timeRange types.TimeRange) ([]*types.LogEntry, error)

	// Log analysis
	DetectAnomalies(timeRange types.TimeRange) ([]*types.LogAnomaly, error)
	GetLogPatterns(timeRange types.TimeRange) ([]*types.LogPattern, error)
	GetErrorTrends(timeRange types.TimeRange) (*types.ErrorTrends, error)

	// Log metrics
	GenerateLogMetrics(timeRange types.TimeRange) (*LogMetrics, error)
	GetLogStats() (*LogStats, error)
}

// TraceCollector defines the interface for distributed tracing
type TraceCollector interface {
	// Trace management
	StartTrace(operationName string, parentSpanID *string) (Trace, error)
	FinishTrace(traceID string) error
	GetTrace(traceID string) (*Trace, error)
	SearchTraces(query *TraceQuery) ([]*Trace, error)

	// Span management within traces
	AddSpan(traceID string, span *Span) error
	UpdateSpan(traceID string, spanID string, updates map[string]interface{}) error
	GetSpan(traceID string, spanID string) (*Span, error)

	// Trace analysis
	GetTraceAnalytics(timeRange types.TimeRange) (*TraceAnalytics, error)
	GetServiceMap(timeRange types.TimeRange) (*ServiceMap, error)
	GetTraceDependencies() (*TraceDependencies, error)

	// Performance insights from traces
	GetSlowTraces(threshold time.Duration, limit int) ([]*Trace, error)
	GetErrorTraces(timeRange types.TimeRange, limit int) ([]*Trace, error)
}

// Data structures

// Timer interface for timing operations
type Timer interface {
	Stop() time.Duration
	Reset() Timer
	GetDuration() time.Duration
}

// MetricDefinition defines a metric
type MetricDefinition struct {
	Name        string            `json:"name"`
	Type        MetricType        `json:"type"`
	Description string            `json:"description"`
	Unit        string            `json:"unit"`
	Labels      []string          `json:"labels"`
	Buckets     []float64         `json:"buckets,omitempty"`    // for histograms
	Quantiles   []float64         `json:"quantiles,omitempty"`  // for summaries
	MaxAge      time.Duration     `json:"maxAge,omitempty"`     // for summaries
	AgeBuckets  int               `json:"ageBuckets,omitempty"` // for summaries
	Properties  map[string]string `json:"properties,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
}

// MetricType defines types of metrics
type MetricType string

const (
	MetricTypeCounter    MetricType = "counter"
	MetricTypeGauge      MetricType = "gauge"
	MetricTypeHistogram  MetricType = "histogram"
	MetricTypeSummary    MetricType = "summary"
	MetricTypeTimeSeries MetricType = "timeseries"
)

// MetricMetadata contains metadata about a metric
type MetricMetadata struct {
	Name        string            `json:"name"`
	Type        MetricType        `json:"type"`
	Description string            `json:"description"`
	Unit        string            `json:"unit"`
	SampleCount int64             `json:"sampleCount"`
	LastUpdated time.Time         `json:"lastUpdated"`
	Labels      []string          `json:"labels"`
	Properties  map[string]string `json:"properties"`
}

// MetricRecord represents a single metric record
type MetricRecord struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels"`
	Type      MetricType        `json:"type"`
}

// HistogramData represents histogram metric data
type HistogramData struct {
	Count       uint64            `json:"count"`
	Sum         float64           `json:"sum"`
	Buckets     []HistogramBucket `json:"buckets"`
	SampleCount uint64            `json:"sampleCount"`
}

// HistogramBucket represents a histogram bucket
type HistogramBucket struct {
	UpperBound float64 `json:"upperBound"`
	Count      uint64  `json:"count"`
}

// SummaryData represents summary metric data
type SummaryData struct {
	Count       uint64            `json:"count"`
	Sum         float64           `json:"sum"`
	Quantiles   []SummaryQuantile `json:"quantiles"`
	SampleCount uint64            `json:"sampleCount"`
}

// SummaryQuantile represents a summary quantile
type SummaryQuantile struct {
	Quantile float64 `json:"quantile"`
	Value    float64 `json:"value"`
}

// DataPoint represents a time series data point
type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// ExportFormat defines metric export formats
type ExportFormat string

const (
	ExportFormatPrometheus ExportFormat = "prometheus"
	ExportFormatJSON       ExportFormat = "json"
	ExportFormatCSV        ExportFormat = "csv"
	ExportFormatInfluxDB   ExportFormat = "influxdb"
)

// QueryResult represents metric query result
type QueryResult struct {
	Query         string            `json:"query"`
	Data          []QueryResultData `json:"data"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	ExecutionTime time.Duration     `json:"executionTime"`
	Timestamp     time.Time         `json:"timestamp"`
}

// QueryResultData represents individual query result data
type QueryResultData struct {
	Metric map[string]string `json:"metric"`
	Values []DataPoint       `json:"values"`
}

// DiskUsage represents disk usage information
type DiskUsage struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"usedPercent"`
}

// NetworkUsage represents network usage information
type NetworkUsage struct {
	BytesReceived   uint64 `json:"bytesReceived"`
	BytesSent       uint64 `json:"bytesSent"`
	PacketsReceived uint64 `json:"packetsReceived"`
	PacketsSent     uint64 `json:"packetsSent"`
	ErrorsReceived  uint64 `json:"errorsReceived"`
	ErrorsSent      uint64 `json:"errorsSent"`
}

// ConnectionPoolStats defined in resource_management.go

// types.TimeRange represents a time range
// types.TimeRange is defined in types/common.go

// PerformanceReport represents performance analysis results
type PerformanceReport struct {
	TimeRange            types.TimeRange              `json:"timeRange"`
	OverallPerformance   *OverallPerformance          `json:"overallPerformance"`
	OperationPerformance map[string]*OpPerformance    `json:"operationPerformance"`
	ResourceUsage        *ResourceUsage               `json:"resourceUsage"`
	Bottlenecks          []*Bottleneck                `json:"bottlenecks"`
	Recommendations      []*PerformanceRecommendation `json:"recommendations"`
	GeneratedAt          time.Time                    `json:"generatedAt"`
}

// OverallPerformance represents overall system performance
type OverallPerformance struct {
	AverageResponseTime time.Duration `json:"averageResponseTime"`
	P95ResponseTime     time.Duration `json:"p95ResponseTime"`
	P99ResponseTime     time.Duration `json:"p99ResponseTime"`
	Throughput          float64       `json:"throughput"`
	ErrorRate           float64       `json:"errorRate"`
	Availability        float64       `json:"availability"`
}

// OpPerformance represents operation-specific performance
type OpPerformance struct {
	OperationName       string        `json:"operationName"`
	CallCount           int64         `json:"callCount"`
	AverageResponseTime time.Duration `json:"averageResponseTime"`
	MinResponseTime     time.Duration `json:"minResponseTime"`
	MaxResponseTime     time.Duration `json:"maxResponseTime"`
	P50ResponseTime     time.Duration `json:"p50ResponseTime"`
	P95ResponseTime     time.Duration `json:"p95ResponseTime"`
	P99ResponseTime     time.Duration `json:"p99ResponseTime"`
	ErrorCount          int64         `json:"errorCount"`
	ErrorRate           float64       `json:"errorRate"`
	Throughput          float64       `json:"throughput"`
}

// ResourceUsage represents system resource usage
type ResourceUsage struct {
	CPU     *CPUUsage     `json:"cpu"`
	Memory  *MemoryUsage  `json:"memory"`
	Disk    *DiskUsage    `json:"disk"`
	Network *NetworkUsage `json:"network"`
}

// CPUUsage represents CPU usage information
type CPUUsage struct {
	Usage      float64 `json:"usage"`
	UserTime   float64 `json:"userTime"`
	SystemTime float64 `json:"systemTime"`
	IdleTime   float64 `json:"idleTime"`
	Cores      int     `json:"cores"`
	LoadAvg1   float64 `json:"loadAvg1"`
	LoadAvg5   float64 `json:"loadAvg5"`
	LoadAvg15  float64 `json:"loadAvg15"`
}

// MemoryUsage is defined in types/common.go

// Bottleneck represents a performance bottleneck
type Bottleneck struct {
	Type        BottleneckType     `json:"type"`
	Resource    string             `json:"resource"`
	Operation   string             `json:"operation,omitempty"`
	Severity    types.Severity     `json:"severity"`
	Impact      float64            `json:"impact"`
	Description string             `json:"description"`
	Metrics     map[string]float64 `json:"metrics"`
	DetectedAt  time.Time          `json:"detectedAt"`
}

// BottleneckType defines types of bottlenecks
type BottleneckType string

const (
	BottleneckTypeCPU      BottleneckType = "cpu"
	BottleneckTypeMemory   BottleneckType = "memory"
	BottleneckTypeDisk     BottleneckType = "disk"
	BottleneckTypeNetwork  BottleneckType = "network"
	BottleneckTypeDatabase BottleneckType = "database"
	BottleneckTypeQueue    BottleneckType = "queue"
	BottleneckTypeExternal BottleneckType = "external"
)

// Severity defined in types/common.go

// PerformanceRecommendation represents a performance improvement recommendation
type PerformanceRecommendation struct {
	Type        RecommendationType `json:"type"`
	Priority    types.Priority     `json:"priority"`
	Description string             `json:"description"`
	Action      string             `json:"action"`
	Impact      ImpactLevel        `json:"impact"`
	Effort      EffortLevel        `json:"effort"`
	Metrics     map[string]float64 `json:"metrics,omitempty"`
}

// RecommendationType defines types of recommendations
type RecommendationType string

const (
	RecommendationTypeScaling       RecommendationType = "scaling"
	RecommendationTypeOptimization  RecommendationType = "optimization"
	RecommendationTypeConfiguration RecommendationType = "configuration"
	RecommendationTypeArchitecture  RecommendationType = "architecture"
)

// Priority defines recommendation priorities
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// ImpactLevel defines impact levels
type ImpactLevel string

const (
	ImpactLevelLow    ImpactLevel = "low"
	ImpactLevelMedium ImpactLevel = "medium"
	ImpactLevelHigh   ImpactLevel = "high"
)

// EffortLevel defines effort levels
type EffortLevel string

const (
	EffortLevelLow    EffortLevel = "low"
	EffortLevelMedium EffortLevel = "medium"
	EffortLevelHigh   EffortLevel = "high"
)

// BottleneckAnalysis represents bottleneck analysis results
type BottleneckAnalysis struct {
	TimeRange   types.TimeRange    `json:"timeRange"`
	Bottlenecks []*Bottleneck      `json:"bottlenecks"`
	Summary     *BottleneckSummary `json:"summary"`
	Trends      []*BottleneckTrend `json:"trends"`
	GeneratedAt time.Time          `json:"generatedAt"`
}

// BottleneckSummary represents summary of bottlenecks
type BottleneckSummary struct {
	TotalBottlenecks      int                    `json:"totalBottlenecks"`
	BottlenecksByType     map[BottleneckType]int `json:"bottlenecksByType"`
	BottlenecksBySeverity map[types.Severity]int `json:"bottlenecksBySeverity"`
	MostCritical          *Bottleneck            `json:"mostCritical"`
	FrequentBottlenecks   []*Bottleneck          `json:"frequentBottlenecks"`
}

// BottleneckTrend represents bottleneck trends over time
type BottleneckTrend struct {
	Type       BottleneckType `json:"type"`
	Resource   string         `json:"resource"`
	Trend      TrendDirection `json:"trend"`
	Change     float64        `json:"change"`
	DataPoints []DataPoint    `json:"dataPoints"`
}

// TrendDirection defines trend directions
type TrendDirection string

const (
	TrendDirectionUp       TrendDirection = "up"
	TrendDirectionDown     TrendDirection = "down"
	TrendDirectionStable   TrendDirection = "stable"
	TrendDirectionVolatile TrendDirection = "volatile"
)

// TrendAnalysis represents trend analysis results
type TrendAnalysis struct {
	Metric      string          `json:"metric"`
	TimeRange   types.TimeRange `json:"timeRange"`
	Trend       TrendDirection  `json:"trend"`
	TrendScore  float64         `json:"trendScore"` // -1 to 1
	DataPoints  []DataPoint     `json:"dataPoints"`
	Forecast    []DataPoint     `json:"forecast,omitempty"`
	Seasonality *Seasonality    `json:"seasonality,omitempty"`
	Anomalies   []Anomaly       `json:"anomalies,omitempty"`
	GeneratedAt time.Time       `json:"generatedAt"`
}

// Seasonality represents seasonal patterns in data
type Seasonality struct {
	Detected bool              `json:"detected"`
	Period   time.Duration     `json:"period,omitempty"`
	Strength float64           `json:"strength,omitempty"`
	Patterns []SeasonalPattern `json:"patterns,omitempty"`
}

// SeasonalPattern represents a seasonal pattern
type SeasonalPattern struct {
	Period    time.Duration `json:"period"`
	Amplitude float64       `json:"amplitude"`
	Phase     float64       `json:"phase"`
}

// Anomaly represents an anomaly in data
type Anomaly struct {
	Timestamp   time.Time      `json:"timestamp"`
	Value       float64        `json:"value"`
	Expected    float64        `json:"expected"`
	Deviation   float64        `json:"deviation"`
	Confidence  float64        `json:"confidence"`
	Severity    types.Severity `json:"severity"`
	Description string         `json:"description,omitempty"`
}

// Threshold represents a threshold for alerting
type Threshold struct {
	Value    float64           `json:"value"`
	Operator ThresholdOperator `json:"operator"`
	Duration time.Duration     `json:"duration,omitempty"`
	Severity types.Severity    `json:"severity"`
}

// ThresholdOperator defines threshold operators
type ThresholdOperator string

const (
	ThresholdOperatorGreaterThan  ThresholdOperator = "gt"
	ThresholdOperatorGreaterEqual ThresholdOperator = "gte"
	ThresholdOperatorLessThan     ThresholdOperator = "lt"
	ThresholdOperatorLessEqual    ThresholdOperator = "lte"
	ThresholdOperatorEqual        ThresholdOperator = "eq"
	ThresholdOperatorNotEqual     ThresholdOperator = "neq"
)

// Alert represents an alert
type Alert struct {
	ID          string                 `json:"id"`
	RuleID      string                 `json:"ruleId"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Severity    types.Severity         `json:"severity"`
	Status      AlertStatus            `json:"status"`
	Labels      map[string]string      `json:"labels"`
	Annotations map[string]string      `json:"annotations"`
	StartTime   time.Time              `json:"startTime"`
	EndTime     *time.Time             `json:"endTime,omitempty"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	Value       float64                `json:"value"`
	Threshold   Threshold              `json:"threshold"`
	Context     map[string]interface{} `json:"context,omitempty"`
}

// AlertStatus defines alert statuses
type AlertStatus string

const (
	AlertStatusFiring       AlertStatus = "firing"
	AlertStatusResolved     AlertStatus = "resolved"
	AlertStatusAcknowledged AlertStatus = "acknowledged"
	AlertStatusSuppressed   AlertStatus = "suppressed"
)

// AlertRule defines an alert rule
type AlertRule struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Query       string            `json:"query"`
	Condition   AlertCondition    `json:"condition"`
	Severity    types.Severity    `json:"severity"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Enabled     bool              `json:"enabled"`
	Interval    time.Duration     `json:"interval"`
	For         time.Duration     `json:"for,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
	CreatedBy   string            `json:"createdBy"`
}

// AlertCondition defines conditions for triggering alerts
type AlertCondition struct {
	Threshold   Threshold     `json:"threshold"`
	Aggregation string        `json:"aggregation,omitempty"`
	TimeWindow  time.Duration `json:"timeWindow,omitempty"`
	GroupBy     []string      `json:"groupBy,omitempty"`
}

// RuleEvaluation represents an alert rule evaluation
type RuleEvaluation struct {
	RuleID    string        `json:"ruleId"`
	Timestamp time.Time     `json:"timestamp"`
	Value     float64       `json:"value"`
	Triggered bool          `json:"triggered"`
	Duration  time.Duration `json:"duration"`
	Error     string        `json:"error,omitempty"`
}

// NotificationChannel defines a notification channel
// NotificationChannel is defined in types/common.go

// AlertStats represents alert statistics
type AlertStats struct {
	TimeRange         types.TimeRange        `json:"timeRange"`
	TotalAlerts       int                    `json:"totalAlerts"`
	AlertsByStatus    map[AlertStatus]int    `json:"alertsByStatus"`
	AlertsBySeverity  map[types.Severity]int `json:"alertsBySeverity"`
	AlertsByRule      map[string]int         `json:"alertsByRule"`
	MTTA              time.Duration          `json:"mtta"` // Mean Time To Acknowledge
	MTTR              time.Duration          `json:"mttr"` // Mean Time To Resolve
	FalsePositiveRate float64                `json:"falsePositiveRate"`
}

// HealthCheckFunc defines a health check function
type HealthCheckFunc func(ctx context.Context) error

// HealthStatus represents health status - using types.HealthStatus instead

// SystemHealth represents overall system health
type SystemHealth struct {
	Status       types.HealthStatus             `json:"status"`
	Components   map[string]*types.HealthStatus `json:"components"`
	Dependencies map[string]*types.HealthStatus `json:"dependencies"`
	Summary      *HealthSummary                 `json:"summary"`
	LastUpdated  time.Time                      `json:"lastUpdated"`
}

// ServiceHealth represents service-specific health
type ServiceHealth struct {
	ServiceName  string                         `json:"serviceName"`
	Status       types.HealthStatus             `json:"status"`
	Instances    map[string]*types.HealthStatus `json:"instances"`
	LoadBalancer *types.HealthStatus            `json:"loadBalancer,omitempty"`
	Database     *types.HealthStatus            `json:"database,omitempty"`
	Cache        *types.HealthStatus            `json:"cache,omitempty"`
	MessageQueue *types.HealthStatus            `json:"messageQueue,omitempty"`
	Summary      *HealthSummary                 `json:"summary"`
	LastUpdated  time.Time                      `json:"lastUpdated"`
}

// HealthSummary represents health summary
type HealthSummary struct {
	TotalComponents     int                        `json:"totalComponents"`
	HealthyComponents   int                        `json:"healthyComponents"`
	UnhealthyComponents int                        `json:"unhealthyComponents"`
	DegradedComponents  int                        `json:"degradedComponents"`
	StatusDistribution  map[types.HealthStatus]int `json:"statusDistribution"`
	WorstStatus         types.HealthStatus         `json:"worstStatus"`
}

// HealthRecord represents a health check record - using types.HealthRecord instead

// HealthThreshold defines health thresholds
type HealthThreshold struct {
	ResponseTime        time.Duration `json:"responseTime"`
	ErrorRate           float64       `json:"errorRate"`
	Availability        float64       `json:"availability"`
	ConsecutiveFailures int           `json:"consecutiveFailures"`
}

// HealthAlert represents a health-related alert
type HealthAlert struct {
	ID            string             `json:"id"`
	ComponentName string             `json:"componentName"`
	Status        types.HealthStatus `json:"status"`
	Message       string             `json:"message"`
	Severity      types.Severity     `json:"severity"`
	Timestamp     time.Time          `json:"timestamp"`
	Resolved      bool               `json:"resolved"`
	ResolvedAt    *time.Time         `json:"resolvedAt,omitempty"`
}

// Callback function types - HealthCheckFunc defined above

// Additional monitoring types

// LogMetrics represents log-based metrics - defined in types/common.go
type LogMetrics = types.LogMetrics

// LogStats represents overall log statistics - defined in types/common.go
type LogStats = types.LogStats

// Trace represents a distributed trace
type Trace struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime"`
	Duration  time.Duration     `json:"duration"`
	Status    string            `json:"status"`
	Tags      map[string]string `json:"tags"`
	Spans     []*Span           `json:"spans"`
}

// Span represents a span within a trace
type Span struct {
	ID        string            `json:"id"`
	TraceID   string            `json:"traceId"`
	ParentID  string            `json:"parentId,omitempty"`
	Name      string            `json:"name"`
	StartTime time.Time         `json:"startTime"`
	EndTime   time.Time         `json:"endTime"`
	Duration  time.Duration     `json:"duration"`
	Tags      map[string]string `json:"tags"`
	Logs      []SpanLog         `json:"logs,omitempty"`
}

// SpanLog represents a log entry within a span
type SpanLog struct {
	Timestamp time.Time              `json:"timestamp"`
	Fields    map[string]interface{} `json:"fields"`
}

// TraceQuery represents a query for traces
type TraceQuery struct {
	Service     string            `json:"service,omitempty"`
	Operation   string            `json:"operation,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	MinDuration time.Duration     `json:"minDuration,omitempty"`
	MaxDuration time.Duration     `json:"maxDuration,omitempty"`
	TimeRange   types.TimeRange   `json:"timeRange"`
	Limit       int               `json:"limit,omitempty"`
}

// TraceAnalytics represents trace analytics results
type TraceAnalytics struct {
	TimeRange      types.TimeRange  `json:"timeRange"`
	TotalTraces    int              `json:"totalTraces"`
	AverageLatency time.Duration    `json:"averageLatency"`
	P95Latency     time.Duration    `json:"p95Latency"`
	P99Latency     time.Duration    `json:"p99Latency"`
	ErrorRate      float64          `json:"errorRate"`
	TopOperations  []OperationStats `json:"topOperations"`
	TopServices    []ServiceStats   `json:"topServices"`
}

// OperationStats represents statistics for an operation
type OperationStats struct {
	Name       string        `json:"name"`
	Count      int           `json:"count"`
	AvgLatency time.Duration `json:"avgLatency"`
	ErrorCount int           `json:"errorCount"`
	ErrorRate  float64       `json:"errorRate"`
}

// ServiceStats represents statistics for a service
type ServiceStats struct {
	Name       string        `json:"name"`
	Count      int           `json:"count"`
	AvgLatency time.Duration `json:"avgLatency"`
	ErrorCount int           `json:"errorCount"`
	ErrorRate  float64       `json:"errorRate"`
}

// ServiceMap represents service dependency map
type ServiceMap struct {
	Services    []ServiceNode `json:"services"`
	Connections []ServiceEdge `json:"connections"`
	GeneratedAt time.Time     `json:"generatedAt"`
}

// ServiceNode represents a service in the service map
type ServiceNode struct {
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	RequestCount int           `json:"requestCount"`
	ErrorRate    float64       `json:"errorRate"`
	AvgLatency   time.Duration `json:"avgLatency"`
}

// ServiceEdge represents a connection between services
type ServiceEdge struct {
	From         string        `json:"from"`
	To           string        `json:"to"`
	RequestCount int           `json:"requestCount"`
	ErrorRate    float64       `json:"errorRate"`
	AvgLatency   time.Duration `json:"avgLatency"`
}

// TraceDependencies represents trace dependency information
type TraceDependencies struct {
	Services     []string               `json:"services"`
	Dependencies map[string][]string    `json:"dependencies"`
	Metrics      map[string]interface{} `json:"metrics"`
	GeneratedAt  time.Time              `json:"generatedAt"`
}
