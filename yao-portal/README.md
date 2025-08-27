# YaoPortal - API网关服务

YaoPortal (曜阙) 是 YaoVerse 区块链系统的 API 网关服务，提供以太坊兼容的 JSON-RPC 接口。

## 功能特性

### 🚀 核心功能
- **以太坊兼容**: 完整支持以太坊 JSON-RPC 协议，包括所有标准方法
- **高性能**: 支持 HTTP 和 WebSocket 连接，支持批量请求处理
- **交易池管理**: 内置交易池用于交易缓存和验证
- **限流保护**: 基于客户端 IP 的智能限流机制
- **后端服务管理**: Kubernetes原生的后端服务连接管理
- **WebSocket订阅**: 完整的以太坊WebSocket订阅支持（newHeads, logs, pendingTransactions等）

### 📊 监控和管理
- **实时监控**: 全面的性能指标和统计信息
- **健康检查**: 内置健康检查端点
- **日志记录**: 结构化日志支持
- **配置管理**: YAML 配置文件支持

### 🔐 安全特性
- **认证系统**: API 密钥和 JWT 令牌验证
- **请求验证**: 严格的请求参数验证
- **连接管理**: WebSocket 连接生命周期管理
- **CORS 支持**: 跨域请求控制

## 项目结构

```
yao-portal/
├── cmd/
│   └── main.go                 # 主程序入口
├── internal/
│   ├── config/
│   │   └── config.go          # 配置管理
│   └── server/
│       ├── server.go          # 主服务器实现
│       ├── jsonrpc.go         # JSON-RPC 服务实现
│       ├── txpool.go          # 交易池服务实现
│       └── ratelimiter.go     # 限流器服务实现
├── interfaces/
│   └── portal.go              # 接口定义
├── configs/
│   └── portal.yaml            # 配置文件模板
├── go.mod                     # Go 模块定义
├── go.sum                     # 依赖锁定文件
└── README.md                  # 项目文档
```

## 核心组件

### 1. JSONRPCService (JSON-RPC服务)
- 实现完整的以太坊 JSON-RPC 协议
- 支持交易提交、状态查询、区块查询等
- 集成交易池和状态读取器
- 提供详细的请求/响应指标

### 2. TxPoolService (交易池服务)
- 交易验证和管理
- 支持待处理和队列交易分离
- 自动清理过期交易
- 丰富的交易状态跟踪

### 3. RateLimiter (限流器)
- 基于令牌桶算法的限流
- 支持全局和客户端级别限流
- 自动封禁恶意客户端
- 可配置的限流策略

### 4. ConnectionManager (连接管理器)
- WebSocket 连接生命周期管理
- 订阅/取消订阅功能
- 消息广播机制
- 连接监控和清理

### 5. AuthenticationService (认证服务)
- API 密钥管理和验证
- JWT 令牌支持
- 请求来源验证
- 细粒度权限控制

## 配置说明

### 基础配置
```yaml
# 节点标识
nodeId: "portal-node-1"
nodeName: "YaoPortal Node 1"
chainId: 1337

# 服务器配置
httpBindAddress: "0.0.0.0"
httpPort: 8545
wsBindAddress: "0.0.0.0"
wsPort: 8546

# 交易池配置
txPoolCapacity: 10000
txPoolTimeout: "300s"

# 后端服务配置（Kubernetes原生）
backendServices:
  consensusService:
    name: "yao-vulcan-service"
    namespace: "yao-verse"
    endpoints:
      - "http://yao-vulcan-service.yao-verse.svc.cluster.local:8080"
    healthCheckPath: "/health"
  
  stateService:
    name: "yao-oracle-service"
    namespace: "yao-verse"
    endpoints:
      - "http://yao-oracle-service.yao-verse.svc.cluster.local:8080"
    healthCheckPath: "/health"

# 限流配置
enableRateLimit: true
requestsPerSecond: 1000
burstSize: 2000
```

