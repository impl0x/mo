package mo



type HTTPErrorHandler func(*Context, error)

func DefaultHTTPErrorHandler(c *Context, err error) {
	// we sort the error on the basis of if its a HttpError
	// or a client side error and send the appropriate message
}