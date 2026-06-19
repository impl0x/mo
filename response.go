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

func (r *Response) write(c *Context){
	r.ContentType.formatter.Serialize(c,r.Body)
}