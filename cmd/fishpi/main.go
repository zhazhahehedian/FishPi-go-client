package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"dpbug/fishpi/go-client/internal/config"
	"dpbug/fishpi/go-client/pkg/fishpi"
	"dpbug/fishpi/go-client/pkg/fishpi/models"
	fishpiws "dpbug/fishpi/go-client/pkg/fishpi/websocket"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/term"
)

func main() {
	fmt.Println("🐟 摸鱼派 Go 客户端")
	fmt.Println("==================")
	fmt.Println()

	// 创建日志记录器（仅用于错误日志）
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("创建日志记录器失败: %v", err)
	}
	defer logger.Sync()

	// 创建客户端（启用静默模式，不输出调试日志）
	client := fishpi.NewClient(
		fishpi.WithLogger(logger),
		fishpi.WithSilent(true), // true-启用静默模式, false-调试模式
	)

	// 检查是否已有保存的API Key
	savedAPIKey, _ := config.GetAPIKey()
	var user *models.User
	if savedAPIKey != "" {
		fmt.Println("检测到已保存的API Key，尝试使用...")
		userData, err := client.LoginWithKey(savedAPIKey)
		if err == nil {
			fmt.Printf("✓ 使用已保存的API Key登录成功!\n")
			user = userData
			// 继续执行后续操作（活跃度、签到等）
		} else {
			fmt.Printf("⚠ 已保存的API Key无效: %v\n", err)
			fmt.Println("需要重新登录")
			fmt.Println()
		}
	}

	// 如果没有使用已保存的API Key登录，则需要交互式登录
	if user == nil {
		// 从命令行交互式获取登录凭据
		reader := bufio.NewReader(os.Stdin)

		// 获取用户名
		fmt.Print("请输入用户名: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("读取用户名失败: %v", err)
		}
		username = strings.TrimSpace(username)
		if username == "" {
			log.Fatal("用户名不能为空")
		}

		fmt.Print("请输入密码: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Fatalf("读取密码失败: %v", err)
		}
		fmt.Println() // 换行
		password := strings.TrimSpace(string(passwordBytes))
		if password == "" {
			log.Fatal("密码不能为空")
		}

		// 获取二重验证令牌（可选）
		fmt.Print("请输入二重验证令牌（如未开启请直接回车）: ")
		mfaCode, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("读取二重验证令牌失败: %v", err)
		}
		mfaCode = strings.TrimSpace(mfaCode)

		// 执行登录
		fmt.Printf("\n正在登录用户: %s\n", username)
		apiKey, err := client.Login(username, password, mfaCode)
		if err != nil {
			log.Fatalf("登录失败: %v", err)
		}

		fmt.Printf("✓ 登录成功! API Key: %s...\n", apiKey[:8])

		// 保存API Key到配置文件
		if err := config.SaveAPIKey(apiKey); err != nil {
			fmt.Printf("⚠ 警告: 保存API Key失败: %v\n", err)
		} else {
			fmt.Println("✓ API Key已保存到配置文件")
		}

		// 获取用户信息
		fmt.Println("\n正在获取用户信息...")
		user, err = client.GetUser()
		if err != nil {
			log.Fatalf("获取用户信息失败: %v", err)
		}
	}

	// 显示用户信息
	printUserInfo(user)

	// 获取活跃度
	fmt.Println("\n正在获取活跃度...")
	liveness, err := client.GetLiveness()
	if err != nil {
		fmt.Printf("⚠ 获取活跃度失败: %v\n", err)
	} else {
		fmt.Printf("✓ 当前活跃度: %.2f\n", liveness)
	}

	// 获取签到状态
	fmt.Println("\n正在获取签到状态...")
	checkedIn, err := client.GetCheckInStatus()
	if err != nil {
		fmt.Printf("⚠ 获取签到状态失败: %v\n", err)
	} else {
		if checkedIn {
			fmt.Println("✓ 今日已签到")
		} else {
			fmt.Println("× 今日未签到")
		}
	}

	// 尝试领取昨日活跃奖励
	fmt.Println("\n正在领取昨日活跃奖励...")
	reward, err := client.ClaimYesterdayLivenessReward()
	if err != nil {
		fmt.Printf("⚠ 领取昨日活跃奖励失败: %v\n", err)
	} else {
		if reward == -1 {
			fmt.Println("× 昨日活跃奖励已领取")
		} else {
			fmt.Printf("✓ 成功领取昨日活跃奖励: %d 积分\n", reward)
		}
	}

	// 显示完整统计信息
	printSummary(user, liveness, checkedIn, reward)

	// 进入主菜单
	showMainMenu(client, user)
}

