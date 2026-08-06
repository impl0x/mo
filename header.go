package mo

import (
	"errors"
	"net/http"
	"reflect"
)

type headersConfig struct {
	Date           bool
	Content_type   bool
	Content_length bool
}

// Change the values to disable or enable the default headers that net/http sets itself.
var DefaultHeadersConfig = headersConfig{true, true, true}

func (hc *headersConfig) writeHeaders(headers http.Header) {
	if !hc.Content_length {
		headers.Del("content-length")
	} else if !hc.Content_type {
		headers.Del("content-type")
	} else if !hc.Date {
		headers.Del("date")
	}
}

// This is used to manage the headers, like adding, binding, setting, etc.
//
// Remember this is *NOT* goroutine safe, use your own mutexes if you are changing global headers in request handlers.
type HeadersManager struct {
	headers map[string]string
}

func NewHeadersManager() HeadersManager {
	return HeadersManager{}
}

// Adds a header using key and value
func (h *HeadersManager) Add(key, value string) {
	if h.headers == nil {
		h.headers = make(map[string]string)
	}
	h.headers[key] = value
}

// Gets a header value from the key
func (h *HeadersManager) Get(key string) (value string, ok bool) {
	value, ok = h.headers[key]
	return
}

// Sets the map to be the headers, faster alternative than binding a struct
func (h *HeadersManager) SetMap(headers map[string]string) {
	h.headers = headers
}

var ErrNotAStruct = errors.New("cannot bind headers to a non struct object")

// Binds a struct to set the headers globally for every request.
// Must contain `header` tag and be exported.
//
// example:
//
//	XTraceId string `header:"x-trace-id"`
func (h *HeadersManager) Bind(target any) error {
	if h.headers == nil {
		h.headers = make(map[string]string)
	}
	rv := reflect.ValueOf(target)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	rt := rv.Type()
	if rv.Kind() != reflect.Struct {
		return ErrNotAStruct
	}
	for i := range rv.NumField() {
		v := rv.Field(i)
		t := rt.Field(i)
		if !t.IsExported() {
			continue
		}
		if v.Kind() != reflect.String {
			continue // headers values must be strings strictly
		}
		tag, ok := t.Tag.Lookup("header")
		if !ok {
			continue
		}
		h.headers[tag] = v.String()
	}
	return nil
}

func (h *HeadersManager) writeHeaders(headers http.Header) {
	for k, v := range h.headers {
		headers.Set(k, v)
	}
}
