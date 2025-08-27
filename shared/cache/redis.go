package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/interfaces"
	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/redis/go-redis/v9"
)

// RedisCache implements L2 shared cache using Redis Cluster
type RedisCache struct {
	client *redis.ClusterClient
	config *RedisConfig
	stats  *interfaces.OracleCacheStats
}

// RedisConfig defines configuration for Redis cluster cache
type RedisConfig struct {
	Endpoints []string `json:"endpoints"`
	Password  string   `json:"password"`
	DB        int      `json:"db"` // DB is ignored in cluster mode but kept for compatibility
	PoolSize  int      `json:"poolSize"`
	TTL       int      `json:"ttl"` // in seconds
}

// NewRedisCache creates a new Redis cluster cache instance
func NewRedisCache(config *RedisConfig) (*RedisCache, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    config.Endpoints,
		Password: config.Password,
		PoolSize: config.PoolSize,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis cluster: %w", err)
	}

	return &RedisCache{
		client: client,
		config: config,
		stats:  &interfaces.OracleCacheStats{},
	}, nil
}

// GetAccount retrieves an account from cache
func (r *RedisCache) GetAccount(ctx context.Context, address common.Address) (*types.Account, error) {
	key := r.accountKey(address)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.stats.Misses++
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get account from Redis: %w", err)
	}

	var account types.Account
	err = json.Unmarshal([]byte(data), &account)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal account: %w", err)
	}

	r.stats.Hits++
	return &account, nil
}

// SetAccount stores an account in cache
func (r *RedisCache) SetAccount(ctx context.Context, address common.Address, account *types.Account) error {
	key := r.accountKey(address)
	data, err := json.Marshal(account)
	if err != nil {
		return fmt.Errorf("failed to marshal account: %w", err)
	}

	ttl := time.Duration(r.config.TTL) * time.Second
	err = r.client.Set(ctx, key, string(data), ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set account in Redis: %w", err)
	}

	return nil
}

// GetStorage retrieves storage value from cache
func (r *RedisCache) GetStorage(ctx context.Context, address common.Address, key common.Hash) (common.Hash, error) {
	cacheKey := r.storageKey(address, key)
	data, err := r.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			r.stats.Misses++
			return common.Hash{}, nil // Cache miss
		}
		return common.Hash{}, fmt.Errorf("failed to get storage from Redis: %w", err)
	}

	r.stats.Hits++
	return common.HexToHash(data), nil
}

// SetStorage stores storage value in cache
func (r *RedisCache) SetStorage(ctx context.Context, address common.Address, key common.Hash, value common.Hash) error {
	cacheKey := r.storageKey(address, key)
	ttl := time.Duration(r.config.TTL) * time.Second
	err := r.client.Set(ctx, cacheKey, value.Hex(), ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set storage in Redis: %w", err)
	}

	return nil
}

// GetCode retrieves contract code from cache
func (r *RedisCache) GetCode(ctx context.Context, address common.Address) ([]byte, error) {
	key := r.codeKey(address)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.stats.Misses++
			return nil, nil // Cache miss
		}
		return nil, fmt.Errorf("failed to get code from Redis: %w", err)
	}

	r.stats.Hits++
	return common.FromHex(data), nil
}

// SetCode stores contract code in cache
func (r *RedisCache) SetCode(ctx context.Context, address common.Address, code []byte) error {
	key := r.codeKey(address)
	ttl := time.Duration(r.config.TTL) * time.Second
	err := r.client.Set(ctx, key, common.Bytes2Hex(code), ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set code in Redis: %w", err)
	}

	return nil
}

// GetCodeHash retrieves contract code hash from cache
func (r *RedisCache) GetCodeHash(ctx context.Context, address common.Address) (common.Hash, error) {
	key := r.codeHashKey(address)
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			r.stats.Misses++
			return common.Hash{}, nil // Cache miss
		}
		return common.Hash{}, fmt.Errorf("failed to get code hash from Redis: %w", err)
	}

	r.stats.Hits++
	return common.HexToHash(data), nil
}

