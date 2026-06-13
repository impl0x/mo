package ratelimiters

import (
	"net/http"
	"time"
)

// ? Configs
type TokenBucketConfig struct {
	MaxTokens  uint8
	RefillRate uint8 // per sec
	RetryAfter bool  // sends the "retry_after" in json response, defaults to true
}
func NewTokenBucket(config TokenBucketConfig)Ratelimiter{
	return &tokenBucket{
		config: config,
	}
}


// ? Internal state
type tokenBucket struct {
	config TokenBucketConfig
	ips    []UserIp
}

type UserIp struct {
	ip       string
	tokens   uint8
	lastSeen time.Time
}

func (rl *tokenBucket) Allow(r *http.Request) bool {
	// Todo: finish logic
	return true
}

