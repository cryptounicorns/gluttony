package errors

import (
	"fmt"
)

// ErrEndpoint is an error indicating that request to endpoint
// resulted in error.
type ErrEndpoint struct {
	url   string
	error string
	code  int
	want  int
}

func (e *ErrEndpoint) Error() string {
	return fmt.Sprintf(
		"Endpoint request to '%s' finished with error '%s' and code '%d' instead of '%d'",
		e.url,
		e.error,
		e.code,
		e.want,
	)
}

// NewErrEndpoint creates new ErrEndpoint.
func NewErrEndpoint(url string, error string, code int, want int) error {
	return &ErrEndpoint{
		url,
		error,
		code,
		want,
	}
}