// SetCodeHash stores contract code hash in cache
func (r *RedisCache) SetCodeHash(ctx context.Context, address common.Address, codeHash common.Hash) error {
	key := r.codeHashKey(address)
	ttl := time.Duration(r.config.TTL) * time.Second
	err := r.client.Set(ctx, key, codeHash.Hex(), ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to set code hash in Redis: %w", err)
	}

	return nil
}

// Delete removes a key from cache
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("failed to delete keys from Redis: %w", err)
	}

	return nil
}

// DeleteByPattern removes keys matching a pattern (cluster-safe implementation)
func (r *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	// For Redis cluster, we need to scan all master nodes
	var allKeys []string
	err := r.client.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
		keys, err := client.Keys(ctx, pattern).Result()
		if err != nil {
			return err
		}
		allKeys = append(allKeys, keys...)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to get keys by pattern: %w", err)
	}

	if len(allKeys) == 0 {
		return nil
	}

	// Delete keys in batches to avoid too large commands
	batchSize := 100
	for i := 0; i < len(allKeys); i += batchSize {
		end := i + batchSize
		if end > len(allKeys) {
			end = len(allKeys)
		}

		batch := allKeys[i:end]
		err = r.client.Del(ctx, batch...).Err()
		if err != nil {
			return fmt.Errorf("failed to delete keys by pattern: %w", err)
		}
	}

	return nil
}

// InvalidateKeys invalidates multiple cache keys
func (r *RedisCache) InvalidateKeys(ctx context.Context, keys []*types.CacheKey) error {
	if len(keys) == 0 {
		return nil
	}

	cacheKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		switch key.Type {
		case types.CacheKeyTypeAccount:
			cacheKeys = append(cacheKeys, r.accountKey(common.HexToAddress(key.Address)))
		case types.CacheKeyTypeStorage:
			cacheKeys = append(cacheKeys, r.storageKey(common.HexToAddress(key.Address), common.HexToHash(key.Key)))
		case types.CacheKeyTypeCode:
			cacheKeys = append(cacheKeys, r.codeKey(common.HexToAddress(key.Address)))
		}
	}

	return r.Delete(ctx, cacheKeys...)
}

// Clear removes all keys from cache (cluster-safe implementation)
func (r *RedisCache) Clear(ctx context.Context) error {
	// For Redis cluster, we need to flush all master nodes
	return r.client.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
		return client.FlushDB(ctx).Err()
	})
}

// GetStats returns cache statistics
func (r *RedisCache) GetStats() *interfaces.OracleCacheStats {
	r.updateStats()
	return r.stats
}

// Ping tests the connection to Redis cluster
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// Helper methods for generating cache keys

func (r *RedisCache) accountKey(address common.Address) string {
	return fmt.Sprintf("yao:account:%s", address.Hex())
}

func (r *RedisCache) storageKey(address common.Address, key common.Hash) string {
	return fmt.Sprintf("yao:storage:%s:%s", address.Hex(), key.Hex())
}

func (r *RedisCache) codeKey(address common.Address) string {
	return fmt.Sprintf("yao:code:%s", address.Hex())
}

func (r *RedisCache) codeHashKey(address common.Address) string {
	return fmt.Sprintf("yao:codehash:%s", address.Hex())
}

// updateStats updates cache statistics
func (r *RedisCache) updateStats() {
	// Calculate hit ratio
	total := r.stats.Hits + r.stats.Misses
	if total > 0 {
		r.stats.HitRatio = float64(r.stats.Hits) / float64(total)
	}

	// Get cluster info for size estimation
	ctx := context.Background()
	var totalMemory uint64
	r.client.ForEachMaster(ctx, func(ctx context.Context, client *redis.Client) error {
		info, err := client.Info(ctx, "memory").Result()
		if err == nil {
			// Parse memory usage from info string (simplified)
			totalMemory += r.parseMemoryUsage(info)
		}
		return nil // Continue even if one node fails
	})

	r.stats.Size = totalMemory
}

// parseMemoryUsage parses memory usage from Redis INFO command output
func (r *RedisCache) parseMemoryUsage(info string) uint64 {
	// This is a simplified implementation
	// In a real implementation, you would parse the INFO output properly
	return 0
}
