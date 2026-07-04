package validator

import (
	"fmt"

	"github.com/impl0x/mo/modules/logger"
)

var ReturnUserErrors bool     // change to true if you want validation User errors to be returned in the [GroupedValidationError].
var LogUserErrors bool = true // logs the user errors.

type ValidationError interface {
	Error() string
	JsonFormat() map[string]any
}

type GroupedValidationError struct {
	parent string
	Errors []error
}

func NewGroupedValidationError() *GroupedValidationError {
	return &GroupedValidationError{}
}

func (gve *GroupedValidationError) Error() string {
	return "Validation error, range over .Errors to get individual errors."
}

func (gve *GroupedValidationError) Append(elems ...error) {
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
		switch e := err.(type) {
		case *ValidateError:
			jsonList = append(jsonList, e.JsonFormat())
		case *GroupedValidationError:
			for _,gve:= range e.Errors{
				gve:=gve.(*ValidateError)
				gve.field=e.parent+"."+gve.field
				jsonList = append(jsonList, gve.JsonFormat())
			}
		}
		jsonList = append(jsonList, err.(ValidationError).JsonFormat())
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

// failed to validate the field
type ValidateError struct {
	Message  string
	field    string
	param    string
	tag      string
	value    any
	strValue string
}

func NewValidateError(msg, field, param, tag string, value any, strValue string) *ValidateError {
	return &ValidateError{
		msg, field, param, tag, value, strValue,
	}
}

func (ve *ValidateError) Error() string {
	return fmt.Sprintf("%v | value: %v | field: %v | tag: %v | param: %v", ve.Message, ve.value, ve.field, ve.tag, ve.param)
}

func (ve *ValidateError) JsonFormat() map[string]any {
	return map[string]any{
		"message": ve.Message,
		"field":   ve.field,
	}
}

func (ve *ValidateError) Field() string {
	return ve.field
}
func (ve *ValidateError) Param() string {
	return ve.param
}
func (ve *ValidateError) Tag() string {
	return ve.tag
}
func (ve *ValidateError) Value() any {
	return ve.value
}
func (ve *ValidateError) StrValue() string {
	return ve.strValue
}
