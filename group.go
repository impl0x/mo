package mo

import "net/http"

type Grouped struct {
	prefix      string
	Middlewares []Middleware
	m           *Mo
}

func (g *Grouped) add(path string, method string, handler HandlerFunc, mi []Middleware) *Route {
	return g.m.add(g.prefix+path, method, handler, append(g.Middlewares, mi...))
}

func (g *Grouped) GET(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodGet, handler, mi)
}
func (g *Grouped) POST(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodPost, handler, mi)
}
func (g *Grouped) PATCH(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodPatch, handler, mi)
}
func (g *Grouped) PUT(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodPut, handler, mi)
}
func (g *Grouped) OPTIONS(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodOptions, handler, mi)
}
func (g *Grouped) DELETE(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodDelete, handler, mi)
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
func (g *Grouped) Group(prefix string, mi ...Middleware) *Grouped {
	return &Grouped{
		prefix:      g.prefix + prefix,
		Middlewares: append(g.Middlewares, mi...),
		m:           g.m,
	}
}
