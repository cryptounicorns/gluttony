package influxdb

import (
	"reflect"
	"strings"
	"time"

	"github.com/corpix/loggers"
	client "github.com/influxdata/influxdb/client/v2"

	"github.com/cryptounicorns/gluttony/databases/errors"
	"github.com/cryptounicorns/gluttony/databases/record"
)

const (
	Name = "influxdb"
)

type InfluxDB struct {
	config         Config
	timestampField reflect.Value
	tagNames       map[string]bool
	fieldNames     map[string]bool
	client         client.Client
	log            loggers.Logger
	points         chan *client.Point
	done           chan struct{}
}

func (d *InfluxDB) timeoutBatchWriter() {
	var (
		interval = d.config.Writer.Batch.FlushInterval.Duration()
	)

	if interval == 0 {
		interval = 5 * time.Second
	}

	for {
		select {
		case <-d.done:
			return
		case <-time.After(interval):
			d.writeBatch()
		}
	}
}

func (d *InfluxDB) writeBatch() error {
	d.log.Debug("Preparing batch...")
	var (
		bp  client.BatchPoints
		ps  []*client.Point
		err error
	)

loop:
	for {
		select {
		case v := <-d.points:
			ps = append(
				ps,
				v,
			)
		default:
			break loop
		}

		if uint(len(ps)) >= d.config.Writer.Batch.Size {
			break loop
		}
	}

	if len(ps) == 0 {
		d.log.Debug("Nothing to write, waiting for new data")
		return nil
	}

	bp, err = client.NewBatchPoints(d.config.Writer.Batch.Points)
	if err != nil {
		return err
	}

	bp.AddPoints(ps)

	d.log.Debug("Writing prepared batch of length ", len(ps))
	err = d.client.Write(bp)
	if err != nil {
		d.log.Error("Failed to write batch with error ", err.Error())

		// FIXME: Influxdb module failed to provide types for errors
		// https://github.com/influxdata/influxdb/issues/9507
		// Matching error messages here, be careful with influxdb client updates
		// and blame Go for inventing fmt.Errorf and errors.New crap ¯\_(ツ)_/¯

		databaseNotFoundErrorPrefix := `{"error":"database not found: `
		if errStr := err.Error(); strings.HasPrefix(errStr, databaseNotFoundErrorPrefix) {
			databaseName := errStr[len(databaseNotFoundErrorPrefix)+2 : len(errStr)-5]
			return errors.NewErrDatabaseNotFound(databaseName)
		}

		return err
	}

	return nil
}

func (d *InfluxDB) tags(rk reflect.Kind, rv reflect.Value) (map[string]string, error) {
	var (
		res = map[string]string{}
		raw interface{}
		k   string
		v   string
		ok  bool
	)

	switch rk {
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			k, ok = key.Interface().(string)
			if !ok {
				return nil, NewErrUnexpectedMapKeyType(
					"string",
					key.Interface(),
				)
			}

			if !d.tagNames[k] {
				continue
			}

			raw = rv.MapIndex(key).Interface()
			if raw == nil {
				continue
			}

			v, ok = raw.(string)
			if !ok {
				return nil, NewErrUnexpectedMapValueType(
					"string",
					rv.MapIndex(key).Interface(),
				)
			}

			res[k] = v
		}
	default:
		return nil, NewErrUnsupportedKind(rk)
	}

	return res, nil
}

func (d *InfluxDB) fields(rk reflect.Kind, rv reflect.Value) (map[string]interface{}, error) {
	var (
		res = map[string]interface{}{}
		k   string
		ok  bool
	)

	switch rk {
	case reflect.Map:
		for _, key := range rv.MapKeys() {
			k, ok = key.Interface().(string)
			if !ok {
				return nil, NewErrUnexpectedMapKeyType(
					"string",
					key.Interface(),
				)
			}

			if !d.fieldNames[k] {
				continue
			}

			res[k] = rv.MapIndex(key).Interface()
		}
	default:
		return nil, NewErrUnsupportedKind(rk)
	}

	return res, nil
}

func (d *InfluxDB) timestamp(rv reflect.Value) (time.Time, error) {
	if d.config.Writer.Point.Timestamp == "" {
		return time.Now(), nil
	}

	var (
		v = rv.
			MapIndex(d.timestampField).
			Interface()
		ts      float64
		sec     int64
		nanoSec int64
		t       time.Time
		ok      bool
	)

	// XXX: This is because of Lua :'(
	// Otherwise we will be dealing with uint64
	ts, ok = v.(float64)
	if !ok {
		return t, NewErrUnexpectedMapKeyType(
			"float64",
			v,
		)
	}

	switch d.config.Writer.Point.TimestampPrecision {
	case "nanosecond":
		sec = int64(ts) / 1000000000
		nanoSec = int64(ts) - sec*1000000000
	case "microsecond":
		sec = int64(ts) / 1000000
		nanoSec = int64(ts) - sec*1000000
	case "millisecond":
		sec = int64(ts) / 1000
		nanoSec = int64(ts) - sec*1000
	case "second":
		sec = int64(ts)
		nanoSec = 0
	default:
		return t, NewErrUnsupportedTimestampPrecision(
			d.config.Writer.Point.TimestampPrecision,
		)
	}

	return time.Unix(sec, nanoSec), nil
}

func (d *InfluxDB) Create(r record.Record) error {
	var (
		rk        = reflect.TypeOf(r).Kind()
		rv        = reflect.ValueOf(r)
		p         *client.Point
		tags      map[string]string
		fields    map[string]interface{}
		timestamp time.Time
		err       error
	)

	tags, err = d.tags(rk, rv)
	if err != nil {
		return err
	}

	fields, err = d.fields(rk, rv)
	if err != nil {
		return err
	}

	timestamp, err = d.timestamp(rv)
	if err != nil {
		return err
	}

	p, err = client.NewPoint(
		d.config.Writer.Point.Name,
		tags,
		fields,
		timestamp,
	)
	if err != nil {
		return err
	}

retry:
	select {
	case d.points <- p:
	default:
		err = d.writeBatch()
		if err != nil {
			return err
		}
		goto retry
	}

	return nil
}

func (d *InfluxDB) Close() error {
	close(d.done)

	return nil
}

func FromConfig(c Config, cl client.Client, l loggers.Logger) (*InfluxDB, error) {
	var (
		tagNames   = map[string]bool{}
		fieldNames = map[string]bool{}
	)

	for _, tag := range c.Writer.Point.Tags {
		tagNames[tag] = true
	}

	for _, field := range c.Writer.Point.Fields {
		fieldNames[field] = true
	}

	var (
		db = &InfluxDB{
			config:         c,
			timestampField: reflect.ValueOf(c.Writer.Point.Timestamp),
			tagNames:       tagNames,
			fieldNames:     fieldNames,
			client:         cl,
			log:            l,
			points: make(
				chan *client.Point,
				c.Writer.Batch.Size,
			),
			done: make(chan struct{}),
		}
	)

	go db.timeoutBatchWriter()

	return db, nil
}
