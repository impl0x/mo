package mo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/impl0x/mo/errs"
)

type JSONSerializer struct{}

func (j JSONSerializer) Serialize(c *Context, target any) error {
	return json.NewEncoder(c.GetResponse()).Encode(target)
}
func (j JSONSerializer) Deserialize(c *Context, target any) error {
	err := json.NewDecoder(c.GetRequest().Body).Decode(target)
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError
	switch {
	// serialization errors
	case errors.Is(err, io.EOF):
		return NewError(JSON, errs.Serialization{Message:"Request body cannot be empty"})
	case errors.Is(err, io.ErrUnexpectedEOF):
		return NewError(JSON, errs.Serialization{Message:"Malformed JSON: Request body was cut short"})
	case errors.Is(err, syntaxError):
		msg := fmt.Sprintf("Malformed JSON: Syntax error at byte position %d", syntaxError.Offset)
		return NewError(JSON, errs.Serialization{Message: msg})
	// validation errors
	case errors.As(err, &unmarshalTypeError):
		return NewError(
			JSON,
			[...]errs.Validation{
				errs.Validation{
					Field:   unmarshalTypeError.Field,
					Message: fmt.Sprintf("Expected %s, Received %s", unmarshalTypeError.Type, unmarshalTypeError.Value),
				},
			},
		)

	case errors.As(err, &invalidUnmarshalError):
		// This is a 500 error because the developer broke the code backend logic
		return &HttpError{&Response{TEXT, "Internal server error", http.StatusInternalServerError}}
	default:
		return NewError(TEXT, "Failed to parse request payload")
	}
}

var DefaultJsonSerializer = JSONSerializer{}

var JSON = ContentType{
	value:     "application/json",
	formatter: DefaultJsonSerializer,
}
