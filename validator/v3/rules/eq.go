package rules

import (
	"reflect"
	"slices"
	"strings"
)

type EqRule struct {
	Name             string
	FieldTypes       []reflect.Kind
	ParamTypes       []reflect.Kind
	FnErrMsgTemplate func(param string) string
}

func (er EqRule) String() string {
	return er.Name
}

func (er EqRule) AcceptedFieldTypes() []reflect.Kind {
	return er.FieldTypes
}

func (er EqRule) IsNil() bool {
	return er.Name == ""
}

const (
	RuleMin        = "min"
	RuleMax        = "max"
	RuleLte        = "lte"
	RuleGte        = "gte"
	RuleLt         = "lt"
	RuleGt         = "gt"
	RuleLen        = "len"
	RuleStartswith = "startswith"
	RuleEndswith   = "endswith"
	RuleOneof      = "oneof"
)

var (
	// value/length must be minimum of [param]
	Min = EqRule{
		RuleMin,
		slices.Concat(TypeString, TypeCollection, TypeNumeric),
		TypeNumeric,
		func(param string) string { return "must be minimum of " + param },
	}
	// value/length must be maximum of [param]
	Max = EqRule{
		RuleMax,
		slices.Concat(TypeString, TypeCollection, TypeNumeric),
		TypeNumeric,
		func(param string) string { return "must be maximum of " + param },
	}
	// value/length must be less than or equal to [param].
	Lte = EqRule{
		RuleLte,
		slices.Concat(TypeString, TypeCollection, TypeNumeric),
		TypeNumeric,
		func(param string) string { return "must be less than or equal to " + param },
	}
	// value/length must be greater than or equal to [param].
	Gte = EqRule{
		RuleGte,
		slices.Concat(TypeString, TypeCollection, TypeNumeric),
		TypeNumeric,
		func(param string) string { return "must be greater than or equal to" + param },
	}
	// value/length must be less than [param].
	Lt = EqRule{
		RuleLt,
		slices.Concat(TypeString, TypeCollection, TypeNumeric),
		TypeNumeric,
		func(param string) string { return "must be less than " + param },
	}
	// value/length must be greater than [param] length/value.
	Gt = EqRule{
		RuleGt,
		slices.Concat(TypeString, TypeCollection, TypeNumeric),
		TypeNumeric,
		func(param string) string { return "must be greater than " + param },
	}
	// length must be equal to [param].
	Len = EqRule{
		RuleLen,
		slices.Concat(TypeString, TypeCollection),
		TypeUInt,
		func(param string) string { return "Length must be equal to " + param },
	}
	// string must start with [param].
	StartsWith = EqRule{
		RuleStartswith,
		TypeString,
		TypeString,
		func(param string) string { return "String must start with " + param },
	}
	// string must end with [param].
	EndsWith = EqRule{
		RuleEndswith,
		TypeString,
		TypeString,
		func(param string) string { return "String must end with " + param },
	}
	// string must be one of [param], separated by single space
	Oneof = EqRule{
		RuleOneof,
		TypeString,
		TypeString,
		func(param string) string { return "String must be one of " + param },
	}
)

// ? ----- Validation functions -----
//
// these funcs are just one line bool returns but i wrote this to have them rule logic separated from reflection and validator logic
// value is the field value in the struct and param is the tag parameter
// example: in a field username string `validate:"min=2", where username is populated with "test" in a instance,
// value will be len("test") = 4 and param will be 2 from the "min=2" tag.
var (
	FnMin = func(value float64, param float64) bool {
		return value > param
	}
	FnMax = func(value float64, param float64) bool {
		return value < param
	}
	FnGte = FnMin // logically the same function
	FnLte = FnMax // ~
	FnLt  = func(value float64, param float64) bool {
		return value <= param
	}
	FnGt = func(value float64, param float64) bool {
		return value >= param
	}
	FnLen = func(value uint, param uint) bool {
		return value == param
	}
	FnStartswith = func(value string, param string) bool {
		return strings.HasPrefix(value, param)
	}
	FnEndswith = func(value string, param string) bool {
		return strings.HasSuffix(value, param)
	}
	// splits the string based on spaces, multiple spaces are handled in case of typos
	FnOneof = func(value string, param string) bool {
		st := 0
		param = param + " "
		for i, c := range param {
			if c == ' ' {
				if st < i && param[st:i] == value {
					return true
				}
				st = i + 1
			}
		}
		return false
	}
)

func EqHit(ruleName string) EqRule {
	switch ruleName {
	case RuleMin:
		return Min
	case RuleMax:
		return Max
	case RuleGte:
		return Gte
	case RuleLte:
		return Lte
	case RuleLt:
		return Lt
	case RuleGt:
		return Gt
	case RuleLen:
		return Len
	case RuleStartswith:
		return StartsWith
	case RuleEndswith:
		return EndsWith
	case RuleOneof:
		return Oneof
	default:
		return EqRule{}
	}
}
