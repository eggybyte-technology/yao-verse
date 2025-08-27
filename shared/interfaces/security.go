package interfaces

import (
	"context"
	"crypto"
	"crypto/tls"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/types"
)

// AuthenticationProvider defines interface for authentication providers
type AuthenticationProvider interface {
	// Authenticate authenticates a user with credentials
	Authenticate(ctx context.Context, credentials *Credentials) (*AuthResult, error)

	// ValidateToken validates an authentication token
	ValidateToken(ctx context.Context, token string) (*TokenInfo, error)

	// RefreshToken refreshes an authentication token
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error)

	// RevokeToken revokes an authentication token
	RevokeToken(ctx context.Context, token string) error

	// GetUserInfo returns user information for a token
	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)

	// HealthCheck performs authentication provider health check
	HealthCheck(ctx context.Context) error
}

// AuthorizationProvider defines interface for authorization providers
type AuthorizationProvider interface {
	// Authorize checks if a user has permission for a resource/action
	Authorize(ctx context.Context, subject *Subject, resource string, action string) (bool, error)

	// GetPermissions returns all permissions for a subject
	GetPermissions(ctx context.Context, subject *Subject) ([]Permission, error)

	// GrantPermission grants a permission to a subject
	GrantPermission(ctx context.Context, subject *Subject, permission Permission) error

	// RevokePermission revokes a permission from a subject
	RevokePermission(ctx context.Context, subject *Subject, permission Permission) error

	// CreateRole creates a new role
	CreateRole(ctx context.Context, role *Role) error

	// DeleteRole deletes a role
	DeleteRole(ctx context.Context, roleID string) error

	// AssignRole assigns a role to a subject
	AssignRole(ctx context.Context, subject *Subject, roleID string) error

	// UnassignRole unassigns a role from a subject
	UnassignRole(ctx context.Context, subject *Subject, roleID string) error

	// GetRoles returns all roles for a subject
	GetRoles(ctx context.Context, subject *Subject) ([]Role, error)

	// HealthCheck performs authorization provider health check
	HealthCheck(ctx context.Context) error
}

// CryptoProvider defines interface for cryptographic operations
type CryptoProvider interface {
	// GenerateKeyPair generates a public/private key pair
	GenerateKeyPair(ctx context.Context, algorithm string, keySize int) (*KeyPair, error)

	// Sign signs data with a private key
	Sign(ctx context.Context, data []byte, privateKey crypto.PrivateKey) ([]byte, error)

	// Verify verifies a signature with a public key
	Verify(ctx context.Context, data []byte, signature []byte, publicKey crypto.PublicKey) (bool, error)

	// Encrypt encrypts data with a public key
	Encrypt(ctx context.Context, data []byte, publicKey crypto.PublicKey) ([]byte, error)

	// Decrypt decrypts data with a private key
	Decrypt(ctx context.Context, encryptedData []byte, privateKey crypto.PrivateKey) ([]byte, error)

	// Hash computes a hash of data
	Hash(ctx context.Context, data []byte, algorithm string) ([]byte, error)

	// GenerateRandomBytes generates cryptographically secure random bytes
	GenerateRandomBytes(ctx context.Context, length int) ([]byte, error)

	// DeriveKey derives a key from a password using PBKDF2
	DeriveKey(ctx context.Context, password []byte, salt []byte, iterations int, keyLength int) ([]byte, error)

	// HealthCheck performs crypto provider health check
	HealthCheck(ctx context.Context) error
}

// CertificateManager defines interface for certificate management
type CertificateManager interface {
	// GenerateCertificate generates a new certificate
	GenerateCertificate(ctx context.Context, request *CertificateRequest) (*Certificate, error)

	// GetCertificate retrieves a certificate by ID
	GetCertificate(ctx context.Context, certID string) (*Certificate, error)

	// ListCertificates lists all certificates
	ListCertificates(ctx context.Context, filter *CertificateFilter) ([]Certificate, error)

	// RevokeCertificate revokes a certificate
	RevokeCertificate(ctx context.Context, certID string, reason RevocationReason) error

	// ValidateCertificate validates a certificate
	ValidateCertificate(ctx context.Context, cert *Certificate) (*ValidationResult, error)

	// RenewCertificate renews an existing certificate
	RenewCertificate(ctx context.Context, certID string) (*Certificate, error)

	// ExportCertificate exports a certificate in specified format
	ExportCertificate(ctx context.Context, certID string, format CertificateFormat) ([]byte, error)

	// ImportCertificate imports a certificate
	ImportCertificate(ctx context.Context, certData []byte, format CertificateFormat) (*Certificate, error)

	// GetCertificateChain returns the certificate chain for a certificate
	GetCertificateChain(ctx context.Context, certID string) ([]*Certificate, error)

	// HealthCheck performs certificate manager health check
	HealthCheck(ctx context.Context) error
}

