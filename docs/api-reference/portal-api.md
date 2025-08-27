# YaoPortal JSON-RPC API 参考

## 概述

YaoPortal提供完整的以太坊兼容JSON-RPC API，支持所有标准的以太坊客户端和工具。本文档详细描述了所有可用的API方法。

## 连接信息

### HTTP 端点
```
POST http://portal-host:8545/
Content-Type: application/json
```

### WebSocket 端点
```
ws://portal-host:8545/ws
```

### 基础请求格式
```json
{
    "jsonrpc": "2.0",
    "method": "method_name",
    "params": [...],
    "id": 1
}
```

### 基础响应格式
```json
{
    "jsonrpc": "2.0",
    "result": ...,
    "id": 1
}
```

### 错误响应格式
```json
{
    "jsonrpc": "2.0",
    "error": {
        "code": -32000,
        "message": "Error description",
        "data": "Additional error information"
    },
    "id": 1
}
```

## 网络信息方法

### net_version
返回当前网络 ID。

**参数**: 无

**返回值**:
- `String` - 网络 ID

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "net_version",
    "params": [],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "1337",
    "id": 1
}
```

### eth_chainId
返回链 ID，用于交易签名以防止重放攻击。

**参数**: 无

**返回值**:
- `QUANTITY` - 链 ID，十六进制格式

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_chainId",
    "params": [],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x539",
    "id": 1
}
```

### eth_gasPrice
返回当前燃料价格。

**参数**: 无

**返回值**:
- `QUANTITY` - 当前燃料价格，单位为 wei

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_gasPrice",
    "params": [],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x9184e72a000",
    "id": 1
}
```

## 区块信息方法

### eth_blockNumber
返回最新区块号。

**参数**: 无

**返回值**:
- `QUANTITY` - 最新区块号

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_blockNumber",
    "params": [],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x4b7",
    "id": 1
}
```

### eth_getBlockByNumber
通过区块号获取区块信息。

**参数**:
1. `QUANTITY|TAG` - 区块号，或字符串 "earliest", "latest", "pending"
2. `Boolean` - 如果为 true，返回完整的交易对象；如果为 false，仅返回交易哈希

**返回值**:
- `Object` - 区块对象，如果区块不存在则返回 null

**区块对象结构**:
```json
{
    "number": "0x1b4",           // 区块号
    "hash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
    "parentHash": "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
    "timestamp": "0x5a",         // 时间戳
    "gasLimit": "0x47e7c4",     // 燃料限制
    "gasUsed": "0x38658",       // 已使用燃料
    "stateRoot": "0xd5855eb08b3387c0af375e9cdb6acfc05eb8f519e419b874b6ff2ffda7ed1dff",
    "receiptHash": "0x4895c15c4d8f0c6d9bc0c6e3e5c1c4c65c5c5c5c5c5c5c5c5c5c5c5c5c5c",
    "transactions": [...]        // 交易数组
}
```

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_getBlockByNumber",
    "params": ["0x1b4", true],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": {
        "number": "0x1b4",
        "hash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
        "parentHash": "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
        "timestamp": "0x5a",
        "gasLimit": "0x47e7c4",
        "gasUsed": "0x38658",
        "stateRoot": "0xd5855eb08b3387c0af375e9cdb6acfc05eb8f519e419b874b6ff2ffda7ed1dff",
        "receiptHash": "0x4895c15c4d8f0c6d9bc0c6e3e5c1c4c65c5c5c5c5c5c5c5c5c5c5c5c5c5c",
        "transactions": [
            {
                "hash": "0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b",
                "blockNumber": "0x1b4",
                "from": "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
                "to": "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
                "value": "0xf3dbb76162000",
                "gas": "0x76c0",
                "gasPrice": "0x9184e72a000",
                "nonce": "0x117",
                "data": "0x",
                "v": "0x25",
                "r": "0xc4b6c7b7f1a4e6a6c7b7a3a4c6b7c7b7f1a4e6a6c7b7a3a4c6b7c7b7f1a4e6a6",
                "s": "0x7f4e6a6c7b7a3a4c6b7c7b7f1a4e6a6c7b7a3a4c6b7c7b7f1a4e6a6c7b7a3a4"
            }
        ]
    },
    "id": 1
}
```

### eth_getBlockByHash
通过区块哈希获取区块信息。

**参数**:
1. `DATA` - 区块哈希
2. `Boolean` - 如果为 true，返回完整的交易对象；如果为 false，仅返回交易哈希

**返回值**:
- `Object` - 区块对象，如果区块不存在则返回 null

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_getBlockByHash",
    "params": ["0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331", false],
    "id": 1
}
```

