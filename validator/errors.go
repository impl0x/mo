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

// only 2 structs satisfy this interface, that is [UserError] / [FieldValidateError]
type ValidationError interface{
	error
	Namespace() string
}

type GroupedValidationError struct {
	Errors []ValidationError // you can type assert for [UserError] / [FieldValidateError] safely, if ReturnUserErrors singleton bool is false then only FieldValidateError will be present.
}

func (gve GroupedValidationError) Error() string {
	return "Validation error, please loop over Errors to see each error."
}

func (gve *GroupedValidationError) Append(elems ...ValidationError) {
	gve.Errors = append(gve.Errors, elems...)
}

// returns a slice of error structs which are compatible with json marshalling, can be safely given to json encoder
func (gve GroupedValidationError) ToJsonStructList() []ValidationErrorJson {
	structList := make([]ValidationErrorJson, len(gve.Errors))
	for i, err := range gve.Errors {
		if _,ok:=err.(UserError);ok{
			if ErrorConfig.LogUserErrors{
				logger.Validator("user error: "+err.Error())
			}
			if !ErrorConfig.ReturnUserErrors{
				continue
			}
		}
		structList[i].Field=err.Namespace()
		structList[i].Message=err.Error()
	}
	return structList
}

type ValidationErrorJson struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (vej ValidationErrorJson) Error() string {
	return vej.Message
}

// The reason [UserError] and [FieldValidateError] are separate and have the same fields is to signify the type of the error

// syntax error in tag formatting
type UserError struct {
	parent    string
	fieldName string
	detail    string
}

func newUserError(parent, fieldName, detail string) UserError {
	return UserError{
		detail: detail,
	}
}

func (ue UserError) Error() string {
	return ue.detail
}

func (ue UserError) Namespace() string {
	fName := ue.fieldName
	if DefaultNameSpaceSettings.UseLowerCase {
		fName = strings.ToLower(fName)
	}
	return ue.parent + fName
}

// Contains all the information about the failed validation for the field
type FieldValidateError struct {
	Message string
	param   string
	parent  string
	f       *field
}

func newFieldValidateError(msg, param, parent string, field field) FieldValidateError {
	return FieldValidateError{
		msg, param, parent, &field,
	}
}

func (ve FieldValidateError) Error() string {
	return ve.Message
}

func (ve FieldValidateError) Tag() string {
	return string(ve.f.t.Tag)
}

// returns just the field
//
// ex: Age
func (ve FieldValidateError) Field() string {

	return ve.f.t.Name
}

// returns parent struct + field name
//
// ex: User.Age
func (ve FieldValidateError) Namespace() string {
	fName := ve.f.fieldName
	if DefaultNameSpaceSettings.UseLowerCase {
		fName = strings.ToLower(fName)
	}
	return ve.parent + fName
}

func (ve FieldValidateError) Value() any {
	return ve.f.v.Interface()
}

func (ve FieldValidateError) Param() string {
	return ve.param
}

// Kind returns the Field's reflect Kind
//
// eg. time.Time's kind is a struct
func (ve FieldValidateError) Kind() reflect.Kind {
	return ve.f.kind
}

// Type returns the Field's reflect Type
//
// eg. time.Time's type is time.Time
func (ve FieldValidateError) Type() reflect.Type {
	return ve.Type()
}
