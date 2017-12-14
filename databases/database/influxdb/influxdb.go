package influxdb

import (
	"reflect"
	"time"

	"github.com/corpix/loggers"
	client "github.com/influxdata/influxdb/client/v2"

	"github.com/cryptounicorns/gluttony/databases/record"
)

const (
	Name = "influxdb"
)

type InfluxDB struct {
	config     Config
	tagNames   map[string]bool
	fieldNames map[string]bool
	client     client.Client
	log        loggers.Logger
	points     chan *client.Point
	done       chan struct{}
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
		return nil
	}

	bp, err = client.NewBatchPoints(d.config.Writer.Batch.Points)
	if err != nil {
		return err
	}

	bp.AddPoints(ps)

	d.log.Debug("Writing prepared batch of length ", len(ps))
	return d.client.Write(bp)
}

func (d *InfluxDB) tags(r record.Record) (map[string]string, error) {
	var (
		rk  = reflect.TypeOf(r).Kind()
		rv  = reflect.ValueOf(r)
		res = map[string]string{}
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

			v, ok = rv.MapIndex(key).Interface().(string)
			if !ok {
				return nil, NewErrUnexpectedMapKeyType(
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

func (d *InfluxDB) fields(r record.Record) (map[string]interface{}, error) {
	var (
		rk  = reflect.TypeOf(r).Kind()
		rv  = reflect.ValueOf(r)
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

func (d *InfluxDB) Create(r record.Record) error {
	var (
		p      *client.Point
		tags   map[string]string
		fields map[string]interface{}
		err    error
	)

	tags, err = d.tags(r)
	if err != nil {
		return err
	}

	fields, err = d.fields(r)
	if err != nil {
		return err
	}

	// FIXME: time
	p, err = client.NewPoint(
		d.config.Writer.Point.Name,
		tags,
		fields,
		time.Now(),
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

func New(c Config, cl client.Client, l loggers.Logger) (*InfluxDB, error) {
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
			config:     c,
			tagNames:   tagNames,
			fieldNames: fieldNames,
			client:     cl,
			log:        l,
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
