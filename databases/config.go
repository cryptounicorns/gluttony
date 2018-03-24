package databases

import (
	"github.com/cryptounicorns/gluttony/databases/database/clickhouse"
	"github.com/cryptounicorns/gluttony/databases/database/influxdb"
)

type Config struct {
	Type       string `validate:"required"`
	Influxdb   *influxdb.Config
	Clickhouse *clickhouse.Config
}
