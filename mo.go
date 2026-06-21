package mo

import (
	"net/http"
)

type HandlerFunc func(c *Context) error
type Middleware func(HandlerFunc) HandlerFunc

type Mo struct {
	router           Router
	HTTPErrorHandler HTTPErrorHandler // Error handler must also handle nil, because every handler return is at the end handed over to the errorHandler even if its a nil
	Middlewares      []Middleware
}

func New() *Mo {
	return &Mo{
		router:           NewSlowRouter(),
		HTTPErrorHandler: DefaultHTTPErrorHandler(false),
	}
}

func (m *Mo) Start(addr string) error {
	return http.ListenAndServe(":8080", m)
}

func (m *Mo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		request:  r,
		response: &Response{w, false},
		Mo:       m,
	}
	route, err := m.router.Find(r.URL.Path, r.Method)
	if err != nil {
		m.HTTPErrorHandler(c, err) // either Method wrong or path Not found
		return
	}
	h := route.Handler
	for _, mi := range m.Middlewares {
		h = mi(h)
	}
	print(len(route.Middlewares))
	for _, mi := range route.Middlewares {
		h = mi(h)
	}
	m.HTTPErrorHandler(c, h(c))
}

func (m *Mo) GET(path string, handler HandlerFunc, mi ...Middleware) *Route {
	r := &Route{path, http.MethodGet, handler, mi}
	m.router.Add(r)
	return r
}
func (m *Mo) POST(path string, handler HandlerFunc, mi ...Middleware) *Route {
	r := &Route{path, http.MethodPost, handler, mi}
	m.router.Add(r)
	return r
}
func (m *Mo) PATCH(path string, handler HandlerFunc, mi ...Middleware) *Route {
	r := &Route{path, http.MethodPatch, handler, mi}
	m.router.Add(r)
	return r
}
func (m *Mo) PUT(path string, handler HandlerFunc, mi ...Middleware) *Route {
	r := &Route{path, http.MethodPut, handler, mi}
	m.router.Add(r)
	return r
}
func (m *Mo) OPTIONS(path string, handler HandlerFunc, mi ...Middleware) *Route {
	r := &Route{path, http.MethodOptions, handler, mi}
	m.router.Add(r)
	return r
}
func (m *Mo) DELETE(path string, handler HandlerFunc, mi ...Middleware) *Route {
	r := &Route{path, http.MethodDelete, handler, mi}
	m.router.Add(r)
	return r
}
