package axiom

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"pandora_box_console_ids_ips/internal/metrics"
)

// UIServer Axiom UI伺服器
type UIServer struct {
	logger         *logrus.Logger
	metricsClient  *metrics.PrometheusMetrics
	websocketConns map[string]*websocket.Conn
	upgrader       websocket.Upgrader
	startTime      time.Time
}

// SystemStatus 系統狀態
type SystemStatus struct {
	Agent struct {
		Status     string  `json:"status"`
		LastSeen   string  `json:"lastSeen"`
		Uptime     string  `json:"uptime"`
		Version    string  `json:"version"`
		CPUUsage   float64 `json:"cpuUsage"`
		MemoryUsage float64 `json:"memoryUsage"`
		DiskUsage  float64 `json:"diskUsage"`
	} `json:"agent"`
	Console struct {
		Status            string  `json:"status"`
		Requests          int     `json:"requests"`
		ResponseTime      float64 `json:"responseTime"`
		ActiveConnections int     `json:"activeConnections"`
		ErrorRate         float64 `json:"errorRate"`
	} `json:"console"`
	Network struct {
		Blocked    bool    `json:"blocked"`
		BlockTime  string  `json:"blockTime"`
		UnlockTime string  `json:"unlockTime"`
		TotalTraffic int64 `json:"totalTraffic"`
		BlockedIPs  int    `json:"blockedIPs"`
		AllowedIPs  int    `json:"allowedIPs"`
	} `json:"network"`
	Security struct {
		TotalAlerts     int    `json:"totalAlerts"`
		CriticalAlerts  int    `json:"criticalAlerts"`
		WarningAlerts   int    `json:"warningAlerts"`
		InfoAlerts      int    `json:"infoAlerts"`
		ThreatsBlocked  int    `json:"threatsBlocked"`
		LastThreat      string `json:"lastThreat"`
	} `json:"security"`
	Monitoring struct {
		Prometheus   bool `json:"prometheus"`
		Grafana      bool `json:"grafana"`
		Loki         bool `json:"loki"`
		AlertManager bool `json:"alertmanager"`
	} `json:"monitoring"`
	Devices struct {
		Total        int    `json:"total"`
		Online       int    `json:"online"`
		Offline      int    `json:"offline"`
		LastActivity string `json:"lastActivity"`
	} `json:"devices"`
}

// Alert 告警結構
type Alert struct {
	ID        string    `json:"id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"`
	Resolved  bool      `json:"resolved"`
}

// WebSocketMessage WebSocket訊息結構
type WebSocketMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// NewUIServer 建立新的UI伺服器
func NewUIServer(logger *logrus.Logger, metricsClient *metrics.PrometheusMetrics) *UIServer {
	return &UIServer{
		logger:         logger,
		metricsClient:  metricsClient,
		websocketConns: make(map[string]*websocket.Conn),
		startTime:      time.Now(),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 在生產環境中應該檢查來源
			},
		},
	}
}

// StartUIServer 啟動UI伺服器
func (ui *UIServer) StartUIServer(port string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(ui.corsMiddleware())
	router.Use(ui.loggingMiddleware())

	// 靜態檔案服務
	router.Static("/static", "./web/static")
	router.StaticFile("/", "/app/web/public/index.html")
	router.StaticFile("/favicon.ico", "/app/web/public/favicon.ico")
	
	// 添加調試路由
	router.GET("/debug", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Debug endpoint working",
			"path": "/app/web/public/index.html",
			"uptime": time.Since(ui.startTime).String(),
		})
	})

	// Swagger 文檔
	router.GET("/swagger.json", ui.getSwaggerJSON)
	router.GET("/swagger", ui.getSwaggerUI)
	router.GET("/api-docs", ui.getSwaggerUI)

	// API 路由
	api := router.Group("/api/v1")
	{
		// 系統狀態
		api.GET("/status", ui.getSystemStatus)
		api.GET("/health", ui.getHealth)
		
		// 儀表板數據
		api.GET("/dashboard", ui.getDashboardData)
		
		// 告警管理
		api.GET("/alerts", ui.getAlerts)
		api.POST("/alerts/:id/resolve", ui.resolveAlert)
		
		// 安全監控
		api.GET("/security/threats", ui.getThreatEvents)
		api.GET("/security/stats", ui.getSecurityStats)
		api.POST("/security/threats/:id/block", ui.blockThreatSource)
		
		// 網路管理
		api.GET("/network/stats", ui.getNetworkStats)
		api.GET("/network/blocked-ips", ui.getBlockedIPs)
		api.DELETE("/network/blocked-ips/:ip", ui.unblockIP)
		api.GET("/network/interfaces", ui.getNetworkInterfaces)
		
		// 設備管理
		api.GET("/devices", ui.getDevices)
		api.GET("/devices/:id", ui.getDeviceDetail)
		api.POST("/devices/:id/restart", ui.restartDevice)
		api.PUT("/devices/:id/config", ui.updateDeviceConfig)
		
		// 報表生成
		api.GET("/reports/security", ui.generateSecurityReport)
		api.GET("/reports/network", ui.generateNetworkReport)
		api.GET("/reports/system", ui.generateSystemReport)
		api.POST("/reports/custom", ui.generateCustomReport)
		
		// 事件管理
		api.GET("/events", ui.getEvents)
		api.GET("/events/:id", ui.getEvent)
		
		// 網路控制
		api.POST("/control/network", ui.controlNetwork)
		api.GET("/control/network/status", ui.getNetworkStatus)
		
		// 裝置控制
		api.POST("/control/device", ui.controlDevice)
		api.GET("/control/device/status", ui.getDeviceStatus)
		
		// 指標數據
		api.GET("/metrics", ui.getMetrics)
		api.GET("/metrics/prometheus", ui.getPrometheusMetrics)
		
		// 監控服務狀態
		api.GET("/monitoring/services", ui.getMonitoringServices)
		api.GET("/monitoring/services/:service/status", ui.getServiceStatus)
	}

	// WebSocket 路由
	router.GET("/ws", ui.handleWebSocket)

	// 啟動定期數據推送
	go ui.startPeriodicUpdates()

	ui.logger.Infof("Axiom UI伺服器啟動於端口: %s", port)
	return router.Run(":" + port)
}

