package mo

import (
	"errors"
	"net/http"
)

var ErrResponseAlreadyCommitted = errors.New("Headers already written")

type Response struct {
	http.ResponseWriter
	committed              bool
	statusCode             int
	defaultHeaders         *HeadersManager
	RequestSpecificHeaders *HeadersManager
}

func newResponse(w http.ResponseWriter, defaultHeaders *HeadersManager) *Response {
	return &Response{
		w,
		false,
		200,
		defaultHeaders,
		DefaultHeadersManager(),
	}
}

// we cache the statusCode and then send it on the first write call
func (r *Response) WriteHeader(statusCode int) {

	r.statusCode = statusCode
}

// here is where we take care of all the headers.
func (r *Response) Write(b []byte) (int, error) {
	if r.committed {
		return 0, ErrResponseAlreadyCommitted
	}
	headers := r.Header()
	DefaultHeadersConfig.writeHeaders(headers)
	r.defaultHeaders.writeHeaders(headers)
	r.RequestSpecificHeaders.writeHeaders(headers)

	r.ResponseWriter.WriteHeader(r.statusCode)
	r.committed = true
	return r.ResponseWriter.Write(b)
}
