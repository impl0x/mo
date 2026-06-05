package config

import (
	"go-backend/server/ratelimit"
)

type ServerConfig struct {
	Host string
	Port uint16

	LogRequests bool
	Ratelimit   ratelimit.RatelimitConfig
}

var DefaultRatelimitConfig = ratelimit.RatelimitConfig{
	Global: ratelimit.FixedWindowCounterConfig{
		MaxRequests: 100,
	},
	PerIp: ratelimit.TokenBucketConfig{
		MaxTokens:      5,
		RefillRate:     2,
		RetryAfter: true,
	},
	ErrorMessage: "Too many requests!",
}
