package wrapify

import (
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
		customFields: map[string]interface{}{},
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

// NewWrap creates a new instance of the `wrapper` struct.
//
// This function initializes a `wrapper` struct with its default values,
// including an empty map for the `Debug` field.
//
// Returns:
//   - A pointer to a newly created `wrapper` instance with initialized fields.
func NewWrap() *wrapper {
	w := &wrapper{
		meta: NewMeta().
			WithApiVersion("v0.0.1").
			WithLocale("en_US"). // vi_VN, en_US
			WithRequestID(unify4g.GenerateCryptoID()),
	}
	return w
}

// Json serializes the `wrapper` instance into a compact JSON string.
//
// This function uses the `unify4g.JsonN` utility to generate a JSON representation
// of the `wrapper` instance. The output is a compact JSON string with no additional
// whitespace or formatting.
//
// Returns:
//   - A compact JSON string representation of the `wrapper` instance.
func (w *wrapper) Json() string {
	return unify4g.JsonN(w.Respond())
}

// JsonPretty serializes the `wrapper` instance into a prettified JSON string.
//
// This function uses the `unify4g.JsonPrettyN` utility to generate a JSON representation
// of the `wrapper` instance. The output is a human-readable JSON string with
// proper indentation and formatting for better readability.
//
// Returns:
//   - A prettified JSON string representation of the `wrapper` instance.
func (w *wrapper) JsonPretty() string {
	return unify4g.JsonPrettyN(w.Respond())
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
	p.page = v
	return p
}

// WithPerPage sets the number of items per page for the `pagination` instance.
//
// This function updates the `perPage` field of the `pagination` and
// returns the modified `pagination` instance to allow method chaining.
//
// Parameters:
//   - `v`: An integer representing the number of items per page to set.
//
// Returns:
//   - A pointer to the modified `pagination` instance (enabling method chaining).
func (p *pagination) WithPerPage(v int) *pagination {
	p.perPage = v
	return p
}

// WithTotalPages sets the total number of pages for the `pagination` instance.
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
	p.totalPages = v
	return p
}

// WithTotalItems sets the total number of items for the `pagination` instance.
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
func (m *meta) WithCustomFields(values map[string]interface{}) *meta {
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
func (m *meta) WithCustomFieldKV(key string, value interface{}) *meta {
	if !m.IsCustomFieldPresent() {
		m.customFields = make(map[string]interface{})
	}
	m.customFields[key] = value
	return m
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
	h.Type = v
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

// WithStatusCode sets the HTTP status code for the `wrapper` instance.
//
// This function updates the `statusCode` field of the `wrapper` and
// returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `code`: An integer representing the HTTP status code to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithStatusCode(code int) *wrapper {
	w.statusCode = code
	return w
}

// WithTotal sets the total number of items for the `wrapper` instance.
//
// This function updates the `total` field of the `wrapper` and
// returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `total`: An integer representing the total number of items to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithTotal(total int) *wrapper {
	w.total = total
	return w
}

// WithMessage sets a message for the `wrapper` instance.
//
// This function updates the `message` field of the `wrapper` with the provided string
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `message`: A string message to be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithMessage(message string) *wrapper {
	w.message = message
	return w
}

// WithBody sets the body data for the `wrapper` instance.
//
// This function updates the `data` field of the `wrapper` with the provided value
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: The value to be set as the body data, which can be any type.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithBody(v interface{}) *wrapper {
	w.data = v
	return w
}

// WithPath sets the request path for the `wrapper` instance.
//
// This function updates the `path` field of the `wrapper` with the provided string
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A string representing the request path.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithPath(v string) *wrapper {
	w.path = v
	return w
}

// WithHeader sets the header for the `wrapper` instance.
//
// This function updates the `header` field of the `wrapper` with the provided `header`
// instance and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A pointer to a `header` struct that will be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithHeader(v *header) *wrapper {
	w.header = v
	return w
}

// WithMeta sets the metadata for the `wrapper` instance.
//
// This function updates the `meta` field of the `wrapper` with the provided `meta`
// instance and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A pointer to a `meta` struct that will be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithMeta(v *meta) *wrapper {
	w.meta = v
	return w
}

