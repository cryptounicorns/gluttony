package cli

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/urfave/cli"
)

// Ticker ok this is embarrassing will move it ASAP
type Ticker struct {
	High         float64 `json:"high"`
	Low          float64 `json:"low"`
	Vol          float64 `json:"vol"`
	Last         float64 `json:"last"`
	Buy          float64 `json:"buy"`
	Sell         float64 `json:"sell"`
	Timestamp    uint64  `json:"timestamp"`
	CurrencyPair string  `json:"currencyPair"`
	Market       string  `json:"market"`
}

var (
	// RootCommands is a list of subcommands for the application.
	RootCommands = []cli.Command{}

	// RootFlags is a list of flags for the application.
	RootFlags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "application configuration file",
			EnvVar: "CONFIG",
			Value:  "config.json",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "add this flag to enable debug mode",
		},
	}
)

// FromNanotime this is just a convenience method
func FromNanotime(ts uint64) time.Time {
	sec := ts / 1000000000
	nsec := ts - sec*1000000000
	return time.Unix(int64(sec), int64(nsec))
}

// RootAction is executing when program called without any subcommand.
func RootAction(c *cli.Context) error {
	file, err := os.Open(Config.JSONInput.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	bp, err := client.NewBatchPoints(Config.InfluxDB.Batch)
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		record := Ticker{}
		err = json.Unmarshal(scanner.Bytes(), &record)
		if err != nil {
			log.Fatal(err)
		}
		tags := map[string]string{
			"currency_pair": record.CurrencyPair,
			"market":        record.Market,
		}
		fields := map[string]interface{}{
			"high": record.High,
			"low":  record.Low,
			"vol":  record.Vol,
			"last": record.Last,
			"buy":  record.Buy,
			"sell": record.Sell,
		}
		pt, err := client.NewPoint("ticker", tags, fields, FromNanotime(record.Timestamp))
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	influxClient, err := client.NewHTTPClient(Config.InfluxDB.Client)
	if err != nil {
		log.Fatal(err)
	}

	// Write the batch
	if err := influxClient.Write(bp); err != nil {
		log.Fatal(err)
	}

	return nil
}
