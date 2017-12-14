package stores

import (
	"fmt"
)

type ErrUnknownStoreType struct {
	t string
}

func (e *ErrUnknownStoreType) Error() string {
	return fmt.Sprintf(
		"Unknown store type '%s'",
		e.t,
	)
}
func NewErrUnknownStoreType(t string) error {
	return &ErrUnknownStoreType{t}
}
