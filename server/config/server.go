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

// Uses the default value for some fields
// and takes the integral values from parameters
func RatelimitConfigWithDefaults(globalMaxRequests uint16, localMaxTokens, localRefillRate uint8) *ratelimit.RatelimitConfig {
	return &ratelimit.RatelimitConfig{
		Global: ratelimit.FixedWindowCounterConfig{
			MaxRequests: globalMaxRequests,
		},
		Local: ratelimit.TokenBucketConfig{
			MaxTokens:  localMaxTokens,
			RefillRate: localRefillRate,
			RetryAfter: true,
		},
		ErrorMessage: "Too many requests!",
	}
}
