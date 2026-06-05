package ratelimit

import (
	"go-backend/modules/logger"
	"strings"
)

//? configs


type RatelimitConfig struct {
	Global       FixedWindowCounterConfig
	PerIp        TokenBucketConfig
	ErrorMessage string
}

//? ratelimiters




type ratelimit struct{
	global FixedWindowCounter
	perIp TokenBucket
	errorMsg string
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
	ratelimiter.global.Config=config.Global
	ratelimiter.perIp.Config=config.PerIp
	ratelimiter.errorMsg=config.ErrorMessage
}

