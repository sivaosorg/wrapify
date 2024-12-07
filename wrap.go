package wrapify

import (
	"time"

	"github.com/sivaosorg/unify4g"
)

// Available checks whether the `wrapper` instance is non-nil.
//
// This function ensures that the `wrapper` object exists and is not nil.
// It serves as a safety check to avoid null pointer dereferences when accessing the instance's fields or methods.
//
// Returns:
//   - A boolean value indicating whether the `wrapper` instance is non-nil:
//   - `true` if the `wrapper` instance is non-nil.
//   - `false` if the `wrapper` instance is nil.
func (w *wrapper) Available() bool {
	return w != nil
}

// Error retrieves the error associated with the `wrapper` instance.
//
// This function returns the `errors` field of the `wrapper`, which contains
// any errors encountered during the operation of the `wrapper`.
//
// Returns:
//   - An error object, or `nil` if no errors are present.
func (w *wrapper) Error() string {
	if !w.Available() {
		return ""
	}
	return w.errors.Error()
}

// StateCode retrieves the HTTP status code associated with the `wrapper` instance.
//
// This function returns the `statusCode` field of the `wrapper`, which represents
// the HTTP status code for the response, indicating the outcome of the request.
//
// Returns:
//   - An integer representing the HTTP status code.
func (w *wrapper) StateCode() int {
	if !w.Available() {
		return 0
	}
	return w.statusCode
}

// Message retrieves the message associated with the `wrapper` instance.
//
// This function returns the `message` field of the `wrapper`, which typically
// provides additional context or a description of the operation's outcome.
//
// Returns:
//   - A string representing the message.
func (w *wrapper) Message() string {
	if !w.Available() {
		return ""
	}
	return w.message
}

// Total retrieves the total number of items associated with the `wrapper` instance.
//
// This function returns the `total` field of the `wrapper`, which indicates
// the total number of items available, often used in paginated responses.
//
// Returns:
//   - An integer representing the total number of items.
func (w *wrapper) Total() int {
	if !w.Available() {
		return 0
	}
	return w.total
}

// Body retrieves the body data associated with the `wrapper` instance.
//
// This function returns the `data` field of the `wrapper`, which contains
// the primary data payload of the response.
//
// Returns:
//   - The body data (of any type), or `nil` if no body data is present.
func (w *wrapper) Body() interface{} {
	if !w.Available() {
		return nil
	}
	return w.data
}

// Debugging retrieves the debugging information from the `wrapper` instance.
//
// This function checks if the `wrapper` instance is available (non-nil) before returning
// the value of the `debug` field. If the `wrapper` is not available, it returns an
// empty map to ensure safe usage.
//
// Returns:
//   - A `map[string]interface{}` containing the debugging information.
//   - An empty map if the `wrapper` instance is not available.
func (w *wrapper) Debugging() map[string]interface{} {
	if !w.Available() {
		return nil
	}
	return w.debug
}

// OnKeyDebugging retrieves the value of a specific debugging key from the `wrapper` instance.
//
// This function checks if the `wrapper` is available (non-nil) and if the specified debugging key
// is present in the `debug` map. If both conditions are met, it returns the value associated with
// the specified key. Otherwise, it returns `nil` to indicate the key is not available.
//
// Parameters:
//   - `key`: A string representing the debugging key to retrieve.
//
// Returns:
//   - The value associated with the specified debugging key if it exists.
//   - `nil` if the `wrapper` is unavailable or the key is not present in the `debug` map.
func (w *wrapper) OnKeyDebugging(key string) interface{} {
	if !w.Available() || !w.IsDebuggingKeyPresent(key) {
		return nil
	}
	return w.debug[key]
}

// Pagination retrieves the `pagination` instance associated with the `wrapper`.
//
// This function returns the `pagination` field of the `wrapper`, allowing access to
// pagination details such as the current page, total pages, and total items. If no
// pagination information is available, it returns `nil`.
//
// Returns:
//   - A pointer to the `pagination` instance if available.
//   - `nil` if the `pagination` field is not set.
func (w *wrapper) Pagination() *pagination {
	return w.pagination
}

// Meta retrieves the `meta` information from the `wrapper` instance.
//
// This function returns the `meta` field, which contains metadata related to the response or data
// in the `wrapper` instance. If no `meta` information is set, it returns `nil`.
//
// Returns:
//   - A pointer to the `meta` instance associated with the `wrapper`.
//   - `nil` if no `meta` information is available.
func (w *wrapper) Meta() *meta {
	return w.meta
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
	return w.Available() && w.debug != nil && len(w.debug) > 0
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
	return w.Available() && w.data != nil
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
	return w.Available() && w.header != nil
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
	return w.Available() && w.meta != nil
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
	return w.Available() && w.pagination != nil
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
	return w.Available() && w.errors != nil
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
	return w.Available() && w.total >= 0
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
	return w.Available() && w.statusCode > 0
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
	return w.Available() && (200 <= w.statusCode) && (w.statusCode <= 299)
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
	return w.Available() && (300 <= w.statusCode) && (w.statusCode <= 399)
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
	return w.Available() && (400 <= w.statusCode) && (w.statusCode <= 499)
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
	return w.Available() && (500 <= w.statusCode) && (w.statusCode <= 599)
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
	return w.Available() && w.IsPagingPresent() && w.pagination.IsLast()
}

// Available checks whether the `pagination` instance is non-nil.
//
// This function ensures that the `pagination` object exists and is not nil.
// It serves as a safety check to avoid null pointer dereferences when accessing the instance's fields or methods.
//
// Returns:
//   - A boolean value indicating whether the `pagination` instance is non-nil:
//   - `true` if the `pagination` instance is non-nil.
//   - `false` if the `pagination` instance is nil.
func (p *pagination) Available() bool {
	return p != nil
}

