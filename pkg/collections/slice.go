package collections

import (
	"fmt"

	"github.com/sivaosorg/wrapify/pkg/encoding"
)

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

// Reduce applies an accumulator function over a slice, producing a single accumulated result.
//
// This function iterates over each element in the input slice `slice`, applying the
// `accumulator` function to combine each element with an accumulated result. It starts
// with an initial value `initialValue` and successively updates the result by applying
// `accumulator` to each element in `slice`. The final accumulated result is returned
// once all elements have been processed.
//
// The function is generic, allowing it to operate on slices of any type `T` and produce
// an output of any type `U`. This enables flexible aggregation operations such as
// summing, counting, or accumulating data into more complex structures.
//
// Parameters:
//   - `slice`: The slice of elements to reduce. It can contain elements of any type `T`.
//   - `accumulator`: A function that takes the current accumulated result of type `U`
//     and an element of type `T`, then returns the updated accumulated result of type `U`.
//   - `initialValue`: The initial value for the accumulator, of type `U`. This is the
//     starting point for the reduction process.
//
// Returns:
//   - The final accumulated result of type `U` after applying `accumulator` to each element
//     in `slice`.
//
// Example:
//
//	// Summing integer values in a slice
//	numbers := []int{1, 2, 3, 4}
//	sum := Reduce(numbers, func(acc, n int) int { return acc + n }, 0)
//	// sum will be 10 as each integer is added to the accumulated result
//
//	// Concatenating strings in a slice
//	words := []string{"go", "is", "fun"}
//	sentence := Reduce(words, func(acc, word string) string { return acc + " " + word }, "")
//	// sentence will be " go is fun" (note leading space due to initial value being "")
//
//	// Using a custom struct and custom accumulator
//	type Product struct {
//	    Name  string
//	    Price float64
//	}
//	products := []Product{{Name: "apple", Price: 0.99}, {Name: "banana", Price: 1.29}}
//	totalPrice := Reduce(products, func(total float64, p Product) float64 { return total + p.Price }, 0.0)
//	// totalPrice will be 2.28 as each product's price is added to the accumulated total
func Reduce[T any, U any](slice []T, accumulator func(U, T) U, initialValue U) U {
	result := initialValue
	for _, item := range slice {
		result = accumulator(result, item)
	}
	return result
}

// IndexOf searches for a specific element in a slice and returns its index if found.
//
// This function iterates over each element in the input slice `slice` to find the first
// occurrence of the specified `item`. If `item` is found, the function returns the index
// of `item` within `slice`. If `item` is not present in the slice, it returns -1.
//
// The function is generic, allowing it to operate on slices of any comparable type `T`
// (e.g., int, string, or other types that support equality comparison).
//
// Parameters:
//   - `slice`: The slice in which to search for `item`. It can contain elements of any
//     comparable type `T`.
//   - `item`: The item to search for within `slice`. It should be of the same type `T`
//     as the elements in `slice`.
//
// Returns:
//   - The zero-based index of `item` in `slice` if it exists; otherwise, -1.
//
// Example:
//
//	// Searching for an integer in a slice
//	numbers := []int{1, 2, 3, 4}
//	index := IndexOf(numbers, 3)
//	// index will be 2, as 3 is located at index 2 in the slice
//
//	// Searching for a string in a slice
//	words := []string{"apple", "banana", "cherry"}
//	index = IndexOf(words, "banana")
//	// index will be 1, as "banana" is at index 1 in the slice
//
//	// Item not found in the slice
//	index = IndexOf(words, "date")
//	// index will be -1, as "date" is not in the slice
func IndexOf[T comparable](slice []T, item T) int {
	for i, value := range slice {
		if value == item {
			return i
		}
	}
	return -1
}

