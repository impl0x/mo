package mo

import (
	"net/http"

)

type HttpError struct {
	StatusCode   int
	Response *Response
}

func NewError() HttpError {
	return HttpError{
		http.StatusBadRequest,
		&Response{
			ContentType: JSON,
			Body: map[string]any{
				"error":"bad request",
			},
			StatusCode: http.StatusBadRequest,
		},
	}
}

func (h *HttpError) Write(c *Context){
	
	
}