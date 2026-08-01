package validator

import (
	"reflect"
	"strings"

	"github.com/impl0x/mo/modules/logger"
)

type errorConfig struct {
	ReturnUserErrors, // change to true if you want validation User errors to be returned in the [GroupedValidationError].
	LogUserErrors bool // logs the user errors.
}

// Config for some error settings
var ErrorConfig = errorConfig{false, true}

// It is either a [UserError] or a [FieldValidateError]
type ValidationError interface {
	JsonFormat() map[string]any
}

type GroupedValidationError struct {
	Errors []ValidationError // you can type assert for [UserError] / [FieldValidateError], if ReturnUserErrors singleton bool is false then only FieldValidateError will be present.
}

func NewGroupedValidationError() *GroupedValidationError {
	return &GroupedValidationError{}
}

func (gve *GroupedValidationError) Error() string {
	return "Validation error, please loop over Errors to see each error."
}

func (gve *GroupedValidationError) Append(elems ...ValidationError) {
	gve.Errors = append(gve.Errors, elems...)
}

func (gve *GroupedValidationError) JsonFormat() []map[string]any {
	jsonList := make([]map[string]any, 0, len(gve.Errors))
	for _, err := range gve.Errors {
		if e, ok := err.(*UserError); ok {
			if ErrorConfig.LogUserErrors {
				logger.Validator(e.Error())
			}
			if !ErrorConfig.ReturnUserErrors {
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
	f       field
}

func NewFieldValidateError(msg, param, parent string, field field) *FieldValidateError {
	return &FieldValidateError{
		msg, param, parent, field,
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
	fName := ve.f.t.Name
	if DefaultNameSpaceSettings.UseLowerCase {
		fName = strings.ToLower(fName)
	}
	return ve.parent + fName
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
