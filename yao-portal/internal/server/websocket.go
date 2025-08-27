package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/eggybyte-technology/yao-portal/interfaces"
	"github.com/gorilla/websocket"
)

// WebSocketServiceImpl implements the WebSocket service interface
type WebSocketServiceImpl struct {
	mu       sync.RWMutex
	config   *interfaces.PortalConfig
	upgrader *websocket.Upgrader

	// Connection management
	connections   map[string]*WebSocketConnection
	subscriptions map[string]map[string]bool // topic -> connectionId -> bool

	// Metrics
	metrics     *interfaces.WebSocketMetrics
	stopChannel chan struct{}
}

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection struct {
	ID            string
	Conn          *websocket.Conn
	RemoteAddr    string
	UserAgent     string
	ConnectedAt   time.Time
	LastActivity  time.Time
	Subscriptions map[string]bool
	SendChannel   chan []byte
	StopChannel   chan struct{}
	IsActive      bool
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService(config *interfaces.PortalConfig) (interfaces.WebSocketService, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	service := &WebSocketServiceImpl{
		config: config,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// Check CORS origins from config
				origin := r.Header.Get("Origin")
				if len(config.CORSOrigins) == 0 {
					return true
				}
				for _, allowed := range config.CORSOrigins {
					if allowed == "*" || allowed == origin {
						return true
					}
				}
				return false
			},
		},
		connections:   make(map[string]*WebSocketConnection),
		subscriptions: make(map[string]map[string]bool),
		metrics: &interfaces.WebSocketMetrics{
			SubscriptionsByType: make(map[string]int),
		},
		stopChannel: make(chan struct{}),
	}

	// Start cleanup routine
	go service.cleanupLoop()

	return service, nil
}

// HandleWebSocketConnection handles a new WebSocket connection
func (ws *WebSocketServiceImpl) HandleWebSocketConnection(ctx context.Context, conn interface{}) error {
	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		return fmt.Errorf("invalid connection type")
	}

	// Generate unique connection ID
	connectionID := fmt.Sprintf("ws_%d", time.Now().UnixNano())

	connection := &WebSocketConnection{
		ID:            connectionID,
		Conn:          wsConn,
		RemoteAddr:    wsConn.RemoteAddr().String(),
		UserAgent:     "", // TODO: Extract from headers
		ConnectedAt:   time.Now(),
		LastActivity:  time.Now(),
		Subscriptions: make(map[string]bool),
		SendChannel:   make(chan []byte, 256),
		StopChannel:   make(chan struct{}),
		IsActive:      true,
	}

	// Register connection
	ws.mu.Lock()
	ws.connections[connectionID] = connection
	ws.metrics.ActiveConnections++
	ws.metrics.TotalConnections++
	ws.mu.Unlock()

	// Start goroutines for handling the connection
	go ws.handleConnectionReads(connection)
	go ws.handleConnectionWrites(connection)

	// Wait for connection to close
	<-connection.StopChannel

	// Cleanup connection
	ws.unregisterConnection(connectionID)

	return nil
}

// handleConnectionReads handles incoming messages from a WebSocket connection
func (ws *WebSocketServiceImpl) handleConnectionReads(conn *WebSocketConnection) {
	defer func() {
		conn.IsActive = false
		close(conn.StopChannel)
	}()

	// Set read deadline
	conn.Conn.SetReadDeadline(time.Now().Add(ws.config.WSReadTimeout))

	for {
		select {
		case <-conn.StopChannel:
			return
		default:
			var message json.RawMessage
			err := conn.Conn.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					fmt.Printf("WebSocket read error: %v\n", err)
				}
				return
			}

			conn.LastActivity = time.Now()
			conn.Conn.SetReadDeadline(time.Now().Add(ws.config.WSReadTimeout))

			// Parse and handle JSON-RPC message
			go ws.handleJSONRPCMessage(conn, message)
		}
	}
}

