package mo

import (
	"net/http"

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
	Headers          HeadersManager   // Headers, sent in every request
	Config           MoConfig
}

type MoConfig struct {
	PrintStartMsg bool
	LogErrors     bool
}

func DefaultConfig() MoConfig {
	return MoConfig{
		true, true,
	}
}

// # Returns a new instance of Mo with the default configurations
func New() *Mo {
	return &Mo{
		router:           NewRadixRouter(),
		HTTPErrorHandler: DefaultHTTPErrorHandler(false),
		Headers:          DefaultHeadersManager(),
		Config:           DefaultConfig(),
	}
}

// # Allows you to pass all the configuration that Mo uses on your own.
//
// router: [RadixRouter] / [BasicRouter]
// header: [HeadersManager]
// errorHandler: [DefaultHTTPErrorHandler] / your own implementation. you can implement the [HTTPErrorHandler] function.
// config: [MoConfig]
//
// Make sure to use the constructor functions and not pass in a raw struct directly, for example call NewRadixRouter and not pass in RadixRouter{} by yourself.
//
// although the compiler is satisfied do not pass in a uninitialized struct value, as most data structures need proper initial data to work.
func NewWithConfig(router Router, header HeadersManager, errorHandler HTTPErrorHandler, config MoConfig) *Mo {
	return &Mo{
		router:           router,
		HTTPErrorHandler: errorHandler,
		Headers:          header,
		Config:           config,
	}
}

// Starts listening on the address specified
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
	// Acquire a context from the pool and initialize with values
	c := contextPool.Get().(*Context)
	c.request = r
	c.response = newResponse(w, &m.Headers)
	c.ResponseHeaders = DefaultHeadersManager()
	c.Mo = m
	c.Store = make(map[string]any)
	c.params = make(map[string]string)

	route, err := m.router.Find(c, r.URL.Path, r.Method)
	if err != nil {
		m.HTTPErrorHandler(c, err) // either Method wrong or path Not found
	} else {
		h := route.Handler
		for i := len(m.Middlewares) - 1; i >= 0; i-- { // wrapping with global middlewares
			h = m.Middlewares[i](h)
		}
		for i := len(route.Middlewares) - 1; i >= 0; i-- { // wrapping with route specific middlewares
			h = route.Middlewares[i](h)
		}
		m.HTTPErrorHandler(c, h(c)) // finally we run the handler and pass the result to the error handler
	}
	for i := len(m.PostMiddlewares) - 1; i >= 0; i-- { // running all the post middlewares
		m.PostMiddlewares[i](c) // we run post middlewares no matter the failure or status of the request, especially for logging purposes.
	}
	// we do not bother cleaning the state because Get just overwrites the state.
	contextPool.Put(c)
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

// Used to group several paths together.
//
// example:
//
//	m := mo.New() // new instance
//	v1Group := m.Group("/api/v1") 		// creates a new group for paths starting with /api/v1.
//	authGroup := v1Group.Group("/auth")	// same as above but for /auth, /api/v1/auth in total.
//	authGroup.POST("/login", loginHandler) 	// registers a path finally for POST /api/v1/auth/login.
//	m.Start(":8080") // starts the server
//
// Add middlewares using "Use" before registering paths
func (m *Mo) Group(prefix string, mi ...Middleware) *Grouped {
	return &Grouped{
		prefix:      prefix,
		Middlewares: mi,
		m:           m,
	}
}
