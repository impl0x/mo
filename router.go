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

// todo: figure out a way to add a dynamic router which supports dynamic and wildcard paths
