package interfaces

import (
	"context"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/types"
)

// ConfigManager defines the interface for configuration management
type ConfigManager interface {
	// Initialize initializes the config manager
	Initialize(ctx context.Context, options *ConfigOptions) error

	// LoadConfig loads configuration from sources
	LoadConfig(ctx context.Context) error

	// ReloadConfig reloads configuration from sources
	ReloadConfig(ctx context.Context) error

	// GetConfig returns the current configuration
	GetConfig(ctx context.Context) (interface{}, error)

	// GetConfigSection returns a specific section of configuration
	GetConfigSection(ctx context.Context, section string) (interface{}, error)

	// SetConfig sets configuration values
	SetConfig(ctx context.Context, config interface{}) error

	// SetConfigSection sets a specific section of configuration
	SetConfigSection(ctx context.Context, section string, value interface{}) error

	// ValidateConfig validates configuration
	ValidateConfig(ctx context.Context, config interface{}) error

	// WatchConfig watches for configuration changes
	WatchConfig(ctx context.Context, callback ConfigChangeCallback) error

	// GetConfigVersion returns current configuration version
	GetConfigVersion(ctx context.Context) (string, error)

	// GetConfigHistory returns configuration history
	GetConfigHistory(ctx context.Context, limit int) ([]ConfigVersion, error)

	// RollbackConfig rolls back to a previous configuration version
	RollbackConfig(ctx context.Context, version string) error

	// ExportConfig exports configuration in specified format
	ExportConfig(ctx context.Context, format ConfigFormat) ([]byte, error)

	// ImportConfig imports configuration from data
	ImportConfig(ctx context.Context, data []byte, format ConfigFormat) error

	// GetMetrics returns configuration manager metrics
	GetMetrics(ctx context.Context) (*ConfigMetrics, error)

	// HealthCheck performs configuration manager health check
	HealthCheck(ctx context.Context) error

	// Close closes the configuration manager
	Close() error
}

// ConfigProvider defines the interface for configuration providers
type ConfigProvider interface {
	// Name returns the provider name
	Name() string

	// Load loads configuration from the provider
	Load(ctx context.Context) (map[string]interface{}, error)

	// Watch watches for configuration changes
	Watch(ctx context.Context, callback ProviderChangeCallback) error

	// Set sets a configuration value (if writable)
	Set(ctx context.Context, key string, value interface{}) error

	// Delete deletes a configuration value (if writable)
	Delete(ctx context.Context, key string) error

	// IsWritable returns true if the provider supports writing
	IsWritable() bool

	// GetMetadata returns provider metadata
	GetMetadata() *ProviderMetadata

	// Close closes the provider
	Close() error
}

// ConfigValidator defines the interface for configuration validation
type ConfigValidator interface {
	// ValidateConfig validates the entire configuration
	ValidateConfig(ctx context.Context, config interface{}) (*types.ValidationResult, error)

	// ValidateSection validates a configuration section
	ValidateSection(ctx context.Context, section string, value interface{}) (*types.ValidationResult, error)

	// ValidateField validates a single configuration field
	ValidateField(ctx context.Context, field string, value interface{}) (*types.ValidationResult, error)

	// GetValidationRules returns validation rules
	GetValidationRules() []ValidationRule

	// AddValidationRule adds a custom validation rule
	AddValidationRule(rule ValidationRule) error

	// RemoveValidationRule removes a validation rule
	RemoveValidationRule(ruleID string) error
}

// ConfigTransformer defines the interface for configuration transformation
type ConfigTransformer interface {
	// Transform transforms configuration data
	Transform(ctx context.Context, input interface{}) (interface{}, error)

	// GetTransformationType returns the transformation type
	GetTransformationType() TransformationType
}

// ConfigEncryption defines the interface for configuration encryption
type ConfigEncryption interface {
	// EncryptValue encrypts a configuration value
	EncryptValue(ctx context.Context, value []byte) ([]byte, error)

	// DecryptValue decrypts a configuration value
	DecryptValue(ctx context.Context, encryptedValue []byte) ([]byte, error)

	// EncryptConfig encrypts sensitive fields in configuration
	EncryptConfig(ctx context.Context, config interface{}) (interface{}, error)

	// DecryptConfig decrypts sensitive fields in configuration
	DecryptConfig(ctx context.Context, config interface{}) (interface{}, error)

	// GetEncryptedFields returns list of fields that should be encrypted
	GetEncryptedFields() []string

	// SetEncryptedFields sets list of fields that should be encrypted
	SetEncryptedFields(fields []string) error
}

