# YaoVerse 快速开始指南

## 概述

本指南将帮助您在15分钟内快速部署并运行一个完整的YaoVerse区块链网络。我们将使用Docker Compose来简化部署过程，让您能够快速体验YaoVerse的高性能特性。

## 环境要求

### 最小配置要求

- **CPU**: 4核心以上
- **内存**: 8GB以上
- **存储**: 50GB可用空间
- **网络**: 良好的网络连接

### 推荐配置

- **CPU**: 8核心以上
- **内存**: 16GB以上
- **存储**: 100GB SSD
- **网络**: 1Gbps网络

### 软件依赖

- Docker 24.0+
- Docker Compose 2.0+
- Git 2.30+
- curl 或 wget（用于测试）

```bash
# 验证Docker安装
docker --version
docker-compose --version

# 验证Docker运行状态
docker ps
```

## 快速部署

### 1. 获取项目代码

```bash
# 克隆项目
git clone https://github.com/eggybyte-technology/yao-verse.git
cd yao-verse

# 检查项目结构
ls -la
```

### 2. 配置环境

```bash
# 复制环境配置模板
cp .env.example .env

# 编辑配置文件（可选，默认配置已可用）
# vim .env
```

默认配置内容：

```bash
# .env
# 网络配置
CHAIN_ID=1337
NETWORK_ID=1337

# 数据库配置
TIDB_PASSWORD=yaoverse_password
REDIS_PASSWORD=redis_password

# RocketMQ配置
ROCKETMQ_ACCESS_KEY=your_access_key
ROCKETMQ_SECRET_KEY=your_secret_key

# 服务端口
PORTAL_PORT=8545
PORTAL_MGMT_PORT=9090
CHRONOS_MGMT_PORT=9091
VULCAN_MGMT_PORT=9092
ARCHIVE_MGMT_PORT=9093

# 监控端口
PROMETHEUS_PORT=9000
GRAFANA_PORT=3000
```

### 3. 启动基础设施

```bash
# 启动基础设施服务（TiDB、Redis、RocketMQ）
docker-compose -f docker-compose.infrastructure.yml up -d

# 等待服务启动（约30秒）
echo "等待基础设施启动..."
sleep 30

# 检查服务状态
docker-compose -f docker-compose.infrastructure.yml ps
```

### 4. 初始化数据库

```bash
# 运行数据库初始化脚本
docker-compose -f docker-compose.infrastructure.yml exec tidb-server \
    mysql -h127.0.0.1 -P4000 -uroot < scripts/init-database.sql

# 验证数据库初始化
docker-compose -f docker-compose.infrastructure.yml exec tidb-server \
    mysql -h127.0.0.1 -P4000 -uroot -e "SHOW DATABASES;"
```

### 5. 启动YaoVerse服务

```bash
# 启动YaoVerse核心服务
docker-compose -f docker-compose.yaoverse.yml up -d

# 等待服务启动
echo "等待YaoVerse服务启动..."
sleep 60

# 检查服务状态
docker-compose -f docker-compose.yaoverse.yml ps
```

### 6. 验证部署

```bash
# 检查所有服务状态
./scripts/check-cluster-health.sh

# 验证Portal JSON-RPC接口
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_chainId",
        "params": [],
        "id": 1
    }'

# 预期输出：
# {"jsonrpc":"2.0","result":"0x539","id":1}
```

## 基础使用

### 1. 获取网络信息

```bash
# 获取链ID
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_chainId",
        "params": [],
        "id": 1
    }'

# 获取当前区块号
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_blockNumber",
        "params": [],
        "id": 1
    }'

# 获取Gas价格
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_gasPrice",
        "params": [],
        "id": 1
    }'
```

### 2. 创建测试账户

我们提供了一些预配置的测试账户：

```bash
# 测试账户信息
ACCOUNT1_ADDRESS="0x1234567890123456789012345678901234567890"
ACCOUNT1_PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

ACCOUNT2_ADDRESS="0x0987654321098765432109876543210987654321"  
ACCOUNT2_PRIVATE_KEY="0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d"

# 查询账户余额
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_getBalance",
        "params": ["'$ACCOUNT1_ADDRESS'", "latest"],
        "id": 1
    }'
```

### 3. 发送交易

