package fishpi

import (
	"fmt"
	"net/http"

	"dpbug/fishpi/go-client/pkg/fishpi/models"
	"go.uber.org/zap"
)

// 查询指定用户信息
func (c *Client) GetMemberInfo(username string) (*models.User, error) {
	c.Logger.Info("获取用户信息", zap.String("username", username))

	// 构建路径
	path := fmt.Sprintf("/user/%s", username)
	if c.APIKey != "" {
		path += "?apiKey=" + c.APIKey
	}

	// 发送请求
	resp, err := c.doRequest(http.MethodGet, path, nil, false)
	if err != nil {
		return nil, err
	}

	// 直接返回的用户对象，不是嵌套在data字段中
	var user models.User
	if err := c.parseResponse(resp, &user); err != nil {
		return nil, err
	}

	c.Logger.Info("获取用户信息成功", zap.String("username", user.UserName))

	return &user, nil
}

// 获取用户活跃度
func (c *Client) GetLiveness() (float64, error) {
	if c.APIKey == "" {
		return 0, fmt.Errorf("API Key未设置，请先登录")
	}

	c.Logger.Info("获取活跃度")

	// 发送请求
	resp, err := c.doRequest(http.MethodGet, "/user/liveness", nil, true)
	if err != nil {
		return 0, err
	}

	// 解析响应
	var liveness models.Liveness
	if err := c.parseResponse(resp, &liveness); err != nil {
		return 0, err
	}

	c.Logger.Info("获取活跃度成功", zap.Float64("liveness", liveness.Liveness))

	return liveness.Liveness, nil
}

// GetCheckInStatus 获取签到状态
func (c *Client) GetCheckInStatus() (bool, error) {
	if c.APIKey == "" {
		return false, fmt.Errorf("API Key未设置，请先登录")
	}

	c.Logger.Info("获取签到状态")

	// 发送请求
	resp, err := c.doRequest(http.MethodGet, "/user/checkedIn", nil, true)
	if err != nil {
		return false, err
	}

	// 解析响应
	var checkIn models.CheckIn
	if err := c.parseResponse(resp, &checkIn); err != nil {
		return false, err
	}

	c.Logger.Info("获取签到状态成功", zap.Bool("checked_in", checkIn.CheckedIn))

	return checkIn.CheckedIn, nil
}

// ClaimYesterdayLivenessReward 领取昨日活跃奖励
func (c *Client) ClaimYesterdayLivenessReward() (int, error) {
	if c.APIKey == "" {
		return 0, fmt.Errorf("API Key未设置，请先登录")
	}

	c.Logger.Info("领取昨日活跃奖励")

	// 发送请求
	resp, err := c.doRequest(http.MethodGet, "/activity/yesterday-liveness-reward-api", nil, true)
	if err != nil {
		return 0, err
	}

	// 解析响应
	var reward models.LivenessReward
	if err := c.parseResponse(resp, &reward); err != nil {
		return 0, err
	}

	if reward.Sum == -1 {
		c.Logger.Info("昨日活跃奖励已领取")
		return -1, nil
	}

	c.Logger.Info("领取昨日活跃奖励成功", zap.Int("points", reward.Sum))

	return reward.Sum, nil
}

// IsCollectedLiveness 查询昨日奖励领取状态
func (c *Client) IsCollectedLiveness() (bool, error) {
	if c.APIKey == "" {
		return false, fmt.Errorf("API Key未设置，请先登录")
	}

	c.Logger.Info("查询昨日奖励领取状态")

	// 发送请求
	resp, err := c.doRequest(http.MethodGet, "/api/activity/is-collected-liveness", nil, true)
	if err != nil {
		return false, err
	}

	// 解析响应
	var status models.LivenessCollectedStatus
	if err := c.parseResponse(resp, &status); err != nil {
		return false, err
	}

	c.Logger.Info("查询昨日奖励领取状态成功",
		zap.Bool("is_collected", status.IsCollectedYesterdayLivenessReward))

	return status.IsCollectedYesterdayLivenessReward, nil
}

