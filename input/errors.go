package input

import (
	"fmt"
)

type ErrUnknownInputType struct {
	t string
}

func (e *ErrUnknownInputType) Error() string {
	return fmt.Sprintf(
		"Unknown input type '%s'",
		e.t,
	)
}
func NewErrUnknownInputType(t string) error {
	return &ErrUnknownInputType{t}
}
