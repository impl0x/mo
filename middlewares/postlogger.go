package middlewares

import (
	"github.com/impl0x/mo"
	"github.com/impl0x/mo/modules/reqlogger"
)

func LoggerWithResponseCode(c *mo.Context) {
	r := c.Request()
	reqlogger.RequestLog(r.RemoteAddr, r.Method, r.URL.Path, c.Response().StatusCode())
}
