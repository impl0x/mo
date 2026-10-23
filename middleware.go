package mo

// Registers a middleware that will run on all paths
func (m *Mo) Use(mi ...Middleware) {
	m.Middlewares = append(m.Middlewares, mi...)
}

// Adds a middleware that runs after all the middlewares and handlers 
// are done running and response has been committed.
func (m *Mo) AddPostMiddleware(mi ...PrePostMiddleware) {
	m.PostMiddlewares = append(m.PostMiddlewares, mi...)
}

// Adds a middleware that runs before anything has ran, before router, middleware or handler.  
func (m *Mo) AddPreMiddleware(mi ...PrePostMiddleware) {
	m.PreMiddlewares = append(m.PreMiddlewares, mi...)
}

// Adds a middleware to this group, all handler registered after adding 
// this middleware will have it run before the handler runs 
func (g *Grouped) Use(mi ...Middleware) {
	g.Middlewares = append(g.Middlewares, mi...)
}