// corsMiddleware CORS中間件
func (ui *UIServer) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// loggingMiddleware 日誌中間件
func (ui *UIServer) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		ui.logger.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency":    latency,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"body_size":  bodySize,
		}).Info("HTTP Request")
	}
}

// getSystemStatus 取得系統狀態
func (ui *UIServer) getSystemStatus(c *gin.Context) {
	status := SystemStatus{
		Agent: struct {
			Status     string  `json:"status"`
			LastSeen   string  `json:"lastSeen"`
			Uptime     string  `json:"uptime"`
			Version    string  `json:"version"`
			CPUUsage   float64 `json:"cpuUsage"`
			MemoryUsage float64 `json:"memoryUsage"`
			DiskUsage  float64 `json:"diskUsage"`
		}{
			Status:      "online",
			LastSeen:    time.Now().Format("2006-01-02 15:04:05"),
			Uptime:      time.Since(ui.startTime).String(),
			Version:     "3.0.0",
			CPUUsage:    15.5,
			MemoryUsage: 42.3,
			DiskUsage:   28.7,
		},
		Console: struct {
			Status            string  `json:"status"`
			Requests          int     `json:"requests"`
			ResponseTime      float64 `json:"responseTime"`
			ActiveConnections int     `json:"activeConnections"`
			ErrorRate         float64 `json:"errorRate"`
		}{
			Status:            "online",
			Requests:          1250,
			ResponseTime:      45.2,
			ActiveConnections: 8,
			ErrorRate:         0.1,
		},
		Network: struct {
			Blocked    bool    `json:"blocked"`
			BlockTime  string  `json:"blockTime"`
			UnlockTime string  `json:"unlockTime"`
			TotalTraffic int64 `json:"totalTraffic"`
			BlockedIPs  int    `json:"blockedIPs"`
			AllowedIPs  int    `json:"allowedIPs"`
		}{
			Blocked:     false,
			BlockTime:   "",
			UnlockTime:  "",
			TotalTraffic: 1024000,
			BlockedIPs:  12,
			AllowedIPs:  156,
		},
		Security: struct {
			TotalAlerts     int    `json:"totalAlerts"`
			CriticalAlerts  int    `json:"criticalAlerts"`
			WarningAlerts   int    `json:"warningAlerts"`
			InfoAlerts      int    `json:"infoAlerts"`
			ThreatsBlocked  int    `json:"threatsBlocked"`
			LastThreat      string `json:"lastThreat"`
		}{
			TotalAlerts:    45,
			CriticalAlerts: 3,
			WarningAlerts:  12,
			InfoAlerts:     30,
			ThreatsBlocked: 8,
			LastThreat:     "2025-01-14 10:30:15",
		},
		Monitoring: struct {
			Prometheus   bool `json:"prometheus"`
			Grafana      bool `json:"grafana"`
			Loki         bool `json:"loki"`
			AlertManager bool `json:"alertmanager"`
		}{
			Prometheus:   true,
			Grafana:      true,
			Loki:         true,
			AlertManager: true,
		},
		Devices: struct {
			Total        int    `json:"total"`
			Online       int    `json:"online"`
			Offline      int    `json:"offline"`
			LastActivity string `json:"lastActivity"`
		}{
			Total:        5,
			Online:       4,
			Offline:      1,
			LastActivity: "2025-01-14 10:25:30",
		},
	}

	c.JSON(http.StatusOK, status)
}

// getHealth 健康檢查
func (ui *UIServer) getHealth(c *gin.Context) {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(ui.startTime).String(),
		"version":   "3.0.0",
		"services": map[string]string{
			"axiom-ui":      "healthy",
			"prometheus":    "healthy",
			"grafana":       "healthy",
			"loki":          "healthy",
			"alertmanager":  "healthy",
		},
	}

	c.JSON(http.StatusOK, health)
}

// getDashboardData 取得儀表板資料
func (ui *UIServer) getDashboardData(c *gin.Context) {
	data := map[string]interface{}{
		"system_status": ui.getSystemStatusData(),
		"alerts":       ui.getAlertsData(),
		"events":       ui.getEventsData(),
		"metrics":      ui.getMetricsData(),
		"timestamp":    time.Now().Unix(),
	}

	c.JSON(http.StatusOK, data)
}

