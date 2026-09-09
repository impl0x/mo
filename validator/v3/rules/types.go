package rules

import "go/types"

// ? ----- Documentation helper types -----
// These variables are of no use to the package itself but are just present for documentation purposes only

// Represents the types accepted as [TypeCollection] in rules
//   - Slice
//   - Array
//   - Map
//
// If a rule accepts only [TypeCollection] then anything else will return an user error and not be validated
var TypeCollection = [...]types.Type{
	&types.Slice{},
	&types.Array{},
	&types.Map{},
}

// Represents the types accepted as [TypeNumeric] in rules
//   - Uint (all uints, such as uint8, uint16, ...)
//   - Int (all ints, such as int8, int16, ...)
//   - Float32, Float64
//
// If a rule accepts only [TypeNumeric] then anything else will return an user error and not be validated
var TypeNumeric = [...]types.BasicKind{
	types.Int,
	types.Int8,
	types.Int16,
	types.Int32,
	types.Int64,
	types.Uint,
	types.Uint8,
	types.Uint16,
	types.Uint32,
	types.Uint64,
	types.Float32,
	types.Float64,
	types.Uintptr,
	types.UntypedInt,
}
