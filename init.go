package wrapify

import (
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
		CustomFields: map[string]interface{}{},
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
	w := &wrapper{}
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
		w.debug = map[string]interface{}{}
	}
	w.debug[key] = value
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
	if w == nil {
		return m
	}
	if w.IsBodyPresent() {
		m["data"] = w.data
	}
	if w.IsHeaderPresent() {
		m["headers"] = w.header
	}
	if w.IsMetaPresent() {
		m["meta"] = w.meta
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
	if p == nil {
		return m
	}
	m["page"] = p.page
	m["per_page"] = p.perPage
	m["total_pages"] = p.totalPages
	m["total_items"] = p.totalItems
	m["is_last"] = p.isLast
	return m
}

// IsDebuggingPresent checks whether debugging information is present in the `wrapper` instance.
//
// This function verifies if the `debug` field of the `wrapper` is not nil and contains at least one entry.
// It returns `true` if debugging information is available; otherwise, it returns `false`.
//
// Returns:
//   - A boolean value indicating whether debugging information is present:
//   - `true` if `debug` is not nil and contains data.
//   - `false` if `debug` is nil or empty.
func (w *wrapper) IsDebuggingPresent() bool {
	return len(w.debug) > 0
}

// IsDebuggingKeyPresent checks whether a specific key exists in the `debug` information.
//
// This function first checks if debugging information is present using `IsDebuggingPresent()`.
// Then it uses `unify4g.MapContainsKey` to verify if the given key is present within the `debug` map.
//
// Parameters:
//   - `key`: The key to search for within the `debug` field.
//
// Returns:
//   - A boolean value indicating whether the specified key is present in the `debug` map:
//   - `true` if the `debug` field is present and contains the specified key.
//   - `false` if `debug` is nil or does not contain the key.
func (w *wrapper) IsDebuggingKeyPresent(key string) bool {
	return w.IsDebuggingPresent() && unify4g.MapContainsKey(w.debug, key)
}

// IsBodyPresent checks whether the body data is present in the `wrapper` instance.
//
// This function checks if the `data` field of the `wrapper` is not nil, indicating that the body contains data.
//
// Returns:
//   - A boolean value indicating whether the body data is present:
//   - `true` if `data` is not nil.
//   - `false` if `data` is nil.
func (w *wrapper) IsBodyPresent() bool {
	return w.data != nil
}

// IsHeaderPresent checks whether header information is present in the `wrapper` instance.
//
// This function checks if the `header` field of the `wrapper` is not nil, indicating that header information is included.
//
// Returns:
//   - A boolean value indicating whether header information is present:
//   - `true` if `header` is not nil.
//   - `false` if `header` is nil.
func (w *wrapper) IsHeaderPresent() bool {
	return w.header != nil
}

// IsMetaPresent checks whether metadata information is present in the `wrapper` instance.
//
// This function checks if the `meta` field of the `wrapper` is not nil, indicating that metadata is available.
//
// Returns:
//   - A boolean value indicating whether metadata is present:
//   - `true` if `meta` is not nil.
//   - `false` if `meta` is nil.
func (w *wrapper) IsMetaPresent() bool {
	return w.meta != nil
}

// IsPagingPresent checks whether pagination information is present in the `wrapper` instance.
//
// This function checks if the `pagination` field of the `wrapper` is not nil, indicating that pagination details are included.
//
// Returns:
//   - A boolean value indicating whether pagination information is present:
//   - `true` if `pagination` is not nil.
//   - `false` if `pagination` is nil.
func (w *wrapper) IsPagingPresent() bool {
	return w.pagination != nil
}

// IsErrorPresent checks whether an error is present in the `wrapper` instance.
//
// This function checks if the `errors` field of the `wrapper` is not nil, indicating that an error has occurred.
//
// Returns:
//   - A boolean value indicating whether an error is present:
//   - `true` if `errors` is not nil.
//   - `false` if `errors` is nil.
func (w *wrapper) IsErrorPresent() bool {
	return w.errors != nil
}

// IsTotalPresent checks whether the total number of items is present in the `wrapper` instance.
//
// This function checks if the `total` field of the `wrapper` is greater than or equal to 0,
// indicating that a valid total number of items has been set.
//
// Returns:
//   - A boolean value indicating whether the total is present:
//   - `true` if `total` is greater than or equal to 0.
//   - `false` if `total` is negative (indicating no total value).
func (w *wrapper) IsTotalPresent() bool {
	return w.total >= 0
}

