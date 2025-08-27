package server

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/eggybyte-technology/yao-portal/interfaces"
	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
)

// TxPoolServiceImpl implements the transaction pool service interface
type TxPoolServiceImpl struct {
	mu     sync.RWMutex
	config *interfaces.PortalConfig

	// Transaction storage
	pendingTxs   map[common.Hash]*types.Transaction
	queuedTxs    map[common.Hash]*types.Transaction
	txStatus     map[common.Hash]*interfaces.TxStatus
	txTimestamps map[common.Hash]time.Time // Track submission times

	// Pool management
	capacity      int
	minGasPrice   *big.Int
	maxTxLifetime time.Duration

	// Metrics
	metrics   *interfaces.TxPoolMetrics
	startTime time.Time

	// Channels for background processing
	addTxChan    chan *types.Transaction
	removeTxChan chan common.Hash
	stopChan     chan struct{}
}

// NewTxPoolService creates a new transaction pool service instance
func NewTxPoolService(config *interfaces.PortalConfig) (*TxPoolServiceImpl, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	pool := &TxPoolServiceImpl{
		config:        config,
		pendingTxs:    make(map[common.Hash]*types.Transaction),
		queuedTxs:     make(map[common.Hash]*types.Transaction),
		txStatus:      make(map[common.Hash]*interfaces.TxStatus),
		txTimestamps:  make(map[common.Hash]time.Time),
		capacity:      config.TxPoolCapacity,
		minGasPrice:   big.NewInt(1000000000), // 1 gwei default
		maxTxLifetime: config.MaxTxLifetime,
		startTime:     time.Now(),
		addTxChan:     make(chan *types.Transaction, 1000),
		removeTxChan:  make(chan common.Hash, 1000),
		stopChan:      make(chan struct{}),
		metrics: &interfaces.TxPoolMetrics{
			AverageGasPrice: big.NewInt(0),
			MinGasPrice:     big.NewInt(1000000000),
			MaxGasPrice:     big.NewInt(0),
		},
	}

	// Start background processing goroutine
	go pool.processingLoop()

	// Start cleanup goroutine
	go pool.cleanupLoop()

	return pool, nil
}

// AddTransaction adds a transaction to the pool
func (t *TxPoolServiceImpl) AddTransaction(ctx context.Context, tx *types.Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction cannot be nil")
	}

	// Basic validation
	if err := t.ValidateTransaction(ctx, tx); err != nil {
		t.updateErrorMetrics("validation", err)
		return fmt.Errorf("transaction validation failed: %w", err)
	}

	select {
	case t.addTxChan <- tx:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		t.updateErrorMetrics("pool_full", fmt.Errorf("transaction pool is full"))
		return fmt.Errorf("transaction pool is full")
	}
}

// RemoveTransaction removes a transaction from the pool
func (t *TxPoolServiceImpl) RemoveTransaction(ctx context.Context, hash common.Hash) error {
	select {
	case t.removeTxChan <- hash:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("failed to remove transaction")
	}
}

// GetTransaction returns a transaction from the pool
func (t *TxPoolServiceImpl) GetTransaction(ctx context.Context, hash common.Hash) (*types.Transaction, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Check pending transactions first
	if tx, exists := t.pendingTxs[hash]; exists {
		return tx, nil
	}

	// Check queued transactions
	if tx, exists := t.queuedTxs[hash]; exists {
		return tx, nil
	}

	return nil, fmt.Errorf("transaction not found in pool")
}

// GetPendingTransactions returns all pending transactions
func (t *TxPoolServiceImpl) GetPendingTransactions(ctx context.Context) ([]*types.Transaction, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	txs := make([]*types.Transaction, 0, len(t.pendingTxs))
	for _, tx := range t.pendingTxs {
		txs = append(txs, tx)
	}

	return txs, nil
}

// GetQueuedTransactions returns all queued transactions
func (t *TxPoolServiceImpl) GetQueuedTransactions(ctx context.Context) ([]*types.Transaction, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	txs := make([]*types.Transaction, 0, len(t.queuedTxs))
	for _, tx := range t.queuedTxs {
		txs = append(txs, tx)
	}

	return txs, nil
}

