package memoryttl

import (
	"github.com/corpix/time"
)

type Config struct {
	TTL        time.Duration
	Resolution time.Duration
}
