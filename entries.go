package wrapify

import (
	"time"

	"github.com/sivaosorg/unify4g"
)

// Parse takes a JSON string as input and maps it into a 'wrapper' struct, populating its fields
// with data extracted from the JSON. The function handles different parts of the response such as
// status code, message, pagination, headers, metadata, and debugging information. It uses
// `unify4g.UnmarshalFromStringN` for unmarshaling the JSON string into a map, then assigns the
// respective values from the map to the fields of the `wrapper` struct.
//
// The function extracts information from the provided JSON, handles nested
// objects, and populates the `wrapper` struct with relevant values for:
//   - statusCode: The HTTP status code of the response
//   - total: The total number of items in the response
//   - message: A message associated with the response
//   - data: The primary data in the response (can be any type)
//   - path: The request path of the API that generated the response
//   - debug: Debugging information, if available
//   - header: An optional header with various status-related fields
//   - meta: Metadata about the API response such as version, locale, etc.
//   - pagination: Pagination details, such as page number and total items
//
// Returns:
//   - A pointer to a `wrapper` struct containing the parsed data
//   - An error if the JSON string cannot be parsed or is invalid
//
// This function is responsible for taking raw JSON data and converting it into a
// structured format, making it easier to work with in Go. It handles various
// nested objects (like `header`, `meta`, and `pagination`), checking for the
// presence of keys and converting them into their appropriate types.
func Parse(json string) (w *wrapper, err error) {
	var data map[string]any
	err = unify4g.UnmarshalFromStringN(json, &data)
	if err != nil {
		return nil, WithErrStack(err)
	}
	if len(data) == 0 {
		return nil, WithErrorf("the wrapper response is empty with JSON string: %v", json)
	}
	w = &wrapper{}
	if value, exists := data["status_code"].(float64); exists {
		w.statusCode = int(value)
	}
	if value, exists := data["total"].(float64); exists {
		w.total = int(value)
	}
	if value, exists := data["message"].(string); exists {
		w.message = value
	}
	if value, exists := data["data"]; exists {
		w.data = value
	}
	if value, exists := data["path"].(string); exists {
		w.path = value
	}
	if value, exists := data["debug"].(map[string]any); exists {
		w.debug = value
	}
	if values, exists := data["header"].(map[string]any); exists {
		header := &header{}
		if value, exists := values["code"].(float64); exists {
			header.code = int(value)
		}
		if value, exists := values["text"].(string); exists {
			header.text = value
		}
		if value, exists := values["type"].(string); exists {
			header.typez = value
		}
		if value, exists := values["description"].(string); exists {
			header.description = value
		}
		w.header = header
	}
	if values, exists := data["meta"].(map[string]any); exists {
		meta := &meta{}
		if value, exists := values["api_version"].(string); exists {
			meta.apiVersion = value
		}
		if value, exists := values["locale"].(string); exists {
			meta.locale = value
		}
		if value, exists := values["request_id"].(string); exists {
			meta.requestID = value
		}
		if customFields, exists := values["custom_fields"].(map[string]any); exists {
			meta.customFields = customFields
		}
		if value, exists := values["requested_time"].(string); exists {
			if t, err := time.Parse(time.RFC3339, value); err == nil {
				meta.requestedTime = t
			}
		}
		w.meta = meta
	}
	if values, exists := data["pagination"].(map[string]any); exists {
		pagination := &pagination{}
		if value, exists := values["page"].(float64); exists {
			pagination.page = int(value)
		}
		if value, exists := values["per_page"].(float64); exists {
			pagination.perPage = int(value)
		}
		if value, exists := values["total_pages"].(float64); exists {
			pagination.totalPages = int(value)
		}
		if value, exists := values["total_items"].(float64); exists {
			pagination.totalItems = int(value)
		}
		if value, exists := values["is_last"].(bool); exists {
			pagination.isLast = value
		}
		w.pagination = pagination
	}
	return w, nil
}

// Deserialize converts a map containing API response data into a `wrapper` struct
// by serializing the map into JSON format and then parsing it.
//
// The function is a helper that bridges between raw map data (e.g., deserialized JSON
// or other dynamic input) and the strongly-typed `wrapper` struct used in the codebase.
// It first converts the input map into a JSON string using `unify4g.JsonN`, then calls
// the `Parse` function to handle the deserialization and field mapping to the `wrapper`.
//
// Parameters:
//   - data: A map[string]interface{} containing the API response data.
//     The map should include keys like "status_code", "message", "meta", etc.,
//     that conform to the expected structure of a `wrapper`.
//
// Returns:
//   - A pointer to a `wrapper` struct populated with data from the map.
//   - An error if the map is empty or if the JSON serialization/parsing fails.
//
// Error Handling:
//   - If the input map is empty or nil, the function returns an error
//     indicating that the data is invalid.
//   - If serialization or parsing fails, the error from `Parse` or `unify4g.JsonN`
//     is propagated, providing context about the failure.
//
// Usage:
// This function is particularly useful when working with raw data maps (e.g., from
// dynamic inputs or unmarshaled data) that need to be converted into the `wrapper`
// struct for further processing.
//
// Example:
//
//	rawData := map[string]interface{}{
//	    "status_code": 200,
//	    "message": "Success",
//	    "data": "response body",
//	}
//	wrapper, err := wrapify.Deserialize(rawData)
//	if err != nil {
//	    log.Println("Error extracting wrapper:", err)
//	} else {
//	    log.Println("Wrapper:", wrapper)
//	}
func Deserialize(data map[string]any) (w *wrapper, err error) {
	if len(data) == 0 {
		return nil, WithErrorf("data is nil/null")
	}
	json := unify4g.JsonN(data)
	return Parse(json)
}
