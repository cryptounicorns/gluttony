package influxdb

import (
	"fmt"
)

type ErrUnexpectedMapValueType struct {
	Want string
	Got  interface{}
}

func (e *ErrUnexpectedMapValueType) Error() string {
	return fmt.Sprintf(
		"Unexpected map value type, want '%s', got '%#v'",
		e.Want,
		e.Got,
	)
}

func NewErrUnexpectedMapValueType(want string, got interface{}) *ErrUnexpectedMapValueType {
	return &ErrUnexpectedMapValueType{
		Want: want,
		Got:  got,
	}
}
