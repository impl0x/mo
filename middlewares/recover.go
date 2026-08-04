package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/impl0x/mo"
	"github.com/impl0x/mo/modules/logger"
)

// set logError to true to log the panics in your code. it is recommended to do so.
func Recover(logError bool) mo.Middleware {
	return func(next mo.HandlerFunc) mo.HandlerFunc {
		return func(c *mo.Context) error {
			defer func() {
				if err := recover(); err != nil {
					if logError {
						logger.Panic(fmt.Sprintf("%v", err))
						println("[STACK TRACE]\n" + string(debug.Stack()))
					}
					c.JSON(http.StatusInternalServerError, mo.ErrInternalServerError)
				}
			}()
			return next(c)
		}
	}
}