// IsStatusCodePresent checks whether a valid status code is present in the `wrapper` instance.
//
// This function checks if the `statusCode` field of the `wrapper` is greater than 0,
// indicating that a valid HTTP status code has been set.
//
// Returns:
//   - A boolean value indicating whether the status code is present:
//   - `true` if `statusCode` is greater than 0.
//   - `false` if `statusCode` is less than or equal to 0.
func (w *wrapper) IsStatusCodePresent() bool {
	return w.statusCode > 0
}

// IsError checks whether there is an error present in the `wrapper` instance.
//
// This function returns `true` if the `wrapper` contains an error, which can be any of the following:
//   - An error present in the `errors` field.
//   - A client error (4xx status code) or a server error (5xx status code).
//
// Returns:
//   - A boolean value indicating whether there is an error:
//   - `true` if there is an error present, either in the `errors` field or as an HTTP client/server error.
//   - `false` if no error is found.
func (w *wrapper) IsError() bool {
	return w.IsErrorPresent() || w.IsClientError() || w.IsServerError()
}

// IsSuccess checks whether the HTTP status code indicates a successful response.
//
// This function checks if the `statusCode` is between 200 and 299, inclusive, which indicates a successful HTTP response.
//
// Returns:
//   - A boolean value indicating whether the HTTP response was successful:
//   - `true` if the status code is between 200 and 299 (inclusive).
//   - `false` if the status code is outside of this range.
func (w *wrapper) IsSuccess() bool {
	return (200 <= w.statusCode) && (w.statusCode <= 299)
}

// IsRedirection checks whether the HTTP status code indicates a redirection response.
//
// This function checks if the `statusCode` is between 300 and 399, inclusive, which indicates a redirection HTTP response.
//
// Returns:
//   - A boolean value indicating whether the HTTP response is a redirection:
//   - `true` if the status code is between 300 and 399 (inclusive).
//   - `false` if the status code is outside of this range.
func (w *wrapper) IsRedirection() bool {
	return (300 <= w.statusCode) && (w.statusCode <= 399)
}

// IsClientError checks whether the HTTP status code indicates a client error.
//
// This function checks if the `statusCode` is between 400 and 499, inclusive, which indicates a client error HTTP response.
//
// Returns:
//   - A boolean value indicating whether the HTTP response is a client error:
//   - `true` if the status code is between 400 and 499 (inclusive).
//   - `false` if the status code is outside of this range.
func (w *wrapper) IsClientError() bool {
	return (400 <= w.statusCode) && (w.statusCode <= 499)
}

// IsServerError checks whether the HTTP status code indicates a server error.
//
// This function checks if the `statusCode` is between 500 and 599, inclusive, which indicates a server error HTTP response.
//
// Returns:
//   - A boolean value indicating whether the HTTP response is a server error:
//   - `true` if the status code is between 500 and 599 (inclusive).
//   - `false` if the status code is outside of this range.
func (w *wrapper) IsServerError() bool {
	return (500 <= w.statusCode) && (w.statusCode <= 599)
}

// IsLast checks whether the current pagination represents the last page.
//
// This function checks the `isLast` field of the `pagination` instance to determine if the current page is the last one.
// The `isLast` field is typically set to `true` when there are no more pages of data available.
//
// Returns:
//   - A boolean value indicating whether the current page is the last:
//   - `true` if `isLast` is `true`, indicating this is the last page of results.
//   - `false` if `isLast` is `false`, indicating more pages are available.
func (p *pagination) IsLast() bool {
	return p.isLast
}

// IsLastPage checks whether the current page is the last page of results.
//
// This function verifies that pagination information is present and then checks if the current page is the last page.
// It combines the checks of `IsPagingPresent()` and `IsLast()` to ensure that the pagination structure exists
// and that it represents the last page.
//
// Returns:
//   - A boolean value indicating whether the current page is the last page:
//   - `true` if pagination is present and the current page is the last one.
//   - `false` if pagination is not present or the current page is not the last.
func (w *wrapper) IsLastPage() bool {
	return w.IsPagingPresent() && w.pagination.IsLast()
}
