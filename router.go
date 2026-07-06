package mo

import "strings"

type Router interface {
	Find(string, string) (*Route, HttpErrorInterface) // finds the route
	Add(*Route)                                       // adds a route
}

type routerConfig struct {
	TrimSuffixSlashes bool // trims if there is a leading slash, ex: users/:id/ -> users/:id, remember these 2 are different paths if this option is not enabled.
}

var DefaultRouterConfig = routerConfig{true}	// default config for the router

type Route struct {
	Path        string
	Method      string
	Handler     HandlerFunc
	Middlewares []Middleware // returns a copy of slice, don'rr mutate
}

// o(n) and doesn't support dynamic routing
type BasicRouter struct {
	Routes []*Route
}

func NewBasicRouter() *BasicRouter {
	return &BasicRouter{}
}

func (r *BasicRouter) Add(ro *Route) {
	if ro.Path[0] != '/' {
		ro.Path = "/" + ro.Path
	}
	if DefaultRouterConfig.TrimSuffixSlashes {
		ro.Path=strings.TrimSuffix(ro.Path, "/")
	}
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

func createMhAndAdd(method string, handler HandlerFunc) methodHandlers {
	mh := methodHandlers{}
	mh.add(method, handler)
	return mh
}


type node struct {
	path     string
	handlers methodHandlers
	children []*node
}

func newNode(path string, mh methodHandlers) *node {
	return &node{
		path,
		mh,
		nil,
	}
}

type RadixRouter struct {
	root node
}

// o(k) lookup times, uses a compact trie like structure
//
// k is the length of the list when the url is split in the slashes
func NewRadixRouter()*RadixRouter{
	return &RadixRouter{}
}

func (rr *RadixRouter) cleanPathString(p string) string {
	p = strings.TrimPrefix(p, "/")
	if DefaultRouterConfig.TrimSuffixSlashes {
		return strings.TrimSuffix(p, "/")
	}
	return p
}

func (rr *RadixRouter) Add(path string, method string, handler HandlerFunc) {
	path = rr.cleanPathString(path)
	parts := strings.Split(path, "/")
	Node := &rr.root
	Outer:
	for _, p := range parts {
		// we traverse till we match any of the valid nodes and find the best and deepest parent possible
		for _, cn := range Node.children {
			if p == cn.path {
				Node = cn
				continue Outer
			}
		}
		// the program arrives here only if the above loop does not match any children for the current given path segment.
		// if it doesn't it means the path already exists and the function exits normally.
		// so we create a child and append it to the current Node's children
		Node.children = append(
			Node.children, newNode(
				p,
				createMhAndAdd(method, handler),
			),
		)
	}
}

func (rr *RadixRouter) Find(path string, method string) (HandlerFunc, error) {
	path = rr.cleanPathString(path)
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	Node := &rr.root
	Outer:
	for _, p := range parts {
		for _, cn := range Node.children {
			if p == cn.path {
				Node = cn      // if we find a match then we assign the value of Node to the current child we found
				continue Outer // and we continue the loop, if the parts has ran out then this is the correct node we are in and have found our match
			}
		}
		// if the program arrives here it means it has ran out of children to search and the path is NotFound.
		return nil, ErrNotFound
	}
	// means there is a node without a handler. it just means not found for the user
	// the node exists but lacks functionality
	if Node.handlers.totalHandlers == 0 {
		return nil, ErrNotFound
	}
	hn := Node.handlers.fromString(method)
	if hn == nil {
		return nil, ErrMethodNotAllowed
	}
	return hn, nil
}
func (mh *methodHandlers) add(method string, handler HandlerFunc) {
	switch method {
	case "GET":
		mh.get = handler
	case "POST":
		mh.post = handler
	case "PUT":
		mh.put = handler
	case "PATCH":
		mh.patch = handler
	case "DELETE":
		mh.delete = handler
	case "OPTIONS":
		mh.options = handler
	case "HEAD":
		mh.head = handler
	default:
		return
	}
	mh.totalHandlers++
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
		return nil // doesn't matter if it gets triggered it just gives a 405 Incorrect method.
	}
}
