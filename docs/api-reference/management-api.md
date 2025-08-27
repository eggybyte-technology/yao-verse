# YaoVerse 管理 API 参考

## 概述

YaoVerse各组件都提供统一的管理API接口，用于监控、健康检查、配置管理和运维操作。本文档描述了所有管理API的详细信息。

## 通用接口规范

所有组件都实现以下通用管理接口：

### 基础端点

每个组件都在独立的端口提供管理API：

- **YaoPortal**: `http://portal-host:9090/api/v1/`
- **YaoChronos**: `http://chronos-host:9091/api/v1/`
- **YaoVulcan**: `http://vulcan-host:9092/api/v1/`
- **YaoArchive**: `http://archive-host:9093/api/v1/`

### 认证

管理API使用Bearer Token认证：

```http
GET /api/v1/health
Authorization: Bearer <management-token>
```

## 通用管理接口

### 健康检查

#### GET /health
检查组件健康状态。

**参数**: 无

**响应**:
```json
{
    "status": "healthy|degraded|unhealthy",
    "timestamp": "2024-01-15T10:30:45Z",
    "component": "yao-portal",
    "version": "1.0.0",
    "checks": [
        {
            "name": "database",
            "status": "healthy",
            "message": "Connection established",
            "lastChecked": "2024-01-15T10:30:44Z",
            "duration": "5ms"
        },
        {
            "name": "message_queue",
            "status": "healthy", 
            "message": "RocketMQ accessible",
            "lastChecked": "2024-01-15T10:30:44Z",
            "duration": "12ms"
        }
    ],
    "uptime": "72h30m45s"
}
```

#### GET /health/ready
检查组件是否准备就绪。

**响应**:
```json
{
    "ready": true,
    "message": "Component is ready to serve requests"
}
```

#### GET /health/live
检查组件是否存活。

**响应**:
```json
{
    "alive": true,
    "message": "Component is alive"
}
```

### 指标监控

#### GET /metrics
获取组件性能指标（Prometheus格式）。

**响应**: Prometheus格式的指标数据
```
# HELP yao_portal_requests_total Total number of requests
# TYPE yao_portal_requests_total counter
yao_portal_requests_total{method="eth_sendRawTransaction"} 1543
yao_portal_requests_total{method="eth_getBalance"} 2987

# HELP yao_portal_request_duration_seconds Request duration in seconds
# TYPE yao_portal_request_duration_seconds histogram
yao_portal_request_duration_seconds_bucket{method="eth_sendRawTransaction",le="0.1"} 1234
yao_portal_request_duration_seconds_bucket{method="eth_sendRawTransaction",le="0.5"} 1543
```

#### GET /metrics/json
获取组件性能指标（JSON格式）。

**响应**: 详见各组件具体指标结构

### 信息接口

#### GET /info
获取组件基本信息。

**响应**:
```json
{
    "component": "yao-portal",
    "version": "1.0.0",
    "buildTime": "2024-01-15T08:00:00Z",
    "gitCommit": "abc123def456",
    "goVersion": "go1.21.5",
    "nodeId": "portal-node-1",
    "nodeName": "Portal Primary",
    "startTime": "2024-01-14T10:00:00Z",
    "uptime": "24h30m45s",
    "pid": 12345,
    "config": {
        "host": "0.0.0.0",
        "port": 8545,
        "chainId": 1337
    }
}
```

#### GET /status
获取组件运行状态。

**响应**:
```json
{
    "status": "running|starting|stopping|stopped|error",
    "message": "Component is running normally",
    "lastStateChange": "2024-01-14T10:00:00Z",
    "connections": {
        "database": "connected",
        "messageQueue": "connected", 
        "cache": "connected"
    },
    "performance": {
        "requestsPerSecond": 125.6,
        "averageResponseTime": 45,
        "errorRate": 0.01
    }
}
```

## YaoPortal 专用接口

### 交易池管理

#### GET /txpool/status
获取交易池状态。

**响应**:
```json
{
    "pending": 45,
    "queued": 12,
    "capacity": 1000,
    "usage": "5.7%",
    "oldestTransaction": {
        "hash": "0xabc123...",
        "submittedAt": "2024-01-15T10:25:00Z",
        "age": "5m45s"
    }
}
```

#### GET /txpool/inspect
检查交易池详细信息。

**响应**:
```json
{
    "pending": {
        "0xa1b2c3...": {
            "nonce": 15,
            "transactions": 3,
            "oldest": "2024-01-15T10:25:00Z"
        }
    },
    "queued": {
        "0xd4e5f6...": {
            "nonce": 20,
            "transactions": 1,
            "oldest": "2024-01-15T10:28:00Z"
        }
    }
}
```

#### DELETE /txpool/transaction/{hash}
从交易池中移除指定交易。

**参数**:
- `hash`: 交易哈希

**响应**:
```json
{
    "success": true,
    "message": "Transaction removed from pool"
}
```

### 连接管理

#### GET /connections
获取当前连接信息。

