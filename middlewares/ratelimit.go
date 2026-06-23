package middlewares

import (
	"net/http"
	"github.com/impl0x/mo"
)

type Ratelimter interface {
	Allow(r *http.Request) bool
}

func Ratelimit(r Ratelimter) mo.Middleware {
	return func(next mo.HandlerFunc) mo.HandlerFunc {
		return func(c *mo.Context) error {
			if !r.Allow(c.Request()) {
				return mo.ErrTooManyRequests
			}
			return next(c)
		}
	}
}