## 账户信息方法

### eth_getBalance
获取指定地址在指定区块的余额。

**参数**:
1. `DATA` - 地址
2. `QUANTITY|TAG` - 区块号，或字符串 "earliest", "latest", "pending"

**返回值**:
- `QUANTITY` - 余额，单位为 wei

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": ["0xa7d9ddbe1f17865597fbd27ec712455208b6b76d", "latest"],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x0234c8a3397aab58",
    "id": 1
}
```

### eth_getTransactionCount
获取指定地址的交易计数（nonce）。

**参数**:
1. `DATA` - 地址
2. `QUANTITY|TAG` - 区块号，或字符串 "earliest", "latest", "pending"

**返回值**:
- `QUANTITY` - 交易计数

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionCount",
    "params": ["0xa7d9ddbe1f17865597fbd27ec712455208b6b76d", "latest"],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x1",
    "id": 1
}
```

### eth_getCode
获取指定地址的合约代码。

**参数**:
1. `DATA` - 地址
2. `QUANTITY|TAG` - 区块号，或字符串 "earliest", "latest", "pending"

**返回值**:
- `DATA` - 合约代码

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_getCode",
    "params": ["0xa7d9ddbe1f17865597fbd27ec712455208b6b76d", "0x2"],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x600160008035811a818181146012578301005b601b6001356025565b8060005260206000f25b600060078202905091905056",
    "id": 1
}
```

### eth_getStorageAt
获取指定地址在指定存储位置的值。

**参数**:
1. `DATA` - 地址
2. `QUANTITY` - 存储位置
3. `QUANTITY|TAG` - 区块号，或字符串 "earliest", "latest", "pending"

**返回值**:
- `DATA` - 存储位置的值

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_getStorageAt",
    "params": ["0x295a70b2de5e3953354a6a8344e616ed314d7251", "0x0", "latest"],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x00000000000000000000000000000000000000000000000000000000000004d2",
    "id": 1
}
```

## 交易方法

### eth_sendRawTransaction
发送已签名的交易。

**参数**:
1. `DATA` - 已签名的交易数据

**返回值**:
- `DATA` - 交易哈希

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_sendRawTransaction",
    "params": ["0xf869018203e882520894f17f52151ebef6c7334fad080c5704d77216b732881bc16d674ec80000801ba02da1c48b670996dcb1f447ef9ef00b33033c48a4fe938f420bec3e56bfd24071a062e0aa78a81bf0290afbc3a9d8e9a068e6d74caa66c5e0a8a91a5b3c7b7b8b2"],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
    "id": 1
}
```

### eth_getTransactionByHash
通过交易哈希获取交易信息。

**参数**:
1. `DATA` - 交易哈希

**返回值**:
- `Object` - 交易对象，如果交易不存在则返回 null

**交易对象结构**:
```json
{
    "hash": "0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b",
    "blockNumber": "0x1b4",
    "blockHash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
    "transactionIndex": "0x1",
    "from": "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
    "to": "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
    "value": "0xf3dbb76162000",
    "gas": "0x76c0",
    "gasPrice": "0x9184e72a000",
    "nonce": "0x117",
    "data": "0x",
    "v": "0x25",
    "r": "0xc4b6c7b7f1a4e6a6c7b7a3a4c6b7c7b7f1a4e6a6c7b7a3a4c6b7c7b7f1a4e6a6",
    "s": "0x7f4e6a6c7b7a3a4c6b7c7b7f1a4e6a6c7b7a3a4c6b7c7b7f1a4e6a6c7b7a3a4"
}
```

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionByHash",
    "params": ["0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b"],
    "id": 1
}
```

