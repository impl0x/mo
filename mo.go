package mo

import (
	"net/http"
	"strings"

	"github.com/impl0x/mo/modules/logger"
)

type HandlerFunc func(*Context) error
type Middleware func(HandlerFunc) HandlerFunc
type PostMiddleware func(*Context)

type Mo struct {
	router           Router           // root router
	HTTPErrorHandler HTTPErrorHandler // Error handler must also handle nil, because every handler return is at the end handed over to the errorHandler even if its a nil
	Middlewares      []Middleware
	PostMiddlewares  []PostMiddleware // runs after all the middlewares and handlers have been ran. used to logging or cleaning up, Don't use this to write to response or set status. This also runs when theres a routing error and no handler or middlewares run.
	Headers          *HeadersManager  // Headers, sent in every request
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
		router:           NewRadixRouter(),
		HTTPErrorHandler: DefaultHTTPErrorHandler(false),
		Headers:          DefaultHeadersManager(),
		Config:           DefaultConfig(),
	}
}

func NewWithConfig(router Router, header *HeadersManager, errorHandler HTTPErrorHandler, config *MoConfig) *Mo {
	return &Mo{
		router:           router,
		HTTPErrorHandler: errorHandler,
		Headers:          header,
		Config:           config,
	}
}

func (m *Mo) Start(addr string) error {
	if m.Config.PrintStartMsg {
		logger.Mo("Started Mo HTTP Server.")
	}
	return http.ListenAndServe(addr, m)
}

// the request flow looks like this
//
// r -> global middlewares -> route middlewares -> handler -> error handler -> post middlewares -x-
func (m *Mo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseHeaders := DefaultHeadersManager()
	c := &Context{
		request:         r,
		response:        newResponse(w, m.Headers),
		ResponseHeaders: responseHeaders,
		Mo:              m,
		Store:           make(map[string]any),
		params:          make(map[string]string),
	}
	route, err := m.router.Find(c, strings.TrimSuffix(r.URL.Path, "/"), r.Method)
	if err != nil {
		m.HTTPErrorHandler(c, err) // either Method wrong or path Not found
	} else {
		h := route.handler
		for i := len(m.Middlewares) - 1; i >= 0; i-- {
			h = m.Middlewares[i](h)
		}
		for i := len(route.Middlewares) - 1; i >= 0; i-- {
			h = route.Middlewares[i](h)
		}
		m.HTTPErrorHandler(c, h(c))
	}
	for i := len(m.PostMiddlewares) - 1; i >= 0; i-- {
		m.PostMiddlewares[i](c) // we run post middlewares no matter the failure or status of the request, especially for logging purposes.
	}
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