// getAlerts 取得告警列表
func (ui *UIServer) getAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	level := c.Query("level")
	resolved := c.Query("resolved")

	alerts := ui.getAlertsData()
	
	// 過濾告警
	filtered := make([]Alert, 0)
	for _, alert := range alerts {
		if level != "" && alert.Level != level {
			continue
		}
		if resolved != "" {
			resolvedBool, _ := strconv.ParseBool(resolved)
			if alert.Resolved != resolvedBool {
				continue
			}
		}
		filtered = append(filtered, alert)
	}

	// 分頁
	start := offset
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	result := filtered[start:end]

	c.JSON(http.StatusOK, gin.H{
		"alerts": result,
		"total":  len(filtered),
		"limit":  limit,
		"offset": offset,
	})
}

// resolveAlert 解決告警
func (ui *UIServer) resolveAlert(c *gin.Context) {
	alertID := c.Param("id")
	
	ui.logger.Infof("解決告警: %s", alertID)
	
	// 這裡應該更新資料庫中的告警狀態
	ui.broadcastWebSocketMessage("alert_resolved", map[string]interface{}{
		"alert_id": alertID,
		"timestamp": time.Now().Unix(),
	})

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": "告警已解決",
		"alert_id": alertID,
	})
}

// getEvents 取得事件列表
func (ui *UIServer) getEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	eventType := c.Query("type")

	events := ui.getEventsData()
	
	// 過濾事件
	filtered := make([]map[string]interface{}, 0)
	for _, event := range events {
		if eventType != "" && event["type"] != eventType {
			continue
		}
		filtered = append(filtered, event)
	}

	// 分頁
	start := offset
	if start > len(filtered) {
		start = len(filtered)
	}
	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	result := filtered[start:end]

	c.JSON(http.StatusOK, gin.H{
		"events": result,
		"total":  len(filtered),
		"limit":  limit,
		"offset": offset,
	})
}

// getEvent 取得單個事件
func (ui *UIServer) getEvent(c *gin.Context) {
	eventID := c.Param("id")
	
	// 這裡應該從資料庫獲取事件
	event := map[string]interface{}{
		"id":        eventID,
		"type":      "threat_detection",
		"severity":  "high",
		"message":   "偵測到惡意IP連接嘗試",
		"timestamp": time.Now().Unix(),
		"data": map[string]interface{}{
			"source_ip":   "192.168.1.100",
			"threat_type": "brute_force",
		},
	}

	c.JSON(http.StatusOK, event)
}

