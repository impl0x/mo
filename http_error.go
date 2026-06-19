package mo

import (
	"encoding/json"
	"net/http"
)

type HttpError struct {
	*Response
}

func NewError(c ContentType, body any) *HttpError {
	var b any
	switch c {
	case TEXT:
		b = body
	case JSON:
		b = map[string]any{"errors": body}
	}
	return &HttpError{
		&Response{
			ContentType: c,
			Body:        b,
			StatusCode:  http.StatusBadRequest,
		},
	}
}

func (h *HttpError) Error() string {
	val, err := json.Marshal(h.Body)
	if err != nil {
		return "unknown error"
	}
	return string(val)
}
