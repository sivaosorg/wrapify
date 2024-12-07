package wrapify

import "time"

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
	Type        string // Type or category of the status, e.g., "info", "error".
	description string // Detailed description of the status.
}

// wrapper is the main structure for wrapping API responses, including metadata, data, and debugging information.
type wrapper struct {
	statusCode int                    // HTTP status code for the response.
	total      int                    // Total number of items (used in non-paginated responses).
	message    string                 // A message providing additional context about the response.
	data       interface{}            // The primary data payload of the response.
	path       string                 // Request path for which the response is generated.
	header     *header                // Structured header details for the response.
	meta       *meta                  // Metadata about the API response.
	pagination *pagination            // Pagination details, if applicable.
	debug      map[string]interface{} // Debugging information (useful for development).
	errors     error                  // Internal errors (not exposed in JSON responses).
}
