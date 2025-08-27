# YaoVerse 测试指南

## 概述

本文档描述了YaoVerse项目的测试策略、测试类型和最佳实践。我们采用多层次的测试方法来确保系统的可靠性和性能。

## 测试金字塔

YaoVerse遵循测试金字塔原则，包含以下层次：

```
        /\
       /  \
      / E2E\     <- 端到端测试 (少量)
     /      \
    /--------\
   / 集成测试  \   <- 集成测试 (适量)
  /           \
 /-------------\
/  单元测试      \  <- 单元测试 (大量)
\______________/
```

## 测试环境搭建

### 本地测试环境

```bash
# 安装测试依赖
go mod download
go install github.com/golang/mock/mockgen@latest
go install github.com/onsi/ginkgo/v2/ginkgo@latest

# 启动测试依赖服务（使用Docker Compose）
docker-compose -f docker-compose.test.yml up -d

# 等待服务启动
sleep 30

# 运行所有测试
make test-all
```

### CI/CD 测试环境

```yaml
# .github/workflows/test.yml
name: Test Suite

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run unit tests
      run: |
        make test-unit
        make test-coverage

  integration-tests:
    runs-on: ubuntu-latest
    services:
      redis:
        image: redis:7-alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      tidb:
        image: pingcap/tidb:latest
        ports:
          - 4000:4000
          
    steps:
    - uses: actions/checkout@v3
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run integration tests
      run: make test-integration
```

## 单元测试

### 测试结构

单元测试遵循以下目录结构：

```
pkg/
├── portal/
│   ├── server.go
│   ├── server_test.go
│   ├── handler.go
│   ├── handler_test.go
│   └── testdata/
│       ├── valid_transaction.json
│       └── invalid_signature.json
├── sequencer/
│   ├── chronos.go
│   ├── chronos_test.go
│   └── mocks/
│       ├── mock_storage.go
│       └── mock_message_queue.go
```

### 测试规范

#### 1. 命名约定

```go
// 测试函数命名格式: TestMethodName_Scenario_ExpectedBehavior
func TestBlockProcessor_ProcessValidBlock_ReturnsSuccess(t *testing.T) {}
func TestBlockProcessor_ProcessInvalidBlock_ReturnsError(t *testing.T) {}
func TestTransactionPool_AddDuplicateTransaction_ReturnsError(t *testing.T) {}

// 基准测试命名格式: BenchmarkMethodName_Scenario
func BenchmarkBlockProcessor_ProcessBlock_LargeBlock(b *testing.B) {}
func BenchmarkStateDB_GetAccount_CacheHit(b *testing.B) {}
```

#### 2. 表格驱动测试

```go
func TestTransactionValidator_ValidateTransaction(t *testing.T) {
    tests := []struct {
        name        string
        transaction *types.Transaction
        wantErr     bool
        expectedErr error
        setup       func(*testing.T) *TransactionValidator
        cleanup     func(*testing.T)
    }{
        {
            name: "valid_transaction",
            transaction: &types.Transaction{
                From:     common.HexToAddress("0x1234..."),
                To:       common.HexToAddress("0x5678..."),
                Value:    big.NewInt(1000),
                Gas:      21000,
                GasPrice: big.NewInt(1e9),
                Nonce:    1,
                Data:     []byte{},
            },
            wantErr: false,
            setup: func(t *testing.T) *TransactionValidator {
                return NewTransactionValidator(testConfig)
            },
        },
        {
            name: "insufficient_balance",
            transaction: &types.Transaction{
                From:     common.HexToAddress("0x1234..."),
                To:       common.HexToAddress("0x5678..."),
                Value:    big.NewInt(999999999), // 余额不足
                Gas:      21000,
                GasPrice: big.NewInt(1e9),
                Nonce:    1,
            },
            wantErr:     true,
            expectedErr: ErrInsufficientBalance,
            setup: func(t *testing.T) *TransactionValidator {
                return NewTransactionValidator(testConfig)
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            validator := tt.setup(t)
            if tt.cleanup != nil {
                defer tt.cleanup(t)
            }

            err := validator.ValidateTransaction(context.Background(), tt.transaction)

            if tt.wantErr {
                assert.Error(t, err)
                if tt.expectedErr != nil {
                    assert.ErrorIs(t, err, tt.expectedErr)
                }
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

#### 3. Mock 使用

```go
//go:generate mockgen -source=interfaces/state_db.go -destination=mocks/mock_state_db.go

