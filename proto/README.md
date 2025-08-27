# YaoVerse gRPC Interfaces

本目录包含YaoVerse系统的所有gRPC服务定义，基于现代微服务架构设计，支持高性能区块链处理。

## 架构概览

YaoVerse采用统一的gRPC接口设计，所有组件间通信都通过类型安全的Protocol Buffers定义。系统包含四个核心组件：

- **YaoPortal (曜阙)** - API网关和用户接口层
- **YaoChronos (曜史)** - 中心化序列器和共识管理
- **YaoVulcan (曜焱)** - 区块执行引擎和状态计算
- **YaoArchive (曜典)** - 数据持久化和存储管理

## 目录结构

```
proto/
├── common/           # 通用类型定义
│   └── types.proto   # 共享的数据结构
├── portal/           # YaoPortal服务定义
│   ├── state_service.proto   # 状态查询服务
│   ├── block_service.proto   # 区块查询服务
│   └── health_service.proto  # 健康检查服务
├── chronos/          # YaoChronos服务定义
│   ├── sequencer_service.proto # 序列器服务
│   └── consensus_service.proto # 共识服务
├── vulcan/           # YaoVulcan服务定义
│   ├── execution_service.proto # 执行服务
│   └── oracle_service.proto    # Oracle缓存服务
├── archive/          # YaoArchive服务定义
│   └── storage_service.proto   # 存储服务
├── buf.yaml          # Buf配置文件
├── buf.gen.yaml      # 代码生成配置
├── Makefile         # 构建脚本
└── README.md        # 本文档
```

## 核心设计原则

### 1. 类型安全
所有接口都使用Protocol Buffers定义，确保编译时类型检查和跨语言兼容性。

### 2. 统一元数据
所有请求/响应都包含标准化的元数据：
- `RequestMetadata`: 请求ID、时间戳、客户端版本
- `ResponseMetadata`: 响应ID、处理时间、服务器版本

### 3. 错误处理
采用gRPC标准错误模型，支持详细的错误信息和状态码。

### 4. HTTP网关支持
通过gRPC-Gateway自动生成RESTful HTTP接口，支持JSON格式。

### 5. 流式处理
支持Server Streaming，用于实时数据推送和大数据传输。

## 服务详解

### YaoPortal Services

**StateService** - 区块链状态查询
- `GetAccount`: 获取账户信息
- `GetBalance`: 查询账户余额
- `GetCode`: 获取合约代码
- `GetStorageAt`: 查询合约存储
- `EstimateGas`: 估算Gas消耗
- `Call`: 执行只读合约调用

**BlockService** - 区块和交易查询
- `GetBlockByNumber/Hash`: 获取区块信息
- `GetTransaction`: 查询交易详情
- `GetTransactionReceipt`: 获取交易回执
- `GetLogs`: 查询事件日志
- `SubscribeNewBlocks`: 订阅新区块 (流式)
- `SubscribeLogs`: 订阅日志事件 (流式)

**HealthService** - 健康检查和服务发现
- `Check`: 基础健康检查
- `CheckDeep`: 深度健康检查（包含依赖项）
- `GetVersion`: 获取版本信息
- `GetServiceInfo`: 获取服务详情
- `ListServices`: 列出可用服务

### YaoChronos Services

**SequencerService** - 交易排序和区块构建
- `SubmitTransaction`: 提交交易到序列器
- `GetPendingTransactions`: 查询待处理交易
- `GetBlockTemplate`: 获取当前区块模板
- `WatchBlocks`: 订阅区块事件 (流式)

**ConsensusService** - 共识状态管理
- `GetConsensusStatus`: 获取共识状态
- `GetLeaderInfo`: 查询Leader信息
- `GetValidators`: 获取验证器列表
- `WatchConsensusEvents`: 监听共识事件 (流式)

### YaoVulcan Services

**ExecutionService** - 区块和交易执行
- `ExecuteBlock`: 执行完整区块
- `ExecuteTransaction`: 执行单个交易
- `ValidateTransaction`: 验证交易
- `SimulateTransaction`: 模拟执行（Gas估算）
- `GetExecutionStatus`: 获取执行引擎状态
- `WatchExecution`: 监听执行事件 (流式)