**响应**:
```json
{
    "total": 25,
    "active": 18,
    "websockets": 7,
    "http": 11,
    "byRemoteIP": {
        "192.168.1.100": 5,
        "192.168.1.101": 3
    },
    "rateLimitStats": {
        "requestsInWindow": 1250,
        "requestsLimit": 6000,
        "windowDuration": "1m"
    }
}
```

## YaoChronos 专用接口

### 领导选举管理

#### GET /election/status
获取领导选举状态。

**响应**:
```json
{
    "isLeader": true,
    "leaderID": "chronos-node-1",
    "leaderEndpoint": "http://chronos-1:9091",
    "leaderTerm": 15,
    "lastElection": "2024-01-14T12:00:00Z",
    "participants": [
        {
            "nodeId": "chronos-node-1",
            "status": "leader",
            "lastSeen": "2024-01-15T10:30:00Z"
        },
        {
            "nodeId": "chronos-node-2", 
            "status": "follower",
            "lastSeen": "2024-01-15T10:29:55Z"
        }
    ]
}
```

#### POST /election/step-down
主动放弃领导权（仅限当前leader）。

**响应**:
```json
{
    "success": true,
    "message": "Successfully stepped down from leadership"
}
```

### 区块生产管理

#### GET /sequencer/status
获取序列器状态。

**响应**:
```json
{
    "status": "active|paused|error",
    "currentBlock": 12345,
    "pendingTransactions": 156,
    "averageBlockTime": "2.3s",
    "blocksPerSecond": 0.43,
    "lastBlockProduced": "2024-01-15T10:29:58Z",
    "orderingStrategy": "fifo"
}
```

#### POST /sequencer/pause
暂停区块生产。

**响应**:
```json
{
    "success": true,
    "message": "Block production paused"
}
```

#### POST /sequencer/resume
恢复区块生产。

**响应**:
```json
{
    "success": true,
    "message": "Block production resumed"
}
```

## YaoVulcan 专用接口

### 执行管理

#### GET /executor/status
获取执行器状态。

**响应**:
```json
{
    "status": "running|paused|error",
    "currentProcessing": {
        "blockNumber": 12345,
        "startTime": "2024-01-15T10:29:45Z",
        "progress": "75%",
        "transactionsExecuted": 18,
        "transactionsTotal": 24
    },
    "lastProcessedBlock": 12344,
    "averageExecutionTime": "1.2s",
    "executionQueue": {
        "pending": 3,
        "processing": 1
    }
}
```

### 缓存管理

#### GET /cache/stats
获取缓存统计信息。

**响应**:
```json
{
    "l1Cache": {
        "hits": 145673,
        "misses": 12456,
        "hitRatio": 0.921,
        "size": 67108864,
        "capacity": 134217728,
        "usage": "50%"
    },
    "l2Cache": {
        "hits": 11234,
        "misses": 1222,
        "hitRatio": 0.902,
        "connections": 5,
        "commandsProcessed": 156789
    },
    "l3Storage": {
        "queries": 1089,
        "averageQueryTime": "15ms",
        "slowQueries": 23,
        "connections": 3
    }
}
```

#### POST /cache/clear
清空缓存。

**参数**:
```json
{
    "level": "l1|l2|all",
    "keys": ["optional", "specific", "keys"]
}
```

**响应**:
```json
{
    "success": true,
    "message": "Cache cleared successfully",
    "keysCleared": 1543
}
```

#### POST /cache/warmup
执行缓存预热。

**参数**:
```json
{
    "accounts": ["0xa1b2c3...", "0xd4e5f6..."],
    "contracts": ["0x123456..."],
    "blockRange": {
        "from": 12300,
        "to": 12350
    }
}
```

**响应**:
```json
{
    "success": true,
    "message": "Cache warmup started",
    "jobId": "warmup-job-123",
    "estimatedDuration": "2m30s"
}
```

## YaoArchive 专用接口

### 持久化管理

#### GET /persister/status
获取持久化器状态。

**响应**:
```json
{
    "status": "running|paused|error",
    "currentCommitting": {
        "blockNumber": 12345,
        "startTime": "2024-01-15T10:29:58Z",
        "stage": "writing_block|writing_state|committing"
    },
    "lastCommittedBlock": 12344,
    "commitQueue": {
        "pending": 2,
        "processing": 1,
        "failed": 0
    },
    "performance": {
        "averageCommitTime": "850ms",
        "transactionThroughput": "2.3/s",
        "storageUsage": "1.2TB"
    }
}
```

### 数据库管理

#### GET /storage/stats
获取存储统计信息。

**响应**:
```json
{
    "database": {
        "connections": 8,
        "activeQueries": 3,
        "slowQueries": 12,
        "averageQueryTime": "25ms",
        "storage": {
            "totalSize": "1.2TB",
            "freeSpace": "800GB",
            "usage": "60%"
        }
    },
    "tables": [
        {
            "name": "blocks",
            "rows": 12345,
            "size": "156MB",
            "lastUpdated": "2024-01-15T10:29:58Z"
        },
        {
            "name": "transactions", 
            "rows": 456789,
            "size": "2.3GB",
            "lastUpdated": "2024-01-15T10:29:58Z"
        }
    ]
}
```

