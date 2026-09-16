package validator

import "github.com/impl0x/go-utils/cache"

type CustomValidatorFunc func(v any, param string) error

var customValidations = cache.NewSyncMapCache[string, CustomValidatorFunc]() // using sync map because its better for heavy reads

// Adds a validation rule to the validator global instance
// param is passed if they exist else empty string is passed
// because no parameters are provided to the function.
//
// the function is only provided the raw interface{} value at runtime. it is the user's job to validate type and return error accordingly,
// the returned error is only used for its .Error() method to be used in the Message parameter for [FieldValidateError],
// do not return User errors as these are treated the same as field errors and will be returned if converted to json struct later.
// if you come up across things such as wrong type being present, then throw a panic because its your codebase that used the tag wrong.
// errors returned can and should be a errors.New instance, but it is recommended to not return a raw errors.New on every fn call,
// because it causes a heap allocation every time errors.New is called, hence used sentinel error variables declared once, either using global variables or closures.
//
// recommended to use only at startup and not at runtime, and not to use duplicate keys globally
func AddCustomValidation(ruleName string, fn CustomValidatorFunc) {
	if fn == nil {
		panic("validator: function passed is nil")
	}
	customValidations.Add(ruleName, fn)
}
