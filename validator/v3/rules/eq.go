package rules

import (
	"slices"
	"strings"
)

// ? ----- Eq Rules -----
// ! Eq rules mean the tag validation rules which have an additional parameter to them after an equals to symbol,
// for example: "min=2,max=15"

type EqRule = string

const (
	Min        EqRule = "min"        // INFO: value/length must be minimum of [param]. TYPES, field: string | [TypeCollection] | [TypeNumeric], param: [TypeNumeric]
	Max        EqRule = "max"        // INFO: value/length must be maximum of [param]. TYPES, field: string | [TypeCollection] | [TypeNumeric], param: [TypeNumeric]
	Lte        EqRule = "lte"        // INFO: value/length must be less than or equal to [param]. TYPES, field: string | [TypeCollection] | [TypeNumeric], param: [TypeNumeric]
	Gte        EqRule = "gte"        // INFO: value/length must be greater than or equal to [param]. TYPES, field: string | [TypeCollection] | [TypeNumeric], param: [TypeNumeric]
	Lt         EqRule = "lt"         // INFO: value/length must be less than [param]. TYPES, field: string | [TypeCollection] | [TypeNumeric], param: [TypeNumeric]
	Gt         EqRule = "gt"         // INFO: value/length must be greater than [param] length/value. TYPES, field: string | [TypeCollection] | [TypeNumeric], param: [TypeNumeric]
	Len        EqRule = "len"        // INFO: length must be equal to [param]. TYPES, field: string | [TypeCollection], param: uint
	Startswith EqRule = "startswith" // INFO: string must start with [param]. TYPES, field: string, param: string
	Endswith   EqRule = "endswith"   // INFO: string must end with [param]. TYPES, field: string, param: string
	Oneof      EqRule = "oneof"      // INFO: string must one of [param]. TYPES, field: string, param: string
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
	FnLen = func(value int, param int) bool {
		return value == param
	}
	FnStartswith = func(value string, param string) bool {
		return strings.HasPrefix(value, param)
	}
	FnEndswith = func(value string, param string) bool {
		return strings.HasSuffix(value, param)
	}
	FnOneof = func(value string, param []string) bool {
		return slices.Contains(param, value)
	}
)
