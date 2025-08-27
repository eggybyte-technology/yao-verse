package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/eggybyte-technology/yao-portal/interfaces"
)

// BackendServiceManagerImpl implements the BackendServiceManager interface
type BackendServiceManagerImpl struct {
	mu     sync.RWMutex
	config *interfaces.PortalConfig

	// Service endpoints mapping
	services map[string]*ServiceInfo

	// HTTP client for health checks
	httpClient *http.Client

	// Metrics
	metrics     *interfaces.BackendServiceMetrics
	stopChannel chan struct{}
}

// ServiceInfo holds information about a backend service
type ServiceInfo struct {
	Name         string
	Namespace    string
	Endpoints    []string
	HealthPath   string
	HealthStatus map[string]bool
	LastCheck    map[string]time.Time
	RequestCount map[string]uint64
	ResponseTime map[string]time.Duration
	ErrorCount   map[string]uint64
}

// NewBackendServiceManager creates a new backend service manager
func NewBackendServiceManager(config *interfaces.PortalConfig) (interfaces.BackendServiceManager, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if config.BackendServices == nil {
		return nil, fmt.Errorf("backend services configuration is required")
	}

	manager := &BackendServiceManagerImpl{
		config:   config,
		services: make(map[string]*ServiceInfo),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		metrics: &interfaces.BackendServiceMetrics{
			ServiceHealth:        make(map[string]bool),
			ServiceRequestCounts: make(map[string]uint64),
			ServiceResponseTimes: make(map[string]time.Duration),
			ServiceErrorRates:    make(map[string]float64),
			LastHealthCheckTime:  make(map[string]time.Time),
		},
		stopChannel: make(chan struct{}),
	}

	// Initialize services from configuration
	if err := manager.initializeServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	// Start health check routine
	go manager.healthCheckLoop()

	return manager, nil
}

// initializeServices initializes service configurations
func (b *BackendServiceManagerImpl) initializeServices() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	serviceConfigs := map[string]*interfaces.BackendServiceConfig{
		"consensus": b.config.BackendServices.ConsensusService,
		"execution": b.config.BackendServices.ExecutionService,
		"state":     b.config.BackendServices.StateService,
		"archive":   b.config.BackendServices.ArchiveService,
	}

	for serviceType, config := range serviceConfigs {
		if config == nil {
			continue
		}

		serviceInfo := &ServiceInfo{
			Name:         config.Name,
			Namespace:    config.Namespace,
			Endpoints:    config.Endpoints,
			HealthPath:   config.HealthCheckPath,
			HealthStatus: make(map[string]bool),
			LastCheck:    make(map[string]time.Time),
			RequestCount: make(map[string]uint64),
			ResponseTime: make(map[string]time.Duration),
			ErrorCount:   make(map[string]uint64),
		}

		// Initialize health status for all endpoints
		for _, endpoint := range config.Endpoints {
			serviceInfo.HealthStatus[endpoint] = false
		}

		b.services[serviceType] = serviceInfo
	}

	return nil
}

// GetService returns the first healthy endpoint for a service type
func (b *BackendServiceManagerImpl) GetService(ctx context.Context, serviceType string) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	serviceInfo, exists := b.services[serviceType]
	if !exists {
		return "", fmt.Errorf("service type %s not found", serviceType)
	}

	// Return the first endpoint (for simplicity, since we don't implement load balancing)
	if len(serviceInfo.Endpoints) > 0 {
		return serviceInfo.Endpoints[0], nil
	}

	return "", fmt.Errorf("no endpoints configured for service type %s", serviceType)
}

// GetHealthyService returns a healthy endpoint for a service type
func (b *BackendServiceManagerImpl) GetHealthyService(ctx context.Context, serviceType string) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	serviceInfo, exists := b.services[serviceType]
	if !exists {
		return "", fmt.Errorf("service type %s not found", serviceType)
	}

	// Find the first healthy endpoint
	for _, endpoint := range serviceInfo.Endpoints {
		if isHealthy, ok := serviceInfo.HealthStatus[endpoint]; ok && isHealthy {
			return endpoint, nil
		}
	}

	// If no healthy endpoint found, return the first one as fallback
	if len(serviceInfo.Endpoints) > 0 {
		return serviceInfo.Endpoints[0], nil
	}

	return "", fmt.Errorf("no healthy endpoints available for service type %s", serviceType)
}

// CheckServiceHealth checks the health of a specific service endpoint
func (b *BackendServiceManagerImpl) CheckServiceHealth(ctx context.Context, serviceType string, endpoint string) (bool, error) {
	b.mu.RLock()
	serviceInfo, exists := b.services[serviceType]
	b.mu.RUnlock()

	if !exists {
		return false, fmt.Errorf("service type %s not found", serviceType)
	}

	healthURL := endpoint + serviceInfo.HealthPath
	start := time.Now()

	// Create request with context timeout
	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		b.recordHealthCheckResult(serviceType, endpoint, false, time.Since(start))
		return false, fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := b.httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		b.recordHealthCheckResult(serviceType, endpoint, false, duration)
		return false, err
	}
	defer resp.Body.Close()

	isHealthy := resp.StatusCode >= 200 && resp.StatusCode < 300
	b.recordHealthCheckResult(serviceType, endpoint, isHealthy, duration)

	return isHealthy, nil
}

