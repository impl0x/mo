package mo

import "net/http"

type Response struct {
	http.ResponseWriter
	committed              bool
	defaultHeaders         *HeadersManager
	RequestSpecificHeaders *HeadersManager
}

// here is where we take care of all the headers.
func (r *Response) WriteHeader(statusCode int) {
	headers := r.Header()
	DefaultHeadersConfig.writeHeaders(headers)
	r.defaultHeaders.writeHeaders(headers)
	r.RequestSpecificHeaders.writeHeaders(headers)
	r.ResponseWriter.WriteHeader(statusCode)
	r.committed = true
}

// prettiest code I've ever written probably
