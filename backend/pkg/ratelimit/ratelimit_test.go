package ratelimit

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	rl := New(3, 100*time.Millisecond)
	defer rl.Stop()

	if !rl.Allow("key1") {
		t.Error("first request should be allowed")
	}
	if !rl.Allow("key1") {
		t.Error("second request should be allowed")
	}
	if !rl.Allow("key1") {
		t.Error("third request should be allowed")
	}
	if rl.Allow("key1") {
		t.Error("fourth request should be blocked")
	}
}

func TestRateLimiter_Expiry(t *testing.T) {
	rl := New(1, 50*time.Millisecond)
	defer rl.Stop()

	if !rl.Allow("key2") {
		t.Error("first request should be allowed")
	}
	if rl.Allow("key2") {
		t.Error("second request should be blocked")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow("key2") {
		t.Error("request after expiry should be allowed")
	}
}

func TestRateLimiter_DifferentKeys(t *testing.T) {
	rl := New(1, time.Minute)
	defer rl.Stop()

	if !rl.Allow("user:a") {
		t.Error("user a first request should be allowed")
	}
	if !rl.Allow("user:b") {
		t.Error("user b first request should be allowed")
	}
	if rl.Allow("user:a") {
		t.Error("user a second request should be blocked")
	}
	if rl.Allow("user:b") {
		t.Error("user b second request should be blocked")
	}
}
