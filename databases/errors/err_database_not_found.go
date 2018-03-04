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
