package errs

type Validation struct{
	Field string `json:"field"`
	Message string `json:"message"`
}

type Serialization struct{
	Message string `json:"message"`
}