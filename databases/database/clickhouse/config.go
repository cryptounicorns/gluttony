package clickhouse

import "time"

type Config struct {
	Host string
	Port int

	User string
	Pass string

	Compress bool
	Debug    bool

	FlushInterval time.Duration
}
