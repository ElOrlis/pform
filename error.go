package pform

import "fmt"

func newParseError(f, t, v string, e error) error {
	return ParseError{f, t, v, e}
}

type ParseError struct {
	Field string
	Type  string
	Value string
	Err   error
}

func (e ParseError) Error() string {
	return fmt.Sprintf(
		"failed to parse %q as %s for field %q: %v",
		e.Value,
		e.Type,
		e.Field,
		e.Err,
	)
}

func (e ParseError) Unwrap() error {
	return e.Err
}

func newRequiredFieldError(f string) error {
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
