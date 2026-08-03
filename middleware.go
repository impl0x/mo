package mo

// Registers a middleware that will run on all paths before the respective handler runs
func (m *Mo) Use(mi ...Middleware) {
	m.Middlewares = append(m.Middlewares, mi...)
}

// Adds a middleware that runs after all the middlewares and handlers are done running and response has been committed.
func (m *Mo) AddPostMiddleware(mi ...PostMiddleware) {
	m.PostMiddlewares = append(m.PostMiddlewares, mi...)
}

// Registers a middleware that will run on this path before the handler runs
func (r *Route) Use(mi ...Middleware) {
	r.Middlewares = append(r.Middlewares, mi...)
}

// Make sure to add middlewares first then add the method handlers
func (g *Grouped) Use(mi ...Middleware) {
	g.Middlewares = append(g.Middlewares, mi...)
}
