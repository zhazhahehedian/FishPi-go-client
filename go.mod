module dpbug/fishpi/go-client

go 1.24.0

require go.uber.org/zap v1.27.0 // 日志库依赖

require (
	github.com/gorilla/websocket v1.5.1 // WebSocket 支持
	golang.org/x/term v0.36.0 // 终端交互依赖
)

require (
	// zap 依赖的多错误聚合
	go.uber.org/multierr v1.10.0 // indirect
	// x/term 及其他依赖使用的底层系统调用
	golang.org/x/sys v0.37.0 // indirect
)

require golang.org/x/net v0.17.0 // indirect
