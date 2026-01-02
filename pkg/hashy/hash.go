package hashy

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/sivaosorg/wrapify/pkg/strutil"
)

// HashValue generates a 64-bit hash value for a single value with options.
// This is the primary hashing function.
//
// Parameters:
//   - value: Any Go value (struct, slice, map, primitive, etc.)
//   - options: Optional configuration (nil uses defaults)
//
// Returns:
//   - uint64: The computed hash value (never zero for valid inputs)
//   - error: Non-nil if hashing fails
//
// Example:
//
//	value := 1
//	hash, err := HashValue(value, nil)
//	fmt.Println(hash, err) // 1 nil
func HashValue(value any, options *hashOptions) (uint64, error) {
	if options == nil {
		options = DefaultOptions()
	}

	if err := options.validate(); err != nil {
		return 0, err
	}

	hasher := &hasher{
		hash:            options.Hasher,
		tagName:         options.TagName,
		treatNilAsZero:  options.ZeroNil,
		ignoreZeroValue: options.IgnoreZeroValue,
		slicesAsSets:    options.SlicesAsSets,
		useStringer:     options.UseStringer,
	}

	hasher.hash.Reset()
	return hasher.hashValue(reflect.ValueOf(value), nil)
}

// Hash generates a 64-bit hash value for the given data.
// It accepts variadic arguments - if only one argument is provided, it hashes
// that value. If multiple arguments are provided, it hashes them as a tuple.
// The last argument can optionally be *Options.
//
// Examples:
//
//	hash, err := Hash(myStruct)                    // Single value
//	hash, err := Hash(val1, val2, val3)            // Multiple values
//	hash, err := Hash(myStruct, opts)              // With options
//	hash, err := Hash(val1, val2, val3, opts)      // Multiple values with options
//
// The hash is deterministic: identical values always produce identical hashes.
//
// Returns:
//   - uint64: The computed hash value (never zero for valid inputs)
//   - error: Non-nil if hashing fails
func Hash(data ...any) (uint64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("pkg.hash: no data provided")
	}

	var opts *hashOptions
	var values []any

	if len(data) > 0 {
		if o, ok := data[len(data)-1].(*hashOptions); ok {
			opts = o
			values = data[:len(data)-1]
		} else {
			values = data
		}
	}

	if len(values) == 0 {
		return 0, fmt.Errorf("pkg.hash: no data provided")
	}

	// If single value, hash it directly
	if len(values) == 1 {
		return HashValue(values[0], opts)
	}

	return HashValue(values, opts)
}

// Hash256 generates a 256-bit hash string for the given data.
// It accepts variadic arguments - if only one argument is provided, it hashes
// that value. If multiple arguments are provided, it hashes them as a tuple.
// The last argument can optionally be *Options.
//
// Returns:
//   - string: The computed hash string (never empty for valid inputs)
//   - error: Non-nil if hashing fails
//
// Example:
//
//	hash, err := Hash256(myStruct)                    // Single value
//	hash, err := Hash256(val1, val2, val3)            // Multiple values
//	hash, err := Hash256(myStruct, opts)              // With options
//	hash, err := Hash256(val1, val2, val3, opts)      // Multiple values with options
//
// The hash is deterministic: identical values always produce identical hashes.
func Hash256(data ...any) (string, error) {
	hash, err := Hash(data...)
	if err != nil {
		return "", err
	}
	return strutil.Hash256(fmt.Sprintf("%v", hash)), nil
}

// HashBase10 generates a decimal hash string for the given data.
// It accepts variadic arguments - if only one argument is provided, it hashes
// that value. If multiple arguments are provided, it hashes them as a tuple.
// The last argument can optionally be *Options.
//
// Returns:
//   - string: The computed hash string (never empty for valid inputs)
//   - error: Non-nil if hashing fails
//
// Example:
//
//	hash, err := HashBase10(myStruct)                    // Single value
//	hash, err := HashBase10(val1, val2, val3)            // Multiple values
//	hash, err := HashBase10(myStruct, opts)              // With options
//	hash, err := HashBase10(val1, val2, val3, opts)      // Multiple values with options
//
// The hash is deterministic: identical values always produce identical hashes.
func HashBase10(data ...any) (string, error) {
	hash, err := Hash(data...)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(hash, 10), nil
}

// HashBase16 generates a hexadecimal hash string for the given data.
// It accepts variadic arguments - if only one argument is provided, it hashes
// that value. If multiple arguments are provided, it hashes them as a tuple.
// The last argument can optionally be *Options.
//
// Returns:
//   - string: The computed hash string (never empty for valid inputs)
//   - error: Non-nil if hashing fails
//
// Example:
//
//	hash, err := HashBase16(myStruct)                    // Single value
//	hash, err := HashBase16(val1, val2, val3)            // Multiple values
//	hash, err := HashBase16(myStruct, opts)              // With options
//	hash, err := HashBase16(val1, val2, val3, opts)      // Multiple values with options
//
// The hash is deterministic: identical values always produce identical hashes.
func HashBase16(data ...any) (string, error) {
	hash, err := Hash(data...)
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(hash, 16), nil
}
