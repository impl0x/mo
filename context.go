package mo

import (
	"encoding/json"
	"net/http"

	"github.com/impl0x/mo/modules/logger"
)

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	Mo       *Mo
}

func (c *Context) writeContentType(value string) {
	header := c.response.Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

// Redirect redirects the request to a provided URL with status code.
func (c *Context) Redirect(code int, url string)  {
	if code < 300 || code > 308 {
		c.Mo.HTTPErrorHandler(c,ErrInternalServerError)
		return
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	
}

// NoContent sends a response with no body and a status code.
func (c *Context) NoContent(code int) {
	c.response.WriteHeader(code)
}

// Blob sends a blob response with status code and content type.
func (c *Context) Blob(code int, contentType string, b []byte) {
	c.writeContentType(contentType)
	c.response.WriteHeader(code)
	writeResp(c.response, b)
}

// JSON sends a JSON response with status code.
func (c *Context) JSON(code int, target any) {
	v,err:=json.Marshal(target)
	if err != nil {
		c.Mo.HTTPErrorHandler(c,ErrInternalServerError)
		return 
	}
	c.writeContentType(MIMEApplicationJSON)
	c.response.WriteHeader(code)
	writeResp(c.response,v)
}

func writeResp(resp http.ResponseWriter,b []byte){
	_, err := resp.Write(b)
	if err != nil {
		logger.Error("Client disconnected! couldn't write response")
		return 
	}
}
