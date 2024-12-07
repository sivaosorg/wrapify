package wrapify

import "time"

// pagination represents pagination details for paginated API responses.
type pagination struct {
	Page       int  `json:"page,omitempty"`        // Current page number.
	PerPage    int  `json:"per_page,omitempty"`    // Number of items per page.
	TotalPages int  `json:"total_pages,omitempty"` // Total number of pages.
	TotalItems int  `json:"total_items,omitempty"` // Total number of items available.
	IsLast     bool `json:"is_last,omitempty"`     // Indicates whether this is the last page.
}

// meta represents metadata information about an API response.
type meta struct {
	ApiVersion    string                 `json:"api_version,omitempty"`    // API version used for the request.
	RequestId     string                 `json:"request_id,omitempty"`     // Unique identifier for the request, useful for debugging.
	Locale        string                 `json:"locale,omitempty"`         // Locale used for the request, e.g., "en-US".
	RequestedTime time.Time              `json:"requested_time,omitempty"` // Timestamp when the request was made.
	CustomFields  map[string]interface{} `json:"custom_fields,omitempty"`  // Additional custom metadata fields.
}

// header represents a structured header for API responses.
type header struct {
	Code        int    `json:"code,omitempty"`        // Application-specific status code.
	Text        string `json:"text,omitempty"`        // Human-readable status text.
	Type        string `json:"type,omitempty"`        // Type or category of the status, e.g., "info", "error".
	Description string `json:"description,omitempty"` // Detailed description of the status.
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
