package validator

import (
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/impl0x/go-utils/cache"
	"github.com/impl0x/mo/validator/v3/rules"
)

const TagName = "validate" // the tag name used in struct fields, change if you want. this is used to parse the tags

var structCache = cache.NewSyncMapCache[reflect.Type, structData]()

// represents a struct type metadata
type structData struct {
	fields []structFieldData
}

type structFieldData struct {
	index      int
	name       string
	ruleFuncs  []ValidateFunc
	fieldKind  reflect.Kind
	isRequired bool
	isOptional bool
	isDive     bool
}

// info:
//   - the goal of this function is to not be cheap, this function is expensive and does
//     a lot of expensive operations, but it only runs once per struct and does that with
//     the result of making cheap functions for the actual validator at runtime who can
//     just call the rule funcs without any issues as the hard work is done here,
//     this function has a lot of allocations but the validation at runtime will have close to zero allocations
//
// assumptions:
//   - type provided is of a valid kind struct and not of a pointer to one or else.
//   - param types are supposed to be in sync with the rules package, it is too complex and not worth it to make then dynamic
//
// panics on:
//   - if assumptions fail.
func newStructDataWithCache(structType reflect.Type) (structData, *UserError) {
	if structType.Kind() != reflect.Struct {
		panic("validator cache initialization: structType is not of type struct")
	}
	sd := structData{
		fields: make([]structFieldData, 0), // allow it to grow
	}
	for i := range structType.NumField() {
		currField := structType.Field(i)
		// if field is not exported we dont validate it as we dont have access to the value itself on validation
		if !currField.IsExported() {
			continue
		}
		kind := currField.Type.Kind()
		// if nested struct we recursively cache it again if it doesn't exist in cache already
		if kind == reflect.Struct {
			_, ok := structCache.Get(currField.Type)
			if !ok {
				sd, err := newStructDataWithCache(currField.Type)
				if err != nil {
					return structData{}, err
				}
				structCache.Add(currField.Type, sd)
			}
			continue
		}
		tag := currField.Tag.Get(TagName)
		// if tag not present or empty we skip this field
		if tag == "" {
			continue
		}
		// parse rules now
		fieldRules := strings.Split(tag, ",") // "required,email" -> ["required","email"]

		// create fieldData for each valid field of the struct
		fieldData := structFieldData{
			index:     i,
			name:      currField.Name,
			ruleFuncs: make([]ValidateFunc, 0), // let it grow, as this function runs only once per struct its okay to have expensive operations.
			fieldKind: kind,
		}

		// loop over rules
		for _, rule := range fieldRules {
			// first we need to check for required, optional and dive, because they are special cases,
			// and just set the bool flag for them to true, the validator will handle the rest while validating.
			switch rule {
			case rules.RuleRequired:
				fieldData.isRequired = true
			case rules.RuleOptional:
				fieldData.isOptional = true
			case rules.RuleDive:
				if !slices.Contains(rules.Dive.FieldTypes, kind) {
					return structData{}, newUserError("Invalid field type for rule "+rules.RuleDive, fieldData.name)
				}
				fieldData.isDive = true
			}
			if fieldData.isRequired || fieldData.isOptional || fieldData.isDive {
				continue
			}

			// the plan here is to check what kind of rule it is first
			// and according to that prepare a validate function which
			// will be appended to ruleFuncs at last.

			// variable to store the prepared function, this function will use closures,
			// as closures are just a pointer to memory on the heap the only expensive
			// operation for validator will be to dereference it
			var ruleFunc ValidateFunc

			nonEqRule := rules.NonEqHit(rule) // does not hit for required optional or dive as we already checked them earlier
			// non eq rule hit
			if !nonEqRule.IsNil() {
				// we dont validate field types dynamically as of now as every rule is as string field, so it is hardcoded as of now.
				if kind != reflect.String {
					return structData{}, newUserError(
						"cannot validate a "+rule+" rule against a "+kind.String(),
						fieldData.name,
					)
				}
				// preparing the func
				ruleFunc = func(v reflect.Value) *FieldValidateError {
					if !nonEqRule.Validate(v.String()) {
						return newFieldValidateError(nonEqRule.ErrMsg, rule, "", fieldData, v)
					}
					return nil
				}

				continue
			}
			// either eq rule, or custom rule
			// splitting ruleName and param
			splitStr := strings.Split(rule, "=")
			ruleName := splitStr[0]
			var param string // can be empty string
			if len(splitStr) != 1 {
				param = splitStr[1]
			}
			// checking for custom validations
			customFn, ok := customValidations.Get(ruleName)
			// custom validation hit
			if ok {
				ruleFunc = func(v reflect.Value) *FieldValidateError {
					err := customFn(v.Interface(), param)
					if err != nil {
						return newFieldValidateError(err.Error(), ruleName, param, fieldData, v)
					}
					return nil
				}
				continue
			}

			// try eq rules
			// if param is empty and eq rules cannot have param empty we return error
			if param != "" {
				return structData{}, newUserError("invalid rule", fieldData.name)
			}

			eqRule := rules.EqHit(ruleName)
			// eq rule hit
			if !eqRule.IsNil() {
				// validate field type
				if !slices.Contains(eqRule.FieldTypes, kind) {
					return structData{}, newUserError("Invalid field type for eq rule "+ruleName, fieldData.name)
				}
				switch eqRule.Name {
				case rules.RuleMin, rules.RuleMax, rules.RuleLte, rules.RuleGte, rules.RuleGt, rules.RuleLt: // all these have the same function sig
					// param validation
					if slices.Compare(eqRule.ParamTypes, rules.TypeNumeric) != 0 { // as they have only type numeric, we panic as we cannot validate then
						panic("validator cache initialization: internal error, param types does not equal type numeric for first case")
					}
					floatParam, err := strconv.ParseFloat(param, 64)
					if err != nil {
						return structData{}, newUserError("Param type must be numeric for rule "+ruleName, fieldData.name)
					}

					// creating conversion functions for value
					var convFn func(v reflect.Value) float64
					var errMsgPrefix string
					if slices.ContainsFunc(eqRule.FieldTypes, func(t reflect.Kind) bool {
						return slices.Contains(slices.Concat(rules.TypeString, rules.TypeCollection), t)
					}) {
						errMsgPrefix = "Length"
						convFn = func(v reflect.Value) float64 {
							return float64(v.Len())
						}
					} else if slices.ContainsFunc(eqRule.FieldTypes, func(t reflect.Kind) bool { return slices.Contains(rules.TypeNumeric, t) }) {
						errMsgPrefix = "Value"
						convFn = func(v reflect.Value) float64 {
							return v.Convert(reflect.TypeFor[float64]()).Float()
						}
					}

					// finding the function for the rule
					var fn func(value float64, param float64) bool
					switch eqRule.Name {
					case rules.RuleMin:
						fn = rules.FnMin
					case rules.RuleMax:
						fn = rules.FnMax
					case rules.RuleLte:
						fn = rules.FnLte
					case rules.RuleGte:
						fn = rules.FnGte
					case rules.RuleGt:
						fn = rules.FnGt
					case rules.RuleLt:
						fn = rules.FnLt
					}

					// finally preparing our function
					ruleFunc = func(v reflect.Value) *FieldValidateError {
						if !fn(convFn(v), floatParam) {
							return newFieldValidateError(errMsgPrefix+eqRule.FnErrMsgTemplate(param), ruleName, param, fieldData, v)
						}
						return nil
					}

				case rules.RuleLen:
					// param validation
					if slices.Compare(eqRule.ParamTypes, rules.TypeUInt) != 0 { // as they have only type uint, we panic as we cannot validate then
						panic("validator cache initialization: internal error, param types does not equal type uint for second case")
					}
					t, err := strconv.ParseUint(param, 10, 64)
					if err != nil {
						return structData{}, newUserError("Param type must be unsigned int for rule "+ruleName, fieldData.name)
					}
					uintParam := uint(t)

					// creating conversion functions for value
					convFn := func(v reflect.Value) uint { return uint(v.Len()) }

					ruleFunc = func(v reflect.Value) *FieldValidateError {
						if !rules.FnLen(convFn(v), uintParam) {
							return newFieldValidateError(eqRule.FnErrMsgTemplate(param), ruleName, param, fieldData, v)
						}
						return nil
					}

				case rules.RuleStartswith, rules.RuleEndswith, rules.RuleOneof:
					// param validation
					if slices.Compare(eqRule.ParamTypes, rules.TypeString) != 0 {
						panic("validator cache initialization: internal error, param types does not equal type string for third case")
					}
					strParam := param

					// creating conversion functions for value
					convFn := func(v reflect.Value) string { return v.String() }

					// finding the function we need to use
					var fn func(value string, param string) bool
					switch eqRule.Name {
					case rules.RuleStartswith:
						fn = rules.FnStartswith
					case rules.RuleEndswith:
						fn = rules.FnEndswith
					case rules.RuleOneof:
						fn = rules.FnOneof
					}

					ruleFunc = func(v reflect.Value) *FieldValidateError {
						if !fn(convFn(v), strParam) {
							return newFieldValidateError(eqRule.FnErrMsgTemplate(param), ruleName, param, fieldData, v)
						}
						return nil
					}
				}
			} else {
				// if not even eq rule then it is a invalid rule
				return structData{}, newUserError("invalid rule", fieldData.name)
			}
			// at this point ruleFunc is set with the appropriate validating function, so we append it to ruleFuncs
			fieldData.ruleFuncs = append(fieldData.ruleFuncs, ruleFunc)
		}
		// append the fieldData to the structData
		sd.fields = append(sd.fields, fieldData)
	}
	return sd, nil
}