// SecretManager defines interface for secret management
type SecretManager interface {
	// StoreSecret stores a secret
	StoreSecret(ctx context.Context, key string, secret *Secret) error

	// GetSecret retrieves a secret
	GetSecret(ctx context.Context, key string) (*Secret, error)

	// DeleteSecret deletes a secret
	DeleteSecret(ctx context.Context, key string) error

	// ListSecrets lists all secret keys
	ListSecrets(ctx context.Context, filter *SecretFilter) ([]string, error)

	// RotateSecret rotates a secret
	RotateSecret(ctx context.Context, key string) (*Secret, error)

	// CreateSecretVersion creates a new version of a secret
	CreateSecretVersion(ctx context.Context, key string, secret *Secret) (int, error)

	// GetSecretVersion retrieves a specific version of a secret
	GetSecretVersion(ctx context.Context, key string, version int) (*Secret, error)

	// SetSecretExpiry sets expiry time for a secret
	SetSecretExpiry(ctx context.Context, key string, expiry time.Time) error

	// HealthCheck performs secret manager health check
	HealthCheck(ctx context.Context) error
}

// TLSManager defines interface for TLS certificate and connection management
type TLSManager interface {
	// GetTLSConfig returns TLS configuration
	GetTLSConfig(ctx context.Context, serverName string) (*tls.Config, error)

	// GetClientTLSConfig returns client TLS configuration
	GetClientTLSConfig(ctx context.Context, serverName string) (*tls.Config, error)

	// LoadCertificate loads a certificate from file or store
	LoadCertificate(ctx context.Context, certPath string, keyPath string) (*tls.Certificate, error)

	// CreateTLSCertificate creates a TLS certificate
	CreateTLSCertificate(ctx context.Context, request *TLSCertificateRequest) (*tls.Certificate, error)

	// ValidateTLSConnection validates a TLS connection
	ValidateTLSConnection(ctx context.Context, conn *tls.Conn) error

	// GetCipherSuites returns supported cipher suites
	GetCipherSuites() []uint16

	// HealthCheck performs TLS manager health check
	HealthCheck(ctx context.Context) error
}

// AuditLogger defines interface for security audit logging
type AuditLogger interface {
	// LogAuthAttempt logs an authentication attempt
	LogAuthAttempt(ctx context.Context, event *AuthEvent) error

	// LogAuthorizationDecision logs an authorization decision
	LogAuthorizationDecision(ctx context.Context, event *AuthzEvent) error

	// LogSecurityEvent logs a general security event
	LogSecurityEvent(ctx context.Context, event *SecurityEvent) error

	// LogCertificateEvent logs a certificate-related event
	LogCertificateEvent(ctx context.Context, event *CertificateEvent) error

	// QueryAuditLogs queries audit logs
	QueryAuditLogs(ctx context.Context, query *AuditQuery) ([]AuditLogEntry, error)

	// GetAuditStats returns audit statistics
	GetAuditStats(ctx context.Context, timeRange *types.TimeRange) (*AuditStats, error)

	// HealthCheck performs audit logger health check
	HealthCheck(ctx context.Context) error
}

// Data structures

// Credentials represents authentication credentials
type Credentials struct {
	Type       CredentialType         `json:"type"`
	Username   string                 `json:"username,omitempty"`
	Password   string                 `json:"password,omitempty"`
	Token      string                 `json:"token,omitempty"`
	ApiKey     string                 `json:"apiKey,omitempty"`
	ClientCert *tls.Certificate       `json:"clientCert,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// CredentialType defines types of credentials
type CredentialType string

const (
	CredentialTypePassword   CredentialType = "password"
	CredentialTypeToken      CredentialType = "token"
	CredentialTypeApiKey     CredentialType = "api_key"
	CredentialTypeClientCert CredentialType = "client_cert"
	CredentialTypeOAuth      CredentialType = "oauth"
	CredentialTypeJWT        CredentialType = "jwt"
	CredentialTypeSAML       CredentialType = "saml"
	CredentialTypeOIDC       CredentialType = "oidc"
)

// AuthResult represents the result of authentication
type AuthResult struct {
	Success      bool          `json:"success"`
	Token        string        `json:"token,omitempty"`
	RefreshToken string        `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time     `json:"expiresAt,omitempty"`
	UserInfo     *UserInfo     `json:"userInfo,omitempty"`
	Error        string        `json:"error,omitempty"`
	MFA          *MFAChallenge `json:"mfa,omitempty"`
}