// GetTransactionStatus returns the status of a transaction
func (t *TxPoolServiceImpl) GetTransactionStatus(ctx context.Context, hash common.Hash) (*interfaces.TxStatus, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if status, exists := t.txStatus[hash]; exists {
		return status, nil
	}

	return nil, fmt.Errorf("transaction status not found")
}

// ValidateTransaction validates a transaction
func (t *TxPoolServiceImpl) ValidateTransaction(ctx context.Context, tx *types.Transaction) error {
	// Basic validation checks
	if tx.Gas == 0 {
		return fmt.Errorf("gas limit cannot be zero")
	}

	if tx.GasPrice == nil || tx.GasPrice.Cmp(t.minGasPrice) < 0 {
		return fmt.Errorf("gas price too low: minimum %s", t.minGasPrice.String())
	}

	// Check transaction size
	// TODO: Implement transaction size validation

	// Check nonce validation
	// TODO: Implement nonce validation against current state

	// Check balance validation
	// TODO: Implement balance validation against current state

	return nil
}

// GetPoolStatus returns the current pool status
func (t *TxPoolServiceImpl) GetPoolStatus(ctx context.Context) (*interfaces.TxPoolStatus, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	pendingCount := len(t.pendingTxs)
	queuedCount := len(t.queuedTxs)
	totalCount := pendingCount + queuedCount

	usagePercent := float64(totalCount) / float64(t.capacity) * 100

	return &interfaces.TxPoolStatus{
		PendingCount:  pendingCount,
		QueuedCount:   queuedCount,
		Capacity:      t.capacity,
		UsagePercent:  usagePercent,
		MinGasPrice:   new(big.Int).Set(t.minGasPrice),
		MaxTxLifetime: int64(t.maxTxLifetime.Seconds()),
	}, nil
}

// SetGasPrice sets the minimum gas price
func (t *TxPoolServiceImpl) SetGasPrice(gasPrice *big.Int) error {
	if gasPrice == nil || gasPrice.Sign() <= 0 {
		return fmt.Errorf("gas price must be positive")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.minGasPrice = new(big.Int).Set(gasPrice)
	t.metrics.MinGasPrice = new(big.Int).Set(gasPrice)

	return nil
}

// GetGasPrice returns the minimum gas price
func (t *TxPoolServiceImpl) GetGasPrice() *big.Int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return new(big.Int).Set(t.minGasPrice)
}

// Clear clears all transactions from the pool
func (t *TxPoolServiceImpl) Clear(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Update metrics before clearing
	t.metrics.DroppedTransactions += uint64(len(t.pendingTxs) + len(t.queuedTxs))

	// Clear all maps
	t.pendingTxs = make(map[common.Hash]*types.Transaction)
	t.queuedTxs = make(map[common.Hash]*types.Transaction)
	t.txStatus = make(map[common.Hash]*interfaces.TxStatus)
	t.txTimestamps = make(map[common.Hash]time.Time)

	return nil
}

// GetMetrics returns transaction pool metrics
func (t *TxPoolServiceImpl) GetMetrics() *interfaces.TxPoolMetrics {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Calculate uptime and transaction rate
	uptime := time.Since(t.startTime)
	tps := float64(t.metrics.TotalTransactions) / uptime.Seconds()

	metrics := *t.metrics
	metrics.TransactionsPerSecond = tps
	metrics.PendingTransactions = len(t.pendingTxs)
	metrics.QueuedTransactions = len(t.queuedTxs)
	metrics.PoolUtilization = float64(len(t.pendingTxs)+len(t.queuedTxs)) / float64(t.capacity) * 100

	return &metrics
}

// processingLoop handles background transaction processing
func (t *TxPoolServiceImpl) processingLoop() {
	for {
		select {
		case tx := <-t.addTxChan:
			t.processAddTransaction(tx)
		case hash := <-t.removeTxChan:
			t.processRemoveTransaction(hash)
		case <-t.stopChan:
			return
		}
	}
}

