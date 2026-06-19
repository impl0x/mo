package mo


type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type HandlerFunc func(c *Context) error

type Mo struct{

}


