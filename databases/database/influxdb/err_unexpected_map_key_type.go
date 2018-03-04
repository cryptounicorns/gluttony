package influxdb

import (
	"fmt"
)

type ErrUnexpectedMapKeyType struct {
	Want string
	Got  interface{}
}

func (e *ErrUnexpectedMapKeyType) Error() string {
	return fmt.Sprintf(
		"Unexpected map key type, want '%s', got '%#v'",
		e.Want,
		e.Got,
	)
}

func NewErrUnexpectedMapKeyType(want string, got interface{}) *ErrUnexpectedMapKeyType {
	return &ErrUnexpectedMapKeyType{
		Want: want,
		Got:  got,
	}
}
