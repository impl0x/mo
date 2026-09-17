package mo

// Registers a middleware that will run on all paths
func (m *Mo) Use(mi ...Middleware) {
	m.Middlewares = append(m.Middlewares, mi...)
}

// Adds a middleware that runs after all the middlewares and handlers 
// are done running and response has been committed.
func (m *Mo) AddPostMiddleware(mi ...PostMiddleware) {
	m.PostMiddlewares = append(m.PostMiddlewares, mi...)
}

// Adds a middleware to this group, all handler registered after adding 
// this middleware will have it run before the handler runs 
func (g *Grouped) Use(mi ...Middleware) {
	g.Middlewares = append(g.Middlewares, mi...)
}
