package fishpi

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

const (
	// DefaultBaseURL 默认API基础URL
	DefaultBaseURL = "https://fishpi.cn"
	// DefaultUserAgent 默认User-Agent
	DefaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"
	// MinRequestInterval 最小请求间隔（秒）
	MinRequestInterval = 30
)

// Client 摸鱼派客户端
type Client struct {
	BaseURL       string
	HTTPClient    *http.Client
	UserAgent     string
	APIKey        string
	Logger        *zap.Logger
	lastReqByPath map[string]time.Time // 记录每个接口端点的上次请求时间
	mu            sync.Mutex           // 保护lastReqByPath的并发访问
}

// ClientOption 客户端配置选项
type ClientOption func(*Client)

// WithBaseURL 设置基础URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithUserAgent 设置User-Agent
func WithUserAgent(ua string) ClientOption {
	return func(c *Client) {
		c.UserAgent = ua
	}
}

// WithLogger 设置日志记录器
func WithLogger(logger *zap.Logger) ClientOption {
	return func(c *Client) {
		c.Logger = logger
	}
}

// WithHTTPClient 设置HTTP客户端
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// NewClient 创建新的客户端实例
func NewClient(opts ...ClientOption) *Client {
	logger, _ := zap.NewProduction()

	client := &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		UserAgent:     DefaultUserAgent,
		Logger:        logger,
		lastReqByPath: make(map[string]time.Time),
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// CommonResponse 通用响应结构
type CommonResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

func (c *Client) doRequest(method, path string, body interface{}, needsAuth bool) (*http.Response, error) {
	// 构建完整URL
	url := c.BaseURL + path

	// 序列化请求体
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// 创建请求
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", c.UserAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 如果需要认证，添加API Key到请求中
	if needsAuth && c.APIKey != "" {
		// API Key可以通过查询参数传递
		q := req.URL.Query()
		q.Add("apiKey", c.APIKey)
		req.URL.RawQuery = q.Encode()
	}

	// 请求频率控制 - 针对单个接口端点
	c.mu.Lock()
	lastReq, exists := c.lastReqByPath[path]
	if exists {
		elapsed := time.Since(lastReq)
		if elapsed < MinRequestInterval*time.Second {
			waitTime := MinRequestInterval*time.Second - elapsed
			c.mu.Unlock()
			c.Logger.Info("等待接口请求间隔",
				zap.String("path", path),
				zap.Duration("wait_time", waitTime),
			)
			time.Sleep(waitTime)
			c.mu.Lock()
		}
	}
	c.mu.Unlock()

	// 记录请求
	c.Logger.Info("发送请求",
		zap.String("method", method),
		zap.String("url", url),
		zap.String("path", path),
		zap.Bool("needs_auth", needsAuth),
	)

	// 发送请求
	resp, err := c.HTTPClient.Do(req)

	// 记录本次请求时间
	c.mu.Lock()
	c.lastReqByPath[path] = time.Now()
	c.mu.Unlock()

	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	return resp, nil
}

func (c *Client) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %w", err)
	}

	c.Logger.Debug("收到响应",
		zap.Int("status_code", resp.StatusCode),
		zap.String("body", string(body)),
	)

	if resp.StatusCode != http.StatusOK {
		c.Logger.Error("HTTP请求失败",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(body)),
		)
		return fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, target); err != nil {
		c.Logger.Error("JSON解析失败",
			zap.Error(err),
			zap.String("response", string(body)),
		)
		return fmt.Errorf("解析响应失败: %w, 响应: %s", err, string(body))
	}

	return nil
}

// SetAPIKey 设置API Key
func (c *Client) SetAPIKey(apiKey string) {
	c.APIKey = apiKey
	c.Logger.Info("API Key已设置")
}

// GetAPIKey 获取当前API Key
func (c *Client) GetAPIKey() string {
	return c.APIKey
}

// MD5Hash 计算MD5哈希（用于密码加密）
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