// recordHealthCheckResult records the result of a health check
func (b *BackendServiceManagerImpl) recordHealthCheckResult(serviceType, endpoint string, isHealthy bool, responseTime time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()

	serviceInfo, exists := b.services[serviceType]
	if !exists {
		return
	}

	serviceInfo.HealthStatus[endpoint] = isHealthy
	serviceInfo.LastCheck[endpoint] = time.Now()
	serviceInfo.ResponseTime[endpoint] = responseTime

	if !isHealthy {
		serviceInfo.ErrorCount[endpoint]++
	}

	// Update global metrics
	b.metrics.ServiceHealth[serviceType+"."+endpoint] = isHealthy
	b.metrics.ServiceResponseTimes[serviceType+"."+endpoint] = responseTime
	b.metrics.LastHealthCheckTime[serviceType+"."+endpoint] = time.Now()
}

// GetAllServiceEndpoints returns all endpoints for a service type
func (b *BackendServiceManagerImpl) GetAllServiceEndpoints(ctx context.Context, serviceType string) ([]string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	serviceInfo, exists := b.services[serviceType]
	if !exists {
		return nil, fmt.Errorf("service type %s not found", serviceType)
	}

	return serviceInfo.Endpoints, nil
}

// UpdateServiceHealth manually updates the health status of a service endpoint
func (b *BackendServiceManagerImpl) UpdateServiceHealth(ctx context.Context, serviceType string, endpoint string, isHealthy bool) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	serviceInfo, exists := b.services[serviceType]
	if !exists {
		return fmt.Errorf("service type %s not found", serviceType)
	}

	serviceInfo.HealthStatus[endpoint] = isHealthy
	serviceInfo.LastCheck[endpoint] = time.Now()

	// Update global metrics
	b.metrics.ServiceHealth[serviceType+"."+endpoint] = isHealthy
	b.metrics.LastHealthCheckTime[serviceType+"."+endpoint] = time.Now()

	return nil
}

// GetServiceMetrics returns metrics for backend services
func (b *BackendServiceManagerImpl) GetServiceMetrics(ctx context.Context) (*interfaces.BackendServiceMetrics, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Update error rates
	for serviceType, serviceInfo := range b.services {
		for _, endpoint := range serviceInfo.Endpoints {
			key := serviceType + "." + endpoint
			totalRequests := serviceInfo.RequestCount[endpoint]
			errorCount := serviceInfo.ErrorCount[endpoint]

			if totalRequests > 0 {
				b.metrics.ServiceErrorRates[key] = float64(errorCount) / float64(totalRequests) * 100
			} else {
				b.metrics.ServiceErrorRates[key] = 0
			}

			b.metrics.ServiceRequestCounts[key] = totalRequests
		}
	}

	b.metrics.Timestamp = time.Now()

	// Return a copy of the metrics
	metricsCopy := &interfaces.BackendServiceMetrics{
		ServiceHealth:        make(map[string]bool),
		ServiceRequestCounts: make(map[string]uint64),
		ServiceResponseTimes: make(map[string]time.Duration),
		ServiceErrorRates:    make(map[string]float64),
		LastHealthCheckTime:  make(map[string]time.Time),
		Timestamp:            b.metrics.Timestamp,
	}

	for k, v := range b.metrics.ServiceHealth {
		metricsCopy.ServiceHealth[k] = v
	}
	for k, v := range b.metrics.ServiceRequestCounts {
		metricsCopy.ServiceRequestCounts[k] = v
	}
	for k, v := range b.metrics.ServiceResponseTimes {
		metricsCopy.ServiceResponseTimes[k] = v
	}
	for k, v := range b.metrics.ServiceErrorRates {
		metricsCopy.ServiceErrorRates[k] = v
	}
	for k, v := range b.metrics.LastHealthCheckTime {
		metricsCopy.LastHealthCheckTime[k] = v
	}

	return metricsCopy, nil
}

// healthCheckLoop runs periodic health checks for all services
func (b *BackendServiceManagerImpl) healthCheckLoop() {
	ticker := time.NewTicker(30 * time.Second) // Health check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			b.performHealthChecks()
		case <-b.stopChannel:
			return
		}
	}
}

// performHealthChecks performs health checks on all registered services
func (b *BackendServiceManagerImpl) performHealthChecks() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	b.mu.RLock()
	services := make(map[string]*ServiceInfo)
	for k, v := range b.services {
		services[k] = v
	}
	b.mu.RUnlock()

	for serviceType, serviceInfo := range services {
		for _, endpoint := range serviceInfo.Endpoints {
			go func(sType, ep string) {
				_, err := b.CheckServiceHealth(ctx, sType, ep)
				if err != nil {
					// Log error - in production, use proper logging
					fmt.Printf("Health check failed for %s at %s: %v\n", sType, ep, err)
				}
			}(serviceType, endpoint)
		}
	}
}

// RecordRequest records a request made to a service endpoint
func (b *BackendServiceManagerImpl) RecordRequest(serviceType, endpoint string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	serviceInfo, exists := b.services[serviceType]
	if !exists {
		return
	}

	serviceInfo.RequestCount[endpoint]++
	b.metrics.ServiceRequestCounts[serviceType+"."+endpoint] = serviceInfo.RequestCount[endpoint]
}

// Stop stops the backend service manager
func (b *BackendServiceManagerImpl) Stop() {
	close(b.stopChannel)
}
