package validator

import "github.com/impl0x/mo/modules/logger"

var ReturnUserErrors bool     // change to true if you want validation User errors to be returned in the [GroupedValidationError].
var LogUserErrors bool = true // logs the user errors.

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

// failed to validate the field
type ValidateError struct {
	Message string
	Field   string
}

func NewValidateError(msg, field string) *ValidateError {
	return &ValidateError{
		msg, field,
	}
}

func (ve *ValidateError) Error() string {
	return ve.Message
}

func (ve *ValidateError) JsonFormat() map[string]any {
	return map[string]any{
		"message": ve.Message,
		"field":   ve.Field,
	}
}