// Unique returns a new slice containing only the unique elements from the input slice,
// preserving their original order.
//
// This function iterates over each element in the input slice `slice` and uses a map
// `uniqueMap` to track elements that have already been encountered. If an element has
// not been seen before, it is added to both the `uniqueValues` result slice and the
// map. This ensures that only the first occurrence of each unique element is kept in
// the final slice, while duplicates are ignored.
//
// The function is generic, allowing it to operate on slices of any comparable type `T`.
// The elements must be of a comparable type to allow them to be used as keys in the map.
//
// Parameters:
//   - `slice`: The input slice from which unique elements are extracted. It can contain
//     elements of any comparable type `T`.
//
// Returns:
//   - A new slice of type `[]T` containing only the unique elements from `slice` in the
//     order of their first appearance.
//
// Example:
//
//	// Extracting unique integers from a slice
//	numbers := []int{1, 2, 2, 3, 4, 4, 5}
//	uniqueNumbers := Unique(numbers)
//	// uniqueNumbers will be []int{1, 2, 3, 4, 5}
//
//	// Extracting unique strings from a slice
//	words := []string{"apple", "banana", "apple", "cherry"}
//	uniqueWords := Unique(words)
//	// uniqueWords will be []string{"apple", "banana", "cherry"}
//
//	// An empty slice will return an empty result
//	empty := []int{}
//	uniqueEmpty := Unique(empty)
//	// uniqueEmpty will be []int{}
func Unique[T comparable](slice []T) []T {
	uniqueMap := make(map[T]bool)
	uniqueValues := make([]T, 0)
	for _, value := range slice {
		if _, found := uniqueMap[value]; !found {
			uniqueValues = append(uniqueValues, value)
			uniqueMap[value] = true
		}
	}
	return uniqueValues
}

// Flatten takes a slice of potentially nested elements and returns a new slice
// containing all elements of type `T` in a flat structure.
//
// This function recursively processes each element in the input slice `s`, checking if
// it is a nested slice (`[]interface{}`). If a nested slice is found, `Flatten` is called
// recursively to flatten it and append its elements to the `result` slice. If an element
// is of type `T`, it is directly appended to `result`. Elements that are neither `[]interface{}`
// nor of type `T` are ignored.
//
// The function is generic, allowing it to work with any element type `T`, which must be
// specified when calling the function. This makes `Flatten` useful for flattening slices
// with nested structures while filtering only the elements of a specified type.
//
// Parameters:
//   - `s`: A slice of `interface{}`, which can contain nested slices (`[]interface{}`) or
//     elements of any type. Nested slices may contain more nested slices at arbitrary depths.
//
// Returns:
//   - A new slice of type `[]T` containing all elements of type `T` from `s`, flattened
//     into a single level.
//
// Example:
//
//	// Flattening a nested slice of integers
//	nestedInts := []interface{}{1, []interface{}{2, 3}, []interface{}{[]interface{}{4, 5}}}
//	flatInts := Flatten[int](nestedInts)
//	// flatInts will be []int{1, 2, 3, 4, 5}
//
//	// Flattening a nested slice with mixed types, extracting only strings
//	mixedNested := []interface{}{"apple", []interface{}{"banana", 1, []interface{}{"cherry"}}}
//	flatStrings := Flatten[string](mixedNested)
//	// flatStrings will be []string{"apple", "banana", "cherry"}
//
//	// Flattening an empty slice
//	empty := []interface{}{}
//	flatEmpty := Flatten[int](empty)
//	// flatEmpty will be []int{}
func Flatten[T any](s []any) []T {
	result := make([]T, 0)
	for _, v := range s {
		switch val := v.(type) {
		case []any:
			result = append(result, Flatten[T](val)...)
		default:
			if _, ok := val.(T); ok {
				result = append(result, val.(T))
			}
		}
	}
	return result
}

// FlattenDeep takes a nested structure of arbitrary depth and returns a flat slice
// containing all elements in a single level.
//
// This function recursively processes each element in `arr`. If an element is itself a
// slice (`[]interface{}`), `FlattenDeep` calls itself to flatten that nested slice and
// appends its elements to the `result` slice. If the element is not a slice, it is directly
// added to `result`. The function allows flattening of complex nested structures while
// maintaining all elements in a single-level output.
//
// This function operates with values of type `interface{}`, making it flexible enough
// to handle mixed types in the input. It returns a slice of `interface{}`, which may
// contain elements of varying types from the original nested structure.
//
// Parameters:
//   - `arr`: The input slice, which can contain nested slices of arbitrary depth and elements
//     of any type.
//
// Returns:
//   - A slice of `[]interface{}` containing all elements from `arr` flattened into a single level.
//
// Example:
//
//	// Flattening a nested structure of mixed values
//	nested := []interface{}{1, []interface{}{2, 3, []interface{}{4, []interface{}{5}}}}
//	flat := FlattenDeep(nested)
//	// flat will be []interface{}{1, 2, 3, 4, 5}
//
//	// Flattening a deeply nested structure with varied types
//	mixedNested := []interface{}{"apple", []interface{}{"banana", 1, []interface{}{"cherry"}}}
//	flatMixed := FlattenDeep(mixedNested)
//	// flatMixed will be []interface{}{"apple", "banana", 1, "cherry"}
//
//	// Flattening a non-nested input returns the input as-is
//	nonNested := 5
//	flatNonNested := FlattenDeep(nonNested)
//	// flatNonNested will be []interface{}{5}
func FlattenDeep(arr any) []any {
	result := make([]any, 0)
	switch v := arr.(type) {
	case []any:
		for _, val := range v {
			result = append(result, FlattenDeep(val)...)
		}
	case any:
		result = append(result, v)
	}
	return result
}

