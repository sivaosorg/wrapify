package hash

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"reflect"
	"time"
)

// WithHasher sets the hash function to use.
//
// Parameters:
//   - h: The hash function to use.
//
// Returns:
//   - A pointer to the `OptionsBuilder` struct.
//
// Example:
//
//	builder := NewOptions().WithHasher(fnv.New64a())
//	opts := builder.Build()
func (b *OptionsBuilder) WithHasher(h hash.Hash64) *OptionsBuilder {
	b.opts.Hasher = h
	return b
}

// WithTagName sets the struct tag to look at when hashing the structure.
//
// Parameters:
//   - name: The name of the struct tag to look at.
//
// Returns:
//   - A pointer to the `OptionsBuilder` struct.
//
// Example:
//
//	builder := NewOptions().WithTagName("json")
//	opts := builder.Build()
func (b *OptionsBuilder) WithTagName(name string) *OptionsBuilder {
	b.opts.TagName = name
	return b
}

// WithZeroNil sets whether nil pointer should be treated equal to a zero value of pointed type.
//
// Parameters:
//   - zeroNil: A boolean indicating whether nil pointer should be treated equal to a zero value of pointed type.
//
// Returns:
//   - A pointer to the `OptionsBuilder` struct.
//
// Example:
//
//	builder := NewOptions().WithZeroNil(true)
//	opts := builder.Build()
func (b *OptionsBuilder) WithZeroNil(zeroNil bool) *OptionsBuilder {
	b.opts.ZeroNil = zeroNil
	return b
}

// WithIgnoreZeroValue sets whether zero value fields should be ignored for hash calculation.
//
// Parameters:
//   - ignore: A boolean indicating whether zero value fields should be ignored for hash calculation.
//
// Returns:
//   - A pointer to the `OptionsBuilder` struct.
//
// Example:
//
//	builder := NewOptions().WithIgnoreZeroValue(true)
//	opts := builder.Build()
func (b *OptionsBuilder) WithIgnoreZeroValue(ignore bool) *OptionsBuilder {
	b.opts.IgnoreZeroValue = ignore
	return b
}

// WithSlicesAsSets sets whether slices should be treated as sets.
//
// Parameters:
//   - asSets: A boolean indicating whether slices should be treated as sets.
//
// Returns:
//   - A pointer to the `OptionsBuilder` struct.
//
// Example:
//
//	builder := NewOptions().WithSlicesAsSets(true)
//	opts := builder.Build()
func (b *OptionsBuilder) WithSlicesAsSets(asSets bool) *OptionsBuilder {
	b.opts.SlicesAsSets = asSets
	return b
}

// WithUseStringer sets whether fmt.Stringer should be used always.
//
// Parameters:
//   - useStringer: A boolean indicating whether fmt.Stringer should be used always.
//
// Returns:
//   - A pointer to the `OptionsBuilder` struct.
//
// Example:
//
//	builder := NewOptions().WithUseStringer(true)
//	opts := builder.Build()
func (b *OptionsBuilder) WithUseStringer(useStringer bool) *OptionsBuilder {
	b.opts.UseStringer = useStringer
	return b
}

// Build builds the options.
//
// Returns:
//   - A pointer to the `Options` struct.
//
// Example:
//
//	builder := NewOptions().Build()
//	opts := builder.Build()
func (b *OptionsBuilder) Build() *Options {
	return &b.opts
}

// unwrapValue unwraps the value of a field.
//
// Parameters:
//   - value: The value of the field to unwrap.
//
// Returns:
//   - The unwrapped value.
//
// Example:
//
//	value := reflect.ValueOf(1)
//	unwrappedValue := h.unwrapValue(value)
//	fmt.Println(unwrappedValue) // 1
func (h *hasher) unwrapValue(value reflect.Value) reflect.Value {
	targetType := reflect.TypeOf(0)

	for {
		if value.Kind() == reflect.Interface {
			value = value.Elem()
			continue
		}

		if value.Kind() == reflect.Ptr {
			if h.treatNilAsZero {
				targetType = value.Type().Elem()
			}
			value = reflect.Indirect(value)
			continue
		}

		break
	}

	// If it is nil, treat it like a zero value
	if !value.IsValid() {
		value = reflect.Zero(targetType)
	}

	return value
}

// normalizeValue normalizes the value of a field.
//
// Parameters:
//   - value: The value of the field to normalize.
//
// Returns:
//   - The normalized value.
//
// Example:
//
//	value := reflect.ValueOf(1)
//	normalizedValue := h.normalizeValue(value)
//	fmt.Println(normalizedValue) // 1
func (h *hasher) normalizeValue(value reflect.Value) reflect.Value {
	switch value.Kind() {
	case reflect.Int:
		return reflect.ValueOf(int64(value.Int()))
	case reflect.Uint:
		return reflect.ValueOf(uint64(value.Uint()))
	case reflect.Bool:
		var val int8
		if value.Bool() {
			val = 1
		}
		return reflect.ValueOf(val)
	}
	return value
}