// Page retrieves the current page number from the `pagination` instance.
//
// This function checks if the `pagination` instance is available (non-nil) before
// returning the value of the `page` field. If the instance is not available, it
// returns a default value of `0`.
//
// Returns:
//   - An integer representing the current page number.
//   - `0` if the `pagination` instance is not available.
func (p *pagination) Page() int {
	if !p.Available() {
		return 0
	}
	return p.page
}

// PerPage retrieves the number of items per page from the `pagination` instance.
//
// This function checks if the `pagination` instance is available (non-nil) before
// returning the value of the `perPage` field. If the instance is not available, it
// returns a default value of `0`.
//
// Returns:
//   - An integer representing the number of items per page.
//   - `0` if the `pagination` instance is not available.
func (p *pagination) PerPage() int {
	if !p.Available() {
		return 0
	}
	return p.perPage
}

// TotalPages retrieves the total number of pages from the `pagination` instance.
//
// This function checks if the `pagination` instance is available (non-nil) before
// returning the value of the `totalPages` field. If the instance is not available, it
// returns a default value of `0`.
//
// Returns:
//   - An integer representing the total number of pages.
//   - `0` if the `pagination` instance is not available.
func (p *pagination) TotalPages() int {
	if !p.Available() {
		return 0
	}
	return p.totalPages
}

// TotalItems retrieves the total number of items from the `pagination` instance.
//
// This function checks if the `pagination` instance is available (non-nil) before
// returning the value of the `totalItems` field. If the instance is not available, it
// returns a default value of `0`.
//
// Returns:
//   - An integer representing the total number of items.
//   - `0` if the `pagination` instance is not available.
func (p *pagination) TotalItems() int {
	if !p.Available() {
		return 0
	}
	return p.totalItems
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
	if !p.Available() {
		return true
	}
	return p.isLast
}

// Available checks whether the `meta` instance is available (non-nil).
//
// This function ensures that the `meta` instance is valid and not nil
// to safely access its fields and methods.
//
// Returns:
//   - `true` if the `meta` instance is not nil.
//   - `false` if the `meta` instance is nil.
func (m *meta) Available() bool {
	return m != nil
}

// IsApiVersionPresent checks whether the API version is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `apiVersion` field contains a non-empty value.
//
// Returns:
//   - `true` if `apiVersion` is non-empty.
//   - `false` if `meta` is unavailable or `apiVersion` is empty.
func (m *meta) IsApiVersionPresent() bool {
	return m.Available() && unify4g.IsNotEmpty(m.apiVersion)
}

// IsRequestIDPresent checks whether the request ID is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `requestID` field contains a non-empty value.
//
// Returns:
//   - `true` if `requestID` is non-empty.
//   - `false` if `meta` is unavailable or `requestID` is empty.
func (m *meta) IsRequestIDPresent() bool {
	return m.Available() && unify4g.IsNotEmpty(m.requestID)
}

// IsLocalePresent checks whether the locale information is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `locale` field contains a non-empty value.
//
// Returns:
//   - `true` if `locale` is non-empty.
//   - `false` if `meta` is unavailable or `locale` is empty.
func (m *meta) IsLocalePresent() bool {
	return m.Available() && unify4g.IsNotEmpty(m.locale)
}

// IsRequestedTimePresent checks whether the requested time is present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `requestedTime` field is not the zero value of `time.Time`.
//
// Returns:
//   - `true` if `requestedTime` is not the zero value of `time.Time`.
//   - `false` if `meta` is unavailable or `requestedTime` is uninitialized.
func (m *meta) IsRequestedTimePresent() bool {
	return m.Available() && m.requestedTime != time.Time{}
}

// IsCustomFieldPresent checks whether custom fields are present in the `meta` instance.
//
// This function verifies that the `meta` instance is available and that
// the `customFields` field is non-nil and contains at least one entry.
//
// Returns:
//   - `true` if `customFields` is non-nil and has a non-zero length.
//   - `false` if `meta` is unavailable, `customFields` is nil, or it is empty.
func (m *meta) IsCustomFieldPresent() bool {
	return m.Available() && m.customFields != nil && len(m.customFields) > 0
}

// IsCustomFieldKeyPresent checks whether a specific key is present in the custom fields of the `meta` instance.
//
// This function first verifies that the `customFields` field is available and contains data using
// `IsCustomFieldPresent`. If so, it checks if the specified key exists within the `customFields` map.
//
// Parameters:
//   - `key`: A string representing the key to search for in the `customFields` map.
//
// Returns:
//   - `true` if the `customFields` map is available and contains the specified key.
//   - `false` if `customFields` is nil, empty, or does not contain the specified key.
func (m *meta) IsCustomFieldKeyPresent(key string) bool {
	return m.IsCustomFieldPresent() && unify4g.MapContainsKey(m.customFields, key)
}

// OnKeyCustomField retrieves the value associated with a specific key in the custom fields of the `meta` instance.
//
// This function checks whether the `meta` instance is available and whether the specified key exists
// in the `customFields` map. If both conditions are met, it returns the corresponding value for the key.
// If the `meta` instance is unavailable or the key is not present, it returns `nil`.
//
// Parameters:
//   - `key`: A string representing the key to search for in the `customFields` map.
//
// Returns:
//   - The value associated with the specified key in the `customFields` map if the key is present.
//   - `nil` if the `meta` instance is unavailable or the key does not exist in the `customFields` map.
func (m *meta) OnKeyCustomField(key string) interface{} {
	if !m.Available() || !m.IsCustomFieldKeyPresent(key) {
		return nil
	}
	return m.customFields[key]
}
