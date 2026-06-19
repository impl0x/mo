package mo

type TextSerializer struct{}

func (t TextSerializer) Serialize(c *Context, target any)error{
	
}
func (t TextSerializer) Deserialize(c *Context, target any)error{
	
}
var DefaultTextSerializer = TextSerializer{}

var TEXT=ContentType{
	value: "text/plain",
	formatter: DefaultTextSerializer,
}