// WithPagination sets the pagination information for the `wrapper` instance.
//
// This function updates the `pagination` field of the `wrapper` with the provided `pagination`
// instance and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A pointer to a `pagination` struct that will be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithPagination(v *pagination) *wrapper {
	w.pagination = v
	return w
}

// WithDebugging sets the debugging information for the `wrapper` instance.
//
// This function updates the `debug` field of the `wrapper` with the provided map of debugging data
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `v`: A map containing debugging information to be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithDebugging(v map[string]interface{}) *wrapper {
	w.debug = v
	return w
}

// WithError sets an error for the `wrapper` instance.
//
// This function updates the `errors` field of the `wrapper` with the provided error
// and returns the modified `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `err`: An error object to be set in the `wrapper`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithError(err error) *wrapper {
	w.errors = err
	return w
}

// WithDebuggingKV adds a key-value pair to the debugging information in the `wrapper` instance.
//
// This function checks if debugging information is already present. If it is not, it initializes
// an empty map. Then it adds the given key-value pair to the `debug` map and returns the modified
// `wrapper` instance to allow method chaining.
//
// Parameters:
//   - `key`: The key for the debugging information to be added.
//   - `value`: The value associated with the key to be added to the `debug` map.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithDebuggingKV(key string, value interface{}) *wrapper {
	if !w.IsDebuggingPresent() {
		w.debug = make(map[string]interface{})
	}
	w.debug[key] = value
	return w
}

// WithApiVersion sets the API version in the `meta` field of the `wrapper` instance.
//
// This function checks if the `meta` information is present in the `wrapper`. If it is not,
// a new `meta` instance is created. Then, it calls the `WithApiVersion` method on the `meta`
// instance to set the API version.
//
// Parameters:
//   - `v`: A string representing the API version to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithApiVersion(v string) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = NewMeta()
	}
	w.meta.WithApiVersion(v)
	return w
}

// WithRequestID sets the request ID in the `meta` field of the `wrapper` instance.
//
// This function ensures that if `meta` information is not already set in the `wrapper`, a new
// `meta` instance is created. Then, it calls the `WithRequestID` method on the `meta` instance
// to set the request ID.
//
// Parameters:
//   - `v`: A string representing the request ID to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithRequestID(v string) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = NewMeta()
	}
	w.meta.WithRequestID(v)
	return w
}

// WithLocale sets the locale in the `meta` field of the `wrapper` instance.
//
// This function ensures the `meta` field is present, creating a new instance if needed, and
// sets the locale in the `meta` using the `WithLocale` method.
//
// Parameters:
//   - `v`: A string representing the locale to set.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithLocale(v string) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = NewMeta()
	}
	w.meta.WithLocale(v)
	return w
}

// WithRequestedTime sets the requested time in the `meta` field of the `wrapper` instance.
//
// This function ensures that the `meta` field exists, and if not, creates a new one. It then
// sets the requested time in the `meta` using the `WithRequestedTime` method.
//
// Parameters:
//   - `v`: A `time.Time` value representing the requested time.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithRequestedTime(v time.Time) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = NewMeta()
	}
	w.meta.WithRequestedTime(v)
	return w
}

// WithCustomFields sets the custom fields in the `meta` field of the `wrapper` instance.
//
// This function checks if the `meta` field is present. If not, it creates a new `meta` instance
// and sets the provided custom fields using the `WithCustomFields` method.
//
// Parameters:
//   - `values`: A map representing the custom fields to set in the `meta`.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithCustomFields(values map[string]interface{}) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = NewMeta()
	}
	w.meta.WithCustomFields(values)
	return w
}

// WithCustomFieldKV sets a specific custom field key-value pair in the `meta` field of the `wrapper` instance.
//
// This function ensures that if the `meta` field is not already set, a new `meta` instance is created.
// It then adds the provided key-value pair to the custom fields of `meta` using the `WithCustomFieldKV` method.
//
// Parameters:
//   - `key`: A string representing the custom field key to set.
//   - `value`: The value associated with the custom field key.
//
// Returns:
//   - A pointer to the modified `wrapper` instance (enabling method chaining).
func (w *wrapper) WithCustomFieldKV(key string, value interface{}) *wrapper {
	if !w.IsMetaPresent() {
		w.meta = NewMeta()
	}
	w.meta.WithCustomFieldKV(key, value)
	return w
}

