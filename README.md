# FishPi Go Client

摸鱼派的一个 Go 客户端程序（搬门弄斧了...）

> **📢 项目说明**  
> 想着练手 Go 项目，反正闲着也是闲着，尝试搓一个看看能不能跑起来。  
> 基于摸鱼派社区开放 API 使用文档 V2.1.6 开发。

## 🎯 项目概况

### 已实现功能 ✅
- ✅ **登录**
  - 用户名密码登录，API Key本地存储复用

- ✅ **用户信息**
  - 用户信息查询和展示
  - 签到状态检查
  - 昨日活跃奖励自动领取
  - 摸鱼总结统计面板

- ✅ **聊天室**
  - 实时消息和发送聊天消息
  - WebSocket 及 自动心跳机制（3 分钟间隔）
  - 红包自动领取（支持猜拳红包，也是3分钟间隔）

- ✅ **清风明月**
  - 获取清风明月列表（支持分页）
  - 发布清风明月（支持 Markdown）

- ✅ **基础框架**
  - HTTP 客户端封装
  - WebSocket
  - 请求频率控制
  - 日志记录（zap）- 支持静默模式和调试模式
  - 配置文件自动管理（JSON）

### 计划实现/待解决的功能 🚧
- 🚧 聊天室吞消息的问题（找了AI也看不出来什么问题）
- 🚧 CLI 命令行工具（cobra）, 直接进入聊天室

## 快速开始

### 环境要求

- Go 1.19 或更高版本

### 日志模式

客户端支持两种日志模式：

**实际使用模式（默认）**：
- 只显示用户友好的输出信息
- 不显示技术性的 HTTP 请求/响应日志
- 适合日常使用

**调试模式**：
- 显示详细的请求/响应日志
- 适合开发和调试

如需启用调试模式，修改 [cmd/fishpi/main.go](cmd/fishpi/main.go#L31-L42)：

```go
// 调试模式：显示详细日志
logger, err := zap.NewDevelopment()
client := fishpi.NewClient(
    fishpi.WithLogger(logger),
    fishpi.WithSilent(false), // 关闭静默模式
)
```

### 编译运行

```bash
# 克隆项目
git clone https://github.com/zhazhahehedian/fishpi-go-client.git
cd FishPi-go-client

# 安装依赖
go mod download

# 编译程序
# Windows:
go build -o fishpi.exe ./cmd/fishpi

# Linux/macOS:
go build -o fishpi ./cmd/fishpi

# 运行程序
./fishpi        # Linux/macOS
.\fishpi.exe    # Windows
```

或者直接运行（无需编译）：

```bash
go run ./cmd/fishpi
```

### 使用方式

#### 首次登录

运行程序后，会提示你输入登录信息：

```
🐟 摸鱼派 Go 客户端
==================

请输入用户名: your_username
请输入密码: [隐藏输入]
请输入二重验证令牌（如未开启请直接回车）: 

正在登录用户: your_username
✓ 登录成功! API Key: oXTQTD4l...
✓ API Key已保存到配置文件

正在获取用户信息...

=== 用户信息 ===
用户名: your_username
昵称: 你的昵称
用户编号: 12345
积分: 5000
在线时长: 1234 分钟
个人主页: https://fishpi.cn/member/your_username
城市: 北京
在线状态: true
个性签名: 摸鱼使我快乐

正在获取活跃度...
✓ 当前活跃度: 85.23

正在获取签到状态...
✓ 今日已签到

正在领取昨日活跃奖励...
✓ 成功领取昨日活跃奖励: 300 积分

==================================================
📊 今日摸鱼总结
==================================================
👤 用户: 你的昵称 (your_username)
💰 当前积分: 5300
⚡ 活跃度: 85.23
✅ 签到状态: 今日已签到
💎 昨日奖励: +300 积分 (刚刚领取)
==================================================
🐟 继续摸鱼吧！
```

#### 后续使用

**首次登录后，API Key 会自动保存到 `~/.fishpi/config.yaml`，下次运行会自动使用保存的 API Key，无需再次输入密码。**

```
🐟 摸鱼派 Go 客户端
==================

检测到已保存的API Key，尝试使用...
✓ 使用已保存的API Key登录成功!

=== 用户信息 ===
...
```

## 📋 功能列表

### 已实现 ✅

**认证模块**
- ✅ `POST /api/getKey` - 登录获取 API Key
- ✅ `GET /api/user` - 验证 API Key 并获取用户信息
- ✅ 支持 MFA 双因素认证

**用户模块**
- ✅ `GET /api/user` - 获取当前用户信息
- ✅ `GET /user/liveness` - 获取活跃度
- ✅ `GET /user/checkedIn` - 获取签到状态
- ✅ `GET /activity/yesterday-liveness-reward-api` - 领取昨日活跃奖励

**聊天室模块**
- ✅ `GET /chat-room/node/get` - 获取 WebSocket 节点
- ✅ `POST /chat-room/send` - 发送聊天消息
- ✅ `POST /chat-room/red-packet/open` - 领取红包
- ✅ WebSocket 实时连接 - 消息接收和显示

**清风明月模块**
- ✅ `GET /api/breezemoons` - 获取清风明月列表
- ✅ `POST /breezemoon` - 发布清风明月

**配置管理**
- ✅ API Key 自动保存和加载
- ✅ 配置文件持久化（`~/.fishpi/config.json`）

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

如果你有任何建议或发现了 bug，请随时提出。

## 📜 许可证

MIT License

## 🔗 相关链接

- [摸鱼派社区](https://fishpi.cn)
- [摸鱼派 API 文档](https://fishpi.cn/article/1636516552191)

---