// handleConnectionWrites handles outgoing messages to a WebSocket connection
func (ws *WebSocketServiceImpl) handleConnectionWrites(conn *WebSocketConnection) {
	ticker := time.NewTicker(ws.config.WSPingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-conn.StopChannel:
			return
		case message := <-conn.SendChannel:
			conn.Conn.SetWriteDeadline(time.Now().Add(ws.config.WSWriteTimeout))
			if err := conn.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Printf("WebSocket write error: %v\n", err)
				return
			}
			conn.LastActivity = time.Now()
		case <-ticker.C:
			conn.Conn.SetWriteDeadline(time.Now().Add(ws.config.WSWriteTimeout))
			if err := conn.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleJSONRPCMessage handles incoming JSON-RPC messages
func (ws *WebSocketServiceImpl) handleJSONRPCMessage(conn *WebSocketConnection, message json.RawMessage) {
	var request interfaces.JSONRPCRequest
	if err := json.Unmarshal(message, &request); err != nil {
		ws.sendError(conn, nil, -32700, "Parse error", nil)
		return
	}

	switch request.Method {
	case "eth_subscribe":
		ws.handleSubscribe(conn, &request)
	case "eth_unsubscribe":
		ws.handleUnsubscribe(conn, &request)
	default:
		// Handle regular JSON-RPC methods
		// This would integrate with the JSON-RPC service
		ws.sendError(conn, request.ID, -32601, "Method not found", nil)
	}
}

// handleSubscribe handles subscription requests
func (ws *WebSocketServiceImpl) handleSubscribe(conn *WebSocketConnection, request *interfaces.JSONRPCRequest) {
	params, ok := request.Params.([]interface{})
	if !ok || len(params) < 1 {
		ws.sendError(conn, request.ID, -32602, "Invalid params", nil)
		return
	}

	subscriptionType, ok := params[0].(string)
	if !ok {
		ws.sendError(conn, request.ID, -32602, "Invalid subscription type", nil)
		return
	}

	// Validate subscription type
	validSubscriptions := []string{"newHeads", "logs", "pendingTransactions", "syncing"}
	isValid := false
	for _, valid := range validSubscriptions {
		if subscriptionType == valid {
			isValid = true
			break
		}
	}

	if !isValid {
		ws.sendError(conn, request.ID, -32602, "Invalid subscription type", nil)
		return
	}

	// Generate subscription ID
	subscriptionID := fmt.Sprintf("sub_%s_%d", subscriptionType, time.Now().UnixNano())

	// Register subscription
	if err := ws.Subscribe(context.Background(), conn.ID, subscriptionID); err != nil {
		ws.sendError(conn, request.ID, -32603, "Internal error", err.Error())
		return
	}

	// Add subscription type to connection
	conn.Subscriptions[subscriptionType] = true

	// Send success response
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      request.ID,
		"result":  subscriptionID,
	}

	responseBytes, _ := json.Marshal(response)
	select {
	case conn.SendChannel <- responseBytes:
	default:
		// Channel full, connection might be slow
	}
}

// handleUnsubscribe handles unsubscription requests
func (ws *WebSocketServiceImpl) handleUnsubscribe(conn *WebSocketConnection, request *interfaces.JSONRPCRequest) {
	params, ok := request.Params.([]interface{})
	if !ok || len(params) < 1 {
		ws.sendError(conn, request.ID, -32602, "Invalid params", nil)
		return
	}

	subscriptionID, ok := params[0].(string)
	if !ok {
		ws.sendError(conn, request.ID, -32602, "Invalid subscription ID", nil)
		return
	}

	// Unregister subscription
	if err := ws.Unsubscribe(context.Background(), conn.ID, subscriptionID); err != nil {
		ws.sendError(conn, request.ID, -32603, "Internal error", err.Error())
		return
	}

	// Send success response
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      request.ID,
		"result":  true,
	}

	responseBytes, _ := json.Marshal(response)
	select {
	case conn.SendChannel <- responseBytes:
	default:
		// Channel full, connection might be slow
	}
}

// sendError sends an error response to a WebSocket connection
func (ws *WebSocketServiceImpl) sendError(conn *WebSocketConnection, id interface{}, code int, message string, data interface{}) {
	errorResponse := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"data":    data,
		},
	}

	responseBytes, _ := json.Marshal(errorResponse)
	select {
	case conn.SendChannel <- responseBytes:
	default:
		// Channel full, connection might be slow
	}
}

// Subscribe subscribes a connection to a topic
func (ws *WebSocketServiceImpl) Subscribe(ctx context.Context, connectionID string, subscription string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	conn, exists := ws.connections[connectionID]
	if !exists {
		return fmt.Errorf("connection not found")
	}

	conn.Subscriptions[subscription] = true

	// Track subscription by topic
	if ws.subscriptions[subscription] == nil {
		ws.subscriptions[subscription] = make(map[string]bool)
	}
	ws.subscriptions[subscription][connectionID] = true

	// Update metrics
	ws.metrics.ActiveSubscriptions++

	// Update subscription type metrics (extract type from subscription ID)
	if len(subscription) > 4 && subscription[:4] == "sub_" {
		subscriptionType := subscription[4:] // Remove "sub_" prefix
		if idx := len(subscriptionType); idx > 0 {
			// Find the subscription type before the timestamp
			for i, char := range subscriptionType {
				if char >= '0' && char <= '9' {
					subscriptionType = subscriptionType[:i-1] // Remove timestamp part
					break
				}
			}
			ws.metrics.SubscriptionsByType[subscriptionType]++
		}
	}

	return nil
}

