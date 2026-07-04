package validator

import "reflect"

type ValidationErrors []FieldError

type FieldError interface {
	Tag() string
	Field() string
	Value() any
	Param() string
	Kind() reflect.Kind
	Type() reflect.Type
	Error() string
}

type fieldError struct{
	errorMessage string
	parent string
	param string
	f *field
}

func NewFieldError(errorMessage, parent, param string, field *field) FieldError {
	return &fieldError{
		errorMessage,parent,param,field,
	}
}

func (fe *fieldError) Tag()string {
	return string(fe.f.t.Tag)
}
func (fe *fieldError) Field()string {
	return fe.parent+fe.f.t.Name
}
func (fe *fieldError) Value()any {
	return fe.f.v.Interface()
}
func (fe *fieldError) Param()string {
	return fe.param
}
func (fe *fieldError) Kind()reflect.Kind {
	return fe.f.kind
}
func (fe *fieldError) Type()reflect.Type {
	return fe.f.t.Type
}
func (fe *fieldError) Error()string {
	return fe.errorMessage
}

type InvalidValidationError struct{
	Message string
}

func newInvalidValidationError(message string) InvalidValidationError {
	return InvalidValidationError{
		message,
	}
}

func (ive InvalidValidationError) Error()string {
	return ive.Message
}
