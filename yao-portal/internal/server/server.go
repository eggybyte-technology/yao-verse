package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/eggybyte-technology/yao-portal/interfaces"
	sharedinterfaces "github.com/eggybyte-technology/yao-verse-shared/interfaces"
)

// PortalServer implements the main YaoPortal server
type PortalServer struct {
	mu     sync.RWMutex
	config *interfaces.PortalConfig

	// Core services
	jsonrpcService   *JSONRPCServiceImpl
	txPoolService    *TxPoolServiceImpl
	rateLimiter      *RateLimiterImpl
	backendManager   *BackendServiceManagerImpl
	websocketService *WebSocketServiceImpl

	// HTTP server
	httpServer *http.Server
	wsServer   *http.Server

	// Server state
	status    interfaces.ServerStatus
	startTime time.Time

	// Metrics
	metrics *interfaces.PortalMetrics
}

// NewPortalServer creates a new portal server instance
func NewPortalServer(config *interfaces.PortalConfig) (interfaces.PortalServer, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Create JSON-RPC service
	jsonrpcService, err := NewJSONRPCService(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create JSON-RPC service: %w", err)
	}

	// Create transaction pool service
	txPoolService, err := NewTxPoolService(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction pool service: %w", err)
	}

	// Create rate limiter
	rateLimiter, err := NewRateLimiter(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create rate limiter: %w", err)
	}

	// Create backend service manager
	backendManager, err := NewBackendServiceManager(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create backend service manager: %w", err)
	}

	// Create WebSocket service
	websocketService, err := NewWebSocketService(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create WebSocket service: %w", err)
	}

	server := &PortalServer{
		config:           config,
		jsonrpcService:   jsonrpcService,
		txPoolService:    txPoolService,
		rateLimiter:      rateLimiter,
		backendManager:   backendManager.(*BackendServiceManagerImpl),
		websocketService: websocketService.(*WebSocketServiceImpl),
		status:           interfaces.ServerStatusStopped,
		startTime:        time.Now(),
		metrics: &interfaces.PortalMetrics{
			ServerMetrics:    &interfaces.ServerMetrics{},
			JSONRPCMetrics:   &interfaces.JSONRPCMetrics{},
			TxPoolMetrics:    &interfaces.TxPoolMetrics{},
			WebSocketMetrics: &interfaces.WebSocketMetrics{},
			MessageMetrics:   &interfaces.MessageMetrics{},
		},
	}

	// Wire services together
	jsonrpcService.SetTxPoolService(txPoolService)

	return server, nil
}

// Start starts the portal server
func (s *PortalServer) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status == interfaces.ServerStatusRunning {
		return fmt.Errorf("server is already running")
	}

	s.status = interfaces.ServerStatusStarting

	fmt.Printf("Starting YaoPortal server on %s:%d\n", s.config.HTTPBindAddress, s.config.HTTPPort)

	// TODO: Initialize additional components:
	// - RocketMQ producer for message publishing
	// - YaoOracle for state reading
	// - WebSocket server for subscriptions
	// - Metrics collection system
	// - Health check endpoints

	// Create HTTP server
	mux := http.NewServeMux()

	// Register JSON-RPC endpoint
	mux.HandleFunc("/", s.handleJSONRPC)

	// Register health check endpoint
	mux.HandleFunc("/health", s.handleHealthCheck)

	// Register metrics endpoint
	mux.HandleFunc("/metrics", s.handleMetrics)

	// Register WebSocket endpoint
	mux.HandleFunc("/ws", s.handleWebSocket)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.HTTPBindAddress, s.config.HTTPPort),
		Handler:      mux,
		ReadTimeout:  s.config.RequestTimeout,
		WriteTimeout: s.config.RequestTimeout,
	}

	// Start HTTP server in background
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// If WebSocket is enabled and on different port, start separate WebSocket server
	if s.config.EnableWS && s.config.WSPort != s.config.HTTPPort {
		wsMux := http.NewServeMux()
		wsMux.HandleFunc("/", s.handleWebSocket)

		s.wsServer = &http.Server{
			Addr:         fmt.Sprintf("%s:%d", s.config.WSBindAddress, s.config.WSPort),
			Handler:      wsMux,
			ReadTimeout:  s.config.WSReadTimeout,
			WriteTimeout: s.config.WSWriteTimeout,
		}

		go func() {
			fmt.Printf("Starting WebSocket server on %s:%d\n", s.config.WSBindAddress, s.config.WSPort)
			if err := s.wsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				fmt.Printf("WebSocket server error: %v\n", err)
			}
		}()
	}

	s.status = interfaces.ServerStatusRunning
	s.startTime = time.Now()

	fmt.Println("YaoPortal server started successfully")

	// Block until context is cancelled
	<-ctx.Done()

	return nil
}

// Stop stops the portal server gracefully
func (s *PortalServer) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.status != interfaces.ServerStatusRunning {
		return fmt.Errorf("server is not running")
	}

	s.status = interfaces.ServerStatusStopping

	fmt.Println("Stopping YaoPortal server...")

	// Stop HTTP server
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down HTTP server: %v\n", err)
		}
	}

	// Stop WebSocket server
	if s.wsServer != nil {
		if err := s.wsServer.Shutdown(ctx); err != nil {
			fmt.Printf("Error shutting down WebSocket server: %v\n", err)
		}
	}

	// Stop background services
	if s.txPoolService != nil {
		s.txPoolService.Stop()
	}

	if s.rateLimiter != nil {
		s.rateLimiter.Stop()
	}

	if s.backendManager != nil {
		s.backendManager.Stop()
	}

	if s.websocketService != nil {
		s.websocketService.Stop()
	}

	s.status = interfaces.ServerStatusStopped

	fmt.Println("YaoPortal server stopped successfully")

	return nil
}

