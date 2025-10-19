# Axiom Backend V3 API 完整文檔

> **版本**: 3.0.0  
> **日期**: 2025-10-16  
> **Base URL**: `http://localhost:3001`

---

## 📋 目錄

1. [系統健康](#系統健康)
2. [Prometheus APIs](#prometheus-apis)
3. [Loki APIs](#loki-apis)
4. [Quantum APIs](#quantum-apis)
5. [Nginx APIs](#nginx-apis)
6. [Windows Logs APIs](#windows-logs-apis)
7. [認證](#認證)
8. [錯誤處理](#錯誤處理)

---

## 系統健康

### GET /health

系統健康檢查

**Response**:
```json
{
  "status": "healthy",
  "service": "axiom-backend-v3",
  "version": "3.0.0",
  "time": "2025-10-16T10:00:00Z"
}
```

---

## Prometheus APIs

Base path: `/api/v2/prometheus`

### POST /api/v2/prometheus/query

執行 PromQL 即時查詢

**Request Body**:
```json
{
  "query": "up{job='node-exporter'}",
  "time": "2025-10-16T10:00:00Z" // 可選
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "success",
    "data": {
      "result_type": "vector",
      "result": [
        {
          "metric": {"job": "node-exporter"},
          "value": [1697452800, "1"]
        }
      ]
    },
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/prometheus/query-range

執行 PromQL 範圍查詢

**Request Body**:
```json
{
  "query": "rate(http_requests_total[5m])",
  "start": "2025-10-16T09:00:00Z",
  "end": "2025-10-16T10:00:00Z",
  "step": "15s"
}
```

### GET /api/v2/prometheus/rules

獲取所有告警規則

**Response**:
```json
{
  "success": true,
  "data": {
    "groups": [
      {
        "name": "example-rules",
        "file": "/etc/prometheus/rules.yml",
        "rules": [
          {
            "name": "HighErrorRate",
            "query": "rate(errors_total[5m]) > 0.05",
            "state": "firing"
          }
        ]
      }
    ]
  }
}
```

### GET /api/v2/prometheus/targets

獲取所有抓取目標

### GET /api/v2/prometheus/health

Prometheus 服務健康檢查

### GET /api/v2/prometheus/status

獲取 Prometheus 服務狀態

---

## Loki APIs

Base path: `/api/v2/loki`

### GET /api/v2/loki/query

查詢 Loki 日誌

**Query Parameters**:
- `query` (required): LogQL 查詢語句
- `limit` (optional): 限制返回數量
- `start` (optional): 開始時間 (RFC3339)
- `end` (optional): 結束時間 (RFC3339)

**Example**:
```
GET /api/v2/loki/query?query={job="varlogs"}|="error"&limit=100
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "success",
    "data": {
      "resultType": "streams",
      "result": [...]
    }
  }
}
```

### GET /api/v2/loki/labels

獲取所有可用標籤

**Response**:
```json
{
  "success": true,
  "data": {
    "labels": ["job", "filename", "level"]
  }
}
```

### GET /api/v2/loki/labels/{label}/values

獲取指定標籤的所有值

**Example**:
```
GET /api/v2/loki/labels/job/values
```

### GET /api/v2/loki/health

Loki 服務健康檢查

---

## Quantum APIs

Base path: `/api/v2/quantum`

### POST /api/v2/quantum/qkd/generate

生成量子密鑰分發

**Request Body**:
```json
{
  "key_length": 256,
  "backend": "simulator", // "simulator" or "ibm_quantum"
  "shots": 1024
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "job_id": "qkd-abc12345",
    "key": "base64_encoded_key",
    "key_length": 256,
    "status": "completed",
    "submitted_at": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/quantum/qsvm/classify

執行量子支持向量機分類

**Request Body**:
```json
{
  "features": [0.1, 0.2, 0.3, 0.4],
  "backend": "simulator",
  "feature_dim": 4,
  "shots": 1024
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "job_id": "qsvm-def67890",
    "prediction": 1,
    "probability": 0.87,
    "confidence": 0.95,
    "status": "completed",
    "submitted_at": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/quantum/zerotrust/predict

執行 Zero Trust 安全預測

**Request Body**:
```json
{
  "user_id": "user123",
  "ip_address": "192.168.1.100",
  "device_id": "device_456",
  "features": {
    "login_time": "10:00",
    "location": "taipei"
  },
  "use_quantum": true
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "trust_score": 0.85,
    "risk_level": "low",
    "decision": "allow",
    "confidence": 0.92,
    "factors": {
      "user_behavior": 0.9,
      "device_trust": 0.8,
      "location_risk": 0.85
    },
    "used_quantum": true,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/quantum/jobs

列出所有量子作業

**Query Parameters**:
- `type`: 作業類型 (qkd, qsvm, qaoa, etc.)
- `status`: 作業狀態 (pending, running, completed, failed)
- `page`: 頁碼
- `page_size`: 每頁數量

### GET /api/v2/quantum/jobs/{jobId}

獲取單個量子作業詳情

### GET /api/v2/quantum/stats

獲取量子作業統計

**Response**:
```json
{
  "success": true,
  "data": {
    "total_jobs": 1523,
    "completed_jobs": 1420,
    "failed_jobs": 53,
    "running_jobs": 50,
    "success_rate": 0.93,
    "jobs_by_type": {
      "qkd": 520,
      "qsvm": 380,
      "qaoa": 250,
      "zerotrust": 373
    }
  }
}
```

### GET /api/v2/quantum/health

Quantum 服務健康檢查

---

## Nginx APIs

Base path: `/api/v2/nginx`

### GET /api/v2/nginx/status

獲取 Nginx 狀態

**Response**:
```json
{
  "success": true,
  "data": {
    "name": "nginx",
    "status": "healthy",
    "active_connections": 42,
    "accepts": 15230,
    "handled": 15230,
    "requests": 30460,
    "reading": 0,
    "writing": 3,
    "waiting": 39,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/nginx/config

獲取 Nginx 配置

**Response**:
```json
{
  "success": true,
  "data": {
    "config": "user nginx;\nworker_processes auto;\n...",
    "config_path": "/etc/nginx/nginx.conf",
    "last_modified": "2025-10-16T09:00:00Z",
    "size": 4096,
    "valid": true
  }
}
```

### PUT /api/v2/nginx/config

更新 Nginx 配置

**Request Body**:
```json
{
  "config": "user nginx;\nworker_processes auto;\n...",
  "validate": true,
  "backup": true
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "config_path": "/etc/nginx/nginx.conf",
    "last_modified": "2025-10-16T10:05:00Z",
    "valid": true
  }
}
```

### POST /api/v2/nginx/reload

重載 Nginx 配置

**Response**:
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Nginx reloaded successfully",
    "duration": 125,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

---

## Windows Logs APIs

Base path: `/api/v2/logs/windows`

### POST /api/v2/logs/windows/batch

批量接收 Windows 日誌

**Request Body**:
```json
{
  "agent_id": "agent-001",
  "computer": "WIN-SERVER-01",
  "logs": [
    {
      "log_type": "Security",
      "source": "Microsoft-Windows-Security-Auditing",
      "event_id": 4624,
      "level": "Information",
      "message": "An account was successfully logged on",
      "time_created": "2025-10-16T10:00:00Z",
      "user_id": "S-1-5-21-...",
      "process_id": 1234,
      "thread_id": 5678,
      "metadata": {}
    }
  ]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "received_count": 1,
    "saved_count": 1,
    "failed_count": 0,
    "errors": [],
    "timestamp": "2025-10-16T10:00:00Z",
    "message": "Logs processed successfully"
  }
}
```

### GET /api/v2/logs/windows

查詢 Windows 日誌

**Query Parameters**:
- `agent_id`: Agent ID 過濾
- `log_type`: 日誌類型 (System, Security, Application, Setup)
- `level`: 日誌級別 (Critical, Error, Warning, Information)
- `event_id`: 事件 ID
- `keyword`: 關鍵字搜索
- `start_time`: 開始時間 (RFC3339)
- `end_time`: 結束時間 (RFC3339)
- `page`: 頁碼 (默認 1)
- `page_size`: 每頁數量 (默認 50, 最大 1000)

**Response**:
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 12345,
        "agent_id": "agent-001",
        "log_type": "Security",
        "source": "Microsoft-Windows-Security-Auditing",
        "event_id": 4624,
        "level": "Information",
        "message": "An account was successfully logged on",
        "time_created": "2025-10-16T10:00:00Z",
        "received_at": "2025-10-16T10:00:05Z",
        "computer": "WIN-SERVER-01"
      }
    ],
    "total": 1523,
    "page": 1,
    "page_size": 50,
    "total_pages": 31,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/logs/windows/stats

獲取 Windows 日誌統計

**Query Parameters**:
- `time_range`: 時間範圍 (1h, 24h, 7d, 30d)

**Response**:
```json
{
  "success": true,
  "data": {
    "total_logs": 15230,
    "logs_by_type": {
      "System": 5420,
      "Security": 6830,
      "Application": 2980
    },
    "logs_by_level": {
      "Critical": 12,
      "Error": 145,
      "Warning": 523,
      "Information": 14550
    },
    "critical_count": 12,
    "error_count": 145,
    "warning_count": 523,
    "time_range": "24h"
  }
}
```

---

## 認證

### Header 認證

所有 API 請求應包含認證 Header：

```http
Authorization: Bearer <JWT_TOKEN>
```

或使用 API Key：

```http
X-API-Key: <API_KEY>
```

---

## 錯誤處理

### 錯誤響應格式

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request",
    "details": "field 'query' is required"
  }
}
```

### 錯誤碼列表

| 錯誤碼 | HTTP 狀態 | 說明 |
|--------|----------|------|
| `INTERNAL_ERROR` | 500 | 內部服務器錯誤 |
| `VALIDATION_ERROR` | 400 | 請求驗證失敗 |
| `NOT_FOUND` | 404 | 資源不存在 |
| `UNAUTHORIZED` | 401 | 未授權 |
| `FORBIDDEN` | 403 | 禁止訪問 |
| `SERVICE_UNAVAILABLE` | 503 | 服務不可用 |
| `TIMEOUT` | 408 | 請求超時 |

---

## 速率限制

- 默認限制：100 requests/minute
- Burst 限制：200 requests
- Header: `X-RateLimit-Remaining`, `X-RateLimit-Reset`

---

## 最佳實踐

### 1. 使用分頁

對於返回列表的 API，始終使用分頁：

```
GET /api/v2/quantum/jobs?page=1&page_size=50
```

### 2. 處理錯誤

始終檢查 `success` 欄位：

```typescript
const response = await fetch('/api/v2/prometheus/query', {
  method: 'POST',
  body: JSON.stringify({ query: 'up' }),
});

const data = await response.json();
if (!data.success) {
  console.error(data.error);
}
```

### 3. 使用時間範圍

對於時間序列查詢，明確指定時間範圍以提升性能。

---

## 版本控制

API 版本通過 URL 路徑指定：

- V2: `/api/v2/*` (當前穩定版)
- V3: `/api/v3/*` (未來版本)
- Experimental: `/api/v2/experimental/*`

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16



> **版本**: 3.0.0  
> **日期**: 2025-10-16  
> **Base URL**: `http://localhost:3001`

---

## 📋 目錄

1. [系統健康](#系統健康)
2. [Prometheus APIs](#prometheus-apis)
3. [Loki APIs](#loki-apis)
4. [Quantum APIs](#quantum-apis)
5. [Nginx APIs](#nginx-apis)
6. [Windows Logs APIs](#windows-logs-apis)
7. [認證](#認證)
8. [錯誤處理](#錯誤處理)

---

## 系統健康

### GET /health

系統健康檢查

**Response**:
```json
{
  "status": "healthy",
  "service": "axiom-backend-v3",
  "version": "3.0.0",
  "time": "2025-10-16T10:00:00Z"
}
```

---

## Prometheus APIs

Base path: `/api/v2/prometheus`

### POST /api/v2/prometheus/query

執行 PromQL 即時查詢

**Request Body**:
```json
{
  "query": "up{job='node-exporter'}",
  "time": "2025-10-16T10:00:00Z" // 可選
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "success",
    "data": {
      "result_type": "vector",
      "result": [
        {
          "metric": {"job": "node-exporter"},
          "value": [1697452800, "1"]
        }
      ]
    },
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/prometheus/query-range

執行 PromQL 範圍查詢

**Request Body**:
```json
{
  "query": "rate(http_requests_total[5m])",
  "start": "2025-10-16T09:00:00Z",
  "end": "2025-10-16T10:00:00Z",
  "step": "15s"
}
```

### GET /api/v2/prometheus/rules

獲取所有告警規則

**Response**:
```json
{
  "success": true,
  "data": {
    "groups": [
      {
        "name": "example-rules",
        "file": "/etc/prometheus/rules.yml",
        "rules": [
          {
            "name": "HighErrorRate",
            "query": "rate(errors_total[5m]) > 0.05",
            "state": "firing"
          }
        ]
      }
    ]
  }
}
```

### GET /api/v2/prometheus/targets

獲取所有抓取目標

### GET /api/v2/prometheus/health

Prometheus 服務健康檢查

### GET /api/v2/prometheus/status

獲取 Prometheus 服務狀態

---

## Loki APIs

Base path: `/api/v2/loki`

### GET /api/v2/loki/query

查詢 Loki 日誌

**Query Parameters**:
- `query` (required): LogQL 查詢語句
- `limit` (optional): 限制返回數量
- `start` (optional): 開始時間 (RFC3339)
- `end` (optional): 結束時間 (RFC3339)

**Example**:
```
GET /api/v2/loki/query?query={job="varlogs"}|="error"&limit=100
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "success",
    "data": {
      "resultType": "streams",
      "result": [...]
    }
  }
}
```

### GET /api/v2/loki/labels

獲取所有可用標籤

**Response**:
```json
{
  "success": true,
  "data": {
    "labels": ["job", "filename", "level"]
  }
}
```

### GET /api/v2/loki/labels/{label}/values

獲取指定標籤的所有值

**Example**:
```
GET /api/v2/loki/labels/job/values
```

### GET /api/v2/loki/health

Loki 服務健康檢查

---

## Quantum APIs

Base path: `/api/v2/quantum`

### POST /api/v2/quantum/qkd/generate

生成量子密鑰分發

**Request Body**:
```json
{
  "key_length": 256,
  "backend": "simulator", // "simulator" or "ibm_quantum"
  "shots": 1024
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "job_id": "qkd-abc12345",
    "key": "base64_encoded_key",
    "key_length": 256,
    "status": "completed",
    "submitted_at": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/quantum/qsvm/classify

執行量子支持向量機分類

**Request Body**:
```json
{
  "features": [0.1, 0.2, 0.3, 0.4],
  "backend": "simulator",
  "feature_dim": 4,
  "shots": 1024
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "job_id": "qsvm-def67890",
    "prediction": 1,
    "probability": 0.87,
    "confidence": 0.95,
    "status": "completed",
    "submitted_at": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/quantum/zerotrust/predict

執行 Zero Trust 安全預測

**Request Body**:
```json
{
  "user_id": "user123",
  "ip_address": "192.168.1.100",
  "device_id": "device_456",
  "features": {
    "login_time": "10:00",
    "location": "taipei"
  },
  "use_quantum": true
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "trust_score": 0.85,
    "risk_level": "low",
    "decision": "allow",
    "confidence": 0.92,
    "factors": {
      "user_behavior": 0.9,
      "device_trust": 0.8,
      "location_risk": 0.85
    },
    "used_quantum": true,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/quantum/jobs

列出所有量子作業

**Query Parameters**:
- `type`: 作業類型 (qkd, qsvm, qaoa, etc.)
- `status`: 作業狀態 (pending, running, completed, failed)
- `page`: 頁碼
- `page_size`: 每頁數量

### GET /api/v2/quantum/jobs/{jobId}

獲取單個量子作業詳情

### GET /api/v2/quantum/stats

獲取量子作業統計

**Response**:
```json
{
  "success": true,
  "data": {
    "total_jobs": 1523,
    "completed_jobs": 1420,
    "failed_jobs": 53,
    "running_jobs": 50,
    "success_rate": 0.93,
    "jobs_by_type": {
      "qkd": 520,
      "qsvm": 380,
      "qaoa": 250,
      "zerotrust": 373
    }
  }
}
```

### GET /api/v2/quantum/health

Quantum 服務健康檢查

---

## Nginx APIs

Base path: `/api/v2/nginx`

### GET /api/v2/nginx/status

獲取 Nginx 狀態

**Response**:
```json
{
  "success": true,
  "data": {
    "name": "nginx",
    "status": "healthy",
    "active_connections": 42,
    "accepts": 15230,
    "handled": 15230,
    "requests": 30460,
    "reading": 0,
    "writing": 3,
    "waiting": 39,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/nginx/config

獲取 Nginx 配置

**Response**:
```json
{
  "success": true,
  "data": {
    "config": "user nginx;\nworker_processes auto;\n...",
    "config_path": "/etc/nginx/nginx.conf",
    "last_modified": "2025-10-16T09:00:00Z",
    "size": 4096,
    "valid": true
  }
}
```

### PUT /api/v2/nginx/config

更新 Nginx 配置

**Request Body**:
```json
{
  "config": "user nginx;\nworker_processes auto;\n...",
  "validate": true,
  "backup": true
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "config_path": "/etc/nginx/nginx.conf",
    "last_modified": "2025-10-16T10:05:00Z",
    "valid": true
  }
}
```

### POST /api/v2/nginx/reload

重載 Nginx 配置

**Response**:
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Nginx reloaded successfully",
    "duration": 125,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

---

## Windows Logs APIs

Base path: `/api/v2/logs/windows`

### POST /api/v2/logs/windows/batch

批量接收 Windows 日誌

**Request Body**:
```json
{
  "agent_id": "agent-001",
  "computer": "WIN-SERVER-01",
  "logs": [
    {
      "log_type": "Security",
      "source": "Microsoft-Windows-Security-Auditing",
      "event_id": 4624,
      "level": "Information",
      "message": "An account was successfully logged on",
      "time_created": "2025-10-16T10:00:00Z",
      "user_id": "S-1-5-21-...",
      "process_id": 1234,
      "thread_id": 5678,
      "metadata": {}
    }
  ]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "received_count": 1,
    "saved_count": 1,
    "failed_count": 0,
    "errors": [],
    "timestamp": "2025-10-16T10:00:00Z",
    "message": "Logs processed successfully"
  }
}
```

### GET /api/v2/logs/windows

查詢 Windows 日誌

**Query Parameters**:
- `agent_id`: Agent ID 過濾
- `log_type`: 日誌類型 (System, Security, Application, Setup)
- `level`: 日誌級別 (Critical, Error, Warning, Information)
- `event_id`: 事件 ID
- `keyword`: 關鍵字搜索
- `start_time`: 開始時間 (RFC3339)
- `end_time`: 結束時間 (RFC3339)
- `page`: 頁碼 (默認 1)
- `page_size`: 每頁數量 (默認 50, 最大 1000)

**Response**:
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 12345,
        "agent_id": "agent-001",
        "log_type": "Security",
        "source": "Microsoft-Windows-Security-Auditing",
        "event_id": 4624,
        "level": "Information",
        "message": "An account was successfully logged on",
        "time_created": "2025-10-16T10:00:00Z",
        "received_at": "2025-10-16T10:00:05Z",
        "computer": "WIN-SERVER-01"
      }
    ],
    "total": 1523,
    "page": 1,
    "page_size": 50,
    "total_pages": 31,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/logs/windows/stats

獲取 Windows 日誌統計

**Query Parameters**:
- `time_range`: 時間範圍 (1h, 24h, 7d, 30d)

**Response**:
```json
{
  "success": true,
  "data": {
    "total_logs": 15230,
    "logs_by_type": {
      "System": 5420,
      "Security": 6830,
      "Application": 2980
    },
    "logs_by_level": {
      "Critical": 12,
      "Error": 145,
      "Warning": 523,
      "Information": 14550
    },
    "critical_count": 12,
    "error_count": 145,
    "warning_count": 523,
    "time_range": "24h"
  }
}
```

---

## 認證

### Header 認證

所有 API 請求應包含認證 Header：

```http
Authorization: Bearer <JWT_TOKEN>
```

或使用 API Key：

```http
X-API-Key: <API_KEY>
```

---

## 錯誤處理

### 錯誤響應格式

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request",
    "details": "field 'query' is required"
  }
}
```

### 錯誤碼列表

| 錯誤碼 | HTTP 狀態 | 說明 |
|--------|----------|------|
| `INTERNAL_ERROR` | 500 | 內部服務器錯誤 |
| `VALIDATION_ERROR` | 400 | 請求驗證失敗 |
| `NOT_FOUND` | 404 | 資源不存在 |
| `UNAUTHORIZED` | 401 | 未授權 |
| `FORBIDDEN` | 403 | 禁止訪問 |
| `SERVICE_UNAVAILABLE` | 503 | 服務不可用 |
| `TIMEOUT` | 408 | 請求超時 |

---

## 速率限制

- 默認限制：100 requests/minute
- Burst 限制：200 requests
- Header: `X-RateLimit-Remaining`, `X-RateLimit-Reset`

---

## 最佳實踐

### 1. 使用分頁

對於返回列表的 API，始終使用分頁：

```
GET /api/v2/quantum/jobs?page=1&page_size=50
```

### 2. 處理錯誤

始終檢查 `success` 欄位：

```typescript
const response = await fetch('/api/v2/prometheus/query', {
  method: 'POST',
  body: JSON.stringify({ query: 'up' }),
});

const data = await response.json();
if (!data.success) {
  console.error(data.error);
}
```

### 3. 使用時間範圍

對於時間序列查詢，明確指定時間範圍以提升性能。

---

## 版本控制

API 版本通過 URL 路徑指定：

- V2: `/api/v2/*` (當前穩定版)
- V3: `/api/v3/*` (未來版本)
- Experimental: `/api/v2/experimental/*`

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16


> **版本**: 3.0.0  
> **日期**: 2025-10-16  
> **Base URL**: `http://localhost:3001`

---

## 📋 目錄

1. [系統健康](#系統健康)
2. [Prometheus APIs](#prometheus-apis)
3. [Loki APIs](#loki-apis)
4. [Quantum APIs](#quantum-apis)
5. [Nginx APIs](#nginx-apis)
6. [Windows Logs APIs](#windows-logs-apis)
7. [認證](#認證)
8. [錯誤處理](#錯誤處理)

---

## 系統健康

### GET /health

系統健康檢查

**Response**:
```json
{
  "status": "healthy",
  "service": "axiom-backend-v3",
  "version": "3.0.0",
  "time": "2025-10-16T10:00:00Z"
}
```

---

## Prometheus APIs

Base path: `/api/v2/prometheus`

### POST /api/v2/prometheus/query

執行 PromQL 即時查詢

**Request Body**:
```json
{
  "query": "up{job='node-exporter'}",
  "time": "2025-10-16T10:00:00Z" // 可選
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "success",
    "data": {
      "result_type": "vector",
      "result": [
        {
          "metric": {"job": "node-exporter"},
          "value": [1697452800, "1"]
        }
      ]
    },
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/prometheus/query-range

執行 PromQL 範圍查詢

**Request Body**:
```json
{
  "query": "rate(http_requests_total[5m])",
  "start": "2025-10-16T09:00:00Z",
  "end": "2025-10-16T10:00:00Z",
  "step": "15s"
}
```

### GET /api/v2/prometheus/rules

獲取所有告警規則

**Response**:
```json
{
  "success": true,
  "data": {
    "groups": [
      {
        "name": "example-rules",
        "file": "/etc/prometheus/rules.yml",
        "rules": [
          {
            "name": "HighErrorRate",
            "query": "rate(errors_total[5m]) > 0.05",
            "state": "firing"
          }
        ]
      }
    ]
  }
}
```

### GET /api/v2/prometheus/targets

獲取所有抓取目標

### GET /api/v2/prometheus/health

Prometheus 服務健康檢查

### GET /api/v2/prometheus/status

獲取 Prometheus 服務狀態

---

## Loki APIs

Base path: `/api/v2/loki`

### GET /api/v2/loki/query

查詢 Loki 日誌

**Query Parameters**:
- `query` (required): LogQL 查詢語句
- `limit` (optional): 限制返回數量
- `start` (optional): 開始時間 (RFC3339)
- `end` (optional): 結束時間 (RFC3339)

**Example**:
```
GET /api/v2/loki/query?query={job="varlogs"}|="error"&limit=100
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "success",
    "data": {
      "resultType": "streams",
      "result": [...]
    }
  }
}
```

### GET /api/v2/loki/labels

獲取所有可用標籤

**Response**:
```json
{
  "success": true,
  "data": {
    "labels": ["job", "filename", "level"]
  }
}
```

### GET /api/v2/loki/labels/{label}/values

獲取指定標籤的所有值

**Example**:
```
GET /api/v2/loki/labels/job/values
```

### GET /api/v2/loki/health

Loki 服務健康檢查

---

## Quantum APIs

Base path: `/api/v2/quantum`

### POST /api/v2/quantum/qkd/generate

生成量子密鑰分發

**Request Body**:
```json
{
  "key_length": 256,
  "backend": "simulator", // "simulator" or "ibm_quantum"
  "shots": 1024
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "job_id": "qkd-abc12345",
    "key": "base64_encoded_key",
    "key_length": 256,
    "status": "completed",
    "submitted_at": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/quantum/qsvm/classify

執行量子支持向量機分類

**Request Body**:
```json
{
  "features": [0.1, 0.2, 0.3, 0.4],
  "backend": "simulator",
  "feature_dim": 4,
  "shots": 1024
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "job_id": "qsvm-def67890",
    "prediction": 1,
    "probability": 0.87,
    "confidence": 0.95,
    "status": "completed",
    "submitted_at": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/quantum/zerotrust/predict

執行 Zero Trust 安全預測

**Request Body**:
```json
{
  "user_id": "user123",
  "ip_address": "192.168.1.100",
  "device_id": "device_456",
  "features": {
    "login_time": "10:00",
    "location": "taipei"
  },
  "use_quantum": true
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "trust_score": 0.85,
    "risk_level": "low",
    "decision": "allow",
    "confidence": 0.92,
    "factors": {
      "user_behavior": 0.9,
      "device_trust": 0.8,
      "location_risk": 0.85
    },
    "used_quantum": true,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/quantum/jobs

列出所有量子作業

**Query Parameters**:
- `type`: 作業類型 (qkd, qsvm, qaoa, etc.)
- `status`: 作業狀態 (pending, running, completed, failed)
- `page`: 頁碼
- `page_size`: 每頁數量

### GET /api/v2/quantum/jobs/{jobId}

獲取單個量子作業詳情

### GET /api/v2/quantum/stats

獲取量子作業統計

**Response**:
```json
{
  "success": true,
  "data": {
    "total_jobs": 1523,
    "completed_jobs": 1420,
    "failed_jobs": 53,
    "running_jobs": 50,
    "success_rate": 0.93,
    "jobs_by_type": {
      "qkd": 520,
      "qsvm": 380,
      "qaoa": 250,
      "zerotrust": 373
    }
  }
}
```

### GET /api/v2/quantum/health

Quantum 服務健康檢查

---

## Nginx APIs

Base path: `/api/v2/nginx`

### GET /api/v2/nginx/status

獲取 Nginx 狀態

**Response**:
```json
{
  "success": true,
  "data": {
    "name": "nginx",
    "status": "healthy",
    "active_connections": 42,
    "accepts": 15230,
    "handled": 15230,
    "requests": 30460,
    "reading": 0,
    "writing": 3,
    "waiting": 39,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/nginx/config

獲取 Nginx 配置

**Response**:
```json
{
  "success": true,
  "data": {
    "config": "user nginx;\nworker_processes auto;\n...",
    "config_path": "/etc/nginx/nginx.conf",
    "last_modified": "2025-10-16T09:00:00Z",
    "size": 4096,
    "valid": true
  }
}
```

### PUT /api/v2/nginx/config

更新 Nginx 配置

**Request Body**:
```json
{
  "config": "user nginx;\nworker_processes auto;\n...",
  "validate": true,
  "backup": true
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "config_path": "/etc/nginx/nginx.conf",
    "last_modified": "2025-10-16T10:05:00Z",
    "valid": true
  }
}
```

### POST /api/v2/nginx/reload

重載 Nginx 配置

**Response**:
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Nginx reloaded successfully",
    "duration": 125,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

---

## Windows Logs APIs

Base path: `/api/v2/logs/windows`

### POST /api/v2/logs/windows/batch

批量接收 Windows 日誌

**Request Body**:
```json
{
  "agent_id": "agent-001",
  "computer": "WIN-SERVER-01",
  "logs": [
    {
      "log_type": "Security",
      "source": "Microsoft-Windows-Security-Auditing",
      "event_id": 4624,
      "level": "Information",
      "message": "An account was successfully logged on",
      "time_created": "2025-10-16T10:00:00Z",
      "user_id": "S-1-5-21-...",
      "process_id": 1234,
      "thread_id": 5678,
      "metadata": {}
    }
  ]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "received_count": 1,
    "saved_count": 1,
    "failed_count": 0,
    "errors": [],
    "timestamp": "2025-10-16T10:00:00Z",
    "message": "Logs processed successfully"
  }
}
```

### GET /api/v2/logs/windows

查詢 Windows 日誌

**Query Parameters**:
- `agent_id`: Agent ID 過濾
- `log_type`: 日誌類型 (System, Security, Application, Setup)
- `level`: 日誌級別 (Critical, Error, Warning, Information)
- `event_id`: 事件 ID
- `keyword`: 關鍵字搜索
- `start_time`: 開始時間 (RFC3339)
- `end_time`: 結束時間 (RFC3339)
- `page`: 頁碼 (默認 1)
- `page_size`: 每頁數量 (默認 50, 最大 1000)

**Response**:
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 12345,
        "agent_id": "agent-001",
        "log_type": "Security",
        "source": "Microsoft-Windows-Security-Auditing",
        "event_id": 4624,
        "level": "Information",
        "message": "An account was successfully logged on",
        "time_created": "2025-10-16T10:00:00Z",
        "received_at": "2025-10-16T10:00:05Z",
        "computer": "WIN-SERVER-01"
      }
    ],
    "total": 1523,
    "page": 1,
    "page_size": 50,
    "total_pages": 31,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/logs/windows/stats

獲取 Windows 日誌統計

**Query Parameters**:
- `time_range`: 時間範圍 (1h, 24h, 7d, 30d)

**Response**:
```json
{
  "success": true,
  "data": {
    "total_logs": 15230,
    "logs_by_type": {
      "System": 5420,
      "Security": 6830,
      "Application": 2980
    },
    "logs_by_level": {
      "Critical": 12,
      "Error": 145,
      "Warning": 523,
      "Information": 14550
    },
    "critical_count": 12,
    "error_count": 145,
    "warning_count": 523,
    "time_range": "24h"
  }
}
```

---

## 認證

### Header 認證

所有 API 請求應包含認證 Header：

```http
Authorization: Bearer <JWT_TOKEN>
```

或使用 API Key：

```http
X-API-Key: <API_KEY>
```

---

## 錯誤處理

### 錯誤響應格式

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request",
    "details": "field 'query' is required"
  }
}
```

### 錯誤碼列表

| 錯誤碼 | HTTP 狀態 | 說明 |
|--------|----------|------|
| `INTERNAL_ERROR` | 500 | 內部服務器錯誤 |
| `VALIDATION_ERROR` | 400 | 請求驗證失敗 |
| `NOT_FOUND` | 404 | 資源不存在 |
| `UNAUTHORIZED` | 401 | 未授權 |
| `FORBIDDEN` | 403 | 禁止訪問 |
| `SERVICE_UNAVAILABLE` | 503 | 服務不可用 |
| `TIMEOUT` | 408 | 請求超時 |

---

## 速率限制

- 默認限制：100 requests/minute
- Burst 限制：200 requests
- Header: `X-RateLimit-Remaining`, `X-RateLimit-Reset`

---

## 最佳實踐

### 1. 使用分頁

對於返回列表的 API，始終使用分頁：

```
GET /api/v2/quantum/jobs?page=1&page_size=50
```

### 2. 處理錯誤

始終檢查 `success` 欄位：

```typescript
const response = await fetch('/api/v2/prometheus/query', {
  method: 'POST',
  body: JSON.stringify({ query: 'up' }),
});

const data = await response.json();
if (!data.success) {
  console.error(data.error);
}
```

### 3. 使用時間範圍

對於時間序列查詢，明確指定時間範圍以提升性能。

---

## 版本控制

API 版本通過 URL 路徑指定：

- V2: `/api/v2/*` (當前穩定版)
- V3: `/api/v3/*` (未來版本)
- Experimental: `/api/v2/experimental/*`

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16



> **版本**: 3.0.0  
> **日期**: 2025-10-16  
> **Base URL**: `http://localhost:3001`

---

## 📋 目錄

1. [系統健康](#系統健康)
2. [Prometheus APIs](#prometheus-apis)
3. [Loki APIs](#loki-apis)
4. [Quantum APIs](#quantum-apis)
5. [Nginx APIs](#nginx-apis)
6. [Windows Logs APIs](#windows-logs-apis)
7. [認證](#認證)
8. [錯誤處理](#錯誤處理)

---

## 系統健康

### GET /health

系統健康檢查

**Response**:
```json
{
  "status": "healthy",
  "service": "axiom-backend-v3",
  "version": "3.0.0",
  "time": "2025-10-16T10:00:00Z"
}
```

---

## Prometheus APIs

Base path: `/api/v2/prometheus`

### POST /api/v2/prometheus/query

執行 PromQL 即時查詢

**Request Body**:
```json
{
  "query": "up{job='node-exporter'}",
  "time": "2025-10-16T10:00:00Z" // 可選
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "success",
    "data": {
      "result_type": "vector",
      "result": [
        {
          "metric": {"job": "node-exporter"},
          "value": [1697452800, "1"]
        }
      ]
    },
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/prometheus/query-range

執行 PromQL 範圍查詢

**Request Body**:
```json
{
  "query": "rate(http_requests_total[5m])",
  "start": "2025-10-16T09:00:00Z",
  "end": "2025-10-16T10:00:00Z",
  "step": "15s"
}
```

### GET /api/v2/prometheus/rules

獲取所有告警規則

**Response**:
```json
{
  "success": true,
  "data": {
    "groups": [
      {
        "name": "example-rules",
        "file": "/etc/prometheus/rules.yml",
        "rules": [
          {
            "name": "HighErrorRate",
            "query": "rate(errors_total[5m]) > 0.05",
            "state": "firing"
          }
        ]
      }
    ]
  }
}
```

### GET /api/v2/prometheus/targets

獲取所有抓取目標

### GET /api/v2/prometheus/health

Prometheus 服務健康檢查

### GET /api/v2/prometheus/status

獲取 Prometheus 服務狀態

---

## Loki APIs

Base path: `/api/v2/loki`

### GET /api/v2/loki/query

查詢 Loki 日誌

**Query Parameters**:
- `query` (required): LogQL 查詢語句
- `limit` (optional): 限制返回數量
- `start` (optional): 開始時間 (RFC3339)
- `end` (optional): 結束時間 (RFC3339)

**Example**:
```
GET /api/v2/loki/query?query={job="varlogs"}|="error"&limit=100
```

**Response**:
```json
{
  "success": true,
  "data": {
    "status": "success",
    "data": {
      "resultType": "streams",
      "result": [...]
    }
  }
}
```

### GET /api/v2/loki/labels

獲取所有可用標籤

**Response**:
```json
{
  "success": true,
  "data": {
    "labels": ["job", "filename", "level"]
  }
}
```

### GET /api/v2/loki/labels/{label}/values

獲取指定標籤的所有值

**Example**:
```
GET /api/v2/loki/labels/job/values
```

### GET /api/v2/loki/health

Loki 服務健康檢查

---

## Quantum APIs

Base path: `/api/v2/quantum`

### POST /api/v2/quantum/qkd/generate

生成量子密鑰分發

**Request Body**:
```json
{
  "key_length": 256,
  "backend": "simulator", // "simulator" or "ibm_quantum"
  "shots": 1024
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "job_id": "qkd-abc12345",
    "key": "base64_encoded_key",
    "key_length": 256,
    "status": "completed",
    "submitted_at": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/quantum/qsvm/classify

執行量子支持向量機分類

**Request Body**:
```json
{
  "features": [0.1, 0.2, 0.3, 0.4],
  "backend": "simulator",
  "feature_dim": 4,
  "shots": 1024
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "job_id": "qsvm-def67890",
    "prediction": 1,
    "probability": 0.87,
    "confidence": 0.95,
    "status": "completed",
    "submitted_at": "2025-10-16T10:00:00Z"
  }
}
```

### POST /api/v2/quantum/zerotrust/predict

執行 Zero Trust 安全預測

**Request Body**:
```json
{
  "user_id": "user123",
  "ip_address": "192.168.1.100",
  "device_id": "device_456",
  "features": {
    "login_time": "10:00",
    "location": "taipei"
  },
  "use_quantum": true
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "trust_score": 0.85,
    "risk_level": "low",
    "decision": "allow",
    "confidence": 0.92,
    "factors": {
      "user_behavior": 0.9,
      "device_trust": 0.8,
      "location_risk": 0.85
    },
    "used_quantum": true,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/quantum/jobs

列出所有量子作業

**Query Parameters**:
- `type`: 作業類型 (qkd, qsvm, qaoa, etc.)
- `status`: 作業狀態 (pending, running, completed, failed)
- `page`: 頁碼
- `page_size`: 每頁數量

### GET /api/v2/quantum/jobs/{jobId}

獲取單個量子作業詳情

### GET /api/v2/quantum/stats

獲取量子作業統計

**Response**:
```json
{
  "success": true,
  "data": {
    "total_jobs": 1523,
    "completed_jobs": 1420,
    "failed_jobs": 53,
    "running_jobs": 50,
    "success_rate": 0.93,
    "jobs_by_type": {
      "qkd": 520,
      "qsvm": 380,
      "qaoa": 250,
      "zerotrust": 373
    }
  }
}
```

### GET /api/v2/quantum/health

Quantum 服務健康檢查

---

## Nginx APIs

Base path: `/api/v2/nginx`

### GET /api/v2/nginx/status

獲取 Nginx 狀態

**Response**:
```json
{
  "success": true,
  "data": {
    "name": "nginx",
    "status": "healthy",
    "active_connections": 42,
    "accepts": 15230,
    "handled": 15230,
    "requests": 30460,
    "reading": 0,
    "writing": 3,
    "waiting": 39,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/nginx/config

獲取 Nginx 配置

**Response**:
```json
{
  "success": true,
  "data": {
    "config": "user nginx;\nworker_processes auto;\n...",
    "config_path": "/etc/nginx/nginx.conf",
    "last_modified": "2025-10-16T09:00:00Z",
    "size": 4096,
    "valid": true
  }
}
```

### PUT /api/v2/nginx/config

更新 Nginx 配置

**Request Body**:
```json
{
  "config": "user nginx;\nworker_processes auto;\n...",
  "validate": true,
  "backup": true
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "config_path": "/etc/nginx/nginx.conf",
    "last_modified": "2025-10-16T10:05:00Z",
    "valid": true
  }
}
```

### POST /api/v2/nginx/reload

重載 Nginx 配置

**Response**:
```json
{
  "success": true,
  "data": {
    "success": true,
    "message": "Nginx reloaded successfully",
    "duration": 125,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

---

## Windows Logs APIs

Base path: `/api/v2/logs/windows`

### POST /api/v2/logs/windows/batch

批量接收 Windows 日誌

**Request Body**:
```json
{
  "agent_id": "agent-001",
  "computer": "WIN-SERVER-01",
  "logs": [
    {
      "log_type": "Security",
      "source": "Microsoft-Windows-Security-Auditing",
      "event_id": 4624,
      "level": "Information",
      "message": "An account was successfully logged on",
      "time_created": "2025-10-16T10:00:00Z",
      "user_id": "S-1-5-21-...",
      "process_id": 1234,
      "thread_id": 5678,
      "metadata": {}
    }
  ]
}
```

**Response**:
```json
{
  "success": true,
  "data": {
    "received_count": 1,
    "saved_count": 1,
    "failed_count": 0,
    "errors": [],
    "timestamp": "2025-10-16T10:00:00Z",
    "message": "Logs processed successfully"
  }
}
```

### GET /api/v2/logs/windows

查詢 Windows 日誌

**Query Parameters**:
- `agent_id`: Agent ID 過濾
- `log_type`: 日誌類型 (System, Security, Application, Setup)
- `level`: 日誌級別 (Critical, Error, Warning, Information)
- `event_id`: 事件 ID
- `keyword`: 關鍵字搜索
- `start_time`: 開始時間 (RFC3339)
- `end_time`: 結束時間 (RFC3339)
- `page`: 頁碼 (默認 1)
- `page_size`: 每頁數量 (默認 50, 最大 1000)

**Response**:
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 12345,
        "agent_id": "agent-001",
        "log_type": "Security",
        "source": "Microsoft-Windows-Security-Auditing",
        "event_id": 4624,
        "level": "Information",
        "message": "An account was successfully logged on",
        "time_created": "2025-10-16T10:00:00Z",
        "received_at": "2025-10-16T10:00:05Z",
        "computer": "WIN-SERVER-01"
      }
    ],
    "total": 1523,
    "page": 1,
    "page_size": 50,
    "total_pages": 31,
    "timestamp": "2025-10-16T10:00:00Z"
  }
}
```

### GET /api/v2/logs/windows/stats

獲取 Windows 日誌統計

**Query Parameters**:
- `time_range`: 時間範圍 (1h, 24h, 7d, 30d)

**Response**:
```json
{
  "success": true,
  "data": {
    "total_logs": 15230,
    "logs_by_type": {
      "System": 5420,
      "Security": 6830,
      "Application": 2980
    },
    "logs_by_level": {
      "Critical": 12,
      "Error": 145,
      "Warning": 523,
      "Information": 14550
    },
    "critical_count": 12,
    "error_count": 145,
    "warning_count": 523,
    "time_range": "24h"
  }
}
```

---

## 認證

### Header 認證

所有 API 請求應包含認證 Header：

```http
Authorization: Bearer <JWT_TOKEN>
```

或使用 API Key：

```http
X-API-Key: <API_KEY>
```

---

## 錯誤處理

### 錯誤響應格式

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request",
    "details": "field 'query' is required"
  }
}
```

### 錯誤碼列表

| 錯誤碼 | HTTP 狀態 | 說明 |
|--------|----------|------|
| `INTERNAL_ERROR` | 500 | 內部服務器錯誤 |
| `VALIDATION_ERROR` | 400 | 請求驗證失敗 |
| `NOT_FOUND` | 404 | 資源不存在 |
| `UNAUTHORIZED` | 401 | 未授權 |
| `FORBIDDEN` | 403 | 禁止訪問 |
| `SERVICE_UNAVAILABLE` | 503 | 服務不可用 |
| `TIMEOUT` | 408 | 請求超時 |

---

## 速率限制

- 默認限制：100 requests/minute
- Burst 限制：200 requests
- Header: `X-RateLimit-Remaining`, `X-RateLimit-Reset`

---

## 最佳實踐

### 1. 使用分頁

對於返回列表的 API，始終使用分頁：

```
GET /api/v2/quantum/jobs?page=1&page_size=50
```

### 2. 處理錯誤

始終檢查 `success` 欄位：

```typescript
const response = await fetch('/api/v2/prometheus/query', {
  method: 'POST',
  body: JSON.stringify({ query: 'up' }),
});

const data = await response.json();
if (!data.success) {
  console.error(data.error);
}
```

### 3. 使用時間範圍

對於時間序列查詢，明確指定時間範圍以提升性能。

---

## 版本控制

API 版本通過 URL 路徑指定：

- V2: `/api/v2/*` (當前穩定版)
- V3: `/api/v3/*` (未來版本)
- Experimental: `/api/v2/experimental/*`

---

**文檔版本**: 3.0.0  
**最後更新**: 2025-10-16