// GroupBy groups elements of a slice into a map based on a specified key.
//
// This function iterates over each element in the input slice `slice`, applies the
// provided `getKey` function to extract a key for each element, and groups elements
// that share the same key into a slice. The function then returns a map where each
// key maps to a slice of elements that correspond to that key.
//
// The function is generic, allowing it to work with slices of any type `T` and to
// generate keys of any comparable type `K`. This makes `GroupBy` useful for organizing
// data based on shared attributes, such as grouping items by category or organizing
// records by a specific field.
//
// Parameters:
//   - `slice`: The input slice containing elements to be grouped. It can be of any type `T`.
//   - `getKey`: A function that takes an element of type `T` and returns a key of type `K`,
//     which is used to group the element in the resulting map.
//
// Returns:
//   - A map of type `map[K][]T`, where each key is associated with a slice of elements
//     that share that key.
//
// Example:
//
//	// Grouping integers by even and odd
//	numbers := []int{1, 2, 3, 4, 5}
//	grouped := GroupBy(numbers, func(n int) string {
//		if n%2 == 0 {
//			return "even"
//		}
//		return "odd"
//	})
//	// grouped will be map[string][]int{"even": {2, 4}, "odd": {1, 3, 5}}
//
//	// Grouping people by age
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//	people := []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}, {Name: "Charlie", Age: 30}}
//	groupedByAge := GroupBy(people, func(p Person) int { return p.Age })
//	// groupedByAge will be map[int][]Person{30: {{Name: "Alice", Age: 30}, {Name: "Charlie", Age: 30}}, 25: {{Name: "Bob", Age: 25}}}
//
//	// Grouping strings by their length
//	words := []string{"apple", "pear", "banana", "peach"}
//	groupedByLength := GroupBy(words, func(word string) int { return len(word) })
//	// groupedByLength will be map[int][]string{5: {"apple", "peach"}, 4: {"pear"}, 6: {"banana"}}
func GroupBy[T any, K comparable](slice []T, getKey func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := getKey(item)
		result[key] = append(result[key], item)
	}
	return result
}

// Join concatenates the string representation of each element in a slice into a single
// string, with a specified separator between each element.
//
// This function iterates over each element in the input slice `slice`, converts each
// element to a string using `fmt.Sprintf` with the `%v` format, and appends it to the
// `result` string. A separator string `separator` is inserted between elements in the
// final concatenated result. If the slice has only one element, no separator is added.
// The function is generic and can work with slices containing elements of any type `T`.
//
// Parameters:
//   - `slice`: The input slice containing elements to be joined. It can contain elements
//     of any type `T`.
//   - `separator`: A string that will be inserted between each element in the final result.
//
// Returns:
//   - A single string that is the result of concatenating all elements in `slice` with the
//     specified `separator` in between.
//
// Example:
//
//	// Joining integers with a comma separator
//	numbers := []int{1, 2, 3}
//	joinedNumbers := Join(numbers, ", ")
//	// joinedNumbers will be "1, 2, 3"
//
//	// Joining strings with a space separator
//	words := []string{"Go", "is", "awesome"}
//	joinedWords := Join(words, " ")
//	// joinedWords will be "Go is awesome"
//
//	// Joining an empty slice returns an empty string
//	emptySlice := []int{}
//	joinedEmpty := Join(emptySlice, ",")
//	// joinedEmpty will be ""
func Join[T any](slice []T, separator string) string {
	result := ""
	for i, item := range slice {
		if i > 0 {
			result += separator
		}
		result += fmt.Sprintf("%v", encoding.JsonSafe(item))
	}
	return result
}
