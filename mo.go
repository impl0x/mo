package mo

import (
	"net/http"

	"github.com/impl0x/mo/modules/logger"
)

type HandlerFunc func(c *Context) error
type Middleware func(HandlerFunc) HandlerFunc

type Mo struct {
	router           Router           // root router
	HTTPErrorHandler HTTPErrorHandler // Error handler must also handle nil, because every handler return is at the end handed over to the errorHandler even if its a nil
	Middlewares      []Middleware
	Headers          map[string]string // default headers
	Config           *MoConfig
}

type MoConfig struct {
	PrintStartMsg bool
	LogErrors     bool
}

func DefaultConfig() *MoConfig {
	return &MoConfig{
		true, true,
	}
}

func New() *Mo {
	return &Mo{
		router:           NewSlowRouter(),
		HTTPErrorHandler: DefaultHTTPErrorHandler(false),
		Headers:          map[string]string{},
		Config:           DefaultConfig(),
	}
}

func NewWithConfig(router Router, errorHandler HTTPErrorHandler, config *MoConfig) *Mo {
	return &Mo{
		router:           router,
		HTTPErrorHandler: errorHandler,
		Headers:          map[string]string{},
		Config:           config,
	}
}

func (m *Mo) Start(addr string) error {
	if m.Config.PrintStartMsg {
		logger.Mo("Started Mo HTTP Server.")
	}
	return http.ListenAndServe(":8080", m)
}

func (m *Mo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := &Context{
		request:  r,
		response: &Response{w, false, m.Headers, map[string]string{}},
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
	for _, mi := range route.Middlewares {
		h = mi(h)
	}
	m.HTTPErrorHandler(c, h(c))
}

func (m *Mo) add(path string, method string, handler HandlerFunc, mi []Middleware) *Route {
	r := &Route{path, method, handler, mi}
	m.router.Add(r)
	return r
}

func (m *Mo) GET(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return m.add(path, http.MethodGet, handler, mi)
}
func (m *Mo) POST(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return m.add(path, http.MethodPost, handler, mi)
}
func (m *Mo) PATCH(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return m.add(path, http.MethodPatch, handler, mi)
}
func (m *Mo) PUT(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return m.add(path, http.MethodPut, handler, mi)
}
func (m *Mo) OPTIONS(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return m.add(path, http.MethodOptions, handler, mi)
}
func (m *Mo) DELETE(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return m.add(path, http.MethodDelete, handler, mi)
}
func (m *Mo) HEAD(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return m.add(path, http.MethodHead, handler, mi)
}

// Add middlewares using "Use" before registering paths
func (m *Mo) Group(prefix string, mi ...Middleware) *Grouped {
	return &Grouped{
		prefix:      prefix,
		Middlewares: mi,
		m:           m,
	}
}
