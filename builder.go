package wrapify

import (
	"fmt"
	"time"

	"github.com/sivaosorg/unify4g"
)

// NewPagination creates a new instance of the `pagination` struct.
//
// This function initializes a `pagination` struct with its default values.
//
// Returns:
//   - A pointer to a newly created `pagination` instance.
func NewPagination() *pagination {
	p := &pagination{}
	return p
}

// NewMeta creates a new instance of the `meta` struct.
//
// This function initializes a `meta` struct with its default values,
// including an empty `CustomFields` map.
//
// Returns:
//   - A pointer to a newly created `meta` instance with initialized fields.
func NewMeta() *meta {
	m := &meta{
		customFields: map[string]any{},
	}
	return m
}

// NewHeader creates a new instance of the `header` struct.
//
// This function initializes a `header` struct with its default values.
//
// Returns:
//   - A pointer to a newly created `header` instance.
func NewHeader() *header {
	h := &header{}
	return h
}

// New creates a new instance of the `wrapper` struct.
//
// This function initializes a `wrapper` struct with its default values,
// including an empty map for the `Debug` field.
//
// Returns:
//   - A pointer to a newly created `wrapper` instance with initialized fields.
func New() *wrapper {
	w := &wrapper{
		meta: defaultMetaValues(),
	}
	return w
}

// Json serializes the `meta` instance into a compact JSON string.
//
// This function uses the `unify4g.JsonN` utility to create a compact JSON representation
// of the `meta` instance. The resulting string is formatted without additional whitespace,
// suitable for efficient storage or transmission of metadata.
//
// Returns:
//   - A compact JSON string representation of the `meta` instance.
func (m *meta) Json() string {
	return unify4g.JsonN(m.Respond())
}

// JsonPretty serializes the `meta` instance into a prettified JSON string.
//
// This function calls the `unify4g.JsonPrettyN` utility to produce a formatted, human-readable
// JSON string representation of the `meta` instance. The output is useful for debugging
// or inspecting metadata in a more structured format.
//
// Returns:
//   - A prettified JSON string representation of the `meta` instance.
func (m *meta) JsonPretty() string {
	return unify4g.JsonPrettyN(m.Respond())
}

// WithApiVersion sets the API version for the `meta` instance.
//
// This function updates the `apiVersion` field of the `meta` instance with the specified value
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `v`: A string representing the API version to set.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithApiVersion(v string) *meta {
	m.apiVersion = v
	return m
}

// WithApiVersionf sets the API version for the `meta` instance using a formatted string.
//
// This function constructs a formatted string for the API version using the provided `format` string
// and arguments (`args`). It then assigns the formatted value to the `apiVersion` field of the `meta` instance.
// The method supports method chaining by returning a pointer to the modified `meta` instance.
//
// Parameters:
//   - format: A format string to construct the API version.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithApiVersionf(format string, args ...any) *meta {
	return m.WithApiVersion(fmt.Sprintf(format, args...))
}

// WithRequestID sets the request ID for the `meta` instance.
//
// This function updates the `requestID` field of the `meta` instance with the specified value
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `v`: A string representing the request ID to set.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithRequestID(v string) *meta {
	m.requestID = v
	return m
}

// WithRequestIDf sets the request ID for the `meta` instance using a formatted string.
//
// This function constructs a formatted string for the request ID using the provided `format` string
// and arguments (`args`). It then assigns the formatted value to the `requestID` field of the `meta` instance.
// The method supports method chaining by returning a pointer to the modified `meta` instance.
//
// Parameters:
//   - format: A format string to construct the request ID.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithRequestIDf(format string, args ...any) *meta {
	return m.WithRequestID(fmt.Sprintf(format, args...))
}

// WithLocale sets the locale for the `meta` instance.
//
// This function updates the `locale` field of the `meta` instance with the specified value
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `v`: A string representing the locale to set.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithLocale(v string) *meta {
	m.locale = v
	return m
}

