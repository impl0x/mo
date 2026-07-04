package validator

import (
	"reflect"

	"github.com/impl0x/mo/modules/logger"
)

var ReturnUserErrors bool     // change to true if you want validation User errors to be returned in the [GroupedValidationError].
var LogUserErrors bool = true // logs the user errors.

// It is either a [UserError] or a [ValidateError]
type ValidationError interface {
	JsonFormat() map[string]any
}

type GroupedValidationError struct {
	Errors []ValidationError
}

func NewGroupedValidationError() *GroupedValidationError {
	return &GroupedValidationError{}
}

func (gve *GroupedValidationError) Error() string {
	return "Validation error"
}

func (gve *GroupedValidationError) Append(elems ...ValidationError) {
	gve.Errors = append(gve.Errors, elems...)
}

func (gve *GroupedValidationError) JsonFormat() []map[string]any {
	jsonList := make([]map[string]any, 0, len(gve.Errors))
	for _, err := range gve.Errors {
		if e, ok := err.(*UserError); ok {
			if LogUserErrors {
				logger.Validator(e.Error())
			}
			if !ReturnUserErrors {
				continue
			}
		}
		jsonList = append(jsonList, err.JsonFormat())
	}
	return jsonList
}

// syntax error in tag formatting
type UserError struct {
	detail string
}

func newUserError(detail string) *UserError {
	return &UserError{
		detail: detail,
	}
}

func (ue *UserError) Error() string {
	return ue.detail
}
func (ue *UserError) JsonFormat() map[string]any {
	return map[string]any{
		"message": ue.detail,
	}
}

// Contains all the information about the failed validation for the field
type FieldValidateError struct {
	Message string
	param   string
	parent  string
	f       *field
}

func NewFieldValidateError(msg, param, parent string, field *field) *FieldValidateError {
	return &FieldValidateError{
		msg, param,parent, field,
	}
}

func (ve *FieldValidateError) Error() string {
	return ve.Message
}

// Formats the error into a Map for sending as a json response
//
// format: {"message":"String length too short","field":"username"}
func (ve *FieldValidateError) JsonFormat() map[string]any {
	return map[string]any{
		"message": ve.Message,
		"field":   ve.Namespace(),
	}
}

func (ve *FieldValidateError) Tag() string {
	return string(ve.f.t.Tag)
}

// returns just the field
//
// ex: Age
func (ve *FieldValidateError) Field() string {
	return ve.f.t.Name
}

// returns parent struct + field name
//
// ex: User.Age
func (ve *FieldValidateError) Namespace() string {
	return ve.parent + ve.f.t.Name
}

func (ve *FieldValidateError) Value() any {
	return ve.f.v.Interface()
}

func (ve *FieldValidateError) Param() string {
	return ve.param
}

// Kind returns the Field's reflect Kind
//
// eg. time.Time's kind is a struct
func (ve *FieldValidateError) Kind() reflect.Kind {
	return ve.f.kind
}

// Type returns the Field's reflect Type
//
// eg. time.Time's type is time.Time
func (ve *FieldValidateError) Type() reflect.Type {
	return ve.Type()
}
