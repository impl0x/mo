package ratelimiters

import (
	"net/http"
	"time"
)

type TokenBucketConfig struct {
	MaxTokens  uint8
	RefillRate uint8 // per sec
	RetryAfter bool  // sends the "retry_after" in json response, defaults to true
}

type TokenBucket struct {
	Config TokenBucketConfig
	ips    []UserIp
}

type UserIp struct {
	ip       string
	tokens   uint8
	lastSeen time.Time
}

func (rl *TokenBucket) Allow(next http.Handler) http.Handler {
	return next
}