// TokenInfo represents information about a token
type TokenInfo struct {
	Valid     bool      `json:"valid"`
	ExpiresAt time.Time `json:"expiresAt"`
	Subject   string    `json:"subject"`
	Issuer    string    `json:"issuer"`
	Audience  string    `json:"audience"`
	Scopes    []string  `json:"scopes"`
}

// UserInfo represents user information
type UserInfo struct {
	ID          string                 `json:"id"`
	Username    string                 `json:"username"`
	Email       string                 `json:"email,omitempty"`
	DisplayName string                 `json:"displayName,omitempty"`
	Groups      []string               `json:"groups,omitempty"`
	Roles       []string               `json:"roles,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// Subject represents an entity that can be authorized
type Subject struct {
	Type       SubjectType            `json:"type"`
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Groups     []string               `json:"groups,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

// SubjectType defines types of subjects
type SubjectType string

const (
	SubjectTypeUser    SubjectType = "user"
	SubjectTypeService SubjectType = "service"
	SubjectTypeDevice  SubjectType = "device"
	SubjectTypeGroup   SubjectType = "group"
	SubjectTypeRole    SubjectType = "role"
)

// Permission represents a permission
type Permission struct {
	ID         string                 `json:"id"`
	Resource   string                 `json:"resource"`
	Action     string                 `json:"action"`
	Effect     PermissionEffect       `json:"effect"`
	Conditions []PermissionCondition  `json:"conditions,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}

// PermissionEffect defines the effect of a permission
type PermissionEffect string

const (
	PermissionEffectAllow PermissionEffect = "allow"
	PermissionEffectDeny  PermissionEffect = "deny"
)

// PermissionCondition represents a condition for a permission
type PermissionCondition struct {
	Type      string      `json:"type"`
	Attribute string      `json:"attribute"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
}

// Role represents a role
type Role struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Permissions []Permission           `json:"permissions"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// KeyPair represents a public/private key pair
type KeyPair struct {
	Algorithm  string            `json:"algorithm"`
	KeySize    int               `json:"keySize"`
	PublicKey  crypto.PublicKey  `json:"publicKey"`
	PrivateKey crypto.PrivateKey `json:"privateKey"`
	CreatedAt  time.Time         `json:"createdAt"`
}

// Certificate represents a digital certificate
type Certificate struct {
	ID           string                 `json:"id"`
	Subject      string                 `json:"subject"`
	Issuer       string                 `json:"issuer"`
	SerialNumber string                 `json:"serialNumber"`
	NotBefore    time.Time              `json:"notBefore"`
	NotAfter     time.Time              `json:"notAfter"`
	KeyUsage     []string               `json:"keyUsage"`
	SANs         []string               `json:"sans"`
	CertData     []byte                 `json:"certData"`
	Status       CertificateStatus      `json:"status"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// CertificateRequest represents a certificate signing request
type CertificateRequest struct {
	Subject      string                 `json:"subject"`
	SANs         []string               `json:"sans,omitempty"`
	KeyUsage     []string               `json:"keyUsage"`
	ValidityDays int                    `json:"validityDays"`
	KeyAlgorithm string                 `json:"keyAlgorithm"`
	KeySize      int                    `json:"keySize"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
}

// CertificateStatus defines certificate status
type CertificateStatus string

const (
	CertificateStatusActive    CertificateStatus = "active"
	CertificateStatusExpired   CertificateStatus = "expired"
	CertificateStatusRevoked   CertificateStatus = "revoked"
	CertificateStatusSuspended CertificateStatus = "suspended"
)

// CertificateFormat defines certificate formats
type CertificateFormat string

const (
	CertificateFormatPEM    CertificateFormat = "pem"
	CertificateFormatDER    CertificateFormat = "der"
	CertificateFormatPKCS12 CertificateFormat = "pkcs12"
	CertificateFormatJWK    CertificateFormat = "jwk"
)

// RevocationReason defines reasons for certificate revocation
type RevocationReason string

const (
	RevocationReasonUnspecified          RevocationReason = "unspecified"
	RevocationReasonKeyCompromise        RevocationReason = "key_compromise"
	RevocationReasonCACompromise         RevocationReason = "ca_compromise"
	RevocationReasonAffiliationChanged   RevocationReason = "affiliation_changed"
	RevocationReasonSuperseded           RevocationReason = "superseded"
	RevocationReasonCessationOfOperation RevocationReason = "cessation_of_operation"
	RevocationReasonCertificateHold      RevocationReason = "certificate_hold"
)

// ValidationResult represents certificate validation result
type ValidationResult struct {
	Valid       bool                   `json:"valid"`
	Errors      []string               `json:"errors,omitempty"`
	Warnings    []string               `json:"warnings,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	ValidatedAt time.Time              `json:"validatedAt"`
}

// CertificateFilter defines filter for certificate listing
type CertificateFilter struct {
	Subject    string            `json:"subject,omitempty"`
	Issuer     string            `json:"issuer,omitempty"`
	Status     CertificateStatus `json:"status,omitempty"`
	ExpiryFrom *time.Time        `json:"expiryFrom,omitempty"`
	ExpiryTo   *time.Time        `json:"expiryTo,omitempty"`
	KeyUsage   []string          `json:"keyUsage,omitempty"`
	Limit      int               `json:"limit,omitempty"`
	Offset     int               `json:"offset,omitempty"`
}

// Secret represents a secret
type Secret struct {
	Key         string                 `json:"key"`
	Value       []byte                 `json:"value"`
	Version     int                    `json:"version"`
	Type        SecretType             `json:"type"`
	Description string                 `json:"description,omitempty"`
	ExpiresAt   *time.Time             `json:"expiresAt,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

// SecretType defines types of secrets
type SecretType string

const (
	SecretTypeGeneric  SecretType = "generic"
	SecretTypePassword SecretType = "password"
	SecretTypeApiKey   SecretType = "api_key"
	SecretTypeCert     SecretType = "certificate"
	SecretTypeKey      SecretType = "private_key"
	SecretTypeToken    SecretType = "token"
	SecretTypeDatabase SecretType = "database"
)

// SecretFilter defines filter for secret listing
type SecretFilter struct {
	Type       SecretType `json:"type,omitempty"`
	ExpiryFrom *time.Time `json:"expiryFrom,omitempty"`
	ExpiryTo   *time.Time `json:"expiryTo,omitempty"`
	Limit      int        `json:"limit,omitempty"`
	Offset     int        `json:"offset,omitempty"`
}

// TLSCertificateRequest represents a TLS certificate request
type TLSCertificateRequest struct {
	CommonName   string   `json:"commonName"`
	Organization string   `json:"organization,omitempty"`
	Country      string   `json:"country,omitempty"`
	SANs         []string `json:"sans,omitempty"`
	ValidityDays int      `json:"validityDays"`
	KeyAlgorithm string   `json:"keyAlgorithm"`
	KeySize      int      `json:"keySize"`
}

// MFAChallenge represents a multi-factor authentication challenge
type MFAChallenge struct {
	Type      MFAType   `json:"type"`
	Challenge string    `json:"challenge"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// MFAType defines types of MFA challenges
type MFAType string

const (
	MFATypeTOTP  MFAType = "totp"
	MFATypeSMS   MFAType = "sms"
	MFATypeEmail MFAType = "email"
	MFATypePush  MFAType = "push"
)

// Audit event structures

// AuthEvent represents an authentication event
type AuthEvent struct {
	EventID       string                 `json:"eventId"`
	Timestamp     time.Time              `json:"timestamp"`
	Type          AuthEventType          `json:"type"`
	Subject       string                 `json:"subject"`
	Success       bool                   `json:"success"`
	FailureReason string                 `json:"failureReason,omitempty"`
	ClientIP      string                 `json:"clientIp"`
	UserAgent     string                 `json:"userAgent,omitempty"`
	SessionID     string                 `json:"sessionId,omitempty"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}

// AuthEventType defines types of authentication events
type AuthEventType string

const (
	AuthEventTypeLogin        AuthEventType = "login"
	AuthEventTypeLogout       AuthEventType = "logout"
	AuthEventTypeTokenRefresh AuthEventType = "token_refresh"
	AuthEventTypeTokenRevoke  AuthEventType = "token_revoke"
	AuthEventTypeMFAChallenge AuthEventType = "mfa_challenge"
	AuthEventTypeMFAVerify    AuthEventType = "mfa_verify"
)

// AuthzEvent represents an authorization event
type AuthzEvent struct {
	EventID    string                 `json:"eventId"`
	Timestamp  time.Time              `json:"timestamp"`
	Subject    string                 `json:"subject"`
	Resource   string                 `json:"resource"`
	Action     string                 `json:"action"`
	Decision   string                 `json:"decision"`
	Reason     string                 `json:"reason,omitempty"`
	ClientIP   string                 `json:"clientIp"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// SecurityEvent represents a general security event
type SecurityEvent struct {
	EventID     string                 `json:"eventId"`
	Timestamp   time.Time              `json:"timestamp"`
	Type        SecurityEventType      `json:"type"`
	Severity    SecurityEventSeverity  `json:"severity"`
	Subject     string                 `json:"subject,omitempty"`
	Resource    string                 `json:"resource,omitempty"`
	Action      string                 `json:"action,omitempty"`
	Description string                 `json:"description"`
	ClientIP    string                 `json:"clientIp,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

// SecurityEventType defines types of security events
type SecurityEventType string

const (
	SecurityEventTypeSecurityBreach      SecurityEventType = "security_breach"
	SecurityEventTypeUnauthorizedAccess  SecurityEventType = "unauthorized_access"
	SecurityEventTypeSuspiciousActivity  SecurityEventType = "suspicious_activity"
	SecurityEventTypeConfigurationChange SecurityEventType = "configuration_change"
	SecurityEventTypeSystemCompromise    SecurityEventType = "system_compromise"
)

// SecurityEventSeverity defines severity levels for security events
type SecurityEventSeverity string

const (
	SecurityEventSeverityLow      SecurityEventSeverity = "low"
	SecurityEventSeverityMedium   SecurityEventSeverity = "medium"
	SecurityEventSeverityHigh     SecurityEventSeverity = "high"
	SecurityEventSeverityCritical SecurityEventSeverity = "critical"
)

// CertificateEvent represents a certificate-related event
type CertificateEvent struct {
	EventID       string                 `json:"eventId"`
	Timestamp     time.Time              `json:"timestamp"`
	Type          CertificateEventType   `json:"type"`
	CertificateID string                 `json:"certificateId"`
	Subject       string                 `json:"subject"`
	Issuer        string                 `json:"issuer"`
	Action        string                 `json:"action"`
	Result        string                 `json:"result"`
	Reason        string                 `json:"reason,omitempty"`
	ClientIP      string                 `json:"clientIp,omitempty"`
	Properties    map[string]interface{} `json:"properties,omitempty"`
}

// CertificateEventType defines types of certificate events
type CertificateEventType string

const (
	CertificateEventTypeIssued    CertificateEventType = "issued"
	CertificateEventTypeRevoked   CertificateEventType = "revoked"
	CertificateEventTypeRenewed   CertificateEventType = "renewed"
	CertificateEventTypeExpired   CertificateEventType = "expired"
	CertificateEventTypeValidated CertificateEventType = "validated"
)

// AuditLogEntry represents an audit log entry
type AuditLogEntry struct {
	ID         string                 `json:"id"`
	Timestamp  time.Time              `json:"timestamp"`
	EventType  string                 `json:"eventType"`
	Subject    string                 `json:"subject"`
	Resource   string                 `json:"resource,omitempty"`
	Action     string                 `json:"action"`
	Result     string                 `json:"result"`
	ClientIP   string                 `json:"clientIp,omitempty"`
	UserAgent  string                 `json:"userAgent,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// AuditQuery represents a query for audit logs
type AuditQuery struct {
	EventType string     `json:"eventType,omitempty"`
	Subject   string     `json:"subject,omitempty"`
	Resource  string     `json:"resource,omitempty"`
	Action    string     `json:"action,omitempty"`
	Result    string     `json:"result,omitempty"`
	ClientIP  string     `json:"clientIp,omitempty"`
	TimeFrom  *time.Time `json:"timeFrom,omitempty"`
	TimeTo    *time.Time `json:"timeTo,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	Offset    int        `json:"offset,omitempty"`
	SortBy    string     `json:"sortBy,omitempty"`
	SortOrder string     `json:"sortOrder,omitempty"`
}

// AuditStats represents audit statistics
type AuditStats struct {
	TotalEvents     int64                 `json:"totalEvents"`
	EventsByType    map[string]int64      `json:"eventsByType"`
	EventsByResult  map[string]int64      `json:"eventsByResult"`
	EventsBySubject map[string]int64      `json:"eventsBySubject"`
	TopResources    []ResourceAccessCount `json:"topResources"`
	TopClients      []ClientAccessCount   `json:"topClients"`
	FailureRate     float64               `json:"failureRate"`
	TimeRange       types.TimeRange       `json:"timeRange"`
	GeneratedAt     time.Time             `json:"generatedAt"`
}

// ResourceAccessCount represents resource access statistics
type ResourceAccessCount struct {
	Resource string `json:"resource"`
	Count    int64  `json:"count"`
}

// ClientAccessCount represents client access statistics
type ClientAccessCount struct {
	ClientIP  string `json:"clientIp"`
	Count     int64  `json:"count"`
	UserAgent string `json:"userAgent,omitempty"`
}
