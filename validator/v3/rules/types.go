package rules

import (
	"reflect"
)

// ? ----- Documentation helper types -----
// These variables are of no use to the package itself but are just present for documentation purposes only

// Represents the types accepted as [TypeCollection] in rules
//   - Slice
//   - Array
//   - Map
//
// If a rule accepts only [TypeCollection] then anything else will return an user error and not be validated
var TypeCollection = [...]reflect.Kind{
	reflect.Slice,
	reflect.Array,
	reflect.Map,
}

// Represents the types accepted as [TypeNumeric] in rules
//   - Uint (all uints, such as uint8, uint16, ...)
//   - Int (all ints, such as int8, int16, ...)
//   - Float32, Float64
//
// If a rule accepts only [TypeNumeric] then anything else will return an user error and not be validated
var TypeNumeric = [...]reflect.Kind{
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Float32,
	reflect.Float64,
	reflect.Uintptr,
}
