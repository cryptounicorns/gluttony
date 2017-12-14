package databases

import (
	"github.com/cryptounicorns/gluttony/databases/database/influxdb"
)

type Config struct {
	Type     string          `validator:"required"`
	Influxdb influxdb.Config `validator:"required,dive"`
}
