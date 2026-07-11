package mo

import (
	"strings"
)

type Router interface {
	Find(c *Context, path string, method string) (*Route, HttpErrorInterface) // finds the route, takes path, method
	Add(*Route)                                                               // adds a route
}

type routerConfig struct {
	TrimSuffixSlashes bool // trims if there is a leading slash, ex: users/:id/ -> users/:id, remember these 2 are different paths if this option is not enabled.
}

var DefaultRouterConfig = routerConfig{true} // default config for the router

type Route struct {
	path        string
	method      string
	handler     HandlerFunc
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
	if ro.path[0] != '/' {
		ro.path = "/" + ro.path
	}
	if DefaultRouterConfig.TrimSuffixSlashes {
		ro.path = strings.TrimSuffix(ro.path, "/")
	}
	r.Routes = append(r.Routes, ro)
}
func (r *BasicRouter) Find(_ *Context, path string, method string) (*Route, HttpErrorInterface) {
	for _, v := range r.Routes {
		if path == v.path {
			if method == v.method {
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

type node struct {
	path       string
	handlers   methodHandlers
	middleware []Middleware
	children   []*node
}

type RadixRouter struct {
	root node
}

// o(k) lookup times, uses a compact trie like structure
//
// k is the length of the list when the url is split in the slashes
func NewRadixRouter() *RadixRouter {
	return &RadixRouter{}
}

func (rr *RadixRouter) cleanPathString(p string) string {
	if DefaultRouterConfig.TrimSuffixSlashes {
		return strings.TrimSuffix(p, "/")
	}
	return p
}

func (rr *RadixRouter) Add(r *Route) {
	r.path = rr.cleanPathString(r.path)
	parts := strings.Split(r.path, "/")
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
		// if it doesn't it means the path already exists and the loop exits after finishing.
		// so we create a child and append it to the current Node's children
		nn := &node{
			path: p,
		}
		Node.children = append(
			Node.children, nn,
		)
		Node = nn // as there was no child available we assign the current node to the new child we made
	}
	// We add the handlers to the deepest node.
	Node.handlers.add(r.method, r.handler)
	// we also add the middlewares with it.
	Node.middleware = append(Node.middleware, r.Middlewares...)
	// Note: if a user adds another handler for the same path and method then the previous one gets overwritten.
}

func (rr *RadixRouter) Find(c *Context, path string, method string) (*Route, HttpErrorInterface) {
	path = rr.cleanPathString(path)
	parts := strings.Split(path, "/")
	Node := &rr.root
	var wildcard *node
	var param *node

Outer:
	for _, p := range parts { // we loop over the parts, ex: [users,:id,posts]
		param = nil
		for _, cn := range Node.children { // we check each node's children to find a match, if children slice is nil it automatically doesn't start the loop
			if p == cn.path {
				Node = cn      // if we find a match then we assign the value of Node to the current child we found
				continue Outer // and we continue the loop, if the parts has ran out then this is the correct node we are in and have found our match
			} else if cn.path[0] == ':' { // this is a param path, ex: ":id"
				c.params[cn.path[1:]] = p // we store the parameter value in the map, (trimming the starting colon ofc.)
				param = cn                // we do not continue the loop here or in wildcard because we are not finishing searching through it, we only continue in the static child because that is the best possible scenario
			} else if cn.path == "*" { // this is a wildcard path, ex: "*"
				wildcard = cn // we just store the node to the type of node we found, and then check it later if we found any to then act accordingly.
			}
		}
		// if the program arrives here it means it has ran out of static children to search.
		if param != nil { // Now, if there is a param we assign node to that and continue
			Node = param
			continue Outer
		} else if wildcard != nil { // and if there is no static and no param, then we assign the node to be the last wildcard we found.
			Node = wildcard
			continue Outer
		}
		// if theres no static, no param and no wildcards. Then its a dead end.
		return nil, ErrNotFound
	}
	// means there is a node without a handler. it just means not found for the user
	// the node exists but lacks functionality
	if Node.handlers.totalHandlers == 0 {
		return nil, ErrNotFound // no handlers were ever registered for this node.
	}
	hn := Node.handlers.fromString(method)
	if hn == nil {
		return nil, ErrMethodNotAllowed // if there is no handler returned then we can assume its a wrong method
	}
	return &Route{path, method, hn, Node.middleware}, nil
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
	// well there is a possible overflow where if the user adds 256 handlers for the same path.... (the variable is of 8 bits)
	// I mean there isn't 256 methods so he/she will just have to keep overwriting the handlers.
	// and well if it overflows to 0 then it will return a 404 not found. just saying.
}

// Returns the method from the string name.
// As the method string is usually passed from the http package containing the const names for methods,
// we check the capitalized versions.
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
