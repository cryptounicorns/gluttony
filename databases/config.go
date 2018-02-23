package databases

import (
	"github.com/cryptounicorns/gluttony/databases/database/influxdb"
)

type Config struct {
	Type     string `validate:"required"`
	Influxdb *influxdb.Config
}
