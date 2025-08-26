package config

import (
	"os"

	"github.com/eggybyte-technology/yao-portal/interfaces"
)

// LoadConfig loads configuration from file
func LoadConfig(configPath string) (*interfaces.PortalConfig, error) {
	// Check if environment variable is set
	if envPath := os.Getenv("YAO_PORTAL_CONFIG"); envPath != "" {
		configPath = envPath
	}

	// For now, return a default configuration
	// TODO: Implement actual configuration loading from YAML/JSON file
	return getDefaultConfig(), nil
}

// getDefaultConfig returns the default configuration
func getDefaultConfig() *interfaces.PortalConfig {
	return &interfaces.PortalConfig{
		// Server configuration
		Host: "0.0.0.0",
		Port: 8545,

		// TLS configuration
		EnableTLS: false,
		CertFile:  "",
		KeyFile:   "",

		// CORS configuration
		EnableCORS:     true,
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},

		// Rate limiting
		EnableRateLimit: true,
		RateLimit:       1000, // 1000 requests per second
		BurstLimit:      2000,

		// Transaction pool configuration
		TxPoolSize:    10000,
		TxPoolTimeout: 300,         // 5 minutes
		MaxTxDataSize: 1024 * 1024, // 1MB

		// Message queue configuration
		RocketMQEndpoints: []string{"localhost:9876"},
		RocketMQAccessKey: "",
		RocketMQSecretKey: "",
		ProducerGroup:     "yao-portal-producer",

		// Chain configuration
		ChainID:     1337,
		NetworkID:   1337,
		GenesisHash: "0x0000000000000000000000000000000000000000000000000000000000000000",

		// Performance tuning
		MaxConcurrentRequests: 1000,
		RequestTimeout:        30, // 30 seconds
		EnableMetrics:         true,
		MetricsInterval:       10, // 10 seconds

		// Logging configuration
		LogLevel:  "info",
		LogFormat: "json",
		LogFile:   "logs/portal.log",
	}
}