// controlNetwork 控制網路狀態
func (ui *UIServer) controlNetwork(c *gin.Context) {
	var request struct {
		Action string `json:"action" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch request.Action {
	case "block":
		ui.logger.Info("透過UI請求阻斷網路")
		ui.broadcastWebSocketMessage("network_status", map[string]interface{}{
			"blocked": true,
			"action":  "block",
			"timestamp": time.Now().Unix(),
		})
	case "unblock":
		ui.logger.Info("透過UI請求解除網路阻斷")
		ui.broadcastWebSocketMessage("network_status", map[string]interface{}{
			"blocked": false,
			"action":  "unblock",
			"timestamp": time.Now().Unix(),
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的動作"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success", 
		"action": request.Action,
		"timestamp": time.Now().Unix(),
	})
}

// getNetworkStatus 取得網路狀態
func (ui *UIServer) getNetworkStatus(c *gin.Context) {
	status := map[string]interface{}{
		"blocked":      false,
		"block_time":   "",
		"unlock_time":  "",
		"total_traffic": 1024000,
		"blocked_ips":   12,
		"allowed_ips":   156,
		"last_update":   time.Now().Unix(),
	}

	c.JSON(http.StatusOK, status)
}

// controlDevice 控制裝置
func (ui *UIServer) controlDevice(c *gin.Context) {
	var request struct {
		Action string                 `json:"action" binding:"required"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ui.logger.Infof("透過UI請求裝置操作: %s", request.Action)

	ui.broadcastWebSocketMessage("device_status", map[string]interface{}{
		"action": request.Action,
		"data":   request.Data,
		"timestamp": time.Now().Unix(),
	})

	c.JSON(http.StatusOK, gin.H{
		"status": "success", 
		"action": request.Action,
		"timestamp": time.Now().Unix(),
	})
}

// getDeviceStatus 取得裝置狀態
func (ui *UIServer) getDeviceStatus(c *gin.Context) {
	status := map[string]interface{}{
		"total_devices":   5,
		"online_devices":  4,
		"offline_devices": 1,
		"last_activity":   time.Now().Format("2006-01-02 15:04:05"),
		"devices": []map[string]interface{}{
			{
				"id":       "device_001",
				"name":     "USB-SERIAL CH340",
				"status":   "online",
				"port":     "/dev/ttyUSB0",
				"last_seen": time.Now().Format("2006-01-02 15:04:05"),
			},
			{
				"id":       "device_002",
				"name":     "Network Interface",
				"status":   "online",
				"port":     "eth0",
				"last_seen": time.Now().Format("2006-01-02 15:04:05"),
			},
		},
	}

	c.JSON(http.StatusOK, status)
}

// getMetrics 取得指標數據
func (ui *UIServer) getMetrics(c *gin.Context) {
	metrics := map[string]interface{}{
		"network_blocked":     0,
		"device_connected":    1,
		"threats_detected":    45,
		"packets_inspected":   1000000,
		"blocked_connections": 12,
		"active_sessions":     8,
		"system_uptime":       time.Since(ui.startTime).Seconds(),
		"cpu_usage":           15.5,
		"memory_usage":        42.3,
		"disk_usage":          28.7,
		"timestamp":           time.Now().Unix(),
	}

	c.JSON(http.StatusOK, metrics)
}

// getPrometheusMetrics 取得Prometheus指標
func (ui *UIServer) getPrometheusMetrics(c *gin.Context) {
	// 這裡應該從Prometheus獲取實際指標
	metrics := map[string]interface{}{
		"pandora_system_uptime_seconds": time.Since(ui.startTime).Seconds(),
		"pandora_threats_detected_total": 45,
		"pandora_packets_inspected_total": 1000000,
		"pandora_blocked_connections_total": 12,
		"pandora_active_sessions": 8,
		"pandora_cpu_usage_percent": 15.5,
		"pandora_memory_usage_percent": 42.3,
		"pandora_disk_usage_percent": 28.7,
	}

	c.JSON(http.StatusOK, metrics)
}

// getMonitoringServices 取得監控服務狀態
func (ui *UIServer) getMonitoringServices(c *gin.Context) {
	services := map[string]interface{}{
		"prometheus": map[string]interface{}{
			"status":    "healthy",
			"url":       "http://localhost:9090",
			"uptime":    "2h 15m",
			"last_check": time.Now().Unix(),
		},
		"grafana": map[string]interface{}{
			"status":    "healthy",
			"url":       "http://localhost:3000",
			"uptime":    "2h 15m",
			"last_check": time.Now().Unix(),
		},
		"loki": map[string]interface{}{
			"status":    "healthy",
			"url":       "http://localhost:3100",
			"uptime":    "2h 15m",
			"last_check": time.Now().Unix(),
		},
		"alertmanager": map[string]interface{}{
			"status":    "healthy",
			"url":       "http://localhost:9093",
			"uptime":    "2h 15m",
			"last_check": time.Now().Unix(),
		},
	}

	c.JSON(http.StatusOK, services)
}

// getServiceStatus 取得單個服務狀態
func (ui *UIServer) getServiceStatus(c *gin.Context) {
	serviceName := c.Param("service")
	
	status := map[string]interface{}{
		"service":     serviceName,
		"status":      "healthy",
		"uptime":      "2h 15m",
		"last_check":  time.Now().Unix(),
		"response_time": "45ms",
	}

	c.JSON(http.StatusOK, status)
}

// handleWebSocket 處理WebSocket連接
func (ui *UIServer) handleWebSocket(c *gin.Context) {
	conn, err := ui.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		ui.logger.Errorf("WebSocket升級失敗: %v", err)
		return
	}
	defer conn.Close()

	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = fmt.Sprintf("client_%d", time.Now().UnixNano())
	}

	ui.websocketConns[clientID] = conn
	defer delete(ui.websocketConns, clientID)

	ui.logger.Infof("WebSocket客戶端連接: %s", clientID)

	// 發送初始數據
	if err := ui.sendWebSocketMessage(conn, "connected", map[string]interface{}{
		"client_id": clientID,
		"timestamp": time.Now().Unix(),
		"server_time": time.Now().Format("2006-01-02 15:04:05"),
	}); err != nil {
		ui.logger.Errorf("發送初始數據失敗: %v", err)
	}

	// 監聽客戶端訊息
	for {
		var message WebSocketMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			ui.logger.Debugf("WebSocket讀取錯誤: %v", err)
			break
		}

		ui.handleWebSocketMessage(clientID, message)
	}
}

// handleWebSocketMessage 處理WebSocket訊息
func (ui *UIServer) handleWebSocketMessage(clientID string, message WebSocketMessage) {
	ui.logger.Debugf("收到WebSocket訊息: %s - %s", clientID, message.Type)

	switch message.Type {
	case "ping":
		conn := ui.websocketConns[clientID]
		if conn != nil {
			if err := ui.sendWebSocketMessage(conn, "pong", map[string]interface{}{
				"timestamp": time.Now().Unix(),
			}); err != nil {
				ui.logger.Errorf("發送pong訊息失敗: %v", err)
			}
		}
	case "subscribe":
		ui.logger.Debugf("客戶端 %s 訂閱: %v", clientID, message.Data)
	case "get_status":
		conn := ui.websocketConns[clientID]
		if conn != nil {
			if err := ui.sendWebSocketMessage(conn, "status_update", ui.getSystemStatusData()); err != nil {
				ui.logger.Errorf("發送狀態更新失敗: %v", err)
			}
		}
	}
}

// sendWebSocketMessage 發送WebSocket訊息
func (ui *UIServer) sendWebSocketMessage(conn *websocket.Conn, msgType string, data interface{}) error {
	message := WebSocketMessage{
		Type: msgType,
		Data: data,
	}
	return conn.WriteJSON(message)
}

