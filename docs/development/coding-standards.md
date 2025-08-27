# YaoVerse 代码规范

## 概述

本文档定义了YaoVerse项目的代码标准和开发规范，旨在确保代码的一致性、可读性和可维护性。所有团队成员都应遵循这些规范。

## 语言规范

### 基础要求

- **编程语言**: Go 1.21+
- **代码风格**: 严格遵循官方Go代码规范
- **格式化工具**: 使用 `gofmt` 和 `goimports`
- **静态检查**: 使用 `golangci-lint` 进行代码质量检查

### Go 语言特定规范

#### 1. 命名规范

```go
// Package names: 简短、小写、无下划线
package portal

// Interface names: 以 -er 结尾或描述性名词
type StateReader interface {}
type BlockProcessor interface {}

// Struct names: 大写开头，驼峰命名
type ExecutionResult struct {}
type BlockBuilder struct {}

// Method names: 大写开头，动词开头
func (p *Portal) StartServer() error {}
func (s *StateDB) GetAccount() (*Account, error) {}

// Variable names: 驼峰命名，有意义的名称
var blockNumber uint64
var transactionHash common.Hash
var executionResults []*ExecutionResult

// Constants: 大写，下划线分隔
const (
    DefaultPort = 8545
    MaxBlockSize = 1024 * 1024
    TopicTxIngress = "tx-ingress-topic"
)

// Private fields/methods: 小写开头
type server struct {
    config    *Config
    logger    *Logger
    isRunning bool
}

func (s *server) handleRequest() {}
```

#### 2. 错误处理

```go
// 优先返回错误，使用标准错误处理模式
func ProcessTransaction(tx *Transaction) (*Receipt, error) {
    if tx == nil {
        return nil, fmt.Errorf("transaction is nil")
    }
    
    receipt, err := executeTransaction(tx)
    if err != nil {
        return nil, fmt.Errorf("failed to execute transaction: %w", err)
    }
    
    return receipt, nil
}

// 使用 errors.Is() 和 errors.As() 进行错误判断
import "errors"

if errors.Is(err, ErrTransactionNotFound) {
    // handle specific error
}

// 自定义错误类型
type ValidationError struct {
    Field string
    Value interface{}
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Msg)
}
```

#### 3. Context 使用

```go
// 所有可能长时间运行的操作都应接受 context
func (s *Server) Start(ctx context.Context) error {}
func (db *StateDB) GetAccount(ctx context.Context, address common.Address) (*Account, error) {}

// Context 应该是第一个参数
func processBlock(ctx context.Context, block *Block, options *ProcessOptions) error {}

// 检查 context 是否已取消
func longRunningOperation(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // do work
        }
    }
}
```

#### 4. 接口设计

```go
// 接口应该小而精准
type StateReader interface {
    GetAccount(ctx context.Context, address common.Address) (*Account, error)
    GetBalance(ctx context.Context, address common.Address) (*big.Int, error)
    GetNonce(ctx context.Context, address common.Address) (uint64, error)
}

// 避免过大的接口，优先组合
type StateDB interface {
    StateReader
    StateWriter
    StateManager
}

// 接口应该定义在使用它的包中，而不是实现它的包中
```

#### 5. 并发安全

```go
// 使用 sync.Mutex 保护共享状态
type SafeCounter struct {
    mu    sync.RWMutex
    count int64
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Get() int64 {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.count
}

// 优先使用 channel 进行 goroutine 间通信
func worker(jobs <-chan Job, results chan<- Result, done <-chan struct{}) {
    for {
        select {
        case job := <-jobs:
            result := processJob(job)
            results <- result
        case <-done:
            return
        }
    }
}
```

## 项目结构规范

### 1. 目录结构