func TestBlockProcessor_ProcessBlock_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // 创建 mock 对象
    mockStateDB := mocks.NewMockStateDB(ctrl)
    mockOracle := mocks.NewMockYaoOracle(ctrl)

    // 设置期望的调用
    mockStateDB.EXPECT().
        GetAccount(gomock.Any(), gomock.Any()).
        Return(&types.Account{Balance: big.NewInt(1000)}, nil).
        Times(2)

    mockStateDB.EXPECT().
        SetAccount(gomock.Any(), gomock.Any(), gomock.Any()).
        Return(nil).
        Times(1)

    mockOracle.EXPECT().
        GetStateDiff(gomock.Any(), gomock.Any()).
        Return(&types.StateDiff{}, nil)

    // 创建处理器实例
    processor := &BlockProcessor{
        stateDB: mockStateDB,
        oracle:  mockOracle,
    }

    // 执行测试
    block := createTestBlock()
    result, err := processor.ProcessBlock(context.Background(), block)

    // 断言结果
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, block.Number, result.BlockNumber)
}
```

#### 4. 测试辅助函数

```go
// testutil/helpers.go
package testutil

import (
    "math/big"
    "testing"
    
    "github.com/ethereum/go-ethereum/common"
    "github.com/eggybyte-technology/yao-verse-shared/types"
)

// CreateTestTransaction 创建测试交易
func CreateTestTransaction(t *testing.T, opts ...TransactionOption) *types.Transaction {
    t.Helper()
    
    tx := &types.Transaction{
        From:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
        To:       common.HexToAddress("0x0987654321098765432109876543210987654321"),
        Value:    big.NewInt(1000),
        Gas:      21000,
        GasPrice: big.NewInt(1e9),
        Nonce:    1,
        Data:     []byte{},
    }
    
    for _, opt := range opts {
        opt(tx)
    }
    
    return tx
}

// TransactionOption 交易选项函数
type TransactionOption func(*types.Transaction)

func WithValue(value *big.Int) TransactionOption {
    return func(tx *types.Transaction) {
        tx.Value = value
    }
}

func WithGasPrice(gasPrice *big.Int) TransactionOption {
    return func(tx *types.Transaction) {
        tx.GasPrice = gasPrice
    }
}

func WithNonce(nonce uint64) TransactionOption {
    return func(tx *types.Transaction) {
        tx.Nonce = nonce
    }
}

