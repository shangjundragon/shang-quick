package model

import "encoding/json"

// OnlineUser 在线用户（缓存中维护，非持久化存储）
type OnlineUser struct {
	TokenId        string `json:"tokenId"`                  // Token 哈希 ID（SHA256）
	UserId         int64  `json:"userId,string"`            // 用户 ID
	Username       string `json:"username"`                  // 用户名
	Nickname       string `json:"nickname"`                  // 昵称
	IpAddr         string `json:"ipAddr"`                    // 登录 IP
	Browser        string `json:"browser"`                   // 浏览器类型
	Os             string `json:"os"`                        // 操作系统
	LoginTime      int64  `json:"loginTime"`                 // 登录时间（毫秒时间戳）
	LastActiveTime int64  `json:"lastActiveTime"`            // 最后活跃时间（毫秒时间戳）
	ExpireTime     int64  `json:"expireTime"`                // Token 过期时间（毫秒时间戳）
}

func (o *OnlineUser) ToJSON() string {
	data, _ := json.Marshal(o)
	return string(data)
}

func OnlineUserFromJSON(jsonStr string) *OnlineUser {
	var o OnlineUser
	json.Unmarshal([]byte(jsonStr), &o)
	return &o
}
