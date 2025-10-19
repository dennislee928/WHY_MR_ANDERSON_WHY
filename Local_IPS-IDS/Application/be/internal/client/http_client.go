package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient HTTP 客戶端封裝
type HTTPClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// Config HTTP 客戶端配置
type Config struct {
	BaseURL        string
	Timeout        time.Duration
	Headers        map[string]string
	MaxIdleConns   int
	IdleConnTimeout time.Duration
}

// NewHTTPClient 創建新的 HTTP 客戶端
func NewHTTPClient(cfg *Config) *HTTPClient {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 100
	}
	if cfg.IdleConnTimeout == 0 {
		cfg.IdleConnTimeout = 90 * time.Second
	}

	transport := &http.Transport{
		MaxIdleConns:        cfg.MaxIdleConns,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     cfg.IdleConnTimeout,
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout:   cfg.Timeout,
			Transport: transport,
		},
		baseURL: cfg.BaseURL,
		headers: cfg.Headers,
	}
}

// Get 發送 GET 請求
func (c *HTTPClient) Get(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil)
}

// Post 發送 POST 請求
func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPost, path, body)
}

// Put 發送 PUT 請求
func (c *HTTPClient) Put(ctx context.Context, path string, body interface{}) ([]byte, error) {
	return c.doRequest(ctx, http.MethodPut, path, body)
}

// Delete 發送 DELETE 請求
func (c *HTTPClient) Delete(ctx context.Context, path string) ([]byte, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil)
}

// doRequest 執行 HTTP 請求
func (c *HTTPClient) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// 設置默認 headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 設置自定義 headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetJSON 發送 GET 請求並解析 JSON
func (c *HTTPClient) GetJSON(ctx context.Context, path string, result interface{}) error {
	data, err := c.Get(ctx, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

// PostJSON 發送 POST 請求並解析 JSON
func (c *HTTPClient) PostJSON(ctx context.Context, path string, body, result interface{}) error {
	data, err := c.Post(ctx, path, body)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(data, result)
	}
	return nil
}

// SetHeader 設置 HTTP Header
func (c *HTTPClient) SetHeader(key, value string) {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[key] = value
}

// SetBaseURL 設置基礎 URL
func (c *HTTPClient) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

