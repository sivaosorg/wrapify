package hash

import (
	"hash"
	"reflect"
	"time"
)

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

// OptionsBuilder provides a fluent interface for building Options.
// It allows for chaining of method calls to configure the Options struct.
type OptionsBuilder struct {
	opts Options
}

// HashInclude is a function that can be used to check if a field should be included in the hash.
// It returns a boolean indicating whether the field should be included in the hash.
// If the function returns an error, the field will not be included in the hash.
type HashInclude func(field string, value any) (bool, error)

// HashIncludeMap is a function that can be used to check if a map field should be included in the hash.
// It returns a boolean indicating whether the map field should be included in the hash.
// If the function returns an error, the map field will not be included in the hash.
type HashIncludeMap func(field string, k, v any) (bool, error)

// Includable is an interface that can optionally be implemented by
// a struct. It will be called for each field in the struct to check whether
// it should be included in the hash.
type Includable interface {
	HashInclude() HashInclude
}

// IncludableMap is an interface that can optionally be implemented by
// a struct. It will be called for each map field in the struct to check whether
// it should be included in the hash.
type IncludableMap interface {
	HashIncludeMap() HashIncludeMap
}

// Hashable is an interface that can optionally be implemented by
// a struct. It will be called to get the hash of the struct.
// It returns a string representing the hash of the struct.
// If the function returns an error, the hash will not be included in the hash.
type Hashable interface {
	Hash() (uint64, error)
}

// ErrNotStringer is returned when there's an error with hash:"string"
//
// Parameters:
//   - Field: The name of the field that caused the error.
//
// Returns:
//   - A pointer to the `ErrNotStringer` struct.
type ErrNotStringer struct {
	Field string
}

// hasher traverses Go values and computes their hash.
//
// Parameters:
//   - hash: The hash function to use.
//   - tagName: The struct tag to look at when hashing the structure.
//   - treatNilAsZero: Determines if nil pointer should be treated equal to a zero value of pointed type.
//   - ignoreZeroValue: Determines if zero value fields should be ignored for hash calculation.
//   - slicesAsSets: Determines if slices should be treated as sets.
//   - useStringer: Determines if fmt.Stringer should be used always.
//
// Returns:
//   - A pointer to the `hasher` struct.
type hasher struct {
	hash            hash.Hash64
	tagName         string
	treatNilAsZero  bool
	ignoreZeroValue bool
	slicesAsSets    bool
	useStringer     bool
}

// visitFlag is a flag that determines the behavior of the visitor.
type visitFlag uint

// visitOptions contains context for visiting a value.
//
// Parameters:
//   - flags: The flags to use for the visit.
//   - structValue: The value of the struct being visited.
//   - fieldName: The name of the field being visited.
//
// Returns:
//   - A pointer to the `visitOptions` struct.
type visitOptions struct {
	flags       visitFlag
	structValue any
	fieldName   string
}

// timeType is the type of time.Time.
var timeType = reflect.TypeOf(time.Time{})
