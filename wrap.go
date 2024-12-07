package wrapify

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
		return map[string]interface{}{}
	}
	return w.debug
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
