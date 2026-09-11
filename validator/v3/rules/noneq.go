package rules

import "regexp"

// ? ----- Non Eq Rules -----
// ! Non Eq rules mean the tag validation rules which have only one word,
// for example: "required,email"

// these are all nonEq rules
type NonEqRule = string

const (
	// required, optional and dive are special and have to be validated
	// by the validator on the go at runtime on each validation request
	// and cannot be expressed as a function directly

	Required NonEqRule = "required" // INFO: field must be present and not have its zero value. TYPE: any
	Optional NonEqRule = "optional" // INFO: skips validation if empty. TYPE: any

	Email    NonEqRule = "email"    // INFO: must satisfy email format. TYPE: string
	E164     NonEqRule = "e.164"    // INFO: must satisfy phone number format. TYPE: string
	Url      NonEqRule = "url"      // INFO: must satisfy url format. TYPE: string
	Uuid     NonEqRule = "uuid"     // INFO: must satisfy uuid format. TYPE: string
	Alpha    NonEqRule = "alpha"    // INFO: must be only alphabets. TYPE: string
	Alphanum NonEqRule = "alphanum" // INFO: must be only alphabets or numbers. TYPE: string
	Numeric  NonEqRule = "numeric"  // INFO: must be only a number. TYPE: string
	Ipv4     NonEqRule = "ipv4"     // INFO: must satisfy ipv4 format. TYPE: string
	Ipv6     NonEqRule = "ipv6"     // INFO: must satisfy ipv6 format. TYPE: string

	Dive NonEqRule = "dive" // INFO: dives into a slice/array and validates all other rules. TYPE: slice | array
)

// ? ----- Validation functions -----
// all the non eq string validation regexes/manual validations
var (
	FnEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9]*[a-zA-Z][a-zA-Z0-9]*$`).MatchString
	FnE164  = regexp.MustCompile(`^\+[1-9]\d{1,14}$`).MatchString
	FnUrl   = regexp.MustCompile(`^https?:\/\/(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(?::\d{1,5})?(?:\/[^\s]*)?$`).MatchString
	FnUuid  = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`).MatchString

	// Logic for the functions below are basically looping over
	// the string using loopStr util func and passing in
	// the condition which will return on the first true,
	// so we pass in the not of our condition so if even one rune fails
	// the loop stops and returns true to which we put a not over the result again
	// and return the opposite as we need to fail if the fail condition satisfied.

	FnAlpha = func(s string) bool {
		return !loopStr(s, func(c rune) bool {
			return !(runeIsLower(c) || runeIsUpper(c)) // if not both lower or upper
		})
	}
	FnAlphanum = func(s string) bool {
		return !loopStr(s, func(c rune) bool {
			return !(runeIsLower(c) || runeIsUpper(c) || runeIsNum(c)) // if not lower, upper or number
		})
	}
	FnNumeric = func(s string) bool {
		return !loopStr(s, func(c rune) bool {
			return !runeIsNum(c) // if not number
		})
	}
	FnIpv4 = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.){3}(25[0-5]|(2[0-4]|1\d|[1-9]|)\d)$`).MatchString
	FnIpv6 = regexp.MustCompile(`^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`).MatchString
)

// ? ----- Utils -----
func runeIsUpper(r rune) bool {
	return r >= 65 && r <= 90
}

func runeIsLower(r rune) bool {
	return r >= 97 && r <= 122
}

func runeIsNum(r rune) bool {
	return r >= 48 && r <= 57
}

// Loops over a string and calls fn with each rune and returns on first true
func loopStr(s string, fn func(c rune) bool) bool {
	var b bool
	for _, c := range s {
		if b = fn(c); b {
			return b
		}
	}
	return b
}

// ? ----- Switch Func-----

// Returns the specific validator function for the rule. Panics on invalid rule
//
// Only works on rules which take a string for input and return bool depending on the validation
func NonEqRuleToFunc(nr NonEqRule) func(s string) bool {
	switch nr {
	case Email:
		return FnEmail
	case E164:
		return FnE164
	case Url:
		return FnUrl
	case Uuid:
		return FnUuid
	case Alpha:
		return FnAlpha
	case Alphanum:
		return FnAlphanum
	case Numeric:
		return FnNumeric
	case Ipv4:
		return FnIpv4
	case Ipv6:
		return FnIpv6
	default:
		panic("rules.nonEqRule.Validator: invalid nonEqRule passed, function does not exist for " + nr)
	}
}
