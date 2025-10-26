package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"dpbug/fishpi/go-client/internal/config"
	"dpbug/fishpi/go-client/pkg/fishpi"
	"dpbug/fishpi/go-client/pkg/fishpi/models"

	"go.uber.org/zap"
	"golang.org/x/term"
)

func main() {
	fmt.Println("🐟 摸鱼派 Go 客户端")
	fmt.Println("==================")
	fmt.Println()

	// 创建日志记录器
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("创建日志记录器失败: %v", err)
	}
	defer logger.Sync()

	// 创建客户端
	client := fishpi.NewClient(
		fishpi.WithLogger(logger),
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

		// 获取密码（隐藏输入）
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
