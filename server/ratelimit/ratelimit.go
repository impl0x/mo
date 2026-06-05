package ratelimit

import (
	"go-backend/modules/logger"
	"strings"
	"time"
)
//? configs
type FixedWindowCounterConfig struct { // Global rl
	MaxRequests uint16 // per sec
}
type TokenBucketConfig struct { // PerIp rl
	MaxTokens  uint8
	RefillRate uint8 // per sec
	RetryAfter bool  // sends the "retry_after" in json response, defaults to true
}
type RatelimitConfig struct {
	Global       FixedWindowCounterConfig
	PerIp        TokenBucketConfig
	ErrorMessage string
}

//? ratelimiters
type fixedWindowCounter struct {
	config  FixedWindowCounterConfig
	counter uint16
}

type tokenBucket struct {
	config TokenBucketConfig
	ips    []userIp
}

type ratelimit struct{
	global fixedWindowCounter
	perIp tokenBucket
	errorMsg string
}

type userIp struct {
	ip       string
	tokens   uint8
	lastSeen time.Time
}

//? main functions
var ratelimiter ratelimit

func Init(config RatelimitConfig) {
	// checks for invalid config data
	if config.Global.MaxRequests == 0 {
		logger.Fatal("Max requests cannot be 0 !", "ratelimit error", "global")
	} else if config.PerIp.MaxTokens == 0 {
		logger.Fatal("Max tokens cannot be 0 !", "ratelimit error", "per ip")
	} else if config.PerIp.RefillRate == 0 {
		logger.Fatal("Refill rate cannot be 0 !", "ratelimit error", "per ip")
	} else if strings.TrimSpace(config.ErrorMessage) == "" {
		logger.Warn("Error message as receivied in config is empty space! Will fall back to the default message", "ratelimit error", "ErrorMessage value")
	}
	ratelimiter.global.config=config.Global
	ratelimiter.perIp.config=config.PerIp
	ratelimiter.errorMsg=config.ErrorMessage
}

