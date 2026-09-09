package validator

import (
	"reflect"
)

type StructValidationEngine struct {
	parent     string
	errs       GroupedValidationError
	structType reflect.Type
	structData structData
}


type field struct{
	t reflect.StructField
}

func (sve *StructValidationEngine) Validate(v any) {
	val:=reflect.ValueOf(v)
	// nvm its too hard for me.
	// i giveup on this project
	
}

func (vd *StructValidationEngine) validateNonEqRuleStr(ruleName nonEqRule, svd nonEqStrValidator) ValidationError {
	if vd.f.kind != reflect.String {
		return newUserError(fmt.Sprintf("cannot validate \"%s\" rule against a %s", ruleName, vd.f.kind.String()), vd.parent, vd.f.fieldName)
	}
	if !svd.Validate(vd.f.v.String()) {
		return newFieldValidateError("Not a valid "+ruleName, "", vd.parent, vd.f)
	}
	return nil
}