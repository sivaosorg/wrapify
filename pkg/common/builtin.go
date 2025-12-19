package common

import "reflect"

// DeepEqual compares two values of any comparable type to determine if they are deeply equal.
//
// This function uses the `reflect.DeepEqual` function from the `reflect` package to compare
// two values `a` and `b`. It checks for deep equality, meaning it considers nested structures,
// such as slices, maps, or structs, and compares them element-by-element or field-by-field.
// If the values are deeply equal, the function returns `true`; otherwise, it returns `false`.
//
// The function is generic, allowing it to work with any type `T` that is comparable, including
// basic types (e.g., integers, strings) as well as complex types with nested structures.
//
// Parameters:
//   - `a`: The first value to compare. It can be of any comparable type `T`.
//   - `b`: The second value to compare. It must be of the same type `T` as `a`.
//
// Returns:
//   - `true` if `a` and `b` are deeply equal; `false` otherwise.
//
// Example:
//
//	// Comparing two integer values
//	isEqual := DeepEqual(5, 5)
//	// isEqual will be true as both integers are equal
//
//	// Comparing two slices with the same elements
//	sliceA := []int{1, 2, 3}
//	sliceB := []int{1, 2, 3}
//	isEqual = DeepEqual(sliceA, sliceB)
//	// isEqual will be true as both slices have identical elements in the same order
//
//	// Comparing two different maps
//	mapA := map[string]int{"a": 1, "b": 2}
//	mapB := map[string]int{"a": 1, "b": 3}
//	isEqual = DeepEqual(mapA, mapB)
//	// isEqual will be false as the values for key "b" differ between the maps
func DeepEqual[T comparable](a, b T) bool {
	return reflect.DeepEqual(a, b)
}
