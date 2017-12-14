package writer

import (
	"fmt"
	"time"
)

type ErrScheduleTimeout time.Duration

func (e ErrScheduleTimeout) Error() string {
	return fmt.Sprintf(
		"Schedule timeouted after '%s'",
		time.Duration(e),
	)
}

func NewErrScheduleTimeout(d time.Duration) ErrScheduleTimeout {
	return ErrScheduleTimeout(d)
}
