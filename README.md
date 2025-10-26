# FishPi Go Client

摸鱼派的一个 Go 客户端程序（G门弄斧了...）

> **📢 项目说明**  
> 想着练手 Go 项目，反正闲着也是闲着，尝试搓一个看看能不能跑起来。  
> 基于摸鱼派社区开放 API 使用文档 V2.1.6 开发。

## 📊 项目进度

**当前阶段**: 第一、二阶段完成 ✅ 第三阶段准备中 🚧

- ✅ **阶段一**: 项目初始化与基础架构 (100%)
- ✅ **阶段二**: 用户功能模块 (100%)
- 🚧 **阶段三**: 聊天室功能 (规划中)
- 📋 **阶段四**: 清风明月 (待开发)
- 📋 **阶段五**: CLI 工具开发 (待开发)

**完成日期**: 2025-10-26 | **版本**: v0.2.0-alpha

## 🎯 项目特性

### 已实现功能 ✅
- ✅ **认证系统**
  - 用户名密码登录
  - API Key 自动获取和验证
  - MFA 双因素认证支持
  - API Key 持久化存储和自动复用
  
- ✅ **用户模块**
  - 用户信息查询和展示
  - 实时活跃度获取（支持 API 频率控制）
  - 签到状态检查
  - 昨日活跃奖励自动领取
  - 摸鱼总结统计面板
  
- ✅ **基础设施**
  - HTTP 客户端封装
  - 请求频率自动控制（防止被限流）
  - 日志记录（zap）
  - 配置文件自动管理
  - 友好的命令行输出（emoji + 格式化）

### 计划实现的功能 🚧
- 🚧 聊天室实时通信（WebSocket）
- 🚧 清风明月（查看+发布）
- 🚧 红包功能（精简版）
- 🚧 CLI 命令行工具（cobra）

## 快速开始

### 环境要求

- Go 1.19 或更高版本

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

## 📋 API 功能列表

### 已实现 API ✅

#### 认证模块
- [x] `POST /api/getKey` - 登录获取 API Key
- [x] `GET /api/user` - 验证 API Key 并获取用户信息
- [x] 支持 MFA 双因素认证

#### 用户模块
- [x] `GET /api/user` - 获取当前用户信息
- [x] `GET /user/<username>` - 查询其他用户信息
- [x] `GET /user/liveness` - 获取活跃度（含频率控制）
- [x] `GET /user/checkedIn` - 获取签到状态
- [x] `GET /activity/yesterday-liveness-reward-api` - 领取昨日活跃奖励
- [x] `GET /api/activity/is-collected-liveness` - 查询奖励领取状态

#### 配置管理
- [x] API Key 自动保存和加载
- [x] 配置文件持久化
- [x] 跨平台路径支持

### 计划实现 API 🚧

#### 聊天室模块
- [ ] `GET /chat-room/more` - 获取聊天历史
- [ ] `POST /chat-room/send` - 发送消息
- [ ] `DELETE /chat-room/revoke/<oId>` - 撤回消息
- [ ] `WSS /chat-room-channel` - WebSocket 实时连接
- [ ] `POST /chat-room/red-packet/open` - 打开红包

#### 清风明月模块
- [ ] `GET /api/breezemoons` - 获取清风明月列表
- [ ] `POST /breezemoon` - 发布清风明月

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

如果你有任何建议或发现了 bug，请随时提出。

## 📜 许可证

MIT License

## 🔗 相关链接

- [摸鱼派社区](https://fishpi.cn)
- [摸鱼派 API 文档](https://fishpi.cn/article/1636516552191)

**注意**: 只是成为了调接口侠，我太菜了Orz。

---