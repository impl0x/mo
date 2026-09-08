package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strings"
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

// ? ----- Non Eq Rules -----
type nonEqStrValidator interface {
	Validate(s string) bool
}

type regexValidator struct {
	regEx *regexp.Regexp
}

func (rv regexValidator) Validate(s string) bool {
	return rv.regEx.MatchString(s)
}

type manualValidator struct {
	fn func(s string) bool
}

func (mv manualValidator) Validate(s string) bool {
	return mv.fn(s)
}

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

// all the non eq string validation regexes/manual validations
var nonEqStrRules = map[nonEqRule]nonEqStrValidator{
	ruleEmail: regexValidator{regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$`)},
	ruleE164:  regexValidator{regexp.MustCompile(`^\+[1-9]\d{1,14}$`)},
	ruleUrl:   regexValidator{regexp.MustCompile(`^https?:\/\/(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(?::\d{1,5})?(?:\/[^\s]*)?$`)},
	ruleUuid:  regexValidator{regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)},
	ruleAlpha: manualValidator{func(s string) bool {
		return loopStr(s, func(c rune) bool {
			if !(runeIsLower(c) || runeIsUpper(c)) { // if not both lower or upper
				return false
			}
			return true
		})
	}},
	ruleAlphanum: manualValidator{func(s string) bool {
		return loopStr(s, func(c rune) bool {
			if !(runeIsLower(c) || runeIsUpper(c) || runeIsNum(c)) { // if not all lower, upper or num
				return false
			}
			return true
		})
	}},
	ruleNumeric: manualValidator{func(s string) bool {
		return loopStr(s, func(c rune) bool {
			if !runeIsUpper(c) {
				return false
			}
			return true
		})
	}},
	ruleIpv4: regexValidator{regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.){3}(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)$`)},
	ruleIpv6: regexValidator{regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`)},
}

func (vd *validator) validateNonEqRuleStr(ruleName nonEqRule, svd nonEqStrValidator) ValidationError {
	if vd.f.kind != reflect.String {
		return newUserError(fmt.Sprintf("cannot validate \"%s\" rule against a %s", ruleName, vd.f.kind.String()), vd.parent, vd.f.fieldName)
	}
	if !svd.Validate(vd.f.v.String()) {
		return newFieldValidateError("Not a valid "+ruleName, "", vd.parent, vd.f)
	}
	return nil
}

// ? ----- Eq Rules -----
// ! Eq rules mean the tag validation rules which have an additional parameter to them after an equals to symbol,
// for example: "min=2,max=15"

type EqRule = string

// the eq rules, i.e. requires a equal to sign
const (
	_              EqRule = ``           // type this will work on 			 	type of param value
	ruleMin        EqRule = "min"        // string | collection | numeric		float64		? the largest type it "can" be, a min value can also contain min=1.52
	ruleMax        EqRule = "max"        // string | collection | numeric 		float64
	ruleLte        EqRule = "lte"        // string | collection | numeric		float64
	ruleGte        EqRule = "gte"        // string | collection | numeric 		float64
	ruleLt         EqRule = "lt"         // string | collection | numeric 		float64
	ruleGt         EqRule = "gt"         // string | collection | numeric 		float64
	ruleLen        EqRule = "len"        // string | collection					int
	ruleStartswith EqRule = "startswith" // string								string
	ruleEndswith   EqRule = "endswith"   // string 								string
	ruleOneof      EqRule = "oneof"      // string					  			[]string
)

// ! eq rule funcs
// these funcs are just one line bool returns but i wrote this to have them rule logic separated from reflection and validator logic
// value is the field value in the struct and param is the tag parameter
// example: in a field username string `validate:"min=2", where username is populated with "test" in a instance,
// value will be len("test") = 4 and param will be 2 from the "min=2" tag.
var (
	eqRuleMin = func(value float64, param float64) bool {
		return value > param
	}
	eqRuleMax = func(value float64, param float64) bool {
		return value < param
	}
	eqRuleGte = eqRuleMin // logically the same function
	eqRuleLte = eqRuleMax // ~
	eqRuleLt  = func(value float64, param float64) bool {
		return value <= param
	}
	eqRuleGt = func(value float64, param float64) bool {
		return value >= param
	}
	eqRuleLen = func(value int, param int) bool {
		return value == param
	}
	eqRuleStartswith = func(value string, param string) bool {
		return strings.HasPrefix(value, param)
	}
	eqRuleEndswith = func(value string, param string) bool {
		return strings.HasSuffix(value, param)
	}
	eqRuleOneof = func(value string, param []string) bool {
		return slices.Contains(param, value)
	}
)
