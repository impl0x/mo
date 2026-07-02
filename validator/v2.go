package validator

import (
	"reflect"
	"slices"
	"strings"
)

type validatorRule = string

const (
	required validatorRule = "required"
	email    validatorRule = "email"
	url      validatorRule = "url"
	uuid     validatorRule = "uuid"
	ipv4     validatorRule = "ipv4"
	ipv6     validatorRule = "ipv6"
)

type validator struct {
	target any
	err    *GroupedValidationError

	rv reflect.Value
	rt reflect.Type
}

type field struct {
	v    reflect.Value
	t    reflect.StructField
	kind reflect.Kind

	rules []string
}

func (vd *validator) init() {
	vd.rv = reflect.ValueOf(vd.target)
	if vd.rv.Kind() == reflect.Pointer { // if v is a pointer then we dereference it
		vd.rv = vd.rv.Elem()
	}
	vd.rt = vd.rv.Type()
	if vd.rt.Kind() != reflect.Struct {
		vd.err.Append(newUserError("Not a struct")) // if not struct we immediately return an error
	}
}

func (vd *validator) loop() {
	for i := range vd.rv.NumField() {
		f := field{
			v: vd.rv.Field(i),
			t: vd.rt.Field(i),
		}
		f.kind = f.v.Kind()

		if f.kind == reflect.Pointer {
			f.v = f.v.Elem()
		}
		if f.kind == reflect.Struct { // recursively validates any nested structs
			vd.err.Append(Validate(f.v.Interface()).Errors...)
			continue
		}
		if !f.t.IsExported() {
			continue // if field isn't exported we skip it
		}
		tag, ok := f.t.Tag.Lookup(validatorTag)
		if !ok {
			continue // if our tag isn't present we skip the field
		}

		f.rules = strings.Split(tag, ",")

		if f.v.IsZero() { // if zero we skip
			if ok := slices.Contains(f.rules, required); ok { // if required we add a error
				vd.err.Append(NewValidateError("Required field not found", f.t.Name))
			}
			continue
		}

		for _, s := range f.rules {
			err := f.handleRules(s)
			if err != nil {
				vd.err.Append(err)
			}
		}

	}
}

func (f *field) handleRules(s string) *ValidateError {
	// TODO: handle the different rules of validation like email,url,etc. using methods defined on field struct itself.
}

func Validate_(target any) *GroupedValidationError {
	v := &validator{
		target: target,
		err:    NewGroupedValidationError(),
	}
	v.init()
	v.loop()
	return v.err

}
