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

type internalErrorJson struct {
	HTTPError
	Error string `json:"error,omitempty"`
}
type validationErrorJson struct {
	HttpError
	Errors []validator.ValidationErrorJson `json:"errors,omitempty"`
}

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
				logger.Mo("cannot write error, response already sent!", "err", err.Error())
			}
			return
		}
		switch e := err.(type) {
		case HTTPError:
			c.JSON(e.StatusCode(), e)
		case *validator.GroupedValidationError:
			c.JSON(http.StatusBadRequest, validationErrorJson{HttpError: ErrBadRequest, Errors: e.ToJsonStructList()})
		case *json.SyntaxError:
			c.JSON(http.StatusUnprocessableEntity, HttpError{
				Code:    http.StatusUnprocessableEntity,
				Message: fmt.Sprintf("JSON syntax error at offset %d", e.Offset),
			})
		case *json.UnmarshalTypeError:
			c.JSON(http.StatusExpectationFailed, HttpError{
				Code:    http.StatusExpectationFailed,
				Message: fmt.Sprintf("Wrong type used for field %s", e.Field),
			})
		case nil:
			c.NoContent(http.StatusNoContent)
		default:
			if e.Error() == "EOF" { // rare case because json parsing returns a errorString of EOF when a body is empty.
				c.JSON(http.StatusUnprocessableEntity, HttpError{
					Code:    http.StatusUnprocessableEntity,
					Message: "EOF",
				})
				return
			}
			resp := internalErrorJson{HTTPError: ErrInternalServerError}
			if exposeError {
				resp.Error = e.Error()
				logger.Mo("Internal error: "+e.Error(), "errorType", fmt.Sprintf("%T", e))
			}
			c.JSON(http.StatusInternalServerError, resp)
		}
	}
}
