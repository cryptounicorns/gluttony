package influxdb

import (
	"fmt"
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
