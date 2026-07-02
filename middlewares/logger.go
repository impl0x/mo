package middlewares

import (
	"fmt"
	"net/http"

	"github.com/impl0x/mo"
	"github.com/impl0x/mo/modules/logger"
)

func Logger() mo.Middleware {
	return func(next mo.HandlerFunc) mo.HandlerFunc {
		return func(c *mo.Context) error {
			req := c.Request()
			v := fmt.Sprintf(`%v -> "%v"`, req.RemoteAddr, req.URL.Path)
			switch req.Method {
			case http.MethodGet:
				logger.Get(v)
			case http.MethodPost:
				logger.Post(v)
			case http.MethodPatch:
				logger.Patch(v)
			case http.MethodPut:
				logger.Put(v)
			case http.MethodDelete:
				logger.Delete(v)
			default:
				logger.Default(req.Method, v)
			}
			return next(c)
		}
	}
}
