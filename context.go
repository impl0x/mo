package mo

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/impl0x/mo/modules/logger"
)

type Response struct {
	http.ResponseWriter
	committed bool
}

func (r *Response) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.committed = true
}

type Context struct {
	request  *http.Request
	response *Response
	Mo       *Mo            // original Mo instance
	Store    map[string]any // stores context values
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
func (c *Context) Response() http.ResponseWriter {
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
	writeResp(c.response, b)
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

func writeResp(resp http.ResponseWriter, b []byte) {
	_, err := resp.Write(b)
	if err != nil {
		logger.Mo("Client disconnected! couldn't write response")
		return
	}
}

// ErrNonExistentKey is error that is returned when key does not exist
var ErrNonExistentKey = errors.New("non existent key")

// ErrInvalidKeyType is error that is returned when the value is not castable to expected type.
var ErrInvalidKeyType = errors.New("invalid key type")

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
