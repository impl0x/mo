package mo

import (
	"net/http"
)

type HandlerFunc func(c *Context) error
type Middleware func(HandlerFunc) HandlerFunc

type Mo struct {
	router           Router
	HTTPErrorHandler HTTPErrorHandler // Error handler must also handle nil, because every handler return is at the end handed over to the errorHandler even if its a nil
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
	newContext := &Context{
		request:  r,
		response: &Response{w, false},
		Mo:       m,
	}
	route, err := m.router.Find(r.URL.Path, r.Method)
	if err != nil {
		m.HTTPErrorHandler(newContext, err) // either Method wrong or path Not found
		return
	}
	m.HTTPErrorHandler(newContext, route.Handler(newContext))
}

func (m *Mo) GET(path string, handler HandlerFunc, mi ...Middleware) {
	m.router.Add(Route{path, http.MethodGet, handler, mi})
}
func (m *Mo) POST(path string, handler HandlerFunc, mi ...Middleware) {
	m.router.Add(Route{path, http.MethodPost, handler, mi})
}
func (m *Mo) PATCH(path string, handler HandlerFunc, mi ...Middleware) {
	m.router.Add(Route{path, http.MethodPatch, handler, mi})
}
func (m *Mo) PUT(path string, handler HandlerFunc, mi ...Middleware) {
	m.router.Add(Route{path, http.MethodPut, handler, mi})
}
func (m *Mo) OPTIONS(path string, handler HandlerFunc, mi ...Middleware) {
	m.router.Add(Route{path, http.MethodOptions, handler, mi})
}
func (m *Mo) DELETE(path string, handler HandlerFunc, mi ...Middleware) {
	m.router.Add(Route{path, http.MethodDelete, handler, mi})
}
