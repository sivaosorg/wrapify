package hash

import "reflect"

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
func HashValue(value any, options *Options) (uint64, error) {
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