### eth_getTransactionReceipt
获取交易回执。

**参数**:
1. `DATA` - 交易哈希

**返回值**:
- `Object` - 交易回执对象，如果交易不存在则返回 null

**回执对象结构**:
```json
{
    "transactionHash": "0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b",
    "transactionIndex": "0x1",
    "blockNumber": "0x1b4",
    "blockHash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
    "cumulativeGasUsed": "0x33bc",
    "gasUsed": "0x4dc",
    "contractAddress": null,
    "logs": [
        {
            "address": "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
            "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"],
            "data": "0x000000000000000000000000000000000000000000000000000000000000007b",
            "blockNumber": "0x1b4",
            "transactionHash": "0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b",
            "transactionIndex": "0x1",
            "blockHash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
            "logIndex": "0x1",
            "removed": false
        }
    ],
    "status": "0x1",
    "effectiveGasPrice": "0x9184e72a000"
}
```

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_getTransactionReceipt",
    "params": ["0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b"],
    "id": 1
}
```

## 合约调用方法

### eth_call
执行一个消息调用交易，该交易直接在节点的VM中执行，但不会被挖掘到区块链中。

**参数**:
1. `Object` - 调用对象
   - `from`: `DATA` - （可选）发送交易的地址
   - `to`: `DATA` - 交易的目标地址
   - `gas`: `QUANTITY` - （可选）为交易执行提供的燃料
   - `gasPrice`: `QUANTITY` - （可选）用于每个付费燃料的gasPrice
   - `value`: `QUANTITY` - （可选）与此交易一起发送的值
   - `data`: `DATA` - （可选）方法签名和编码参数的哈希值
2. `QUANTITY|TAG` - 区块号，或字符串 "earliest", "latest", "pending"

**返回值**:
- `DATA` - 执行结果

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_call",
    "params": [
        {
            "to": "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
            "data": "0x70a08231000000000000000000000000a7d9ddbe1f17865597fbd27ec712455208b6b76d"
        },
        "latest"
    ],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x000000000000000000000000000000000000000000000000000000000000007b",
    "id": 1
}
```

### eth_estimateGas
预估交易所需的燃料。

**参数**:
1. `Object` - 交易对象，格式与 eth_call 相同

**返回值**:
- `QUANTITY` - 预估的燃料数量

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_estimateGas",
    "params": [
        {
            "from": "0xa7d9ddbe1f17865597fbd27ec712455208b6b76d",
            "to": "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
            "value": "0x9184e72a"
        }
    ],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x5208",
    "id": 1
}
```

## 日志和过滤器方法

### eth_getLogs
根据过滤器获取日志。

**参数**:
1. `Object` - 过滤器对象
   - `fromBlock`: `QUANTITY|TAG` - （可选）开始区块号
   - `toBlock`: `QUANTITY|TAG` - （可选）结束区块号
   - `address`: `DATA|Array` - （可选）合约地址或地址数组
   - `topics`: `Array of DATA` - （可选）主题数组
   - `blockHash`: `DATA` - （可选）区块哈希，与 fromBlock/toBlock 互斥

**返回值**:
- `Array` - 日志对象数组

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_getLogs",
    "params": [
        {
            "fromBlock": "0x1",
            "toBlock": "0x2",
            "address": "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
            "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]
        }
    ],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": [
        {
            "address": "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
            "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"],
            "data": "0x000000000000000000000000000000000000000000000000000000000000007b",
            "blockNumber": "0x1b4",
            "transactionHash": "0xc6ef2fc5426d6ad6fd9e2a26abeab0aa2411b7ab17f30a99d3cb96aed1d1055b",
            "transactionIndex": "0x1",
            "blockHash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
            "logIndex": "0x1",
            "removed": false
        }
    ],
    "id": 1
}
```

## 管理方法

### yao_getNodeInfo
获取节点信息（YaoVerse 扩展方法）。

**参数**: 无

