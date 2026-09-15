package mo

import (
	"net/http"
	"strings"
)

type Router interface {
	Find(c *Context, path string, method string) (RouteInfo, HttpError) // finds the route, takes path, method
	Add(RouteInfo)                                                      // adds a route
}

type RouterConfig struct {
	TrimSuffixSlashes bool // trims if there is a leading slash, ex: users/:id/ -> users/:id, remember these 2 are different paths if this option is not enabled.
}

var DefaultRouterConfig = RouterConfig{true} // default config for the router

type RouteInfo struct {
	Path    string      // the path for this route, read only do not mutate
	Method  string      // the method for this route, read only do not mutate
	Handler HandlerFunc // the handler which will be called, read only do not mutate
}

// // o(n) and doesn't support dynamic routing
// type BasicRouter struct {
// 	Routes []*Route
// }

// func NewBasicRouter() *BasicRouter {
// 	return &BasicRouter{}
// }

// func (r *BasicRouter) Add(ro *Route) {
// 	if ro.Path[0] != '/' {
// 		ro.Path = "/" + ro.Path
// 	}
// 	if DefaultRouterConfig.TrimSuffixSlashes {
// 		ro.Path = strings.TrimSuffix(ro.Path, "/")
// 	}
// 	r.Routes = append(r.Routes, ro)
// }
// func (r *BasicRouter) Find(_ *Context, path, method string) (*Route, HTTPError) {
// 	for _, v := range r.Routes {
// 		if path == v.Path {
// 			if method == v.Method {
// 				return v, nil
// 			}
// 			return nil, ErrMethodNotAllowed
// 		}
// 	}
// 	return nil, ErrNotFound
// }

type methodHandlers struct {
	get     HandlerFunc
	post    HandlerFunc
	put     HandlerFunc
	patch   HandlerFunc
	delete  HandlerFunc
	head    HandlerFunc
	options HandlerFunc
	connect HandlerFunc
	trace   HandlerFunc
	any     HandlerFunc
}

// if not a valid http method then a handler is set for "any" method, which triggers on any invalid http method.
func (mh *methodHandlers) add(method string, handler HandlerFunc) {
	switch method {
	case http.MethodGet:
		mh.get = handler
	case http.MethodPost:
		mh.post = handler
	case http.MethodPut:
		mh.put = handler
	case http.MethodPatch:
		mh.patch = handler
	case http.MethodDelete:
		mh.delete = handler
	case http.MethodHead:
		mh.head = handler
	case http.MethodOptions:
		mh.options = handler
	case http.MethodConnect:
		mh.connect = handler
	case http.MethodTrace:
		mh.trace = handler
	default:
		mh.any = handler
	}
}

// We use pass by value in this method because we are not mutating 
// anything to the original instance and the struct is small enough
// to be passed by value, and also because this method is used on
// the find method by the router so it reduces a pointer lookup.

// Returns the method from the string name.
func (mh methodHandlers) fromString(method string) HandlerFunc {
	switch method {
	case http.MethodGet:
		return mh.get
	case http.MethodPost:
		return mh.post
	case http.MethodPut:
		return mh.put
	case http.MethodPatch:
		return mh.patch
	case http.MethodDelete:
		return mh.delete
	case http.MethodHead:
		return mh.head
	case http.MethodOptions:
		return mh.options
	case http.MethodConnect:
		return mh.connect
	case http.MethodTrace:
		return mh.trace
	default:
		return mh.any
	}
}

