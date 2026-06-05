package ratelimit

import (
	"time"
)


type FixedWindowCounterConfig struct {
	MaxRequests    uint16 // per sec
}
type TokenBucketConfig struct {
	MaxTokens  uint8
	RefillRate uint8 // per sec
	SendRetryAfter bool   // sends the "retry_after" in json response if set to true
}
type RatelimitConfig struct {
	Global FixedWindowCounterConfig
	PerIp  TokenBucketConfig
	ErrorMessage string
}
type fixedWindowCounter struct {
	config  FixedWindowCounterConfig
	counter uint16
}


type tokenBucket struct {
	config TokenBucketConfig
	ips    *[]userIp
}

type userIp struct {
	ip       string
	tokens   uint8
	lastSeen time.Time
}


func initRatelimiter(config RatelimitConfig){
	// checks for invalid data
	// global
	if config.Global.MaxRequests==0{
		print("Invalid data")
		// os.Exit(1)
		// TODO: implement a good logger here which does fatal
	}
}