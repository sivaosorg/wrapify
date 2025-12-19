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

// Filter returns a new slice containing only the elements from the input slice
// that satisfy a specified condition.
//
// This function iterates over each element in the input slice `list` and applies
// the provided `condition` function to it. If the `condition` function returns `true`
// for an element, that element is added to the `filtered` slice. At the end, the
// function returns the `filtered` slice containing only the elements that met the
// condition.
//
// The function is generic, allowing it to work with slices of any type `T` and
// any condition function that takes a `T` and returns a boolean.
//
// Parameters:
//   - `list`: The slice of elements to filter. It can contain elements of any type `T`.
//   - `condition`: A function that defines the filtering criteria. It takes an element
//     of type `T` as input and returns `true` if the element should be included in
//     the result, or `false` if it should be excluded.
//
// Returns:
//   - A new slice of type `[]T` containing only the elements from `list` for which
//     the `condition` function returned `true`.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	oddNumbers := Filter(numbers, func(n int) bool { return n%2 != 0 })
//	// oddNumbers will be []int{1, 3, 5} as only the odd numbers satisfy the condition
//
//	words := []string{"apple", "banana", "cherry"}
//	longWords := Filter(words, func(word string) bool { return len(word) > 5 })
//	// longWords will be []string{"banana", "cherry"} as they are longer than 5 characters
func Filter[T any](list []T, condition func(T) bool) []T {
	filtered := make([]T, 0)
	for _, item := range list {
		if condition(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// Concat returns a new slice that is the result of concatenating multiple input slices
// into a single slice.
//
// This function takes a variable number of slices as input and combines them into
// one contiguous slice. It first calculates the total length needed for the resulting
// slice, then copies each input slice into the appropriate position within the
// resulting slice.
//
// The function is generic, allowing it to concatenate slices of any type `T`.
//
// Parameters:
//   - `slices`: A variadic parameter representing the slices to concatenate. Each slice
//     can contain elements of any type `T`, and they will be concatenated in the order
//     they are provided.
//
// Returns:
//   - A new slice of type `[]T` containing all elements from each input slice in sequence.
//
// Example:
//
//	// Concatenating integer slices
//	a := []int{1, 2}
//	b := []int{3, 4}
//	c := []int{5, 6}
//	combined := Concat(a, b, c)
//	// combined will be []int{1, 2, 3, 4, 5, 6}
//
//	// Concatenating string slices
//	words1 := []string{"hello", "world"}
//	words2 := []string{"go", "lang"}
//	concatenatedWords := Concat(words1, words2)
//	// concatenatedWords will be []string{"hello", "world", "go", "lang"}
func Concat[T any](slices ...[]T) []T {
	totalLen := 0
	for _, s := range slices {
		totalLen += len(s)
	}
	result := make([]T, totalLen)
	i := 0
	for _, s := range slices {
		copy(result[i:], s)
		i += len(s)
	}
	return result
}

// Sum calculates the sum of elements in a slice after transforming each element to a float64.
//
// This function iterates over each element in the input slice `slice`, applies a transformation
// function `transformer` to convert the element to a float64, and adds the result to a running
// total. The final sum is returned as a float64.
//
// The function is generic, allowing it to operate on slices of any type `T`. The `transformer`
// function is used to convert each element to a float64, enabling flexible summation of
// different types (e.g., integers, custom types with numeric properties).
//
// Parameters:
//   - `slice`: The slice of elements to sum. It can contain elements of any type `T`.
//   - `transformer`: A function that takes an element of type `T` and returns a float64
//     representation, which will be used in the summation.
//
// Returns:
//   - A float64 representing the sum of the transformed elements.
//
// Example:
//
//	// Summing integer slice values
//	numbers := []int{1, 2, 3, 4}
//	total := Sum(numbers, func(n int) float64 { return float64(n) })
//	// total will be 10.0 as each integer is converted to float64 and summed
//
//	// Summing custom struct values
//	type Product struct {
//	    Price float64
//	}
//	products := []Product{{Price: 9.99}, {Price: 19.99}, {Price: 5.0}}
//	totalPrice := Sum(products, func(p Product) float64 { return p.Price })
//	// totalPrice will be 34.98 as the prices are summed
func Sum[T any](slice []T, transformer func(T) float64) float64 {
	sum := 0.0
	for _, item := range slice {
		sum += transformer(item)
	}
	return sum
}

// Equal checks if two slices are equal in both length and elements.
//
// This function compares two slices `a` and `b` of any comparable type `T`. It first
// checks if the lengths of the two slices are the same. If they are not, it returns `false`.
// If the lengths match, it then iterates through each element in `a` and `b` to check
// if corresponding elements are equal. If all elements are equal, the function returns `true`;
// otherwise, it returns `false`.
//
// The function is generic and can be used with slices of any comparable type, such as
// integers, strings, or other types that support equality comparison.
//
// Parameters:
//   - `a`: The first slice to compare. It should contain elements of a comparable type `T`.
//   - `b`: The second slice to compare. It should also contain elements of type `T`.
//
// Returns:
//   - `true` if both slices have the same length and identical elements at each position;
//     `false` otherwise.
//
// Example:
//
//	// Comparing integer slices
//	a := []int{1, 2, 3}
//	b := []int{1, 2, 3}
//	isEqual := Equal(a, b)
//	// isEqual will be true as both slices contain the same elements in the same order
//
//	c := []int{1, 2, 4}
//	isEqual = Equal(a, c)
//	// isEqual will be false as the elements differ
//
//	// Comparing string slices
//	names1 := []string{"Alice", "Bob"}
//	names2 := []string{"Alice", "Bob"}
//	isEqual = Equal(names1, names2)
//	// isEqual will be true since the slices have identical elements
func Equal[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
