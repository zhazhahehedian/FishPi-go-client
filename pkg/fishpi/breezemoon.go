package fishpi

import (
	"fmt"

	"dpbug/fishpi/go-client/pkg/fishpi/models"

	"go.uber.org/zap"
)

// GetBreezemoons 获取清风明月列表
// page: 页码（从1开始）
// size: 每页显示数量
func (c *Client) GetBreezemoons(page, size int) (*models.BreezemoonListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}

	path := fmt.Sprintf("/api/breezemoons?p=%d&size=%d", page, size)
	c.Logger.Info("获取清风明月列表",
		zap.Int("page", page),
		zap.Int("size", size),
	)

	resp, err := c.doRequest("GET", path, nil, false)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result models.BreezemoonListResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	c.Logger.Info("获取清风明月列表成功",
		zap.Int("count", len(result.Breezemoons)),
	)

	return &result, nil
}

// PostBreezemoon 发布清风明月
// content: 清风明月内容（可以是Markdown格式）
func (c *Client) PostBreezemoon(content string) error {
	if content == "" {
		return fmt.Errorf("清风明月内容不能为空")
	}

	if c.APIKey == "" {
		return fmt.Errorf("API Key未设置，请先登录")
	}

	c.Logger.Info("发布清风明月",
		zap.Int("content_length", len(content)),
	)

	reqBody := map[string]interface{}{
		"apiKey":            c.APIKey,
		"breezemoonContent": content,
	}

	resp, err := c.doRequest("POST", "/breezemoon", reqBody, false)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result models.BreezemoonPostResponse
	if err := c.parseResponse(resp, &result); err != nil {
		return err
	}

	c.Logger.Info("发布清风明月成功")
	return nil
}
