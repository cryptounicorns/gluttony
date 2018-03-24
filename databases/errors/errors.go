package errors

import (
	"fmt"
)

type ErrDatabaseNotFound struct {
	Name string
}

func (e ErrDatabaseNotFound) Error() string {
	return fmt.Sprintf(
		"Database '%s' was not found",
		e.Name,
	)
}
func NewErrDatabaseNotFound(name string) error {
	return ErrDatabaseNotFound{name}
}

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
