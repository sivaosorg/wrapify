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

// OptionsConfig defines the configuration options for pretty-printing JSON data.
// It allows customization of width, prefix, indentation, and sorting of keys.
// These options control how the JSON output will be formatted.
//
// Fields:
//   - Width: The maximum column width for single-line arrays. This prevents arrays from becoming too wide.
//     Default is 80 characters.
//   - Prefix: A string that will be prepended to each line of the output. Useful for adding custom prefixes
//     or structuring the output with additional information. Default is an empty string.
//   - Indent: The string used for indentation in nested JSON structures. Default is two spaces ("  ").
//   - SortKeys: A flag indicating whether the keys in JSON objects should be sorted alphabetically. Default is false.
type OptionsConfig struct {
	// Width is an max column width for single line arrays
	// Default is 80
	Width int `json:"width"`
	// Prefix is a prefix for all lines
	// Default is an empty string
	Prefix string `json:"prefix"`
	// Indent is the nested indentation
	// Default is two spaces
	Indent string `json:"indent"`
	// SortKeys will sort the keys alphabetically
	// Default is false
	SortKeys bool `json:"sort_keys"`
}

// Style is the color style
type Style struct {
	Key, String, Number [2]string
	True, False, Null   [2]string
	Escape              [2]string
	Brackets            [2]string
	Append              func(dst []byte, c byte) []byte
}

// DefaultOptionsConfig is a pre-configured default set of options for pretty-printing JSON.
// This configuration uses a width of 80, an empty prefix, two-space indentation, and does not sort keys.
// It is used when no custom options are provided in the PrettyOptions function.
var DefaultOptionsConfig = &OptionsConfig{Width: 80, Prefix: "", Indent: "  ", SortKeys: false}

// TerminalStyle is for terminals
var TerminalStyle *Style

type result int
type byKind int

// jsonType represents the different types of JSON values.
//
// This enumeration defines constants representing various JSON data types, including `null`, `boolean`, `number`,
// `string`, and `JSON object or array`. These constants are used by the `getJsonType` function to identify the type
// of a given JSON value based on its first character.
type jsonType int

type pair struct {
	keyStart, keyEnd     int
	valueStart, valueEnd int
}

// byKeyVal is a struct that provides a way to sort JSON key-value pairs.
// It contains the JSON data, a buffer to hold trimmed values, and a list of pairs to be sorted.
type byKeyVal struct {
	sorted bool   // indicates whether the pairs are sorted
	json   []byte // original JSON data
	buf    []byte // buffer used for processing values
	pairs  []pair // list of key-value pairs to sort
}
