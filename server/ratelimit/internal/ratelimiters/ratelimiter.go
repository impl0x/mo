package ratelimiters

import "net/http"

type Ratelimiter interface {
	Allow(r *http.Request) bool
}

type RatelimiterConfig interface {
	Default() *Ratelimiter
}

func New[T Ratelimiter]() {
	
}
