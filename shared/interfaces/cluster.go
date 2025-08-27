package interfaces

import (
	"context"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/types"
)

// ClusterManager defines the interface for cluster management
type ClusterManager interface {
	// Initialize initializes the cluster manager
	Initialize(ctx context.Context, config *ClusterConfig) error

	// Join joins the cluster
	Join(ctx context.Context) error

	// Leave leaves the cluster gracefully
	Leave(ctx context.Context) error

	// GetNodes returns all nodes in the cluster
	GetNodes(ctx context.Context) ([]Node, error)

	// GetActiveNodes returns only active nodes in the cluster
	GetActiveNodes(ctx context.Context) ([]Node, error)

	// GetNode returns information about a specific node
	GetNode(ctx context.Context, nodeID string) (*Node, error)

	// RemoveNode removes a node from the cluster
	RemoveNode(ctx context.Context, nodeID string) error

	// GetClusterInfo returns overall cluster information
	GetClusterInfo(ctx context.Context) (*ClusterInfo, error)

	// RegisterNodeChangeCallback registers callback for node changes
	RegisterNodeChangeCallback(callback NodeChangeCallback) error

	// GetMetrics returns cluster metrics
	GetMetrics(ctx context.Context) (*ClusterMetrics, error)

	// HealthCheck performs cluster health check
	HealthCheck(ctx context.Context) error
}

// LeaderElector defines the interface for leader election
type LeaderElector interface {
	// Start starts the leader election process
	Start(ctx context.Context) error

	// Stop stops the leader election process
	Stop(ctx context.Context) error

	// IsLeader returns true if this node is the current leader
	IsLeader() bool

	// GetLeader returns information about the current leader
	GetLeader(ctx context.Context) (*LeaderInfo, error)

	// StepDown causes the current leader to step down
	StepDown(ctx context.Context) error

	// RegisterLeadershipCallback registers callback for leadership changes
	RegisterLeadershipCallback(callback LeadershipCallback) error

	// GetElectionStatus returns the current election status
	GetElectionStatus() *ElectionStatus

	// GetMetrics returns election metrics
	GetMetrics() *ElectionMetrics

	// HealthCheck performs leader election health check
	HealthCheck(ctx context.Context) error
}

// NodeDiscovery defines the interface for node discovery
type NodeDiscovery interface {
	// Start starts the node discovery service
	Start(ctx context.Context) error

	// Stop stops the node discovery service
	Stop(ctx context.Context) error

	// RegisterNode registers this node with the discovery service
	RegisterNode(ctx context.Context, node *NodeRegistration) error

	// UnregisterNode unregisters this node from the discovery service
	UnregisterNode(ctx context.Context) error

	// DiscoverNodes discovers nodes in the cluster
	DiscoverNodes(ctx context.Context, filter *NodeFilter) ([]Node, error)

	// WatchNodes watches for node changes
	WatchNodes(ctx context.Context, callback NodeChangeCallback) error

	// UpdateNodeStatus updates the status of this node
	UpdateNodeStatus(ctx context.Context, status NodeStatus) error

	// GetMetrics returns discovery metrics
	GetMetrics() *DiscoveryMetrics

	// HealthCheck performs node discovery health check
	HealthCheck(ctx context.Context) error
}

// LoadBalancer defines the interface for load balancing
type LoadBalancer interface {
	// SelectNode selects a node for load balancing
	SelectNode(ctx context.Context, criteria *SelectionCriteria) (*Node, error)

	// UpdateNodeWeights updates the weights for load balancing
	UpdateNodeWeights(ctx context.Context, weights map[string]int) error

	// GetNodeWeights returns current node weights
	GetNodeWeights(ctx context.Context) (map[string]int, error)

	// SetStrategy sets the load balancing strategy
	SetStrategy(strategy LoadBalancingStrategy) error

	// GetStrategy returns the current load balancing strategy
	GetStrategy() LoadBalancingStrategy

	// GetMetrics returns load balancer metrics
	GetMetrics() *LoadBalancerMetrics

	// HealthCheck performs load balancer health check
	HealthCheck(ctx context.Context) error
}

