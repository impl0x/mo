package validator

import "errors"

const validatorTag = "validate"
var ErrValidation = errors.New("Validation error")
func Validate(v any) error {
	 
}
