// Package ratelimit 提供基于滑动窗口策略的内存限流器
package ratelimit

import (
	"sync"
	"time"
)

type entry struct {
	count    int       // 当前窗口内请求计数
	expireAt time.Time // 窗口过期时间
}

// RateLimiter 基于滑动窗口策略的内存限流器，带后台 goroutine 清理过期条目
type RateLimiter struct {
	mu       sync.RWMutex
	entries  map[string]*entry
	limit    int
	window   time.Duration
	cleanupInterval time.Duration
	stopCh   chan struct{}
}

func New(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		entries:  make(map[string]*entry),
		limit:    limit,
		window:   window,
		cleanupInterval: window * 2,
		stopCh:   make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	e, ok := rl.entries[key]
	if !ok || now.After(e.expireAt) {
		rl.entries[key] = &entry{count: 1, expireAt: now.Add(rl.window)}
		return true
	}
	if e.count >= rl.limit {
		return false
	}
	e.count++
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			for k, e := range rl.entries {
				if now.After(e.expireAt) {
					delete(rl.entries, k)
				}
			}
			rl.mu.Unlock()
		case <-rl.stopCh:
			return
		}
	}
}

func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}