func showMainMenu(client *fishpi.Client, user *models.User) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("📋 主菜单")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("1. 进入聊天室")
		fmt.Println("2. 进入清风明月")
		fmt.Println("3. 查看个人信息")
		fmt.Println("0. 退出")
		fmt.Print("\n请选择功能: ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("⚠ 读取输入失败: %v\n", err)
			continue
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			enterChatRoom(client)
		case "2":
			enterBreezemoon(client)
		case "3":
			printUserInfo(user)
		case "0":
			fmt.Println("\n👋 再见！继续摸鱼吧~")
			return
		default:
			fmt.Println("\n⚠ 无效的选择，请重新输入")
		}
	}
}

func enterChatRoom(client *fishpi.Client) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("💬 进入聊天室")
	fmt.Println(strings.Repeat("=", 50))

	// 连接 WebSocket
	fmt.Println("\n正在连接聊天室 WebSocket...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := client.ConnectChatRoom(ctx)
	if err != nil {
		fmt.Printf("⚠ 连接聊天室失败: %v\n", err)
		fmt.Println("按回车键返回主菜单...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}
	defer conn.Close()

	fmt.Println("✓ 聊天室连接成功！")
	fmt.Println("\n使用说明：")
	fmt.Println("- 直接输入文字发送消息")
	fmt.Println("- 红包会自动领取（30秒间隔，猜拳随机出拳）")
	fmt.Println("- 输入 /exit 或 /quit 退出聊天室")
	fmt.Println()

	// 启动消息接收协程（带自动领取红包）
	stopReceive := make(chan struct{})
	go receiveChatMessagesWithClient(conn, client, stopReceive)

	// 主循环：发送消息
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\n⚠ 读取输入失败: %v\n", err)
			break
		}
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		// 处理命令
		if strings.HasPrefix(input, "/") {
			if input == "/exit" || input == "/quit" {
				fmt.Println("\n👋 正在退出聊天室...")
				close(stopReceive)
				time.Sleep(100 * time.Millisecond)
				return
			} else if input == "/help" {
				fmt.Println("\n可用命令：")
				fmt.Println("  /exit, /quit - 退出聊天室")
				fmt.Println()
				continue
			} else {
				fmt.Printf("⚠ 未知命令: %s (输入 /help 查看帮助)\n", input)
				continue
			}
		}

		// 发送消息
		if err := client.SendChatMessage(input); err != nil {
			fmt.Printf("⚠ 发送消息失败: %v\n", err)
		}
	}

	close(stopReceive)
}

func receiveChatMessagesWithClient(conn *fishpiws.ChatRoomConn, client *fishpi.Client, stop chan struct{}) {
	var lastRedPacketTime time.Time
	var redPacketMu sync.Mutex
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-stop:
			return
		default:
			_, message, err := conn.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					fmt.Printf("\n⚠ WebSocket 连接断开: %v\n", err)
				}
				return
			}

			var msg models.ChatMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				fmt.Printf("\r\033[K[解析失败] 错误: %v\n原始消息: %s\n> ", err, string(message))
				continue
			}

			// 清除当前行并打印消息
			fmt.Print("\r\033[K")
			printChatMessage(&msg)

			// 自动领取红包（30s间隔限制）
			if msg.IsRedPacket() {
				redPacketMu.Lock()
				elapsed := time.Since(lastRedPacketTime)
				if elapsed < 30*time.Second {
					redPacketMu.Unlock()
					fmt.Printf("⚠ 红包领取冷却中，还需等待 %.0f 秒\n", (30*time.Second - elapsed).Seconds())
					fmt.Print("> ")
					continue
				}
				lastRedPacketTime = time.Now()
				redPacketMu.Unlock()

				// 解析红包内容以确定类型
				rp, err := msg.GetRedPacket()
				if err != nil {
					fmt.Printf("⚠ 红包解析失败，跳过自动领取: %v\n", err)
					fmt.Print("> ")
					continue
				}

				// 根据红包类型决定gesture参数
				gesture := -1
				if rp.Type == "rockPaperScissors" {
					gesture = rng.Intn(3) // 随机出拳：0=石头，1=剪刀，2=布
				}

				// 异步领取红包
				go func(oId string, g int) {
					time.Sleep(100 * time.Millisecond)
					if result, err := client.OpenRedPacket(oId, g); err == nil {
						gestureName := ""
						if g >= 0 {
							gestureNames := []string{"石头", "剪刀", "布"}
							gestureName = fmt.Sprintf("，出了%s", gestureNames[g])
						}
						fmt.Printf("\r\033[K✓ 自动领取红包成功%s！祝福语: %s\n> ", gestureName, result.Data.Msg)
					} else {
						fmt.Printf("\r\033[K⚠ 自动领取红包失败: %v\n> ", err)
					}
				}(msg.OID, gesture)
			}
			fmt.Print("> ")
		}
	}
}

