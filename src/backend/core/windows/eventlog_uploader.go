package windows

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// EventLogUploader 事件日誌上傳器
type EventLogUploader struct {
	logger         *logrus.Logger
	axiomBackendURL string
	agentID        string
	computerName   string
	httpClient     *http.Client
	retryAttempts  int
	retryDelay     time.Duration
}

// NewEventLogUploader 創建事件日誌上傳器
func NewEventLogUploader(logger *logrus.Logger, axiomBackendURL, agentID string) *EventLogUploader {
	computerName, _ := os.Hostname()
	
	return &EventLogUploader{
		logger:          logger,
		axiomBackendURL: axiomBackendURL,
		agentID:         agentID,
		computerName:    computerName,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryAttempts: 3,
		retryDelay:    5 * time.Second,
	}
}

// UploadLogs 上傳日誌到 Axiom Backend
func (u *EventLogUploader) UploadLogs(logs []WindowsEventLog) error {
	if len(logs) == 0 {
		return nil
	}

	// 構建上傳請求
	requestBody := map[string]interface{}{
		"agent_id": u.agentID,
		"computer": u.computerName,
		"logs":     logs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v2/logs/windows/batch", u.axiomBackendURL)

	// 重試機制
	var lastErr error
	for i := 0; i < u.retryAttempts; i++ {
		if i > 0 {
			u.logger.Infof("Retrying upload (%d/%d)...", i+1, u.retryAttempts)
			time.Sleep(u.retryDelay)
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("create request: %w", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := u.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("execute request: %w", err)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			u.logger.Infof("Successfully uploaded %d logs to Axiom Backend", len(logs))
			return nil
		}

		// 讀取錯誤響應
		var errorResp bytes.Buffer
		errorResp.ReadFrom(resp.Body)
		resp.Body.Close()
		
		lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp.String())
	}

	return fmt.Errorf("failed after %d attempts: %w", u.retryAttempts, lastErr)
}

// SetRetryAttempts 設置重試次數
func (u *EventLogUploader) SetRetryAttempts(attempts int) {
	u.retryAttempts = attempts
}

// SetRetryDelay 設置重試延遲
func (u *EventLogUploader) SetRetryDelay(delay time.Duration) {
	u.retryDelay = delay
}



import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// EventLogUploader 事件日誌上傳器
type EventLogUploader struct {
	logger         *logrus.Logger
	axiomBackendURL string
	agentID        string
	computerName   string
	httpClient     *http.Client
	retryAttempts  int
	retryDelay     time.Duration
}

// NewEventLogUploader 創建事件日誌上傳器
func NewEventLogUploader(logger *logrus.Logger, axiomBackendURL, agentID string) *EventLogUploader {
	computerName, _ := os.Hostname()
	
	return &EventLogUploader{
		logger:          logger,
		axiomBackendURL: axiomBackendURL,
		agentID:         agentID,
		computerName:    computerName,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryAttempts: 3,
		retryDelay:    5 * time.Second,
	}
}

// UploadLogs 上傳日誌到 Axiom Backend
func (u *EventLogUploader) UploadLogs(logs []WindowsEventLog) error {
	if len(logs) == 0 {
		return nil
	}

	// 構建上傳請求
	requestBody := map[string]interface{}{
		"agent_id": u.agentID,
		"computer": u.computerName,
		"logs":     logs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v2/logs/windows/batch", u.axiomBackendURL)

	// 重試機制
	var lastErr error
	for i := 0; i < u.retryAttempts; i++ {
		if i > 0 {
			u.logger.Infof("Retrying upload (%d/%d)...", i+1, u.retryAttempts)
			time.Sleep(u.retryDelay)
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("create request: %w", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := u.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("execute request: %w", err)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			u.logger.Infof("Successfully uploaded %d logs to Axiom Backend", len(logs))
			return nil
		}

		// 讀取錯誤響應
		var errorResp bytes.Buffer
		errorResp.ReadFrom(resp.Body)
		resp.Body.Close()
		
		lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp.String())
	}

	return fmt.Errorf("failed after %d attempts: %w", u.retryAttempts, lastErr)
}

// SetRetryAttempts 設置重試次數
func (u *EventLogUploader) SetRetryAttempts(attempts int) {
	u.retryAttempts = attempts
}

// SetRetryDelay 設置重試延遲
func (u *EventLogUploader) SetRetryDelay(delay time.Duration) {
	u.retryDelay = delay
}


import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// EventLogUploader 事件日誌上傳器
type EventLogUploader struct {
	logger         *logrus.Logger
	axiomBackendURL string
	agentID        string
	computerName   string
	httpClient     *http.Client
	retryAttempts  int
	retryDelay     time.Duration
}

// NewEventLogUploader 創建事件日誌上傳器
func NewEventLogUploader(logger *logrus.Logger, axiomBackendURL, agentID string) *EventLogUploader {
	computerName, _ := os.Hostname()
	
	return &EventLogUploader{
		logger:          logger,
		axiomBackendURL: axiomBackendURL,
		agentID:         agentID,
		computerName:    computerName,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryAttempts: 3,
		retryDelay:    5 * time.Second,
	}
}

// UploadLogs 上傳日誌到 Axiom Backend
func (u *EventLogUploader) UploadLogs(logs []WindowsEventLog) error {
	if len(logs) == 0 {
		return nil
	}

	// 構建上傳請求
	requestBody := map[string]interface{}{
		"agent_id": u.agentID,
		"computer": u.computerName,
		"logs":     logs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v2/logs/windows/batch", u.axiomBackendURL)

	// 重試機制
	var lastErr error
	for i := 0; i < u.retryAttempts; i++ {
		if i > 0 {
			u.logger.Infof("Retrying upload (%d/%d)...", i+1, u.retryAttempts)
			time.Sleep(u.retryDelay)
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("create request: %w", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := u.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("execute request: %w", err)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			u.logger.Infof("Successfully uploaded %d logs to Axiom Backend", len(logs))
			return nil
		}

		// 讀取錯誤響應
		var errorResp bytes.Buffer
		errorResp.ReadFrom(resp.Body)
		resp.Body.Close()
		
		lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp.String())
	}

	return fmt.Errorf("failed after %d attempts: %w", u.retryAttempts, lastErr)
}

// SetRetryAttempts 設置重試次數
func (u *EventLogUploader) SetRetryAttempts(attempts int) {
	u.retryAttempts = attempts
}

// SetRetryDelay 設置重試延遲
func (u *EventLogUploader) SetRetryDelay(delay time.Duration) {
	u.retryDelay = delay
}



import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// EventLogUploader 事件日誌上傳器
type EventLogUploader struct {
	logger         *logrus.Logger
	axiomBackendURL string
	agentID        string
	computerName   string
	httpClient     *http.Client
	retryAttempts  int
	retryDelay     time.Duration
}

// NewEventLogUploader 創建事件日誌上傳器
func NewEventLogUploader(logger *logrus.Logger, axiomBackendURL, agentID string) *EventLogUploader {
	computerName, _ := os.Hostname()
	
	return &EventLogUploader{
		logger:          logger,
		axiomBackendURL: axiomBackendURL,
		agentID:         agentID,
		computerName:    computerName,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		retryAttempts: 3,
		retryDelay:    5 * time.Second,
	}
}

// UploadLogs 上傳日誌到 Axiom Backend
func (u *EventLogUploader) UploadLogs(logs []WindowsEventLog) error {
	if len(logs) == 0 {
		return nil
	}

	// 構建上傳請求
	requestBody := map[string]interface{}{
		"agent_id": u.agentID,
		"computer": u.computerName,
		"logs":     logs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v2/logs/windows/batch", u.axiomBackendURL)

	// 重試機制
	var lastErr error
	for i := 0; i < u.retryAttempts; i++ {
		if i > 0 {
			u.logger.Infof("Retrying upload (%d/%d)...", i+1, u.retryAttempts)
			time.Sleep(u.retryDelay)
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(jsonData))
		if err != nil {
			lastErr = fmt.Errorf("create request: %w", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := u.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("execute request: %w", err)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			resp.Body.Close()
			u.logger.Infof("Successfully uploaded %d logs to Axiom Backend", len(logs))
			return nil
		}

		// 讀取錯誤響應
		var errorResp bytes.Buffer
		errorResp.ReadFrom(resp.Body)
		resp.Body.Close()
		
		lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp.String())
	}

	return fmt.Errorf("failed after %d attempts: %w", u.retryAttempts, lastErr)
}

// SetRetryAttempts 設置重試次數
func (u *EventLogUploader) SetRetryAttempts(attempts int) {
	u.retryAttempts = attempts
}

// SetRetryDelay 設置重試延遲
func (u *EventLogUploader) SetRetryDelay(delay time.Duration) {
	u.retryDelay = delay
}