// Unsubscribe unsubscribes a connection from a topic
func (ws *WebSocketServiceImpl) Unsubscribe(ctx context.Context, connectionID string, subscription string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	conn, exists := ws.connections[connectionID]
	if !exists {
		return fmt.Errorf("connection not found")
	}

	if _, exists := conn.Subscriptions[subscription]; exists {
		delete(conn.Subscriptions, subscription)

		// Remove from topic subscriptions
		if topicSubs, exists := ws.subscriptions[subscription]; exists {
			delete(topicSubs, connectionID)
			if len(topicSubs) == 0 {
				delete(ws.subscriptions, subscription)
			}
		}

		// Update metrics
		ws.metrics.ActiveSubscriptions--
	}

	return nil
}

// Broadcast broadcasts a message to all subscribers of a topic
func (ws *WebSocketServiceImpl) Broadcast(ctx context.Context, topic string, data interface{}) error {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	// Create broadcast message
	message := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_subscription",
		"params": map[string]interface{}{
			"subscription": topic,
			"result":       data,
		},
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %w", err)
	}

	// Find subscribers for this topic
	subscribers, exists := ws.subscriptions[topic]
	if !exists {
		return nil // No subscribers
	}

	// Send to all subscribers
	for connectionID := range subscribers {
		if conn, exists := ws.connections[connectionID]; exists && conn.IsActive {
			select {
			case conn.SendChannel <- messageBytes:
			default:
				// Channel full, skip this connection
			}
		}
	}

	// Update metrics
	ws.metrics.BroadcastsSent++
	ws.metrics.MessagesSent += uint64(len(subscribers))

	return nil
}

// GetSubscriptions returns all subscriptions for a connection
func (ws *WebSocketServiceImpl) GetSubscriptions(ctx context.Context, connectionID string) ([]string, error) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	conn, exists := ws.connections[connectionID]
	if !exists {
		return nil, fmt.Errorf("connection not found")
	}

	subscriptions := make([]string, 0, len(conn.Subscriptions))
	for subscription := range conn.Subscriptions {
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}

// GetWebSocketMetrics returns WebSocket metrics
func (ws *WebSocketServiceImpl) GetWebSocketMetrics(ctx context.Context) (*interfaces.WebSocketMetrics, error) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	// Create a copy of the metrics
	metrics := &interfaces.WebSocketMetrics{
		ActiveConnections:   ws.metrics.ActiveConnections,
		TotalConnections:    ws.metrics.TotalConnections,
		ConnectionsDropped:  ws.metrics.ConnectionsDropped,
		ActiveSubscriptions: ws.metrics.ActiveSubscriptions,
		SubscriptionsByType: make(map[string]int),
		MessagesSent:        ws.metrics.MessagesSent,
		MessagesReceived:    ws.metrics.MessagesReceived,
		BroadcastsSent:      ws.metrics.BroadcastsSent,
		AverageLatency:      ws.metrics.AverageLatency,
		ConnectionUptime:    ws.metrics.ConnectionUptime,
		Timestamp:           time.Now(),
	}

	for k, v := range ws.metrics.SubscriptionsByType {
		metrics.SubscriptionsByType[k] = v
	}

	return metrics, nil
}

// unregisterConnection removes a connection from the service
func (ws *WebSocketServiceImpl) unregisterConnection(connectionID string) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	conn, exists := ws.connections[connectionID]
	if !exists {
		return
	}

	// Remove all subscriptions for this connection
	for subscription := range conn.Subscriptions {
		if topicSubs, exists := ws.subscriptions[subscription]; exists {
			delete(topicSubs, connectionID)
			if len(topicSubs) == 0 {
				delete(ws.subscriptions, subscription)
			}
		}
	}

	// Update metrics
	ws.metrics.ActiveConnections--
	ws.metrics.ConnectionsDropped++
	ws.metrics.ActiveSubscriptions -= len(conn.Subscriptions)

	// Remove connection
	delete(ws.connections, connectionID)

	// Close connection
	conn.Conn.Close()
}

// cleanupLoop periodically cleans up stale connections
func (ws *WebSocketServiceImpl) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ws.cleanupStaleConnections()
		case <-ws.stopChannel:
			return
		}
	}
}

// cleanupStaleConnections removes connections that haven't been active
func (ws *WebSocketServiceImpl) cleanupStaleConnections() {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	staleThreshold := time.Now().Add(-5 * time.Minute)
	staleConnections := make([]string, 0)

	for connectionID, conn := range ws.connections {
		if conn.LastActivity.Before(staleThreshold) || !conn.IsActive {
			staleConnections = append(staleConnections, connectionID)
		}
	}

	// Remove stale connections
	for _, connectionID := range staleConnections {
		ws.unregisterConnection(connectionID)
	}
}

// Stop stops the WebSocket service
func (ws *WebSocketServiceImpl) Stop() {
	close(ws.stopChannel)

	// Close all active connections
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for _, conn := range ws.connections {
		if conn.IsActive {
			close(conn.StopChannel)
		}
	}
}
