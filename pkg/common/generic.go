package common

import "reflect"

// PredicateTransform applies a transformation function to each element of a collection (slice, array, or map) and returns
// a new collection with the transformed values.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice, array, or map,
// and a mapping function (`mapper`). The function applies the `mapper` function to each value in the collection,
// transforming it. For slices and arrays, the callback is applied to each element. For maps, the callback is applied
// to each key and value. The transformed elements (or key-value pairs) are collected into a new result, which is returned.
//
// The function uses reflection to handle different types of collections and constructs a new collection with the transformed values.
//
// Parameters:
//   - `collection`: The collection (slice, array, or map) to iterate over and transform. It can be of any type.
//   - `mapper`: A function that takes a value from the collection and transforms it. The function is applied to each
//     element of the collection (or key-value pair for maps).
//
// Returns:
//   - A new collection of the same type as the original collection, where each element has been transformed using the
//     provided `mapper` function.
//
// Example:
//
//	// Mapping a slice of integers to their squares
//	numbers := []int{1, 2, 3, 4}
//	squared := PredicateTransform(numbers, func(value interface{}) interface{} {
//		return value.(int) * value.(int)
//	})
//	// squared will be []int{1, 4, 9, 16}
//
//	// Mapping a map of strings to their lengths
//	words := map[string]string{"apple": "fruit", "carrot": "vegetable"}
//	lengths := PredicateTransform(words, func(value interface{}) interface{} {
//		return len(value.(string))
//	})
//	// lengths will be a new collection containing key-value pairs, where values represent string lengths.
//
// Notes:
//   - For slices and arrays, the `mapper` function is applied to each element.
//   - For maps, the `mapper` function is applied to both the key and the value. The resulting collection will include
//     both transformed keys and values.
//
// Limitations:
//   - The function creates a new collection based on the results of the `mapper` function, so it does not modify the
//     original collection.
func PredicateTransform(collection any, mapper func(value any) any) any {
	v := reflect.ValueOf(collection)
	result := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(mapper(v.Index(0).Interface()))), 0, 0)

	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			mappedValue := mapper(v.Index(i).Interface())
			result = reflect.Append(result, reflect.ValueOf(mappedValue))
		}
	} else if v.Kind() == reflect.Map {
		keys := v.MapKeys()
		for _, key := range keys {
			mappedKey := mapper(key.Interface())
			mappedValue := mapper(v.MapIndex(key).Interface())
			result = reflect.Append(result, reflect.ValueOf(mappedKey))
			result = reflect.Append(result, reflect.ValueOf(mappedValue))
		}
	}
	return result.Interface()
}

// PredicateFilter filters a collection (slice or array) based on a predicate function and returns a new collection
// containing only the elements that satisfy the condition specified by the predicate.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice or array,
// and a filtering predicate function (`predicate`). The function applies the `predicate` function to each element in
// the collection, and if the predicate returns true, the element is included in the new collection. The function only
// supports slices and arrays as input collections.
//
// The function uses reflection to handle slices and arrays and constructs a new collection containing only the
// elements that pass the filter condition.
//
// Parameters:
//   - `collection`: The collection (slice or array) to filter. It can be of any type, but only slices and arrays
//     are supported.
//   - `predicate`: A function that takes a value from the collection and returns a boolean indicating whether the
//     element should be included in the resulting collection. If it returns true, the element is included.
//
// Returns:
//   - A new collection of the same type as the original collection, containing only the elements that satisfy the
//     condition defined by the `predicate` function.
//
// Example:
//
//	// Filtering a slice of integers to get only even numbers
//	numbers := []int{1, 2, 3, 4, 5, 6}
//	evens := PredicateFilter(numbers, func(value interface{}) bool {
//		return value.(int)%2 == 0
//	})
//	// evens will be []int{2, 4, 6}
//
// Notes:
//   - This function only works with slices or arrays, and will return an empty collection if the input is of another type.
//
// Limitations:
//   - The function creates a new collection based on the results of the `predicate` function, so it does not modify
//     the original collection.
func PredicateFilter(collection any, predicate func(value any) bool) any {
	v := reflect.ValueOf(collection)
	result := reflect.MakeSlice(v.Type(), 0, 0)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if predicate(item) {
				result = reflect.Append(result, reflect.ValueOf(item))
			}
		}
	}
	return result.Interface()
}

// Iterate iterates over a collection (slice, array, or map) and applies a callback function on each element.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice, array, or map,
// and a callback function. The callback function is executed for each element in the collection. For slices and arrays,
// the callback receives the index and the corresponding value. For maps, the callback receives each key and value, with
// the key being passed first followed by the value. For slices and arrays, the index is passed, while for maps, it is
// passed as -1 (since maps are unordered).
//
// The function uses reflection to handle different types of collections and ensures that the correct value is passed
// to the callback.
//
// Parameters:
//   - `collection`: The collection (slice, array, or map) to iterate over. It can be of any type.
//   - `callback`: A function that takes two arguments: the index (for slices and arrays, -1 for maps) and the value
//     from the collection. The callback is executed for each element in the collection.
//
// Example:
//
//	// Iterating over a slice
//	numbers := []int{1, 2, 3, 4}
//	Iterate(numbers, func(index int, value interface{}) {
//		fmt.Printf("Index: %d, Value: %v\n", index, value)
//	})
//	// Output:
//	// Index: 0, Value: 1
//	// Index: 1, Value: 2
//	// Index: 2, Value: 3
//	// Index: 3, Value: 4
//
//	// Iterating over a map
//	colors := map[string]string{"red": "FF0000", "green": "00FF00", "blue": "0000FF"}
//	Iterate(colors, func(index int, value interface{}) {
//		fmt.Printf("Value: %v\n", value)
//	})
//	// Output:
//	// Value: red
//	// Value: FF0000
//	// Value: green
//	// Value: 00FF00
//	// Value: blue
//	// Value: 0000FF
//
// Notes:
//   - For slices and arrays, the callback will receive the index and the value from the collection.
//   - For maps, the callback will be executed twice per key-value pair: once with the key and once with the value,
//     since maps are unordered and the order of key-value pairs cannot be guaranteed.
func Iterate(collection any, callback func(index int, value any)) {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			callback(i, v.Index(i).Interface())
		}
	} else if v.Kind() == reflect.Map {
		keys := v.MapKeys()
		for _, key := range keys {
			callback(-1, key.Interface())
			callback(-1, v.MapIndex(key).Interface())
		}
	}
}
