package main

import "net/http"

type Receiver struct{

}
func (rec *Receiver)ServeHTTP(w http.ResponseWriter, r *http.Request){
	
}

func main() {
	http.ListenAndServe(":8080",nil)
}
