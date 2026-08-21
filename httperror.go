package mo

import (
	"fmt"
	"net/http"
)

// If you are implementing this interface and using the default http error handler then
// make sure that the struct is json compatible because its going to be sent directly to the json encoder
type HTTPError interface {
	StatusCode() int
	error
}

// To return a custom formatted message, return a struct implementing HTTPError
// Or just return c.Json with a statusCode
// Or just define a custom function for yourself, anything works.
func NewHttpError(code int, message string) HttpError {
	return HttpError{
		Code:    code,
		Message: message,
	}
}

// error occurred during request lifecycle
type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h HttpError) StatusCode() int {
	return h.Code
}
func (h HttpError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", h.Code, h.Message)
}

// common http errors with the default status code text
//
// These can be directly passed to context.JSON as these structs are json compatible, feel free to return them directly too if using the default error handler
var (
	ErrBadRequest                  = HttpError{http.StatusBadRequest, http.StatusText(http.StatusBadRequest)}                       // 400
	ErrUnauthorized                = HttpError{http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized)}                   // 401
	ErrForbidden                   = HttpError{http.StatusForbidden, http.StatusText(http.StatusForbidden)}                         // 403
	ErrNotFound                    = HttpError{http.StatusNotFound, http.StatusText(http.StatusNotFound)}                           // 404
	ErrMethodNotAllowed            = HttpError{http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed)}           // 405
	ErrRequestTimeout              = HttpError{http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout)}               // 408
	ErrStatusRequestEntityTooLarge = HttpError{http.StatusRequestEntityTooLarge, http.StatusText(http.StatusRequestEntityTooLarge)} // 413
	ErrUnsupportedMediaType        = HttpError{http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType)}   // 415
	ErrTooManyRequests             = HttpError{http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests)}             // 429
	ErrInternalServerError         = HttpError{http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)}     // 500
	ErrBadGateway                  = HttpError{http.StatusBadGateway, http.StatusText(http.StatusBadGateway)}                       // 502
	ErrServiceUnavailable          = HttpError{http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable)}       // 503
)
