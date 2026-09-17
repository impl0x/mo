package mo

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/impl0x/go-utils/cache"
	"github.com/impl0x/mo/modules/logger"
	"github.com/impl0x/mo/validator/v3"
)

var contextPool = sync.Pool{
	New: func() any {
		return &Context{
			store: ContextStore{
				Store:  make(map[string]any),
				Params: make(map[string]string),
			},
		}
	},
}

// stores context values, usually used inside an context instance. It is not goroutine safe, do not pass same mo context to multiple goroutines. otherwise your pc might blow up
type ContextStore struct {
	Store  map[string]any
	Params map[string]string
}

func (cs *ContextStore) clear() {
	clear(cs.Params)
	clear(cs.Params)
}

type Context struct {
	Mo              *Mo // original Mo instance
	request         *http.Request
	response        Response
	ResponseHeaders HeadersManager // Sends headers with the response for this request
	store           ContextStore
}

func (c *Context) writeContentType(value string) {
	c.response.Header().Set(HeaderContentType, value)
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) Response() *Response {
	return &c.response
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
	c.response.ResponseWriter.WriteHeader(code) // skips the delayed response writer cache, because if we don't call write ourselves then http defaults to writing a 200 ok
	c.response.committed = true
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
	return json.MarshalWrite(&c.response, target)
}

func (c *Context) TEXT(code int, body string) error {
	return c.Blob(code, MIMETextPlain, []byte(body))
}

// use get to retrieve values
func (c *Context) QueryParams() url.Values {
	return c.request.URL.Query()
}

// Returns the url parameter, for example if "/users/:id" was registered as a path
// and a request arrives with the path "/users/123", c.Param("id") will return "123" in this case
//
// this also works for wildcard paths, "users/*", "users/123/comments". c.Param("*")="123/comments"
//   - parameter path key: registered key when adding the path, ":id", ":userid", etc becomes "id","userid"
//   - wildcard path key: "*", returns the whole path received after the last static/param path.
//
// note: not goroutine safe, do not pass same mo context to multiple goroutines, if doing so manage your own lock.
func (c *Context) Param(key string) (string, bool) {
	v, ok := c.store.Params[key]
	return v, ok
}

// ErrNonExistentKey is error that is returned when key does not exist
var ErrNonExistentKey = errors.New("non existent key")

// ErrInvalidKeyType is error that is returned when the value is not castable to expected type.
var ErrInvalidKeyType = errors.New("invalid key type")

// Adds a value to the context storage, not goroutine safe.
func (c *Context) Add(key string, value any) {
	c.store.Store[key] = value
}

// Gets a value from the context storage, not goroutine safe
func (c *Context) Get(key string) (any, bool) {
	v, ok := c.store.Store[key]
	return v, ok
}

// Gets a value from the context storage (typed), not goroutine safe
func (c *Context) GetTyped[T any](key string) (T, error) {
	value, ok := c.store.Store[key]
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

// Deletes a item in Store
//
// although it is not mandatory to delete all the items in store yourself as it is done automatically at the end but you can do it
func (c *Context) Delete(key string) {
	delete(c.store.Store, key)
}

// the map instance is returned, try not to mutate as it might break internal mechanisms
func (c *Context) Store() map[string]any {
	return c.store.Store
}

type bindHeaderStructCacheData struct {
	fieldData []struct {
		index   int
		keyName string
	}
}

var bindHeaderCache = cache.NewSyncMapCache[reflect.Type, bindHeaderStructCacheData]()

// Binds the headers of a request to a struct provided
//
// Key names for the request is used by the field name or if a `header` tag is present that name is used.
// All fields must be strings or that field will be ignored. An "omitempty" can be used for the header tag,
// this ignores any other type of field without logging an error.
func (c *Context) BindHeaders(target any) {
	rv := reflect.ValueOf(target)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		if c.Mo.Config.LogErrors {
			logger.Mo("Cannot bind headers to a non struct object")
		}
		return
	}
	rt := rv.Type()
	sd, ok := bindHeaderCache.Get(rt)
	if !ok {
		// cache the field names and indexes
		for i := range rt.NumField() {
			t := rt.Field(i)
			if !t.IsExported() {
				continue
			}
			keyName, ok := t.Tag.Lookup("header")
			if !ok {
				keyName = t.Name
			}
			if t.Type.Kind() != reflect.String {
				if c.Mo.Config.LogErrors && keyName != "omitempty" {
					logger.Mo("context: binding headers to a struct requires all fields to be strings! But field " + t.Name + " is of type " + t.Type.Name())
				}
				continue // headers values must be strings strictly
			}
			sd.fieldData = append(sd.fieldData, struct {
				index   int
				keyName string
			}{i, keyName})
		}
		// add to cache
		bindHeaderCache.Add(rt, sd)
	}
	for _, fd := range sd.fieldData {
		v := rv.Field(fd.index)
		value := c.request.Header.Get(fd.keyName)
		if value != "" {
			v.SetString(value)
		}
	}
}

// Decodes the request body into a struct
func (c *Context) DecodeBody(target any) error {
	return json.UnmarshalRead(c.request.Body, target)
}

// Decodes the request body into a struct and validates that using [github.com/impl0x/mo/validator]
func (c *Context) DecodeAndValidateBody(target any) error {
	err := json.UnmarshalRead(c.request.Body, target)
	if err != nil {
		return err
	}
	err = validator.Validate(target)
	if err != nil {
		return err
	}
	return nil
}
