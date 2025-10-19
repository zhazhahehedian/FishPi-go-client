package models

// User 用户信息
type User struct {
	OID                  string     `json:"oId"`
	UserNo               int        `json:"userNo"`
	UserName             string     `json:"userName"`
	UserNickname         string     `json:"userNickname"`
	UserRole             string     `json:"userRole"`
	UserAvatarURL        string     `json:"userAvatarURL"`
	UserCity             string     `json:"userCity"`
	UserOnlineFlag       bool       `json:"userOnlineFlag"`
	OnlineMinute         int        `json:"onlineMinute"`
	UserPoint            int        `json:"userPoint"`
	UserAppRole          int        `json:"userAppRole"`
	UserIntro            string     `json:"userIntro"`
	UserURL              string     `json:"userURL"`
	CardBg               string     `json:"cardBg"`
	FollowingUserCount   int        `json:"followingUserCount"`
	SysMetal             string     `json:"sysMetal"`
	UserProvince         string     `json:"userProvince,omitempty"`
	UserUsedPoint        int        `json:"userUsedPoint,omitempty"`
	UserCommentCount     int        `json:"userCommentCount,omitempty"`
	UserArticleCount     int        `json:"userArticleCount,omitempty"`
	UserFollowerStatus   int        `json:"userFollowerStatus,omitempty"`
	UserCommentStatus    int        `json:"userCommentStatus,omitempty"`
	UserOnlineStatus     int        `json:"userOnlineStatus,omitempty"`
	UserUAStatus         int        `json:"userUAStatus,omitempty"`
	UserCurrentCheckinStreak int    `json:"userCurrentCheckinStreak,omitempty"`
}

// UserMetal 用户徽章
type UserMetal struct {
	List []Metal `json:"list"`
}

// Metal 徽章
type Metal struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Attr        string `json:"attr"` // 包含url、backcolor、fontcolor
	Data        string `json:"data"`
}

// Liveness 活跃度
type Liveness struct {
	Liveness float64 `json:"liveness"`
}

// CheckIn 签到信息
type CheckIn struct {
	CheckedIn bool `json:"checkedIn"`
}

// LivenessReward 活跃度奖励
type LivenessReward struct {
	Sum int `json:"sum"` // -1表示已领取
}

// LivenessCollectedStatus 昨日活跃奖励领取状态
type LivenessCollectedStatus struct {
	IsCollectedYesterdayLivenessReward bool `json:"isCollectedYesterdayLivenessReward"`
}

