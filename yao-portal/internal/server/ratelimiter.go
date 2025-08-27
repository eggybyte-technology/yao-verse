package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/eggybyte-technology/yao-portal/interfaces"
	"golang.org/x/time/rate"
)

// RateLimiterImpl implements the rate limiter interface
type RateLimiterImpl struct {
	mu     sync.RWMutex
	config *interfaces.PortalConfig

	// Global rate limiter
	globalLimiter *rate.Limiter

	// Per-client rate limiters
	clientLimiters map[string]*clientLimiter

	// Configuration
	defaultLimit    *interfaces.RateLimitConfig
	cleanupInterval time.Duration

	// Metrics
	metrics   *interfaces.RateLimitStats
	startTime time.Time

	// Background cleanup
	stopChan chan struct{}
}

// clientLimiter represents a rate limiter for a specific client
type clientLimiter struct {
	limiter    *rate.Limiter
	config     *interfaces.RateLimitConfig
	lastAccess time.Time

	// Metrics for this client
	requestCount  int
	rejectedCount int
	windowStart   time.Time
	isBlocked     bool
	blockedUntil  *time.Time
}

// NewRateLimiter creates a new rate limiter instance
func NewRateLimiter(config *interfaces.PortalConfig) (*RateLimiterImpl, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	defaultConfig := &interfaces.RateLimitConfig{
		RequestsPerSecond: config.RequestsPerSecond,
		BurstSize:         config.BurstSize,
		WindowDuration:    time.Minute,
		BlockDuration:     5 * time.Minute,
	}

	rl := &RateLimiterImpl{
		config:          config,
		globalLimiter:   rate.NewLimiter(rate.Limit(config.RequestsPerSecond), config.BurstSize),
		clientLimiters:  make(map[string]*clientLimiter),
		defaultLimit:    defaultConfig,
		cleanupInterval: 10 * time.Minute,
		startTime:       time.Now(),
		stopChan:        make(chan struct{}),
		metrics: &interfaces.RateLimitStats{
			RequestsByEndpoint: make(map[string]uint64),
		},
	}

	// Start cleanup goroutine
	go rl.cleanupLoop()

	return rl, nil
}

// AllowRequest checks if a request should be allowed based on rate limits
func (r *RateLimiterImpl) AllowRequest(ctx context.Context, clientIP string, endpoint string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Update global metrics
	r.metrics.TotalRequests++
	r.metrics.RequestsByEndpoint[endpoint]++

	// Check global rate limit first
	if !r.globalLimiter.Allow() {
		r.metrics.RejectedRequests++
		return false, nil
	}

	// Get or create client limiter
	client := r.getOrCreateClientLimiter(clientIP)

	// Check if client is currently blocked
	if client.isBlocked && client.blockedUntil != nil && time.Now().Before(*client.blockedUntil) {
		r.metrics.RejectedRequests++
		client.rejectedCount++
		return false, nil
	}

	// Clear block if expired
	if client.isBlocked && client.blockedUntil != nil && time.Now().After(*client.blockedUntil) {
		client.isBlocked = false
		client.blockedUntil = nil
		client.requestCount = 0
		client.rejectedCount = 0
		client.windowStart = time.Now()
	}

	// Check client-specific rate limit
	if !client.limiter.Allow() {
		r.metrics.RejectedRequests++
		client.rejectedCount++

		// Consider blocking the client if too many rejections
		if client.rejectedCount >= 10 && client.config.BlockDuration > 0 {
			client.isBlocked = true
			blockUntil := time.Now().Add(client.config.BlockDuration)
			client.blockedUntil = &blockUntil
			r.metrics.BlockedClients++
		}

		return false, nil
	}

	// Update client metrics
	client.requestCount++
	client.lastAccess = time.Now()
	r.metrics.AllowedRequests++

	return true, nil
}

// GetRateLimit returns the current rate limit settings for a client
func (r *RateLimiterImpl) GetRateLimit(ctx context.Context, clientIP string) (*interfaces.RateLimitStatus, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	client := r.clientLimiters[clientIP]
	if client == nil {
		// Return default settings for new clients
		return &interfaces.RateLimitStatus{
			ClientIP:       clientIP,
			RequestCount:   0,
			WindowStart:    time.Now(),
			WindowDuration: r.defaultLimit.WindowDuration,
			Limit:          int(r.defaultLimit.RequestsPerSecond * r.defaultLimit.WindowDuration.Seconds()),
			Remaining:      int(r.defaultLimit.RequestsPerSecond * r.defaultLimit.WindowDuration.Seconds()),
			ResetTime:      time.Now().Add(r.defaultLimit.WindowDuration),
			IsBlocked:      false,
		}, nil
	}

	// Calculate remaining tokens
	tokens := client.limiter.Tokens()
	remaining := int(tokens)
	if remaining < 0 {
		remaining = 0
	}

	// Calculate reset time (next window)
	resetTime := client.windowStart.Add(client.config.WindowDuration)
	if time.Now().After(resetTime) {
		resetTime = time.Now().Add(client.config.WindowDuration)
	}

	status := &interfaces.RateLimitStatus{
		ClientIP:       clientIP,
		RequestCount:   client.requestCount,
		WindowStart:    client.windowStart,
		WindowDuration: client.config.WindowDuration,
		Limit:          client.config.BurstSize,
		Remaining:      remaining,
		ResetTime:      resetTime,
		IsBlocked:      client.isBlocked,
		BlockedUntil:   client.blockedUntil,
	}

	return status, nil
}