```
yao-<component>/
├── cmd/                    # 主程序入口
│   └── yao-<component>/
│       └── main.go
├── internal/               # 私有应用代码
│   ├── config/            # 配置管理
│   ├── server/            # 服务器实现
│   ├── handlers/          # 处理器
│   ├── middleware/        # 中间件
│   └── models/            # 数据模型
├── interfaces/            # 公共接口定义
├── pkg/                   # 可被外部使用的库代码
├── scripts/               # 构建、部署脚本
├── deployments/           # 部署配置
├── docs/                  # 组件文档
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 2. 包组织原则

```go
// 按功能而非类型组织包
// 好的例子：
internal/
├── portal/          # 门户相关功能
├── sequencer/       # 序列器功能
├── executor/        # 执行器功能
└── storage/         # 存储功能

// 避免按类型组织：
internal/
├── models/          # 所有模型
├── handlers/        # 所有处理器
├── services/        # 所有服务
└── utils/           # 工具函数
```

### 3. 文件命名

```bash
# 文件名使用小写和下划线
block_processor.go
transaction_pool.go
state_manager.go

# 测试文件以 _test.go 结尾
block_processor_test.go
transaction_pool_test.go

# 示例文件以 _example.go 结尾
state_manager_example.go
```

## 注释和文档规范

### 1. 包级别注释

```go
// Package portal provides Ethereum-compatible JSON-RPC API gateway functionality.
// It handles transaction submission, validation, and routing to the YaoVerse system.
//
// The portal supports:
//   - Standard Ethereum JSON-RPC methods
//   - WebSocket connections for real-time updates
//   - Rate limiting and load balancing
//   - Health checks and metrics
//
// Example usage:
//   config := &PortalConfig{
//       Host: "0.0.0.0",
//       Port: 8545,
//   }
//   server := portal.NewServer(config)
//   if err := server.Start(context.Background()); err != nil {
//       log.Fatal(err)
//   }
package portal
```

### 2. 类型注释

```go
// StateDB represents the blockchain state database interface.
// It provides methods for reading and writing account states, contract storage,
// and managing state snapshots.
//
// StateDB is thread-safe and can be used concurrently from multiple goroutines.
// All methods that may take significant time accept a context.Context parameter
// for cancellation support.
type StateDB interface {
    // GetAccount retrieves an account by address.
    // Returns ErrAccountNotFound if the account doesn't exist.
    GetAccount(ctx context.Context, address common.Address) (*Account, error)
    
    // SetAccount updates or creates an account.
    // The account will be persisted when Commit() is called.
    SetAccount(ctx context.Context, address common.Address, account *Account) error
}

// ExecutionResult contains the result of block execution.
// It includes the updated state, transaction receipts, and execution metrics.
type ExecutionResult struct {
    BlockNumber  uint64      `json:"blockNumber"`  // The block number that was executed
    BlockHash    common.Hash `json:"blockHash"`    // The hash of the executed block
    StateRoot    common.Hash `json:"stateRoot"`    // The resulting state root
    Receipts     []*Receipt  `json:"receipts"`     // Transaction receipts
    StateDiff    *StateDiff  `json:"stateDiff"`    // State changes made during execution
    GasUsed      uint64      `json:"gasUsed"`      // Total gas used by all transactions
    ExecutedAt   time.Time   `json:"executedAt"`   // When the block was executed
    ExecutorNode string      `json:"executorNode"` // Which node executed the block
}
```

### 3. 函数注释

```go
// ProcessBlock executes all transactions in a block and returns the execution result.
// The block must be validated before calling this function.
//
// The function performs the following steps:
//   1. Create a new EVM instance for the block
//   2. Execute each transaction in order
//   3. Collect receipts and calculate state diff
//   4. Return the complete execution result
//
// If any transaction fails to execute, the entire block execution is aborted
// and an error is returned. The state is automatically reverted to the
// pre-execution snapshot.
//
// Example:
//   result, err := processor.ProcessBlock(ctx, block)
//   if err != nil {
//       log.Printf("Block execution failed: %v", err)
//       return err
//   }
//   log.Printf("Block %d executed successfully, gas used: %d", 
//              result.BlockNumber, result.GasUsed)
func (p *BlockProcessor) ProcessBlock(ctx context.Context, block *Block) (*ExecutionResult, error) {
    // Implementation...
}
```

### 4. 内联注释

```go
func processTransaction(tx *Transaction) (*Receipt, error) {
    // Validate transaction signature and nonce
    if err := validateTransaction(tx); err != nil {
        return nil, fmt.Errorf("transaction validation failed: %w", err)
    }
    
    // Create execution context
    execCtx := &ExecutionContext{
        Transaction: tx,
        BlockNumber: getCurrentBlockNumber(),
        Timestamp:   time.Now().Unix(),
    }
    
    // Execute the transaction
    // Note: State changes are applied to a temporary state until commit
    result, err := executeInEVM(execCtx)
    if err != nil {
        // Execution failed - state is automatically reverted
        return nil, fmt.Errorf("EVM execution failed: %w", err)
    }
    
    // Generate receipt with logs and gas usage
    receipt := generateReceipt(tx, result)
    return receipt, nil
}
```

## 日志规范

### 1. 日志级别

```go
import "go.uber.org/zap"

