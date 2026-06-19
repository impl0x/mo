package router

import (
	
	"net/http"
	
)

type Router struct {
	routes Routes
}

func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	
}



