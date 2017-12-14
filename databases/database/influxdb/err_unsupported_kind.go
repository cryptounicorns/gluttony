package influxdb

import (
	"fmt"
	"reflect"
)

type ErrUnsupportedKind struct {
	Kind reflect.Kind
}

func (e *ErrUnsupportedKind) Error() string {
	return fmt.Sprintf(
		"Unsupported kind '%s'",
		e.Kind,
	)
}

func NewErrUnsupportedKind(k reflect.Kind) *ErrUnsupportedKind {
	return &ErrUnsupportedKind{
		Kind: k,
	}
}