// tryHashable tries to get the hash of a value.
//
// Parameters:
//   - value: The value to try to get the hash of.
//
// Returns:
//   - The hash of the value, a boolean indicating if the value is hashable, and an error if there is an error.
//
// Example:
//
//	value := reflect.ValueOf(1)
//	hash, ok, err := h.tryHashable(value)
//	fmt.Println(hash, ok, err) // 1 true nil
func (h *hasher) tryHashable(value reflect.Value) (uint64, bool, error) {
	// Try direct implementation
	if hashable, ok := value.Interface().(Hashable); ok {
		hash, err := hashable.Hash()
		return hash, true, err
	}

	// Try pointer implementation
	if value.CanAddr() {
		ptr := value.Addr()
		if hashable, ok := ptr.Interface().(Hashable); ok {
			hash, err := hashable.Hash()
			return hash, true, err
		}
	}

	return 0, false, errors.New("pkg.hash: value is not hashable")
}

// hashZeroValue hashes a zero value.
//
// Returns:
//   - The hash of the zero value.
//
// Example:
//
//	hash, err := h.hashZeroValue()
//	fmt.Println(hash, err) // 0 nil
func (h *hasher) hashZeroValue() (uint64, error) {
	h.hash.Reset()
	if err := binary.Write(h.hash, binary.LittleEndian, int64(0)); err != nil {
		return 0, err
	}
	return h.hash.Sum64(), nil
}

// hashNumeric hashes a numeric value.
//
// Parameters:
//   - value: The value to hash.
//
// Returns:
//   - The hash of the value.
//
// Example:
//
//	value := reflect.ValueOf(1)
//	hash, err := h.hashNumeric(value)
//	fmt.Println(hash, err) // 1 nil
func (h *hasher) hashNumeric(value reflect.Value) (uint64, error) {
	h.hash.Reset()
	if err := binary.Write(h.hash, binary.LittleEndian, value.Interface()); err != nil {
		return 0, err
	}
	return h.hash.Sum64(), nil
}

// hashTime hashes a time value.
//
// Parameters:
//   - value: The value to hash.
//
// Returns:
//   - The hash of the value.
//
// Example:
//
//	value := reflect.ValueOf(time.Now())
//	hash, err := h.hashTime(value)
//	fmt.Println(hash, err) // time.Now().Unix() nil
func (h *hasher) hashTime(value reflect.Value) (uint64, error) {
	h.hash.Reset()
	timeVal := value.Interface().(time.Time)
	data, err := timeVal.MarshalBinary()
	if err != nil {
		return 0, err
	}
	if err := binary.Write(h.hash, binary.LittleEndian, data); err != nil {
		return 0, err
	}
	return h.hash.Sum64(), nil
}

// hashString hashes a string value.
//
// Parameters:
//   - value: The value to hash.
//
// Returns:
//   - The hash of the value.
//
// Example:
//
//	value := reflect.ValueOf("hello")
//	hash, err := h.hashString(value)
//	fmt.Println(hash, err) // "hello" nil
func (h *hasher) hashString(value reflect.Value) (uint64, error) {
	h.hash.Reset()
	if _, err := h.hash.Write([]byte(value.String())); err != nil {
		return 0, err
	}
	return h.hash.Sum64(), nil
}

// hashUpdateOrdered hashes two values in order.
//
// Parameters:
//   - a: The first value to hash.
//   - b: The second value to hash.
//
// Returns:
//   - The hash of the two values.
//
// Example:
//
//	a := 1
//	b := 2
//	hash, err := h.hashUpdateOrdered(a, b)
//	fmt.Println(hash, err) // 3 nil
func (h *hasher) hashUpdateOrdered(a, b uint64) uint64 {
	h.hash.Reset()
	if err := binary.Write(h.hash, binary.LittleEndian, a); err != nil {
		panic(fmt.Sprintf("hash write failed: %v", err))
	}
	if err := binary.Write(h.hash, binary.LittleEndian, b); err != nil {
		panic(fmt.Sprintf("hash write failed: %v", err))
	}
	return h.hash.Sum64()
}

// hashFinishUnordered "hardens" the XOR result to prevent cancellation issues.
// After mixing a group of unique hashes with hashUpdateUnordered, it's necessary
// to call hashFinishUnordered. This prevents issues where XOR operations can
// cancel out when the same hash appears in different contexts.
//
// Parameters:
//   - a: The value to hash.
//
// Returns:
//   - The hash of the value.
//
// Example:
//
//	a := 1
//	hash, err := h.hashFinishUnordered(a)
//	fmt.Println(hash, err) // 1 nil
func (h *hasher) hashFinishUnordered(a uint64) uint64 {
	h.hash.Reset()
	if err := binary.Write(h.hash, binary.LittleEndian, a); err != nil {
		panic(fmt.Sprintf("hash write failed: %v", err))
	}
	return h.hash.Sum64()
}

// hashUpdateUnordered hashes two values in unordered.
//
// Parameters:
//   - a: The first value to hash.
//   - b: The second value to hash.
//
// Returns:
//   - The hash of the two values.
//
// Example:
//
//	a := 1
//	b := 2
//	hash, err := h.hashUpdateUnordered(a, b)
//	fmt.Println(hash, err) // 3 nil
func hashUpdateUnordered(a, b uint64) uint64 {
	return a ^ b
}