// CreateTestBlock 创建测试区块
func CreateTestBlock(t *testing.T, txCount int) *types.Block {
    t.Helper()
    
    var transactions []*types.Transaction
    for i := 0; i < txCount; i++ {
        tx := CreateTestTransaction(t, WithNonce(uint64(i)))
        transactions = append(transactions, tx)
    }
    
    return &types.Block{
        Number:       1,
        Hash:         common.HexToHash("0xabc123"),
        ParentHash:   common.HexToHash("0x000000"),
        Timestamp:    uint64(time.Now().Unix()),
        GasLimit:     8000000,
        Transactions: transactions,
    }
}
```

### 覆盖率要求

- **总体覆盖率**: 最低80%
- **关键组件覆盖率**: 最低90%
  - YaoOracle状态库
  - 交易验证逻辑
  - 区块处理逻辑
  - 缓存管理逻辑

```bash
# 生成覆盖率报告
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 检查覆盖率
go tool cover -func=coverage.out | grep total
```

## 集成测试

### 测试分类

#### 1. 组件间集成测试

测试不同组件之间的交互：

```go
// test/integration/portal_chronos_test.go
func TestPortalChronosIntegration(t *testing.T) {
    // 启动测试环境
    env := setupIntegrationEnvironment(t)
    defer env.Cleanup()

    // 启动 Portal 和 Chronos
    portal := startPortalServer(t, env.Config.Portal)
    chronos := startChronosServer(t, env.Config.Chronos)
    
    defer portal.Stop()
    defer chronos.Stop()

    // 等待服务就绪
    waitForHealthy(t, portal, chronos)

    // 测试交易提交流程
    t.Run("submit_transaction_flow", func(t *testing.T) {
        client := NewJSONRPCClient(portal.Address())
        
        // 提交交易
        txHash, err := client.SendRawTransaction(validTxData)
        require.NoError(t, err)
        require.NotEmpty(t, txHash)

        // 验证 Chronos 收到交易
        eventually(t, func() bool {
            status := chronos.GetSequencerStatus()
            return status.PendingTransactions > 0
        }, 5*time.Second)

        // 验证区块生成
        eventually(t, func() bool {
            blockNumber := client.BlockNumber()
            return blockNumber > 0
        }, 10*time.Second)
    })
}
```

#### 2. 数据库集成测试

```go
// test/integration/database_test.go  
func TestDatabaseIntegration(t *testing.T) {
    db := setupTestTiDB(t)
    defer db.Close()

    storage := NewL3Storage(db)
    
    t.Run("account_operations", func(t *testing.T) {
        ctx := context.Background()
        address := common.HexToAddress("0x1234...")
        
        // 测试账户不存在
        exists, err := storage.HasAccount(ctx, address)
        assert.NoError(t, err)
        assert.False(t, exists)

        // 创建账户
        account := &types.Account{
            Address: address,
            Balance: big.NewInt(1000),
            Nonce:   0,
        }
        
        err = storage.SetAccount(ctx, address, account)
        assert.NoError(t, err)

        // 验证账户创建成功
        retrieved, err := storage.GetAccount(ctx, address)
        assert.NoError(t, err)
        assert.Equal(t, account.Balance, retrieved.Balance)
    })
}
```

#### 3. 消息队列集成测试

```go
// test/integration/message_queue_test.go
func TestMessageQueueIntegration(t *testing.T) {
    mq := setupTestRocketMQ(t)
    defer mq.Cleanup()

    producer := NewMessageProducer(mq.Config)
    consumer := NewMessageConsumer(mq.Config)
    
    defer producer.Close()
    defer consumer.Close()

    t.Run("transaction_message_flow", func(t *testing.T) {
        // 创建消息处理器
        var receivedMessages []*types.TxSubmissionMessage
        handler := &testMessageHandler{
            onTxSubmission: func(msg *types.TxSubmissionMessage) error {
                receivedMessages = append(receivedMessages, msg)
                return nil
            },
        }

        // 启动消费者
        err := consumer.Subscribe(types.TopicTxIngress, types.ConsumerGroupChronos, handler)
        require.NoError(t, err)
        
        go consumer.Start()

        // 发送测试消息
        tx := createTestTransaction(t)
        msg := &types.TxSubmissionMessage{
            Type:        types.MessageTypeTxSubmission,
            Transaction: tx,
            SubmittedAt: time.Now(),
            MessageID:   generateMessageID(),
        }

        err = producer.SendTxSubmission(msg)
        require.NoError(t, err)

        // 验证消息被接收
        eventually(t, func() bool {
            return len(receivedMessages) == 1
        }, 5*time.Second)
        
        assert.Equal(t, msg.MessageID, receivedMessages[0].MessageID)
    })
}
```

### 测试环境管理

```go
// test/integration/environment.go
type IntegrationEnvironment struct {
    TiDB     *sql.DB
    Redis    *redis.Client
    RocketMQ *rocketmq.Client
    Config   *IntegrationConfig
    cleanup  []func()
}

func SetupIntegrationEnvironment(t *testing.T) *IntegrationEnvironment {
    t.Helper()
    
    env := &IntegrationEnvironment{}
    
    // 启动 TiDB
    tidbContainer, tidbDB := startTiDBContainer(t)
    env.TiDB = tidbDB
    env.cleanup = append(env.cleanup, func() {
        tidbDB.Close()
        tidbContainer.Terminate(context.Background())
    })
    
    // 启动 Redis
    redisContainer, redisClient := startRedisContainer(t)
    env.Redis = redisClient
    env.cleanup = append(env.cleanup, func() {
        redisClient.Close()
        redisContainer.Terminate(context.Background())
    })
    
    // 初始化数据库模式
    err := initializeSchema(env.TiDB)
    require.NoError(t, err)
    
    return env
}

