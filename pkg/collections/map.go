package collections

// MapContainsKey checks if a specified key is present within a given map.
//
// This function takes a map with keys of any comparable type `K` and values of
// any type `V`. It checks if the specified `key` exists in the map `m`. If the key
// is found, it returns `true`; otherwise, it returns `false`.
//
// The function is generic and can be used with maps that have keys of any type
// that supports comparison (e.g., int, string). The value type `V` can be any type.
//
// Parameters:
//   - `m`: The map in which to search for the key. The map has keys of type `K`
//     and values of type `V`.
//   - `key`: The key to search for within `m`. It should be of the same type `K` as
//     the keys in `m`.
//
// Returns:
//   - `true` if `key` is found in `m`, `false` otherwise.
//
// Example:
//
//	ages := map[string]int{"Alice": 30, "Bob": 25}
//	isPresent := MapContainsKey(ages, "Alice") // isPresent will be true as "Alice" is a key in the map
//
//	prices := map[int]float64{1: 9.99, 2: 19.99}
//	isPresent := MapContainsKey(prices, 3) // isPresent will be false as 3 is not a key in the map
func MapContainsKey[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]
	return ok
}

// Map returns a new slice where each element is the result of applying a specified
// transformation function to each element in the input slice.
//
// This function iterates over each element in the input slice `list`, applies the
// provided transformation function `f` to it, and stores the result in the new slice
// `result`. The length of the resulting slice is the same as the input slice, and
// each element in `result` corresponds to a transformed element from `list`.
//
// The function is generic, allowing it to work with slices of any type `T` and
// apply a transformation function that converts each element of type `T` to a
// new type `U`.
//
// Parameters:
//   - `list`: The slice of elements to transform. It can contain elements of any type `T`.
//   - `f`: A function that defines the transformation. It takes an element of type `T`
//     as input and returns a transformed value of type `U`.
//
// Returns:
//   - A new slice of type `[]U` where each element is the result of applying `f`
//     to the corresponding element in `list`.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	squaredNumbers := Map(numbers, func(n int) int { return n * n })
//	// squaredNumbers will be []int{1, 4, 9, 16, 25} as each number is squared
//
//	words := []string{"apple", "banana", "cherry"}
//	wordLengths := Map(words, func(word string) int { return len(word) })
//	// wordLengths will be []int{5, 6, 6} as each word's length is calculated
func Map[T any, U any](list []T, f func(T) U) []U {
	result := make([]U, len(list))
	for i, item := range list {
		result[i] = f(item)
	}
	return result
}
