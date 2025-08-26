package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/eggybyte-technology/yao-verse-shared/interfaces"
	"github.com/eggybyte-technology/yao-verse-shared/oracle"
	"github.com/eggybyte-technology/yao-verse-shared/types"
	"github.com/ethereum/go-ethereum/common"
)

// Example demonstrating how to use YaoOracle
func main() {
	// Create YaoOracle instance
	yaoOracle, err := oracle.NewYaoOracle()
	if err != nil {
		log.Fatal("Failed to create YaoOracle:", err)
	}
	defer yaoOracle.Close()

	// Configuration for YaoOracle
	config := &interfaces.OracleConfig{
		// L1 Cache configuration
		L1CacheSize: 10000, // 10K items
		L1CacheTTL:  300,   // 5 minutes

		// L2 Cache (Redis) configuration
		RedisEndpoints: []string{"localhost:6379"},
		RedisPassword:  "",
		RedisDB:        0,
		RedisPoolSize:  10,
		L2CacheTTL:     1800, // 30 minutes

		// L3 Storage (TiDB) configuration
		TiDBDSN:         "root:@tcp(localhost:4000)/yaoverse",
		TiDBMaxOpenConn: 100,
		TiDBMaxIdleConn: 10,
		TiDBMaxLifetime: 3600, // 1 hour

		// Performance tuning
		BatchSize:       100,
		EnableMetrics:   true,
		MetricsInterval: 10, // 10 seconds
	}

	// Initialize YaoOracle
	ctx := context.Background()
	err = yaoOracle.Initialize(ctx, config)
	if err != nil {
		log.Fatal("Failed to initialize YaoOracle:", err)
	}

	// Start background services
	err = yaoOracle.Start(ctx)
	if err != nil {
		log.Fatal("Failed to start YaoOracle:", err)
	}

	// Example 1: Account operations
	log.Println("=== Account Operations ===")

	// Create test account
	testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	testAccount := &types.Account{
		Address:     testAddr,
		Nonce:       1,
		Balance:     big.NewInt(1000000000000000000), // 1 ETH
		CodeHash:    common.Hash{},
		StorageRoot: common.Hash{},
	}

	// Set account (will cache in L1 and L2)
	err = yaoOracle.SetAccount(ctx, testAddr, testAccount)
	if err != nil {
		log.Printf("Failed to set account: %v", err)
	} else {
		log.Printf("Set account %s with balance %s", testAddr.Hex(), testAccount.Balance.String())
	}

	// Get account (should hit L1 cache for fast access)
	retrievedAccount, err := yaoOracle.GetAccount(ctx, testAddr)
	if err != nil {
		log.Printf("Failed to get account: %v", err)
	} else if retrievedAccount != nil {
		log.Printf("Retrieved account %s with nonce %d and balance %s",
			retrievedAccount.Address.Hex(), retrievedAccount.Nonce, retrievedAccount.Balance.String())
	}

	// Example 2: Contract storage operations
	log.Println("\n=== Contract Storage Operations ===")

	// Contract address and storage key
	contractAddr := common.HexToAddress("0xabcdefabcdefabcdefabcdefabcdefabcdefabcd")
	storageKey := common.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	storageValue := common.HexToHash("0x2222222222222222222222222222222222222222222222222222222222222222")

	// Set storage value
	err = yaoOracle.SetStorageAt(ctx, contractAddr, storageKey, storageValue)
	if err != nil {
		log.Printf("Failed to set storage: %v", err)
	} else {
		log.Printf("Set storage for contract %s at key %s with value %s",
			contractAddr.Hex(), storageKey.Hex(), storageValue.Hex())
	}

	// Get storage value (should hit cache)
	retrievedValue, err := yaoOracle.GetStorageAt(ctx, contractAddr, storageKey)
	if err != nil {
		log.Printf("Failed to get storage: %v", err)
	} else {
		log.Printf("Retrieved storage value %s from contract %s at key %s",
			retrievedValue.Hex(), contractAddr.Hex(), storageKey.Hex())
	}

	// Example 3: Contract code operations
	log.Println("\n=== Contract Code Operations ===")

	// Sample contract code (ERC20 token bytecode snippet)
	contractCode := common.Hex2Bytes("608060405234801561001057600080fd5b50600436106100365760003560e01c8063a9059cbb1461003b578063dd62ed3e14610057575b600080fd5b6100556004803603810190610050919061007d565b610073565b005b61005f610075565b60405161006c91906100bd565b60405180910390f35b5050565b60008054905090565b60008135905061008f816100d8565b92915050565b6000813590506100a4816100ef565b92915050565b6100b3816100d8565b82525050565b60006020820190506100ce60008301846100aa565b92915050565b6100dd816100d8565b81146100e857600080fd5b50565b6100f4816100d8565b81146100ff57600080fd5b5056fea2646970667358221220...")

	// Set contract code
	err = yaoOracle.SetCode(ctx, contractAddr, contractCode)
	if err != nil {
		log.Printf("Failed to set contract code: %v", err)
	} else {
		log.Printf("Set contract code for address %s, code length: %d bytes",
			contractAddr.Hex(), len(contractCode))
	}

	// Get contract code
	retrievedCode, err := yaoOracle.GetCode(ctx, contractAddr)
	if err != nil {
		log.Printf("Failed to get contract code: %v", err)
	} else {
		log.Printf("Retrieved contract code from address %s, code length: %d bytes",
			contractAddr.Hex(), len(retrievedCode))
	}

	// Example 4: Cache invalidation simulation
	log.Println("\n=== Cache Invalidation ===")

	// Simulate cache invalidation message (would normally come from YaoArchive)
	invalidationMsg := &types.CacheInvalidationMessage{
		BlockNumber: 1000,
		Keys: []*types.CacheKey{
			{
				Type:    types.CacheKeyTypeAccount,
				Address: testAddr.Hex(),
			},
			{
				Type:    types.CacheKeyTypeStorage,
				Address: contractAddr.Hex(),
				Key:     storageKey.Hex(),
			},
		},
		Timestamp: time.Now(),
	}

	// Handle cache invalidation
	handler := yaoOracle.GetCacheInvalidationHandler()
	err = handler.HandleCacheInvalidation(invalidationMsg)
	if err != nil {
		log.Printf("Failed to handle cache invalidation: %v", err)
	} else {
		log.Printf("Successfully invalidated %d cache keys", len(invalidationMsg.Keys))
	}

	// Example 5: Performance metrics
	log.Println("\n=== Performance Metrics ===")

	metrics := yaoOracle.GetMetrics()
	log.Printf("Total requests: %d", metrics.TotalRequests)
	log.Printf("Average response time: %d ns", metrics.AverageResponseTime)
	log.Printf("L3 queries: %d", metrics.L3Queries)
	log.Printf("Total errors: %d", metrics.TotalErrors)
	log.Printf("Invalidated keys: %d", metrics.InvalidatedKeys)

	if metrics.L1Stats != nil {
		log.Printf("L1 Cache - Hits: %d, Misses: %d, Hit Ratio: %.2f%%",
			metrics.L1Stats.Hits, metrics.L1Stats.Misses, metrics.L1Stats.HitRatio*100)
	}

	if metrics.L2Stats != nil {
		log.Printf("L2 Cache - Hits: %d, Misses: %d, Hit Ratio: %.2f%%",
			metrics.L2Stats.Hits, metrics.L2Stats.Misses, metrics.L2Stats.HitRatio*100)
	}

	log.Println("\n=== YaoOracle Example Completed Successfully ===")
}
