package validator

import (
	"reflect"

	"github.com/impl0x/mo/validator/v3/rules"
)

type NameSpaceSettings struct {
	UseLowerCase      bool // true: user.age, false: User.Age
	UseRootStructName bool // true: User.Address.City, false: Address.City
}

// Read [NameSpaceSettings] for docs, this is used in [FieldValidateError].Namespace()
//
// do not mutate while running, only mutate at the start because this is a global variable
//
// else manage your own mutex
var DefaultNameSpaceSettings = NameSpaceSettings{true, false}

type ValidateFunc func(v reflect.Value) *FieldValidateError

// input can either be a struct or a pointer to a struct
//
// output can be either a [*UserError] or [GroupedValidationError]
func Validate(value any) error {
	var errs GroupedValidationError
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return newUserError("value passed is not a struct", "")
	}
	rt := rv.Type()
	sd, ok := structCache.Get(rt)
	var err *UserError
	if !ok {
		// generate new structData for the struct
		sd, err = newStructDataWithCache(rt)
		if err != nil {
			return err
		}
		// add to cache
		structCache.Add(rt, sd)
	}

	// takes a value and validateFunc slice and loops over all the functions to apply it to the value,
	// and appends any errors to the closure errs variable
	fnWrapper := func(v reflect.Value, ruleFuncs []ValidateFunc) {
		for _, fn := range ruleFuncs {
			err := fn(v)
			if err != nil {
				errs.Append(err)
			}
		}
	}

	for _, field := range sd.fields {
		v := rv.Field(field.index)
		if field.isRequired && v.IsZero() {
			// if required and not present we give error
			return newFieldValidateError(rules.Required.ErrMsg, rules.RuleRequired, "", field, v)
		} else if field.isOptional && v.IsZero() {
			// if optional and value not present we skip the field for validation
			continue
		} else if field.isDive {
			// if dive then we perform validation on all elements
			switch field.fieldKind {
			case reflect.Array, reflect.Slice:
				for i := range v.Len() {
					fnWrapper(v.Index(i), field.ruleFuncs)
				}
			case reflect.Map:
				iter := v.MapRange()
				for iter.Next() {
					fnWrapper(iter.Value(), field.ruleFuncs)
				}
			}
			continue
		}
		// apply all rule validations on this field
		fnWrapper(v, field.ruleFuncs)
	}
	if errs == nil {
		// we cannot directly return "errs" from here even if it is nil, because of
		// the way interfaces are implemented in go, in this function we return an 
		// variable of type error as the second value and thus compiler has to box 
		// our [GroupedValidationError] type to an error interface. An interface is 
		// implemented by having two pointers, one to the type and one to the value,
		// and the way nil-ability is treated, that is if a interface is nil is by 
		// checking if the type pointer is nil. So even if the variable having its 
		// value as nil it still has a type, in our case of type [GroupedValidationError].
		// Therefore checking the returned error variable from this function for nil 
		// will always result true no matter the value of the variable, just because 
		// the variable has a type.
		return nil
	}
	// return errors if present.
	return errs
}
