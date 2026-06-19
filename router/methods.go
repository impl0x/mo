package router

type httpMethod struct {
	name string
}

var (
	GET     = httpMethod{"GET"}
	HEAD    = httpMethod{"HEAD"}
	POST    = httpMethod{"POST"}
	PUT     = httpMethod{"PUT"}
	PATCH   = httpMethod{"PATCH"}
	DELETE  = httpMethod{"DELETE"}
	CONNECT = httpMethod{"CONNECT"}
	OPTIONS = httpMethod{"OPTIONS"}
	TRACE   = httpMethod{"TRACE"}
)