```bash
# 使用提供的脚本发送测试交易
./scripts/send-test-transaction.sh

# 或手动构造交易（需要签名工具）
# 这里展示已签名的转账交易示例
SIGNED_TX="0xf86c808504a817c800825208940987654321098765432109876543210987654321880de0b6b3a76400008025a0c4b6c7b7f1a4e6a6c7b7a3a4c6b7c7b7f1a4e6a6c7b7a3a4c6b7c7b7f1a4e6a6a07f4e6a6c7b7a3a4c6b7c7b7f1a4e6a6c7b7a3a4c6b7c7b7f1a4e6a6c7b7a3a4"

curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_sendRawTransaction",
        "params": ["'$SIGNED_TX'"],
        "id": 1
    }'
```

### 4. 查询交易状态

```bash
# 假设上一步返回的交易哈希为
TX_HASH="0xabc123def456..."

# 查询交易信息
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_getTransactionByHash",
        "params": ["'$TX_HASH'"],
        "id": 1
    }'

# 查询交易回执
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_getTransactionReceipt",
        "params": ["'$TX_HASH'"],
        "id": 1
    }'
```

## 智能合约部署

### 1. 准备合约代码

```solidity
// SimpleStorage.sol
pragma solidity ^0.8.0;

contract SimpleStorage {
    uint256 private storedData;
    
    event DataStored(uint256 data);
    
    function set(uint256 x) public {
        storedData = x;
        emit DataStored(x);
    }
    
    function get() public view returns (uint256) {
        return storedData;
    }
}
```

编译后的字节码（示例）：
```
BYTECODE="0x608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806360fe47b11461003b5780636d4ce63c14610057575b600080fd5b610055600480360381019061005091906100a3565b610075565b005b61005f61007f565b60405161006c91906100df565b60405180910390f35b8060008190555050565b60008054905090565b600080fd5b6000819050919050565b61009d8161008a565b81146100a857600080fd5b50565b6000813590506100ba81610094565b92915050565b6000602082840312156100d6576100d5610085565b5b60006100e4848285016100ab565b91505092915050565b6100f68161008a565b82525050565b600060208201905061011160008301846100ed565b9291505056fea264697066735822122012345678901234567890123456789012345678901234567890123456789012345664736f6c63430008110033"
```

### 2. 部署合约

```bash
# 构造合约部署交易
DEPLOY_DATA="0x608060405234801561001057600080fd5b50610150806100206000396000f3fe..."

# 发送部署交易（需要正确签名）
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_sendRawTransaction",
        "params": ["'$SIGNED_DEPLOY_TX'"],
        "id": 1
    }'

# 获取合约地址（从交易回执中）
CONTRACT_ADDRESS="0x5FbDB2315678afecb367f032d93F642f64180aa3"
```

### 3. 调用合约

```bash
# 调用合约的 get() 方法
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_call",
        "params": [{
            "to": "'$CONTRACT_ADDRESS'",
            "data": "0x6d4ce63c"
        }, "latest"],
        "id": 1
    }'

# 调用合约的 set(42) 方法（需要发送交易）
# data = set(uint256) selector + 42的编码
SET_42_DATA="0x60fe47b1000000000000000000000000000000000000000000000000000000000000002a"

curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_sendRawTransaction",
        "params": ["'$SIGNED_SET_TX'"],
        "id": 1
    }'
```

## 监控和管理

### 1. 访问管理界面

```bash
# 查看Prometheus监控
open http://localhost:9000

# 查看Grafana仪表板
open http://localhost:3000
# 默认登录: admin/admin
```

### 2. 检查系统状态

```bash
# 检查Portal状态
curl -s http://localhost:9090/api/v1/health | jq

# 检查Chronos状态
curl -s http://localhost:9091/api/v1/health | jq

# 检查Vulcan状态
curl -s http://localhost:9092/api/v1/health | jq

# 检查Archive状态
curl -s http://localhost:9093/api/v1/health | jq
```

### 3. 查看性能指标

```bash
# 获取Portal性能指标
curl -s http://localhost:9090/api/v1/metrics/json | jq

# 获取系统整体状态
./scripts/cluster-status.sh
```

## 性能测试

### 1. 基础性能测试

```bash
# 运行内置的性能测试脚本
./scripts/performance-test.sh

# 测试交易吞吐量
./scripts/throughput-test.sh --duration=60s --tps=100
```

### 2. 自定义负载测试

```bash
# 安装测试工具
npm install -g ethereum-load-tester

# 运行负载测试
ethereum-load-tester \
    --rpc-url http://localhost:8545 \
    --private-key $ACCOUNT1_PRIVATE_KEY \
    --to $ACCOUNT2_ADDRESS \
    --value 0.001 \
    --concurrent 10 \
    --duration 60
```

## 故障排除

