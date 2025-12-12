package utilities

// Pair represents a pair of values, used by Zip function.
//
// Fields:
//   - First: The first value of type T.
//   - Second: The second value of type U.
type Pair[T any, U any] struct {
	First  T
	Second U
}