// GetJSONRPCService returns the JSON-RPC service instance
func (s *PortalServer) GetJSONRPCService() interfaces.JSONRPCService {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jsonrpcService
}

// GetTxPoolService returns the transaction pool service instance
func (s *PortalServer) GetTxPoolService() interfaces.TxPoolService {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.txPoolService
}

// GetMessageProducer returns the message producer for publishing transactions
func (s *PortalServer) GetMessageProducer() sharedinterfaces.MessageProducer {
	// TODO: Return actual message producer implementation
	return nil
}

// GetStateReader returns the state reader for queries
func (s *PortalServer) GetStateReader() sharedinterfaces.StateReader {
	// TODO: Return actual state reader implementation
	return nil
}

// RegisterWebSocketClient registers a WebSocket client for subscriptions
func (s *PortalServer) RegisterWebSocketClient(ctx context.Context, clientID string, subscriptions []string) error {
	// TODO: Implement WebSocket client registration
	return fmt.Errorf("not implemented")
}

// UnregisterWebSocketClient unregisters a WebSocket client
func (s *PortalServer) UnregisterWebSocketClient(ctx context.Context, clientID string) error {
	// TODO: Implement WebSocket client unregistration
	return fmt.Errorf("not implemented")
}

// BroadcastToClients broadcasts a message to all subscribed clients
func (s *PortalServer) BroadcastToClients(ctx context.Context, subscription string, data interface{}) error {
	// TODO: Implement broadcasting to WebSocket clients
	return fmt.Errorf("not implemented")
}

// GetStatus returns the current server status
func (s *PortalServer) GetStatus(ctx context.Context) *interfaces.PortalStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	txPoolStatus, _ := s.txPoolService.GetPoolStatus(ctx)

	return &interfaces.PortalStatus{
		ServerStatus:        s.status,
		ConnectedClients:    0, // TODO: Get actual count
		ActiveSubscriptions: 0, // TODO: Get actual count
		RequestsPerSecond:   0, // TODO: Calculate actual RPS
		TxPoolStatus:        txPoolStatus,
		Uptime:              time.Since(s.startTime),
		StartTime:           s.startTime,
		Version:             "v1.0.0", // TODO: Get from build info
	}
}

// GetMetrics returns comprehensive server metrics
func (s *PortalServer) GetMetrics(ctx context.Context) *interfaces.PortalMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Collect metrics from all services
	serverMetrics := &interfaces.ServerMetrics{
		Status:          s.status,
		Uptime:          time.Since(s.startTime),
		CPUUsage:        0, // TODO: Implement actual CPU monitoring
		MemoryUsage:     0, // TODO: Implement actual memory monitoring
		GoroutineCount:  0, // TODO: Get goroutine count
		GCPauseTime:     0, // TODO: Get GC pause time
		RequestsHandled: 0, // TODO: Track requests handled
		ErrorsCount:     0, // TODO: Track errors
		Timestamp:       time.Now(),
	}

	jsonrpcMetrics := s.jsonrpcService.GetMetrics()
	txPoolMetrics := s.txPoolService.GetMetrics()

	// TODO: Collect WebSocket and Message metrics

	return &interfaces.PortalMetrics{
		ServerMetrics:    serverMetrics,
		JSONRPCMetrics:   jsonrpcMetrics,
		TxPoolMetrics:    txPoolMetrics,
		WebSocketMetrics: &interfaces.WebSocketMetrics{}, // TODO: Implement
		MessageMetrics:   &interfaces.MessageMetrics{},   // TODO: Implement
	}
}

// HealthCheck performs comprehensive health check
func (s *PortalServer) HealthCheck(ctx context.Context) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.status != interfaces.ServerStatusRunning {
		return fmt.Errorf("server is not running")
	}

	// TODO: Implement comprehensive health checks:
	// - Check RocketMQ connectivity
	// - Check YaoOracle health
	// - Check transaction pool health
	// - Check rate limiter health

	return nil
}

// HTTP handlers

// handleJSONRPC handles JSON-RPC requests
func (s *PortalServer) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	// Check rate limits
	clientIP := getClientIP(r)
	allowed, err := s.rateLimiter.AllowRequest(r.Context(), clientIP, r.URL.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rate limit error: %v", err), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Delegate to JSON-RPC service
	s.jsonrpcService.HandleHTTPRequest(w, r)
}

// handleHealthCheck handles health check requests
func (s *PortalServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.HealthCheck(r.Context()); err != nil {
		http.Error(w, fmt.Sprintf("Health check failed: %v", err), http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}

// handleMetrics handles metrics requests
func (s *PortalServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metricsData := s.GetMetrics(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// TODO: Serialize metrics to JSON
	_ = metricsData // Acknowledge the variable is used
	fmt.Fprintf(w, `{"message":"Metrics endpoint - TODO: implement JSON serialization"}`)
}

// handleWebSocket handles WebSocket connections
func (s *PortalServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Check rate limits
	clientIP := getClientIP(r)
	allowed, err := s.rateLimiter.AllowRequest(r.Context(), clientIP, r.URL.Path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Rate limit error: %v", err), http.StatusInternalServerError)
		return
	}
	if !allowed {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := s.websocketService.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	// Handle the WebSocket connection
	if err := s.websocketService.HandleWebSocketConnection(r.Context(), conn); err != nil {
		fmt.Printf("WebSocket connection error: %v\n", err)
	}
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}

	// Check X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
