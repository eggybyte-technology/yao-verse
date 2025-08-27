# YaoVerse (曜界) - 高性能区块链系统

## 项目概述

**YaoVerse (曜界)** 是一个专为受控内网环境设计的下一代区块链系统。其核心目标是突破传统区块链的性能瓶颈，提供一个兼容以太坊生态、具备极致处理能力和高可用性的中心化调度平台。

在"YaoVerse"的宇宙中，我们做出了一个根本性的架构抉择：**在完全可控和可信的环境中，以中心化的范式换取去中心化共识所无法企及的极致性能**。

## 核心特性

- 🚀 **极致性能**: 纳秒级状态访问，支持数万TPS
- 🔧 **以太坊兼容**: 完整支持JSON-RPC接口和EVM执行
- 🏗️ **微服务架构**: 职责分离，独立扩展
- 💾 **三级缓存**: L1本地缓存 + L2 Redis + L3 TiDB
- 🔄 **实时同步**: 基于RocketMQ的缓存失效机制
- 📊 **可观测性**: 全面的监控指标和健康检查
- ⚡ **水平扩展**: 无状态设计，支持动态扩缩容

## 系统架构

YaoVerse采用微服务架构，包含四个核心组件：

### 1. YaoPortal (曜阙) - API网关
- 提供以太坊兼容的JSON-RPC接口
- 交易验证和路由
- 负载均衡和限流
- 只读查询处理

### 2. YaoChronos (曜史) - 中心化序列器
- 全局唯一的交易排序器
- 区块构建和分发
- Leader选举和高可用
- 确保系统确定性

### 3. YaoVulcan (曜焱) - 执行节点集群
- 智能合约执行引擎
- 内置EVM和状态管理
- 嵌入式YaoOracle状态库
- 支持水平扩展

### 4. YaoArchive (曜典) - 持久化服务
- 数据最终落地
- ACID事务保证
- 缓存失效通知
- 数据完整性验证

## 技术栈

- **编程语言**: Go 1.21+
- **容器编排**: Kubernetes
- **数据库**: TiDB (分布式关系型数据库)
- **缓存**: Redis Cluster
- **消息队列**: Apache RocketMQ
- **监控**: Prometheus + Grafana
- **日志**: ELK Stack

## 快速开始

### 环境要求

- Kubernetes 1.26+
- Go 1.21+
- Docker 24.0+

### 部署

1. **克隆项目**
```bash
git clone https://github.com/eggybyte-technology/yao-verse.git
cd yao-verse
```

2. **构建镜像**
```bash
make build-images
```

3. **部署到Kubernetes**
```bash
kubectl apply -f deployments/
```

4. **验证部署**
```bash
kubectl get pods -n yao-verse
```

### 配置

每个组件都有独立的配置文件，详情请参阅：
- [YaoPortal配置](./components/portal.md#配置)
- [YaoChronos配置](./components/chronos.md#配置)
- [YaoVulcan配置](./components/vulcan.md#配置)
- [YaoArchive配置](./components/archive.md#配置)

## 文档结构

```
docs/
├── README.md                 # 项目概述 (本文档)
├── architecture.md           # 系统架构详述
├── api-reference/            # API接口文档
│   ├── portal-api.md         # YaoPortal JSON-RPC API
│   ├── chronos-api.md        # YaoChronos 管理API
│   ├── vulcan-api.md         # YaoVulcan 管理API
│   └── archive-api.md        # YaoArchive 管理API
├── components/               # 组件详细文档
│   ├── portal.md             # YaoPortal 组件文档
│   ├── chronos.md            # YaoChronos 组件文档
│   ├── vulcan.md             # YaoVulcan 组件文档
│   ├── archive.md            # YaoArchive 组件文档
│   └── oracle.md             # YaoOracle 状态库文档
├── deployment/               # 部署指南
│   ├── kubernetes.md         # Kubernetes 部署
│   ├── configuration.md      # 配置管理
│   └── monitoring.md         # 监控配置
├── development/              # 开发指南
│   ├── coding-standards.md   # 代码规范
│   ├── testing.md           # 测试指南
│   ├── contributing.md      # 贡献指南
│   └── debugging.md         # 调试指南
└── operations/               # 运维指南
    ├── troubleshooting.md    # 故障排除
    ├── performance.md        # 性能调优
    └── security.md           # 安全指南
```

## 性能指标

在标准测试环境下，YaoVerse能够实现：

- **TPS**: 50,000+ 交易每秒
- **延迟**: 
  - L1缓存命中: < 10ns
  - L2缓存命中: < 1ms
  - L3数据库查询: < 10ms
- **吞吐量**: 单节点处理 10,000+ 区块每小时
- **可用性**: 99.9%+ 系统可用性

## 贡献

我们欢迎社区贡献！请参阅[贡献指南](./development/contributing.md)了解详情。

## 许可证

本项目采用 Apache 2.0 许可证 - 详见 [LICENSE](../LICENSE) 文件。

## 联系我们

- 项目主页: https://github.com/eggybyte-technology/yao-verse
- 问题反馈: https://github.com/eggybyte-technology/yao-verse/issues
- 邮件: yao-verse@eggybyte.tech

---

© 2024 EggyByte Technology. All rights reserved. 