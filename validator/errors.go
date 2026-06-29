package validator

import "fmt"

type ValidationError interface {
	error
	JsonFormat() map[string]any
}

// syntax error in tag formatting
type UserError struct {
	detail string
}

func (de *UserError) Error() string {
	return de.detail
}
func (de *UserError) JsonFormat() map[string]any {
	return map[string]any{
		"message": de.detail,
	}
}

func newUserError(detail string) *UserError {
	return &UserError{
		detail: detail,
	}
}


// required field not found error
type MissingFieldError struct {
	detail      string
	jsonDetails struct {
		field string
	}
}

func newMissingFieldError(field string) *MissingFieldError {
	return &MissingFieldError{
		detail:      fmt.Sprintf("field \"%v\" not found", field),
		jsonDetails: struct{ field string }{field},
	}
}
func (mfe *MissingFieldError) Error() string {
	return mfe.detail
}
func (mfe *MissingFieldError) JsonFormat() map[string]any {
	return map[string]any{
		"message": "missing required field",
		"field":   mfe.jsonDetails.field,
	}
}

// incorrect type for given tag requirements
//
// for example if email requirement has int type
// this is an indirect user error
type TypeError struct {
	detail      string
	jsonDetails struct {
		field    string
		got      string
		expected string
	}
}

func newTypeError(got, expected, field string) *TypeError {
	return &TypeError{
		detail: fmt.Sprintf("field \"%v\" has incorrect type", field),
		jsonDetails: struct {
			field    string
			got      string
			expected string
		}{field, got, expected},
	}
}
func (te *TypeError) Error() string {
	return te.detail
}
func (te *TypeError) JsonFormat() map[string]any {
	return map[string]any{
		"message":  "incorrect type provided",
		"field":    te.jsonDetails.field,
		"provided": te.jsonDetails.got,
		"expected": te.jsonDetails.expected,
	}
}

// failed to validate the field
type ValidateError struct {
	detail      string
	jsonDetails struct {
		message string
		field   string
	}
}

func newValidateError(msg, field string) *ValidateError {
	return &ValidateError{
		detail: "Validation error",
		jsonDetails: struct{ message, field string }{
			msg, field,
		},
	}
}

func (ve *ValidateError) Error() string {
	return ve.detail
}

func (ve *ValidateError) JsonFormat() map[string]any {
	return map[string]any{
		"message": ve.jsonDetails.message,
		"field":   ve.jsonDetails.field,
	}
}