// HealthMonitor defines the interface for cluster health monitoring
type HealthMonitor interface {
	// Start starts health monitoring
	Start(ctx context.Context) error

	// Stop stops health monitoring
	Stop(ctx context.Context) error

	// CheckNodeHealth checks the health of a specific node
	CheckNodeHealth(ctx context.Context, nodeID string) (*types.HealthStatus, error)

	// CheckClusterHealth checks the overall cluster health
	CheckClusterHealth(ctx context.Context) (*ClusterHealthStatus, error)

	// RegisterHealthCallback registers callback for health changes
	RegisterHealthCallback(callback HealthChangeCallback) error

	// SetHealthCheckInterval sets the health check interval
	SetHealthCheckInterval(interval time.Duration) error

	// GetHealthHistory returns health history for a node
	GetHealthHistory(ctx context.Context, nodeID string, duration time.Duration) ([]types.HealthRecord, error)

	// GetMetrics returns health monitor metrics
	GetMetrics() *HealthMonitorMetrics

	// HealthCheck performs health monitor self-check
	HealthCheck(ctx context.Context) error
}

// ServiceRegistry defines the interface for service registration and discovery
type ServiceRegistry interface {
	// RegisterService registers a service
	RegisterService(ctx context.Context, service *ServiceRegistration) error

	// UnregisterService unregisters a service
	UnregisterService(ctx context.Context, serviceID string) error

	// DiscoverServices discovers services by name or filter
	DiscoverServices(ctx context.Context, serviceName string, filter *ServiceFilter) ([]Service, error)

	// WatchServices watches for service changes
	WatchServices(ctx context.Context, serviceName string, callback ServiceChangeCallback) error

	// UpdateServiceHealth updates service health status
	UpdateServiceHealth(ctx context.Context, serviceID string, health types.HealthStatus) error

	// GetServiceEndpoints returns endpoints for a service
	GetServiceEndpoints(ctx context.Context, serviceName string) ([]ServiceEndpoint, error)

	// GetMetrics returns service registry metrics
	GetMetrics() *ServiceRegistryMetrics

	// HealthCheck performs service registry health check
	HealthCheck(ctx context.Context) error
}

// FailoverManager defines the interface for automatic failover management
type FailoverManager interface {
	// Start starts the failover manager
	Start(ctx context.Context) error

	// Stop stops the failover manager
	Stop(ctx context.Context) error

	// RegisterFailoverTarget registers a service for failover monitoring
	RegisterFailoverTarget(target *FailoverTarget) error

	// UnregisterFailoverTarget unregisters a failover target
	UnregisterFailoverTarget(targetID string) error

	// TriggerFailover manually triggers failover for a service
	TriggerFailover(ctx context.Context, targetID string, reason string) error

	// GetFailoverStatus returns failover status for all targets
	GetFailoverStatus(ctx context.Context) (map[string]*FailoverStatus, error)

	// RegisterFailoverCallback registers callback for failover events
	RegisterFailoverCallback(callback FailoverCallback) error

	// GetMetrics returns failover manager metrics
	GetMetrics() *FailoverMetrics

	// HealthCheck performs failover manager health check
	HealthCheck(ctx context.Context) error
}

// Data structures

