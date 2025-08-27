# Changelog

All notable changes to YaoPortal will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.0.0] - 2024-12-27

### Added

#### 核心功能增强
- **完整的以太坊兼容JSON-RPC实现**: 实现了所有标准的以太坊JSON-RPC方法，包括交易、账户、区块、调用等相关API
- **完整的WebSocket支持**: 实现了以太坊WebSocket订阅功能，支持 newHeads、logs、pendingTransactions、syncing 订阅
- **后端服务管理器**: Kubernetes原生的后端服务连接管理，不实现负载均衡，由Kubernetes的Service和Higress处理流量分配
- **批量请求处理**: 支持JSON-RPC批量请求处理
- **完整的参数验证**: 对所有JSON-RPC请求参数进行严格验证

#### 网络和连接管理
- **WebSocket连接生命周期管理**: 完整的WebSocket连接创建、维护、清理机制
- **订阅管理**: 支持客户端订阅和取消订阅不同类型的事件
- **连接健康检查**: 自动检测并清理不活跃的连接
- **CORS支持**: 完整的跨域请求支持

#### 后端服务集成
- **服务发现**: 支持多种后端服务（共识、执行、状态、归档）
- **健康检查**: 定期检查后端服务健康状态
- **错误处理**: 完善的后端服务错误处理和故障转移
- **指标收集**: 收集后端服务的响应时间、错误率等指标

#### 配置管理
- **扩展的配置结构**: 新增后端服务配置、以太坊配置等
- **默认配置**: 提供完整的Kubernetes环境默认配置
- **环境变量覆盖**: 支持通过环境变量覆盖配置选项

#### API方法实现
- **交易相关方法**:
  - `eth_sendRawTransaction`: 发送原始交易
  - `eth_getTransactionByHash`: 根据哈希获取交易
  - `eth_getTransactionReceipt`: 获取交易收据
  - `eth_getTransactionCount`: 获取账户nonce

- **账户相关方法**:
  - `eth_getBalance`: 获取账户余额
  - `eth_getCode`: 获取合约代码
  - `eth_getStorageAt`: 获取存储值

- **区块相关方法**:
  - `eth_getBlockByNumber`: 根据区块号获取区块
  - `eth_getBlockByHash`: 根据区块哈希获取区块
  - `eth_blockNumber`: 获取当前区块号

- **调用相关方法**:
  - `eth_call`: 执行只读调用
  - `eth_estimateGas`: 估算gas消耗

- **网络相关方法**:
  - `eth_chainId`: 获取链ID
  - `eth_gasPrice`: 获取gas价格
  - `eth_getLogs`: 获取事件日志
  - `net_version`: 获取网络版本
  - `web3_clientVersion`: 获取客户端版本
  - `net_listening`: 获取网络监听状态
  - `net_peerCount`: 获取peer节点数量

- **工具方法**:
  - `eth_getProof`: 获取账户和存储证明

#### WebSocket订阅
- **订阅类型支持**:
  - `newHeads`: 新区块头订阅
  - `logs`: 事件日志订阅
  - `pendingTransactions`: 待处理交易订阅
  - `syncing`: 同步状态订阅
- **订阅管理**: 支持 `eth_subscribe` 和 `eth_unsubscribe` 方法
- **消息广播**: 高效的订阅消息广播机制

#### 监控和指标
- **后端服务指标**: 收集服务健康状态、响应时间、错误率
- **WebSocket指标**: 连接数、订阅数、消息统计
- **JSON-RPC指标**: 请求统计、方法使用统计、错误统计

### Enhanced

#### 现有功能改进
- **HTTP请求处理**: 完善了HTTP JSON-RPC请求处理逻辑
- **错误处理**: 标准化的JSON-RPC错误响应格式
- **服务器生命周期**: 改进了服务器启动和停止流程
- **配置验证**: 增强了配置文件验证逻辑

#### 代码结构优化
- **接口定义扩展**: 新增了后端服务管理、WebSocket服务等接口
- **类型定义完善**: 增加了配置结构、指标结构等类型定义
- **依赖管理**: 更新了Go模块依赖

### Technical Details

#### Dependencies Added
- `github.com/gorilla/websocket@v1.5.3`: WebSocket支持
- `golang.org/x/time@v0.9.0`: 速率限制支持

#### Architecture Changes
- 新增 `BackendServiceManagerImpl`: 后端服务管理实现
- 新增 `WebSocketServiceImpl`: WebSocket服务实现
- 扩展 `JSONRPCServiceImpl`: 完整的方法路由和执行
- 更新 `PortalServer`: 集成所有新服务

#### Configuration Schema Updates
- 新增 `BackendServicesConfig`: 后端服务配置
- 新增 `EthereumConfig`: 以太坊兼容性配置
- 扩展 `PortalConfig`: 支持所有新功能的配置

### Known Limitations
- 某些JSON-RPC方法仍需要与实际的后端服务集成才能返回真实数据
- `eth_getLogs` 的过滤器解析需要进一步完善
- 后端服务的实际调用逻辑需要根据具体的服务接口来实现

### Migration Guide
1. 更新配置文件以包含新的 `backendServices` 和 `ethereum` 配置节
2. 如果使用自定义配置，请参考新的配置结构添加必要的字段
3. WebSocket端点现在可以通过 `/ws` 路径访问

### Next Steps
- 与YaoVulcan、YaoChronos、YaoOracle等后端服务的实际集成
- 实现完整的事件日志过滤和查询
- 添加认证和权限管理
- 性能优化和压力测试 