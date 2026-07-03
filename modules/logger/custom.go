package logger

import "strings"

var mo = Level{
	Color: Colors["light_blue"],
	Label: "Mo",
}

func Mo(v ...string) {
	log(mo, v)
}

var validator = Level{
	Color: Colors["light_blue"],
	Label: "Validator",
}

func Validator(v ...string){
	log(validator,v)
}

var get = Level{
	Color: Colors["cyan"],
	Label: "GET",
}

var post = Level{
	Color: Colors["orange"],
	Label: "POST",
}

var patch = Level{
	Color: Colors["purple"],
	Label: "PATCH",
}

var put = Level{
	Color: Colors["magenta"],
	Label: "PATCH",
}

var delete = Level{
	Color: Colors["yellow"],
	Label: "DELETE",
}

func Get(v ...string) {
	log(get, v)
}

func Post(v ...string) {
	log(post, v)
}

func Put(v ...string) {
	log(put, v)
}

func Patch(v ...string) {
	log(patch, v)
}

func Delete(v ...string) {
	log(delete, v)
}

func Default(method string, v ...string) {
	log(Level{
		Color: Colors["grey"],
		Label: strings.ToUpper(method)},
		v,
	)
}
