package models

import (
	"encoding/json"
	"strings"
)

// ChatRoomNode 聊天室节点信息
type ChatRoomNode struct {
	Code      int    `json:"code"` // 显示用户当前所在的区服中文名称
	Msg       string `json:"msg,omitempty"`
	Data      string `json:"data,omitempty"` // 自动分配的 WebSocket 地址，请取用该地址连接到聊天室频道
	Avaliable []struct {
		Node   string `json:"node"`   // 节点 WSS 地址
		Name   string `json:"name"`   // 节点中文名称
		Online int    `json:"online"` // 节点当前在线人数
		Weight int    `json:"weight"` // 节点权重
	} `json:"avaliable,omitempty"`
	ApiKey string `json:"apiKey"` // 自动生成的 ApiKey，用于在手动选择节点时和 avaliable 中的 node 拼接，以连接到自定义节点的WebSocket服务器
}

// ChatMessage 聊天室消息
type ChatMessage struct {
	OID              string `json:"oId"`
	Type             string `json:"type"`    // 消息类型，如 "msg"
	UserOId          int64  `json:"userOId"` // 用户OID（数字类型）
	UserName         string `json:"userName"`
	UserNickname     string `json:"userNickname"`
	UserAvatarURL    string `json:"userAvatarURL"`
	Content          string `json:"content"`            // 普通消息为HTML，红包消息为JSON字符串
	Time             string `json:"time"`               // 时间（字符串格式："2025-10-29 10:49:55"）
	MD               string `json:"md"`                 // Markdown 格式内容（红包消息无此字段）
	SysMetal         string `json:"sysMetal,omitempty"`
	Client           string `json:"client,omitempty"`   // 客户端标识
	UserCardBg       string `json:"userCardBg,omitempty"`
	UserOnlineFlag   bool   `json:"userOnlineFlag,omitempty"`
	UserAvatarURL210 string `json:"userAvatarURL210,omitempty"`
	UserAvatarURL48  string `json:"userAvatarURL48,omitempty"`
	UserAvatarURL20  string `json:"userAvatarURL20,omitempty"`
}

// RedPacketContent 红包消息内容（存在 Content 字段中）
type RedPacketContent struct {
	MsgType  string `json:"msgType"`  // 固定为 "redPacket"
	Msg      string `json:"msg"`      // 红包祝福语
	SenderId string `json:"senderId"` // 发送者ID
	Recivers string `json:"recivers"` // 接收者列表（专属红包），API返回的是JSON字符串，如 "[]" 或 "[\"user1\"]"
	Money    int    `json:"money"`    // 红包总金额（积分）
	Count    int    `json:"count"`    // 红包数量
	Type     string `json:"type"`     // 红包类型: random, average, specify, heartbeat, rockPaperScissors
	Got      int    `json:"got"`      // 已领取数量
	Who      []struct {
		UserName  string `json:"userName"`
		Avatar    string `json:"avatar"`
		UserMoney int    `json:"userMoney"`
		Time      string `json:"time"`
	} `json:"who"` // 已领取者信息
}

// GetReciverList 解析接收者列表（将JSON字符串转为数组）
func (rp *RedPacketContent) GetReciverList() ([]string, error) {
	var recivers []string
	if rp.Recivers == "" {
		return recivers, nil
	}
	if err := json.Unmarshal([]byte(rp.Recivers), &recivers); err != nil {
		return nil, err
	}
	return recivers, nil
}

// IsRedPacket 判断消息是否为红包消息
func (m *ChatMessage) IsRedPacket() bool {
	// 红包消息的 Content 是 JSON 格式，且以 { 开头
	return strings.HasPrefix(strings.TrimSpace(m.Content), "{")
}

// GetRedPacket 解析红包内容
func (m *ChatMessage) GetRedPacket() (*RedPacketContent, error) {
	var rp RedPacketContent
	if err := json.Unmarshal([]byte(m.Content), &rp); err != nil {
		return nil, err
	}
	return &rp, nil
}

// FormatDisplay 格式化消息用于展示
func (m *ChatMessage) FormatDisplay() string {
	if m.IsRedPacket() {
		rp, err := m.GetRedPacket()
		if err != nil {
			return m.UserNickname + ": [红包解析失败]"
		}
		return m.UserNickname + ": [红包] " + rp.Msg + " (ID: " + m.OID + ")"
	}
	// 移除 MD 中的 emoji 和特殊字符，保留纯文本
	content := strings.ReplaceAll(m.MD, "\n", " ")
	return m.UserNickname + ": " + content
}

// ChatHistory 聊天历史记录响应
type ChatHistory struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg,omitempty"`
	Data []ChatMessage `json:"data,omitempty"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	Content string `json:"content"` // 消息内容(Markdown 格式)
	Client  string `json:"client"`  // 客户端标识
}

// SendMessageResponse 发送消息响应
type SendMessageResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

// RedPacket 红包信息
type RedPacket struct {
	OID         string `json:"oId"`
	UserName    string `json:"userName"`
	Msg         string `json:"msg"`                   // 红包祝福语
	Recivers    string `json:"recivers"`              // 接收者列表（逗号分隔）
	Count       int    `json:"count"`                 // 红包数量
	Got         int    `json:"got"`                   // 已领取数量
	Money       int    `json:"money"`                 // 红包总金额（积分）
	Type        string `json:"type"`                  // 红包类型: random(拼手气), average(平分), specify(专属), heartbeat(心跳), rockPaperScissors(猜拳)
	Who         string `json:"who,omitempty"`         // 专属红包接收者
	Time        int64  `json:"time"`                  // 时间戳
	GestureType int    `json:"gestureType,omitempty"` // 猜拳类型
}

// RedPacketInfo 红包详情
type RedPacketInfo struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg,omitempty"`
	Data RedPacket `json:"data,omitempty"`
}
