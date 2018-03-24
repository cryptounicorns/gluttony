package influxdb

import (
	"fmt"
	"reflect"
)

type ErrUnsupportedTimestampPrecision struct {
	TimestampPrecision string
}

func (e *ErrUnsupportedTimestampPrecision) Error() string {
	return fmt.Sprintf(
		"Unsupported timestamp precision '%s'",
		e.TimestampPrecision,
	)
}

func NewErrUnsupportedTimestampPrecision(s string) *ErrUnsupportedTimestampPrecision {
	return &ErrUnsupportedTimestampPrecision{
		TimestampPrecision: s,
	}
}

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
