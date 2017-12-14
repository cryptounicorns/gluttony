package errors

import (
	"fmt"
)

type ErrKeyNotFound struct {
	Key string
}

func (e *ErrKeyNotFound) Error() string {
	return fmt.Sprintf(
		"Key '%s' was not found in the store",
		e.Key,
	)
}

func NewErrKeyNotFound(key string) *ErrKeyNotFound {
	return &ErrKeyNotFound{
		Key: key,
	}
}
