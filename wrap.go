package wrapify

import (
	"fmt"
	"net/http"
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
	return w.Cause().Error()
}

// Cause traverses the error chain and returns the underlying cause of the error
// associated with the `wrapper` instance.
//
// This function checks if the error stored in the `wrapper` is itself another
// `wrapper` instance. If so, it recursively calls `Cause` on the inner error
// to find the ultimate cause. Otherwise, it returns the current error.
//
// Returns:
//   - The underlying cause of the error, which can be another error or the original error.
func (w *wrapper) Cause() error {
	// Traverse through wrapped errors.
	// We will use Unwrap() method to unwrap errors instead of checking for *wrapper explicitly.
	// This way, we can traverse to the innermost cause regardless of error type.
	cause := w.errors
	for cause != nil {
		// If the error has a Cause method, we use it to find the underlying error.
		if err, ok := cause.(interface{ Cause() error }); ok {
			cause = err.Cause()
		} else {
			// No Cause() method, so return the error itself.
			break
		}
	}
	return cause
}

// StatusCode retrieves the HTTP status code associated with the `wrapper` instance.
//
// This function returns the `statusCode` field of the `wrapper`, which represents
// the HTTP status code for the response, indicating the outcome of the request.
//
// Returns:
//   - An integer representing the HTTP status code.
func (w *wrapper) StatusCode() int {
	if !w.Available() {
		return 0
	}
	return w.statusCode
}

