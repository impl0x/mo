package middlewares

import (
	"fmt"
	"runtime/debug"

	"github.com/impl0x/mo"
	"github.com/impl0x/mo/modules/logger"
)

func Recover(logError bool) mo.Middleware {
	return func(next mo.HandlerFunc) mo.HandlerFunc {
		return func(c *mo.Context) error {
			defer func() {
				if err := recover(); err != nil {
					logger.Panic(fmt.Sprintf("%v", err))
					println("[STACK TRACE]\n" + string(debug.Stack()))
					c.JSON(500, mo.ErrInternalServerError.JsonFormat())
				}
			}()
			return next(c)
		}
	}
}
