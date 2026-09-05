package validator

import (
	"fmt"
	"reflect"
	"regexp"
)

// ? ----- Utils -----
func runeIsUpper(r rune) (b bool) {
	if r >= 65 && r <= 90 {
		b = true
	}
	return
}

func runeIsLower(r rune) (b bool) {
	if r >= 97 && r <= 122 {
		b = true
	}
	return
}

func runeIsNum(r rune) (b bool) {
	if r >= 48 && r <= 57 {
		b = true
	}
	return
}

// applies the func fn to every character in the string and returns true on the first true returned from fn
func loopStr(s string, fn func(c rune) bool) bool {
	var b bool
	for _, c := range s {
		if b = fn(c); b {
			return b
		}
	}
	return b
}

type strValidator interface {
	Validate(s string) bool
}

type regexVal struct {
	regEx *regexp.Regexp
}

func (rv regexVal) Validate(s string) bool {
	return rv.regEx.MatchString(s)
}

type manVal struct {
	fn func(s string) bool
}

func (mv manVal) Validate(s string) bool {
	return mv.fn(s)
}

// ? ----- Rules -----

type nonEqRule = string

const ( // these are all nonEq rules
	ruleRequired nonEqRule = "required"
	ruleOptional nonEqRule = "optional"

	ruleEmail    nonEqRule = "email"
	ruleE164     nonEqRule = "e.164"
	ruleUrl      nonEqRule = "url"
	ruleUuid     nonEqRule = "uuid"
	ruleAlpha    nonEqRule = "alpha"
	ruleAlphanum nonEqRule = "alphanum"
	ruleNumeric  nonEqRule = "numeric"

	ruleIpv4 nonEqRule = "ipv4"
	ruleIpv6 nonEqRule = "ipv6"

	ruleDive nonEqRule = "dive"
)

// all the string validation regexes/manual validations
var nonEqStrRules = map[nonEqRule]strValidator{
	ruleEmail: regexVal{regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$`)},
	ruleE164:  regexVal{regexp.MustCompile(`^\+[1-9]\d{1,14}$`)},
	ruleUrl:   regexVal{regexp.MustCompile(`^https?:\/\/(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(?::\d{1,5})?(?:\/[^\s]*)?$`)},
	ruleUuid:  regexVal{regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)},
	ruleAlpha: manVal{func(s string) bool {
		return loopStr(s, func(c rune) bool {
			if !(runeIsLower(c) || runeIsUpper(c)) { // if not both lower or upper
				return false
			}
			return true
		})
	}},
	ruleAlphanum: manVal{func(s string) bool {
		return loopStr(s, func(c rune) bool {
			if !(runeIsLower(c) || runeIsUpper(c) || runeIsNum(c)) { // if not all lower, upper or num
				return false
			}
			return true
		})
	}},
	ruleNumeric: manVal{func(s string) bool {
		return loopStr(s, func(c rune) bool {
			if !runeIsUpper(c) {
				return false
			}
			return true
		})
	}},
	ruleIpv4: regexVal{regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.){3}(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)$`)},
	ruleIpv6: regexVal{regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)},
}

func (vd *validator) validateNonEqRuleStr(ruleName nonEqRule, svd strValidator) ValidationError {
	if vd.f.kind != reflect.String {
		return newUserError(fmt.Sprintf("cannot validate \"%s\" rule against a %s", ruleName, vd.f.kind.String()), vd.parent, vd.f.fieldName)
	}
	if !svd.Validate(vd.f.v.String()) {
		return newFieldValidateError("Not a valid "+ruleName, "", vd.parent, vd.f)
	}
	return nil
}

type EqRule = string

// the eq rules, i.e. requires a equal to sign
const (
	_          EqRule = ``           // type this will work on 			 	type of rule value
	min_       EqRule = "min"        // string | collection | numeric		float64		? the largest type it "can" be, a min value can also contain min=1.52
	lte        EqRule = "lte"        // string | collection | numeric		float64
	lt         EqRule = "lt"         // string | collection | numeric 		float64
	max_       EqRule = "max"        // string | collection | numeric 		float64
	gte        EqRule = "gte"        // string | collection | numeric 		float64
	gt         EqRule = "gt"         // string | collection | numeric 		float64
	len_       EqRule = "len"        // string | collection					int
	oneof      EqRule = "oneof"      // string					  			[]string
	startswith EqRule = "startswith" // string								string
	endswith   EqRule = "endswith"   // string 								string
)
