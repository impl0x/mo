package mo

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
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
		var jsonSyntaxErr *jsontext.SyntacticError
		var jsonSemanticErr *json.SemanticError
		var httpErr HTTPError
		var vdErr validator.GroupedValidationError
		if c.response.committed {
			if err == nil {
				return
			}
			if c.Mo.Config.LogErrors {
				logger.Mo("cannot write error, response already sent!", "err", err.Error())
			}
			return
		}
		switch {
		case errors.As(err, &httpErr):
			c.JSON(httpErr.StatusCode(), httpErr)
		case errors.As(err, &vdErr):
			c.JSON(http.StatusBadRequest, validationErrorJson{HttpError: ErrBadRequest, Errors: vdErr.ToJsonStructList()})
		case errors.As(err, &jsonSyntaxErr):
			c.JSON(http.StatusUnprocessableEntity, HttpError{
				Code:    http.StatusUnprocessableEntity,
				Message: fmt.Sprintf("JSON syntax error at offset %d", jsonSyntaxErr.ByteOffset),
			})
		case errors.As(err, &jsonSemanticErr):
			c.JSON(http.StatusBadRequest, HttpError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Wrong type used for field %s, expected type: %s", jsonSemanticErr.JSONPointer.LastToken(), jsonSemanticErr.GoType.String()),
			})
		case err == nil:
			c.NoContent(http.StatusNoContent)
		default:
			if err.Error() == "EOF" {
				c.JSON(http.StatusUnprocessableEntity, HttpError{
					Code:    http.StatusUnprocessableEntity,
					Message: "Syntax Error: JSON cut off unexpectedly",
				})
				return
			}
			resp := internalErrorJson{HTTPError: ErrInternalServerError}
			if exposeError {
				resp.Error = err.Error()
				logger.Mo("Internal error: "+err.Error(), "errorType", fmt.Sprintf("%T", err))
			}
			c.JSON(http.StatusInternalServerError, resp)
		}
	}
}
