package dto

import "fmt"

func errRequired(field string) error {
	return fmt.Errorf("%s is required", field)
}

func errInvalid(field, reason string) error {
	return fmt.Errorf("%s is invalid: %s", field, reason)
}
