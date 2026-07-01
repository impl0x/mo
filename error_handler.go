package mo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/impl0x/mo/modules/logger"
	"github.com/impl0x/mo/validator"
)

// Error Handler must handle nil, HttpErrorInterface and error. (internal)
type HTTPErrorHandler func(*Context, error)

// if err==nil, returns
// if response already committed (headers written), returns and logs
// if error is of type HttpErrorInterface, calls c.Json() with e.StatusCode() and e.JsonFormat()
// if you want a custom error message returned, implement the HTTPErrorInterface.
// Then return a valid json from JsonFormat() method and a valid status-code from StatusCode()
func DefaultHTTPErrorHandler(exposeError bool) HTTPErrorHandler {
	return func(c *Context, err error) {
		if err == nil {
			return
		}
		if c.response.committed {
			if c.Mo.Config.LogErrors {
				logger.Mo("Cannot write error, response already sent!", "err", err.Error())
			}
			return
		}
		switch e := err.(type) {
		case HttpErrorInterface:
			c.JSON(e.StatusCode(), e.JsonFormat())
		case *validator.GroupedValidationError:
			c.JSON(http.StatusBadRequest, map[string]any{
				"code":    http.StatusBadRequest,
				"message": "Validation error",
				"errors":  e.JsonFormat(),
			})
		case *json.SyntaxError:
			c.JSON(http.StatusUnprocessableEntity, map[string]any{
				"code":    http.StatusUnprocessableEntity,
				"message": fmt.Sprintf("JSON syntax error at offset %d", e.Offset),
			})
		case *json.UnmarshalTypeError:
			c.JSON(http.StatusExpectationFailed, map[string]any{
				"code":    http.StatusExpectationFailed,
				"message": fmt.Sprintf("Wrong type used for field %s", e.Field),
			})
		default:
			resp := map[string]any{
				"code":    http.StatusInternalServerError,
				"message": http.StatusText(http.StatusInternalServerError),
			}
			if exposeError {
				resp["error"] = e.Error()
			}
			c.JSON(http.StatusInternalServerError, resp)
		}
	}
}