func (e *IntegrationEnvironment) Cleanup() {
    for i := len(e.cleanup) - 1; i >= 0; i-- {
        e.cleanup[i]()
    }
}
```

## 性能测试

### 基准测试

#### 1. 单组件性能测试

```go
// benchmark/portal_benchmark_test.go
func BenchmarkPortalServer_HandleJSONRPC(b *testing.B) {
    server := setupBenchmarkPortalServer(b)
    defer server.Stop()

    client := setupJSONRPCClient(b, server.Address())

    b.Run("eth_getBalance", func(b *testing.B) {
        address := common.HexToAddress("0x1234...")
        
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                _, err := client.GetBalance(address, "latest")
                if err != nil {
                    b.Fatal(err)
                }
            }
        })
    })

    b.Run("eth_sendRawTransaction", func(b *testing.B) {
        txData := generateValidTransactionData(b)
        
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _, err := client.SendRawTransaction(txData)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}

// benchmark/oracle_benchmark_test.go
func BenchmarkYaoOracle_StateAccess(b *testing.B) {
    oracle := setupBenchmarkOracle(b)
    defer oracle.Close()

    address := common.HexToAddress("0x1234...")
    
    // 预热缓存
    oracle.GetAccount(context.Background(), address)

    b.Run("L1_cache_hit", func(b *testing.B) {
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                _, err := oracle.GetAccount(context.Background(), address)
                if err != nil {
                    b.Fatal(err)
                }
            }
        })
    })

    b.Run("L2_cache_hit", func(b *testing.B) {
        // 清空 L1 缓存
        oracle.ClearL1()
        
        b.ResetTimer()
        b.RunParallel(func(pb *testing.PB) {
            for pb.Next() {
                _, err := oracle.GetAccount(context.Background(), address)
                if err != nil {
                    b.Fatal(err)
                }
            }
        })
    })
}
```

#### 2. 端到端性能测试

```go
// benchmark/e2e_benchmark_test.go
func BenchmarkE2E_TransactionThroughput(b *testing.B) {
    // 启动完整的 YaoVerse 集群
    cluster := startBenchmarkCluster(b)
    defer cluster.Stop()

    client := NewJSONRPCClient(cluster.Portal.Address())

    b.Run("sequential_transactions", func(b *testing.B) {
        txData := generateTransactionBatch(b, b.N)
        
        b.ResetTimer()
        start := time.Now()
        
        for i := 0; i < b.N; i++ {
            _, err := client.SendRawTransaction(txData[i])
            if err != nil {
                b.Fatal(err)
            }
        }
        
        duration := time.Since(start)
        tps := float64(b.N) / duration.Seconds()
        
        b.ReportMetric(tps, "tps")
        b.ReportMetric(duration.Seconds()*1000/float64(b.N), "ms/tx")
    })

    b.Run("concurrent_transactions", func(b *testing.B) {
        const concurrency = 100
        txData := generateTransactionBatch(b, b.N)
        
        semaphore := make(chan struct{}, concurrency)
        var wg sync.WaitGroup
        
        b.ResetTimer()
        start := time.Now()
        
        for i := 0; i < b.N; i++ {
            wg.Add(1)
            go func(tx []byte) {
                defer wg.Done()
                semaphore <- struct{}{}
                defer func() { <-semaphore }()
                
                _, err := client.SendRawTransaction(tx)
                if err != nil {
                    b.Error(err)
                }
            }(txData[i])
        }
        
        wg.Wait()
        duration := time.Since(start)
        tps := float64(b.N) / duration.Seconds()
        
        b.ReportMetric(tps, "tps")
        b.ReportMetric(duration.Seconds()*1000/float64(b.N), "ms/tx")
    })
}
```

### 负载测试

使用专门的负载测试工具：

```go
// test/load/load_test.go
package load

