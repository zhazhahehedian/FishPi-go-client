module dpbug/fishpi/go-client

go 1.24.0

require go.uber.org/zap v1.24.0 // 日志库依赖

require golang.org/x/term v0.36.0 // 终端交互依赖

require (
	// zap 依赖的原子类型封装
	go.uber.org/atomic v1.7.0 // indirect
	// zap 依赖的多错误聚合
	go.uber.org/multierr v1.6.0 // indirect
	// x/term 及其他依赖使用的底层系统调用
	golang.org/x/sys v0.37.0 // indirect
)
