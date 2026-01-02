package hash

import "hash"

// Options contains configuration for hash generation.
type Options struct {
	// Hasher is the hash function to use. If this isn't set, it will
	// default to FNV-1a.
	Hasher hash.Hash64

	// TagName is the struct tag to look at when hashing the structure.
	// By default this is "hash".
	TagName string

	// ZeroNil determines if nil pointer should be treated equal
	// to a zero value of pointed type. By default this is false.
	ZeroNil bool

	// IgnoreZeroValue determines if zero value fields should be
	// ignored for hash calculation.
	IgnoreZeroValue bool

	// SlicesAsSets assumes that a `set` tag is always present for slices.
	// Default is false (in which case the tag is used instead)
	SlicesAsSets bool

	// UseStringer will attempt to use fmt.Stringer always. If the struct
	// doesn't implement fmt.Stringer, it'll fall back to trying usual tricks.
	// If this is true, and the "string" tag is also set, the tag takes
	// precedence (meaning that if the type doesn't implement fmt.Stringer, we
	// return an error)
	UseStringer bool
}
