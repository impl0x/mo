package mo

import "net/http"

type Response struct {
	ContentType ContentType
	Body        any
	StatusCode  int
}

func DefaultResponse() *Response {
	return &Response{
		StatusCode:  http.StatusOK,
	}
}
