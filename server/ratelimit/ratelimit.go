package ratelimit

import (
	// "go-backend/server/core/servertypes"
	"net/http"
)

type Ratelimiter interface{
	Allow(next http.Handler)http.Handler
	
}


//? configs
type RatelimitConfig [G any, L any] struct {
	Global       G // Global refers to the entire server application, if max requests in a minute second exceeds it will serve 429 to everyone and every endpoint.
	Local        L	// Local refers to per session ratelimiting, based on IP and uses token bucket algorithm
	ErrorMessage string
}



//? main ratelimiter
type ratelimit struct{
	global ratelimiters.FixedWindowCounter
	local ratelimiters.TokenBucket
	errorMsg string
}


//? main functions
var ratelimiter ratelimit

// func Init  [G any, L any] (config RatelimitConfig[G, L]) {
	// checks for invalid config data
	// if config.Global.MaxRequests == 0 {
	// 	logger.Fatal("Max requests cannot be 0 !", "ratelimit error", "global")
	// } else if config.Local.MaxTokens == 0 {
	// 	logger.Fatal("Max tokens cannot be 0 !", "ratelimit error", "per ip")
	// } else if config.Local.RefillRate == 0 {
	// 	logger.Fatal("Refill rate cannot be 0 !", "ratelimit error", "per ip")
	// } else if strings.TrimSpace(config.ErrorMessage) == "" {
	// 	logger.Warn("Error message as receivied in config is empty space! Will fall back to the default message", "ratelimit error", "ErrorMessage value")
	// }
	// ratelimiter.global.Config=config.Global
	// ratelimiter.local.Config=config.Local
	// ratelimiter.errorMsg=config.ErrorMessage
// }

func Init[T Ratelimiter](){
	
}