// 使用结构化日志
logger := zap.NewProduction()
defer logger.Sync()

// DEBUG: 详细的调试信息
logger.Debug("Processing transaction", 
    zap.String("txHash", tx.Hash().Hex()),
    zap.Uint64("blockNumber", blockNumber),
    zap.Duration("processingTime", duration))

// INFO: 一般信息
logger.Info("Block execution completed",
    zap.Uint64("blockNumber", block.Number),
    zap.Int("txCount", len(block.Transactions)),
    zap.Uint64("gasUsed", block.GasUsed))

// WARN: 警告信息
logger.Warn("High memory usage detected",
    zap.Uint64("memoryUsage", memUsage),
    zap.Uint64("threshold", threshold))

// ERROR: 错误信息
logger.Error("Failed to connect to database",
    zap.Error(err),
    zap.String("dsn", config.DatabaseDSN),
    zap.Int("retryAttempt", retryCount))

// FATAL: 致命错误（会导致程序退出）
logger.Fatal("Critical system error",
    zap.Error(err),
    zap.String("component", "yao-portal"))
```

### 2. 日志内容规范

```go
// 包含关键上下文信息
logger.Info("Transaction executed successfully",
    zap.String("component", "yao-vulcan"),
    zap.String("txHash", tx.Hash().Hex()),
    zap.String("from", tx.From().Hex()),
    zap.String("to", tx.To().Hex()),
    zap.Uint64("gasUsed", receipt.GasUsed),
    zap.Duration("executionTime", execTime),
    zap.String("traceId", traceId)) // 用于链路追踪

// 错误日志包含堆栈信息
logger.Error("Block processing failed",
    zap.Error(err),
    zap.Stack("stackTrace"),
    zap.Uint64("blockNumber", block.Number),
    zap.String("blockHash", block.Hash().Hex()),
    zap.String("processorNode", nodeID))
```

## 测试规范

### 1. 测试结构

```go
func TestBlockProcessor_ProcessBlock(t *testing.T) {
    tests := []struct {
        name          string
        block         *Block
        expectedError string
        expectedGas   uint64
        setup         func() // 测试前的准备工作
        cleanup       func() // 测试后的清理工作
    }{
        {
            name: "successful_block_processing",
            block: &Block{
                Number:       1,
                Transactions: []*Transaction{validTx1, validTx2},
            },
            expectedError: "",
            expectedGas:   42000,
            setup:         func() { setupMockStateDB() },
            cleanup:       func() { cleanupMockStateDB() },
        },
        {
            name: "invalid_transaction_in_block",
            block: &Block{
                Number:       2,
                Transactions: []*Transaction{invalidTx},
            },
            expectedError: "transaction validation failed",
            expectedGas:   0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.setup != nil {
                tt.setup()
            }
            defer func() {
                if tt.cleanup != nil {
                    tt.cleanup()
                }
            }()

            processor := NewBlockProcessor(mockConfig)
            result, err := processor.ProcessBlock(context.Background(), tt.block)

            if tt.expectedError != "" {
                assert.Error(t, err)
                assert.Contains(t, err.Error(), tt.expectedError)
                return
            }

            require.NoError(t, err)
            assert.Equal(t, tt.expectedGas, result.GasUsed)
        })
    }
}
```

### 2. 测试覆盖率

```go
// 单元测试要求：
// - 所有公共方法必须有测试
// - 代码覆盖率不低于80%
// - 包含正常情况和异常情况的测试

