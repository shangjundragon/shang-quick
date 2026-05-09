package model

import "encoding/json"

type OnlineUser struct {
	TokenId        string `json:"tokenId"`
	UserId         int64  `json:"userId,string"`
	Username       string `json:"username"`
	Nickname       string `json:"nickname"`
	IpAddr         string `json:"ipAddr"`
	Browser        string `json:"browser"`
	Os             string `json:"os"`
	LoginTime      int64  `json:"loginTime"`
	LastActiveTime int64  `json:"lastActiveTime"`
	ExpireTime     int64  `json:"expireTime"`
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