// processAddTransaction processes adding a transaction to the pool
func (t *TxPoolServiceImpl) processAddTransaction(tx *types.Transaction) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Calculate transaction hash (assuming tx.Hash contains the hash)
	hash := tx.Hash

	// Check if transaction already exists
	if _, exists := t.pendingTxs[hash]; exists {
		return
	}
	if _, exists := t.queuedTxs[hash]; exists {
		return
	}

	// Check capacity
	totalCount := len(t.pendingTxs) + len(t.queuedTxs)
	if totalCount >= t.capacity {
		// Remove oldest transaction to make space
		t.evictOldestTransaction()
	}

	// Add to pending pool (for now, all transactions go to pending)
	t.pendingTxs[hash] = tx

	// Update status
	t.txStatus[hash] = &interfaces.TxStatus{
		Hash:   hash,
		Status: interfaces.TxStatusPending,
	}
	t.txTimestamps[hash] = time.Now()

	// Update metrics
	t.metrics.TotalTransactions++
	if tx.GasPrice != nil {
		t.updateGasMetrics(tx.GasPrice)
	}
}

// processRemoveTransaction processes removing a transaction from the pool
func (t *TxPoolServiceImpl) processRemoveTransaction(hash common.Hash) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Remove from pending
	if _, exists := t.pendingTxs[hash]; exists {
		delete(t.pendingTxs, hash)
		delete(t.txStatus, hash)
		delete(t.txTimestamps, hash)
		return
	}

	// Remove from queued
	if _, exists := t.queuedTxs[hash]; exists {
		delete(t.queuedTxs, hash)
		delete(t.txStatus, hash)
		delete(t.txTimestamps, hash)
		return
	}
}

// cleanupLoop periodically cleans up expired transactions
func (t *TxPoolServiceImpl) cleanupLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.cleanupExpiredTransactions()
		case <-t.stopChan:
			return
		}
	}
}

// cleanupExpiredTransactions removes expired transactions
func (t *TxPoolServiceImpl) cleanupExpiredTransactions() {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	expired := make([]common.Hash, 0)

	// Check pending transactions
	for hash, timestamp := range t.txTimestamps {
		if now.Sub(timestamp) > t.maxTxLifetime {
			expired = append(expired, hash)
		}
	}

	// Remove expired transactions
	for _, hash := range expired {
		delete(t.pendingTxs, hash)
		delete(t.queuedTxs, hash)
		delete(t.txStatus, hash)
		delete(t.txTimestamps, hash)
		t.metrics.DroppedTransactions++
	}
}

// evictOldestTransaction removes the oldest transaction to make space
func (t *TxPoolServiceImpl) evictOldestTransaction() {
	var oldestHash common.Hash
	var oldestTime time.Time

	// Find oldest transaction
	for hash, timestamp := range t.txTimestamps {
		if oldestTime.IsZero() || timestamp.Before(oldestTime) {
			oldestTime = timestamp
			oldestHash = hash
		}
	}

	if oldestHash != (common.Hash{}) {
		delete(t.pendingTxs, oldestHash)
		delete(t.queuedTxs, oldestHash)
		delete(t.txStatus, oldestHash)
		delete(t.txTimestamps, oldestHash)
		t.metrics.DroppedTransactions++
	}
}

// updateGasMetrics updates gas-related metrics
func (t *TxPoolServiceImpl) updateGasMetrics(gasPrice *big.Int) {
	if gasPrice == nil {
		return
	}

	// Update min gas price
	if t.metrics.MinGasPrice.Cmp(gasPrice) > 0 {
		t.metrics.MinGasPrice = new(big.Int).Set(gasPrice)
	}

	// Update max gas price
	if t.metrics.MaxGasPrice.Cmp(gasPrice) < 0 {
		t.metrics.MaxGasPrice = new(big.Int).Set(gasPrice)
	}

	// TODO: Update average gas price calculation
}

// updateErrorMetrics updates error metrics
func (t *TxPoolServiceImpl) updateErrorMetrics(errorType string, err error) {
	switch errorType {
	case "validation":
		t.metrics.ValidationErrors++
	case "insufficient_funds":
		t.metrics.InsufficientFundsErrors++
	case "nonce":
		t.metrics.NonceErrors++
	case "gas_price":
		t.metrics.GasPriceErrors++
	}
}

// Stop stops the transaction pool service
func (t *TxPoolServiceImpl) Stop() {
	close(t.stopChan)
}