// ClusterConfig defines cluster configuration
type ClusterConfig struct {
	ClusterName      string `json:"clusterName"`
	NodeID           string `json:"nodeId"`
	NodeName         string `json:"nodeName"`
	BindAddress      string `json:"bindAddress"`
	AdvertiseAddress string `json:"advertiseAddress"`
	Port             int    `json:"port"`

	// Discovery configuration
	DiscoveryMethod  DiscoveryMethod   `json:"discoveryMethod"`
	SeedNodes        []string          `json:"seedNodes,omitempty"`
	ConsulConfig     *ConsulConfig     `json:"consulConfig,omitempty"`
	EtcdConfig       *EtcdConfig       `json:"etcdConfig,omitempty"`
	KubernetesConfig *KubernetesConfig `json:"kubernetesConfig,omitempty"`

	// Election configuration
	EnableElection   bool          `json:"enableElection"`
	ElectionTimeout  time.Duration `json:"electionTimeout"`
	HeartbeatTimeout time.Duration `json:"heartbeatTimeout"`

	// Health check configuration
	HealthCheckInterval time.Duration `json:"healthCheckInterval"`
	HealthCheckTimeout  time.Duration `json:"healthCheckTimeout"`
	HealthCheckRetries  int           `json:"healthCheckRetries"`
	UnhealthyThreshold  int           `json:"unhealthyThreshold"`

	// Networking configuration
	EnableTLS   bool   `json:"enableTls"`
	TLSCertFile string `json:"tlsCertFile,omitempty"`
	TLSKeyFile  string `json:"tlsKeyFile,omitempty"`
	TLSCAFile   string `json:"tlsCaFile,omitempty"`

	// Security configuration
	AuthToken string `json:"authToken,omitempty"`

	// Performance configuration
	MaxNodes       int           `json:"maxNodes"`
	MessageTimeout time.Duration `json:"messageTimeout"`
	RetryBackoff   time.Duration `json:"retryBackoff"`

	Properties map[string]interface{} `json:"properties,omitempty"`
}

// DiscoveryMethod defines methods for node discovery
type DiscoveryMethod string

const (
	DiscoveryMethodStatic     DiscoveryMethod = "static"
	DiscoveryMethodConsul     DiscoveryMethod = "consul"
	DiscoveryMethodEtcd       DiscoveryMethod = "etcd"
	DiscoveryMethodKubernetes DiscoveryMethod = "kubernetes"
	DiscoveryMethodDNS        DiscoveryMethod = "dns"
	DiscoveryMethodMulticast  DiscoveryMethod = "multicast"
)

// ConsulConfig defines Consul configuration for discovery
type ConsulConfig struct {
	Address     string            `json:"address"`
	Token       string            `json:"token,omitempty"`
	Datacenter  string            `json:"datacenter,omitempty"`
	ServiceName string            `json:"serviceName"`
	Tags        []string          `json:"tags,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
	CheckTTL    time.Duration     `json:"checkTtl,omitempty"`
}

// EtcdConfig defines etcd configuration for discovery
type EtcdConfig struct {
	Endpoints   []string      `json:"endpoints"`
	Username    string        `json:"username,omitempty"`
	Password    string        `json:"password,omitempty"`
	KeyPrefix   string        `json:"keyPrefix"`
	LeaseTTL    time.Duration `json:"leaseTtl"`
	DialTimeout time.Duration `json:"dialTimeout,omitempty"`
}

// KubernetesConfig defines Kubernetes configuration for discovery
type KubernetesConfig struct {
	Namespace     string            `json:"namespace"`
	ServiceName   string            `json:"serviceName"`
	PodName       string            `json:"podName"`
	LabelSelector string            `json:"labelSelector,omitempty"`
	FieldSelector string            `json:"fieldSelector,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
}

// Node represents a node in the cluster
type Node struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Address      string                 `json:"address"`
	Port         int                    `json:"port"`
	Type         NodeType               `json:"type"`
	Status       NodeStatus             `json:"status"`
	Health       types.HealthStatus     `json:"health"`
	Version      string                 `json:"version"`
	Tags         []string               `json:"tags,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	LastSeen     time.Time              `json:"lastSeen"`
	JoinedAt     time.Time              `json:"joinedAt"`
	Capabilities []string               `json:"capabilities,omitempty"`
	Load         *NodeLoad              `json:"load,omitempty"`
}

// NodeType defines types of nodes
type NodeType string

const (
	NodeTypePortal  NodeType = "portal"
	NodeTypeChronos NodeType = "chronos"
	NodeTypeVulcan  NodeType = "vulcan"
	NodeTypeArchive NodeType = "archive"
	NodeTypeGeneric NodeType = "generic"
)

// NodeStatus defines node status
type NodeStatus string

const (
	NodeStatusActive    NodeStatus = "active"
	NodeStatusInactive  NodeStatus = "inactive"
	NodeStatusJoining   NodeStatus = "joining"
	NodeStatusLeaving   NodeStatus = "leaving"
	NodeStatusSuspended NodeStatus = "suspended"
	NodeStatusFailed    NodeStatus = "failed"
)

// HealthStatus defines health status
// HealthStatus and constants are defined in types/common.go

// NodeLoad represents current load on a node
type NodeLoad struct {
	CPUPercent     float64 `json:"cpuPercent"`
	MemoryPercent  float64 `json:"memoryPercent"`
	DiskPercent    float64 `json:"diskPercent"`
	NetworkInMbps  float64 `json:"networkInMbps"`
	NetworkOutMbps float64 `json:"networkOutMbps"`
	Connections    int     `json:"connections"`
	RequestsPerSec float64 `json:"requestsPerSec"`
}

// NodeRegistration represents node registration information
type NodeRegistration struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Type         NodeType               `json:"type"`
	Address      string                 `json:"address"`
	Port         int                    `json:"port"`
	Version      string                 `json:"version"`
	Tags         []string               `json:"tags,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Capabilities []string               `json:"capabilities,omitempty"`
	HealthCheck  *HealthCheckConfig     `json:"healthCheck,omitempty"`
}

