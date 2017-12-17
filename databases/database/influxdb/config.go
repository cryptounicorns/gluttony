package influxdb

import (
	"github.com/corpix/time"
	client "github.com/influxdata/influxdb/client/v2"
)

type PointConfig struct {
	Name               string   `validator:"required"`
	Fields             []string `validator:"required"`
	Tags               []string
	Timestamp          string `validator:"required"`
	TimestampPrecision string `validator:"required,eq=nanosecond|eq=microsecond|eq=millisecond|eq=second"`
}

type BatchConfig struct {
	Points        client.BatchPointsConfig `validator:"required"`
	FlushInterval time.Duration            `validator:"required"`
	Size          uint                     `validator:"required"`
}

type WriterConfig struct {
	Batch BatchConfig `validator:"required,dive"`
	Point PointConfig `validator:"required,dive"`
}

type Config struct {
	Client client.HTTPConfig `validator:"required"`
	Writer WriterConfig      `validator:"required,dive"`
}
