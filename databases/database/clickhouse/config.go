package clickhouse

import "time"

type Config struct {
	Host string
	Port string

	User string
	Pass string

	Compress bool
	Debug    bool

	ReadTimeout  int64
	WriteTimeout int64

	FlushInterval time.Duration
}
