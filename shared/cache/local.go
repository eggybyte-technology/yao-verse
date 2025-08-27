package cache

import (
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/eggybyte-technology/yao-verse-shared/interfaces"
	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
)

// LocalCache implements L1 local in-memory cache using Ristretto
type LocalCache struct {
	cache  *ristretto.Cache[string, interface{}]
	config *LocalCacheConfig
	stats  *interfaces.OracleCacheStats
	mutex  sync.RWMutex
}

// LocalCacheConfig defines configuration for local cache
type LocalCacheConfig struct {
	MaxSize     int64 `json:"maxSize"`     // Maximum cache size in bytes
	NumCounters int64 `json:"numCounters"` // Number of counters for admission policy
	MaxCost     int64 `json:"maxCost"`     // Maximum cost for items
	BufferItems int64 `json:"bufferItems"` // Buffer size for updates
	TTL         int   `json:"ttl"`         // Time to live in seconds
}

// NewLocalCache creates a new local cache instance
func NewLocalCache(config *LocalCacheConfig) (*LocalCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, interface{}]{
		NumCounters: config.NumCounters, // 10x the number of items you expect to keep in cache
		MaxCost:     config.MaxCost,     // Maximum cost of cache (1 for each item)
		BufferItems: config.BufferItems, // 64 is a good default buffer size
	})
	if err != nil {
		return nil, err
	}

	return &LocalCache{
		cache:  cache,
		config: config,
		stats:  &interfaces.OracleCacheStats{},
	}, nil
}

// GetAccount retrieves an account from local cache
func (l *LocalCache) GetAccount(address common.Address) (*types.Account, bool) {
	key := l.accountKey(address)
	value, found := l.cache.Get(key)
	if !found {
		l.mutex.Lock()
		l.stats.Misses++
		l.mutex.Unlock()
		return nil, false
	}

	account, ok := value.(*types.Account)
	if !ok {
		l.mutex.Lock()
		l.stats.Misses++
		l.mutex.Unlock()
		return nil, false
	}

	l.mutex.Lock()
	l.stats.Hits++
	l.mutex.Unlock()
	return account, true
}

// SetAccount stores an account in local cache
func (l *LocalCache) SetAccount(address common.Address, account *types.Account) {
	key := l.accountKey(address)

	// Calculate the cost (approximate size in bytes)
	cost := l.calculateAccountCost(account)

	// Set TTL if configured
	var ttl time.Duration
	if l.config.TTL > 0 {
		ttl = time.Duration(l.config.TTL) * time.Second
	}

	l.cache.SetWithTTL(key, account, cost, ttl)
}

// GetStorage retrieves storage value from local cache
func (l *LocalCache) GetStorage(address common.Address, key common.Hash) (common.Hash, bool) {
	cacheKey := l.storageKey(address, key)
	value, found := l.cache.Get(cacheKey)
	if !found {
		l.mutex.Lock()
		l.stats.Misses++
		l.mutex.Unlock()
		return common.Hash{}, false
	}

	hash, ok := value.(common.Hash)
	if !ok {
		l.mutex.Lock()
		l.stats.Misses++
		l.mutex.Unlock()
		return common.Hash{}, false
	}

	l.mutex.Lock()
	l.stats.Hits++
	l.mutex.Unlock()
	return hash, true
}

// SetStorage stores storage value in local cache
func (l *LocalCache) SetStorage(address common.Address, key common.Hash, value common.Hash) {
	cacheKey := l.storageKey(address, key)

	// Storage entries are typically 32 bytes
	cost := int64(32)

	// Set TTL if configured
	var ttl time.Duration
	if l.config.TTL > 0 {
		ttl = time.Duration(l.config.TTL) * time.Second
	}

	l.cache.SetWithTTL(cacheKey, value, cost, ttl)
}

// GetCode retrieves contract code from local cache
func (l *LocalCache) GetCode(address common.Address) ([]byte, bool) {
	key := l.codeKey(address)
	value, found := l.cache.Get(key)
	if !found {
		l.mutex.Lock()
		l.stats.Misses++
		l.mutex.Unlock()
		return nil, false
	}

	code, ok := value.([]byte)
	if !ok {
		l.mutex.Lock()
		l.stats.Misses++
		l.mutex.Unlock()
		return nil, false
	}

	l.mutex.Lock()
	l.stats.Hits++
	l.mutex.Unlock()
	return code, true
}

