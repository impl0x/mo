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
		if c.response.committed {
			if err == nil {
				return
			}
			if c.Mo.Config.LogErrors {
				logger.Mo("Cannot write error, response already sent!", "err", err.Error())
			}
			return
		}
		switch e := err.(type) {
		case HttpErrorInterface:
			err:=c.JSON(e.StatusCode(), e.JsonFormat())	// can throw error if user returns custom HTTP error where the jsonFormat function does not return a valid json.
			if err!=nil{
				c.JSON(ErrInternalServerError.Code, ErrInternalServerError.JsonFormat())
				logger.Mo("Invalid HTTP error returned from handler!", "err", err.Error())
			}
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
		case nil:
			c.NoContent(http.StatusNoContent)
		default:
			if e.Error() == "EOF" { // rare case because json parsing returns a errorString of EOF when a body is empty.
				c.JSON(http.StatusUnprocessableEntity, map[string]any{
					"code":    http.StatusUnprocessableEntity,
					"message": "EOF",
				})
				return
			}
			resp := ErrInternalServerError.JsonFormat().(map[string]any) // safe type conversion because our HttpError method always returns a map[string]any
			if exposeError {
				resp["error"] = e.Error()
				logger.Mo("Internal error: "+e.Error(), "errorType", fmt.Sprintf("%T", e))
			}
			c.JSON(http.StatusInternalServerError, resp)
		}
	}
}