type node struct {
	segment        string
	methods        methodHandlers
	staticChildren []*node
	paramChild     *node
	wildcardChild  *node
	isHandler      bool
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

func cleanPathString(p string) string {
	if DefaultRouterConfig.TrimSuffixSlashes {
		return strings.TrimSuffix(p, "/")
	}
	return p
}

// Adds a path to the router
func (rr *RadixRouter) Add(r RouteInfo) {
	path := cleanPathString(r.Path)
	currNode := &rr.root
	remainder := path

	wildcardPresent := false // to make sure only one wildcard can exist per url and urls like "users/*/:id/* can be rejected instantly
Outer:
	for {
		// byte traversing logic
		idx := strings.IndexByte(remainder, '/')
		if idx == -1 {
			break
		}
		segment := remainder[:idx]
		if segment == "" {
			panic("mo/router.go/SegmentTreeRouter.Add: empty segment found in URL path, not allowed. Provide valid URLs")
		}
		remainder = remainder[idx+1:]

		// we traverse till we match any static child node
		for _, scn := range currNode.staticChildren {
			if segment == scn.segment {
				currNode = scn
				continue Outer
			}
		}
		// the program arrives here only if the above loop does not match any children for the current given path segment.
		// if it doesn't it means the path already exists and the loop exits after finishing.
		// so we create a child and append it to the current Node's children
		newNode := &node{
			segment: segment,
		}
		switch segment[0] {
		case ':':
			if currNode.paramChild != nil {
				panic(`mo/router.go/SegmentTreeRouter.Add: cannot have more than one parameter type route under one node`)
			}
			if currNode.wildcardChild != nil {
				panic("mo/router.go/SegmentTreeRouter.Add: cannot have a param after a wildcard in a URL")
			}
			currNode.paramChild = newNode
		case '*':
			if currNode.wildcardChild != nil {
				panic("mo/router.go/SegmentTreeRouter.Add: cannot have more than one wildcard type route under one node")
			}
			if wildcardPresent {
				panic(`mo/router.go/SegmentTreeRouter.Add: cannot have more than one wildcard labels in one URL path ("*")`)
			}
			currNode.wildcardChild = newNode
			wildcardPresent = true
		default:
			currNode.staticChildren = append(currNode.staticChildren, newNode)
		}
		currNode = newNode // as there was no child available we assign the current node to the new child we made
	}
	// We add the handlers to the current node.
	// Note: if a user adds another handler for the same path and method then the previous one gets overwritten.
	currNode.methods.add(r.Method, r.Handler)
	if !currNode.isHandler {
		currNode.isHandler = true
	}
}

// Finds a path from the path and method given, returns a [HttpError] instance if not found or wrong method
//
// The returned Route instance is a read only value, do not write to it and expect changes.
func (rr *RadixRouter) Find(c *Context, path, method string) (RouteInfo, HttpError) {
	path = cleanPathString(path)
	remainder := path
	currNode := &rr.root
Outer:
	for {
		// we loop over the parts, ex: [users,:id,posts]
		// byte traversing logic
		idx := strings.IndexByte(remainder, '/')
		if idx == -1 {
			break
		}
		segment := remainder[:idx]
		remainder = remainder[idx+1:]

		// traverse to find any static child first
		for _, scn := range currNode.staticChildren {
			if segment == scn.segment {
				currNode = scn
				continue Outer
			}
		}
		// if not found any static child we look for param child
		if currNode.paramChild != nil {
			currNode = currNode.paramChild
			c.params[currNode.segment[1:]] = segment
			continue Outer
		}
		// if not param we look for wildcard child
		if currNode.wildcardChild != nil {
			currNode = currNode.wildcardChild
			c.params["*"] = segment + remainder // allocates on the heap but its okay as wildcard paths are rare, not important to optimize as of now
			break Outer                         // if it is wildcard match we do not traverse any longer and exit early
		}
		// if there is no static, no param and no wildcards. Then its a dead end.
		return RouteInfo{}, ErrNotFound
	}
	// means there is a node without a handler. it just means not found for the user
	// the node exists but lacks functionality
	if !currNode.isHandler{
		return RouteInfo{}, ErrNotFound // no handlers were ever registered for this node.
	}
	handler := currNode.methods.fromString(method)
	if handler == nil {
		return RouteInfo{}, ErrMethodNotAllowed // if there is no handler returned then we can assume its a wrong method
	}
	return RouteInfo{path, method, handler}, HttpError{}
}
