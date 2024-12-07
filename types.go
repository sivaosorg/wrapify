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
	StatusCode int                     `json:"status_code" binding:"required"` // HTTP status code for the response.
	Total      int                     `json:"total"`                          // Total number of items (used in non-paginated responses).
	Message    string                  `json:"message"`                        // A message providing additional context about the response.
	Data       interface{}             `json:"data,omitempty"`                 // The primary data payload of the response.
	Path       string                  `json:"path,omitempty"`                 // Request path for which the response is generated.
	Header     *header                 `json:"headers,omitempty"`              // Structured header details for the response.
	Meta       *meta                   `json:"meta,omitempty"`                 // Metadata about the API response.
	Pagination *pagination             `json:"pagination,omitempty"`           // Pagination details, if applicable.
	Debug      *map[string]interface{} `json:"debug,omitempty"`                // Debugging information (useful for development).
	Errors     error                   `json:"-"`                              // Internal errors (not exposed in JSON responses).
}
