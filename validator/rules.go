package validator

import (
	"reflect"
	"regexp"
)

type regexRule struct {
	name  string
	regEx *regexp.Regexp
}

type validatorRule = string

const ( // these are all nonEq rules
	required validatorRule = "required"
	optional validatorRule = "optional"

	email    validatorRule = "email"
	e164     validatorRule = "e.164"
	url      validatorRule = "url"
	uuid     validatorRule = "uuid"
	alpha    validatorRule = "alpha"
	alphanum validatorRule = "alphanum"
	numeric  validatorRule = "numeric"

	ipv4 validatorRule = "ipv4"
	ipv6 validatorRule = "ipv6"
)

var emailRx = regexRule{
	email,
	regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$`),
}
var uuidRx = regexRule{
	uuid,
	regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`),
}
var urlRx = regexRule{
	url,
	regexp.MustCompile(`^https?:\/\/(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(?::\d{1,5})?(?:\/[^\s]*)?$`),
}
var e164Rx = regexRule{
	e164,
	regexp.MustCompile(`^\+[1-9]\d{1,14}$`),
}

var alphaRx = regexRule{
	alpha,
	regexp.MustCompile(`^[a-zA-Z]+$`),
}

var alphanumRx = regexRule{
	alphanum,
	regexp.MustCompile(`^[a-zA-Z0-9]+$`),
}

var numericRx = regexRule{
	numeric,
	regexp.MustCompile(`^[0-9]+$`),
}

var ipv4Rx = regexRule{
	ipv4,
	regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.){3}(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)$`),
}

var ipv6Rx = regexRule{
	ipv6,
	regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`),
}

func (r *regexRule) validate(f field, v any) ValidationError {
	if f.kind != reflect.String {
		return newUserError("Cannot validate \"email\" rule against a " + f.kind.String())
	}
	if !r.regEx.MatchString(v.(string)) {
		return NewValidateError("Not a valid "+r.name, f.t.Name)
	}
	return nil
}

// the eq rules, i.e. requires a equal to sign
const (
	_     validatorRule = ``      // type this will work on | type of rule value
	min_  validatorRule = "min"   // string | collection | numeric		float64		the maximum type it "can" be, a min value can also contain min=1.52
	lte   validatorRule = "lte"   // string | collection | numeric		float64
	max_  validatorRule = "max"   // string | collection | numeric 		float64
	gte   validatorRule = "gte"   // string | collection | numeric 		float64
	len_  validatorRule = "len"   // string | collection				int
	oneof validatorRule = "oneof" // string					  			[]string
)
