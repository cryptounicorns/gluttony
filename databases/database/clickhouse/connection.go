package clickhouse

import (
	"log"

	"github.com/corpix/loggers"
	"github.com/kshvakov/clickhouse"
)

func Connect(c Config, l loggers.Logger) (clickhouse.Clickhouse, error) {
	l.Debug("Connecting...")

	connect, err := clickhouse.OpenDirect("tcp://127.0.0.1:9000?username=&debug=true")
	if err != nil {
		log.Fatal(err)
	}
	return connect, nil
}
