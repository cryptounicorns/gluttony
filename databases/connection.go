package databases

import (
	"fmt"
	"io"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"

	"github.com/cryptounicorns/gluttony/databases/database/influxdb"
	"github.com/cryptounicorns/gluttony/databases/errors"
)

type Connection = io.Closer

func Connect(c Config, l loggers.Logger) (Connection, error) {
	var (
		t   = strings.ToLower(c.Type)
		log = prefixwrapper.New(
			fmt.Sprintf("DatabaseConnection %s: ", t),
			l,
		)
	)

	switch t {
	case influxdb.Name:
		return influxdb.Connect(
			*c.Influxdb,
			log,
		)
	default:
		return nil, errors.NewErrUnknownDatabaseType(c.Type)
	}
}