// WithPage sets the current page number in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified page number is then
// applied to the pagination instance.
//
// Parameters:
//   - v: The page number to set.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithPage(v int) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = NewPagination()
	}
	w.pagination.WithPage(v)
	return w
}

// WithPerPage sets the number of items per page in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified items-per-page value
// is then applied to the pagination instance.
//
// Parameters:
//   - v: The number of items per page to set.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithPerPage(v int) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = NewPagination()
	}
	w.pagination.WithPerPage(v)
	return w
}

// WithTotalPages sets the total number of pages in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified total pages value
// is then applied to the pagination instance.
//
// Parameters:
//   - v: The total number of pages to set.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithTotalPages(v int) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = NewPagination()
	}
	w.pagination.WithTotalPages(v)
	return w
}

// WithTotalItems sets the total number of items in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified total items value
// is then applied to the pagination instance.
//
// Parameters:
//   - v: The total number of items to set.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithTotalItems(v int) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = NewPagination()
	}
	w.pagination.WithTotalItems(v)
	return w
}

// WithIsLast sets whether the current page is the last one in the wrapper's pagination.
//
// If the pagination object is not already initialized, it creates a new one
// using the `NewPagination` function. The specified boolean value is then
// applied to indicate whether the current page is the last.
//
// Parameters:
//   - v: A boolean indicating whether the current page is the last.
//
// Returns:
//   - A pointer to the updated `wrapper` instance.
func (w *wrapper) WithIsLast(v bool) *wrapper {
	if !w.IsPagingPresent() {
		w.pagination = NewPagination()
	}
	w.pagination.WithIsLast(v)
	return w
}

// Respond generates a map representation of the `wrapper` instance.
//
// This method collects various fields of the `wrapper` (e.g., `data`, `header`, `meta`, etc.)
// and organizes them into a key-value map. Only non-nil or meaningful fields are added
// to the resulting map to ensure a clean and concise response structure.
//
// Fields included in the response:
//   - `data`: The primary data payload, if present.
//   - `headers`: The structured header details, if present.
//   - `meta`: Metadata about the response, if present.
//   - `pagination`: Pagination details, if applicable.
//   - `debug`: Debugging information, if provided.
//   - `total`: Total number of items, if set to a valid non-negative value.
//   - `status_code`: The HTTP status code, if greater than 0.
//   - `message`: A descriptive message, if not empty.
//   - `path`: The request path, if not empty.
//
// Returns:
//   - A `map[string]interface{}` containing the structured response data.
func (w *wrapper) Respond() map[string]interface{} {
	m := make(map[string]interface{})
	if !w.Available() {
		return m
	}
	if w.IsBodyPresent() {
		m["data"] = w.data
	}
	if w.IsHeaderPresent() {
		m["headers"] = w.header.Respond()
	}
	if w.IsMetaPresent() {
		m["meta"] = w.meta.Respond()
	}
	if w.IsPagingPresent() {
		m["pagination"] = w.pagination.Respond()
	}
	if w.IsDebuggingPresent() {
		m["debug"] = w.debug
	}
	if w.IsTotalPresent() {
		m["total"] = w.total
	}
	if w.IsStatusCodePresent() {
		m["status_code"] = w.statusCode
	}
	if unify4g.IsNotEmpty(w.message) {
		m["message"] = w.message
	}
	if unify4g.IsNotEmpty(w.path) {
		m["path"] = w.path
	}
	return m
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
func (p *pagination) Respond() map[string]interface{} {
	m := make(map[string]interface{})
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
func (m *meta) Respond() map[string]interface{} {
	mk := make(map[string]interface{})
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
func (h *header) Respond() map[string]interface{} {
	m := make(map[string]interface{})
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
		m["type"] = h.Type
	}
	if h.IsDescriptionPresent() {
		m["description"] = h.description
	}
	return m
}
