package wrapify

// Error retrieves the error associated with the `wrapper` instance.
//
// This function returns the `errors` field of the `wrapper`, which contains
// any errors encountered during the operation of the `wrapper`.
//
// Returns:
//   - An error object, or `nil` if no errors are present.
func (w *wrapper) Error() error {
	return w.errors
}

// StateCode retrieves the HTTP status code associated with the `wrapper` instance.
//
// This function returns the `statusCode` field of the `wrapper`, which represents
// the HTTP status code for the response, indicating the outcome of the request.
//
// Returns:
//   - An integer representing the HTTP status code.
func (w *wrapper) StateCode() int {
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
	return w.data
}
