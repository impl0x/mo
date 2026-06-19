package mo

import "errors"

type JSONSerializer struct{}

func (j JSONSerializer) Serialize(c *Context, target any) error {
	return errors.New("test")
}
func (j JSONSerializer) Deserialize(c *Context, target any) error {
	return errors.New("test")
}

var DefaultJsonSerializer = JSONSerializer{}