// WithRequestedTime sets the requested time for the `meta` instance.
//
// This function updates the `requestedTime` field of the `meta` instance with the specified value
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `v`: A `time.Time` object representing the requested time to set.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithRequestedTime(v time.Time) *meta {
	m.requestedTime = v
	return m
}

// WithCustomFields sets the custom fields for the `meta` instance.
//
// This function updates the `customFields` map of the `meta` instance with the provided values
// and returns the updated `meta` instance for method chaining.
//
// Parameters:
//   - `values`: A map of string keys to interface{} values representing custom fields.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithCustomFields(values map[string]any) *meta {
	m.customFields = values
	return m
}

// WithCustomFieldKV sets a specific custom field key-value pair for the `meta` instance.
//
// This function adds or updates a custom field in the `customFields` map of the `meta` instance.
// If the `customFields` map is empty, it is initialized first.
//
// Parameters:
//   - `key`: A string representing the custom field key.
//   - `value`: An interface{} representing the value to associate with the custom field key.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithCustomFieldKV(key string, value any) *meta {
	if !m.IsCustomFieldPresent() {
		m.customFields = make(map[string]interface{})
	}
	m.customFields[key] = value
	return m
}

// WithCustomFieldKVf sets a specific custom field key-value pair for the `meta` instance
// using a formatted value.
//
// This function creates a formatted string value using the provided `format` string and
// `args`. It then calls `WithCustomFieldKV` to add or update the custom field with the
// specified key and the formatted value. The modified `meta` instance is returned for
// method chaining.
//
// Parameters:
//   - key: A string representing the key for the custom field.
//   - format: A format string to construct the value.
//   - args: A variadic list of arguments to be interpolated into the format string.
//
// Returns:
//   - A pointer to the modified `meta` instance, enabling method chaining.
func (m *meta) WithCustomFieldKVf(key string, format string, args ...any) *meta {
	return m.WithCustomFieldKV(key, fmt.Sprintf(format, args...))
}

// Respond generates a map representation of the `meta` instance.
//
// This method collects various fields of the `meta` instance (e.g., `apiVersion`, `requestID`, etc.)
// and organizes them into a key-value map. Only fields that are available and valid
// (e.g., non-empty or initialized) are included in the resulting map.
//
// Fields included in the response:
//   - `api_version`: The API version, if present.
//   - `request_id`: The unique request identifier, if present.
//   - `locale`: The locale information, if present.
//   - `requested_time`: The requested time, if it is initialized.
//   - `custom_fields`: A map of custom fields, if present.
//
// Returns:
//   - A `map[string]interface{}` containing the structured metadata.
func (m *meta) Respond() map[string]any {
	mk := make(map[string]any)
	if !m.Available() {
		return mk
	}
	if m.IsApiVersionPresent() {
		mk["api_version"] = m.apiVersion
	}
	if m.IsRequestIDPresent() {
		mk["request_id"] = m.requestID
	}
	if m.IsLocalePresent() {
		mk["locale"] = m.locale
	}
	if m.IsRequestedTimePresent() {
		mk["requested_time"] = m.requestedTime
	}
	if m.IsCustomFieldPresent() {
		mk["custom_fields"] = m.customFields
	}
	return mk
}

// Get returns buffer from pool
//
// If pool is empty, it creates a new buffer
// of predefined size
//
// Returns:
//   - A byte slice buffer
func (bp *BufferPool) Get() []byte {
	select {
	case buf := <-bp.buffers:
		return buf
	default:
		return make([]byte, bp.size)
	}
}

// Put returns buffer to pool
//
// # If pool is full, the buffer is discarded
//
// Parameters:
//   - buf: A byte slice buffer to be returned to the pool
//
// Returns:
//   - None
func (bp *BufferPool) Put(buf []byte) {
	if buf == nil {
		return
	}
	select {
	case bp.buffers <- buf:
	default:
		// Pool is full, discard
	}
}
