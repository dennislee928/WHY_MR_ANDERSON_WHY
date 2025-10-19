package axiom

// Swagger 文檔結構
const SwaggerDoc = `{
  "swagger": "2.0",
  "info": {
    "title": "Pandora Box Console IDS-IPS API",
    "description": "智慧型入侵偵測與防護系統 RESTful API",
    "version": "3.0.0",
    "contact": {
      "name": "Pandora Security Team",
      "email": "support@pandora-ids.com"
    },
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/licenses/MIT"
    }
  },
  "host": "localhost:3001",
  "basePath": "/api/v1",
  "schemes": ["http", "https"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "securityDefinitions": {
    "ApiKeyAuth": {
      "type": "apiKey",
      "in": "header",
      "name": "Authorization"
    }
  },
  "tags": [
    {
      "name": "System",
      "description": "系統狀態與健康檢查"
    },
    {
      "name": "Dashboard",
      "description": "儀表板數據"
    },
    {
      "name": "Security",
      "description": "安全監控與威脅偵測"
    },
    {
      "name": "Network",
      "description": "網路管理與流量監控"
    },
    {
      "name": "Alerts",
      "description": "告警管理"
    },
    {
      "name": "Events",
      "description": "事件管理"
    },
    {
      "name": "Control",
      "description": "系統控制"
    },
    {
      "name": "Metrics",
      "description": "指標數據"
    },
    {
      "name": "Monitoring",
      "description": "監控服務"
    }
  ],
  "paths": {
    "/status": {
      "get": {
        "tags": ["System"],
        "summary": "取得系統狀態",
        "description": "返回完整的系統狀態資訊，包括 Agent、Console、Network、Security 等",
        "operationId": "getSystemStatus",
        "responses": {
          "200": {
            "description": "成功返回系統狀態",
            "schema": {
              "$ref": "#/definitions/SystemStatus"
            }
          }
        }
      }
    },
    "/health": {
      "get": {
        "tags": ["System"],
        "summary": "健康檢查",
        "description": "檢查系統健康狀態，包括所有服務的運行狀態",
        "operationId": "getHealth",
        "responses": {
          "200": {
            "description": "系統健康",
            "schema": {
              "$ref": "#/definitions/HealthResponse"
            }
          }
        }
      }
    },
    "/dashboard": {
      "get": {
        "tags": ["Dashboard"],
        "summary": "取得儀表板數據",
        "description": "返回儀表板所需的所有數據，包括系統狀態、告警、事件和指標",
        "operationId": "getDashboardData",
        "responses": {
          "200": {
            "description": "成功返回儀表板數據",
            "schema": {
              "$ref": "#/definitions/DashboardData"
            }
          }
        }
      }
    },
    "/security/threats": {
      "get": {
        "tags": ["Security"],
        "summary": "取得威脅事件列表",
        "description": "返回威脅偵測事件列表，支援過濾和分頁",
        "operationId": "getThreatEvents",
        "parameters": [
          {
            "name": "severity",
            "in": "query",
            "description": "威脅嚴重程度過濾",
            "type": "string",
            "enum": ["all", "critical", "high", "medium", "low"],
            "default": "all"
          },
          {
            "name": "time_range",
            "in": "query",
            "description": "時間範圍",
            "type": "string",
            "enum": ["1h", "24h", "7d", "30d"],
            "default": "24h"
          },
          {
            "name": "limit",
            "in": "query",
            "description": "每頁數量",
            "type": "integer",
            "default": 50
          },
          {
            "name": "offset",
            "in": "query",
            "description": "偏移量",
            "type": "integer",
            "default": 0
          }
        ],
        "responses": {
          "200": {
            "description": "成功返回威脅事件",
            "schema": {
              "$ref": "#/definitions/ThreatEventsResponse"
            }
          }
        }
      }
    },
    "/security/stats": {
      "get": {
        "tags": ["Security"],
        "summary": "取得安全統計",
        "description": "返回安全相關的統計數據",
        "operationId": "getSecurityStats",
        "responses": {
          "200": {
            "description": "成功返回安全統計",
            "schema": {
              "$ref": "#/definitions/ThreatStats"
            }
          }
        }
      }
    },
    "/network/stats": {
      "get": {
        "tags": ["Network"],
        "summary": "取得網路統計",
        "description": "返回網路流量和連線統計",
        "operationId": "getNetworkStats",
        "responses": {
          "200": {
            "description": "成功返回網路統計",
            "schema": {
              "$ref": "#/definitions/NetworkStats"
            }
          }
        }
      }
    },
    "/network/blocked-ips": {
      "get": {
        "tags": ["Network"],
        "summary": "取得被阻斷的 IP 列表",
        "description": "返回所有被阻斷的 IP 位址及其資訊",
        "operationId": "getBlockedIPs",
        "responses": {
          "200": {
            "description": "成功返回被阻斷的 IP",
            "schema": {
              "$ref": "#/definitions/BlockedIPsResponse"
            }
          }
        }
      }
    },
    "/network/blocked-ips/{ip}": {
      "delete": {
        "tags": ["Network"],
        "summary": "解除 IP 阻斷",
        "description": "移除指定 IP 的阻斷狀態",
        "operationId": "unblockIP",
        "parameters": [
          {
            "name": "ip",
            "in": "path",
            "description": "IP 位址",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "成功解除阻斷",
            "schema": {
              "$ref": "#/definitions/SuccessResponse"
            }
          }
        }
      }
    },
    "/network/interfaces": {
      "get": {
        "tags": ["Network"],
        "summary": "取得網路介面列表",
        "description": "返回所有網路介面的狀態和統計",
        "operationId": "getNetworkInterfaces",
        "responses": {
          "200": {
            "description": "成功返回網路介面",
            "schema": {
              "$ref": "#/definitions/NetworkInterfacesResponse"
            }
          }
        }
      }
    },
    "/alerts": {
      "get": {
        "tags": ["Alerts"],
        "summary": "取得告警列表",
        "description": "返回系統告警列表，支援過濾和分頁",
        "operationId": "getAlerts",
        "parameters": [
          {
            "name": "level",
            "in": "query",
            "description": "告警級別",
            "type": "string",
            "enum": ["critical", "warning", "info"]
          },
          {
            "name": "resolved",
            "in": "query",
            "description": "是否已解決",
            "type": "boolean"
          },
          {
            "name": "limit",
            "in": "query",
            "description": "每頁數量",
            "type": "integer",
            "default": 50
          },
          {
            "name": "offset",
            "in": "query",
            "description": "偏移量",
            "type": "integer",
            "default": 0
          }
        ],
        "responses": {
          "200": {
            "description": "成功返回告警列表",
            "schema": {
              "$ref": "#/definitions/AlertsResponse"
            }
          }
        }
      }
    },
    "/alerts/{id}/resolve": {
      "post": {
        "tags": ["Alerts"],
        "summary": "解決告警",
        "description": "標記指定告警為已解決",
        "operationId": "resolveAlert",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "description": "告警 ID",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "成功解決告警",
            "schema": {
              "$ref": "#/definitions/SuccessResponse"
            }
          }
        }
      }
    },
    "/events": {
      "get": {
        "tags": ["Events"],
        "summary": "取得事件列表",
        "description": "返回系統事件列表，支援過濾和分頁",
        "operationId": "getEvents",
        "parameters": [
          {
            "name": "type",
            "in": "query",
            "description": "事件類型",
            "type": "string"
          },
          {
            "name": "limit",
            "in": "query",
            "description": "每頁數量",
            "type": "integer",
            "default": 50
          },
          {
            "name": "offset",
            "in": "query",
            "description": "偏移量",
            "type": "integer",
            "default": 0
          }
        ],
        "responses": {
          "200": {
            "description": "成功返回事件列表",
            "schema": {
              "$ref": "#/definitions/EventsResponse"
            }
          }
        }
      }
    },
    "/events/{id}": {
      "get": {
        "tags": ["Events"],
        "summary": "取得單個事件",
        "description": "返回指定 ID 的事件詳情",
        "operationId": "getEvent",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "description": "事件 ID",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "成功返回事件詳情",
            "schema": {
              "$ref": "#/definitions/Event"
            }
          }
        }
      }
    },
    "/control/network": {
      "post": {
        "tags": ["Control"],
        "summary": "控制網路狀態",
        "description": "阻斷或解除網路阻斷",
        "operationId": "controlNetwork",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "description": "控制動作",
            "required": true,
            "schema": {
              "$ref": "#/definitions/NetworkControlRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "成功執行控制動作",
            "schema": {
              "$ref": "#/definitions/SuccessResponse"
            }
          },
          "400": {
            "description": "無效的請求",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/control/network/status": {
      "get": {
        "tags": ["Control"],
        "summary": "取得網路控制狀態",
        "description": "返回當前網路控制狀態",
        "operationId": "getNetworkStatus",
        "responses": {
          "200": {
            "description": "成功返回網路狀態",
            "schema": {
              "$ref": "#/definitions/NetworkControlStatus"
            }
          }
        }
      }
    },
    "/control/device": {
      "post": {
        "tags": ["Control"],
        "summary": "控制裝置",
        "description": "執行裝置控制動作",
        "operationId": "controlDevice",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "description": "控制動作",
            "required": true,
            "schema": {
              "$ref": "#/definitions/DeviceControlRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "成功執行控制動作",
            "schema": {
              "$ref": "#/definitions/SuccessResponse"
            }
          }
        }
      }
    },
    "/control/device/status": {
      "get": {
        "tags": ["Control"],
        "summary": "取得裝置狀態",
        "description": "返回所有裝置的狀態資訊",
        "operationId": "getDeviceStatus",
        "responses": {
          "200": {
            "description": "成功返回裝置狀態",
            "schema": {
              "$ref": "#/definitions/DeviceStatusResponse"
            }
          }
        }
      }
    },
    "/metrics": {
      "get": {
        "tags": ["Metrics"],
        "summary": "取得指標數據",
        "description": "返回系統指標數據",
        "operationId": "getMetrics",
        "responses": {
          "200": {
            "description": "成功返回指標數據",
            "schema": {
              "type": "object"
            }
          }
        }
      }
    },
    "/metrics/prometheus": {
      "get": {
        "tags": ["Metrics"],
        "summary": "取得 Prometheus 指標",
        "description": "返回 Prometheus 格式的指標數據",
        "operationId": "getPrometheusMetrics",
        "responses": {
          "200": {
            "description": "成功返回 Prometheus 指標",
            "schema": {
              "type": "object"
            }
          }
        }
      }
    },
    "/monitoring/services": {
      "get": {
        "tags": ["Monitoring"],
        "summary": "取得監控服務狀態",
        "description": "返回所有監控服務的狀態",
        "operationId": "getMonitoringServices",
        "responses": {
          "200": {
            "description": "成功返回監控服務狀態",
            "schema": {
              "type": "object"
            }
          }
        }
      }
    },
    "/monitoring/services/{service}/status": {
      "get": {
        "tags": ["Monitoring"],
        "summary": "取得單個服務狀態",
        "description": "返回指定監控服務的詳細狀態",
        "operationId": "getServiceStatus",
        "parameters": [
          {
            "name": "service",
            "in": "path",
            "description": "服務名稱",
            "required": true,
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "成功返回服務狀態",
            "schema": {
              "type": "object"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "SystemStatus": {
      "type": "object",
      "properties": {
        "agent": {
          "type": "object",
          "properties": {
            "status": {"type": "string"},
            "lastSeen": {"type": "string"},
            "uptime": {"type": "string"},
            "version": {"type": "string"},
            "cpuUsage": {"type": "number"},
            "memoryUsage": {"type": "number"},
            "diskUsage": {"type": "number"}
          }
        },
        "console": {
          "type": "object"
        },
        "network": {
          "type": "object"
        },
        "security": {
          "type": "object"
        },
        "monitoring": {
          "type": "object"
        },
        "devices": {
          "type": "object"
        }
      }
    },
    "HealthResponse": {
      "type": "object",
      "properties": {
        "status": {"type": "string"},
        "timestamp": {"type": "integer"},
        "uptime": {"type": "string"},
        "version": {"type": "string"},
        "services": {"type": "object"}
      }
    },
    "DashboardData": {
      "type": "object",
      "properties": {
        "system_status": {"type": "object"},
        "alerts": {"type": "array"},
        "events": {"type": "array"},
        "metrics": {"type": "object"},
        "timestamp": {"type": "integer"}
      }
    },
    "ThreatEventsResponse": {
      "type": "object",
      "properties": {
        "threats": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ThreatEvent"
          }
        },
        "total": {"type": "integer"},
        "limit": {"type": "integer"},
        "offset": {"type": "integer"}
      }
    },
    "ThreatEvent": {
      "type": "object",
      "properties": {
        "id": {"type": "string"},
        "timestamp": {"type": "string"},
        "type": {"type": "string"},
        "severity": {"type": "string"},
        "source_ip": {"type": "string"},
        "destination_ip": {"type": "string"},
        "port": {"type": "integer"},
        "protocol": {"type": "string"},
        "action": {"type": "string"},
        "description": {"type": "string"},
        "rule_id": {"type": "string"}
      }
    },
    "ThreatStats": {
      "type": "object",
      "properties": {
        "total_threats": {"type": "integer"},
        "blocked_threats": {"type": "integer"},
        "active_threats": {"type": "integer"},
        "resolved_threats": {"type": "integer"},
        "threat_trend": {"type": "number"},
        "top_threat_types": {"type": "array"},
        "top_source_ips": {"type": "array"}
      }
    },
    "NetworkStats": {
      "type": "object",
      "properties": {
        "total_traffic": {"type": "integer"},
        "inbound_traffic": {"type": "integer"},
        "outbound_traffic": {"type": "integer"},
        "packet_loss": {"type": "number"},
        "latency": {"type": "number"},
        "bandwidth_usage": {"type": "number"},
        "active_connections": {"type": "integer"},
        "blocked_connections": {"type": "integer"}
      }
    },
    "BlockedIPsResponse": {
      "type": "object",
      "properties": {
        "blocked_ips": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/BlockedIP"
          }
        }
      }
    },
    "BlockedIP": {
      "type": "object",
      "properties": {
        "ip": {"type": "string"},
        "reason": {"type": "string"},
        "blocked_at": {"type": "string"},
        "expires_at": {"type": "string"},
        "threat_count": {"type": "integer"}
      }
    },
    "NetworkInterfacesResponse": {
      "type": "object",
      "properties": {
        "interfaces": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/NetworkInterface"
          }
        }
      }
    },
    "NetworkInterface": {
      "type": "object",
      "properties": {
        "name": {"type": "string"},
        "status": {"type": "string"},
        "ip_address": {"type": "string"},
        "mac_address": {"type": "string"},
        "rx_bytes": {"type": "integer"},
        "tx_bytes": {"type": "integer"},
        "rx_packets": {"type": "integer"},
        "tx_packets": {"type": "integer"}
      }
    },
    "AlertsResponse": {
      "type": "object",
      "properties": {
        "alerts": {"type": "array"},
        "total": {"type": "integer"},
        "limit": {"type": "integer"},
        "offset": {"type": "integer"}
      }
    },
    "EventsResponse": {
      "type": "object",
      "properties": {
        "events": {"type": "array"},
        "total": {"type": "integer"},
        "limit": {"type": "integer"},
        "offset": {"type": "integer"}
      }
    },
    "Event": {
      "type": "object",
      "properties": {
        "id": {"type": "string"},
        "type": {"type": "string"},
        "severity": {"type": "string"},
        "message": {"type": "string"},
        "timestamp": {"type": "integer"},
        "data": {"type": "object"}
      }
    },
    "NetworkControlRequest": {
      "type": "object",
      "required": ["action"],
      "properties": {
        "action": {
          "type": "string",
          "enum": ["block", "unblock"]
        }
      }
    },
    "DeviceControlRequest": {
      "type": "object",
      "required": ["action"],
      "properties": {
        "action": {"type": "string"},
        "data": {"type": "object"}
      }
    },
    "NetworkControlStatus": {
      "type": "object",
      "properties": {
        "blocked": {"type": "boolean"},
        "block_time": {"type": "string"},
        "unlock_time": {"type": "string"},
        "total_traffic": {"type": "integer"},
        "blocked_ips": {"type": "integer"},
        "allowed_ips": {"type": "integer"},
        "last_update": {"type": "integer"}
      }
    },
    "DeviceStatusResponse": {
      "type": "object",
      "properties": {
        "total_devices": {"type": "integer"},
        "online_devices": {"type": "integer"},
        "offline_devices": {"type": "integer"},
        "last_activity": {"type": "string"},
        "devices": {"type": "array"}
      }
    },
    "SuccessResponse": {
      "type": "object",
      "properties": {
        "status": {"type": "string"},
        "message": {"type": "string"},
        "timestamp": {"type": "integer"}
      }
    },
    "ErrorResponse": {
      "type": "object",
      "properties": {
        "error": {"type": "string"},
        "message": {"type": "string"},
        "timestamp": {"type": "integer"}
      }
    }
  }
}`