### 高级配置
```yaml
# 监控配置
enableMetrics: true
metricsInterval: "10s"
enableHealthCheck: true

# 日志配置
logLevel: "info"
logFormat: "json"
logFile: "logs/portal.log"

# 消息队列配置
rocketmqEndpoints:
  - "localhost:9876"
producerGroup: "yao-portal-producer"
```

## API 接口

### JSON-RPC 接口
YaoPortal 完全兼容以太坊 JSON-RPC 协议，支持以下方法：

#### 交易相关
- `eth_sendRawTransaction` - 发送原始交易
- `eth_getTransactionByHash` - 根据哈希获取交易
- `eth_getTransactionReceipt` - 获取交易收据
- `eth_getTransactionCount` - 获取账户nonce

#### 账户相关
- `eth_getBalance` - 获取账户余额
- `eth_getCode` - 获取合约代码
- `eth_getStorageAt` - 获取存储值

#### 区块相关
- `eth_getBlockByNumber` - 根据区块号获取区块
- `eth_getBlockByHash` - 根据区块哈希获取区块
- `eth_blockNumber` - 获取当前区块号

#### 调用相关
- `eth_call` - 执行只读调用
- `eth_estimateGas` - 估算gas消耗

#### 网络相关
- `eth_chainId` - 获取链ID
- `eth_gasPrice` - 获取gas价格
- `eth_getLogs` - 获取事件日志
- `net_version` - 获取网络版本
- `web3_clientVersion` - 获取客户端版本
- `net_listening` - 获取网络监听状态
- `net_peerCount` - 获取peer节点数量

#### WebSocket订阅方法
- `eth_subscribe` - 创建订阅
  - `newHeads` - 订阅新区块头
  - `logs` - 订阅事件日志
  - `pendingTransactions` - 订阅待处理交易
  - `syncing` - 订阅同步状态
- `eth_unsubscribe` - 取消订阅

### 管理接口
- `GET /health` - 健康检查
- `GET /metrics` - 指标监控
- `GET /ws` - WebSocket连接端点
- `POST /admin/ratelimit` - 限流管理

## 性能指标

### 关键指标
- **请求处理能力**: 支持每秒1000+请求
- **交易池容量**: 默认10000交易缓存
- **连接数**: 支持1000并发WebSocket连接
- **响应时间**: P99响应时间 < 100ms

### 监控指标
- 总请求数和每秒请求数
- 成功/失败请求统计
- 交易池使用率和状态
- 限流统计和客户端状态
- 系统资源使用情况

## 部署说明

### 环境要求
- Go 1.21+
- RocketMQ 集群
- TiDB 数据库
- Redis 集群（可选）

### 构建和运行
```bash
# 构建
go build -o yao-portal ./cmd

# 运行
./yao-portal -config configs/portal.yaml

# 或使用环境变量
YAO_PORTAL_LOG_LEVEL=debug ./yao-portal
```

### Docker 部署
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o yao-portal ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/yao-portal .
COPY configs/portal.yaml ./configs/
CMD ["./yao-portal"]
```

## 开发指南

### 本地开发
```bash
# 安装依赖
go mod download

# 运行测试
go test ./...

# 代码格式化
go fmt ./...

# 静态检查
go vet ./...
```

### 扩展开发
1. **添加新的 JSON-RPC 方法**: 在 `jsonrpc.go` 中实现
2. **自定义限流策略**: 扩展 `ratelimiter.go`
3. **集成新的认证方式**: 修改 `AuthenticationService`
4. **添加监控指标**: 在相应服务中添加指标收集

## 故障排查

### 常见问题
1. **连接被拒绝**: 检查防火墙和端口配置
2. **限流触发**: 调整限流参数或检查客户端行为
3. **交易池满**: 增加容量或减少交易生存时间
4. **内存使用高**: 检查交易池和连接数配置

### 日志分析
```bash
# 查看错误日志
grep "ERROR" logs/portal.log

# 监控请求频率
grep "JSON-RPC" logs/portal.log | wc -l

# 分析限流情况
grep "rate limit" logs/portal.log
```

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 编写测试
4. 提交 Pull Request

## 许可证

本项目采用 MIT 许可证。 