// broadcastWebSocketMessage 廣播WebSocket訊息
func (ui *UIServer) broadcastWebSocketMessage(msgType string, data interface{}) {
	message := WebSocketMessage{
		Type: msgType,
		Data: data,
	}

	for clientID, conn := range ui.websocketConns {
		if err := conn.WriteJSON(message); err != nil {
			ui.logger.Errorf("發送WebSocket訊息失敗 (客戶端: %s): %v", clientID, err)
			delete(ui.websocketConns, clientID)
		}
	}
}

// startPeriodicUpdates 啟動定期更新
func (ui *UIServer) startPeriodicUpdates() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		ui.sendPeriodicUpdates()
	}
}

// sendPeriodicUpdates 發送定期更新
func (ui *UIServer) sendPeriodicUpdates() {
	data := map[string]interface{}{
		"system_status": ui.getSystemStatusData(),
		"alerts":       ui.getAlertsData(),
		"events":       ui.getEventsData(),
		"timestamp":    time.Now().Unix(),
	}

	ui.broadcastWebSocketMessage("dashboard_update", data)
}

// getSystemStatusData 取得系統狀態數據
func (ui *UIServer) getSystemStatusData() map[string]interface{} {
	return map[string]interface{}{
		"agent": map[string]interface{}{
			"status":       "online",
			"last_seen":    time.Now().Format("2006-01-02 15:04:05"),
			"uptime":       time.Since(ui.startTime).String(),
			"version":      "3.0.0",
			"cpu_usage":    15.5,
			"memory_usage": 42.3,
			"disk_usage":   28.7,
		},
		"console": map[string]interface{}{
			"status":             "online",
			"requests":           1250,
			"response_time":      45.2,
			"active_connections": 8,
			"error_rate":         0.1,
		},
		"network": map[string]interface{}{
			"blocked":       false,
			"block_time":    "",
			"unlock_time":   "",
			"total_traffic": 1024000,
			"blocked_ips":   12,
			"allowed_ips":   156,
		},
		"security": map[string]interface{}{
			"total_alerts":    45,
			"critical_alerts": 3,
			"warning_alerts":  12,
			"info_alerts":     30,
			"threats_blocked": 8,
			"last_threat":     "2025-01-14 10:30:15",
		},
		"monitoring": map[string]interface{}{
			"prometheus":   true,
			"grafana":      true,
			"loki":         true,
			"alertmanager": true,
		},
		"devices": map[string]interface{}{
			"total":         5,
			"online":        4,
			"offline":       1,
			"last_activity": "2025-01-14 10:25:30",
		},
	}
}

// getAlertsData 取得告警數據
func (ui *UIServer) getAlertsData() []Alert {
	return []Alert{
		{
			ID:        "alert_001",
			Level:     "critical",
			Message:   "偵測到惡意IP連接嘗試",
			Timestamp: time.Now().Add(-5 * time.Minute),
			Source:    "threat_detection",
			Resolved:  false,
		},
		{
			ID:        "alert_002",
			Level:     "warning",
			Message:   "網路流量異常",
			Timestamp: time.Now().Add(-10 * time.Minute),
			Source:    "network_monitor",
			Resolved:  false,
		},
		{
			ID:        "alert_003",
			Level:     "info",
			Message:   "系統更新完成",
			Timestamp: time.Now().Add(-15 * time.Minute),
			Source:    "system",
			Resolved:  true,
		},
	}
}

// getEventsData 取得事件數據
func (ui *UIServer) getEventsData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id":        "event_001",
			"type":      "threat_detection",
			"severity":  "high",
			"message":   "偵測到惡意IP連接嘗試",
			"timestamp": time.Now().Unix() - 300,
			"data": map[string]interface{}{
				"source_ip":   "192.168.1.100",
				"threat_type": "brute_force",
			},
		},
		{
			"id":        "event_002",
			"type":      "network_event",
			"severity":  "info",
			"message":   "網路連接已恢復",
			"timestamp": time.Now().Unix() - 600,
			"data": map[string]interface{}{
				"action": "unblock",
			},
		},
		{
			"id":        "event_003",
			"type":      "device_event",
			"severity":  "medium",
			"message":   "USB裝置連接",
			"timestamp": time.Now().Unix() - 900,
			"data": map[string]interface{}{
				"device_id": "device_001",
				"port":      "/dev/ttyUSB0",
			},
		},
	}
}

// getMetricsData 取得指標數據
func (ui *UIServer) getMetricsData() map[string]interface{} {
	return map[string]interface{}{
		"network_blocked":     0,
		"device_connected":    1,
		"threats_detected":    45,
		"packets_inspected":   1000000,
		"blocked_connections": 12,
		"active_sessions":     8,
		"system_uptime":       time.Since(ui.startTime).Seconds(),
		"cpu_usage":           15.5,
		"memory_usage":        42.3,
		"disk_usage":          28.7,
	}
}

