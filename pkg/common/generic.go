package common

import "reflect"

// Transform applies a transformation function to each element of a collection (slice, array, or map) and returns
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
//	squared := Transform(numbers, func(value interface{}) interface{} {
//		return value.(int) * value.(int)
//	})
//	// squared will be []int{1, 4, 9, 16}
//
//	// Mapping a map of strings to their lengths
//	words := map[string]string{"apple": "fruit", "carrot": "vegetable"}
//	lengths := Transform(words, func(value interface{}) interface{} {
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
func Transform(collection any, predicate func(value any) any) any {
	v := reflect.ValueOf(collection)
	result := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(predicate(v.Index(0).Interface()))), 0, 0)

	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			mappedValue := predicate(v.Index(i).Interface())
			result = reflect.Append(result, reflect.ValueOf(mappedValue))
		}
	} else if v.Kind() == reflect.Map {
		keys := v.MapKeys()
		for _, key := range keys {
			mappedKey := predicate(key.Interface())
			mappedValue := predicate(v.MapIndex(key).Interface())
			result = reflect.Append(result, reflect.ValueOf(mappedKey))
			result = reflect.Append(result, reflect.ValueOf(mappedValue))
		}
	}
	return result.Interface()
}

// Filter filters a collection (slice or array) based on a predicate function and returns a new collection
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
//	evens := Filter(numbers, func(value interface{}) bool {
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
func Filter(collection any, predicate func(value any) bool) any {
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

// Reduce reduces a collection (slice or array) to a single value by applying a reducer function
// to each element, combining them into a single result.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice or array,
// and a reducer function (`reducer`). The reducer function is applied to each element of the collection, along with an
// accumulator value, and the result of the reducer is passed as the accumulator to the next iteration. This process
// continues until all elements are processed, resulting in a single accumulated value.
//
// The function uses reflection to handle slices and arrays and iterates over the collection to apply the reducer function.
//
// Parameters:
//   - `collection`: The collection (slice or array) to reduce. It can be of any type, but only slices and arrays
//     are supported.
//   - `reducer`: A function that takes two arguments: an accumulator and an element from the collection. It returns
//     a new accumulator value after combining the accumulator and the element.
//   - `initialValue`: The initial value of the accumulator, used as the starting point for the reduction process.
//
// Returns:
//   - A single value, which is the result of reducing the entire collection using the `reducer` function. The return
//     type is the same as the type of the `initialValue`.
//
// Example:
//
//	// Reducing a slice of integers by summing them
//	numbers := []int{1, 2, 3, 4, 5}
//	sum := Reduce(numbers, func(acc interface{}, value interface{}) interface{} {
//		return acc.(int) + value.(int)
//	}, 0)
//	// sum will be 15
//
// Notes:
//   - This function only works with slices or arrays, and will return the initial value if the input collection is empty.
//
// Limitations:
//   - The function creates a single accumulated result by repeatedly applying the `reducer` function to each element,
//     so it does not modify the original collection.
func Reduce(collection any, reducer func(acc any, value any) any, initialValue any) any {
	v := reflect.ValueOf(collection)
	accumulator := initialValue
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			accumulator = reducer(accumulator, v.Index(i).Interface())
		}
	}
	return accumulator
}

// Find searches for the first element in a collection (slice or array) that satisfies a given predicate
// function and returns it.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice or array,
// and a predicate function (`predicate`). The predicate function is applied to each element of the collection,
// and if it returns true for an element, that element is returned as the result. The function will return the first
// element that satisfies the condition, and if no elements satisfy the condition, it returns `nil`.
//
// The function uses reflection to handle slices and arrays, iterating over the collection to check each element
// with the provided predicate.
//
// Parameters:
//   - `collection`: The collection (slice or array) to search through. It can be of any type, but only slices and arrays
//     are supported.
//   - `predicate`: A function that takes a value from the collection and returns a boolean indicating whether the element
//     satisfies the condition. If it returns true, the element is returned.
//
// Returns:
//   - The first element from the collection that satisfies the `predicate` function. If no elements satisfy the condition,
//     it returns `nil`.
//
// Example:
//
//	// Finding the first even number in a slice of integers
//	numbers := []int{1, 2, 3, 4, 5}
//	even := Find(numbers, func(value interface{}) bool {
//		return value.(int)%2 == 0
//	})
//	// even will be 2 (the first even number)
//
// Notes:
//   - This function only works with slices or arrays. If the collection is of another type, it will return `nil`.
//   - The function returns the first element that matches the predicate and stops searching after finding it.
//
// Limitations:
//   - The function works only with slices and arrays. If no elements satisfy the predicate, the function will return `nil`,
//     even if the collection is non-empty.
func Find(collection any, predicate func(value any) bool) any {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if predicate(item) {
				return item
			}
		}
	}
	return nil
}

