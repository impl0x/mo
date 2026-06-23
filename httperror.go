package mo

import (
	"fmt"
	"net/http"
)

type HttpErrorInterface interface {
	StatusCode() int
	JsonFormat() any
	error
}
// To return a custom formatted message, return a struct implementing HttpErrorInterface
// Or just return c.Json with a statusCode
// Or just define a custom function for yourself, anything works.
func NewHTTPError(code int, message string) HttpErrorInterface {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

// error occurred during request lifecycle
// Used manually
type HTTPError struct {
	Code    int
	Message string
}

func (h *HTTPError) StatusCode() int {
	return h.Code
}
func (h *HTTPError) JsonFormat() any {
	return map[string]any{
		"code":    h.Code,
		"message": h.Message,
	}
}
func (h *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", h.Code, h.Message)
}

// used to store the common errors
// not to be used by the user directly
type HttpError struct {
	Code int
}

var (
	ErrBadRequest                  = &HttpError{http.StatusBadRequest}            // 400
	ErrUnauthorized                = &HttpError{http.StatusUnauthorized}          // 401
	ErrForbidden                   = &HttpError{http.StatusForbidden}             // 403
	ErrNotFound                    = &HttpError{http.StatusNotFound}              // 404
	ErrMethodNotAllowed            = &HttpError{http.StatusMethodNotAllowed}      // 405
	ErrRequestTimeout              = &HttpError{http.StatusRequestTimeout}        // 408
	ErrStatusRequestEntityTooLarge = &HttpError{http.StatusRequestEntityTooLarge} // 413
	ErrUnsupportedMediaType        = &HttpError{http.StatusUnsupportedMediaType}  // 415
	ErrTooManyRequests             = &HttpError{http.StatusTooManyRequests}       // 429
	ErrInternalServerError         = &HttpError{http.StatusInternalServerError}   // 500
	ErrBadGateway                  = &HttpError{http.StatusBadGateway}            // 502
	ErrServiceUnavailable          = &HttpError{http.StatusServiceUnavailable}    // 503
)

func (he HttpError) StatusCode() int {
	return he.Code
}
func (he HttpError) StatusText() string {
	return http.StatusText(he.Code) // does not include status code
}
func (he HttpError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.StatusText())
}
func (he HttpError) JsonFormat() any {
	return map[string]any{
		"code":he.StatusCode(),
		"message":he.StatusText(),
	}
}
