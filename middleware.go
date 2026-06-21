package mo

func (m *Mo) Use(mi... Middleware){
	m.Middlewares = append(m.Middlewares, mi...)
}

func (r *Route) Use(mi... Middleware){
	r.Middlewares = append(r.Middlewares, mi...)
}
