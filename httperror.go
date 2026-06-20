package mo

import (
	"fmt"
	"net/http"
)

type HttpError interface{
	StatusCode()int
	StatusText()string
	error
}

// used to store the common errors
type httpError struct {
	code int
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
	return he.code
}
func (he httpError) StatusText()string{
	return http.StatusText(he.code) // does not include status code
}
func (he httpError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.code, he.StatusText())
}