// SetCode stores contract code in local cache
func (l *LocalCache) SetCode(address common.Address, code []byte) {
	key := l.codeKey(address)

	// Calculate cost based on code length
	cost := int64(len(code))

	// Set TTL if configured
	var ttl time.Duration
	if l.config.TTL > 0 {
		ttl = time.Duration(l.config.TTL) * time.Second
	}

	l.cache.SetWithTTL(key, code, cost, ttl)
}

// GetCodeHash retrieves contract code hash from local cache
func (l *LocalCache) GetCodeHash(address common.Address) (common.Hash, bool) {
	key := l.codeHashKey(address)
	value, found := l.cache.Get(key)
	if !found {
		l.mutex.Lock()
		l.stats.Misses++
		l.mutex.Unlock()
		return common.Hash{}, false
	}

	hash, ok := value.(common.Hash)
	if !ok {
		l.mutex.Lock()
		l.stats.Misses++
		l.mutex.Unlock()
		return common.Hash{}, false
	}

	l.mutex.Lock()
	l.stats.Hits++
	l.mutex.Unlock()
	return hash, true
}

// SetCodeHash stores contract code hash in local cache
func (l *LocalCache) SetCodeHash(address common.Address, codeHash common.Hash) {
	key := l.codeHashKey(address)

	// Hash is 32 bytes
	cost := int64(32)

	// Set TTL if configured
	var ttl time.Duration
	if l.config.TTL > 0 {
		ttl = time.Duration(l.config.TTL) * time.Second
	}

	l.cache.SetWithTTL(key, codeHash, cost, ttl)
}

// Delete removes keys from local cache
func (l *LocalCache) Delete(keys ...string) {
	for _, key := range keys {
		l.cache.Del(key)
	}
}

// InvalidateKeys invalidates multiple cache keys
func (l *LocalCache) InvalidateKeys(keys []*types.CacheKey) {
	cacheKeys := make([]string, len(keys))
	for i, key := range keys {
		switch key.Type {
		case types.CacheKeyTypeAccount:
			cacheKeys[i] = l.accountKey(common.HexToAddress(key.Address))
		case types.CacheKeyTypeStorage:
			cacheKeys[i] = l.storageKey(common.HexToAddress(key.Address), common.HexToHash(key.Key))
		case types.CacheKeyTypeCode:
			cacheKeys[i] = l.codeKey(common.HexToAddress(key.Address))
		default:
			continue
		}
	}

	l.Delete(cacheKeys...)
}

// Clear removes all items from local cache
func (l *LocalCache) Clear() {
	l.cache.Clear()
}

// GetStats returns cache statistics
func (l *LocalCache) GetStats() *interfaces.OracleCacheStats {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	metrics := l.cache.Metrics
	total := l.stats.Hits + l.stats.Misses
	hitRatio := float64(0)
	if total > 0 {
		hitRatio = float64(l.stats.Hits) / float64(total)
	}

	return &interfaces.OracleCacheStats{
		Hits:     l.stats.Hits,
		Misses:   l.stats.Misses,
		HitRatio: hitRatio,
		Size:     uint64(metrics.CostAdded() - metrics.CostEvicted()),
		Capacity: uint64(l.config.MaxCost),
	}
}

// Close closes the cache
func (l *LocalCache) Close() {
	l.cache.Close()
}

// Helper methods for generating cache keys

func (l *LocalCache) accountKey(address common.Address) string {
	return "account:" + address.Hex()
}

func (l *LocalCache) storageKey(address common.Address, key common.Hash) string {
	return "storage:" + address.Hex() + ":" + key.Hex()
}

func (l *LocalCache) codeKey(address common.Address) string {
	return "code:" + address.Hex()
}

func (l *LocalCache) codeHashKey(address common.Address) string {
	return "codehash:" + address.Hex()
}

// calculateAccountCost estimates the memory cost of an account
func (l *LocalCache) calculateAccountCost(account *types.Account) int64 {
	// Address (20 bytes) + Nonce (8 bytes) + Balance (variable, estimate 32 bytes) +
	// CodeHash (32 bytes) + StorageRoot (32 bytes) + overhead
	return 20 + 8 + 32 + 32 + 32 + 50 // 174 bytes approximately
}
