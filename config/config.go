package config

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/corpix/formats"
	"github.com/cryptounicorns/gluttony/logger"
	"github.com/imdario/mergo"
	"github.com/influxdata/influxdb/client/v2"
)

var (
	// LoggerConfig represents default logger config.
	LoggerConfig = logger.Config{
		Level: "info",
	}

	// Default represents default application config.
	Default = Config{
		Logger: LoggerConfig,
		JSONInput: JSONInputConfig{
			Path: "/var/gluttony/tickers.json",
		},
		InfluxDB: InfluxDBConfig{
			Client: client.HTTPConfig{
				Addr: "localhost",
			},
		},
	}
)

// JSONInputConfig input file configuration (sort of)
// @TODO Extend with fsnotify and/or replace with NSQ consumer
type JSONInputConfig struct {
	Path string
}

// InfluxDBConfig InfluxDB-related configuration
type InfluxDBConfig struct {
	Client client.HTTPConfig
	Batch  client.BatchPointsConfig
}

// Config represents application configuration structure.
type Config struct {
	Logger    logger.Config
	JSONInput JSONInputConfig
	InfluxDB  InfluxDBConfig
}

// FromReader fills Config structure `c` passed by reference with
// parsed config data in some `f` from reader `r`.
// It copies `Default` into the target structure before unmarshaling
// config, so it will have default values.
func FromReader(f formats.Format, r io.Reader, c *Config) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = mergo.Merge(c, Default)
	if err != nil {
		return err
	}

	return f.Unmarshal(data, c)
}

// FromFile fills Config structure `c` passed by reference with
// parsed config data from file in `path`.
func FromFile(path string, c *Config) error {
	f, err := formats.NewFromPath(path)
	if err != nil {
		return err
	}

	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	return FromReader(f, r, c)
}
