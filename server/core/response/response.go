package response

import "net/http"

type Response struct {
	ContentType    ContentType
	StatusCode     uint16
	Body           any
	ResponseWriter http.ResponseWriter
}

type ContentType struct {
	headerValue string
}

var (
	JSON      = ContentType{"application/json"}
	Text      = ContentType{"text/plain"}
	NoContent = ContentType{}
)

func NewResponse(cType ContentType, body any, w http.ResponseWriter) *Response {
	 
	if _,ok:= body.(string); cType==Text && !ok{
		panic("Invalid type provided")
	}

	
	return &Response{
		ContentType:    cType,
		StatusCode:     200, // default, if user wants he can change the field after initialization.
		Body:           body,
		ResponseWriter: w,
	}
}

func (r *Response) setReqMetaData() {
	if r.ContentType!=NoContent{
		r.ResponseWriter.Header().Set("Content-Type", r.ContentType.headerValue)
	}
	r.ResponseWriter.WriteHeader(int(r.StatusCode))

}
