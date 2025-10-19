package main

import (
	"crypto/rand"
	"net/http"
	
	"github.com/gin-gonic/gin"
	
	"axiom-backend/internal/agent"
	"axiom-backend/internal/compliance"
	"axiom-backend/internal/database"
	"axiom-backend/internal/handler"
	"axiom-backend/internal/service"
	"axiom-backend/internal/storage"
)

// setupRoutes 設置所有路由
func setupRoutes(router *gin.Engine, db *database.Database, cfg *Config) {
	// ============================================
	// 初始化基礎服務
	// ============================================
	prometheusService := service.NewPrometheusService(cfg.PrometheusURL)
	lokiService := service.NewLokiService(cfg.LokiURL)
	quantumService := service.NewQuantumService(cfg.QuantumURL, db)
	nginxService := service.NewNginxService(cfg.NginxURL, cfg.NginxConfigPath)
	windowsLogService := service.NewWindowsLogService(db)
	
	// ============================================
	// 組合服務
	// ============================================
	combinedService := service.NewCombinedService(db, prometheusService, lokiService, quantumService, windowsLogService)
	
	// ============================================
	// 高級服務
	// ============================================
	timeTravelService := service.NewTimeTravelService(db, prometheusService, lokiService)
	adaptiveSecurityService := service.NewAdaptiveSecurityService(quantumService)
	selfHealingService := service.NewSelfHealingService(db, prometheusService, quantumService)
	apiGovernanceService := service.NewAPIGovernanceService(db)
	dataLineageService := service.NewDataLineageService()
	contextAwareService := service.NewContextAwareService()
	techDebtService := service.NewTechDebtService()
	agentPracticalService := service.NewAgentPracticalService()
	
	// ============================================
	// Agent 管理 (Phase 11)
	// ============================================
	agentManager := agent.NewAgentManager(db)
	
	// ============================================
	// Storage 管理 (Phase 12)
	// ============================================
	tieringPipeline := storage.NewTieringPipeline(db.Redis, db.PG)
	
	// 啟動自動分層管道
	go tieringPipeline.Start()
	
	// ============================================
	// Compliance 引擎 (Phase 13)
	// ============================================
	piiDetector := compliance.NewPIIDetector()
	
	// 生成32字節的加密密鑰
	encryptionKey := make([]byte, 32)
	rand.Read(encryptionKey)
	
	anonymizer := compliance.NewAnonymizer("axiom-salt-key", encryptionKey)
	gdprService := compliance.NewGDPRService(db.PG, anonymizer)
	
	// ============================================
	// 初始化所有處理器
	// ============================================
	prometheusHandler := handler.NewPrometheusHandler(prometheusService)
	lokiHandler := handler.NewLokiHandler(lokiService)
	quantumHandler := handler.NewQuantumHandler(quantumService)
	nginxHandler := handler.NewNginxHandler(nginxService)
	windowsLogHandler := handler.NewWindowsLogHandler(windowsLogService)
	combinedHandler := handler.NewCombinedHandler(combinedService)
	timeTravelHandler := handler.NewTimeTravelHandler(timeTravelService)
	adaptiveSecurityHandler := handler.NewAdaptiveSecurityHandler(adaptiveSecurityService)
	selfHealingHandler := handler.NewSelfHealingHandler(selfHealingService)
	apiGovernanceHandler := handler.NewAPIGovernanceHandler(apiGovernanceService)
	dataLineageHandler := handler.NewDataLineageHandler(dataLineageService)
	contextAwareHandler := handler.NewContextAwareHandler(contextAwareService)
	techDebtHandler := handler.NewTechDebtHandler(techDebtService)
	agentPracticalHandler := handler.NewAgentPracticalHandler(agentPracticalService)
	
	// Phase 11: Agent 處理器
	agentHandler := handler.NewAgentHandler(agentManager)
	
	// Phase 12: Storage 處理器
	storageHandler := handler.NewStorageHandler(tieringPipeline)
	
	// Phase 13: Compliance 處理器
	complianceHandler := handler.NewComplianceHandler(piiDetector, anonymizer)
	gdprHandler := handler.NewGDPRHandler(gdprService)
	
	// ============================================
	// API v2 路由
	// ============================================
	v2 := router.Group("/api/v2")
	{
		// ========== 基礎服務 APIs ==========
		
		// Prometheus 路由
		prom := v2.Group("/prometheus")
		{
			prom.GET("/health", prometheusHandler.HealthCheck)
			prom.GET("/status", prometheusHandler.GetStatus)
			prom.POST("/query", prometheusHandler.Query)
			prom.POST("/query-range", prometheusHandler.QueryRange)
			prom.GET("/rules", prometheusHandler.GetAlertRules)
			prom.GET("/targets", prometheusHandler.GetTargets)
		}

		// Loki 路由
		loki := v2.Group("/loki")
		{
			loki.GET("/health", lokiHandler.HealthCheck)
			loki.GET("/query", lokiHandler.QueryLogs)
			loki.GET("/labels", lokiHandler.GetLabels)
			loki.GET("/labels/:label/values", lokiHandler.GetLabelValues)
		}

		// Quantum 路由
		quantum := v2.Group("/quantum")
		{
			quantum.GET("/health", quantumHandler.HealthCheck)
			quantum.POST("/qkd/generate", quantumHandler.GenerateQKD)
			quantum.POST("/qsvm/classify", quantumHandler.ClassifyQSVM)
			quantum.POST("/zerotrust/predict", quantumHandler.PredictZeroTrust)
			quantum.GET("/jobs", quantumHandler.ListJobs)
			quantum.GET("/jobs/:jobId", quantumHandler.GetJob)
			quantum.GET("/stats", quantumHandler.GetStats)
		}

		// Nginx 路由
		nginx := v2.Group("/nginx")
		{
			nginx.GET("/status", nginxHandler.GetStatus)
			nginx.GET("/config", nginxHandler.GetConfig)
			nginx.PUT("/config", nginxHandler.UpdateConfig)
			nginx.POST("/reload", nginxHandler.Reload)
		}

		// Windows Logs 路由
		logs := v2.Group("/logs/windows")
		{
			logs.POST("/batch", windowsLogHandler.BatchReceive)
			logs.GET("", windowsLogHandler.Query)
			logs.GET("/stats", windowsLogHandler.GetStats)
		}
		
		// ========== Phase 11: Agent 管理 APIs ==========
		
		agentRoutes := v2.Group("/agent")
		{
			agentRoutes.POST("/register", agentHandler.RegisterAgent)
			agentRoutes.POST("/heartbeat", agentHandler.Heartbeat)
			agentRoutes.GET("/list", agentHandler.ListAgents)
			agentRoutes.GET("/:agentId/status", agentHandler.GetAgent)
			agentRoutes.PUT("/:agentId/config", agentHandler.UpdateConfig)
			agentRoutes.DELETE("/:agentId", agentHandler.DeregisterAgent)
			agentRoutes.GET("/health", agentHandler.CheckHealth)
			
			// Agent 實用功能
			practical := agentRoutes.Group("/practical")
			{
				practical.POST("/discover-assets", agentPracticalHandler.DiscoverAssets)
				practical.POST("/check-compliance", agentPracticalHandler.CheckCompliance)
				practical.POST("/execute-command", agentPracticalHandler.ExecuteRemoteCommand)
				practical.GET("/execution/:executionId", agentPracticalHandler.GetExecutionStatus)
			}
		}
		
		// ========== Phase 12: Storage APIs ==========
		
		storageRoutes := v2.Group("/storage")
		{
			storageRoutes.GET("/tiers/stats", storageHandler.GetTierStats)
			storageRoutes.POST("/tier/transfer", storageHandler.TriggerTransfer)
		}
		
		// ========== Phase 13: Compliance APIs ==========
		
		complianceRoutes := v2.Group("/compliance")
		{
			// PII 檢測與匿名化
			pii := complianceRoutes.Group("/pii")
			{
				pii.POST("/detect", complianceHandler.DetectPII)
				pii.POST("/anonymize", complianceHandler.AnonymizeData)
				pii.POST("/depseudonymize", complianceHandler.Depseudonymize)
				pii.GET("/types", complianceHandler.GetSupportedPIITypes)
			}
			
			// GDPR
			gdpr := complianceRoutes.Group("/gdpr")
			{
				gdpr.POST("/deletion-request", gdprHandler.CreateDeletionRequest)
				gdpr.GET("/deletion-request/list", gdprHandler.ListDeletionRequests)
				gdpr.POST("/deletion-request/:requestId/approve", gdprHandler.ApproveDeletionRequest)
				gdpr.POST("/deletion-request/:requestId/execute", gdprHandler.ExecuteDeletion)
				gdpr.GET("/deletion-request/:requestId/verify", gdprHandler.VerifyDeletion)
				gdpr.POST("/data-export", gdprHandler.ExportData)
			}
		}
		
		// ========== Combined APIs ==========
		
		combined := v2.Group("/combined")
		{
			combined.POST("/incident/investigate", combinedHandler.InvestigateIncident)
			combined.POST("/performance/analyze", combinedHandler.AnalyzePerformance)
			combined.GET("/observability/dashboard/unified", combinedHandler.GetUnifiedObservability)
			combined.POST("/alerts/intelligent-grouping", combinedHandler.IntelligentAlertGrouping)
			combined.POST("/compliance/full-audit", combinedHandler.FullComplianceAudit)
			
			// Self Healing
			combined.POST("/self-healing/remediate", selfHealingHandler.Remediate)
			combined.GET("/self-healing/success-rate", selfHealingHandler.GetSuccessRate)
		}
		
		// ========== Time Travel APIs ==========
		
		timeTravel := v2.Group("/time-travel")
		{
			timeTravel.POST("/snapshot/create", timeTravelHandler.CreateSnapshot)
			timeTravel.GET("/snapshot/:snapshotId", timeTravelHandler.GetSnapshot)
			timeTravel.GET("/snapshot/compare", timeTravelHandler.CompareSnapshots)
			timeTravel.POST("/what-if-analysis", timeTravelHandler.WhatIfAnalysis)
		}
		
		// ========== Adaptive Security APIs ==========
		
		adaptiveSecurity := v2.Group("/adaptive-security")
		{
			risk := adaptiveSecurity.Group("/risk")
			{
				risk.POST("/calculate", adaptiveSecurityHandler.CalculateRisk)
			}
			
			access := adaptiveSecurity.Group("/access")
			{
				access.POST("/evaluate", adaptiveSecurityHandler.EvaluateAccess)
				access.GET("/trust-score/:entityId", adaptiveSecurityHandler.GetTrustScore)
			}
			
			honeypot := adaptiveSecurity.Group("/honeypot")
			{
				honeypot.POST("/deploy", adaptiveSecurityHandler.DeployHoneypot)
				honeypot.GET("/:honeypotId/interactions", adaptiveSecurityHandler.GetHoneypotInteractions)
				honeypot.POST("/analyze-attacker", adaptiveSecurityHandler.AnalyzeAttacker)
			}
		}
		
		// ========== API Governance APIs ==========
		
		governance := v2.Group("/governance")
		{
			governance.GET("/api-health/:apiPath", apiGovernanceHandler.GetAPIHealth)
			governance.GET("/api-usage-analytics", apiGovernanceHandler.GetUsageAnalytics)
		}
		
		// ========== Data Lineage APIs ==========
		
		dataLineage := v2.Group("/data-lineage")
		{
			dataLineage.POST("/trace", dataLineageHandler.TraceDataLineage)
			dataLineage.POST("/impact-analysis", dataLineageHandler.AnalyzeImpact)
		}
		
		// ========== Context Aware APIs ==========
		
		contextAware := v2.Group("/context-aware")
		{
			contextAware.POST("/alert-routing", contextAwareHandler.RouteAlert)
		}
		
		// ========== Tech Debt APIs ==========
		
		techDebt := v2.Group("/tech-debt")
		{
			techDebt.POST("/scan", techDebtHandler.ScanTechDebt)
			techDebt.POST("/remediation-roadmap", techDebtHandler.GenerateRoadmap)
		}
	}

	// Swagger UI 路由
	router.GET("/swagger", func(c *gin.Context) {
		swaggerHTML := `
<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Axiom Backend V3 API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui.css" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.9.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/swagger.json',
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
                validatorUrl: null,
                tryItOutEnabled: true,
                supportedSubmitMethods: ['get', 'post', 'put', 'delete', 'patch'],
                docExpansion: "list",
                apisSorter: "alpha",
                operationsSorter: "alpha",
                tagsSorter: "alpha"
            });
        };
    </script>
</body>
</html>`
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, swaggerHTML)
	})

	// Swagger JSON 規格
	router.GET("/swagger.json", func(c *gin.Context) {
		swaggerSpec := gin.H{
			"swagger": "2.0",
			"info": gin.H{
				"title":       "Axiom Backend V3 API",
				"version":     "3.1.0",
				"description": "企業級統一管理平台 API 文檔 - 包含 50+ 端點",
				"contact": gin.H{
					"name": "Axiom Team",
					"url":  "https://github.com/Local_IPS-IDS",
				},
			},
			"host":     "localhost:3001",
			"basePath": "/",
			"schemes":  []string{"http", "https"},
			"consumes": []string{"application/json"},
			"produces": []string{"application/json"},
			"paths": gin.H{
				// 基礎健康檢查
				"/health": gin.H{
					"get": gin.H{
						"tags":        []string{"Health"},
						"summary":     "健康檢查",
						"description": "檢查服務健康狀態",
						"responses": gin.H{
							"200": gin.H{"description": "服務正常", "schema": gin.H{"$ref": "#/definitions/HealthResponse"}},
						},
					},
				},
				
				// ========== Prometheus APIs ==========
				"/api/v2/prometheus/health": gin.H{
					"get": gin.H{
						"tags":        []string{"Prometheus"},
						"summary":     "Prometheus 健康檢查",
						"responses": gin.H{"200": gin.H{"description": "Prometheus 服務狀態"}},
					},
				},
				"/api/v2/prometheus/status": gin.H{
					"get": gin.H{
						"tags":        []string{"Prometheus"},
						"summary":     "Prometheus 狀態",
						"responses": gin.H{"200": gin.H{"description": "Prometheus 詳細狀態"}},
					},
				},
				"/api/v2/prometheus/query": gin.H{
					"post": gin.H{
						"tags":        []string{"Prometheus"},
						"summary":     "Prometheus 查詢",
						"parameters":  []gin.H{{"name": "body", "in": "body", "required": true, "schema": gin.H{"type": "object"}}},
						"responses":   gin.H{"200": gin.H{"description": "查詢結果"}},
					},
				},
				"/api/v2/prometheus/query-range": gin.H{
					"post": gin.H{
						"tags":        []string{"Prometheus"},
						"summary":     "Prometheus 範圍查詢",
						"responses":   gin.H{"200": gin.H{"description": "範圍查詢結果"}},
					},
				},
				"/api/v2/prometheus/rules": gin.H{
					"get": gin.H{
						"tags":        []string{"Prometheus"},
						"summary":     "獲取告警規則",
						"responses":   gin.H{"200": gin.H{"description": "告警規則列表"}},
					},
				},
				"/api/v2/prometheus/targets": gin.H{
					"get": gin.H{
						"tags":        []string{"Prometheus"},
						"summary":     "獲取監控目標",
						"responses":   gin.H{"200": gin.H{"description": "監控目標列表"}},
					},
				},

				// ========== Loki APIs ==========
				"/api/v2/loki/health": gin.H{
					"get": gin.H{
						"tags":        []string{"Loki"},
						"summary":     "Loki 健康檢查",
						"responses":   gin.H{"200": gin.H{"description": "Loki 服務狀態"}},
					},
				},
				"/api/v2/loki/query": gin.H{
					"get": gin.H{
						"tags":        []string{"Loki"},
						"summary":     "日誌查詢",
						"responses":   gin.H{"200": gin.H{"description": "日誌查詢結果"}},
					},
				},
				"/api/v2/loki/labels": gin.H{
					"get": gin.H{
						"tags":        []string{"Loki"},
						"summary":     "獲取標籤",
						"responses":   gin.H{"200": gin.H{"description": "可用標籤列表"}},
					},
				},
				"/api/v2/loki/labels/{label}/values": gin.H{
					"get": gin.H{
						"tags":        []string{"Loki"},
						"summary":     "獲取標籤值",
						"parameters":  []gin.H{{"name": "label", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "標籤值列表"}},
					},
				},

				// ========== Quantum APIs ==========
				"/api/v2/quantum/health": gin.H{
					"get": gin.H{
						"tags":        []string{"Quantum"},
						"summary":     "Quantum 健康檢查",
						"responses":   gin.H{"200": gin.H{"description": "Quantum 服務狀態"}},
					},
				},
				"/api/v2/quantum/qkd/generate": gin.H{
					"post": gin.H{
						"tags":        []string{"Quantum"},
						"summary":     "生成量子金鑰分發",
						"responses":   gin.H{"200": gin.H{"description": "QKD 生成結果"}},
					},
				},
				"/api/v2/quantum/qsvm/classify": gin.H{
					"post": gin.H{
						"tags":        []string{"Quantum"},
						"summary":     "量子支援向量機分類",
						"responses":   gin.H{"200": gin.H{"description": "QSVM 分類結果"}},
					},
				},
				"/api/v2/quantum/zerotrust/predict": gin.H{
					"post": gin.H{
						"tags":        []string{"Quantum"},
						"summary":     "零信任預測",
						"responses":   gin.H{"200": gin.H{"description": "零信任預測結果"}},
					},
				},
				"/api/v2/quantum/jobs": gin.H{
					"get": gin.H{
						"tags":        []string{"Quantum"},
						"summary":     "列出量子任務",
						"responses":   gin.H{"200": gin.H{"description": "任務列表"}},
					},
				},
				"/api/v2/quantum/jobs/{jobId}": gin.H{
					"get": gin.H{
						"tags":        []string{"Quantum"},
						"summary":     "獲取量子任務詳情",
						"parameters":  []gin.H{{"name": "jobId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "任務詳情"}},
					},
				},
				"/api/v2/quantum/stats": gin.H{
					"get": gin.H{
						"tags":        []string{"Quantum"},
						"summary":     "量子服務統計",
						"responses":   gin.H{"200": gin.H{"description": "統計資訊"}},
					},
				},

				// ========== Nginx APIs ==========
				"/api/v2/nginx/status": gin.H{
					"get": gin.H{
						"tags":        []string{"Nginx"},
						"summary":     "Nginx 狀態",
						"responses":   gin.H{"200": gin.H{"description": "Nginx 狀態資訊"}},
					},
				},
				"/api/v2/nginx/config": gin.H{
					"get": gin.H{
						"tags":        []string{"Nginx"},
						"summary":     "獲取 Nginx 配置",
						"responses":   gin.H{"200": gin.H{"description": "配置資訊"}},
					},
					"put": gin.H{
						"tags":        []string{"Nginx"},
						"summary":     "更新 Nginx 配置",
						"responses":   gin.H{"200": gin.H{"description": "配置更新成功"}},
					},
				},
				"/api/v2/nginx/reload": gin.H{
					"post": gin.H{
						"tags":        []string{"Nginx"},
						"summary":     "重新載入 Nginx",
						"responses":   gin.H{"200": gin.H{"description": "重新載入成功"}},
					},
				},

				// ========== Windows Logs APIs ==========
				"/api/v2/logs/windows/batch": gin.H{
					"post": gin.H{
						"tags":        []string{"Windows Logs"},
						"summary":     "批量接收 Windows 日誌",
						"description": "批量接收 Windows 事件日誌",
						"parameters": []gin.H{
							{
								"name":        "body",
								"in":          "body",
								"required":    true,
								"schema": gin.H{
									"type": "object",
									"properties": gin.H{
										"agent_id": gin.H{"type": "string", "example": "agent-001"},
										"computer": gin.H{"type": "string", "example": "WORKSTATION-01"},
										"logs": gin.H{
											"type": "array",
											"items": gin.H{
												"type": "object",
												"properties": gin.H{
													"log_type": gin.H{"type": "string", "example": "System"},
													"source": gin.H{"type": "string", "example": "Service Control Manager"},
													"event_id": gin.H{"type": "integer", "example": 7036},
													"level": gin.H{"type": "string", "example": "Information"},
													"message": gin.H{"type": "string", "example": "Service started successfully"},
													"time_created": gin.H{"type": "string", "format": "date-time"},
												},
												"required": []string{"log_type"},
											},
											"minItems": 1,
											"maxItems": 1000,
										},
									},
									"required": []string{"agent_id", "logs"},
								},
							},
						},
						"responses": gin.H{"200": gin.H{"description": "日誌接收成功"}},
					},
				},
				"/api/v2/logs/windows": gin.H{
					"get": gin.H{
						"tags":        []string{"Windows Logs"},
						"summary":     "查詢 Windows 日誌",
						"description": "查詢 Windows 事件日誌",
						"parameters": []gin.H{
							{
								"name":        "agent_id",
								"in":          "query",
								"type":        "string",
								"description": "Agent ID",
								"required":    false,
							},
							{
								"name":        "log_type",
								"in":          "query",
								"type":        "string",
								"description": "日誌類型 (System, Security, Application)",
								"required":    false,
							},
							{
								"name":        "level",
								"in":          "query",
								"type":        "string",
								"description": "日誌級別 (Information, Warning, Error, Critical)",
								"required":    false,
							},
							{
								"name":        "page",
								"in":          "query",
								"type":        "integer",
								"description": "頁碼",
								"required":    false,
								"default":     1,
							},
							{
								"name":        "page_size",
								"in":          "query",
								"type":        "integer",
								"description": "每頁數量",
								"required":    false,
								"default":     20,
							},
						},
						"responses": gin.H{"200": gin.H{"description": "日誌查詢結果"}},
					},
				},
				"/api/v2/logs/windows/stats": gin.H{
					"get": gin.H{
						"tags":        []string{"Windows Logs"},
						"summary":     "Windows 日誌統計",
						"responses":   gin.H{"200": gin.H{"description": "統計資訊"}},
					},
				},

				// ========== Agent APIs ==========
				"/api/v2/agent/register": gin.H{
					"post": gin.H{
						"tags":        []string{"Agent"},
						"summary":     "Agent 註冊",
						"parameters":  []gin.H{{"name": "body", "in": "body", "required": true, "schema": gin.H{"$ref": "#/definitions/AgentRegistration"}}},
						"responses":   gin.H{"200": gin.H{"description": "註冊成功", "schema": gin.H{"$ref": "#/definitions/AgentResponse"}}},
					},
				},
				"/api/v2/agent/heartbeat": gin.H{
					"post": gin.H{
						"tags":        []string{"Agent"},
						"summary":     "Agent 心跳",
						"responses":   gin.H{"200": gin.H{"description": "心跳成功"}},
					},
				},
				"/api/v2/agent/list": gin.H{
					"get": gin.H{
						"tags":        []string{"Agent"},
						"summary":     "列出所有 Agent",
						"responses":   gin.H{"200": gin.H{"description": "Agent 列表"}},
					},
				},
				"/api/v2/agent/{agentId}/status": gin.H{
					"get": gin.H{
						"tags":        []string{"Agent"},
						"summary":     "獲取 Agent 狀態",
						"parameters":  []gin.H{{"name": "agentId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "Agent 狀態"}},
					},
				},
				"/api/v2/agent/{agentId}/config": gin.H{
					"put": gin.H{
						"tags":        []string{"Agent"},
						"summary":     "更新 Agent 配置",
						"parameters":  []gin.H{{"name": "agentId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "配置更新成功"}},
					},
				},
				"/api/v2/agent/{agentId}": gin.H{
					"delete": gin.H{
						"tags":        []string{"Agent"},
						"summary":     "註銷 Agent",
						"parameters":  []gin.H{{"name": "agentId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "註銷成功"}},
					},
				},
				"/api/v2/agent/health": gin.H{
					"get": gin.H{
						"tags":        []string{"Agent"},
						"summary":     "Agent 健康檢查",
						"responses":   gin.H{"200": gin.H{"description": "健康狀態"}},
					},
				},

				// ========== Agent Practical APIs ==========
				"/api/v2/agent/practical/discover-assets": gin.H{
					"post": gin.H{
						"tags":        []string{"Agent Practical"},
						"summary":     "資產發現",
						"responses":   gin.H{"200": gin.H{"description": "資產發現結果"}},
					},
				},
				"/api/v2/agent/practical/check-compliance": gin.H{
					"post": gin.H{
						"tags":        []string{"Agent Practical"},
						"summary":     "合規性檢查",
						"responses":   gin.H{"200": gin.H{"description": "合規性檢查結果"}},
					},
				},
				"/api/v2/agent/practical/execute-command": gin.H{
					"post": gin.H{
						"tags":        []string{"Agent Practical"},
						"summary":     "執行遠端命令",
						"responses":   gin.H{"200": gin.H{"description": "命令執行結果"}},
					},
				},
				"/api/v2/agent/practical/execution/{executionId}": gin.H{
					"get": gin.H{
						"tags":        []string{"Agent Practical"},
						"summary":     "獲取執行狀態",
						"parameters":  []gin.H{{"name": "executionId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "執行狀態"}},
					},
				},

				// ========== Storage APIs ==========
				"/api/v2/storage/tiers/stats": gin.H{
					"get": gin.H{
						"tags":        []string{"Storage"},
						"summary":     "儲存層統計",
						"responses":   gin.H{"200": gin.H{"description": "統計資訊"}},
					},
				},
				"/api/v2/storage/tier/transfer": gin.H{
					"post": gin.H{
						"tags":        []string{"Storage"},
						"summary":     "觸發資料轉移",
						"responses":   gin.H{"200": gin.H{"description": "轉移成功"}},
					},
				},

				// ========== Compliance APIs ==========
				"/api/v2/compliance/pii/detect": gin.H{
					"post": gin.H{
						"tags":        []string{"Compliance"},
						"summary":     "PII 檢測",
						"parameters":  []gin.H{{"name": "body", "in": "body", "required": true, "schema": gin.H{"$ref": "#/definitions/PIIDetection"}}},
						"responses":   gin.H{"200": gin.H{"description": "檢測結果"}},
					},
				},
				"/api/v2/compliance/pii/anonymize": gin.H{
					"post": gin.H{
						"tags":        []string{"Compliance"},
						"summary":     "PII 匿名化",
						"responses":   gin.H{"200": gin.H{"description": "匿名化結果"}},
					},
				},
				"/api/v2/compliance/pii/depseudonymize": gin.H{
					"post": gin.H{
						"tags":        []string{"Compliance"},
						"summary":     "PII 去偽名化",
						"responses":   gin.H{"200": gin.H{"description": "去偽名化結果"}},
					},
				},
				"/api/v2/compliance/pii/types": gin.H{
					"get": gin.H{
						"tags":        []string{"Compliance"},
						"summary":     "獲取支援的 PII 類型",
						"responses":   gin.H{"200": gin.H{"description": "PII 類型列表"}},
					},
				},

				// ========== GDPR APIs ==========
				"/api/v2/compliance/gdpr/deletion-request": gin.H{
					"post": gin.H{
						"tags":        []string{"GDPR"},
						"summary":     "創建刪除請求",
						"responses":   gin.H{"200": gin.H{"description": "請求創建成功"}},
					},
				},
				"/api/v2/compliance/gdpr/deletion-request/list": gin.H{
					"get": gin.H{
						"tags":        []string{"GDPR"},
						"summary":     "列出刪除請求",
						"responses":   gin.H{"200": gin.H{"description": "請求列表"}},
					},
				},
				"/api/v2/compliance/gdpr/deletion-request/{requestId}/approve": gin.H{
					"post": gin.H{
						"tags":        []string{"GDPR"},
						"summary":     "批准刪除請求",
						"parameters":  []gin.H{{"name": "requestId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "批准成功"}},
					},
				},
				"/api/v2/compliance/gdpr/deletion-request/{requestId}/execute": gin.H{
					"post": gin.H{
						"tags":        []string{"GDPR"},
						"summary":     "執行刪除",
						"parameters":  []gin.H{{"name": "requestId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "刪除成功"}},
					},
				},
				"/api/v2/compliance/gdpr/deletion-request/{requestId}/verify": gin.H{
					"get": gin.H{
						"tags":        []string{"GDPR"},
						"summary":     "驗證刪除",
						"parameters":  []gin.H{{"name": "requestId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "驗證結果"}},
					},
				},
				"/api/v2/compliance/gdpr/data-export": gin.H{
					"post": gin.H{
						"tags":        []string{"GDPR"},
						"summary":     "資料匯出",
						"responses":   gin.H{"200": gin.H{"description": "匯出成功"}},
					},
				},

				// ========== Combined APIs ==========
				"/api/v2/combined/incident/investigate": gin.H{
					"post": gin.H{
						"tags":        []string{"Combined"},
						"summary":     "事件調查",
						"description": "一鍵事件調查，整合 Loki、Prometheus、AlertManager、Agent、AI",
						"parameters": []gin.H{
							{
								"name":        "body",
								"in":          "body",
								"required":    true,
								"schema": gin.H{
									"type": "object",
									"properties": gin.H{
										"alert_id": gin.H{"type": "string", "example": "alert-123456"},
										"time_range": gin.H{"type": "string", "example": "1h", "default": "1h"},
									},
									"required": []string{"alert_id"},
								},
							},
						},
						"responses": gin.H{"200": gin.H{"description": "調查結果"}},
					},
				},
				"/api/v2/combined/performance/analyze": gin.H{
					"post": gin.H{
						"tags":        []string{"Combined"},
						"summary":     "效能分析",
						"responses":   gin.H{"200": gin.H{"description": "分析結果"}},
					},
				},
				"/api/v2/combined/observability/dashboard/unified": gin.H{
					"get": gin.H{
						"tags":        []string{"Combined"},
						"summary":     "統一可觀測性儀表板",
						"responses":   gin.H{"200": gin.H{"description": "儀表板資料"}},
					},
				},
				"/api/v2/combined/alerts/intelligent-grouping": gin.H{
					"post": gin.H{
						"tags":        []string{"Combined"},
						"summary":     "智慧告警分組",
						"responses":   gin.H{"200": gin.H{"description": "分組結果"}},
					},
				},
				"/api/v2/combined/compliance/full-audit": gin.H{
					"post": gin.H{
						"tags":        []string{"Combined"},
						"summary":     "完整合規性稽核",
						"responses":   gin.H{"200": gin.H{"description": "稽核結果"}},
					},
				},
				"/api/v2/combined/self-healing/remediate": gin.H{
					"post": gin.H{
						"tags":        []string{"Self Healing"},
						"summary":     "自我修復",
						"description": "智能診斷並自動修復系統問題",
						"parameters": []gin.H{
							{
								"name":        "body",
								"in":          "body",
								"required":    true,
								"schema": gin.H{
									"type": "object",
									"properties": gin.H{
										"incident_type": gin.H{"type": "string", "example": "high_cpu_usage"},
										"parameters": gin.H{
											"type": "object",
											"additionalProperties": true,
											"example": gin.H{"threshold": 80, "action": "restart_service"},
										},
									},
									"required": []string{"incident_type"},
								},
							},
						},
						"responses": gin.H{"200": gin.H{"description": "修復結果"}},
					},
				},
				"/api/v2/combined/self-healing/success-rate": gin.H{
					"get": gin.H{
						"tags":        []string{"Self Healing"},
						"summary":     "自我修復成功率",
						"responses":   gin.H{"200": gin.H{"description": "成功率統計"}},
					},
				},

				// ========== Time Travel APIs ==========
				"/api/v2/time-travel/snapshot/create": gin.H{
					"post": gin.H{
						"tags":        []string{"Time Travel"},
						"summary":     "創建快照",
						"description": "創建系統狀態快照",
						"parameters": []gin.H{
							{
								"name":        "body",
								"in":          "body",
								"required":    true,
								"schema": gin.H{
									"type": "object",
									"properties": gin.H{
										"name": gin.H{"type": "string", "example": "backup-2025-10-16"},
										"description": gin.H{"type": "string", "example": "系統備份快照"},
									},
									"required": []string{"name"},
								},
							},
						},
						"responses": gin.H{"200": gin.H{"description": "快照創建成功"}},
					},
				},
				"/api/v2/time-travel/snapshot/{snapshotId}": gin.H{
					"get": gin.H{
						"tags":        []string{"Time Travel"},
						"summary":     "獲取快照",
						"parameters":  []gin.H{{"name": "snapshotId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "快照詳情"}},
					},
				},
				"/api/v2/time-travel/snapshot/compare": gin.H{
					"get": gin.H{
						"tags":        []string{"Time Travel"},
						"summary":     "比較快照",
						"responses":   gin.H{"200": gin.H{"description": "比較結果"}},
					},
				},
				"/api/v2/time-travel/what-if-analysis": gin.H{
					"post": gin.H{
						"tags":        []string{"Time Travel"},
						"summary":     "假設分析",
						"description": "執行 What-If 分析",
						"parameters": []gin.H{
							{
								"name":        "body",
								"in":          "body",
								"required":    true,
								"schema": gin.H{
									"type": "object",
									"properties": gin.H{
										"scenario": gin.H{"type": "string", "example": "increase_load_50_percent"},
										"parameters": gin.H{
											"type": "object",
											"additionalProperties": true,
											"example": gin.H{"load_multiplier": 1.5, "duration": "30m"},
										},
									},
									"required": []string{"scenario"},
								},
							},
						},
						"responses": gin.H{"200": gin.H{"description": "分析結果"}},
					},
				},

				// ========== Adaptive Security APIs ==========
				"/api/v2/adaptive-security/risk/calculate": gin.H{
					"post": gin.H{
						"tags":        []string{"Adaptive Security"},
						"summary":     "風險計算",
						"responses":   gin.H{"200": gin.H{"description": "風險評估結果"}},
					},
				},
				"/api/v2/adaptive-security/access/evaluate": gin.H{
					"post": gin.H{
						"tags":        []string{"Adaptive Security"},
						"summary":     "存取評估",
						"responses":   gin.H{"200": gin.H{"description": "評估結果"}},
					},
				},
				"/api/v2/adaptive-security/access/trust-score/{entityId}": gin.H{
					"get": gin.H{
						"tags":        []string{"Adaptive Security"},
						"summary":     "獲取信任分數",
						"parameters":  []gin.H{{"name": "entityId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "信任分數"}},
					},
				},
				"/api/v2/adaptive-security/honeypot/deploy": gin.H{
					"post": gin.H{
						"tags":        []string{"Adaptive Security"},
						"summary":     "部署蜜罐",
						"responses":   gin.H{"200": gin.H{"description": "部署成功"}},
					},
				},
				"/api/v2/adaptive-security/honeypot/{honeypotId}/interactions": gin.H{
					"get": gin.H{
						"tags":        []string{"Adaptive Security"},
						"summary":     "獲取蜜罐互動",
						"parameters":  []gin.H{{"name": "honeypotId", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "互動記錄"}},
					},
				},
				"/api/v2/adaptive-security/honeypot/analyze-attacker": gin.H{
					"post": gin.H{
						"tags":        []string{"Adaptive Security"},
						"summary":     "分析攻擊者",
						"responses":   gin.H{"200": gin.H{"description": "分析結果"}},
					},
				},

				// ========== API Governance APIs ==========
				"/api/v2/governance/api-health/{apiPath}": gin.H{
					"get": gin.H{
						"tags":        []string{"API Governance"},
						"summary":     "API 健康檢查",
						"parameters":  []gin.H{{"name": "apiPath", "in": "path", "required": true, "type": "string"}},
						"responses":   gin.H{"200": gin.H{"description": "API 健康狀態"}},
					},
				},
				"/api/v2/governance/api-usage-analytics": gin.H{
					"get": gin.H{
						"tags":        []string{"API Governance"},
						"summary":     "API 使用分析",
						"responses":   gin.H{"200": gin.H{"description": "使用統計"}},
					},
				},

				// ========== Data Lineage APIs ==========
				"/api/v2/data-lineage/trace": gin.H{
					"post": gin.H{
						"tags":        []string{"Data Lineage"},
						"summary":     "資料血緣追蹤",
						"responses":   gin.H{"200": gin.H{"description": "追蹤結果"}},
					},
				},
				"/api/v2/data-lineage/impact-analysis": gin.H{
					"post": gin.H{
						"tags":        []string{"Data Lineage"},
						"summary":     "影響分析",
						"responses":   gin.H{"200": gin.H{"description": "分析結果"}},
					},
				},

				// ========== Context Aware APIs ==========
				"/api/v2/context-aware/alert-routing": gin.H{
					"post": gin.H{
						"tags":        []string{"Context Aware"},
						"summary":     "告警路由",
						"responses":   gin.H{"200": gin.H{"description": "路由結果"}},
					},
				},

				// ========== Tech Debt APIs ==========
				"/api/v2/tech-debt/scan": gin.H{
					"post": gin.H{
						"tags":        []string{"Tech Debt"},
						"summary":     "技術債務掃描",
						"responses":   gin.H{"200": gin.H{"description": "掃描結果"}},
					},
				},
				"/api/v2/tech-debt/remediation-roadmap": gin.H{
					"post": gin.H{
						"tags":        []string{"Tech Debt"},
						"summary":     "生成修復路線圖",
						"responses":   gin.H{"200": gin.H{"description": "路線圖"}},
					},
				},
			},
			"definitions": gin.H{
				"HealthResponse": gin.H{
					"type": "object",
					"properties": gin.H{
						"status":  gin.H{"type": "string", "example": "healthy"},
						"service": gin.H{"type": "string", "example": "axiom-backend-v3"},
						"version": gin.H{"type": "string", "example": "3.1.0"},
						"time":    gin.H{"type": "string", "format": "date-time"},
					},
				},
				"AgentRegistration": gin.H{
					"type": "object",
					"properties": gin.H{
						"mode":         gin.H{"type": "string", "example": "internal"},
						"hostname":     gin.H{"type": "string", "example": "test-server"},
						"ip_address":   gin.H{"type": "string", "example": "127.0.0.1"},
						"capabilities": gin.H{"type": "array", "items": gin.H{"type": "string"}},
					},
					"required": []string{"mode", "hostname", "ip_address"},
				},
				"AgentResponse": gin.H{
					"type": "object",
					"properties": gin.H{
						"success": gin.H{"type": "boolean", "example": true},
						"data": gin.H{
							"type": "object",
							"properties": gin.H{
								"agent_id": gin.H{"type": "string", "example": "agent_123456"},
							},
						},
					},
				},
				"PIIDetection": gin.H{
					"type": "object",
					"properties": gin.H{
						"text": gin.H{"type": "string", "example": "Contact: test@example.com, Card: 4532-1234-5678-9010"},
					},
					"required": []string{"text"},
				},
				"Error": gin.H{
					"type": "object",
					"properties": gin.H{
						"code":    gin.H{"type": "string"},
						"message": gin.H{"type": "string"},
						"details": gin.H{"type": "object"},
					},
				},
			},
		}
		c.JSON(http.StatusOK, swaggerSpec)
	})
}