func printChatMessage(msg *models.ChatMessage) {
	// 格式化时间 - time 字段已经是格式化后的字符串（如 "2025-10-29 10:49:55"）
	// 只需要提取时间部分（HH:MM:SS）
	timestamp := msg.Time
	if len(timestamp) > 10 {
		// 从 "2025-10-29 10:49:55" 提取 "10:49:55"
		timestamp = timestamp[11:]
	}
	// 如果时间超过8个字符（HH:MM:SS），只保留 HH:MM
	if len(timestamp) > 8 {
		timestamp = timestamp[:5]
	}

	// 判断是否为红包消息
	if msg.IsRedPacket() {
		// 解析红包内容
		rp, err := msg.GetRedPacket()
		if err != nil {
			// 如果解析失败，显示原始内容和错误信息
			nickname := msg.UserNickname
			if nickname == "" {
				nickname = "未知用户"
			}
			fmt.Printf("[%s] [红包] %s: [红包解析失败: %v]\n", timestamp, nickname, err)
			fmt.Printf("原始内容: %s\n", msg.Content)
			return
		}

		// 显示红包信息
		nickname := msg.UserNickname
		if nickname == "" {
			nickname = "未知用户"
		}
		redPacketType := getRedPacketTypeName(rp.Type)
		fmt.Printf("[%s] [红包] %s: %s (%s红包, %d/%d已领取, 总计%d积分)\n",
			timestamp, nickname, rp.Msg, redPacketType, rp.Got, rp.Count, rp.Money)
		return
	}

	// 普通消息 - 优先使用 Markdown 格式，如果没有则使用 HTML 内容
	content := msg.MD
	if content == "" {
		content = msg.Content
	}

	// 过滤空消息（可能是心跳、系统消息等）
	if strings.TrimSpace(content) == "" && strings.TrimSpace(msg.UserNickname) == "" {
		return
	}

	// 如果昵称为空但有内容，使用默认昵称
	nickname := msg.UserNickname
	if nickname == "" {
		nickname = "系统"
	}

	fmt.Printf("[%s] %s: %s\n", timestamp, nickname, content)
}

func getRedPacketTypeName(rpType string) string {
	switch rpType {
	case "random":
		return "拼手气"
	case "average":
		return "平分"
	case "specify":
		return "专属"
	case "heartbeat":
		return "心跳"
	case "rockPaperScissors":
		return "猜拳"
	default:
		return ""
	}
}

func printUserInfo(user interface{}) {
	var userData *fishpi.UserResponse
	switch v := user.(type) {
	case *fishpi.UserResponse:
		userData = v
	case *models.User:
		// 如果是 models.User 类型，需要转换
		userData = &fishpi.UserResponse{
			Data: v,
		}
	}

	if userData == nil || userData.Data == nil {
		fmt.Println("用户信息为空")
		return
	}

	fmt.Println("\n=== 用户信息 ===")
	fmt.Printf("用户名: %s\n", userData.Data.UserName)
	fmt.Printf("昵称: %s\n", userData.Data.UserNickname)
	fmt.Printf("用户编号: %s\n", userData.Data.UserNo)
	fmt.Printf("积分: %d\n", userData.Data.UserPoint)
	fmt.Printf("在线时长: %d 分钟\n", userData.Data.OnlineMinute)
	fmt.Printf("个人主页: %s\n", userData.Data.UserURL)
	fmt.Printf("城市: %s\n", userData.Data.UserCity)
	fmt.Printf("在线状态: %v\n", userData.Data.UserOnlineFlag)
	fmt.Printf("个性签名: %s\n", userData.Data.UserIntro)
}