// getThreatEvents 取得威脅事件列表
func (ui *UIServer) getThreatEvents(c *gin.Context) {
	_ = c.DefaultQuery("severity", "all")
	_ = c.DefaultQuery("time_range", "24h")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// 模擬威脅事件數據
	threats := []map[string]interface{}{
		{
			"id":             "threat_001",
			"timestamp":      time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"type":           "SQL Injection Attempt",
			"severity":       "critical",
			"source_ip":      "192.168.1.100",
			"destination_ip": "10.0.0.5",
			"port":           3306,
			"protocol":       "TCP",
			"action":         "blocked",
			"description":    "偵測到 SQL 注入攻擊嘗試",
			"rule_id":        "RULE-001",
		},
		{
			"id":             "threat_002",
			"timestamp":      time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			"type":           "Brute Force Attack",
			"severity":       "high",
			"source_ip":      "192.168.1.101",
			"destination_ip": "10.0.0.6",
			"port":           22,
			"protocol":       "SSH",
			"action":         "blocked",
			"description":    "偵測到 SSH 暴力破解攻擊",
			"rule_id":        "RULE-002",
		},
		{
			"id":             "threat_003",
			"timestamp":      time.Now().Add(-3 * time.Hour).Format(time.RFC3339),
			"type":           "Port Scanning",
			"severity":       "medium",
			"source_ip":      "192.168.1.102",
			"destination_ip": "10.0.0.7",
			"port":           0,
			"protocol":       "ICMP",
			"action":         "logged",
			"description":    "偵測到端口掃描活動",
			"rule_id":        "RULE-003",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"threats": threats,
		"total":   len(threats),
		"limit":   limit,
		"offset":  offset,
	})
}

// getSecurityStats 取得安全統計
func (ui *UIServer) getSecurityStats(c *gin.Context) {
	stats := map[string]interface{}{
		"total_threats":    156,
		"blocked_threats":  142,
		"active_threats":   8,
		"resolved_threats": 148,
		"threat_trend":     -5.2,
		"top_threat_types": []map[string]interface{}{
			{"type": "SQL Injection", "count": 42},
			{"type": "Brute Force", "count": 35},
			{"type": "Port Scanning", "count": 28},
			{"type": "DDoS", "count": 22},
			{"type": "XSS", "count": 18},
		},
		"top_source_ips": []map[string]interface{}{
			{"ip": "192.168.1.100", "count": 25},
			{"ip": "192.168.1.101", "count": 18},
			{"ip": "192.168.1.102", "count": 15},
			{"ip": "192.168.1.103", "count": 12},
		},
	}

	c.JSON(http.StatusOK, stats)
}

// blockThreatSource 阻斷威脅來源
func (ui *UIServer) blockThreatSource(c *gin.Context) {
	threatID := c.Param("id")
	
	ui.logger.Infof("阻斷威脅來源: %s", threatID)
	
	ui.broadcastWebSocketMessage("threat_blocked", map[string]interface{}{
		"threat_id": threatID,
		"timestamp": time.Now().Unix(),
	})

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "威脅來源已阻斷",
		"threat_id": threatID,
		"timestamp": time.Now().Unix(),
	})
}

// getNetworkStats 取得網路統計
func (ui *UIServer) getNetworkStats(c *gin.Context) {
	stats := map[string]interface{}{
		"total_traffic":       10240000,
		"inbound_traffic":     6144000,
		"outbound_traffic":    4096000,
		"packet_loss":         0.05,
		"latency":             12.5,
		"bandwidth_usage":     65.5,
		"active_connections":  128,
		"blocked_connections": 12,
	}

	c.JSON(http.StatusOK, stats)
}

// getBlockedIPs 取得被阻斷的 IP 列表
func (ui *UIServer) getBlockedIPs(c *gin.Context) {
	blockedIPs := []map[string]interface{}{
		{
			"ip":           "192.168.1.100",
			"reason":       "SQL 注入攻擊嘗試",
			"blocked_at":   time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			"expires_at":   time.Now().Add(22 * time.Hour).Format(time.RFC3339),
			"threat_count": 25,
		},
		{
			"ip":           "192.168.1.101",
			"reason":       "SSH 暴力破解攻擊",
			"blocked_at":   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"expires_at":   time.Now().Add(23 * time.Hour).Format(time.RFC3339),
			"threat_count": 18,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"blocked_ips": blockedIPs,
	})
}

// unblockIP 解除 IP 阻斷
func (ui *UIServer) unblockIP(c *gin.Context) {
	ip := c.Param("ip")
	
	ui.logger.Infof("解除 IP 阻斷: %s", ip)
	
	ui.broadcastWebSocketMessage("ip_unblocked", map[string]interface{}{
		"ip":        ip,
		"timestamp": time.Now().Unix(),
	})

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "IP 阻斷已解除",
		"ip":        ip,
		"timestamp": time.Now().Unix(),
	})
}

// getNetworkInterfaces 取得網路介面列表
func (ui *UIServer) getNetworkInterfaces(c *gin.Context) {
	interfaces := []map[string]interface{}{
		{
			"name":        "eth0",
			"status":      "up",
			"ip_address":  "10.0.0.5",
			"mac_address": "00:1A:2B:3C:4D:5E",
			"rx_bytes":    12345678,
			"tx_bytes":    9876543,
			"rx_packets":  123456,
			"tx_packets":  98765,
		},
		{
			"name":        "eth1",
			"status":      "up",
			"ip_address":  "192.168.1.5",
			"mac_address": "00:1A:2B:3C:4D:5F",
			"rx_bytes":    5432109,
			"tx_bytes":    4321098,
			"rx_packets":  54321,
			"tx_packets":  43210,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"interfaces": interfaces,
	})
}

