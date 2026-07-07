package mo

func (m *Mo) Use(mi ...Middleware) {
	m.Middlewares = append(m.Middlewares, mi...)
}

func (m *Mo) AddPostMiddleware(mi ...PostMiddleware) {
	m.PostMiddlewares = append(m.PostMiddlewares, mi...)
}

func (r *Route) Use(mi ...Middleware) {
	r.Middlewares = append(r.Middlewares, mi...)
}

// Make sure to add middlewares first then add the method handlers
func (g *Grouped) Use(mi ...Middleware) {
	g.Middlewares = append(g.Middlewares, mi...)
}