// StatusText returns a human-readable string representation of the HTTP status.
//
// This function combines the status code with its associated status text, which
// is retrieved using the `http.StatusText` function from the `net/http` package.
// The returned string follows the format "statusCode (statusText)".
//
// For example, if the status code is 200, the function will return "200 (OK)".
// If the status code is 404, it will return "404 (Not Found)".
//
// Returns:
//   - A string formatted as "statusCode (statusText)", where `statusCode` is the
//     numeric HTTP status code and `statusText` is the corresponding textual description.
func (w *wrapper) StatusText() string {
	return fmt.Sprintf("%d (%s)", w.StatusCode(), http.StatusText(w.StatusCode()))
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

// Header retrieves the `header` associated with the `wrapper` instance.
//
// This function returns the `header` field from the `wrapper` instance, which contains
// information about the HTTP response or any other relevant metadata. If the `wrapper`
// instance is correctly initialized, it will return the `header`; otherwise, it may
// return `nil` if the `header` has not been set.
//
// Returns:
//   - A pointer to the `header` instance associated with the `wrapper`.
//   - `nil` if the `header` is not set or the `wrapper` is uninitialized.
func (w *wrapper) Header() *header {
	return w.header
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
	s, ok := w.data.(string)
	if ok {
		return w.Available() &&
			unify4g.IsNotEmpty(s) &&
			unify4g.IsNotBlank(s)
	}
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

// ApiVersion retrieves the API version from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before attempting
// to retrieve the `apiVersion`. If the `meta` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The API version as a string if available.
//   - An empty string if the `meta` instance is unavailable.
func (m *meta) ApiVersion() string {
	if !m.Available() {
		return ""
	}
	return m.apiVersion
}

// RequestID retrieves the request ID from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before retrieving
// the `requestID`. If the `meta` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The request ID as a string if available.
//   - An empty string if the `meta` instance is unavailable.
func (m *meta) RequestID() string {
	if !m.Available() {
		return ""
	}
	if unify4g.IsEmpty(m.requestID) || unify4g.IsBlank(m.requestID) {
		m.WithRequestID(unify4g.GenerateCryptoID())
	}
	return m.requestID
}

// Locale retrieves the locale from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before retrieving
// the `locale`. If the `meta` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The locale as a string if available.
//   - An empty string if the `meta` instance is unavailable.
func (m *meta) Locale() string {
	if !m.Available() {
		return ""
	}
	return m.locale
}

// RequestedTime retrieves the requested time from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before retrieving
// the `requestedTime`. If the `meta` instance is unavailable, it returns the zero value
// of `time.Time` (January 1, year 1, 00:00:00 UTC).
//
// Returns:
//   - The requested time as a `time.Time` object if available.
//   - The zero value of `time.Time` if the `meta` instance is unavailable.
func (m *meta) RequestedTime() time.Time {
	if !m.Available() {
		return time.Time{}
	}
	return m.requestedTime
}

// CustomFields retrieves the custom fields from the `meta` instance.
//
// This function checks if the `meta` instance is available (non-nil) before retrieving
// the `customFields`. If the `meta` instance is unavailable, it returns `nil`.
//
// Returns:
//   - A map of custom fields if available.
//   - `nil` if the `meta` instance is unavailable.
func (m *meta) CustomFields() map[string]interface{} {
	if !m.Available() {
		return nil
	}
	return m.customFields
}

// Available checks if the `header` instance is non-nil.
//
// This function ensures that the `header` instance is not nil before performing any operations.
// It returns `true` if the `header` is non-nil, and `false` if the `header` is nil.
//
// Returns:
//   - `true` if the `header` instance is not nil.
//   - `false` if the `header` instance is nil.
func (h *header) Available() bool {
	return h != nil
}

// IsCodePresent checks if the `code` field in the `header` instance is present and greater than zero.
//
// This function first checks if the `header` is available (non-nil), and then checks if the `code`
// field is greater than zero, indicating that it is present and valid.
//
// Returns:
//   - `true` if the `code` field is greater than zero.
//   - `false` if the `code` field is either not present (nil) or zero.
func (h *header) IsCodePresent() bool {
	return h.Available() && h.code > 0
}

// IsTextPresent checks if the `text` field in the `header` instance is present and not empty.
//
// This function verifies if the `header` is available and if the `text` field is not empty, using
// the `unify4g.IsNotEmpty` utility to ensure the presence of the `text` field.
//
// Returns:
//   - `true` if the `text` field is non-empty.
//   - `false` if the `text` field is either not present (nil) or empty.
func (h *header) IsTextPresent() bool {
	return h.Available() && unify4g.IsNotEmpty(h.text)
}

// IsTypePresent checks if the `Type` field in the `header` instance is present and not empty.
//
// This function checks if the `header` instance is available and if the `Type` field is not empty,
// utilizing the `unify4g.IsNotEmpty` utility to determine whether the `Type` field contains a value.
//
// Returns:
//   - `true` if the `Type` field is non-empty.
//   - `false` if the `Type` field is either not present (nil) or empty.
func (h *header) IsTypePresent() bool {
	return h.Available() && unify4g.IsNotEmpty(h.typez)
}

// IsDescriptionPresent checks if the `description` field in the `header` instance is present and not empty.
//
// This function ensures that the `header` is available and that the `description` field is not empty,
// using `unify4g.IsNotEmpty` to check for non-emptiness.
//
// Returns:
//   - `true` if the `description` field is non-empty.
//   - `false` if the `description` field is either not present (nil) or empty.
func (h *header) IsDescriptionPresent() bool {
	return h.Available() && unify4g.IsNotEmpty(h.description)
}

// Code retrieves the code value from the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) before retrieving
// the `code` field. If the `header` instance is unavailable, it returns 0.
//
// Returns:
//   - The `code` as an integer if available.
//   - 0 if the `header` instance is unavailable.
func (h *header) Code() int {
	if !h.Available() {
		return 0
	}
	return h.code
}

// Text retrieves the text value from the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) before retrieving
// the `text` field. If the `header` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The `text` as a string if available.
//   - An empty string if the `header` instance is unavailable.
func (h *header) Text() string {
	if !h.Available() {
		return ""
	}
	return h.text
}

// Type retrieves the type value from the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) before retrieving
// the `Type` field. If the `header` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The `Type` as a string if available.
//   - An empty string if the `header` instance is unavailable.
func (h *header) Type() string {
	if !h.Available() {
		return ""
	}
	return h.typez
}

// Description retrieves the description value from the `header` instance.
//
// This function checks if the `header` instance is available (non-nil) before retrieving
// the `description` field. If the `header` instance is unavailable, it returns an empty string.
//
// Returns:
//   - The `description` as a string if available.
//   - An empty string if the `header` instance is unavailable.
func (h *header) Description() string {
	if !h.Available() {
		return ""
	}
	return h.description
}

// WrapOk creates a wrapper for a successful HTTP response (200 OK).
//
// This function sets the HTTP status code to 200 (OK) and includes a message and data payload
// in the response body.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapOk(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusOK).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapCreated creates a wrapper for a resource creation response (201 WrapCreated).
//
// This function sets the HTTP status code to 201 (WrapCreated) and includes a message and data payload
// in the response body.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapCreated(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusCreated).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapBadRequest creates a wrapper for a client error response (400 Bad Request).
//
// This function sets the HTTP status code to 400 (Bad Request) and includes a message and data payload
// in the response body.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapBadRequest(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusBadRequest).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapNotFound creates a wrapper for a resource not found response (404 Not Found).
//
// This function sets the HTTP status code to 404 (Not Found) and includes a message and data payload
// in the response body.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapNotFound(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusNotFound).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapNotImplemented creates a wrapper for a response indicating unimplemented functionality (501 Not Implemented).
//
// This function sets the HTTP status code to 501 (Not Implemented) and includes a message and data payload
// in the response body.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapNotImplemented(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusNotImplemented).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapTooManyRequest creates a wrapper for a rate-limiting response (429 Too Many Requests).
//
// This function sets the HTTP status code to 429 (Too Many Requests) and includes a message and data payload
// in the response body.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapTooManyRequest(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusTooManyRequests).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapLocked creates a wrapper for a locked resource response (423 WrapLocked).
//
// This function sets the HTTP status code to 423 (WrapLocked) and includes a message and data payload
// in the response body.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapLocked(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusLocked).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapNoContent creates a wrapper for a successful response without a body (204 No Content).
//
// This function sets the HTTP status code to 204 (No Content) and includes a message and data payload
// in the response body.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapNoContent(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusNoContent).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapProcessing creates a wrapper for a response indicating ongoing processing (102 WrapProcessing).
//
// This function sets the HTTP status code to 102 (WrapProcessing) and includes a message and data payload
// in the response body.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapProcessing(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusProcessing).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapUpgradeRequired creates a wrapper for a response indicating an upgrade is required (426 Upgrade Required).
//
// This function sets the HTTP status code to 426 (Upgrade Required) and includes a message and data payload
// in the response body. It is typically used when the client must switch to a different protocol.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapUpgradeRequired(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusUpgradeRequired).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapServiceUnavailable creates a wrapper for a response indicating the service is temporarily unavailable (503 Service Unavailable).
//
// This function sets the HTTP status code to 503 (Service Unavailable) and includes a message and data payload
// in the response body. It is typically used when the server is unable to handle the request due to temporary overload
// or maintenance.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapServiceUnavailable(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusServiceUnavailable).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapInternalServerError creates a wrapper for a server error response (500 Internal Server Error).
//
// This function sets the HTTP status code to 500 (Internal Server Error) and includes a message and data payload
// in the response body. It is typically used when the server encounters an unexpected condition that prevents it
// from fulfilling the request.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapInternalServerError(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusInternalServerError).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapGatewayTimeout creates a wrapper for a response indicating a gateway timeout (504 Gateway Timeout).
//
// This function sets the HTTP status code to 504 (Gateway Timeout) and includes a message and data payload
// in the response body. It is typically used when the server did not receive a timely response from an upstream server.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapGatewayTimeout(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusGatewayTimeout).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapMethodNotAllowed creates a wrapper for a response indicating the HTTP method is not allowed (405 Method Not Allowed).
//
// This function sets the HTTP status code to 405 (Method Not Allowed) and includes a message and data payload
// in the response body. It is typically used when the server knows the method is not supported for the target resource.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapMethodNotAllowed(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusMethodNotAllowed).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapUnauthorized creates a wrapper for a response indicating authentication is required (401 WrapUnauthorized).
//
// This function sets the HTTP status code to 401 (WrapUnauthorized) and includes a message and data payload
// in the response body. It is typically used when the request has not been applied because it lacks valid
// authentication credentials.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapUnauthorized(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusUnauthorized).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapForbidden creates a wrapper for a response indicating access to the resource is forbidden (403 WrapForbidden).
//
// This function sets the HTTP status code to 403 (WrapForbidden) and includes a message and data payload
// in the response body. It is typically used when the server understands the request but refuses to authorize it.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapForbidden(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusForbidden).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapAccepted creates a wrapper for a response indicating the request has been accepted for processing (202 WrapAccepted).
//
// This function sets the HTTP status code to 202 (WrapAccepted) and includes a message and data payload
// in the response body. It is typically used when the request has been received but processing is not yet complete.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapAccepted(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusAccepted).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapRequestTimeout creates a wrapper for a response indicating the client request has timed out (408 Request Timeout).
//
// This function sets the HTTP status code to 408 (Request Timeout) and includes a message and data payload
// in the response body. It is typically used when the server did not receive a complete request message within the time it was prepared to wait.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapRequestTimeout(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusRequestTimeout).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapRequestEntityTooLarge creates a wrapper for a response indicating the request entity is too large (413 Payload Too Large).
//
// This function sets the HTTP status code to 413 (Payload Too Large) and includes a message and data payload
// in the response body. It is typically used when the server refuses to process a request because the request entity is larger than the server is willing or able to process.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapRequestEntityTooLarge(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusRequestEntityTooLarge).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapUnsupportedMediaType creates a wrapper for a response indicating the media type is not supported (415 Unsupported Media Type).
//
// This function sets the HTTP status code to 415 (Unsupported Media Type) and includes a message and data payload
// in the response body. It is typically used when the server refuses to accept the request because the payload is in an unsupported format.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapUnsupportedMediaType(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusUnsupportedMediaType).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapHTTPVersionNotSupported creates a wrapper for a response indicating the HTTP version is not supported (505 HTTP Version Not Supported).
//
// This function sets the HTTP status code to 505 (HTTP Version Not Supported) and includes a message and data payload
// in the response body. It is typically used when the server does not support the HTTP protocol version used in the request.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapHTTPVersionNotSupported(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusHTTPVersionNotSupported).
		WithMessage(message).
		WithBody(data)
	return w
}

// WrapPaymentRequired creates a wrapper for a response indicating payment is required (402 Payment Required).
//
// This function sets the HTTP status code to 402 (Payment Required) and includes a message and data payload
// in the response body. It is typically used when access to the requested resource requires payment.
//
// Parameters:
//   - message: A string containing the response message.
//   - data: The data payload to include in the response.
//
// Returns:
//   - A pointer to a `wrapper` instance representing the response.
func WrapPaymentRequired(message string, data interface{}) *wrapper {
	w := New().
		WithStatusCode(http.StatusPaymentRequired).
		WithMessage(message).
		WithBody(data)
	return w
}
