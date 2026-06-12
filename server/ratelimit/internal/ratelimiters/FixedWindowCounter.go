package ratelimiters

import (
	"net/http"
)

type FixedWindowCounterConfig struct {
	MaxRequests uint16 // per sec
}

type FixedWindowCounter struct {
	Config  FixedWindowCounterConfig
	counter uint16
}

func (rl *FixedWindowCounter) Allow(next http.Handler) http.Handler {
	return next
}
