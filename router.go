package mo

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


// o(n)
type BasicRouter struct {
	Routes      []*Route
	Middlewares []Middleware
}

func NewSlowRouter() *BasicRouter {
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



// type RadixRouter struct {
// 	root *node
// }

// type pathKind uint8

// const (
// 	static pathKind = iota
// 	param
// 	wildcard
// )

// type methodHandlers struct {
// 	get, post, put, patch, delete, options, head HandlerFunc
// }

// type node struct {
// 	part           string
// 	methodHandlers *methodHandlers
// 	kind           pathKind
// 	parent         *node
// 	child          *node
// }

// func NewRadixRouter() *RadixRouter {
// 	return &RadixRouter{
// 		root: &node{
// 			part:           "/",
// 			methodHandlers: nil,
// 			kind:           static,
// 			parent:         nil,
// 			child:          nil,
// 		},
// 	}
// }

// func newMethodHandler(method string, handler HandlerFunc) *methodHandlers{
// 	mh:=methodHandlers{}
// 	switch method {
// 	case http.MethodGet:
// 		mh.get = handler
// 	case http.MethodPost:
// 		mh.post = handler
// 	case http.MethodPut:
// 		mh.put = handler
// 	case http.MethodPatch:
// 		mh.patch = handler
// 	case http.MethodDelete:
// 		mh.delete = handler
// 	case http.MethodOptions:
// 		mh.options = handler
// 	case http.MethodHead:
// 		mh.head = handler
// 	}
// 	return &mh
// }

// func (r *RadixRouter) Add(ro Route) {
// 	pathSplits:=strings.Split(ro.Path, "/")
// 	if r.root.child == nil {
// 		var n *node
// 		for _,v:=range pathSplits{
// 			n=&node{
// 				part: v,
// 				// kind: ,,
// 			}
// 		}
// 		r.root.child = &node{
// 			part:           ro.Path,
// 			methodHandlers: newMethodHandler(ro.Method,ro.Handler),
// 			// kind:
// 		}
// 	}
// }


