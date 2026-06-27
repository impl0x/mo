package validator

import (
	"errors"
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

type ValidationErrors struct {
	
}

var emailRx = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$`)
var urlRx = regexp.MustCompile(`^https?:\/\/(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(?::\d{1,5})?(?:\/[^\s]*)?$`)
var ErrTypeError = errors.New("Wrong type")
var ErrValidation = errors.New("Validation error")
var ErrMissingField = errors.New("Missing field")
var ErrInvalidStruct = errors.New("Invalid struct")

func Validate(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	rt := reflect.TypeOf(v)
	if rt.Kind() != reflect.Struct {
		return ErrInvalidStruct
	}
	for i := 0; i < rv.NumField(); i++ {
		v := rv.Field(i)
		t := rt.Field(i)
		if !t.IsExported() {
			continue
		}
		tag, ok := t.Tag.Lookup(validatorTag)
		if !ok {
			continue
		}
		requirements := strings.Split(tag, ",")

		if ok := slices.Contains(requirements, "required"); ok {
			if v.IsZero() {
				return ErrMissingField
			}
		} else {
			if v.IsZero() {
				continue
			}
		}
		kind := v.Kind()
		for _, r := range requirements {
			switch r {
			case "email":
				if kind != reflect.String {
					return ErrTypeError
				}
				if !emailRx.MatchString(v.String()) {
					return ErrValidation
				}
			case "url":
				if kind != reflect.String {
					return ErrTypeError
				}
				if !urlRx.MatchString(v.String()) {
					return ErrValidation
				}
			case "required":
				continue
			default:
				cons := strings.Split(r, "=")
				if len(cons) != 2 {
					return ErrValidation
				}
				rule := cons[0]
				cods := cons[1]
				switch rule {
				case "min":
					min, err := strconv.ParseFloat(cods, 64)
					if err != nil {
						return ErrTypeError
					}
					if kind == reflect.String {
						val := len(v.String())
						if val < int(min) {
							return ErrValidation
						}
					} else if slices.Contains(NumTypes, kind) {
						if v.Convert(reflect.TypeFor[float64]()).Float() < min {
							return ErrValidation
						}
					} else {
						return ErrTypeError
					}
				case "max":
					max, err := strconv.ParseFloat(cods, 64)
					if err != nil {
						return ErrTypeError
					}
					if kind == reflect.String {
						val := len(v.String())
						if val > int(max) {
							return ErrValidation
						}
					} else if slices.Contains(NumTypes, kind) {
						if v.Convert(reflect.TypeFor[float64]()).Float() > max {
							return ErrValidation
						}
					} else {
						return ErrTypeError
					}
				case "oneof":
					got := fmt.Sprintf("%v", v.Interface())
					allowed := strings.Split(cods, " ")
					if !slices.Contains(allowed, got) {
						return ErrValidation
					}
				default:
					return errors.New("idk")
				}
			}
		}
	}
	return nil
}
