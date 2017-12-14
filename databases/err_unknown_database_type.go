package databases

import (
	"fmt"
)

type ErrUnknownDatabaseType struct {
	t string
}

func (e *ErrUnknownDatabaseType) Error() string {
	return fmt.Sprintf(
		"Unknown database type '%s'",
		e.t,
	)
}
func NewErrUnknownDatabaseType(t string) error {
	return &ErrUnknownDatabaseType{t}
}
