package mo

import (
	"errors"
	"net/http"
)

var ErrResponseAlreadyCommitted = errors.New("headers already written")

type Response struct {
	http.ResponseWriter
	committed              bool
	statusCode             int
	defaultHeaders         *HeadersManager // pointer to mo.Headers, is present only to write in the [Write] func
	RequestSpecificHeaders HeadersManager
}

// Returns a new Response struct with the default status code 200 OK
func newResponse(w http.ResponseWriter, defaultHeaders *HeadersManager) Response {
	return Response{
		ResponseWriter:         w,
		statusCode:             http.StatusOK,
		defaultHeaders:         defaultHeaders,
		RequestSpecificHeaders: DefaultHeadersManager(),
	}
}

// we cache the statusCode and then send it on the first write call after writing all the headers.
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
	r.committed = true // we set committed to true to mark that the response has been written.
	return r.ResponseWriter.Write(b)
}

func (r *Response) StatusCode() int {
	return r.statusCode
}
