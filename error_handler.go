package mo

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/impl0x/mo/modules/logger"
	"github.com/impl0x/mo/validator/v3"
)

// Error Handler must handle nil, HttpErrorInterface and error. (internal)
type HTTPErrorHandler func(*Context, error)

type internalErrorJson struct {
	HttpError
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
		if err == nil {
			c.NoContent(http.StatusNoContent)
		} else if e, ok := errors.AsType[HttpError](err); ok {
			c.JSON(e.StatusCode(), e)
		} else if e, ok := errors.AsType[validator.GroupedValidationError](err); ok {
			c.JSON(http.StatusBadRequest, validationErrorJson{HttpError: ErrBadRequest, Errors: e.ToJsonStructList()})
		} else if e, ok := errors.AsType[*jsontext.SyntacticError](err); ok {
			c.JSON(http.StatusUnprocessableEntity, HttpError{
				Code:    http.StatusUnprocessableEntity,
				Message: "JSON syntax error at offset " + strconv.FormatInt(e.ByteOffset, 10),
			})
		} else if e, ok := errors.AsType[*json.SemanticError](err); ok {
			c.JSON(http.StatusBadRequest, HttpError{
				Code:    http.StatusBadRequest,
				Message: "Wrong type used for field " + e.JSONPointer.LastToken() + " expected type: " + e.GoType.String(),
			})
		} else {
			if err.Error() == "EOF" {
				c.JSON(http.StatusUnprocessableEntity, HttpError{
					Code:    http.StatusUnprocessableEntity,
					Message: "Syntax Error: JSON cut off unexpectedly",
				})
				return
			}
			resp := internalErrorJson{HttpError: ErrInternalServerError}
			if exposeError {
				resp.Error = err.Error()
				logger.Mo("Internal error: "+err.Error(), "errorType", fmt.Sprintf("%T", err))
			}
			c.JSON(http.StatusInternalServerError, resp)
		}
	}
}