// ConfigAuditor defines the interface for configuration auditing
type ConfigAuditor interface {
	// AuditConfigChange audits a configuration change
	AuditConfigChange(ctx context.Context, event *ConfigChangeEvent) error

	// GetAuditLog returns configuration audit log
	GetAuditLog(ctx context.Context, filter *AuditFilter) ([]ConfigAuditEntry, error)

	// GetAuditStats returns audit statistics
	GetAuditStats(ctx context.Context, timeRange *types.TimeRange) (*ConfigAuditStats, error)
}

// ConfigBackup defines the interface for configuration backup and restore
type ConfigBackup interface {
	// CreateBackup creates a configuration backup
	CreateBackup(ctx context.Context, name string) (*BackupInfo, error)

	// RestoreBackup restores from a configuration backup
	RestoreBackup(ctx context.Context, backupID string) error

	// ListBackups lists all available backups
	ListBackups(ctx context.Context) ([]BackupInfo, error)

	// DeleteBackup deletes a configuration backup
	DeleteBackup(ctx context.Context, backupID string) error

	// GetBackupInfo returns information about a backup
	GetBackupInfo(ctx context.Context, backupID string) (*BackupInfo, error)

	// ExportBackup exports a backup to external storage
	ExportBackup(ctx context.Context, backupID string, destination string) error

	// ImportBackup imports a backup from external storage
	ImportBackup(ctx context.Context, source string) (*BackupInfo, error)
}

// Data structures

// ConfigOptions defines options for configuration manager initialization
type ConfigOptions struct {
	// Provider configurations
	Providers []ProviderConfig `json:"providers"`

	// Merge strategy for multiple providers
	MergeStrategy MergeStrategy `json:"mergeStrategy"`

	// Hot reload settings
	EnableHotReload bool          `json:"enableHotReload"`
	ReloadInterval  time.Duration `json:"reloadInterval,omitempty"`
	ReloadDebounce  time.Duration `json:"reloadDebounce,omitempty"`

	// Validation settings
	EnableValidation     bool   `json:"enableValidation"`
	ValidationSchemaPath string `json:"validationSchemaPath,omitempty"`
	StrictValidation     bool   `json:"strictValidation"`

	// Encryption settings
	EnableEncryption bool   `json:"enableEncryption"`
	EncryptionKey    string `json:"encryptionKey,omitempty"`
	EncryptionAlgo   string `json:"encryptionAlgo,omitempty"`

	// Backup settings
	EnableBackup      bool          `json:"enableBackup"`
	BackupInterval    time.Duration `json:"backupInterval,omitempty"`
	MaxBackups        int           `json:"maxBackups,omitempty"`
	BackupCompression bool          `json:"backupCompression"`

	// Audit settings
	EnableAudit  bool   `json:"enableAudit"`
	AuditLogPath string `json:"auditLogPath,omitempty"`

	// Performance settings
	CacheSize int           `json:"cacheSize"`
	CacheTTL  time.Duration `json:"cacheTtl"`

	// Security settings
	ReadOnly       bool     `json:"readOnly"`
	AllowedSources []string `json:"allowedSources,omitempty"`
	RestrictedKeys []string `json:"restrictedKeys,omitempty"`
}

// ProviderConfig defines configuration for a config provider
type ProviderConfig struct {
	Name     string                 `json:"name"`
	Type     ProviderType           `json:"type"`
	Priority int                    `json:"priority"`
	Config   map[string]interface{} `json:"config"`
	Enabled  bool                   `json:"enabled"`
}

// ProviderType defines types of configuration providers
type ProviderType string

const (
	ProviderTypeFile        ProviderType = "file"
	ProviderTypeEnvironment ProviderType = "environment"
	ProviderTypeConsul      ProviderType = "consul"
	ProviderTypeEtcd        ProviderType = "etcd"
	ProviderTypeVault       ProviderType = "vault"
	ProviderTypeDatabase    ProviderType = "database"
	ProviderTypeKubernetes  ProviderType = "kubernetes"
	ProviderTypeHTTP        ProviderType = "http"
	ProviderTypeS3          ProviderType = "s3"
)

// MergeStrategy defines strategies for merging configurations
type MergeStrategy string

const (
	MergeStrategyOverride MergeStrategy = "override"
	MergeStrategyMerge    MergeStrategy = "merge"
	MergeStrategyAppend   MergeStrategy = "append"
	MergeStrategyPriority MergeStrategy = "priority"
)

// ConfigFormat defines configuration formats
type ConfigFormat string

const (
	ConfigFormatJSON ConfigFormat = "json"
	ConfigFormatYAML ConfigFormat = "yaml"
	ConfigFormatTOML ConfigFormat = "toml"
	ConfigFormatINI  ConfigFormat = "ini"
	ConfigFormatXML  ConfigFormat = "xml"
	ConfigFormatHCL  ConfigFormat = "hcl"
)

