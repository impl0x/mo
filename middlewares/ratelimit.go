package middlewares

import (
	"github.com/impl0x/mo"
	"github.com/impl0x/mo/middlewares/ratelimiters"
)

func Ratelimit(r ratelimiters.Ratelimter) mo.Middleware {
	return func(next mo.HandlerFunc) mo.HandlerFunc {
		return func(c *mo.Context) error {
			if !r.Allow(c.Request()) {
				return mo.ErrTooManyRequests
			}
			return next(c)
		}
	}
}
