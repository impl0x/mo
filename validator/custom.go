package validator

// Interface satisfying the custom validation logic, made to store generic structs in a map without instantiation
type CustomValidator interface {
	Validate(f *field, parent string, value any) ValidationError
}

var customValidations = map[string]CustomValidator{}

type customValidator[T any] struct {
	fn func(value T) error
}

func (cd customValidator[T]) Validate(f *field, parent string, value any) ValidationError {
	val, ok := value.(T)
	if !ok {
		return newUserError("Type does not match according to the custom validation generic type registered", f.fieldName, parent)
	}
	err := cd.fn(val)
	if err != nil {
		return newFieldValidateError(err.Error(), "", parent, *f)
	}
	return nil
}

// Adds a custom validation tag with a function determining the validation
//
//   - tag: the tag name used in the actual struct tags, such as `validator:"tagname"`
//   - fn: the function which performs the logical operation on the field value deciding the validation for it, return an error with to pass and nil to fail validation, remember that the Error method ont he error instance will be used as the final error message
//
// the field type must be same as the generic type in the func fn. Otherwise a user error is returned by default.
//
// This adds it to a global non mutexed map of tag and functions, do not add at runtime
func AddCustomValidation[T any](tag string, fn func(value T) error) {
	customValidations[tag] = customValidator[T]{fn}
}
