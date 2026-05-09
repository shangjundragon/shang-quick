package service

import (
	"backend/model"
	"backend/pkg/cache"
	"backend/pkg/global_vars"
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	onlineUserPrefix     = "online_user:active:"
	blacklistPrefix      = "online_user:blacklist:"
	onlineUserTTL        = 30 * time.Minute
	blacklistPrefixCheck = "online_user:blacklist"
)

type onlineUserService struct{}

var OnlineUserService = &onlineUserService{}

func getOnlineUserKey(tokenId string) string {
	return onlineUserPrefix + tokenId
}

func getBlacklistKey(tokenId string) string {
	return blacklistPrefix + tokenId
}

func (s *onlineUserService) RecordLogin(ctx context.Context, token string, userId int64, username, nickname string) {
	if cache.GlobalCache == nil {
		return
	}

	tokenId := s.extractTokenId(token)
	now := time.Now().UnixMilli()
	expireHours := global_vars.ConfigYml.GetInt("Jwt.ExpireHours")
	if expireHours <= 0 {
		expireHours = 24
	}
	expireTime := now + int64(expireHours)*3600*1000

	onlineUser := &model.OnlineUser{
		TokenId:        tokenId,
		UserId:         userId,
		Username:       username,
		Nickname:       nickname,
		LoginTime:      now,
		LastActiveTime: now,
		ExpireTime:     expireTime,
	}

	cache.GlobalCache.Set(ctx, getOnlineUserKey(tokenId), onlineUser.ToJSON(), onlineUserTTL)
}

func (s *onlineUserService) UpdateActiveTime(ctx context.Context, token string) {
	if cache.GlobalCache == nil {
		return
	}

	tokenId := s.extractTokenId(token)
	key := getOnlineUserKey(tokenId)
	data, _ := cache.GlobalCache.Get(ctx, key)
	if data == "" {
		return
	}

	onlineUser := model.OnlineUserFromJSON(data)
	onlineUser.LastActiveTime = time.Now().UnixMilli()
	cache.GlobalCache.Set(ctx, key, onlineUser.ToJSON(), onlineUserTTL)
}

func (s *onlineUserService) SetUserAgent(ctx context.Context, token string, ipAddr, browser, os string) {
	if cache.GlobalCache == nil {
		return
	}

	tokenId := s.extractTokenId(token)
	key := getOnlineUserKey(tokenId)
	data, _ := cache.GlobalCache.Get(ctx, key)
	if data == "" {
		return
	}

	onlineUser := model.OnlineUserFromJSON(data)
	onlineUser.IpAddr = ipAddr
	onlineUser.Browser = browser
	onlineUser.Os = os
	cache.GlobalCache.Set(ctx, key, onlineUser.ToJSON(), onlineUserTTL)
}

func (s *onlineUserService) IsBlacklisted(ctx context.Context, token string) bool {
	if cache.GlobalCache == nil {
		return false
	}

	tokenId := s.extractTokenId(token)
	data, _ := cache.GlobalCache.Get(ctx, getBlacklistKey(tokenId))
	return data != ""
}

func (s *onlineUserService) KickUser(ctx context.Context, tokenId string) error {
	if cache.GlobalCache == nil {
		return fmt.Errorf("缓存未初始化")
	}

	key := getOnlineUserKey(tokenId)
	data, _ := cache.GlobalCache.Get(ctx, key)
	if data == "" {
		return fmt.Errorf("用户不在线")
	}

	onlineUser := model.OnlineUserFromJSON(data)
	remainTTL := time.Duration(onlineUser.ExpireTime-time.Now().UnixMilli()) * time.Millisecond
	if remainTTL <= 0 {
		remainTTL = onlineUserTTL
	}

	cache.GlobalCache.Set(ctx, getBlacklistKey(tokenId), "1", remainTTL)
	cache.GlobalCache.Delete(ctx, key)

	return nil
}

func (s *onlineUserService) List(ctx context.Context) ([]*model.OnlineUser, error) {
	if cache.GlobalCache == nil {
		return []*model.OnlineUser{}, nil
	}

	var result []*model.OnlineUser

	if memCache, ok := cache.GlobalCache.(*cache.MemoryCache); ok {
		memCache.ForEach(func(key, value string) bool {
			if strings.HasPrefix(key, onlineUserPrefix) {
				onlineUser := model.OnlineUserFromJSON(value)
				result = append(result, onlineUser)
			}
			return true
		})
	}

	return result, nil
}

func (s *onlineUserService) extractTokenId(token string) string {
	if len(token) <= 16 {
		return token
	}
	return token[:16]
}
