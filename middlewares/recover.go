package middlewares

import (
	"github.com/impl0x/mo"
)

func Recover() mo.Middleware {
	return func(next mo.HandlerFunc) mo.HandlerFunc {
		return func(c *mo.Context) error {
			defer func() {
				if err := Recover(); err != nil {
					c.JSON(500, mo.ErrInternalServerError.JsonFormat())
				}
			}()
			return next(c)
		}
	}
}
