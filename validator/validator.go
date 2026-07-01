package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var NumTypes = []reflect.Kind{
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Int,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Uint,
	reflect.Float32,
	reflect.Float64,
}

const validatorTag = "validate"

var emailRx = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$`)
var urlRx = regexp.MustCompile(`^https?:\/\/(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(?::\d{1,5})?(?:\/[^\s]*)?$`)

// Validates structs using the standard validation tags
//
// examples:
//
// `validate:"required"`
//
// `validate:"required,email"`
//
// `validate:"required,min=3,max=10"`
func Validate(target any) *GroupedValidationError {
	gve := NewGroupedValidationError()
	rv := reflect.ValueOf(target)     // stores the value
	rt := reflect.TypeOf(target)      // stores the type
	if rv.Kind() == reflect.Pointer { // if v is a pointer then we dereference it
		rv = rv.Elem()
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		gve.Append(newUserError("Not a struct"))
		return gve // if not struct we immediately return an error
	}
	for i := range rv.NumField() { // loops over all the fields of the struct
		v := rv.Field(i) // value
		t := rt.Field(i) // type
		kind := v.Kind() // storing the type, will be used a lot below
		if kind == reflect.Pointer {
			v = v.Elem()
		}
		if kind == reflect.Struct { // recursively validates any nested structs
			gve.Append(Validate(v.Interface()).Errors...)
			continue
		}
		if !t.IsExported() {
			continue // if field isn't exported we skip it
		}
		tag, ok := t.Tag.Lookup(validatorTag)
		if !ok {
			continue // if our tag isn't present we skip the field
		}

		requirements := strings.Split(tag, ",") // ex: required,min=2,max=8. we split the requirements/rules

		if ok := slices.Contains(requirements, "required"); ok {
			if v.IsZero() { // if theres a required tag and field is initialized to its zero value
				gve.Append(newValidateError("Required field not found", t.Name)) // we store a missing field error
				continue
			}
		} else {
			if v.IsZero() { // if the field isn't required and its a zero value, we skip it
				continue
			}
		}

		for _, s := range requirements {
			err := validateRule(s, v, t, kind)
			if err != nil {
				gve.Append(err)
			}
		}
	}
	return gve
}

func validateRule(s string, v reflect.Value, t reflect.StructField, kind reflect.Kind) ValidationError {
	switch s {
	case "email":
		if kind != reflect.String {
			return newUserError("Cannot validate \"email\" rule against a " + kind.String())
		}
		if emailRx.MatchString(v.String()) {
			return newValidateError("Not a valid email", t.Name)
		}
	case "url":
		if kind != reflect.String {
			return newUserError("Cannot validate \"url\" rule against a " + kind.String())
		}
		if urlRx.MatchString(v.String()) {
			return newValidateError("Not a valid URL", t.Name)
		}
	case "required":
		return nil // we took care of it before looping, so we can skip now
	default:
		cons := strings.Split(s, "=") // [min,2]
		if len(cons) != 2 {
			return newUserError("Syntax error for tag")
		}
		rule := cons[0] // min
		cods := cons[1] // 2
		switch rule {
		case "min":
			min, err := strconv.ParseFloat(cods, 64) // every number is convertible to float64
			if err != nil {
				return newUserError("min condition must be convertible to float64")
			}
			if kind == reflect.String {
				val := len(v.String())
				if val < int(min) {
					return newValidateError(fmt.Sprintf("Length of string must be more than %v", min), t.Name)
				}
			} else if slices.Contains(NumTypes, kind) {
				if v.Convert(reflect.TypeFor[float64]()).Float() < min {
					return newValidateError(fmt.Sprintf("Value must be less than %v", min), t.Name)
				}
			} else {
				return newUserError("The field must be either string or a type of number")
			}
		case "max":
			max, err := strconv.ParseFloat(cods, 64)
			if err != nil {
				return newUserError("max condition must be convertible to float64")
			}
			if kind == reflect.String {
				val := len(v.String())
				if val > int(max) {
					return newValidateError(fmt.Sprintf("Length of string must be less than %v", max), t.Name)
				}
			} else if slices.Contains(NumTypes, kind) {
				if v.Convert(reflect.TypeFor[float64]()).Float() > max {
					return newValidateError(fmt.Sprintf("Value must be more than %v", max), t.Name)
				}
			} else {
				return newUserError("The field must be either string or a type of number")
			}
		case "oneof":
			got := fmt.Sprintf("%v", v.Interface())
			allowed := strings.Split(cods, " ")
			if !slices.Contains(allowed, got) {
				return newValidateError(fmt.Sprintf("Value must be either one of %v", strings.Join(allowed, ", ")), t.Name)
			}
		default:
			return newUserError("Unknown rule") // if no match
		}
	}
	return nil
}
