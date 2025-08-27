package config

import (
	"fmt"
	"os"
	"time"

	"github.com/eggybyte-technology/yao-portal/interfaces"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from file
func LoadConfig(configPath string) (*interfaces.PortalConfig, error) {
	// Check if environment variable is set
	if envPath := os.Getenv("YAO_PORTAL_CONFIG"); envPath != "" {
		configPath = envPath
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If config file doesn't exist, return default config
		return getDefaultConfig(), nil
	}

	// Read configuration file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse YAML configuration
	var config interfaces.PortalConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides
	applyEnvOverrides(&config)

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// getDefaultConfig returns the default configuration
func getDefaultConfig() *interfaces.PortalConfig {
	return &interfaces.PortalConfig{
		// Node configuration
		NodeID:   "portal-node-1",
		NodeName: "YaoPortal Node 1",

		// Server configuration
		HTTPBindAddress: "0.0.0.0",
		HTTPPort:        8545,
		WSBindAddress:   "0.0.0.0",
		WSPort:          8546,
		EnableTLS:       false,
		TLSCert:         "",
		TLSKey:          "",

		// JSON-RPC configuration
		EnableHTTP:     true,
		EnableWS:       true,
		CORSOrigins:    []string{"*"},
		MaxRequestSize: 1024 * 1024, // 1MB
		RequestTimeout: 30 * time.Second,
		BatchLimit:     100,

		// Transaction pool configuration
		TxPoolCapacity: 10000,
		TxPoolTimeout:  300 * time.Second,
		MaxTxLifetime:  600 * time.Second,

		// WebSocket configuration
		WSMaxConnections: 1000,
		WSPingInterval:   30 * time.Second,
		WSWriteTimeout:   10 * time.Second,
		WSReadTimeout:    10 * time.Second,

		// State reading configuration
		StateReadTimeout: 5 * time.Second,
		MaxBlockRange:    1000,
		EnableTracing:    false,

		// Message queue configuration
		RocketMQEndpoints: []string{"localhost:9876"},
		RocketMQAccessKey: "",
		RocketMQSecretKey: "",
		ProducerGroup:     "yao-portal-producer",

		// Rate limiting configuration
		EnableRateLimit:   true,
		RequestsPerSecond: 1000,
		BurstSize:         2000,

		// Authentication configuration
		EnableAuth:     false,
		AllowedOrigins: []string{"*"},
		APIKeys:        []string{},

		// Monitoring configuration
		EnableMetrics:       true,
		MetricsInterval:     10 * time.Second,
		EnableHealthCheck:   true,
		HealthCheckInterval: 30 * time.Second,

		// Logging configuration
		LogLevel:     "info",
		LogFormat:    "json",
		LogFile:      "logs/portal.log",
		LogRequests:  true,
		LogResponses: false,

		// Backend services configuration
		BackendServices: &interfaces.BackendServicesConfig{
			ConsensusService: &interfaces.BackendServiceConfig{
				Name:            "yao-vulcan-service",
				Namespace:       "yao-verse",
				Endpoints:       []string{"http://yao-vulcan-service.yao-verse.svc.cluster.local:8080"},
				HealthCheckPath: "/health",
			},
			ExecutionService: &interfaces.BackendServiceConfig{
				Name:            "yao-chronos-service",
				Namespace:       "yao-verse",
				Endpoints:       []string{"http://yao-chronos-service.yao-verse.svc.cluster.local:8080"},
				HealthCheckPath: "/health",
			},
			StateService: &interfaces.BackendServiceConfig{
				Name:            "yao-oracle-service",
				Namespace:       "yao-verse",
				Endpoints:       []string{"http://yao-oracle-service.yao-verse.svc.cluster.local:8080"},
				HealthCheckPath: "/health",
			},
			ArchiveService: &interfaces.BackendServiceConfig{
				Name:            "yao-archive-service",
				Namespace:       "yao-verse",
				Endpoints:       []string{"http://yao-archive-service.yao-verse.svc.cluster.local:8080"},
				HealthCheckPath: "/health",
			},
		},

		// Ethereum configuration
		Ethereum: &interfaces.EthereumConfig{
			NetworkName:     "YaoVerse",
			ChainId:         1337,
			DefaultGasPrice: "20000000000", // 20 gwei
			DefaultGasLimit: 21000,
			BlockTime:       3 * time.Second,
			MaxBlockSize:    15000000, // 15M gas limit
		},

		// Additional configurations
		ConsumerGroup:       "yao-portal-consumer",
		RateLimitWindowSize: 1 * time.Minute,
		MetricsPort:         9090,
		EnableAccessLog:     true,
	}
}

// applyEnvOverrides applies environment variable overrides
func applyEnvOverrides(config *interfaces.PortalConfig) {
	if logLevel := os.Getenv("YAO_PORTAL_LOG_LEVEL"); logLevel != "" {
		config.LogLevel = logLevel
	}
	if logFormat := os.Getenv("YAO_PORTAL_LOG_FORMAT"); logFormat != "" {
		config.LogFormat = logFormat
	}
	if nodeID := os.Getenv("YAO_PORTAL_NODE_ID"); nodeID != "" {
		config.NodeID = nodeID
	}
	if nodeName := os.Getenv("YAO_PORTAL_NODE_NAME"); nodeName != "" {
		config.NodeName = nodeName
	}
}

// validateConfig validates the configuration
func validateConfig(config *interfaces.PortalConfig) error {
	if config.NodeID == "" {
		return fmt.Errorf("nodeId cannot be empty")
	}
	if config.HTTPPort <= 0 || config.HTTPPort > 65535 {
		return fmt.Errorf("invalid HTTP port: %d", config.HTTPPort)
	}
	if config.WSPort <= 0 || config.WSPort > 65535 {
		return fmt.Errorf("invalid WebSocket port: %d", config.WSPort)
	}
	if config.TxPoolCapacity <= 0 {
		return fmt.Errorf("txPoolCapacity must be greater than 0")
	}
	if config.RequestTimeout <= 0 {
		return fmt.Errorf("requestTimeout must be greater than 0")
	}
	if len(config.RocketMQEndpoints) == 0 {
		return fmt.Errorf("rocketmqEndpoints cannot be empty")
	}
	if config.ProducerGroup == "" {
		return fmt.Errorf("producerGroup cannot be empty")
	}

	return nil
}
