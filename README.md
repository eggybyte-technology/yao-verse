# YaoVerse 曜界

**YaoVerse (曜界)** 是一个专为受控内网环境设计的下一代高性能区块链系统。系统采用中心化调度架构，在完全可控和可信的环境中，以中心化的范式换取去中心化共识所无法企及的极致性能。

## 系统架构

YaoVerse 由四个核心服务组成：

### 1. YaoPortal (曜阙) - API网关
- **职责**: 以太坊兼容的JSON-RPC API网关
- **功能**: 
  - 提供标准的以太坊JSON-RPC接口
  - 交易请求快速卸载到消息队列
  - 只读请求路由和处理
  - 交易池管理

### 2. YaoChronos (曜史) - 中心化序列器
- **职责**: 全局交易排序和区块打包
- **功能**:
  - 唯一的交易排序权威
  - 批量处理交易并打包成区块
  - Leader-Follower高可用架构
  - 支持多种排序策略

### 3. YaoVulcan (曜焱) - 执行节点集群
- **职责**: 智能合约执行和EVM运行
- **功能**:
  - 内置EVM执行交易
  - 集成YaoOracle三级缓存状态库
  - 并行处理不同区块
  - 实时缓存同步

### 4. YaoArchive (曜典) - 持久化服务
- **职责**: 数据最终落地和缓存失效通知
- **功能**:
  - 原子性数据写入TiDB
  - 广播缓存失效通知
  - 数据索引和验证
  - 备份和恢复

## 核心技术特性

### 🚀 三级缓存架构 (YaoOracle)
- **L1**: 进程内纳秒级本地缓存 (Ristretto v2)
- **L2**: Redis Cluster毫秒级共享缓存  
- **L3**: TiDB持久化存储
- **实时同步**: 基于RocketMQ广播的精确缓存失效

### 📊 消息驱动架构
- **tx-ingress-topic**: 交易提交队列
- **block-execute-topic**: 区块执行队列（集群消费）
- **state-commit-topic**: 状态提交队列（集群消费）
- **cache-invalidation-topic**: 缓存失效队列（广播消费）

### 🔧 核心中间件
- **TiDB**: 分布式关系型数据库，支持水平扩展
- **Redis Cluster**: 分布式共享缓存，cluster模式部署
- **RocketMQ**: 企业级消息队列，支持集群和广播消费模式

## 项目结构

```
yao-verse/
├── shared/                 # 共享类型和接口定义
│   ├── types/             # 公共数据类型
│   ├── interfaces/        # 核心接口定义
│   ├── cache/             # 缓存实现 (L1本地缓存 + L2Redis集群)
│   ├── storage/           # 存储实现 (L3 TiDB)
│   ├── mq/                # 消息队列实现 (RocketMQ)
│   ├── oracle/            # YaoOracle完整实现
│   └── examples/          # 使用示例
├── yao-portal/            # API网关服务
│   ├── cmd/               # 主程序入口
│   ├── interfaces/        # Portal特定接口
│   └── internal/          # 内部实现
├── yao-chronos/           # 序列器服务
│   ├── cmd/               # 主程序入口
│   └── interfaces/        # Chronos特定接口
├── yao-vulcan/            # 执行节点服务
│   ├── cmd/               # 主程序入口
│   └── interfaces/        # Vulcan特定接口
├── yao-archive/           # 存档服务
│   ├── cmd/               # 主程序入口
│   └── interfaces/        # Archive特定接口
├── go.work               # Go workspace配置
└── README.md             # 项目说明
```

## 快速开始

### 环境要求
- Go 1.24.5+
- TiDB集群
- Redis Cluster集群
- RocketMQ集群
- Kubernetes (推荐)

### 构建项目
```bash
# 克隆项目
git clone https://github.com/eggybyte-technology/yao-verse.git
cd yao-verse

# 构建所有服务
go work sync
go build ./yao-portal/cmd
go build ./yao-chronos/cmd  
go build ./yao-vulcan/cmd
go build ./yao-archive/cmd
```

### 使用共享库

#### YaoOracle 三级缓存示例
```go
package main

import (
    "context"
    "log"
    "math/big"
    
    "github.com/ethereum/go-ethereum/common"
    "github.com/eggybyte-technology/yao-verse-shared/interfaces"
    "github.com/eggybyte-technology/yao-verse-shared/oracle"
    "github.com/eggybyte-technology/yao-verse-shared/types"
)

func main() {
    // 创建 YaoOracle 实例
    yaoOracle, err := oracle.NewYaoOracle()
    if err != nil {
        log.Fatal("Failed to create YaoOracle:", err)
    }
    defer yaoOracle.Close()

    // 配置三级缓存
    config := &interfaces.OracleConfig{
        // L1本地缓存配置 (Ristretto v2)
        L1CacheSize: 10000,
        L1CacheTTL:  300,
        
        // L2 Redis Cluster配置
        RedisEndpoints: []string{
            "redis-cluster-node1:6379",
            "redis-cluster-node2:6379", 
            "redis-cluster-node3:6379",
        },
        RedisPoolSize: 20,
        L2CacheTTL: 1800,
        
        // L3 TiDB配置
        TiDBDSN: "root:@tcp(tidb-cluster:4000)/yaoverse",
        TiDBMaxOpenConn: 100,
        TiDBMaxIdleConn: 10,
        
        EnableMetrics: true,
        MetricsInterval: 10,
    }

    // 初始化并启动
    ctx := context.Background()
    yaoOracle.Initialize(ctx, config)
    yaoOracle.Start(ctx)

    // 使用示例
    testAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
    
    // 设置账户 (自动缓存到L1和L2)
    account := &types.Account{
        Address: testAddr,
        Nonce:   1,
        Balance: big.NewInt(1000000000000000000), // 1 ETH
    }
    yaoOracle.SetAccount(ctx, testAddr, account)

    // 获取账户 (纳秒级L1缓存命中)
    retrievedAccount, _ := yaoOracle.GetAccount(ctx, testAddr)
    log.Printf("Account balance: %s", retrievedAccount.Balance.String())

    // 查看性能指标
    metrics := yaoOracle.GetMetrics()
    log.Printf("L1 Hit Ratio: %.2f%%", metrics.L1Stats.HitRatio*100)
    log.Printf("L2 Hit Ratio: %.2f%%", metrics.L2Stats.HitRatio*100)
}
```

