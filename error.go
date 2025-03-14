package pf

import "fmt"

func NewRequiredFieldError(f string) error {
	return RequiredFieldError{f}
}

type RequiredFieldError struct {
	field string
}

func (r RequiredFieldError) Error() string {
	return fmt.Sprintf("missing %s field", r.field)
}

func (r RequiredFieldError) Unwrap() error {
	return fmt.Errorf("missing %s field", r.field)
}
