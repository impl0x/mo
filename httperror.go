package mo

import (
	"net/http"
	"strconv"
)

// Returns a new instance with the code and message.
func NewHttpError(code int, message string) HttpError {
	return HttpError{
		Code:    code,
		Message: message,
	}
}

// Used by the router and also by the user if they want to.
// to be used by value
type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// checks if a error is empty
func (h HttpError) IsNil() bool {
	return h.Code == 0
}

func (h HttpError) StatusCode() int {
	return h.Code
}
func (h HttpError) Error() string {
	return "code=" + strconv.Itoa(h.Code) + " message=" + h.Message
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