// getDevices 取得設備列表
func (ui *UIServer) getDevices(c *gin.Context) {
	devices := []map[string]interface{}{
		{
			"id":          "device_001",
			"name":        "USB-SERIAL CH340",
			"type":        "serial",
			"status":      "online",
			"port":        "/dev/ttyUSB0",
			"baud_rate":   115200,
			"last_seen":   time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
			"uptime":      "48h 32m",
			"rx_bytes":    1234567,
			"tx_bytes":    987654,
			"error_count": 2,
		},
		{
			"id":          "device_002",
			"name":        "Network Interface eth0",
			"type":        "network",
			"status":      "online",
			"port":        "eth0",
			"ip_address":  "10.0.0.5",
			"mac_address": "00:1A:2B:3C:4D:5E",
			"last_seen":   time.Now().Add(-1 * time.Minute).Format(time.RFC3339),
			"uptime":      "72h 15m",
			"rx_bytes":    12345678,
			"tx_bytes":    9876543,
			"error_count": 0,
		},
		{
			"id":          "device_003",
			"name":        "Sensor Module",
			"type":        "sensor",
			"status":      "offline",
			"port":        "/dev/ttyUSB1",
			"last_seen":   time.Now().Add(-30 * time.Minute).Format(time.RFC3339),
			"uptime":      "0h 0m",
			"error_count": 5,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
		"total":   len(devices),
		"online":  2,
		"offline": 1,
	})
}

// getDeviceDetail 取得設備詳情
func (ui *UIServer) getDeviceDetail(c *gin.Context) {
	deviceID := c.Param("id")
	
	device := map[string]interface{}{
		"id":          deviceID,
		"name":        "USB-SERIAL CH340",
		"type":        "serial",
		"status":      "online",
		"port":        "/dev/ttyUSB0",
		"baud_rate":   115200,
		"data_bits":   8,
		"stop_bits":   1,
		"parity":      "none",
		"last_seen":   time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		"uptime":      "48h 32m",
		"rx_bytes":    1234567,
		"tx_bytes":    987654,
		"error_count": 2,
		"config": map[string]interface{}{
			"auto_restart":   true,
			"timeout":        30,
			"retry_attempts": 3,
		},
		"statistics": map[string]interface{}{
			"total_connections":   1523,
			"successful_reads":    15230,
			"failed_reads":        2,
			"avg_response_time":   12.5,
		},
	}

	c.JSON(http.StatusOK, device)
}

// restartDevice 重啟設備
func (ui *UIServer) restartDevice(c *gin.Context) {
	deviceID := c.Param("id")
	
	ui.logger.Infof("重啟設備: %s", deviceID)
	
	ui.broadcastWebSocketMessage("device_restarting", map[string]interface{}{
		"device_id": deviceID,
		"timestamp": time.Now().Unix(),
	})

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "設備重啟命令已發送",
		"device_id": deviceID,
		"timestamp": time.Now().Unix(),
	})
}

// updateDeviceConfig 更新設備配置
func (ui *UIServer) updateDeviceConfig(c *gin.Context) {
	deviceID := c.Param("id")
	
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ui.logger.Infof("更新設備配置: %s, 配置: %v", deviceID, config)
	
	ui.broadcastWebSocketMessage("device_config_updated", map[string]interface{}{
		"device_id": deviceID,
		"config":    config,
		"timestamp": time.Now().Unix(),
	})

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "設備配置已更新",
		"device_id": deviceID,
		"config":    config,
		"timestamp": time.Now().Unix(),
	})
}

// generateSecurityReport 生成安全報表
func (ui *UIServer) generateSecurityReport(c *gin.Context) {
	timeRange := c.DefaultQuery("time_range", "24h")
	format := c.DefaultQuery("format", "json")

	report := map[string]interface{}{
		"report_type":  "security",
		"time_range":   timeRange,
		"generated_at": time.Now().Format(time.RFC3339),
		"summary": map[string]interface{}{
			"total_threats":    156,
			"blocked_threats":  142,
			"active_threats":   8,
			"resolved_threats": 148,
		},
		"threat_by_type": []map[string]interface{}{
			{"type": "SQL Injection", "count": 42, "percentage": 26.9},
			{"type": "Brute Force", "count": 35, "percentage": 22.4},
			{"type": "Port Scanning", "count": 28, "percentage": 17.9},
			{"type": "DDoS", "count": 22, "percentage": 14.1},
			{"type": "XSS", "count": 18, "percentage": 11.5},
		},
		"top_source_ips": []map[string]interface{}{
			{"ip": "192.168.1.100", "threats": 25, "severity": "high"},
			{"ip": "192.168.1.101", "threats": 18, "severity": "medium"},
		},
		"timeline": []map[string]interface{}{
			{"hour": "00:00", "threats": 5},
			{"hour": "01:00", "threats": 3},
			{"hour": "02:00", "threats": 8},
		},
	}

	if format == "csv" {
		c.Header("Content-Disposition", "attachment; filename=security_report.csv")
		c.Header("Content-Type", "text/csv")
		// 這裡應該生成實際的 CSV
		c.String(http.StatusOK, "報表數據...")
	} else {
		c.JSON(http.StatusOK, report)
	}
}

