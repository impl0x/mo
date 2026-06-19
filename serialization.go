package mo

type Serializer interface{
	Serialize(*Context, any)error
	Deserialize(*Context, any)error
}

type ContentType struct {
	value string
	formatter Serializer
}

var JSON=ContentType{
	value: "application/json",
	formatter: DefaultJsonSerializer,
}