// HealthCheckConfig defines health check configuration for a node
type HealthCheckConfig struct {
	Type     HealthCheckType   `json:"type"`
	Endpoint string            `json:"endpoint,omitempty"`
	Interval time.Duration     `json:"interval"`
	Timeout  time.Duration     `json:"timeout"`
	Method   string            `json:"method,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Body     string            `json:"body,omitempty"`
}

// HealthCheckType defines types of health checks
type HealthCheckType string

const (
	HealthCheckTypeHTTP HealthCheckType = "http"
	HealthCheckTypeTCP  HealthCheckType = "tcp"
	HealthCheckTypeGRPC HealthCheckType = "grpc"
	HealthCheckTypeExec HealthCheckType = "exec"
)

// NodeFilter defines filter criteria for node discovery
type NodeFilter struct {
	Type         NodeType           `json:"type,omitempty"`
	Status       NodeStatus         `json:"status,omitempty"`
	Health       types.HealthStatus `json:"health,omitempty"`
	Tags         []string           `json:"tags,omitempty"`
	Capabilities []string           `json:"capabilities,omitempty"`
	MinVersion   string             `json:"minVersion,omitempty"`
	MaxLoad      *NodeLoad          `json:"maxLoad,omitempty"`
}

// ClusterInfo represents overall cluster information
type ClusterInfo struct {
	Name          string                     `json:"name"`
	TotalNodes    int                        `json:"totalNodes"`
	ActiveNodes   int                        `json:"activeNodes"`
	LeaderNodeID  string                     `json:"leaderNodeId,omitempty"`
	Version       string                     `json:"version"`
	Status        ClusterStatus              `json:"status"`
	Health        types.HealthStatus         `json:"health"`
	CreatedAt     time.Time                  `json:"createdAt"`
	NodesByType   map[NodeType]int           `json:"nodesByType"`
	NodesByStatus map[NodeStatus]int         `json:"nodesByStatus"`
	NodesByHealth map[types.HealthStatus]int `json:"nodesByHealth"`
	Properties    map[string]interface{}     `json:"properties,omitempty"`
}

// ClusterStatus defines cluster status
type ClusterStatus string

const (
	ClusterStatusForming     ClusterStatus = "forming"
	ClusterStatusActive      ClusterStatus = "active"
	ClusterStatusDegraded    ClusterStatus = "degraded"
	ClusterStatusMaintenance ClusterStatus = "maintenance"
	ClusterStatusShutdown    ClusterStatus = "shutdown"
)

// LeaderInfo represents information about the cluster leader
type LeaderInfo struct {
	NodeID      string    `json:"nodeId"`
	NodeName    string    `json:"nodeName"`
	Address     string    `json:"address"`
	ElectedAt   time.Time `json:"electedAt"`
	Term        int64     `json:"term"`
	LastContact time.Time `json:"lastContact"`
}

// ElectionStatus represents the status of leader election
type ElectionStatus struct {
	IsLeader     bool      `json:"isLeader"`
	LeaderID     string    `json:"leaderId,omitempty"`
	Term         int64     `json:"term"`
	LastElection time.Time `json:"lastElection"`
	Participants []string  `json:"participants"`
	Status       string    `json:"status"`
}

// LoadBalancingStrategy defines load balancing strategies
type LoadBalancingStrategy string

const (
	LoadBalancingStrategyRoundRobin     LoadBalancingStrategy = "round_robin"
	LoadBalancingStrategyWeightedRandom LoadBalancingStrategy = "weighted_random"
	LoadBalancingStrategyLeastLoaded    LoadBalancingStrategy = "least_loaded"
	LoadBalancingStrategyConsistentHash LoadBalancingStrategy = "consistent_hash"
	LoadBalancingStrategyPowerOfTwo     LoadBalancingStrategy = "power_of_two"
)

// SelectionCriteria defines criteria for node selection
type SelectionCriteria struct {
	Type         NodeType               `json:"type,omitempty"`
	Capabilities []string               `json:"capabilities,omitempty"`
	Tags         []string               `json:"tags,omitempty"`
	MaxLoad      *NodeLoad              `json:"maxLoad,omitempty"`
	HashKey      string                 `json:"hashKey,omitempty"`
	Properties   map[string]interface{} `json:"properties,omitempty"`
}

// ClusterHealthStatus represents overall cluster health status
type ClusterHealthStatus struct {
	Status         types.HealthStatus            `json:"status"`
	TotalNodes     int                           `json:"totalNodes"`
	HealthyNodes   int                           `json:"healthyNodes"`
	UnhealthyNodes int                           `json:"unhealthyNodes"`
	DegradedNodes  int                           `json:"degradedNodes"`
	NodeHealth     map[string]types.HealthStatus `json:"nodeHealth"`
	Timestamp      time.Time                     `json:"timestamp"`
	Details        map[string]interface{}        `json:"details,omitempty"`
}

// HealthRecord represents a health check record
// HealthRecord is defined in types/common.go

// Service represents a service in the registry
type Service struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Address      string                 `json:"address"`
	Port         int                    `json:"port"`
	Protocol     string                 `json:"protocol"`
	Health       types.HealthStatus     `json:"health"`
	Tags         []string               `json:"tags,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Endpoints    []ServiceEndpoint      `json:"endpoints,omitempty"`
	RegisteredAt time.Time              `json:"registeredAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

// ServiceRegistration represents service registration information
type ServiceRegistration struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	Address     string                 `json:"address"`
	Port        int                    `json:"port"`
	Protocol    string                 `json:"protocol"`
	Tags        []string               `json:"tags,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Endpoints   []ServiceEndpoint      `json:"endpoints,omitempty"`
	HealthCheck *HealthCheckConfig     `json:"healthCheck,omitempty"`
}

