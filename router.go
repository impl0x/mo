package mo

import "strings"

type Router interface {
	Find(string, string) (*Route, HttpErrorInterface) // finds the route
	Add(*Route)                                       // adds a route
}

type Route struct {
	Path        string
	Method      string
	Handler     HandlerFunc
	Middlewares []Middleware // returns a copy of slice, don't mutate
}

// o(n) and doesn't support dynamic routing
type BasicRouter struct {
	Routes []*Route
}

func NewBasicRouter() *BasicRouter {
	return &BasicRouter{}
}

func (r *BasicRouter) Add(ro *Route) {
	r.Routes = append(r.Routes, ro)
}
func (r *BasicRouter) Find(path string, method string) (*Route, HttpErrorInterface) {
	for _, v := range r.Routes {
		if path == v.Path {
			if method == v.Method {
				return v, nil
			}
			return nil, ErrMethodNotAllowed
		}
	}
	return nil, ErrNotFound
}

type RadixRouter struct {
}

type methodHandlers struct {
	totalHandlers uint8

	get     HandlerFunc
	post    HandlerFunc
	put     HandlerFunc
	patch   HandlerFunc
	options HandlerFunc
	head    HandlerFunc
	delete  HandlerFunc
}

type node struct {
	path     string
	handlers methodHandlers
	parent   *node
	children []*node
}

var root = &node{
	path:   "",
	parent: nil,
	// example data
	children: []*node{
		{
			path: "users",
			children: []*node{
				{
					path:     "posts",
					handlers: methodHandlers{totalHandlers: 1, get: func(ctx *Context) error { return NewHTTPError(200, "test") }},
					children: nil,
				},
				{
					path:     ":id",
					children: nil,
				},
			},
		},
	},
}

type Tree struct {
}

func (t *Tree) Add(path string, method string, handler HandlerFunc) {
	parts := strings.Split(path, "/")
	Node := root
Outer:
	for _, p := range parts {
		for _, cn := range Node.children {
			if p == cn.path {
				Node = cn
				continue Outer
			}
		}
	}
}

func (t *Tree) Find(path string, method string) (HandlerFunc, error) {
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	Node := root
Outer:
	for _, p := range parts {
		for _, cn := range Node.children {
			if p == cn.path {
				Node = cn
				continue Outer
			}
		}
		return nil, ErrNotFound
	}

	if Node.handlers.totalHandlers == 0 {
		return nil, ErrNotFound
	}
	hn:=Node.handlers.fromString(method)
	if hn == nil {
		return nil, ErrMethodNotAllowed
	}
	return hn,nil
}

func (mh *methodHandlers) fromString(method string) HandlerFunc {
	switch method {
	case "GET":
		return mh.get
	case "POST":
		return mh.post
	case "PUT":
		return mh.put
	case "PATCH":
		return mh.patch
	case "DELETE":
		return mh.delete
	case "OPTIONS":
		return mh.options
	case "HEAD":
		return mh.head
	default:
		return nil
	}
}
