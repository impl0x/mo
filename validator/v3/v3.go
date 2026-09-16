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
			return newFieldValidateError(rules.Required.ErrMsg, rules.RuleRequired, "", field, v)
		} else if field.isOptional && v.IsZero() {
			return nil
		} else if field.isDive {
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
		}
		fnWrapper(v, field.ruleFuncs)
	}
	if errs == nil {
		// we cannot directly return errs here because of the way interfaces
		// are implemented in go, here we return an error type so the compiler
		// has to make our GroupedValidationError type into a error interface,
		// and the way go treats interface nil-ability is by seeing if a the
		// interface contains a type or not. In this case even though the value
		// is nil, it still has type data. Therefore checking the result for nil
		// will always result in false even if the value is nil in reality.
		return nil
	}
	return errs
}