// ServiceEndpoint represents a service endpoint
type ServiceEndpoint struct {
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	Method   string            `json:"method,omitempty"`
	Protocol string            `json:"protocol"`
	Port     int               `json:"port,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// ServiceFilter defines filter criteria for service discovery
type ServiceFilter struct {
	Version  string             `json:"version,omitempty"`
	Protocol string             `json:"protocol,omitempty"`
	Health   types.HealthStatus `json:"health,omitempty"`
	Tags     []string           `json:"tags,omitempty"`
	Metadata map[string]string  `json:"metadata,omitempty"`
}

// FailoverTarget represents a service target for failover monitoring
type FailoverTarget struct {
	ID                string                 `json:"id"`
	ServiceName       string                 `json:"serviceName"`
	PrimaryEndpoint   string                 `json:"primaryEndpoint"`
	BackupEndpoints   []string               `json:"backupEndpoints"`
	HealthCheckConfig *HealthCheckConfig     `json:"healthCheckConfig"`
	FailoverPolicy    *FailoverPolicy        `json:"failoverPolicy"`
	Properties        map[string]interface{} `json:"properties,omitempty"`
}

// FailoverPolicy defines failover policy
type FailoverPolicy struct {
	AutoFailover        bool          `json:"autoFailover"`
	FailureThreshold    int           `json:"failureThreshold"`
	RecoveryThreshold   int           `json:"recoveryThreshold"`
	CheckInterval       time.Duration `json:"checkInterval"`
	FailbackDelay       time.Duration `json:"failbackDelay"`
	MaxFailoverAttempts int           `json:"maxFailoverAttempts"`
}

// FailoverStatus represents failover status
type FailoverStatus struct {
	TargetID         string        `json:"targetId"`
	CurrentEndpoint  string        `json:"currentEndpoint"`
	Status           FailoverState `json:"status"`
	FailureCount     int           `json:"failureCount"`
	LastFailoverTime *time.Time    `json:"lastFailoverTime,omitempty"`
	LastRecoveryTime *time.Time    `json:"lastRecoveryTime,omitempty"`
	FailoverAttempts int           `json:"failoverAttempts"`
}

// FailoverState defines failover states
type FailoverState string

const (
	FailoverStateNormal     FailoverState = "normal"
	FailoverStateFailed     FailoverState = "failed"
	FailoverStateFailedOver FailoverState = "failed_over"
	FailoverStateRecovering FailoverState = "recovering"
)

// Callback function types
type NodeChangeCallback func(event NodeChangeEvent)
type LeadershipCallback func(isLeader bool, leaderInfo *LeaderInfo)
type HealthChangeCallback func(nodeID string, oldStatus, newStatus types.HealthStatus)
type ServiceChangeCallback func(event ServiceChangeEvent)
type FailoverCallback func(event FailoverEvent)

// Event structures

// NodeChangeEvent represents a node change event
type NodeChangeEvent struct {
	Type      NodeChangeEventType `json:"type"`
	Node      Node                `json:"node"`
	OldStatus NodeStatus          `json:"oldStatus,omitempty"`
	NewStatus NodeStatus          `json:"newStatus,omitempty"`
	Timestamp time.Time           `json:"timestamp"`
}

// NodeChangeEventType defines types of node change events
type NodeChangeEventType string

const (
	NodeChangeEventTypeJoined  NodeChangeEventType = "joined"
	NodeChangeEventTypeLeft    NodeChangeEventType = "left"
	NodeChangeEventTypeUpdated NodeChangeEventType = "updated"
	NodeChangeEventTypeFailed  NodeChangeEventType = "failed"
)

// ServiceChangeEvent represents a service change event
type ServiceChangeEvent struct {
	Type      ServiceChangeEventType `json:"type"`
	Service   Service                `json:"service"`
	Timestamp time.Time              `json:"timestamp"`
}

// ServiceChangeEventType defines types of service change events
type ServiceChangeEventType string

const (
	ServiceChangeEventTypeRegistered    ServiceChangeEventType = "registered"
	ServiceChangeEventTypeUnregistered  ServiceChangeEventType = "unregistered"
	ServiceChangeEventTypeHealthChanged ServiceChangeEventType = "health_changed"
)

// FailoverEvent represents a failover event
type FailoverEvent struct {
	Type        FailoverEventType `json:"type"`
	TargetID    string            `json:"targetId"`
	ServiceName string            `json:"serviceName"`
	OldEndpoint string            `json:"oldEndpoint,omitempty"`
	NewEndpoint string            `json:"newEndpoint,omitempty"`
	Reason      string            `json:"reason"`
	Timestamp   time.Time         `json:"timestamp"`
}

// FailoverEventType defines types of failover events
type FailoverEventType string

const (
	FailoverEventTypeTriggered FailoverEventType = "triggered"
	FailoverEventTypeRecovered FailoverEventType = "recovered"
	FailoverEventTypeFailed    FailoverEventType = "failed"
)

// Metrics structures

// ClusterMetrics defines cluster-wide metrics
type ClusterMetrics struct {
	NodeCount         int                        `json:"nodeCount"`
	ActiveNodeCount   int                        `json:"activeNodeCount"`
	NodesByType       map[NodeType]int           `json:"nodesByType"`
	NodesByStatus     map[NodeStatus]int         `json:"nodesByStatus"`
	NodesByHealth     map[types.HealthStatus]int `json:"nodesByHealth"`
	AverageLoad       *NodeLoad                  `json:"averageLoad"`
	NetworkPartitions int                        `json:"networkPartitions"`
	MessagesSent      uint64                     `json:"messagesSent"`
	MessagesReceived  uint64                     `json:"messagesReceived"`
	MessageLatencyP99 time.Duration              `json:"messageLatencyP99"`
	ClusterEvents     uint64                     `json:"clusterEvents"`
	Timestamp         time.Time                  `json:"timestamp"`
}

// ElectionMetrics defines leader election metrics
type ElectionMetrics struct {
	ElectionCount      uint64        `json:"electionCount"`
	LeadershipTime     time.Duration `json:"leadershipTime"`
	LastElectionTime   time.Time     `json:"lastElectionTime"`
	ElectionDuration   time.Duration `json:"electionDuration"`
	SplitBrainEvents   uint64        `json:"splitBrainEvents"`
	IsCurrentLeader    bool          `json:"isCurrentLeader"`
	CandidateCount     int           `json:"candidateCount"`
	HeartbeatsSent     uint64        `json:"heartbeatsSent"`
	HeartbeatsReceived uint64        `json:"heartbeatsReceived"`
	Timestamp          time.Time     `json:"timestamp"`
}

// DiscoveryMetrics defines node discovery metrics
type DiscoveryMetrics struct {
	NodesDiscovered     uint64        `json:"nodesDiscovered"`
	NodesRegistered     uint64        `json:"nodesRegistered"`
	NodesUnregistered   uint64        `json:"nodesUnregistered"`
	DiscoveryErrors     uint64        `json:"discoveryErrors"`
	DiscoveryLatencyP99 time.Duration `json:"discoveryLatencyP99"`
	WatchConnections    int           `json:"watchConnections"`
	Timestamp           time.Time     `json:"timestamp"`
}

// LoadBalancerMetrics defines load balancer metrics
type LoadBalancerMetrics struct {
	RequestsBalanced    uint64                `json:"requestsBalanced"`
	RequestsByNode      map[string]uint64     `json:"requestsByNode"`
	SelectionLatencyP99 time.Duration         `json:"selectionLatencyP99"`
	NodesAvailable      int                   `json:"nodesAvailable"`
	Strategy            LoadBalancingStrategy `json:"strategy"`
	Rebalances          uint64                `json:"rebalances"`
	Timestamp           time.Time             `json:"timestamp"`
}

// HealthMonitorMetrics defines health monitor metrics
type HealthMonitorMetrics struct {
	HealthChecksPerformed uint64            `json:"healthChecksPerformed"`
	HealthChecksByNode    map[string]uint64 `json:"healthChecksByNode"`
	HealthCheckLatencyP99 time.Duration     `json:"healthCheckLatencyP99"`
	HealthTransitions     uint64            `json:"healthTransitions"`
	HealthCheckErrors     uint64            `json:"healthCheckErrors"`
	Timestamp             time.Time         `json:"timestamp"`
}

// ServiceRegistryMetrics defines service registry metrics
type ServiceRegistryMetrics struct {
	ServicesRegistered      uint64        `json:"servicesRegistered"`
	ServicesActive          uint64        `json:"servicesActive"`
	ServiceLookups          uint64        `json:"serviceLookups"`
	ServiceLookupLatencyP99 time.Duration `json:"serviceLookupLatencyP99"`
	RegistryErrors          uint64        `json:"registryErrors"`
	WatchConnections        int           `json:"watchConnections"`
	Timestamp               time.Time     `json:"timestamp"`
}

// FailoverMetrics defines failover metrics
type FailoverMetrics struct {
	TargetsMonitored     int           `json:"targetsMonitored"`
	FailoversTriggered   uint64        `json:"failoversTriggered"`
	FailoverAttempts     uint64        `json:"failoverAttempts"`
	SuccessfulFailovers  uint64        `json:"successfulFailovers"`
	FailedFailovers      uint64        `json:"failedFailovers"`
	RecoveryAttempts     uint64        `json:"recoveryAttempts"`
	SuccessfulRecoveries uint64        `json:"successfulRecoveries"`
	AverageFailoverTime  time.Duration `json:"averageFailoverTime"`
	Timestamp            time.Time     `json:"timestamp"`
}
