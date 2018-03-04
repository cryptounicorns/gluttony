package errors

import (
	"fmt"
)

type ErrUnknownDatabaseType struct {
	Type string
}

func (e ErrUnknownDatabaseType) Error() string {
	return fmt.Sprintf(
		"Unknown database type '%s'",
		e.Type,
	)
}
func NewErrUnknownDatabaseType(t string) error {
	return ErrUnknownDatabaseType{t}
}