// 使用 testify 进行断言
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)

// Mock 接口用于测试
type MockStateDB struct {
    mock.Mock
}

func (m *MockStateDB) GetAccount(ctx context.Context, address common.Address) (*Account, error) {
    args := m.Called(ctx, address)
    return args.Get(0).(*Account), args.Error(1)
}
```

### 3. 基准测试

```go
func BenchmarkBlockProcessor_ProcessBlock(b *testing.B) {
    processor := NewBlockProcessor(testConfig)
    block := generateTestBlock(100) // 100个交易的测试块
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := processor.ProcessBlock(context.Background(), block)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// 并发基准测试
func BenchmarkBlockProcessor_ProcessBlock_Parallel(b *testing.B) {
    processor := NewBlockProcessor(testConfig)
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            block := generateTestBlock(10)
            _, err := processor.ProcessBlock(context.Background(), block)
            if err != nil {
                b.Error(err)
            }
        }
    })
}
```

## 配置管理规范

### 1. 配置结构

```go
// Config represents the configuration for YaoVulcan
type Config struct {
    // Node configuration
    NodeID   string `json:"nodeId" validate:"required"`
    NodeName string `json:"nodeName" validate:"required"`
    
    // Server configuration
    Host string `json:"host" validate:"required"`
    Port int    `json:"port" validate:"min=1,max=65535"`
    
    // Database configuration
    Database DatabaseConfig `json:"database" validate:"required"`
    
    // Cache configuration
    Cache CacheConfig `json:"cache" validate:"required"`
    
    // Logging configuration
    Logging LoggingConfig `json:"logging" validate:"required"`
}

type DatabaseConfig struct {
    DSN             string `json:"dsn" validate:"required"`
    MaxOpenConn     int    `json:"maxOpenConn" validate:"min=1"`
    MaxIdleConn     int    `json:"maxIdleConn" validate:"min=1"`
    ConnMaxLifetime int    `json:"connMaxLifetime" validate:"min=1"` // in minutes
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
    config := &Config{}
    
    // Load from file
    if err := loadFromFile(configPath, config); err != nil {
        return nil, fmt.Errorf("failed to load config from file: %w", err)
    }
    
    // Override with environment variables
    if err := loadFromEnv(config); err != nil {
        return nil, fmt.Errorf("failed to load config from env: %w", err)
    }
    
    // Validate configuration
    if err := validateConfig(config); err != nil {
        return nil, fmt.Errorf("config validation failed: %w", err)
    }
    
    return config, nil
}
```

### 2. 环境变量映射

```go
// 使用结构体标签定义环境变量映射
type Config struct {
    Host     string `json:"host" env:"YAO_HOST" envDefault:"0.0.0.0"`
    Port     int    `json:"port" env:"YAO_PORT" envDefault:"8545"`
    LogLevel string `json:"logLevel" env:"YAO_LOG_LEVEL" envDefault:"info"`
}

// 支持嵌套配置
type DatabaseConfig struct {
    Host     string `json:"host" env:"YAO_DB_HOST" validate:"required"`
    Port     int    `json:"port" env:"YAO_DB_PORT" envDefault:"3306"`
    Username string `json:"username" env:"YAO_DB_USERNAME" validate:"required"`
    Password string `json:"password" env:"YAO_DB_PASSWORD" validate:"required"`
    Database string `json:"database" env:"YAO_DB_DATABASE" validate:"required"`
}
```

## 错误处理和重试策略

### 1. 错误定义

```go
// 定义项目特定的错误类型
var (
    ErrBlockNotFound       = errors.New("block not found")
    ErrTransactionNotFound = errors.New("transaction not found")
    ErrInvalidSignature    = errors.New("invalid transaction signature")
    ErrInsufficientBalance = errors.New("insufficient balance")
    ErrGasLimitExceeded    = errors.New("gas limit exceeded")
    ErrNonceTooLow         = errors.New("nonce too low")
    ErrNonceTooHigh        = errors.New("nonce too high")
)

