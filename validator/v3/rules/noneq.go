package rules

import (
	"reflect"
	"regexp"
)

type NonEqRule struct {
	Name       string
	FieldTypes []reflect.Kind
	ErrMsg     string
	Validate   func(s string) bool
}

func (nr NonEqRule) IsNil() bool {
	return nr.Name == ""
}

func (nr NonEqRule) String() string {
	return nr.Name
}

func (nr NonEqRule) AcceptedFieldTypes() []reflect.Kind {
	return nr.FieldTypes
}

const (
	// required, optional and dive are special and have to be validated
	// by the validator on the go at runtime on each validation request
	// and cannot be expressed as a function directly

	RuleRequired = "required"
	RuleOptional = "optional"

	RuleEmail    = "email"
	RuleE164     = "e.164"
	RuleUrl      = "url"
	RuleUuid     = "uuid"
	RuleAlpha    = "alpha"
	RuleAlphanum = "alphanum"
	RuleNumeric  = "numeric"
	RuleIpv4     = "ipv4"
	RuleIpv6     = "ipv6"

	RuleDive = "dive"
)

var (
	// field must be present, and not be zero value
	Required = NonEqRule{
		RuleRequired,
		TypeAny,
		"Required field not found",
		nil,
	}
	// field may not be present, will not be validated if not present
	Optional = NonEqRule{
		RuleOptional,
		TypeAny,
		"",
		nil,
	}
	// dives into a map/slice/array to validate all fields with all the validations on the field
	Dive = NonEqRule{
		RuleDive,
		TypeCollection,
		"",
		nil,
	}
	// must be a valid email address, validated using regex
	Email = NonEqRule{
		RuleEmail,
		TypeString,
		"Invalid email",
		FnEmail,
	}
	// must be a valid e.164 string, validated using regex
	E164 = NonEqRule{
		RuleE164,
		TypeString,
		"Invalid e.164 string",
		FnE164,
	}
	// must be a valid url address, validated using regex
	Url = NonEqRule{
		RuleUrl,
		TypeString,
		"Invalid URL address",
		FnUrl,
	}
	// must be a valid UUID string, , validated using regex
	Uuid = NonEqRule{
		RuleUuid,
		TypeString,
		"Invalid UUID string",
		FnUuid,
	}
	// string must be only alphabets, validated using byte checking
	Alpha = NonEqRule{
		RuleAlpha,
		TypeString,
		"String must only contain alphabets",
		FnAlpha,
	}
	// string must be only alphabets or numbers, validated using byte checking
	Alphanum = NonEqRule{
		RuleAlphanum,
		TypeString,
		"String must only contain alphabets or numbers",
		FnAlphanum,
	}
	// string must be only numbers, validated using byte checking
	Numeric = NonEqRule{
		RuleNumeric,
		TypeString,
		"String must only contain numbers",
		FnNumeric,
	}
	// must be valid IPv4 address, validated using regex
	Ipv4 = NonEqRule{
		RuleIpv4,
		TypeString,
		"Invalid IPv4 address",
		FnIpv4,
	}
	// must be valid IPv6 address, validated using regex
	Ipv6 = NonEqRule{
		RuleIpv6,
		TypeString,
		"Invalid IPv6 address",
		FnIpv6,
	}
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

func NonEqHit(ruleName string) NonEqRule {
	switch ruleName {
	case RuleRequired:
		return Required
	case RuleOptional:
		return Optional
	case RuleEmail:
		return Email
	case RuleE164:
		return E164
	case RuleUrl:
		return Url
	case RuleUuid:
		return Uuid
	case RuleAlpha:
		return Alpha
	case RuleAlphanum:
		return Alphanum
	case RuleNumeric:
		return Numeric
	case RuleIpv4:
		return Ipv4
	case RuleIpv6:
		return Ipv6
	default:
		return NonEqRule{}
	}
}
