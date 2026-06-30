package mo

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/impl0x/mo/modules/logger"
)





type Context struct {
	request  *http.Request
	response *Response
	Headers  map[string]string // stores the headers for specific requests
	Mo       *Mo               // original Mo instance
	Store    map[string]any    // stores context values
}

func (c *Context) writeContentType(value string) {
	header := c.response.Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

func (c *Context) Request() *http.Request {
	return c.request
}
func (c *Context) Response() *Response {
	return c.response
}

// Redirect redirects the request to a provided URL with status code.
func (c *Context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return ErrInternalServerError
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}

// NoContent sends a response with no body and a status code.
func (c *Context) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

// Blob sends a blob response with status code and content type.
func (c *Context) Blob(code int, contentType string, b []byte) error {
	c.writeContentType(contentType)
	c.response.WriteHeader(code)
	_, err := c.response.Write(b)
	if err != nil {
		if c.Mo.Config.LogErrors {
			logger.Mo("Client disconnected! couldn't write response")
		}
	}
	return nil
}

// JSON sends a JSON response with status code.
func (c *Context) JSON(code int, target any) error {
	c.writeContentType(MIMEApplicationJSON)
	c.response.WriteHeader(code)
	return json.NewEncoder(c.response).Encode(target)
}

func (c *Context) TEXT(code int, body string) error {
	return c.Blob(code, MIMETextPlain, []byte(body))
}

// ErrNonExistentKey is error that is returned when key does not exist
var ErrNonExistentKey = errors.New("non existent key")

// ErrInvalidKeyType is error that is returned when the value is not castable to expected type.
var ErrInvalidKeyType = errors.New("invalid key type")

// Adds a value to the context storage
func (c *Context) Add(key string, value any) {
	c.Store[key] = value
}

// Gets a value from the context storage
func (c *Context) Get(key string) (any, bool) {
	v, ok := c.Store[key]
	return v, ok
}

// Gets a value from the context storage (typed)
func ContextGet[T any](c *Context, key string) (T, error) {
	value, ok := c.Store[key]
	if !ok {
		var zero T
		return zero, ErrNonExistentKey
	}
	typed, ok := value.(T)
	if !ok {
		var zero T
		return zero, ErrInvalidKeyType
	}
	return typed, nil
}

// Binds the request headers to a struct
//
// must contain tag `header`
//
// example: 
// 	token string `header:"authorization"`
//
// fields of the struct MUST be strings!
func (c *Context) BindHeaders(target any) {

	headers := c.request.Header

	rv := reflect.ValueOf(target)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	rt := rv.Type()
	if rv.Kind() != reflect.Struct {
		if c.Mo.Config.LogErrors {
			logger.Mo("Cannot bind headers to a non struct object")
		}
		return
	}
	for i := range rv.NumField() {
		v := rv.Field(i)
		t := rt.Field(i)
		if !t.IsExported() {
			continue
		}
		if v.Kind() != reflect.String {
			if c.Mo.Config.LogErrors {
				logger.Mo("Header binding variables must be strictly string")
			}
			continue // headers values must be strings strictly
		}
		tag, ok := t.Tag.Lookup("header")
		if !ok {
			continue
		}
		value := headers.Get(tag)
		if value == "" {
			continue
		}
		v.SetString(value)
	}
}
