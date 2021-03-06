package databases

import (
	"fmt"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	client "github.com/influxdata/influxdb/client/v2"

	"github.com/cryptounicorns/gluttony/databases/database/influxdb"
	"github.com/cryptounicorns/gluttony/databases/errors"
	"github.com/cryptounicorns/gluttony/databases/record"
)

type Database interface {
	Create(record.Record) error
	Close() error
}

func FromConfig(c Config, conn Connection, l loggers.Logger) (Database, error) {
	var (
		t   = strings.ToLower(c.Type)
		log = prefixwrapper.New(
			fmt.Sprintf("Database %s: ", t),
			l,
		)
	)

	switch t {
	case influxdb.Name:
		return influxdb.FromConfig(
			*c.Influxdb,
			conn.(client.Client),
			log,
		)
	default:
		return nil, errors.NewErrUnknownDatabaseType(c.Type)
	}
}
