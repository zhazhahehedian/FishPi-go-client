package fishpi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"dpbug/fishpi/go-client/pkg/fishpi/models"
	"dpbug/fishpi/go-client/pkg/fishpi/websocket"

	"go.uber.org/zap"
)

// 发送消息
// content: 消息内容（支持 Markdown 格式）
func (c *Client) SendChatMessage(content string) error {
	if c.APIKey == "" {
		return fmt.Errorf("API Key未设置，请先登录")
	}

	c.Logger.Info("发消息",
		zap.String("content", content),
		zap.String("client_name", c.ClientName))

	// 构建请求体
	reqBody := models.SendMessageRequest{
		Content: content,
		Client:  c.ClientName,
	}

	// 记录请求体（调试用）
	reqJSON, _ := json.Marshal(reqBody)
	c.Logger.Debug("请求体", zap.String("json", string(reqJSON)))

	// 发送请求
	resp, err := c.doRequest(http.MethodPost, "/chat-room/send", reqBody, true)
	if err != nil {
		return err
	}

	// 解析响应
	var sendResp models.SendMessageResponse
	if err := c.parseResponse(resp, &sendResp); err != nil {
		return err
	}

	// 检查响应状态
	if sendResp.Code != 0 {
		return fmt.Errorf("发送消息失败: %s", sendResp.Msg)
	}

	c.Logger.Info("发送消息成功")

	return nil
}

// 领取红包
// oId: 红包消息的 ID
// gesture: 猜拳红包必须参数，0=石头，1=剪刀，2=布；普通红包传-1
func (c *Client) OpenRedPacket(oId string, gesture int) (*models.RedPacketInfo, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("API Key未设置，请先登录")
	}

	c.Logger.Info("领取红包", zap.String("oId", oId), zap.Int("gesture", gesture))

	// 根据API文档，红包接口需要在请求体中包含 apiKey、oId 和 gesture
	reqBody := map[string]interface{}{
		"apiKey": c.APIKey,
		"oId":    oId,
	}
	if gesture >= 0 {
		reqBody["gesture"] = gesture
	}

	// 记录请求体（调试用）
	reqJSON, _ := json.Marshal(reqBody)
	c.Logger.Debug("领取红包请求体", zap.String("json", string(reqJSON)))

	// 注意：这里传 false，因为 apiKey 已经在请求体中了，不需要再作为查询参数
	resp, err := c.doRequest(http.MethodPost, "/chat-room/red-packet/open", reqBody, false)
	if err != nil {
		return nil, err
	}

	var result models.RedPacketInfo
	if err := c.parseResponse(resp, &result); err != nil {
		return nil, err
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("领取红包失败: %s", result.Msg)
	}

	c.Logger.Info("领取红包成功", zap.Int("points", result.Data.Money))
	return &result, nil
}

// 连接聊天室 WebSocket
// 需要先请求接口拿到websocket节点信息，随后再根据返回的节点信息去连接websocket
func (c *Client) ConnectChatRoom(ctx context.Context, opts ...websocket.ChatRoomConnOption) (*websocket.ChatRoomConn, error) {
	if c.APIKey == "" {
		return nil, fmt.Errorf("API Key未设置，请先登录")
	}

	c.Logger.Info("获取聊天室节点信息")

	// 先请求接口拿到websocket节点信息
	resp, err := c.doRequest(http.MethodGet, "/chat-room/node/get", nil, true)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var nodeResp models.ChatRoomNode
	if err := c.parseResponse(resp, &nodeResp); err != nil {
		return nil, err
	}

	// 检查响应状态
	if nodeResp.Code != 0 {
		return nil, fmt.Errorf("获取聊天室节点信息失败: %s", nodeResp.Msg)
	}

	c.Logger.Info("获取到节点信息", zap.String("node", nodeResp.Data))

	// 使用返回的 url 直接连接
	return websocket.ConnectChatRoom(ctx, nodeResp.Data, c.UserAgent, c.Logger, opts...)
}