#### POST /storage/optimize
执行存储优化。

**参数**:
```json
{
    "tables": ["blocks", "transactions"],
    "operations": ["analyze", "vacuum", "reindex"]
}
```

**响应**:
```json
{
    "success": true,
    "message": "Storage optimization started",
    "jobId": "optimize-job-456",
    "estimatedDuration": "15m"
}
```

## 配置管理

### GET /config
获取当前配置。

**响应**:
```json
{
    "current": {
        "host": "0.0.0.0",
        "port": 8545,
        "logLevel": "info",
        "maxConnections": 100
    },
    "defaults": {
        "host": "127.0.0.1",
        "port": 8545,
        "logLevel": "info",
        "maxConnections": 50
    },
    "lastReloaded": "2024-01-14T10:00:00Z"
}
```

### PUT /config
更新配置（热重载）。

**请求**:
```json
{
    "logLevel": "debug",
    "maxConnections": 200,
    "restart": false
}
```

**响应**:
```json
{
    "success": true,
    "message": "Configuration updated successfully",
    "changedFields": ["logLevel", "maxConnections"],
    "restartRequired": false
}
```

### POST /config/reload
重新加载配置文件。

**响应**:
```json
{
    "success": true,
    "message": "Configuration reloaded from file",
    "changedFields": ["port", "logLevel"],
    "restartRequired": true
}
```

## 日志管理

### GET /logs/level
获取当前日志级别。

**响应**:
```json
{
    "level": "info",
    "availableLevels": ["debug", "info", "warn", "error", "fatal"]
}
```

### PUT /logs/level
设置日志级别。

**请求**:
```json
{
    "level": "debug"
}
```

**响应**:
```json
{
    "success": true,
    "message": "Log level changed to debug",
    "previousLevel": "info"
}
```

### GET /logs/recent
获取最近的日志条目。

**参数**:
- `limit`: 数量限制 (default: 100, max: 1000)
- `level`: 日志级别过滤
- `since`: 时间过滤 (ISO 8601格式)

**响应**:
```json
{
    "logs": [
        {
            "timestamp": "2024-01-15T10:30:45Z",
            "level": "info",
            "component": "server",
            "message": "Request processed successfully",
            "fields": {
                "method": "eth_getBalance",
                "duration": "25ms",
                "address": "0xa1b2c3..."
            }
        }
    ],
    "total": 1,
    "hasMore": false
}
```

## 运维操作

### POST /shutdown
优雅关闭组件。

**参数**:
```json
{
    "gracefulTimeout": "30s",
    "forceKillTimeout": "60s"
}
```

**响应**:
```json
{
    "success": true,
    "message": "Shutdown initiated",
    "estimatedDuration": "30s"
}
```

### POST /restart
重启组件。

**参数**:
```json
{
    "graceful": true,
    "timeout": "60s"
}
```

**响应**:
```json
{
    "success": true,
    "message": "Restart initiated"
}
```

## 错误响应格式

所有管理API的错误响应都使用统一格式：

```json
{
    "error": {
        "code": "INVALID_PARAMETER",
        "message": "Invalid parameter value",
        "details": "Parameter 'level' must be one of: debug, info, warn, error, fatal",
        "timestamp": "2024-01-15T10:30:45Z",
        "requestId": "req-123456"
    }
}
```

### 错误代码

| 错误代码 | HTTP状态 | 说明 |
|---------|---------|------|
| `UNAUTHORIZED` | 401 | 未授权访问 |
| `FORBIDDEN` | 403 | 权限不足 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `INVALID_PARAMETER` | 400 | 参数无效 |
| `OPERATION_FAILED` | 500 | 操作失败 |
| `SERVICE_UNAVAILABLE` | 503 | 服务不可用 |
| `TIMEOUT` | 408 | 请求超时 |

## 安全考虑

1. **访问控制**: 管理API需要有效的Bearer Token
2. **IP白名单**: 只允许特定IP访问管理接口
3. **操作审计**: 所有管理操作都会记录到审计日志
4. **权限分级**: 不同操作需要不同级别的权限

## 使用示例

### 获取集群整体状态

```bash
#!/bin/bash

# 检查所有组件健康状态
for component in portal chronos vulcan archive; do
    echo "Checking $component..."
    curl -H "Authorization: Bearer $MGMT_TOKEN" \
         "http://${component}:909X/api/v1/health" | jq .
done
```

### 监控性能指标

```bash
# 获取Portal性能指标
curl -H "Authorization: Bearer $MGMT_TOKEN" \
     "http://portal:9090/api/v1/metrics/json" | \
     jq '.totalRequests, .requestsPerSecond, .averageResponseTime'
```

### 维护操作

```bash
# 清理Vulcan缓存
curl -X POST \
     -H "Authorization: Bearer $MGMT_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"level": "l1"}' \
     "http://vulcan:9092/api/v1/cache/clear"
```

---

更多详细信息请参考各组件的具体文档和 [YaoVerse 运维指南](../operations/)。 