// generateNetworkReport 生成網路報表
func (ui *UIServer) generateNetworkReport(c *gin.Context) {
	timeRange := c.DefaultQuery("time_range", "24h")
	format := c.DefaultQuery("format", "json")

	report := map[string]interface{}{
		"report_type":  "network",
		"time_range":   timeRange,
		"generated_at": time.Now().Format(time.RFC3339),
		"summary": map[string]interface{}{
			"total_traffic":       "10.24 GB",
			"inbound_traffic":     "6.14 GB",
			"outbound_traffic":    "4.10 GB",
			"active_connections":  128,
			"blocked_connections": 12,
			"avg_latency":         12.5,
			"packet_loss":         0.05,
		},
		"traffic_by_hour": []map[string]interface{}{
			{"hour": "00:00", "inbound": 250, "outbound": 180},
			{"hour": "01:00", "inbound": 220, "outbound": 150},
		},
		"top_destinations": []map[string]interface{}{
			{"ip": "8.8.8.8", "traffic": "1.2 GB", "protocol": "DNS"},
			{"ip": "1.1.1.1", "traffic": "0.8 GB", "protocol": "HTTPS"},
		},
		"protocols": []map[string]interface{}{
			{"protocol": "HTTPS", "traffic": "5.2 GB", "percentage": 50.8},
			{"protocol": "HTTP", "traffic": "2.5 GB", "percentage": 24.4},
			{"protocol": "DNS", "traffic": "1.5 GB", "percentage": 14.6},
		},
	}

	if format == "csv" {
		c.Header("Content-Disposition", "attachment; filename=network_report.csv")
		c.Header("Content-Type", "text/csv")
		c.String(http.StatusOK, "報表數據...")
	} else {
		c.JSON(http.StatusOK, report)
	}
}

// generateSystemReport 生成系統報表
func (ui *UIServer) generateSystemReport(c *gin.Context) {
	timeRange := c.DefaultQuery("time_range", "24h")
	format := c.DefaultQuery("format", "json")

	report := map[string]interface{}{
		"report_type":  "system",
		"time_range":   timeRange,
		"generated_at": time.Now().Format(time.RFC3339),
		"summary": map[string]interface{}{
			"uptime":        "72h 45m",
			"cpu_usage_avg": 15.5,
			"mem_usage_avg": 42.3,
			"disk_usage":    28.7,
		},
		"resource_usage": []map[string]interface{}{
			{"timestamp": "2025-01-14 10:00", "cpu": 12.5, "memory": 40.2, "disk": 28.5},
			{"timestamp": "2025-01-14 11:00", "cpu": 15.8, "memory": 42.5, "disk": 28.7},
		},
		"services_status": []map[string]interface{}{
			{"service": "Axiom UI", "status": "healthy", "uptime": "72h 45m"},
			{"service": "Prometheus", "status": "healthy", "uptime": "72h 45m"},
			{"service": "Grafana", "status": "healthy", "uptime": "72h 45m"},
		},
		"events": []map[string]interface{}{
			{"time": "2025-01-14 10:30", "type": "system", "message": "服務重啟"},
			{"time": "2025-01-14 11:15", "type": "info", "message": "配置更新"},
		},
	}

	if format == "csv" {
		c.Header("Content-Disposition", "attachment; filename=system_report.csv")
		c.Header("Content-Type", "text/csv")
		c.String(http.StatusOK, "報表數據...")
	} else {
		c.JSON(http.StatusOK, report)
	}
}

// generateCustomReport 生成自訂報表
func (ui *UIServer) generateCustomReport(c *gin.Context) {
	var request struct {
		ReportType string                 `json:"report_type"`
		TimeRange  string                 `json:"time_range"`
		Metrics    []string               `json:"metrics"`
		Filters    map[string]interface{} `json:"filters"`
		Format     string                 `json:"format"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ui.logger.Infof("生成自訂報表: %v", request)

	report := map[string]interface{}{
		"report_type":  request.ReportType,
		"time_range":   request.TimeRange,
		"generated_at": time.Now().Format(time.RFC3339),
		"metrics":      request.Metrics,
		"filters":      request.Filters,
		"data": []map[string]interface{}{
			{"timestamp": "2025-01-14 10:00", "value": 123},
			{"timestamp": "2025-01-14 11:00", "value": 156},
		},
	}

	if request.Format == "csv" {
		c.Header("Content-Disposition", "attachment; filename=custom_report.csv")
		c.Header("Content-Type", "text/csv")
		c.String(http.StatusOK, "報表數據...")
	} else {
		c.JSON(http.StatusOK, report)
	}
}

// getSwaggerJSON 返回 Swagger JSON 文檔
func (ui *UIServer) getSwaggerJSON(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, SwaggerDoc)
}

// getSwaggerUI 返回 Swagger UI 頁面
func (ui *UIServer) getSwaggerUI(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <title>Pandora Box Console API 文檔</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui.css">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
        .topbar {
            display: none;
        }
        .swagger-ui .info {
            margin: 50px 0;
        }
        .swagger-ui .info .title {
            font-size: 36px;
            color: #0284c7;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                defaultModelsExpandDepth: 1,
                defaultModelExpandDepth: 1,
                docExpansion: "list",
                filter: true,
                showExtensions: true,
                showCommonExtensions: true
            });

            window.ui = ui;
        };
    </script>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}