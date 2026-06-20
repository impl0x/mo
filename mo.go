package mo

import (
	"net/http"
)

type Mo struct {
	router           *Router
	HTTPErrorHandler HTTPErrorHandler
}

func New() *Mo {
	return &Mo{
		router: &Router{},
	}
}
func (m *Mo) GET(path string, handler HandlerFunc) {
	m.router.Routes = append(m.router.Routes, Route{path, http.MethodGet, handler})
}
func (m *Mo) Start(addr string) error {
	return http.ListenAndServe(":8080", m)
}

func (m *Mo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newContext := &Context{
		request:  r,
		response: w,
		Mo:       m,
	}
	route, err := m.router.Route(r.URL.Path, r.Method)
	if err != nil {
		println("Routing error")
		m.HTTPErrorHandler(newContext, err) // either Method wrong or path Not found
		return
	}
	route.Handler(newContext)
}

type Route struct {
	Path    string
	Method  string
	Handler HandlerFunc
}
type Router struct {
	Routes []Route
}

var emptyRoute = Route{}

func (r *Router) Route(path string, method string) (Route, HttpError) {
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

type HandlerFunc func(c *Context) error
type Middleware func(HandlerFunc) HandlerFunc