#### Redis Cluster 适配特性
- **集群安全操作**: 自动适配Redis Cluster的分片特性
- **模式检测**: 支持单机和集群模式自适应
- **批量操作**: 针对大批量操作进行优化
- **故障恢复**: 内置集群节点故障处理

#### RocketMQ 消息队列示例
```go
package main

import (
    "log"
    "time"
    
    "github.com/eggybyte-technology/yao-verse-shared/mq"
    "github.com/eggybyte-technology/yao-verse-shared/types"
)

func main() {
    config := &mq.RocketMQConfig{
        Endpoints:     []string{"rocketmq-nameserver1:9876", "rocketmq-nameserver2:9876"},
        ProducerGroup: "yao-producer",
        ConsumerGroup: "yao-consumer",
        RetryTimes:    3,
    }

    // 创建生产者
    producer, err := mq.NewRocketMQProducer(config)
    if err != nil {
        log.Fatal(err)
    }
    defer producer.Close()

    // 发送消息
    msg := &types.TxSubmissionMessage{
        Type: types.MessageTypeTxSubmission,
        MessageID: "tx-001",
        SubmittedAt: time.Now(),
    }
    producer.SendTxSubmission(msg)

    // 创建消费者（支持集群和广播模式）
    consumer, _ := mq.NewRocketMQConsumer(config)
    
    // 集群消费模式（负载均衡）
    consumer.Subscribe(types.TopicTxIngress, config.ConsumerGroup, &MyHandler{})
    
    // 广播消费模式（缓存失效）
    consumer.SubscribeBroadcast(types.TopicCacheInvalidation, &CacheInvalidationHandler{})
    
    consumer.Start()
}
```

更多详细示例请参考 `shared/examples/` 目录。

### 启动服务

#### 1. 启动 YaoPortal (API网关)
```bash
./yao-portal -config configs/portal.yaml
```

#### 2. 启动 YaoChronos (序列器)
```bash  
./yao-chronos -config configs/chronos.yaml
```

#### 3. 启动 YaoVulcan (执行节点)
```bash
./yao-vulcan -config configs/vulcan.yaml
```

#### 4. 启动 YaoArchive (存档服务)
```bash
./yao-archive -config configs/archive.yaml
```

## 性能特点

- **极致TPS**: 中心化架构消除共识开销，实现极高交易吞吐量
- **纳秒级状态访问**: L1本地缓存提供纳秒级状态读取性能  
- **集群级扩展**: Redis Cluster和无状态执行节点支持动态扩缩容
- **实时一致性**: 基于RocketMQ广播的精确缓存同步机制
- **以太坊兼容**: 完全兼容以太坊JSON-RPC接口

## 设计哲学

1. **中心化信任**: 在可控环境中以中心化换取极致性能
2. **职责分离**: 微服务化架构，各组件独立扩展
3. **计算存储分离**: 无状态执行节点，统一状态存储
4. **消息驱动**: 异步消息队列实现组件解耦
5. **简化存储**: 关系型数据库替代复杂的MPT树结构

## 开发状态

当前项目已完成核心共享库的实现：

- ✅ 核心接口定义完成
- ✅ 项目结构搭建完成
- ✅ Go模块配置完成
- ✅ **YaoOracle三级缓存架构实现完成**
  - ✅ L1本地内存缓存 (Ristretto v2 - 泛型支持)
  - ✅ L2分布式缓存 (Redis Cluster适配)
  - ✅ L3持久化存储 (TiDB完整实现)
  - ✅ 实时缓存失效机制
- ✅ **TiDB存储层完整实现**
  - ✅ 完整的区块链数据表结构
  - ✅ 账户、合约、存储数据管理
  - ✅ 事务支持和数据一致性
  - ✅ 集群连接池管理
- ✅ **RocketMQ消息队列实现**
  - ✅ 生产者和消费者实现
  - ✅ 集群和广播消费模式
  - ✅ 四种核心Topic支持
  - ✅ 高可用和重试机制
- ✅ **完整使用示例和文档**
- ✅ **类型系统和接口完善**
- 🚧 各服务具体业务逻辑实现中
- 🚧 单元测试编写中
- 🚧 部署和运维配置中

## 技术栈版本

### 核心依赖
- **Go**: 1.24.5+
- **Ristretto**: v2.3.0 (L1缓存，支持泛型)
- **Redis**: go-redis/v9 (Redis Cluster客户端)
- **RocketMQ**: v2.1.2 (消息队列)
- **Go-Ethereum**: v1.16.2 (以太坊兼容性)

### 中间件
- **TiDB**: 最新稳定版 (分布式数据库)
- **Redis Cluster**: 7.0+ (分布式缓存)
- **Apache RocketMQ**: 4.9+ (消息队列)

## 贡献指南

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系方式

- 项目地址: https://github.com/eggybyte-technology/yao-verse
- 问题反馈: https://github.com/eggybyte-technology/yao-verse/issues

---

**YaoVerse** - 为高性能私有链应用而生的区块链基础设施 