import (
    "context"
    "sync"
    "time"
    
    vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestPortalLoadTest(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping load test in short mode")
    }

    server := startPortalServer(t)
    defer server.Stop()

    // 配置负载测试参数
    rate := vegeta.Rate{Freq: 100, Per: time.Second} // 100 RPS
    duration := 30 * time.Second

    // 创建目标
    targets := []vegeta.Target{
        {
            Method: "POST",
            URL:    server.Address(),
            Header: http.Header{
                "Content-Type": []string{"application/json"},
            },
            Body: []byte(`{
                "jsonrpc": "2.0",
                "method": "eth_blockNumber",
                "params": [],
                "id": 1
            }`),
        },
    }

    targeter := vegeta.NewStaticTargeter(targets...)
    attacker := vegeta.NewAttacker()

    // 执行攻击
    var results vegeta.Results
    for res := range attacker.Attack(targeter, rate, duration, "load-test") {
        results = append(results, res)
    }

    // 分析结果
    metrics := vegeta.NewMetrics(results)
    
    // 断言性能要求
    assert.True(t, metrics.Success > 0.99, "Success rate should be > 99%%")
    assert.True(t, metrics.Latencies.P99 < 100*time.Millisecond, "P99 latency should be < 100ms")
    assert.True(t, metrics.Latencies.Mean < 50*time.Millisecond, "Mean latency should be < 50ms")

    // 输出详细报告
    t.Logf("Load Test Results:")
    t.Logf("Requests: %d", metrics.Requests)
    t.Logf("Success Rate: %.2f%%", metrics.Success*100)
    t.Logf("Mean Latency: %v", metrics.Latencies.Mean)
    t.Logf("P99 Latency: %v", metrics.Latencies.P99)
    t.Logf("Throughput: %.2f req/sec", metrics.Rate)
}
```

## 端到端测试

### 测试场景

```go
// test/e2e/transaction_flow_test.go
func TestE2E_CompleteTransactionFlow(t *testing.T) {
    // 启动完整的 YaoVerse 集群
    cluster := startE2ECluster(t)
    defer cluster.Stop()

    client := NewJSONRPCClient(cluster.Portal.Address())

    t.Run("complete_transaction_lifecycle", func(t *testing.T) {
        // 1. 获取初始状态
        initialBlock, err := client.BlockNumber()
        require.NoError(t, err)

        senderAddr := common.HexToAddress("0x1234...")
        receiverAddr := common.HexToAddress("0x5678...")

        senderBalance, err := client.GetBalance(senderAddr, "latest")
        require.NoError(t, err)
        
        receiverBalance, err := client.GetBalance(receiverAddr, "latest")
        require.NoError(t, err)

        // 2. 提交交易
        tx := createSignedTransaction(t, senderAddr, receiverAddr, big.NewInt(1000))
        txHash, err := client.SendRawTransaction(tx)
        require.NoError(t, err)
        require.NotEmpty(t, txHash)

        // 3. 等待交易被处理
        var receipt *types.Receipt
        eventually(t, func() bool {
            receipt, err = client.GetTransactionReceipt(txHash)
            return err == nil && receipt != nil
        }, 30*time.Second)

        // 4. 验证交易成功
        assert.Equal(t, uint64(1), receipt.Status)
        assert.True(t, receipt.BlockNumber > initialBlock)

        // 5. 验证状态变更
        newSenderBalance, err := client.GetBalance(senderAddr, "latest")
        require.NoError(t, err)
        
        newReceiverBalance, err := client.GetBalance(receiverAddr, "latest")
        require.NoError(t, err)

        // 计算预期余额（考虑gas费用）
        gasUsed := big.NewInt(int64(receipt.GasUsed))
        gasPrice := big.NewInt(1e9) // 假设的 gas price
        gasFee := new(big.Int).Mul(gasUsed, gasPrice)
        
        expectedSenderBalance := new(big.Int).Sub(senderBalance, big.NewInt(1000))
        expectedSenderBalance = new(big.Int).Sub(expectedSenderBalance, gasFee)
        expectedReceiverBalance := new(big.Int).Add(receiverBalance, big.NewInt(1000))

        assert.Equal(t, expectedSenderBalance, newSenderBalance)
        assert.Equal(t, expectedReceiverBalance, newReceiverBalance)

        // 6. 验证区块包含交易
        block, err := client.GetBlockByNumber(receipt.BlockNumber, true)
        require.NoError(t, err)
        
        found := false
        for _, blockTx := range block.Transactions {
            if blockTx.Hash == txHash {
                found = true
                break
            }
        }
        assert.True(t, found, "Transaction should be included in block")
    })
}
```

### 故障测试

```go
// test/e2e/failure_test.go
func TestE2E_FailureRecovery(t *testing.T) {
    cluster := startE2ECluster(t)
    defer cluster.Stop()

    t.Run("chronos_leader_failover", func(t *testing.T) {
        client := NewJSONRPCClient(cluster.Portal.Address())

        // 获取当前 leader
        leaderStatus := cluster.Chronos.GetElectionStatus()
        require.True(t, leaderStatus.IsLeader)

        // 提交交易以确保系统正常工作
        tx1 := createTestTransaction(t)
        hash1, err := client.SendRawTransaction(tx1)
        require.NoError(t, err)

        // 模拟 leader 故障
        cluster.Chronos.SimulateFailure()
        time.Sleep(5 * time.Second) // 等待故障检测

        // 启动新的 Chronos 节点
        newChronos := cluster.StartNewChronos()
        defer newChronos.Stop()

        // 等待新 leader 选出
        eventually(t, func() bool {
            status := newChronos.GetElectionStatus()
            return status.IsLeader
        }, 30*time.Second)

        // 验证系统继续工作
        tx2 := createTestTransaction(t)
        hash2, err := client.SendRawTransaction(tx2)
        require.NoError(t, err)
        require.NotEqual(t, hash1, hash2)

        // 验证两笔交易都被处理
        eventually(t, func() bool {
            receipt1, _ := client.GetTransactionReceipt(hash1)
            receipt2, _ := client.GetTransactionReceipt(hash2)
            return receipt1 != nil && receipt2 != nil
        }, 30*time.Second)
    })

    t.Run("database_connection_loss", func(t *testing.T) {
        client := NewJSONRPCClient(cluster.Portal.Address())

        // 断开数据库连接
        cluster.TiDB.SimulateNetworkPartition()
        
        // 提交交易（应该失败或延迟）
        tx := createTestTransaction(t)
        _, err := client.SendRawTransaction(tx)
        
        // 这里可能成功（如果有缓冲）或失败
        if err != nil {
            t.Logf("Expected failure due to DB partition: %v", err)
        }

        // 恢复数据库连接
        cluster.TiDB.RestoreNetworkConnection()
        time.Sleep(10 * time.Second) // 等待连接恢复

        // 验证系统恢复正常
        tx2 := createTestTransaction(t)
        hash2, err := client.SendRawTransaction(tx2)
        require.NoError(t, err)

        eventually(t, func() bool {
            receipt, _ := client.GetTransactionReceipt(hash2)
            return receipt != nil && receipt.Status == 1
        }, 30*time.Second)
    })
}
```

## 测试工具和辅助函数

### 测试集群管理

```go
// test/cluster/cluster.go
type TestCluster struct {
    Portal   *PortalServer
    Chronos  *ChronosServer
    Vulcan   []*VulcanServer
    Archive  []*ArchiveServer
    TiDB     *TestTiDB
    Redis    *TestRedis
    RocketMQ *TestRocketMQ
    Config   *ClusterConfig
}

