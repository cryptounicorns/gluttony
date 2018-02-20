package influxdb

import (
	"github.com/corpix/time"
	client "github.com/influxdata/influxdb/client/v2"
)

type PointConfig struct {
	Name               string   `validate:"required"`
	Fields             []string `validate:"required"`
	Tags               []string
	Timestamp          string `validate:"required"`
	TimestampPrecision string `validate:"required,eq=nanosecond|eq=microsecond|eq=millisecond|eq=second"`
}

type BatchConfig struct {
	Points        client.BatchPointsConfig `validate:"required"`
	FlushInterval time.Duration            `validate:"required"`
	Size          uint                     `validate:"required"`
}

type WriterConfig struct {
	Batch BatchConfig `validate:"required,dive"`
	Point PointConfig `validate:"required,dive"`
}

type Config struct {
	Client client.HTTPConfig `validate:"required"`
	Writer WriterConfig      `validate:"required,dive"`
}