**返回值**:
- `Object` - 节点信息对象

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "yao_getNodeInfo",
    "params": [],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": {
        "nodeId": "portal-node-1",
        "version": "1.0.0",
        "chainId": 1337,
        "networkId": 1337,
        "uptime": "72h30m45s",
        "connections": 5,
        "pendingTransactions": 12
    },
    "id": 1
}
```

### yao_getMetrics
获取节点指标（YaoVerse 扩展方法）。

**参数**: 无

**返回值**:
- `Object` - 指标对象

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "yao_getMetrics",
    "params": [],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": {
        "totalRequests": 15420,
        "requestsPerSecond": 25.3,
        "averageResponseTime": 45,
        "pendingTransactions": 12,
        "submittedTxs": 8934,
        "droppedTxs": 23,
        "totalErrors": 156,
        "activeConnections": 8,
        "memoryUsage": 234567890,
        "cpuUsage": 12.5
    },
    "id": 1
}
```

## WebSocket 订阅

### eth_subscribe
订阅特定事件（仅适用于 WebSocket 连接）。

**支持的订阅类型**:
- `newHeads` - 新区块头
- `logs` - 匹配特定条件的日志
- `newPendingTransactions` - 新的待处理交易

**参数**:
1. `String` - 订阅类型
2. `Object` - （可选）订阅参数

**返回值**:
- `String` - 订阅 ID

**示例 - 订阅新区块**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_subscribe",
    "params": ["newHeads"],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": "0x9ce59a13059e417087c02d3236a0b1cc",
    "id": 1
}

// 通知
{
    "jsonrpc": "2.0",
    "method": "eth_subscription",
    "params": {
        "subscription": "0x9ce59a13059e417087c02d3236a0b1cc",
        "result": {
            "number": "0x1b4",
            "hash": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
            "parentHash": "0x9646252be9520f6e71339a8df9c55e4d7619deeb018d2a3f2d21fc165dde5eb5",
            "timestamp": "0x5a"
        }
    }
}
```

**示例 - 订阅日志**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_subscribe",
    "params": [
        "logs",
        {
            "address": "0xf02c1c8e6114b1dbe8937a39260b5b0a374432bb",
            "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]
        }
    ],
    "id": 1
}
```

### eth_unsubscribe
取消订阅。

**参数**:
1. `String` - 订阅 ID

**返回值**:
- `Boolean` - 取消成功返回 true

**示例**:
```json
// 请求
{
    "jsonrpc": "2.0",
    "method": "eth_unsubscribe",
    "params": ["0x9ce59a13059e417087c02d3236a0b1cc"],
    "id": 1
}

// 响应
{
    "jsonrpc": "2.0",
    "result": true,
    "id": 1
}
```

## 错误代码

YaoPortal 使用标准的 JSON-RPC 错误代码，以及一些自定义的错误代码：

| 错误代码 | 消息 | 说明 |
|---------|------|------|
| -32700 | Parse error | 无效的 JSON |
| -32600 | Invalid Request | JSON 不是有效的请求对象 |
| -32601 | Method not found | 方法不存在或不可用 |
| -32602 | Invalid params | 无效的方法参数 |
| -32603 | Internal error | JSON-RPC 内部错误 |
| -32000 | Server error | 通用服务器错误 |
| -32001 | Transaction rejected | 交易被拒绝 |
| -32002 | Insufficient balance | 余额不足 |
| -32003 | Gas limit exceeded | 燃料限制超出 |
| -32004 | Nonce too low | Nonce 过低 |
| -32005 | Nonce too high | Nonce 过高 |
| -32006 | Invalid signature | 无效的签名 |

## 性能建议

1. **批量请求**: 使用 JSON-RPC 批量请求减少网络往返
2. **WebSocket**: 对于实时更新，使用 WebSocket 连接
3. **缓存**: 对于静态数据（如已确认区块），启用客户端缓存
4. **分页**: 对于大量日志查询，使用适当的区块范围

## 限流策略

YaoPortal 实现了基于 IP 的限流策略：
- 默认每秒 100 个请求
- WebSocket 连接每个 IP 最多 10 个
- 大量查询操作（如 eth_getLogs）有额外限制

---

更多信息请参考 [YaoVerse 官方文档](https://github.com/eggybyte-technology/yao-verse)。 