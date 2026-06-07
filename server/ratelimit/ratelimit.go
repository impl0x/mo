package ratelimit

import (
	"go-backend/modules/logger"
	"strings"
)

//? configs


type RatelimitConfig struct { // in future i plan to make the global and local be of a type interface of ratelimiter, which will then allow multiple types of ratelimiting to be present
	Global       FixedWindowCounterConfig // Global refers to the entire server application, if max requests in a minute second exceeds it will serve 429 to everyone and every endpoint.
	Local        TokenBucketConfig	// Local refers to per session ratelimiting, based on IP and uses token bucket algorithm
	ErrorMessage string
}

//? ratelimiter
type ratelimit struct{
	global FixedWindowCounter
	local TokenBucket
	errorMsg string
}



//? main functions
var ratelimiter ratelimit

func Init(config RatelimitConfig) {
	// checks for invalid config data
	if config.Global.MaxRequests == 0 {
		logger.Fatal("Max requests cannot be 0 !", "ratelimit error", "global")
	} else if config.Local.MaxTokens == 0 {
		logger.Fatal("Max tokens cannot be 0 !", "ratelimit error", "per ip")
	} else if config.Local.RefillRate == 0 {
		logger.Fatal("Refill rate cannot be 0 !", "ratelimit error", "per ip")
	} else if strings.TrimSpace(config.ErrorMessage) == "" {
		logger.Warn("Error message as receivied in config is empty space! Will fall back to the default message", "ratelimit error", "ErrorMessage value")
	}
	ratelimiter.global.Config=config.Global
	ratelimiter.local.Config=config.Local
	ratelimiter.errorMsg=config.ErrorMessage
}

