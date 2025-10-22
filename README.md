# FishPi Go Client

摸鱼派的一个 Go 客户端程序（G门弄斧了...）

> **📢 项目说明**  
> 想着练手 Go 项目，反正闲着也是闲着，尝试搓一个看看能不能跑起来。  
> 基于摸鱼派社区开放 API 使用文档 V2.1.6 开发。

## 项目特性

- ✅ 交互式登录（支持二重验证）
- ✅ 配置文件自动管理（API Key 持久化）
- ✅ 请求频率自动控制
- ✅ 完善的日志记录

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
积分: 5000
...
```

**首次登录后，API Key 会自动保存到 `~/.fishpi/config.json`，下次运行会自动使用保存的 API Key。**

## API 功能列表

### 已实现功能 ✅

- [x] 用户认证
  - [x] 登录获取 API Key
  - [x] 验证 API Key
  - [x] 使用 API Key 登录
- [x] 用户模块
  - [x] 获取当前用户信息
  - [x] 查询其他用户信息
  - [x] 获取活跃度
  - [x] 获取签到状态
  - [x] 领取昨日活跃奖励
  - [x] 查询奖励领取状态
- [x] 配置管理
  - [x] 自动保存/加载配置
  - [x] API Key 持久化

## 配置说明

客户端会自动在用户目录创建配置文件：`~/.fishpi/config.json`

配置文件示例：

```json
{
  "base_url": "https://fishpi.cn",
  "user_agent": "Mozilla/5.0 ...",
  "api_key": "your-api-key-here"
}
```

### 计划实现的功能 🚧
- [ ] 聊天室功能
- [ ] 帖子浏览与评论
- [ ] 私信功能

## 技术栈

- Go 1.19+
- `net/http` - HTTP 客户端
- `go.uber.org/zap` - 日志库
- `golang.org/x/term` - 终端交互（密码隐藏输入）

## 学习笔记

这是一个用于学习 Go 语言的项目，主要学习内容包括：

1. ✅ Go 项目结构和模块管理 (go.mod)
2. ✅ HTTP 客户端封装和请求处理
3. ✅ 配置文件管理 (JSON)
4. ✅ 终端交互和密码安全输入
5. ✅ 错误处理和日志记录
6. 🚧 WebSocket 实时通信
7. 🚧 并发和协程管理
8. 🚧 CLI 工具开发 (cobra)

## 许可证

MIT License

## 相关链接

- [摸鱼派社区](https://fishpi.cn)
- [摸鱼派 API 文档](https://fishpi.cn/article/1636516552191)

## 致谢

感谢摸鱼派社区提供开放 API！

---

**注意**: 只是成为了调接口侠，我太菜了Orz。

---