// 可包装的错误类型
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
    Cause   error
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
    return e.Cause
}
```

### 2. 重试策略

```go
// 使用指数退避重试
func retryWithBackoff(ctx context.Context, operation func() error, maxRetries int) error {
    backoff := time.Millisecond * 100
    
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        // 检查是否是可重试的错误
        if !isRetryableError(err) {
            return err
        }
        
        if attempt == maxRetries-1 {
            return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, err)
        }
        
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(backoff):
            backoff *= 2 // 指数退避
            if backoff > time.Second*10 {
                backoff = time.Second * 10 // 最大延迟
            }
        }
    }
    
    return nil
}

func isRetryableError(err error) bool {
    // 网络错误、临时错误等可以重试
    var netErr net.Error
    if errors.As(err, &netErr) && netErr.Temporary() {
        return true
    }
    
    // 数据库连接错误可以重试
    if errors.Is(err, sql.ErrConnDone) {
        return true
    }
    
    return false
}
```

## 性能要求

### 1. 响应时间要求

```go
// API 响应时间要求
const (
    MaxAPIResponseTime     = 100 * time.Millisecond  // API最大响应时间
    MaxBlockProcessingTime = 1 * time.Second         // 区块处理最大时间
    MaxStateReadTime       = 10 * time.Millisecond   // 状态读取最大时间
    MaxCacheAccessTime     = 1 * time.Millisecond    // 缓存访问最大时间
)

// 使用 context.WithTimeout 控制操作超时
func (s *StateDB) GetAccount(ctx context.Context, address common.Address) (*Account, error) {
    ctx, cancel := context.WithTimeout(ctx, MaxStateReadTime)
    defer cancel()
    
    // 实际的状态读取逻辑
    return s.storage.GetAccount(ctx, address)
}
```

### 2. 内存使用要求

```go
// 定期检查内存使用情况
func (s *Server) monitorMemoryUsage() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        var m runtime.MemStats
        runtime.GC()
        runtime.ReadMemStats(&m)
        
        memUsageMB := bToMb(m.Alloc)
        
        if memUsageMB > s.config.MaxMemoryUsageMB {
            s.logger.Warn("High memory usage detected",
                zap.Uint64("currentMB", memUsageMB),
                zap.Uint64("limitMB", s.config.MaxMemoryUsageMB))
            
            // 触发内存清理
            s.performMemoryCleanup()
        }
    }
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
```

## 安全规范

### 1. 输入验证

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

// 严格验证所有外部输入
func ValidateTransactionRequest(req *TransactionRequest) error {
    if err := validate.Struct(req); err != nil {
        return fmt.Errorf("request validation failed: %w", err)
    }
    
    // 自定义验证逻辑
    if req.GasLimit > MaxGasLimit {
        return fmt.Errorf("gas limit exceeds maximum: %d", MaxGasLimit)
    }
    
    if req.Value.Cmp(MaxTransactionValue) > 0 {
        return fmt.Errorf("transaction value exceeds maximum")
    }
    
    return nil
}

// 地址验证
func ValidateEthereumAddress(address string) error {
    if !common.IsHexAddress(address) {
        return fmt.Errorf("invalid ethereum address format: %s", address)
    }
    return nil
}
```

### 2. 敏感信息处理

