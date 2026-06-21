package mo

type Router interface {
	Find(string, string) (Route, HttpErrorInterface)
	Add(Route)
}

type Route struct {
	Path        string
	Method      string
	Handler     HandlerFunc
	Middlewares []Middleware
}

var emptyRoute = Route{}

// o(n)
type SlowRouter struct {
	Routes []Route
}

func NewDefaultRouter() *SlowRouter {
	return &SlowRouter{}
}

func (r *SlowRouter) Find(path string, method string) (Route, HttpErrorInterface) {
	for _, v := range r.Routes {
		if path == v.Path {
			if method == v.Method {
				return v, nil
			}
			return emptyRoute, ErrMethodNotAllowed
		}
	}
	return emptyRoute, ErrNotFound
}

func (r *SlowRouter) Add(ro Route) {
	r.Routes = append(r.Routes, ro)
}