// All checks whether all elements in a collection (slice or array) satisfy a given condition.
//
// This function takes a collection of any type (using an empty `interface{}`), which can be a slice or array,
// and a condition function (`condition`). The condition function is applied to each element of the collection, and
// the function returns `true` if all elements satisfy the condition. If any element does not satisfy the condition,
// the function returns `false`. If the collection is empty, the function returns `true` (since the condition is trivially satisfied).
//
// The function uses reflection to handle slices and arrays, iterating over the collection to check each element
// with the provided condition.
//
// Parameters:
//   - `collection`: The collection (slice or array) to check. It can be of any type, but only slices and arrays
//     are supported.
//   - `condition`: A function that takes a value from the collection and returns a boolean indicating whether the
//     element satisfies the condition. If it returns `false` for any element, the function immediately returns `false`.
//
// Returns:
//   - `true` if all elements in the collection satisfy the condition, otherwise `false`.
//   - If the collection is empty, the function returns `true`.
//
// Example:
//
//	// Checking if all elements in a slice of integers are positive
//	numbers := []int{1, 2, 3, 4, 5}
//	allPositive := All(numbers, func(value interface{}) bool {
//		return value.(int) > 0
//	})
//	// allPositive will be true
//
// Notes:
//   - This function only works with slices or arrays. If the collection is of another type, it will return `false`.
//   - The function returns `false` as soon as it finds an element that does not satisfy the condition, making it more
//     efficient for early termination.
//
// Limitations:
//   - The function works only with slices and arrays. If no elements satisfy the condition, the function will return `false`,
//     but if all elements are valid, it will return `true`. An empty collection is considered to trivially satisfy the condition.
func All(collection any, condition func(value any) bool) bool {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if !condition(v.Index(i).Interface()) {
				return false
			}
		}
		return true
	}
	return false
}

// Any checks whether any element in a collection (slice or array) satisfies a given condition.
//
// This function accepts a collection (slice or array) and a condition function. It returns `true` if at least one element in the collection satisfies the condition,
// and `false` if none do. The function stops iterating as soon as a matching element is found, making it more efficient for early termination.
//
// Parameters:
//   - `collection`: A slice or array of any type to check.
//   - `condition`: A function that takes a value from the collection and returns a boolean indicating whether the element satisfies the condition.
//
// Returns:
//   - `true` if any element satisfies the condition, `false` otherwise.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	anyNegative := Any(numbers, func(value interface{}) bool {
//	  return value.(int) < 0
//	})
//	// anyNegative will be false because no element is negative.
func Any(collection any, condition func(value any) bool) bool {
	v := reflect.ValueOf(collection)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if condition(v.Index(i).Interface()) {
				return true
			}
		}
		return false
	}
	return false
}

// Count returns the number of elements in a collection (slice or array) that satisfy a given condition.
//
// This function takes a collection (slice or array) and a condition function. It iterates through the collection, applying the condition to each element.
// It returns the total count of elements that satisfy the condition.
//
// Parameters:
//   - `collection`: A slice or array of any type to check.
//   - `condition`: A function that checks if an element satisfies a condition. The function returns `true` for elements that match the condition, and `false` otherwise.
//
// Returns:
//   - The count of elements in the collection that satisfy the condition.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	countNegative := Count(numbers, func(value interface{}) bool {
//	  return value.(int) < 0
//	})
//	// countNegative will be 0, since no element is negative.
func Count(collection any, condition func(value any) bool) int {
	v := reflect.ValueOf(collection)
	count := 0
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			if condition(v.Index(i).Interface()) {
				count++
			}
		}
	}
	return count
}

// Remove returns a new collection (slice or array) where all elements that satisfy a given condition are removed.
//
// This function takes a collection (slice or array) and a condition function. It iterates through the collection, and for each element that does not satisfy the condition,
// it is added to a new result collection. The function returns a new collection that contains only the elements that do not match the condition.
//
// Parameters:
//   - `collection`: A slice or array of any type to process.
//   - `condition`: A function that checks if an element satisfies a condition. If the element matches the condition, it is removed from the result collection.
//
// Returns:
//   - A new collection (slice or array) with elements that do not satisfy the condition.
//
// Example:
//
//	numbers := []int{1, 2, 3, 4, 5}
//	result := Remove(numbers, func(value interface{}) bool {
//	  return value.(int) % 2 == 0 // Removes even numbers
//	})
//	// result will be []int{1, 3, 5}
func Remove(collection any, condition func(value any) bool) any {
	v := reflect.ValueOf(collection)
	result := reflect.MakeSlice(v.Type(), 0, 0)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			if !condition(item) {
				result = reflect.Append(result, reflect.ValueOf(item))
			}
		}
	}
	return result.Interface()
}
