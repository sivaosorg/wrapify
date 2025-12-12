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

// Json serializes the `pagination` instance into a compact JSON string.
//
// This function uses the `unify4g.JsonN` utility to generate a JSON representation
// of the `pagination` instance. The output is a compact JSON string with no additional
// whitespace or formatting, providing a minimalistic view of the pagination data.
//
// Returns:
//   - A compact JSON string representation of the `pagination` instance.
func (p *pagination) Json() string {
	return unify4g.JsonN(p.Respond())
}

// JsonPretty serializes the `pagination` instance into a prettified JSON string.
//
// This function uses the `unify4g.JsonPrettyN` utility to generate a JSON representation
// of the `pagination` instance. The output is a human-readable JSON string with
// proper indentation and formatting for better readability, which is helpful for
// inspecting pagination data during development or debugging.
//
// Returns:
//   - A prettified JSON string representation of the `pagination` instance.
func (p *pagination) JsonPretty() string {
	return unify4g.JsonPrettyN(p.Respond())
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

// Json serializes the `header` instance into a compact JSON string.
//
// This function uses the `unify4g.JsonN` utility to create a compact JSON representation
// of the `header` instance. The resulting string contains only the key information, formatted
// with minimal whitespace, making it suitable for compact storage or transmission of header data.
//
// Returns:
//   - A compact JSON string representation of the `header` instance.
func (h *header) Json() string {
	return unify4g.JsonN(h.Respond())
}

// JsonPretty serializes the `header` instance into a prettified JSON string.
//
// This function uses the `unify4g.JsonPrettyN` utility to produce a formatted, human-readable
// JSON string representation of the `header` instance. The output is structured with indentation
// and newlines, making it ideal for inspecting header data in a clear, easy-to-read format, especially
// during debugging or development.
//
// Returns:
//   - A prettified JSON string representation of the `header` instance, formatted for improved readability.
func (h *header) JsonPretty() string {
	return unify4g.JsonPrettyN(h.Respond())
}

// WithPage sets the page number for the `pagination` instance.
//
// This function updates the `page` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
//
// Parameters:
//   - `v`: An integer representing the page number to set.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithPage(v int) *pagination {
	if v < 1 {
		v = 1
	}
	p.page = v
	return p
}

// WithPerPage sets the number of items per page for the `pagination` instance.
//
// This function updates the `perPage` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
// Validates that perPage is >= 1, defaults to 10 if invalid value.
//
// Parameters:
//   - `v`: An integer representing the number of items per page to set.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithPerPage(v int) *pagination {
	if v < 1 {
		v = 10
	}
	p.perPage = v
	return p
}

// WithTotalPages sets the total number of pages for the `pagination` instance.
// Ensure that totalPages is >= 0, defaults to 0 if invalid value.
//
// This function updates the `totalPages` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
//
// Parameters:
//   - `v`: An integer representing the total number of pages to set.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithTotalPages(v int) *pagination {
	if v < 0 {
		v = 0
	}
	p.totalPages = v
	return p
}

// WithTotalItems sets the total number of items for the `pagination` instance.
// Ensure that totalItems is >= 0, defaults to 0 if invalid value.
//
// This function updates the `totalItems` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
//
// Parameters:
//   - `v`: An integer representing the total number of items to set.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithTotalItems(v int) *pagination {
	if v < 0 {
		v = 0
	}
	p.totalItems = v
	return p
}

// WithIsLast sets whether this is the last page in the `pagination` instance.
//
// This function updates the `isLast` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
//
// Parameters:
//   - `v`: A boolean value indicating whether this is the last page.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithIsLast(v bool) *pagination {
	p.isLast = v
	return p
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

// WithCode sets the `code` field of the `header` instance.
//
// This function assigns the provided integer value to the `code` field of the `header`
// and returns the updated `header` instance, allowing for method chaining.
//
// Parameters:
//   - `v`: The integer value to set as the HTTP status code.
//
// Returns:
//   - The updated `header` instance with the `code` field set to the provided value.
func (h *header) WithCode(v int) *header {
	h.code = v
	return h
}

// WithText sets the `text` field of the `header` instance.
//
// This function assigns the provided string value to the `text` field of the `header`
// and returns the updated `header` instance, allowing for method chaining.
//
// Parameters:
//   - `v`: The string value to set as the text message.
//
// Returns:
//   - The updated `header` instance with the `text` field set to the provided value.
func (h *header) WithText(v string) *header {
	h.text = v
	return h
}

// WithType sets the `Type` field of the `header` instance.
//
// This function assigns the provided string value to the `Type` field of the `header`
// and returns the updated `header` instance, allowing for method chaining.
//
// Parameters:
//   - `v`: The string value to set as the type of the header.
//
// Returns:
//   - The updated `header` instance with the `Type` field set to the provided value.
func (h *header) WithType(v string) *header {
	h.typez = v
	return h
}

// WithDescription sets the `description` field of the `header` instance.
//
// This function assigns the provided string value to the `description` field of the `header`
// and returns the updated `header` instance, allowing for method chaining.
//
// Parameters:
//   - `v`: The string value to set as the description of the header.
//
// Returns:
//   - The updated `header` instance with the `description` field set to the provided value.
func (h *header) WithDescription(v string) *header {
	h.description = v
	return h
}

// Respond generates a map representation of the `pagination` instance.
//
// This method collects various fields related to pagination (e.g., `page`, `per_page`, etc.)
// and organizes them into a key-value map. It ensures that only valid pagination details
// are included in the response.
//
// The following fields are included in the pagination response:
//   - `page`: The current page number.
//   - `per_page`: The number of items per page.
//   - `total_pages`: The total number of pages available.
//   - `total_items`: The total number of items available across all pages.
//   - `is_last`: A boolean indicating if this is the last page.
//
// Returns:
//   - A `map[string]interface{}` containing the structured pagination data.
func (p *pagination) Respond() map[string]any {
	m := make(map[string]any)
	if !p.Available() {
		return m
	}
	m["page"] = p.page
	m["per_page"] = p.perPage
	m["total_pages"] = p.totalPages
	m["total_items"] = p.totalItems
	m["is_last"] = p.isLast
	return m
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

// Respond generates a map representation of the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) and includes the
// values of its fields in the returned map. Only the fields that are present (i.e., non-empty)
// are added to the map, ensuring a clean and concise response.
//
// Fields included in the response:
//   - `code`: The HTTP status code, if present and greater than 0.
//   - `text`: The associated text message, if present and not empty.
//   - `type`: The type of the header, if present and not empty.
//   - `description`: A description related to the header, if present and not empty.
//
// Returns:
//   - A `map[string]interface{}` containing the fields of the `header` instance that are present.
func (h *header) Respond() map[string]any {
	m := make(map[string]any)
	if !h.Available() {
		return m
	}
	if h.IsCodePresent() {
		m["code"] = h.code
	}
	if h.IsTextPresent() {
		m["text"] = h.text
	}
	if h.IsTypePresent() {
		m["type"] = h.typez
	}
	if h.IsDescriptionPresent() {
		m["description"] = h.description
	}
	return m
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