func StartE2ECluster(t *testing.T) *TestCluster {
    t.Helper()

    cluster := &TestCluster{}
    
    // 启动基础设施
    cluster.TiDB = StartTestTiDB(t)
    cluster.Redis = StartTestRedis(t)
    cluster.RocketMQ = StartTestRocketMQ(t)

    // 创建配置
    cluster.Config = generateClusterConfig(cluster)

    // 启动服务
    cluster.Chronos = StartChronosServer(t, cluster.Config.Chronos)
    
    // 启动多个 Vulcan 节点
    for i := 0; i < 3; i++ {
        vulcan := StartVulcanServer(t, cluster.Config.Vulcan[i])
        cluster.Vulcan = append(cluster.Vulcan, vulcan)
    }

    // 启动多个 Archive 节点
    for i := 0; i < 2; i++ {
        archive := StartArchiveServer(t, cluster.Config.Archive[i])
        cluster.Archive = append(cluster.Archive, archive)
    }

    // 启动 Portal
    cluster.Portal = StartPortalServer(t, cluster.Config.Portal)

    // 等待所有服务就绪
    cluster.WaitForHealthy(t, 60*time.Second)

    return cluster
}

func (c *TestCluster) Stop() {
    if c.Portal != nil {
        c.Portal.Stop()
    }
    
    for _, vulcan := range c.Vulcan {
        vulcan.Stop()
    }
    
    for _, archive := range c.Archive {
        archive.Stop()
    }
    
    if c.Chronos != nil {
        c.Chronos.Stop()
    }
    
    if c.RocketMQ != nil {
        c.RocketMQ.Stop()
    }
    
    if c.Redis != nil {
        c.Redis.Stop()
    }
    
    if c.TiDB != nil {
        c.TiDB.Stop()
    }
}
```

### 断言辅助函数

```go
// test/assertions/assertions.go
package assertions

// Eventually 重试断言直到成功或超时
func Eventually(t *testing.T, condition func() bool, timeout time.Duration, msgAndArgs ...interface{}) {
    t.Helper()
    
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    timer := time.NewTimer(timeout)
    defer timer.Stop()
    
    for {
        select {
        case <-timer.C:
            assert.True(t, condition(), msgAndArgs...)
            return
        case <-ticker.C:
            if condition() {
                return
            }
        }
    }
}

