// Package cache 提供缓存抽象接口和内存实现，支持 Redis 切换
package cache

import (
	"context"
	"backend/pkg/global_vars"
	"sync"
	"time"
)

// Cache 缓存抽象接口，支持内存和 Redis 两种实现切换
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

// MemoryCache 基于 sync.RWMutex + map 的进程内缓存实现
type MemoryCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value      string
	expireTime time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]cacheItem),
	}
}

func (m *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.data[key]
	if !ok {
		return "", nil
	}
	if time.Now().After(item.expireTime) {
		return "", nil
	}
	return item.value, nil
}

func (m *MemoryCache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = cacheItem{
		value:      value,
		expireTime: time.Now().Add(expiration),
	}
	return nil
}

func (m *MemoryCache) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MemoryCache) ForEach(fn func(key, value string) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	now := time.Now()
	for key, item := range m.data {
		if now.After(item.expireTime) {
			continue
		}
		if !fn(key, item.value) {
			break
		}
	}
}

var GlobalCache Cache

func InitCache() {
	useRedis := global_vars.ConfigYml.GetBool("Cache.UseRedis")
	if useRedis {
		GlobalCache = nil
		global_vars.ZapLog.Warn("Redis 缓存尚未实现，请实现 RedisCache 或设置 Cache.UseRedis=false")
	} else {
		GlobalCache = NewMemoryCache()
		global_vars.ZapLog.Info("初始化内存缓存成功")
	}
}
