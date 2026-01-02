package hash

import "fmt"

// validate checks if options are valid.
//
// Parameters:
//   - o: A pointer to the `Options` struct to validate.
//
// Returns:
//   - An error if the options are invalid, otherwise nil.
func (o *Options) validate() error {
	if o.Hasher == nil {
		return fmt.Errorf("options: hasher cannot be nil")
	}
	if o.TagName == "" {
		return fmt.Errorf("options: tag name cannot be empty")
	}
	return nil
}
