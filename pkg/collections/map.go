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
