package validator

import (
	"fmt"
	"reflect"
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

type validator struct {
	target any
	err    *GroupedValidationError

	rv reflect.Value
	rt reflect.Type

	f field
}

type field struct {
	v    reflect.Value
	t    reflect.StructField
	kind reflect.Kind

	rules []string
}

func (vd *validator) init() ValidationError {
	vd.rv = reflect.ValueOf(vd.target)
	if vd.rv.Kind() == reflect.Pointer { // if v is a pointer then we dereference it
		vd.rv = vd.rv.Elem()
	}
	vd.rt = vd.rv.Type()
	if vd.rt.Kind() != reflect.Struct {
		return newUserError("Not a struct") // if not struct we immediately return an error
	}
	return nil
}

func (vd *validator) loop() {
FieldLoop:
	for i := range vd.rv.NumField() {
		vd.f = field{
			v: vd.rv.Field(i),
			t: vd.rt.Field(i),
		}
		vd.f.kind = vd.f.v.Kind()

		if vd.f.kind == reflect.Pointer {
			vd.f.v = vd.f.v.Elem()
		}
		if !vd.f.t.IsExported() {
			continue // if field isn't exported we skip it
		}
		if vd.f.kind == reflect.Struct { // recursively validates any nested structs
			vd.err.Append(Validate(vd.f.v.Interface())) 
			continue
		}
		tag, ok := vd.f.t.Tag.Lookup(validatorTag)
		if !ok {
			continue // if our tag isn't present we skip the field
		}

		vd.f.rules = strings.Split(tag, ",") // "required,email"->["required","email"]

		// checks for field [optional] and [required]
		if vd.f.v.IsZero() {
			for _, ru := range vd.f.rules {
				if ru == required { // i.e. if zero and required we append a error
					vd.err.Append(NewValidateError("Required field not found", vd.f.t.Name, "", string(vd.f.t.Tag), vd.f.v.Interface(), vd.f.v.String()))
					continue FieldLoop // we continue the outer loop
				}
				if ru == optional {
					continue FieldLoop // its optional so we can continue without checks
				}
			}
		}

		for _, rule := range vd.f.rules { // loop over every rule and pass it to the handler
			err := vd.handleNonEqRules(rule)
			if err != nil {
				vd.err.Append(err)
			}
		}

	}
}

// handles the rules without a equal-to sign, required, email, etc.
func (vd *validator) handleNonEqRules(rule string) ValidationError {
	var err ValidationError
	v := vd.f.v.Interface()
	switch rule { // written a theory at the end of this function to make this cleaner
	case required: // we don't deal with required
	case email:
		err = emailRx.validate(vd.f, v)
	case e164:
		err = e164Rx.validate(vd.f, v)
	case url:
		err = urlRx.validate(vd.f, v)
	case uuid:
		err = uuidRx.validate(vd.f, v)
	case alpha:
		err = alphaRx.validate(vd.f, v)
	case alphanum:
		err = alphanumRx.validate(vd.f, v)
	case numeric:
		err = numericRx.validate(vd.f, v)
	case ipv4:
		err = ipv4Rx.validate(vd.f, v)
	case ipv6:
		err = ipv6Rx.validate(vd.f, v)
	default: // eqRules here, min,max,lte,gte,etc.
		err = vd.handleEqRules(rule) // rule: min=2
	}
	return err
}

// there is technically a way to make the above function cleaner and shorter.
// that is by not using a big switch statement and repeating myself with the function calls.
// I could make separate types for nonEq and eq and then loop against every possible rule
// and call the validate function just once. that would work but that would make it more messier
// according to my option that is. So I will let this one slide.
// It is a bit of repeating but it keeps things separated.

func (vd *validator) handleEqRules(eqRule string) ValidationError {
	var isCollection = vd.f.kind == reflect.Slice || vd.f.kind == reflect.Array || vd.f.kind == reflect.Map || vd.f.kind == reflect.String
	split := strings.Split(eqRule, "=")
	if len(split) != 2 {
		return newUserError("Syntax error for tag")
	}
	rule := split[0]
	ruleValueStr := split[1]
	var err ValidationError
	switch rule {
	case min_, max_, gte, lte:
		ruleValue, e := strconv.ParseFloat(ruleValueStr, 64)
		if e != nil {
			return newUserError(fmt.Sprintf("Condition value must be convertible to float64. i.e. ex: min=\"3.14\", the value 3.14 be either uint, int, float64. field: %v", vd.f.t.Name))
		}
		// type checking for the field value here
		if slices.Contains(NumTypes, vd.f.kind) { // checking numTypes, int,uint,float,etc.
			err = vd.handleNumericComparison(rule, vd.f.v.Convert(reflect.TypeFor[float64]()).Float(), ruleValue, "Field value")
		} else if isCollection { // array, slice, string
			err = vd.handleNumericComparison(rule, float64(vd.f.v.Len()), ruleValue, vd.f.kind.String()+" length")
		} else { // unsupported
			err = newUserError(fmt.Sprintf("The field must be either string, collection or numeric. field: %v", vd.f.t.Name))
		}
	case len_:
		ruleValue, e := strconv.Atoi(ruleValueStr)
		if e != nil {
			return newUserError(fmt.Sprintf("len tag value must be int. field: %v", vd.f.t.Name))
		}
		if isCollection {
			if vd.f.v.Len() != ruleValue {
				err = NewValidateError(vd.f.kind.String()+" length must be exactly "+ruleValueStr, vd.f.t.Name, ruleValueStr, string(vd.f.t.Tag), vd.f.v.Interface(), vd.f.v.String())
			}
		} else {
			err = newUserError(fmt.Sprintf("The field must be either string or collection. field: %v", vd.f.t.Name))
		}
	case oneof:
		if vd.f.kind == reflect.String {
			ruleValues := strings.Split(ruleValueStr, " ")
			if !slices.Contains(ruleValues, vd.f.v.String()) {
				err = NewValidateError(fmt.Sprintf("Value must be either one of %v", strings.Join(ruleValues, ", ")), vd.f.t.Name, ruleValueStr, string(vd.f.t.Tag), vd.f.v.Interface(), vd.f.v.String())
			}
		} else {
			err = newUserError("oneof tag must only be used on a string field")
		}

	default: // do oneof and len
		err = newUserError(fmt.Sprintf("Syntax error: Invalid tag value for field %v, rule: %v", vd.f.t.Name, rule))
	}
	return err
}

func (vd *validator) handleNumericComparison(rule string, value float64, ruleValue float64, errorValueName string) ValidationError {
	switch rule {
	case min_, gte:
		if value < ruleValue {
			return NewValidateError(fmt.Sprintf("%v must be more than %v", errorValueName, ruleValue), vd.f.t.Name, fmt.Sprintf("%v", ruleValue), string(vd.f.t.Tag), value, vd.f.v.String())
		}
	case max_, lte:
		if value > ruleValue {
			return NewValidateError(fmt.Sprintf("%v must be less than %v", errorValueName, ruleValue), vd.f.t.Name, fmt.Sprintf("%v", ruleValue), string(vd.f.t.Tag), value, vd.f.v.String())
		}
	}
	return nil
}

func Validate(target any) *GroupedValidationError {
	v := &validator{
		target: target,
		err:    NewGroupedValidationError(),
	}
	err := v.init()
	if err != nil {
		v.err.Append(err)
		return v.err
	}
	v.loop()
	return v.err

}
