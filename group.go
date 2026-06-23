package mo

import "net/http"

type Group struct {
	prefix string
	Middlewares []Middleware
	m      *Mo
}


func (g *Group) add(path string, method string, handler HandlerFunc, mi []Middleware) *Route {
	return g.m.add(g.prefix+path,method,handler, append(g.Middlewares,mi...))
}

func (g *Group) GET(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodGet, handler, mi)
}
func (g *Group) POST(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodPost, handler, mi)
}
func (g *Group) PATCH(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodPatch, handler, mi)
}
func (g *Group) PUT(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodPut, handler, mi)
}
func (g *Group) OPTIONS(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodOptions, handler, mi)
}
func (g *Group) DELETE(path string, handler HandlerFunc, mi ...Middleware) *Route {
	return g.add(path, http.MethodDelete, handler, mi)
}