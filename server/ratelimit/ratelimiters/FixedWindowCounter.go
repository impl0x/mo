package ratelimiters

import (
	"net/http"
)


// ? Configs
type FixedWindowCounterConfig struct {
	MaxRequests uint16 // per sec
}

func NewFixedWindowCounter(config FixedWindowCounterConfig)Ratelimiter{
	return &fixedWindowCounter{
		config: config,
	}
}


// ? Internal state
type fixedWindowCounter struct {
	config  FixedWindowCounterConfig
	counter uint16
}

func (rl *fixedWindowCounter) Allow(r *http.Request) bool {
	// Todo: finish logic
	return true
}

