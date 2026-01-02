package hash

import (
	"fmt"

	"github.com/sivaosorg/wrapify/pkg/strutil"
)

// validate checks if options are valid.
//
// Parameters:
//   - o: A pointer to the `Options` struct to validate.
//
// Returns:
//   - An error if the options are invalid, otherwise nil.
func (o *Options) validate() error {
	if o.Hasher == nil {
		return fmt.Errorf("pkg.hash.options: hasher cannot be nil")
	}
	if strutil.IsEmpty(o.TagName) {
		return fmt.Errorf("pkg.hash.options: tag name cannot be empty")
	}
	return nil
}