// WithinDuration 断言操作在指定时间内完成
func WithinDuration(t *testing.T, expected time.Duration, operation func()) {
    t.Helper()
    
    start := time.Now()
    operation()
    duration := time.Since(start)
    
    assert.True(t, duration <= expected, 
        "Operation took %v, expected <= %v", duration, expected)
}

// TransactionInBlock 断言交易包含在指定区块中
func TransactionInBlock(t *testing.T, client JSONRPCClient, txHash common.Hash, blockNumber uint64) {
    t.Helper()
    
    block, err := client.GetBlockByNumber(blockNumber, true)
    require.NoError(t, err, "Failed to get block %d", blockNumber)
    
    found := false
    for _, tx := range block.Transactions {
        if tx.Hash == txHash {
            found = true
            break
        }
    }
    
    assert.True(t, found, "Transaction %s not found in block %d", txHash.Hex(), blockNumber)
}
```

## 持续集成配置

### Makefile 测试目标

```makefile
# Makefile
.PHONY: test test-unit test-integration test-e2e test-benchmark test-coverage

# 运行所有测试
test: test-unit test-integration

# 单元测试
test-unit:
	go test -v -race -short ./...

# 集成测试
test-integration:
	go test -v -race -tags=integration ./test/integration/...

# 端到端测试
test-e2e:
	go test -v -timeout=10m -tags=e2e ./test/e2e/...

# 基准测试
test-benchmark:
	go test -v -bench=. -benchmem ./...

# 负载测试
test-load:
	go test -v -timeout=5m -tags=load ./test/load/...

# 覆盖率测试
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out | grep total

# 清理测试数据
test-clean:
	docker-compose -f docker-compose.test.yml down -v
	rm -f coverage.out coverage.html

# 测试环境准备
test-setup:
	docker-compose -f docker-compose.test.yml up -d
	sleep 30  # 等待服务启动
	go run ./scripts/setup-test-data.go
```

### Docker Compose 测试配置

```yaml
# docker-compose.test.yml
version: '3.8'

services:
  tidb-test:
    image: pingcap/tidb:v7.5.0
    ports:
      - "4001:4000"
    environment:
      - TIDB_MODE=standalone
    command: 
      - --store=unistore
      - --path=""
      - --log-level=error

  redis-test:
    image: redis:7-alpine
    ports:
      - "6380:6379"
    command: redis-server --appendonly yes

  rocketmq-namesrv-test:
    image: apache/rocketmq:5.1.4
    ports:
      - "9877:9876"
    command: mqnamesrv

  rocketmq-broker-test:
    image: apache/rocketmq:5.1.4
    ports:
      - "10912:10911"
      - "10910:10909"
    depends_on:
      - rocketmq-namesrv-test
    environment:
      - NAMESRV_ADDR=rocketmq-namesrv-test:9876
    command: mqbroker -n rocketmq-namesrv-test:9876 -c /home/rocketmq/rocketmq-5.1.4/conf/broker.conf

networks:
  default:
    name: yaoverse-test
```

## 性能基准

### 目标性能指标

| 组件 | 指标 | 目标值 | 测试方法 |
|------|------|--------|----------|
| YaoPortal | JSON-RPC响应时间 | P99 < 50ms | 负载测试 |
| YaoPortal | 并发连接数 | > 1000 | 连接测试 |
| YaoChronos | 区块生产间隔 | < 2s | 基准测试 |
| YaoVulcan | 交易处理速度 | > 1000 TPS | 性能测试 |
| YaoOracle | L1缓存命中率 | > 95% | 缓存测试 |
| YaoOracle | 状态读取延迟 | L1 < 10ns, L2 < 1ms | 基准测试 |

### 性能回归检测

```bash
#!/bin/bash
# scripts/performance-check.sh

set -e

echo "Running performance regression tests..."

# 运行基准测试并保存结果
go test -bench=. -benchmem -count=5 ./... > benchmark-new.txt

# 如果存在之前的基准测试结果，进行比较
if [ -f benchmark-baseline.txt ]; then
    echo "Comparing with baseline..."
    benchcmp benchmark-baseline.txt benchmark-new.txt || {
        echo "Performance regression detected!"
        exit 1
    }
else
    echo "No baseline found, saving current results as baseline"
    cp benchmark-new.txt benchmark-baseline.txt
fi

echo "Performance check passed!"
```

---

遵循本测试指南将确保YaoVerse系统的高质量和可靠性。定期运行所有类型的测试，并根据结果持续改进代码质量。 