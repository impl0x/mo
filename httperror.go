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
	return &HttpError{
		Code:    code,
		Message: message,
	}
}

// error occurred during request lifecycle
type HttpError struct {
	Code    int
	Message string
}

func (h *HttpError) StatusCode() int {
	return h.Code
}
func (h *HttpError) JsonFormat() any {
	return map[string]any{
		"code":    h.Code,
		"message": h.Message,
	}
}
func (h *HttpError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", h.Code, h.Message)
}

// used to store the common errors
type httpError struct {
	Code int
}

var (
	ErrBadRequest                  = &httpError{http.StatusBadRequest}            // 400
	ErrUnauthorized                = &httpError{http.StatusUnauthorized}          // 401
	ErrForbidden                   = &httpError{http.StatusForbidden}             // 403
	ErrNotFound                    = &httpError{http.StatusNotFound}              // 404
	ErrMethodNotAllowed            = &httpError{http.StatusMethodNotAllowed}      // 405
	ErrRequestTimeout              = &httpError{http.StatusRequestTimeout}        // 408
	ErrStatusRequestEntityTooLarge = &httpError{http.StatusRequestEntityTooLarge} // 413
	ErrUnsupportedMediaType        = &httpError{http.StatusUnsupportedMediaType}  // 415
	ErrTooManyRequests             = &httpError{http.StatusTooManyRequests}       // 429
	ErrInternalServerError         = &httpError{http.StatusInternalServerError}   // 500
	ErrBadGateway                  = &httpError{http.StatusBadGateway}            // 502
	ErrServiceUnavailable          = &httpError{http.StatusServiceUnavailable}    // 503
)

func (he httpError) StatusCode() int {
	return he.Code
}
func (he httpError) StatusText() string {
	return http.StatusText(he.Code) // does not include status code
}
func (he httpError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.StatusText())
}
func (he httpError) JsonFormat() any {
	return map[string]any{
		"code":he.StatusCode(),
		"message":he.StatusText(),
	}
}
