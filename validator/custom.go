package validator

// Interface satisfying the custom validation logic, made to store generic structs in a map without insantiation
type CustomValidator interface {
	Validate(f *field, parent string, value any) ValidationError
}

var customValidations = map[string]CustomValidator{}

type customValidator[T any] struct {
	errMsg string
	fn     func(value T) bool
}

func (cd customValidator[T]) Validate(f *field, parent string, value any) ValidationError {
	val, ok := value.(T)
	if !ok {
		return newUserError("Type does not match according to the custom validation generic type registered", f.fieldName, parent)
	}
	ok = cd.fn(val)
	if ok {
		return nil
	}
	return newFieldValidateError(cd.errMsg, "", parent, *f)
}

// Adds a custom validation tag with a function determining the validation
//
//   - tag: the tag name used in the actual struct tags, such as `validator:"tagname"`
//   - errorMessage: the error message to be returned if validation fails
//   - fn: the function which performs the logical operation on the field value deciding the validation for it, return true to pass and false to fail validation.
//
// the field type must be same as the generic type in the func fn. Otherwise a user error is returned by default.
//
// This adds it to a global non mutexed map of tag and functions, do not add at runtime
func AddCustomValidation[T any](tag string, errorMessage string, fn func(value T) bool) {
	customValidations[tag] = customValidator[T]{errorMessage, fn}
}
