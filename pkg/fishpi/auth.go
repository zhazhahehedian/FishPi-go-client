package fishpi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dpbug/fishpi/go-client/pkg/fishpi/models"
	"go.uber.org/zap"
)

// 登录请求
type LoginRequest struct {
	NameOrEmail  string `json:"nameOrEmail"`
	UserPassword string `json:"userPassword"` // MD5加密后的密码
	MfaCode      string `json:"mfaCode,omitempty"`
}

// 登录响应
type LoginResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Key  string `json:"Key,omitempty"`
}

// 用户信息响应
type UserResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg,omitempty"`
	Data *models.User `json:"data,omitempty"`
}

// Login 用户登录，获取API Key
// nameOrEmail: 用户名或邮箱
// password: 明文密码（函数内部会自动MD5加密）
// mfaCode: 两步验证码（如果未设置则留空）
func (c *Client) Login(nameOrEmail, password, mfaCode string) (string, error) {
	c.Logger.Info("开始登录",
		zap.String("name_or_email", nameOrEmail),
		zap.Bool("has_mfa", mfaCode != ""),
	)

	// 对密码进行MD5加密
	hashedPassword := MD5Hash(password)

	// 构建请求体
	reqBody := LoginRequest{
		NameOrEmail:  nameOrEmail,
		UserPassword: hashedPassword,
		MfaCode:      mfaCode,
	}

	// 发送请求
	resp, err := c.doRequest(http.MethodPost, "/api/getKey", reqBody, false)
	if err != nil {
		return "", err
	}

	// 解析响应
	var loginResp LoginResponse
	if err := c.parseResponse(resp, &loginResp); err != nil {
		return "", err
	}

	// 检查响应状态
	if loginResp.Code != 0 {
		return "", fmt.Errorf("登录失败: %s", loginResp.Msg)
	}

	if loginResp.Key == "" {
		return "", fmt.Errorf("登录响应中未包含API Key")
	}

	// 保存API Key
	c.SetAPIKey(loginResp.Key)

	c.Logger.Info("登录成功", zap.String("api_key", loginResp.Key[:8]+"..."))

	return loginResp.Key, nil
}

// GetUser 获取当前用户信息（同时验证API Key是否有效）
func (c *Client) GetUser() (*models.User, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("API Key未设置，请先登录")
	}

	c.Logger.Info("获取用户信息")

	// 发送请求
	resp, err := c.doRequest(http.MethodGet, "/api/user", nil, true)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var userResp UserResponse
	if err := c.parseResponse(resp, &userResp); err != nil {
		return nil, err
	}

	// 检查响应状态
	if userResp.Code != 0 {
		return nil, fmt.Errorf("获取用户信息失败: %s", userResp.Msg)
	}

	if userResp.Data == nil {
		return nil, fmt.Errorf("用户信息为空")
	}

	c.Logger.Info("获取用户信息成功",
		zap.String("username", userResp.Data.UserName),
		zap.String("nickname", userResp.Data.UserNickname),
	)

	return userResp.Data, nil
}

// ValidateAPIKey 验证API Key是否有效
func (c *Client) ValidateAPIKey() (bool, error) {
	if c.APIKey == "" {
		return false, fmt.Errorf("API Key未设置")
	}

	user, err := c.GetUser()
	if err != nil {
		return false, err
	}

	return user != nil, nil
}

// LoginWithKey 使用已有的API Key登录
func (c *Client) LoginWithKey(apiKey string) (*models.User, error) {
	c.Logger.Info("使用API Key登录")

	// 设置API Key
	c.SetAPIKey(apiKey)

	// 验证API Key
	user, err := c.GetUser()
	if err != nil {
		c.APIKey = "" // 清空无效的API Key
		return nil, fmt.Errorf("API Key无效: %w", err)
	}

	c.Logger.Info("使用API Key登录成功",
		zap.String("username", user.UserName),
	)

	return user, nil
}

// GetMetalList 解析用户徽章列表
func GetMetalList(sysMetal string) (*models.UserMetal, error) {
	if sysMetal == "" {
		return &models.UserMetal{List: []models.Metal{}}, nil
	}

	var metalList models.UserMetal
	if err := json.Unmarshal([]byte(sysMetal), &metalList); err != nil {
		return nil, fmt.Errorf("解析徽章列表失败: %w", err)
	}

	return &metalList, nil
}

