package mo

import "net/http"

type Context struct {
	request *http.Request
	response http.ResponseWriter
	Context map[string]any
}

func NewRequestContext(w http.ResponseWriter, r *http.Request)*Context{
	return &Context{
		r,
		w,
		map[string]any{},
	}
}

func (c *Context) GetResponse()http.ResponseWriter{
	return c.response
}
func (c *Context) GetRequest () *http.Request{
	return c.request
}