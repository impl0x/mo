package validator

import (
	"reflect"
	"strconv"
	"strings"
)

type errorConfig struct {
	ReturnUserErrors, // change to true if you want validation User errors to be returned in the [GroupedValidationError].
	LogUserErrors bool // logs the user errors.
}

// Config for some error settings
var ErrorConfig = errorConfig{false, true}

type GroupedValidationError []*FieldValidateError

// shows the first error
//
//	type assert to []ValidationError and loop over to see each error.
func (gve GroupedValidationError) Error() string {
	var buf strings.Builder
	zeroElem := *gve[0]
	buf.Grow(128) // eyeballing
	buf.WriteString("Field validation errors, on field ")
	buf.WriteRune('"')
	buf.WriteString(zeroElem.t.name)
	buf.WriteRune('"')
	buf.WriteString(", ")
	buf.WriteRune('"')
	buf.WriteString(zeroElem.Message)
	buf.WriteRune('"')
	if len(gve) > 1 {
		buf.WriteString(", and ")
		buf.WriteString(strconv.Itoa(len(gve) - 1))
		buf.WriteString(" more error")
		if len(gve) > 2 {
			buf.WriteString("s")
		}
	}
	return buf.String()
}

func (gve *GroupedValidationError) Append(elems ...*FieldValidateError) {
	*gve = append(*gve, elems...)
}

// returns a slice of error structs which are compatible with json marshalling, can be safely given to json encoder
func (gve GroupedValidationError) ToJsonStructList() []ValidationErrorJson {
	structList := make([]ValidationErrorJson, len(gve))
	for i, err := range gve {
		structList[i] = ValidationErrorJson{err.Namespace(), err.Error()}
	}
	if len(structList) == 0 {
		return nil
	}
	return structList
}

// this struct is used to format and return validation errors as a json
// when used the ToJsonStructList method on [GroupedValidationError]
type ValidationErrorJson struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// syntax or logical error in the validation tag
type UserError struct {
	parent    string
	fieldName string
	detail    string
}

func newUserError(detail, fieldName string) *UserError {
	return &UserError{
		fieldName: fieldName,
		detail:    detail,
	}
}

func (ue *UserError) Error() string {
	return ue.detail
}

func (ue *UserError) Namespace() string {
	fName := ue.fieldName
	if DefaultNameSpaceSettings.UseLowerCase {
		fName = strings.ToLower(fName)
	}
	return ue.parent + fName
}

// Contains all the information about the failed validation for the field
type FieldValidateError struct {
	Message string // message that is displayed in the error
	parent  string // refers to the struct type name.
	rule    string
	param   string
	t       struct {
		index int
		name  string            // field name
		tag   reflect.StructTag // field tag
		typ   reflect.Type      // field type
	}
	v reflect.Value // field value
}

// returns a new pointer to [FieldValidateError]
// set parent separately as it is not set in the constructor
func newFieldValidateError(msg, rule, ruleParam string, fieldTypeData structFieldData, fieldValue reflect.Value) *FieldValidateError {
	return &FieldValidateError{
		Message: msg,
		rule:    rule,
		param:   ruleParam,
		t: struct {
			index int
			name  string
			tag   reflect.StructTag
			typ   reflect.Type
		}{fieldTypeData.index, fieldTypeData.name, fieldTypeData.tag, fieldTypeData.typ},
		v: fieldValue,
	}
}

// returns the error message
func (ve *FieldValidateError) Error() string {
	return ve.Message
}

// returns the index of the field in the struct
func (ve *FieldValidateError) Index() int {
	return ve.t.index
}

// returns just the field
//
// ex: Age
func (ve *FieldValidateError) Field() string {
	return ve.t.name
}

// returns the tag string in the struct
func (ve *FieldValidateError) Tag() string {
	return string(ve.t.tag)
}

// returns parent struct + field name, uses [DefaultNameSpaceSettings] for the formatting by default
//
// ex: User.Age
func (ve *FieldValidateError) Namespace() string {
	fieldName := ve.t.name
	if DefaultNameSpaceSettings.UseLowerCase {
		fieldName = strings.ToLower(fieldName)
	}
	return ve.parent + fieldName
}

// returns the value of the field, warning it is inefficient for smaller values as go has to box it into an interface
func (ve *FieldValidateError) Value() any {
	return ve.v.Interface()
}

// returns the rule name
func (ve *FieldValidateError) Rule() string {
	return ve.rule
}

// returns the parameter for the rule it failed
func (ve *FieldValidateError) Param() string {
	return ve.param
}

// Kind returns the Field's reflect Kind
//
// eg: time.Time's kind is a struct
func (ve *FieldValidateError) Kind() reflect.Kind {
	return ve.t.typ.Kind()
}

// Type returns the Field's reflect Type
//
// eg: time.Time's type is time.Time
func (ve *FieldValidateError) Type() reflect.Type {
	return ve.t.typ
}
