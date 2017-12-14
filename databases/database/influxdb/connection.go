package influxdb

import (
	"github.com/corpix/loggers"
	client "github.com/influxdata/influxdb/client/v2"
)

func Connect(c Config, l loggers.Logger) (client.Client, error) {
	l.Debug("Connecting...")
	return client.NewHTTPClient(c.Client)
}