// ConfigVersion represents a version of configuration
type ConfigVersion struct {
	ID          string                 `json:"id"`
	Version     string                 `json:"version"`
	Config      map[string]interface{} `json:"config"`
	Checksum    string                 `json:"checksum"`
	CreatedAt   time.Time              `json:"createdAt"`
	CreatedBy   string                 `json:"createdBy"`
	Description string                 `json:"description,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// ProviderMetadata contains metadata about a configuration provider
type ProviderMetadata struct {
	Name        string                 `json:"name"`
	Type        ProviderType           `json:"type"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Writable    bool                   `json:"writable"`
	Watchable   bool                   `json:"watchable"`
	Priority    int                    `json:"priority"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Status      ProviderStatus         `json:"status"`
	LastSync    time.Time              `json:"lastSync"`
	ErrorCount  int                    `json:"errorCount"`
}

// ProviderStatus defines provider status
type ProviderStatus string

const (
	ProviderStatusActive   ProviderStatus = "active"
	ProviderStatusInactive ProviderStatus = "inactive"
	ProviderStatusError    ProviderStatus = "error"
	ProviderStatusSyncing  ProviderStatus = "syncing"
)

// ValidationResult represents the result of configuration validation
// ValidationResult, ValidationError, and ValidationWarning are defined in types/common.go

// ValidationRule defines a validation rule
type ValidationRule struct {
	ID       string                 `json:"id"`
	Field    string                 `json:"field"`
	Type     ValidationType         `json:"type"`
	Required bool                   `json:"required"`
	MinValue interface{}            `json:"minValue,omitempty"`
	MaxValue interface{}            `json:"maxValue,omitempty"`
	Pattern  string                 `json:"pattern,omitempty"`
	Enum     []interface{}          `json:"enum,omitempty"`
	Custom   func(interface{}) bool `json:"-"`
	Message  string                 `json:"message"`
	Severity ValidationSeverity     `json:"severity"`
	Enabled  bool                   `json:"enabled"`
}

// ValidationType defines types of validation
type ValidationType string

const (
	ValidationTypeString  ValidationType = "string"
	ValidationTypeNumber  ValidationType = "number"
	ValidationTypeInteger ValidationType = "integer"
	ValidationTypeBoolean ValidationType = "boolean"
	ValidationTypeArray   ValidationType = "array"
	ValidationTypeObject  ValidationType = "object"
	ValidationTypeEmail   ValidationType = "email"
	ValidationTypeURL     ValidationType = "url"
	ValidationTypeRegex   ValidationType = "regex"
	ValidationTypeCustom  ValidationType = "custom"
)

// ValidationSeverity defines validation severity levels
type ValidationSeverity string

const (
	ValidationSeverityError   ValidationSeverity = "error"
	ValidationSeverityWarning ValidationSeverity = "warning"
	ValidationSeverityInfo    ValidationSeverity = "info"
)

// TransformationType defines types of configuration transformation
type TransformationType string

const (
	TransformationTypeTemplate     TransformationType = "template"
	TransformationTypeSubstitution TransformationType = "substitution"
	TransformationTypeMapping      TransformationType = "mapping"
	TransformationTypeFilter       TransformationType = "filter"
	TransformationTypeCustom       TransformationType = "custom"
)

// ConfigChangeEvent represents a configuration change event
type ConfigChangeEvent struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	Type       ConfigChangeType       `json:"type"`
	Source     string                 `json:"source"`
	Field      string                 `json:"field,omitempty"`
	OldValue   interface{}            `json:"oldValue,omitempty"`
	NewValue   interface{}            `json:"newValue,omitempty"`
	Version    string                 `json:"version"`
	User       string                 `json:"user,omitempty"`
	Reason     string                 `json:"reason,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// ConfigChangeType defines types of configuration changes
type ConfigChangeType string

const (
	ConfigChangeTypeCreate   ConfigChangeType = "create"
	ConfigChangeTypeUpdate   ConfigChangeType = "update"
	ConfigChangeTypeDelete   ConfigChangeType = "delete"
	ConfigChangeTypeReload   ConfigChangeType = "reload"
	ConfigChangeTypeRollback ConfigChangeType = "rollback"
)

// ConfigAuditEntry represents a configuration audit log entry
type ConfigAuditEntry struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	Action     string                 `json:"action"`
	User       string                 `json:"user"`
	Source     string                 `json:"source"`
	Target     string                 `json:"target"`
	Success    bool                   `json:"success"`
	Error      string                 `json:"error,omitempty"`
	Changes    []ConfigChange         `json:"changes,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// ConfigChange represents a single configuration change
type ConfigChange struct {
	Field    string      `json:"field"`
	OldValue interface{} `json:"oldValue,omitempty"`
	NewValue interface{} `json:"newValue,omitempty"`
	Action   string      `json:"action"`
}

// AuditFilter defines filter for configuration audit log
type AuditFilter struct {
	User     string     `json:"user,omitempty"`
	Action   string     `json:"action,omitempty"`
	Source   string     `json:"source,omitempty"`
	Target   string     `json:"target,omitempty"`
	Success  *bool      `json:"success,omitempty"`
	TimeFrom *time.Time `json:"timeFrom,omitempty"`
	TimeTo   *time.Time `json:"timeTo,omitempty"`
	Limit    int        `json:"limit,omitempty"`
	Offset   int        `json:"offset,omitempty"`
}

// ConfigAuditStats represents configuration audit statistics
type ConfigAuditStats struct {
	TotalChanges    int64               `json:"totalChanges"`
	ChangesByUser   map[string]int64    `json:"changesByUser"`
	ChangesByAction map[string]int64    `json:"changesByAction"`
	ChangesBySource map[string]int64    `json:"changesBySource"`
	SuccessRate     float64             `json:"successRate"`
	TopUsers        []UserChangeCount   `json:"topUsers"`
	TopSources      []SourceChangeCount `json:"topSources"`
	TimeRange       types.TimeRange     `json:"timeRange"`
	GeneratedAt     time.Time           `json:"generatedAt"`
}

// UserChangeCount represents user change statistics
type UserChangeCount struct {
	User  string `json:"user"`
	Count int64  `json:"count"`
}

// SourceChangeCount represents source change statistics
type SourceChangeCount struct {
	Source string `json:"source"`
	Count  int64  `json:"count"`
}

// BackupInfo represents information about a configuration backup
type BackupInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Version     string                 `json:"version"`
	Size        int64                  `json:"size"`
	Checksum    string                 `json:"checksum"`
	Compressed  bool                   `json:"compressed"`
	Encrypted   bool                   `json:"encrypted"`
	CreatedAt   time.Time              `json:"createdAt"`
	CreatedBy   string                 `json:"createdBy"`
	ExpiresAt   *time.Time             `json:"expiresAt,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Status      BackupStatus           `json:"status"`
}

