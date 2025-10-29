package websocket

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	defaultHeartbeatInterval  = 3 * time.Minute
	defaultHeartbeatPayload   = "-hb-"
	defaultWriteTimeout       = 10 * time.Second
	websocketHandshakeTimeout = 10 * time.Second
)

// ChatRoomConn 封装聊天室 WebSocket 连接，及后台心跳。
type ChatRoomConn struct {
	Conn              *websocket.Conn
	logger            *zap.Logger
	heartbeatInterval time.Duration
	heartbeatPayload  []byte
	stopHeartbeat     chan struct{}
	once              sync.Once
}

// Close 关闭 WebSocket 连接, 停止心跳。
func (c *ChatRoomConn) Close() error {
	var closeErr error
	c.once.Do(func() {
		close(c.stopHeartbeat)
		closeErr = c.Conn.Close()
	})
	return closeErr
}

func (c *ChatRoomConn) startHeartbeatLoop() {
	if c.heartbeatInterval <= 0 || len(c.heartbeatPayload) == 0 {
		return
	}

	ticker := time.NewTicker(c.heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))
			if err := c.Conn.WriteMessage(websocket.TextMessage, c.heartbeatPayload); err != nil {
				c.logger.Warn("聊天室心跳发送失败", zap.Error(err))
				return
			}
		case <-c.stopHeartbeat:
			return
		}
	}
}

type chatRoomConnConfig struct {
	heartbeatInterval time.Duration
	heartbeatPayload  string
}

func defaultChatRoomConnConfig() chatRoomConnConfig {
	return chatRoomConnConfig{
		heartbeatInterval: defaultHeartbeatInterval,
		heartbeatPayload:  defaultHeartbeatPayload,
	}
}

// 自定义聊天室 ws。
type ChatRoomConnOption func(*chatRoomConnConfig)

// 覆盖发送心跳的时间间隔。
// 提供非正值将禁用自动心跳。
func WithChatRoomHeartbeatInterval(interval time.Duration) ChatRoomConnOption {
	return func(cfg *chatRoomConnConfig) {
		cfg.heartbeatInterval = interval
	}
}

// 覆盖心跳消息中使用的载荷。
func WithChatRoomHeartbeatPayload(payload string) ChatRoomConnOption {
	return func(cfg *chatRoomConnConfig) {
		cfg.heartbeatPayload = payload
	}
}

// 打开聊天室频道的 WebSocket 连接。
// 提供的 context 控制拨号超时；记得在返回的连接上调用 Close。
// wsURL 是从 /chat-room/node/get 接口返回的完整 WebSocket 地址（已包含 apiKey 参数）。
func ConnectChatRoom(ctx context.Context, wsURL, userAgent string, logger *zap.Logger, opts ...ChatRoomConnOption) (*ChatRoomConn, error) {
	cfg := defaultChatRoomConnConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	header := http.Header{}
	if ua := strings.TrimSpace(userAgent); ua != "" {
		header.Set("User-Agent", ua)
	}

	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = websocketHandshakeTimeout

	logger.Info("尝试连接聊天室 WebSocket", zap.String("url", wsURL))

	wsConn, resp, err := dialer.DialContext(ctx, wsURL, header)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		var body string
		if resp != nil {
			if data, readErr := io.ReadAll(resp.Body); readErr == nil {
				body = string(data)
			}
		}

		if body != "" {
			return nil, fmt.Errorf("连接聊天室失败: %w (返回响应: %s)", err, body)
		}
		return nil, fmt.Errorf("连接聊天室失败: %w", err)
	}

	wsConn.SetWriteDeadline(time.Time{})

	chatConn := &ChatRoomConn{
		Conn:              wsConn,
		logger:            logger,
		heartbeatInterval: cfg.heartbeatInterval,
		heartbeatPayload:  []byte(cfg.heartbeatPayload),
		stopHeartbeat:     make(chan struct{}),
	}

	defaultCloseHandler := wsConn.CloseHandler()
	wsConn.SetCloseHandler(func(code int, text string) error {
		logger.Info("聊天室连接关闭", zap.Int("code", code), zap.String("text", text))
		if defaultCloseHandler != nil {
			return defaultCloseHandler(code, text)
		}
		return nil
	})

	if cfg.heartbeatInterval > 0 && cfg.heartbeatPayload != "" {
		go chatConn.startHeartbeatLoop()
	}

	logger.Info("聊天室 WebSocket 连接成功")
	return chatConn, nil
}