// SetRateLimit sets a custom rate limit for a client
func (r *RateLimiterImpl) SetRateLimit(ctx context.Context, clientIP string, limit *interfaces.RateLimitConfig) error {
	if limit == nil {
		return fmt.Errorf("limit config cannot be nil")
	}
	if limit.RequestsPerSecond <= 0 {
		return fmt.Errorf("requests per second must be positive")
	}
	if limit.BurstSize <= 0 {
		return fmt.Errorf("burst size must be positive")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Create new client limiter with custom config
	clientLim := &clientLimiter{
		limiter:     rate.NewLimiter(rate.Limit(limit.RequestsPerSecond), limit.BurstSize),
		config:      limit,
		lastAccess:  time.Now(),
		windowStart: time.Now(),
	}

	r.clientLimiters[clientIP] = clientLim

	return nil
}

// ResetRateLimit resets rate limit counters for a client
func (r *RateLimiterImpl) ResetRateLimit(ctx context.Context, clientIP string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	client := r.clientLimiters[clientIP]
	if client == nil {
		return fmt.Errorf("client not found: %s", clientIP)
	}

	// Reset counters and unblock
	client.requestCount = 0
	client.rejectedCount = 0
	client.windowStart = time.Now()
	client.isBlocked = false
	client.blockedUntil = nil

	// Reset the rate limiter
	client.limiter = rate.NewLimiter(rate.Limit(client.config.RequestsPerSecond), client.config.BurstSize)

	return nil
}

// GetGlobalRateStats returns global rate limiting statistics
func (r *RateLimiterImpl) GetGlobalRateStats(ctx context.Context) (*interfaces.RateLimitStats, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Calculate active and blocked clients
	activeClients := 0
	blockedClients := 0
	for _, client := range r.clientLimiters {
		if time.Since(client.lastAccess) <= 10*time.Minute {
			activeClients++
		}
		if client.isBlocked {
			blockedClients++
		}
	}

	// Calculate average latency (placeholder)
	uptime := time.Since(r.startTime)
	avgLatency := time.Duration(0)
	if r.metrics.TotalRequests > 0 {
		avgLatency = uptime / time.Duration(r.metrics.TotalRequests)
	}

	stats := *r.metrics
	stats.ActiveClients = activeClients
	stats.BlockedClients = blockedClients
	stats.AverageLatency = avgLatency
	stats.Timestamp = time.Now()

	return &stats, nil
}

// getOrCreateClientLimiter gets existing or creates new client limiter
func (r *RateLimiterImpl) getOrCreateClientLimiter(clientIP string) *clientLimiter {
	client := r.clientLimiters[clientIP]
	if client != nil {
		return client
	}

	// Create new client limiter with default config
	client = &clientLimiter{
		limiter:     rate.NewLimiter(rate.Limit(r.defaultLimit.RequestsPerSecond), r.defaultLimit.BurstSize),
		config:      r.defaultLimit,
		lastAccess:  time.Now(),
		windowStart: time.Now(),
	}

	r.clientLimiters[clientIP] = client
	return client
}

// cleanupLoop periodically cleans up inactive client limiters
func (r *RateLimiterImpl) cleanupLoop() {
	ticker := time.NewTicker(r.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.cleanup()
		case <-r.stopChan:
			return
		}
	}
}

// cleanup removes inactive client limiters
func (r *RateLimiterImpl) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	cutoff := time.Now().Add(-30 * time.Minute) // Remove clients inactive for 30 minutes
	toDelete := make([]string, 0)

	for clientIP, client := range r.clientLimiters {
		if client.lastAccess.Before(cutoff) && !client.isBlocked {
			toDelete = append(toDelete, clientIP)
		}
	}

	// Delete inactive clients
	for _, clientIP := range toDelete {
		delete(r.clientLimiters, clientIP)
	}
}

// Stop stops the rate limiter
func (r *RateLimiterImpl) Stop() {
	close(r.stopChan)
}
