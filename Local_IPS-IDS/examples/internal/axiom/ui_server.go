package axiom

import (
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
}

// DashboardData 儀表板資料結構
type DashboardData struct {
	SystemStatus    SystemStatus    `json:"system_status"`
	SecurityMetrics SecurityMetrics `json:"security_metrics"`
	NetworkStatus   NetworkStatus   `json:"network_status"`
	RecentEvents    []Event         `json:"recent_events"`
	Timestamp       int64           `json:"timestamp"`
}

// SystemStatus 系統狀態
type SystemStatus struct {
	Uptime          int64  `json:"uptime"`
	DeviceConnected bool   `json:"device_connected"`
	NetworkBlocked  bool   `json:"network_blocked"`
	ActiveSessions  int    `json:"active_sessions"`
	Status          string `json:"status"`
}

// SecurityMetrics 安全指標
type SecurityMetrics struct {
	ThreatsDetected    int64   `json:"threats_detected"`
	PacketsInspected   int64   `json:"packets_inspected"`
	BlockedConnections int     `json:"blocked_connections"`
	SecurityScore      float64 `json:"security_score"`
}

// NetworkStatus 網路狀態
type NetworkStatus struct {
	InboundTraffic  float64 `json:"inbound_traffic"`
	OutboundTraffic float64 `json:"outbound_traffic"`
	Latency         float64 `json:"latency"`
	PacketLoss      float64 `json:"packet_loss"`
}

// Event 事件結構
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Severity  string                 `json:"severity"`
	Message   string                 `json:"message"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
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

	// 靜態檔案服務
	router.Static("/static", "./web/static")
	router.StaticFile("/", "./web/index.html")
	router.StaticFile("/favicon.ico", "./web/favicon.ico")

	// API 路由
	api := router.Group("/api/v1")
	{
		api.GET("/dashboard", ui.getDashboardData)
		api.GET("/events", ui.getEvents)
		api.GET("/metrics", ui.getMetrics)
		api.POST("/control/network", ui.controlNetwork)
		api.POST("/control/device", ui.controlDevice)
		api.GET("/status", ui.getSystemStatus)
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
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// getDashboardData 取得儀表板資料
func (ui *UIServer) getDashboardData(c *gin.Context) {
	data := DashboardData{
		SystemStatus: SystemStatus{
			Uptime:          time.Now().Unix() - 1000000, // 示例數據
			DeviceConnected: true,
			NetworkBlocked:  false,
			ActiveSessions:  5,
			Status:          "healthy",
		},
		SecurityMetrics: SecurityMetrics{
			ThreatsDetected:    42,
			PacketsInspected:   1000000,
			BlockedConnections: 15,
			SecurityScore:      95.5,
		},
		NetworkStatus: NetworkStatus{
			InboundTraffic:  1024.5,
			OutboundTraffic: 512.3,
			Latency:         12.5,
			PacketLoss:      0.1,
		},
		RecentEvents: ui.getRecentEvents(),
		Timestamp:    time.Now().Unix(),
	}

	c.JSON(http.StatusOK, data)
}

// getEvents 取得事件列表
func (ui *UIServer) getEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	eventType := c.Query("type")

	events := ui.getFilteredEvents(limit, offset, eventType)

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"total":  len(events),
		"limit":  limit,
		"offset": offset,
	})
}

// getMetrics 取得指標數據
func (ui *UIServer) getMetrics(c *gin.Context) {
	// 這裡應該從Prometheus獲取實際指標
	metrics := map[string]interface{}{
		"network_blocked":     0,
		"device_connected":    1,
		"threats_detected":    42,
		"packets_inspected":   1000000,
		"blocked_connections": 15,
		"active_sessions":     5,
		"system_uptime":       time.Now().Unix() - 1000000,
	}

	c.JSON(http.StatusOK, metrics)
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
		// 這裡應該調用實際的網路管理器
		ui.broadcastWebSocketMessage("network_status", map[string]interface{}{
			"blocked": true,
			"action":  "block",
		})
	case "unblock":
		ui.logger.Info("透過UI請求解除網路阻斷")
		// 這裡應該調用實際的網路管理器
		ui.broadcastWebSocketMessage("network_status", map[string]interface{}{
			"blocked": false,
			"action":  "unblock",
		})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的動作"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "action": request.Action})
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

	// 這裡應該調用實際的裝置管理器
	ui.broadcastWebSocketMessage("device_status", map[string]interface{}{
		"action": request.Action,
		"data":   request.Data,
	})

	c.JSON(http.StatusOK, gin.H{"status": "success", "action": request.Action})
}

// getSystemStatus 取得系統狀態
func (ui *UIServer) getSystemStatus(c *gin.Context) {
	status := SystemStatus{
		Uptime:          time.Now().Unix() - 1000000,
		DeviceConnected: true,
		NetworkBlocked:  false,
		ActiveSessions:  5,
		Status:          "healthy",
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
		clientID = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	ui.websocketConns[clientID] = conn
	defer delete(ui.websocketConns, clientID)

	ui.logger.Infof("WebSocket客戶端連接: %s", clientID)

	// 發送初始數據
	if err := ui.sendWebSocketMessage(conn, "connected", map[string]interface{}{
		"client_id": clientID,
		"timestamp": time.Now().Unix(),
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
		// 處理訂閱請求
		ui.logger.Debugf("客戶端 %s 訂閱: %v", clientID, message.Data)
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
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		ui.sendPeriodicUpdates()
	}
}

// sendPeriodicUpdates 發送定期更新
func (ui *UIServer) sendPeriodicUpdates() {
	data := DashboardData{
		SystemStatus: SystemStatus{
			Uptime:          time.Now().Unix() - 1000000,
			DeviceConnected: true,
			NetworkBlocked:  false,
			ActiveSessions:  5,
			Status:          "healthy",
		},
		SecurityMetrics: SecurityMetrics{
			ThreatsDetected:    42,
			PacketsInspected:   1000000,
			BlockedConnections: 15,
			SecurityScore:      95.5,
		},
		NetworkStatus: NetworkStatus{
			InboundTraffic:  1024.5,
			OutboundTraffic: 512.3,
			Latency:         12.5,
			PacketLoss:      0.1,
		},
		RecentEvents: ui.getRecentEvents(),
		Timestamp:    time.Now().Unix(),
	}

	ui.broadcastWebSocketMessage("dashboard_update", data)
}

// getRecentEvents 取得最近事件
func (ui *UIServer) getRecentEvents() []Event {
	// 示例事件數據
	events := []Event{
		{
			ID:        "1",
			Type:      "threat_detection",
			Severity:  "high",
			Message:   "偵測到惡意IP連接嘗試",
			Timestamp: time.Now().Unix() - 300,
			Data: map[string]interface{}{
				"source_ip":   "192.168.1.100",
				"threat_type": "brute_force",
			},
		},
		{
			ID:        "2",
			Type:      "network_event",
			Severity:  "info",
			Message:   "網路連接已恢復",
			Timestamp: time.Now().Unix() - 600,
			Data: map[string]interface{}{
				"action": "unblock",
			},
		},
	}

	return events
}

// getFilteredEvents 取得過濾後的事件
func (ui *UIServer) getFilteredEvents(limit, offset int, eventType string) []Event {
	allEvents := ui.getRecentEvents()

	// 簡單的過濾和分頁邏輯
	filtered := make([]Event, 0)
	for _, event := range allEvents {
		if eventType == "" || event.Type == eventType {
			filtered = append(filtered, event)
		}
	}

	start := offset
	if start > len(filtered) {
		start = len(filtered)
	}

	end := start + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end]
}