```go
// 配置中的敏感信息不应记录到日志
type Config struct {
    DatabasePassword string `json:"databasePassword" log:"-"`     // 不记录到日志
    APISecretKey     string `json:"apiSecretKey" log:"-"`         // 不记录到日志
    PublicEndpoint   string `json:"publicEndpoint"`               // 可以记录
}

// 日志脱敏
func (c *Config) String() string {
    cfg := *c
    cfg.DatabasePassword = "***"
    cfg.APISecretKey = "***"
    
    data, _ := json.MarshalIndent(cfg, "", "  ")
    return string(data)
}
```

## 代码质量检查

### 1. 使用 golangci-lint

```yaml
# .golangci.yml
run:
  timeout: 5m
  modules-download-mode: readonly

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck
    - exportloopref
    - nolintlint
    - revive
    - gosec
    - unconvert
    - misspell
    - nakedret
    - prealloc
    - scopelint

linters-settings:
  revive:
    rules:
      - name: exported
        arguments: ["checkPrivateReceivers", "sayRepetitiveInsteadOfStutters"]
  
  gosec:
    excludes:
      - G204 # subprocess with variable
      - G304 # file path provided as taint input
      
  nakedret:
    max-func-lines: 30
```

### 2. Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Running pre-commit checks..."

# Format code
gofmt -w .
goimports -w .

# Run linter
golangci-lint run

# Run tests
go test ./...

# Check if there are any changes after formatting
if ! git diff --exit-code; then
    echo "Code was reformatted. Please stage the changes and commit again."
    exit 1
fi

echo "All checks passed!"
```

## Makefile 规范

```makefile
.PHONY: build test lint clean docker-build docker-push help

# Build variables
VERSION ?= latest
REGISTRY ?= eggybyte
IMAGE_NAME ?= yao-verse
GO_VERSION = 1.21

# Default target
.DEFAULT_GOAL := help

## Build the binary
build:
	go build -o bin/yao-portal ./cmd/yao-portal
	go build -o bin/yao-chronos ./cmd/yao-chronos
	go build -o bin/yao-vulcan ./cmd/yao-vulcan
	go build -o bin/yao-archive ./cmd/yao-archive

## Run tests
test:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## Run linter
lint:
	golangci-lint run

## Format code
fmt:
	gofmt -w .
	goimports -w .

## Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

## Build Docker images
docker-build:
	docker build -t $(REGISTRY)/$(IMAGE_NAME)-portal:$(VERSION) -f deployments/portal/Dockerfile .
	docker build -t $(REGISTRY)/$(IMAGE_NAME)-chronos:$(VERSION) -f deployments/chronos/Dockerfile .
	docker build -t $(REGISTRY)/$(IMAGE_NAME)-vulcan:$(VERSION) -f deployments/vulcan/Dockerfile .
	docker build -t $(REGISTRY)/$(IMAGE_NAME)-archive:$(VERSION) -f deployments/archive/Dockerfile .

## Push Docker images
docker-push: docker-build
	docker push $(REGISTRY)/$(IMAGE_NAME)-portal:$(VERSION)
	docker push $(REGISTRY)/$(IMAGE_NAME)-chronos:$(VERSION)
	docker push $(REGISTRY)/$(IMAGE_NAME)-vulcan:$(VERSION)
	docker push $(REGISTRY)/$(IMAGE_NAME)-archive:$(VERSION)

## Show help
help:
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) printf "  %-20s%s\n", $$1, $$2 \
	}' $(MAKEFILE_LIST)
```

## 版本控制规范

### 1. 提交信息格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

类型 (type):
- `feat`: 新功能
- `fix`: 错误修复
- `docs`: 文档变更
- `style`: 代码格式变更
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例:
```
feat(portal): add WebSocket support for real-time updates

- Implement WebSocket server for transaction status updates
- Add subscription management for multiple clients
- Include connection health monitoring

Closes #123
```

### 2. 分支管理

- `main`: 主分支，包含稳定版本
- `develop`: 开发分支，集成最新功能
- `feature/xxx`: 功能分支
- `fix/xxx`: 修复分支
- `release/xxx`: 发布分支

---

遵循这些代码规范将帮助我们构建高质量、易维护的YaoVerse系统。所有团队成员都应该熟悉并严格遵循这些标准。 