**OracleService** - 三级缓存状态访问 (YaoOracle)
- `GetState`: 通过缓存层次获取状态
- `GetStateBatch`: 批量状态查询
- `PutState`: 存储状态到缓存
- `InvalidateState`: 失效缓存条目
- `WarmupCache`: 缓存预热
- `PurgeCache`: 缓存清理

### YaoArchive Services

**StorageService** - 数据持久化
- `StoreBlock`: 持久化区块数据
- `StoreReceipts`: 存储交易回执
- `StoreState`: 保存状态变更
- `GetHistoricalState`: 查询历史状态
- `ValidateChainIntegrity`: 验证链完整性
- `CreateBackup`: 创建数据备份

## 通信模式

### 同步gRPC调用
适用于：
- 实时查询（账户余额、区块信息）
- 状态验证（交易验证、Gas估算）
- 配置管理（动态配置、服务发现）
- 健康检查（服务状态、可用性探测）

### 异步消息队列
适用于：
- 核心业务流程（交易处理、区块执行）
- 事件驱动（缓存失效、状态变更）
- 削峰填谷（高并发处理）

### 流式处理
适用于：
- 实时通知（新区块、日志事件）
- 大数据传输（历史数据同步）
- 监控告警（系统状态变化）

## 缓存层次 (YaoOracle)

YaoOracle实现三级缓存架构：

1. **L1 Cache** - 进程内存缓存 (Go-Ristretto)
   - 纳秒级访问速度
   - 最热数据存储
   - 自动LRU淘汰

2. **L2 Cache** - 分布式缓存 (Redis)
   - 毫秒级访问速度
   - 节点间共享缓存
   - 支持失效广播

3. **L3 Storage** - 持久化存储 (TiDB)
   - 数十毫秒访问
   - 最终数据源
   - ACID事务保障

## 快速开始

### 1. 安装依赖
```bash
make deps
```

### 2. 生成代码
```bash
make generate
```

### 3. 开发工作流
```bash
make dev  # 格式化 + 检查 + 生成
```

### 4. CI/CD工作流
```bash
make ci   # 检查 + 兼容性 + 生成
```

## 开发指南

### 添加新服务
1. 在对应组件目录下创建`.proto`文件
2. 定义服务接口和消息类型
3. 添加HTTP注解（可选）
4. 运行`make generate`生成代码

### 修改现有接口
1. 确保向后兼容性
2. 运行`make breaking`检查破坏性变更
3. 更新版本号和文档
4. 重新生成代码

### 最佳实践
- 使用语义化的方法命名
- 为所有字段添加注释
- 合理使用oneof用于互斥字段
- 预留字段编号用于未来扩展
- 使用枚举表示状态和类型

## 生成的代码

运行`make generate`后，将生成：

- **Go代码**: `../shared/gen/proto/` - gRPC服务器和客户端代码
- **HTTP网关**: 自动生成的RESTful HTTP接口
- **OpenAPI文档**: `../docs/api-reference/` - Swagger/OpenAPI规范

## 工具链

- **Buf**: Protocol Buffer管理和代码生成
- **protoc-gen-go**: Go语言绑定生成器
- **protoc-gen-go-grpc**: Go gRPC代码生成器
- **grpc-gateway**: HTTP/JSON网关生成器
- **protoc-gen-openapiv2**: OpenAPI文档生成器

## 版本管理

遵循语义化版本控制：
- **主版本号**: 不兼容的API变更
- **次版本号**: 向后兼容的功能新增
- **修订版本号**: 向后兼容的错误修复

## 监控和可观测性

所有gRPC服务都支持：
- **健康检查**: 标准化健康检查接口
- **指标收集**: Prometheus兼容的指标
- **链路追踪**: 分布式追踪支持
- **结构化日志**: 统一的日志格式

---

更多详细信息请参考[YaoVerse架构文档](../docs/architecture.md)。 