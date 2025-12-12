package collections

// Contains checks if a specified item is present within a given slice.
//
// This function iterates over a slice of any type that supports comparison
// and checks if the specified `item` exists within it. It returns `true`
// if the `item` is found and `false` otherwise.
//
// The function is generic, so it can be used with any comparable type,
// including strings, integers, floats, or custom types that implement the
// comparable interface.
//
// Parameters:
//   - `array`: The slice of elements to search through. This slice can
//     contain any type `T` that supports comparison (e.g., int, string).
//   - `item`: The item to search for within `array`. It should be of the
//     same type `T` as the elements in `array`.
//
// Returns:
//   - `true` if `item` is found within `array`, `false` otherwise.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	isPresent := Contains(numbers, 3) // isPresent will be true as 3 is in the slice
//
//	names := []string{"Alice", "Bob", "Charlie"}
//	isPresent := Contains(names, "Eve") // isPresent will be false as "Eve" is not in the slice
func Contains[T comparable](array []T, item T) bool {
	if len(array) == 0 {
		return false
	}
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}
