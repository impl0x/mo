package ratelimit

import (
	"go-backend/server/core/servertypes"
	"go-backend/server/ratelimit/ratelimiters"
	"net/http"
)

type Ratelimit struct {
	StatusCode   int    // default to 429
	ErrorMessage string // default to "Too many requests!"
}

// returns a middleware which implements the ratelimiter
func (r *Ratelimit) NewRatelimiter(rl ratelimiters.Ratelimiter) servertypes.Middleware {
	if r.ErrorMessage == "" {
		r.ErrorMessage = "Too many requests!"
	} else if r.StatusCode == 0 {
		if r.StatusCode > 599 {
			panic("Status code cannot be greater than 599")
		}
		r.StatusCode = 429
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !rl.Allow(r){
				// TODO: return a 429 and ErrorMessage in proper format
			}
		})
	}
}