// BackupStatus defines backup status
type BackupStatus string

const (
	BackupStatusCreating  BackupStatus = "creating"
	BackupStatusCompleted BackupStatus = "completed"
	BackupStatusFailed    BackupStatus = "failed"
	BackupStatusExpired   BackupStatus = "expired"
	BackupStatusCorrupted BackupStatus = "corrupted"
)

// ConfigMetrics defines configuration manager metrics
type ConfigMetrics struct {
	// Configuration metrics
	TotalConfigurations  int            `json:"totalConfigurations"`
	ConfigurationsByType map[string]int `json:"configurationsByType"`
	ConfigSize           int64          `json:"configSize"`
	ConfigVersion        string         `json:"configVersion"`

	// Provider metrics
	ActiveProviders    int                    `json:"activeProviders"`
	ProvidersByType    map[ProviderType]int   `json:"providersByType"`
	ProvidersByStatus  map[ProviderStatus]int `json:"providersByStatus"`
	ProviderSyncErrors uint64                 `json:"providerSyncErrors"`

	// Change metrics
	TotalChanges      uint64                      `json:"totalChanges"`
	ChangesInLastHour uint64                      `json:"changesInLastHour"`
	ChangesInLastDay  uint64                      `json:"changesInLastDay"`
	ChangesByType     map[ConfigChangeType]uint64 `json:"changesByType"`

	// Validation metrics
	ValidationAttempts uint64 `json:"validationAttempts"`
	ValidationFailures uint64 `json:"validationFailures"`
	ValidationErrors   uint64 `json:"validationErrors"`
	ValidationWarnings uint64 `json:"validationWarnings"`

	// Performance metrics
	AverageLoadTime       time.Duration `json:"averageLoadTime"`
	AverageValidationTime time.Duration `json:"averageValidationTime"`
	CacheHitRate          float64       `json:"cacheHitRate"`

	// Backup metrics
	TotalBackups   int        `json:"totalBackups"`
	BackupSize     int64      `json:"backupSize"`
	LastBackupTime *time.Time `json:"lastBackupTime,omitempty"`
	BackupErrors   uint64     `json:"backupErrors"`

	// Error metrics
	TotalErrors  uint64            `json:"totalErrors"`
	ErrorsByType map[string]uint64 `json:"errorsByType"`

	Timestamp time.Time `json:"timestamp"`
}

// Callback function types
type ConfigChangeCallback func(event ConfigChangeEvent)
type ProviderChangeCallback func(provider string, changes map[string]interface{})

// TimeRange represents a time range (reused from other interfaces)
// TimeRange is defined in types/common.go
