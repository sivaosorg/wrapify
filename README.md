# wrapify

**wrapify** is a Go library designed to simplify and standardize API response wrapping for RESTful services. It leverages the Decorator Pattern to dynamically add error handling, metadata, pagination, and other response features in a clean and human-readable format. With Wrapify, you can ensure consistent and extensible API responses with minimal boilerplate. Perfect for building robust, maintainable REST APIs in Go!

### Requirements

- Go version 1.23 or higher

### Installation

To install, you can use the following commands based on your preference:

- For a specific version:

  ```bash
  go get github.com/sivaosorg/wrapify@v0.0.1
  ```

- For the latest version:
  ```bash
  go get -u github.com/sivaosorg/wrapify@latest
  ```

### Getting started

#### Getting wrapify

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```go
import "github.com/sivaosorg/wrapify"
```

#### Usage

> Structure of the wrapper-based standardized API response

```go
// R represents a wrapper around the main `wrapper` struct.
// It is used as a high-level abstraction to provide a simplified interface for handling API responses.
// The `R` type allows for easier manipulation of the wrapped data, metadata,
// and other response components, while maintaining the flexibility of the underlying `wrapper` structure.
type R struct {
    *wrapper
}
// Available checks whether the `wrapper` instance is non-nil.
func (w *wrapper) Available() bool
// Body retrieves the body data associated with the `wrapper` instance.
func (w *wrapper) Body() interface{}
// Cause traverses the error chain and returns the underlying cause of the error
// associated with the `wrapper` instance.
func (w *wrapper) Cause() error
// Debugging retrieves the debugging information from the `wrapper` instance.
func (w *wrapper) Debugging() map[string]interface{}
// Error retrieves the error associated with the `wrapper` instance.
func (w *wrapper) Error() string
// Header retrieves the `header` associated with the `wrapper` instance.
func (w *wrapper) Header() *header
// IsBodyPresent checks whether the body data is present in the `wrapper` instance.
func (w *wrapper) IsBodyPresent() bool
// IsClientError checks whether the HTTP status code indicates a client error.
func (w *wrapper) IsClientError() bool
// IsDebuggingKeyPresent checks whether a specific key exists in the `debug` information.
func (w *wrapper) IsDebuggingKeyPresent(key string) bool
// IsDebuggingPresent checks whether debugging information is present in the `wrapper` instance.
func (w *wrapper) IsDebuggingPresent() bool
// IsError checks whether there is an error present in the `wrapper` instance.
// This function returns `true` if the `wrapper` contains an error, which can be any of the following:
//   - An error present in the `errors` field.
//   - A client error (4xx status code) or a server error (5xx status code).
func (w *wrapper) IsError() bool
// IsErrorPresent checks whether an error is present in the `wrapper` instance.
func (w *wrapper) IsErrorPresent() bool
// IsHeaderPresent checks whether header information is present in the `wrapper` instance.
func (w *wrapper) IsHeaderPresent() bool
// IsLastPage checks whether the current page is the last page of results.
func (w *wrapper) IsLastPage() bool
// IsMetaPresent checks whether metadata information is present in the `wrapper` instance.
func (w *wrapper) IsMetaPresent() bool
// IsPagingPresent checks whether pagination information is present in the `wrapper` instance.
func (w *wrapper) IsPagingPresent() bool
// IsRedirection checks whether the HTTP status code indicates a redirection response.
func (w *wrapper) IsRedirection() bool
// IsServerError checks whether the HTTP status code indicates a server error.
func (w *wrapper) IsServerError() bool
// IsStatusCodePresent checks whether a valid status code is present in the `wrapper` instance.
func (w *wrapper) IsStatusCodePresent() bool
// IsSuccess checks whether the HTTP status code indicates a successful response.
func (w *wrapper) IsSuccess() bool
// IsTotalPresent checks whether the total number of items is present in the `wrapper` instance.
func (w *wrapper) IsTotalPresent() bool
// Json serializes the `wrapper` instance into a compact JSON string.
func (w *wrapper) Json() string
// JsonPretty serializes the `wrapper` instance into a prettified JSON string.
func (w *wrapper) JsonPretty() string
// Message retrieves the message associated with the `wrapper` instance.
func (w *wrapper) Message() string
// Meta retrieves the `meta` information from the `wrapper` instance.
func (w *wrapper) Meta() *meta
// OnKeyDebugging retrieves the value of a specific debugging key from the `wrapper` instance.
func (w *wrapper) OnKeyDebugging(key string) interface{}
// Pagination retrieves the `pagination` instance associated with the `wrapper`.
func (w *wrapper) Pagination() *pagination
// R represents a wrapper around the main `wrapper` struct. It is used as a high-level
// abstraction to provide a simplified interface for handling API responses.
// The `R` type allows for easier manipulation of the wrapped data, metadata, and other
// response components, while maintaining the flexibility of the underlying `wrapper` structure.
func (w *wrapper) Reply() R
// Respond generates a map representation of the `wrapper` instance.
func (w *wrapper) Respond() map[string]interface{}
// StatusCode retrieves the HTTP status code associated with the `wrapper` instance.
func (w *wrapper) StatusCode() int
// StatusText returns a human-readable string representation of the HTTP status.
func (w *wrapper) StatusText() string
// Total retrieves the total number of items associated with the `wrapper` instance.
func (w *wrapper) Total() int
// WithApiVersion sets the API version in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithApiVersion(v string) *wrapper
// WithApiVersionf sets the API version in the `meta` field of the `wrapper` instance using a formatted string.
func (w *wrapper) WithApiVersionf(format string, args ...interface{}) *wrapper
// WithBody sets the body data for the `wrapper` instance.
func (w *wrapper) WithBody(v interface{}) *wrapper
// WithCustomFieldKV sets a specific custom field key-value pair in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithCustomFieldKV(key string, value interface{}) *wrapper
// WithCustomFieldKVf sets a specific custom field key-value pair in the `meta` field of the `wrapper` instance
// using a formatted value.
func (w *wrapper) WithCustomFieldKVf(key string, format string, args ...interface{}) *wrapper
// WithCustomFields sets the custom fields in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithCustomFields(values map[string]interface{}) *wrapper
// WithDebugging sets the debugging information for the `wrapper` instance.
func (w *wrapper) WithDebugging(v map[string]interface{}) *wrapper
// WithDebuggingKV adds a key-value pair to the debugging information in the `wrapper` instance.
func (w *wrapper) WithDebuggingKV(key string, value interface{}) *wrapper
// WithDebuggingKVf adds a formatted key-value pair to the debugging information in the `wrapper` instance.
func (w *wrapper) WithDebuggingKVf(key string, format string, args ...interface{}) *wrapper
// WithErrMessage adds a plain contextual message to an existing error and sets it for the `wrapper` instance.
func (w *wrapper) WithErrMessage(err error, message string) *wrapper
// WithErrMessagef adds a formatted contextual message to an existing error and sets it for the `wrapper` instance.
func (w *wrapper) WithErrMessagef(err error, format string, args ...interface{}) *wrapper
// WithErrSck sets an error with a stack trace for the `wrapper` instance.
func (w *wrapper) WithErrSck(err error) *wrapper
// WithErrWrap wraps an existing error with an additional message and sets it for the `wrapper` instance.
func (w *wrapper) WithErrWrap(err error, message string) *wrapper
// WithErrWrapf wraps an existing error with a formatted message and sets it for the `wrapper` instance.
func (w *wrapper) WithErrWrapf(err error, format string, args ...interface{}) *wrapper
// WithError sets an error for the `wrapper` instance using a plain error message.
func (w *wrapper) WithError(message string) *wrapper
// WithErrorf sets a formatted error for the `wrapper` instance.
func (w *wrapper) WithErrorf(format string, args ...interface{}) *wrapper
// WithHeader sets the header for the `wrapper` instance.
func (w *wrapper) WithHeader(v *header) *wrapper
// WithIsLast sets whether the current page is the last one in the wrapper's pagination.
func (w *wrapper) WithIsLast(v bool) *wrapper
// WithLocale sets the locale in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithLocale(v string) *wrapper
// WithMessage sets a message for the `wrapper` instance.
func (w *wrapper) WithMessage(message string) *wrapper
// WithMessagef sets a formatted message for the `wrapper` instance.
func (w *wrapper) WithMessagef(message string, args ...interface{}) *wrapper
// WithMeta sets the metadata for the `wrapper` instance.
func (w *wrapper) WithMeta(v *meta) *wrapper
// WithPage sets the current page number in the wrapper's pagination.
func (w *wrapper) WithPage(v int) *wrapper
// WithPagination sets the pagination information for the `wrapper` instance.
func (w *wrapper) WithPagination(v *pagination) *wrapper
// WithPath sets the request path for the `wrapper` instance.
func (w *wrapper) WithPath(v string) *wrapper
// WithPathf sets a formatted request path for the `wrapper` instance.
func (w *wrapper) WithPathf(v string, args ...interface{}) *wrapper
// WithPerPage sets the number of items per page in the wrapper's pagination.
func (w *wrapper) WithPerPage(v int) *wrapper
// WithRequestID sets the request ID in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithRequestID(v string) *wrapper
// WithRequestIDf sets the request ID in the `meta` field of the `wrapper` instance using a formatted string.
func (w *wrapper) WithRequestIDf(format string, args ...interface{}) *wrapper
// WithRequestedTime sets the requested time in the `meta` field of the `wrapper` instance.
func (w *wrapper) WithRequestedTime(v time.Time) *wrapper
// WithStatusCode sets the HTTP status code for the `wrapper` instance.
func (w *wrapper) WithStatusCode(code int) *wrapper
// WithTotal sets the total number of items for the `wrapper` instance.
func (w *wrapper) WithTotal(total int) *wrapper
// WithTotalItems sets the total number of items in the wrapper's pagination.
func (w *wrapper) WithTotalItems(v int) *wrapper
// WithTotalPages sets the total number of pages in the wrapper's pagination.
func (w *wrapper) WithTotalPages(v int) *wrapper
```

### Contributing

To contribute to project, follow these steps:

1. Clone the repository:

   ```bash
   git clone --depth 1 https://github.com/sivaosorg/wrapify.git
   ```

2. Navigate to the project directory:

   ```bash
   cd wrapify
   ```

3. Prepare the project environment:
   ```bash
   go mod tidy
   ```