func printSummary(user interface{}, liveness float64, checkedIn bool, reward int) {
	var userData *fishpi.UserResponse
	switch v := user.(type) {
	case *fishpi.UserResponse:
		userData = v
	case *models.User:
		userData = &fishpi.UserResponse{
			Data: v,
		}
	}

	if userData == nil || userData.Data == nil {
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📊 今日摸鱼总结")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("👤 用户: %s (%s)\n", userData.Data.UserNickname, userData.Data.UserName)
	fmt.Printf("💰 当前积分: %d\n", userData.Data.UserPoint)
	fmt.Printf("⚡ 活跃度: %.2f\n", liveness)

	if checkedIn {
		fmt.Printf("✅ 签到状态: 今日已签到\n")
	} else {
		fmt.Printf("❌ 签到状态: 今日未签到\n")
	}

	if reward == -1 {
		fmt.Printf("💎 昨日奖励: 已领取\n")
	} else if reward > 0 {
		fmt.Printf("💎 昨日奖励: +%d 积分 (刚刚领取)\n", reward)
	} else {
		fmt.Printf("💎 昨日奖励: 暂无可领取\n")
	}

	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🐟 继续摸鱼吧！")
	fmt.Println()
}

func enterBreezemoon(client *fishpi.Client) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("🌙 清风明月")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Println("1. 查看清风明月列表")
		fmt.Println("2. 发布清风明月")
		fmt.Println("0. 返回主菜单")
		fmt.Print("\n请选择功能: ")

		choice, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("⚠ 读取输入失败: %v\n", err)
			continue
		}
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			viewBreezemoons(client, reader)
		case "2":
			postBreezemoon(client, reader)
		case "0":
			return
		default:
			fmt.Println("\n⚠ 无效的选择，请重新输入")
		}
	}
}

func viewBreezemoons(client *fishpi.Client, reader *bufio.Reader) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("📜 清风明月列表")
	fmt.Println(strings.Repeat("=", 50))

	// 获取页码
	fmt.Print("\n请输入页码 (默认1): ")
	pageInput, _ := reader.ReadString('\n')
	pageInput = strings.TrimSpace(pageInput)
	page := 1
	if pageInput != "" {
		fmt.Sscanf(pageInput, "%d", &page)
	}

	// 获取每页数量
	fmt.Print("请输入每页显示数量 (默认20): ")
	sizeInput, _ := reader.ReadString('\n')
	sizeInput = strings.TrimSpace(sizeInput)
	size := 20
	if sizeInput != "" {
		fmt.Sscanf(sizeInput, "%d", &size)
	}

	fmt.Printf("\n正在获取第 %d 页清风明月 (每页%d条)...\n", page, size)

	result, err := client.GetBreezemoons(page, size)
	if err != nil {
		fmt.Printf("⚠ 获取清风明月列表失败: %v\n", err)
		fmt.Println("\n按回车键继续...")
		reader.ReadString('\n')
		return
	}

	if len(result.Breezemoons) == 0 {
		fmt.Println("\n暂无清风明月")
	} else {
		fmt.Printf("\n共获取到 %d 条清风明月:\n", len(result.Breezemoons))
		fmt.Println(strings.Repeat("-", 50))

		for i, bm := range result.Breezemoons {
			fmt.Printf("\n[%d] %s (@%s)\n", i+1, bm.TimeAgo, bm.BreezemoonAuthorName)
			if bm.BreezemoonCity != "" {
				fmt.Printf("📍 %s\n", bm.BreezemoonCity)
			}
			fmt.Printf("💬 %s\n", bm.BreezemoonContent)
			fmt.Println(strings.Repeat("-", 50))
		}
	}

	fmt.Println("\n按回车键继续...")
	reader.ReadString('\n')
}

func postBreezemoon(client *fishpi.Client, reader *bufio.Reader) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✍️  发布清风明月")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("\n提示: 支持 Markdown 格式")
	fmt.Println("输入内容后按回车发布 (输入空行取消):")
	fmt.Print("\n> ")

	content, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("⚠ 读取输入失败: %v\n", err)
		return
	}

	content = strings.TrimSpace(content)
	if content == "" {
		fmt.Println("\n已取消发布")
		return
	}

	fmt.Println("\n正在发布清风明月...")
	if err := client.PostBreezemoon(content); err != nil {
		fmt.Printf("⚠ 发布失败: %v\n", err)
	} else {
		fmt.Println("✓ 发布成功！")
	}

	fmt.Println("\n按回车键继续...")
	reader.ReadString('\n')
}
