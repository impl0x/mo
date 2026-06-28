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

func Validate(target any) []ValidationError {
	var errs []ValidationError
	rv := reflect.ValueOf(target) // stores the value
	rt := reflect.TypeOf(target)  // stores the type
	if rv.Kind() == reflect.Ptr { // if v is a pointer then we dereference it
		rv = rv.Elem()
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		errs = append(errs, newUserError("Not a struct"))
		return errs // if not struct we immediately return an error
	}
	for i := 0; i < rv.NumField(); i++ { // loops over all the fields of the struct
		v := rv.Field(i) // value
		t := rt.Field(i) // type
		kind := v.Kind() // storing the type, will be used a lot below
		if kind == reflect.Ptr {
			v = v.Elem()
		}
		if kind == reflect.Struct { // recursively validates any nested structs
			errs = append(errs, Validate(v.Interface())...)
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
			if v.IsZero() { // if theres a required tag and 54 field is initialized to its zero value
				errs = append(errs, newMissingFieldError(t.Name)) // we store a missing field error
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
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func matchRegexRule(rx *regexp.Regexp, v reflect.Value, kind reflect.Kind, fieldName string) ValidationError {
	if kind != reflect.String {
		return newTypeError(kind.String(), reflect.String.String(), fieldName)
	}
	if !rx.MatchString(v.String()) {
		return newValidateError("Not a valid URL", fieldName)
	}
	return nil
}

func validateRule(s string, v reflect.Value, t reflect.StructField, kind reflect.Kind) ValidationError {
	switch s {
	case "email":
		err := matchRegexRule(emailRx, v, kind, t.Name)
		if err != nil {
			return err
		}
	case "url":
		err := matchRegexRule(urlRx, v, kind, t.Name)
		if err != nil {
			return err
		}
	case "required":
		return nil // we took care of it before looping, so we can skip now
	default:
		cons := strings.Split(s, "=")
		if len(cons) != 2 {
			return newUserError("syntax error for tag")
		}
		rule := cons[0]
		cods := cons[1]
		switch rule {
		case "min":
			min, err := strconv.ParseFloat(cods, 64)
			if err != nil {
				return newUserError("min condition must be convertible to float64")
			}
			if kind == reflect.String {
				val := len(v.String())
				if val < int(min) {
					return newValidateError(fmt.Sprintf("length of string must be more than %v", min), t.Name)
				}
			} else if slices.Contains(NumTypes, kind) {
				if v.Convert(reflect.TypeFor[float64]()).Float() < min {
					return newValidateError(fmt.Sprintf("value must be less than %v", min), t.Name)
				}
			} else {
				return newUserError("the field must be either string or a type of number")
			}
		case "max":
			max, err := strconv.ParseFloat(cods, 64)
			if err != nil {
				return newUserError("max condition must be convertible to float64")
			}
			if kind == reflect.String {
				val := len(v.String())
				if val > int(max) {
					return newValidateError(fmt.Sprintf("length of string must be less than %v", max), t.Name)
				}
			} else if slices.Contains(NumTypes, kind) {
				if v.Convert(reflect.TypeFor[float64]()).Float() > max {
					return newValidateError(fmt.Sprintf("value must be more than %v", max), t.Name)
				}
			} else {
				return newUserError("the field must be either string or a type of number")
			}
		case "oneof":
			got := fmt.Sprintf("%v", v.Interface())
			allowed := strings.Split(cods, " ")
			if !slices.Contains(allowed, got) {
				return newValidateError(fmt.Sprintf("value must be either one of %v", strings.Join(allowed, ", ")), t.Name)
			}
		default:
			return newUserError("unknown rule")
		}
	}
	return nil
}
