package ratelimiters

import "net/http"

type Ratelimiter interface {
	Allow(r *http.Request) bool
}