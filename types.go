package wrapify

import (
	"sync"
	"time"
)

// R represents a wrapper around the main `wrapper` struct. It is used as a high-level
// abstraction to provide a simplified interface for handling API responses.
// The `R` type allows for easier manipulation of the wrapped data, metadata, and other
// response components, while maintaining the flexibility of the underlying `wrapper` structure.
type R struct {
	*wrapper
}

// Frame represents a program counter inside a stack frame.
// A `Frame` is essentially a single point in the stack trace,
// representing a program counter (the location in code) at the
// time of a function call. Historically, for compatibility reasons,
// a `Frame` is interpreted as a `uintptr`, but the value stored
// in the `Frame` represents the program counter + 1. This allows
// for distinguishing between an invalid program counter and a valid one.
//
// A `Frame` is typically used within a `StackTrace` to track the
// sequence of function calls leading to the current point in the program.
// A frame is a low-level representation of a specific place in the code,
// helping in debugging by pinpointing the exact line of execution that caused
// an error or event.
//
// Example usage:
//
//	var f Frame = Frame(0x1234567890)
//	fmt.Println(f) // Prints the value of the program counter + 1
type Frame uintptr

// StackTrace represents a stack of Frames, which are ordered from the
// innermost (newest) function call to the outermost (oldest) function call
// in the current stack. It provides a high-level representation of the
// sequence of function calls (the call stack) that led to the current execution
// point, typically used for debugging or error reporting. The `StackTrace`
// contains a slice of `Frame` values, which can be interpreted as program
// counters in the call stack.
//
// The `StackTrace` can be used to generate detailed stack traces when an
// error occurs, helping developers track down the sequence of function
// calls that resulted in the error. For example, it may be used in conjunction
// with the `underlying` and `underlyingStack` types to record where an error
// occurred in the code (using `Callers()` to populate the stack) and provide
// information on the call path leading to the error.
//
// Example usage:
//
//	var trace StackTrace = StackTrace{Frame(0x1234567890), Frame(0x0987654321)}
//	fmt.Println(trace) // Prints the stack trace with the frames
type StackTrace []Frame

// pagination represents pagination details for paginated API responses.
type pagination struct {
	page       int  // Current page number.
	perPage    int  // Number of items per page.
	totalPages int  // Total number of pages.
	totalItems int  // Total number of items available.
	isLast     bool // Indicates whether this is the last page.
}

// meta represents metadata information about an API response.
type meta struct {
	apiVersion    string                 // API version used for the request.
	requestID     string                 // Unique identifier for the request, useful for debugging.
	locale        string                 // Locale used for the request, e.g., "en-US".
	requestedTime time.Time              // Timestamp when the request was made.
	customFields  map[string]interface{} // Additional custom metadata fields.
}

// header represents a structured header for API responses.
type header struct {
	code        int    // Application-specific status code.
	text        string // Human-readable status text.
	typez       string // Type or category of the status, e.g., "info", "error".
	description string // Detailed description of the status.
}

// wrapper is the main structure for wrapping API responses, including metadata, data, and debugging information.
type wrapper struct {
	statusCode int            // HTTP status code for the response.
	total      int            // Total number of items (used in non-paginated responses).
	message    string         // A message providing additional context about the response.
	data       any            // The primary data payload of the response.
	path       string         // Request path for which the response is generated.
	header     *header        // Structured header details for the response.
	meta       *meta          // Metadata about the API response.
	pagination *pagination    // Pagination details, if applicable.
	debug      map[string]any // Debugging information (useful for development).
	errors     error          // Internal errors (not exposed in JSON responses).
	cachedWrap map[string]any // Cached response data for performance optimization.
	cacheHash  string         // Hash of the cached response, used for cache validation.
	cacheMutex sync.RWMutex   // Mutex for synchronizing access to the cached response data.
}

// stack represents a stack of program counters. It is a slice of `uintptr`
// values that store the addresses (or locations) of program counters
// from a function call stack. This is useful for tracking where errors
// or events originated in the code, and it is used by the other error
// types to capture a stack trace at the time of error creation.
//
// The stack is typically populated using a function like `Callers()`,
// which retrieves the current call stack.
type stack []uintptr

// underlying represents an error with an associated message and a stack trace.
// It is used to encapsulate the message (i.e., what went wrong) and the
// program's call stack (i.e., where the error occurred) at the point when
// the error is created. This type is used internally by functions like
// `WithError` and `WithErrorf`.
//
// Fields:
//   - msg: The error message describing what went wrong. This is a string
//     representation of the error context.
//   - stack: A pointer to a stack of program counters that represent the
//     stack trace where the error was created.
type underlying struct {
	msg    string // Message that describes the error
	*stack        // Pointer to the stack trace where the error occurred
}

// underlyingStack is an error type that contains an existing error and its
// associated stack trace. It is used when we need to annotate an existing
// error with a stack trace, preserving the original error while adding
// additional information such as the place where the annotation occurred.
//
// This type is typically used in functions like `Wrap` and `Wrapf`.
//
// Fields:
//   - error: The original error that is being wrapped or annotated.
//   - stack: A pointer to the stack trace captured at the point this
//     annotation was added.
type underlyingStack struct {
	error  // The original error being wrapped or annotated
	*stack // Pointer to the stack trace where the annotation occurred
}

// underlyingMessage represents an error with a cause (another error)
// and a message. This type is used to create errors that have both
// an associated message and a "cause," which can be another error
// that led to the current one. It is commonly used to propagate errors
// and add context to the original error message.
//
// Fields:
//   - cause: The underlying error that caused this error. This could be
//     another error returned from a function, allowing us to trace back
//     the error chain.
//   - msg: A string message describing the context of the error, which
//     can be appended to the original cause message to provide more clarity.
type underlyingMessage struct {
	cause error  // The original error being wrapped or annotated
	msg   string // The message describing the additional context for the error
}
