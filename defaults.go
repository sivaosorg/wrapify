package wrapify

import (
	"time"

	"github.com/sivaosorg/unify4g"
)

// defaultMetaValues returns the default values for the `meta` struct.
//
// This function creates a new `meta` instance with the default values,
// including the locale, API version, requested time, and request ID.
//
// Returns:
//   - A pointer to a newly created `meta` instance with the default values.
func defaultMetaValues() *meta {
	return Meta().
		WithLocale("en_US"). // vi_VN, en_US
		WithApiVersion("v0.0.1").
		WithRequestedTime(time.Now()).
		WithRequestID(unify4g.GenerateCryptoID())
}
