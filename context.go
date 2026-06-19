package mo

import "net/http"

type Context struct {
	request *http.Request
	response http.ResponseWriter
}

func newRequestContext(w http.ResponseWriter, r *http.Request){
	
}