### 常见问题及解决方案

#### 1. 服务无法启动

```bash
# 检查Docker资源使用
docker system df
docker system events

# 检查端口占用
netstat -tulpn | grep :8545

# 查看服务日志
docker-compose -f docker-compose.yaoverse.yml logs portal
docker-compose -f docker-compose.yaoverse.yml logs chronos
```

#### 2. 交易无法确认

```bash
# 检查Chronos是否正常生产区块
curl -s http://localhost:9091/api/v1/sequencer/status | jq

# 检查交易池状态
curl -s http://localhost:9090/api/v1/txpool/status | jq

# 查看最新区块
curl -X POST http://localhost:8545 \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "eth_getBlockByNumber",
        "params": ["latest", false],
        "id": 1
    }'
```

#### 3. 性能问题

```bash
# 检查缓存命中率
curl -s http://localhost:9092/api/v1/cache/stats | jq

# 检查数据库性能
docker-compose -f docker-compose.infrastructure.yml exec tidb-server \
    mysql -h127.0.0.1 -P4000 -uroot -e "SHOW PROCESSLIST;"

# 检查系统资源
docker stats
```

### 日志分析

```bash
# 查看所有服务日志
docker-compose -f docker-compose.yaoverse.yml logs -f

# 查看特定服务日志
docker-compose -f docker-compose.yaoverse.yml logs -f portal

# 搜索错误日志
docker-compose -f docker-compose.yaoverse.yml logs | grep -i error
```

## 进阶配置

### 1. 自定义配置

```bash
# 编辑服务配置
mkdir -p config/custom

# 复制并修改Portal配置
cp config/portal.yml config/custom/portal.yml
vim config/custom/portal.yml

# 重启服务应用新配置
docker-compose -f docker-compose.yaoverse.yml restart portal
```

### 2. 扩展部署

```bash
# 增加Vulcan节点
docker-compose -f docker-compose.yaoverse.yml up -d --scale vulcan=3

# 增加Archive节点
docker-compose -f docker-compose.yaoverse.yml up -d --scale archive=2
```

### 3. 生产环境部署

```bash
# 使用生产配置
docker-compose -f docker-compose.production.yml up -d

# 配置SSL证书
cp ssl/cert.pem config/tls/
cp ssl/key.pem config/tls/

# 启用认证
export ENABLE_AUTH=true
export JWT_SECRET=$(openssl rand -base64 32)
```

## 开发工具集成

### 1. MetaMask配置

```
网络名称: YaoVerse Local
RPC URL: http://localhost:8545
链ID: 1337
货币符号: ETH
区块浏览器URL: http://localhost:3001
```

### 2. Hardhat配置

```javascript
// hardhat.config.js
require("@nomiclabs/hardhat-waffle");

module.exports = {
  solidity: "0.8.19",
  networks: {
    yaoverse: {
      url: "http://localhost:8545",
      accounts: [
        "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
      ],
      chainId: 1337
    }
  }
};
```

### 3. Truffle配置

```javascript
// truffle-config.js
module.exports = {
  networks: {
    yaoverse: {
      host: "localhost",
      port: 8545,
      network_id: 1337,
      gas: 8000000,
      gasPrice: 1000000000
    }
  },
  compilers: {
    solc: {
      version: "0.8.19"
    }
  }
};
```

## 清理环境

```bash
# 停止所有服务
docker-compose -f docker-compose.yaoverse.yml down
docker-compose -f docker-compose.infrastructure.yml down

# 清理数据卷（注意：会删除所有数据）
docker-compose -f docker-compose.yaoverse.yml down -v
docker-compose -f docker-compose.infrastructure.yml down -v

# 清理镜像
docker image prune -f
```

## 下一步

现在您已经成功部署了YaoVerse网络，可以进一步探索：

1. [API参考文档](./api-reference/) - 了解完整的API接口
2. [系统架构](./architecture.md) - 深入理解系统设计
3. [组件文档](./components/) - 学习各组件的详细功能
4. [开发指南](./development/) - 参与YaoVerse开发
5. [运维指南](./operations/) - 生产环境运维

## 获得帮助

- 📖 [官方文档](https://docs.yaoverse.io)
- 💬 [Discord社区](https://discord.gg/yaoverse)
- 🐛 [问题报告](https://github.com/eggybyte-technology/yao-verse/issues)
- 📧 [技术支持](mailto:support@eggybyte.tech)

---

恭喜！您已经成功部署了YaoVerse高性能区块链网络。现在可以开始探索这个强大的区块链平台了。 