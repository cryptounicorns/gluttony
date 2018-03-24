package clickhouse

import (
	"log"
	"net/url"
	"strconv"

	"github.com/corpix/loggers"
	"github.com/kshvakov/clickhouse"
)

func Connect(c Config, l loggers.Logger) (clickhouse.Clickhouse, error) {
	l.Debug("Connecting...")

	url, _ := url.Parse("")
	url.Scheme = "tcp"
	url.Host = c.Host + ":" + c.Port
	values := url.Query()
	{
		// values.Set("database", c.Database)
		values.Set("username", c.User)
		values.Set("password", c.Pass)

		values.Set("compress", strconv.FormatBool(c.Compress))
		values.Set("debug", strconv.FormatBool(c.Debug))

		values.Set("read_timeout", strconv.FormatInt(c.ReadTimeout, 10))
		values.Set("write_timeout", strconv.FormatInt(c.WriteTimeout, 10))
	}
	url.RawQuery = values.Encode()

	connect, err := clickhouse.OpenDirect(url.String())
	if err != nil {
		log.Fatal(err)
	}
	return connect, nil
}
