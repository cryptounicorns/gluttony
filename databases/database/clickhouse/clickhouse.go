package clickhouse

import (
	"database/sql/driver"
	"time"

	"github.com/corpix/loggers"
	"github.com/cryptounicorns/gluttony/databases/record"
	"github.com/kshvakov/clickhouse"
)

const (
	Name = "clickhouse"

	defaultBatchWriteTime = 5 * time.Second
)

type Clickhouse struct {
	client clickhouse.Clickhouse
	config Config
	log    loggers.Logger

	// timestampField reflect.Value

	points chan interface{}
	done   chan struct{}
}

func FromConfig(cfg Config, client clickhouse.Clickhouse, log loggers.Logger) (*Clickhouse, error) {
	{
		tx, _ := client.Begin()
		stmt, _ := client.Prepare(dbCreate)
		if _, err := stmt.Exec([]driver.Value{}); err != nil {
			log.Fatal(err)
		}
		tx.Commit()
	}

	db := &Clickhouse{
		client: client,
		config: cfg,
		log:    log,
		// timestampField: reflect.ValueOf(c.Writer.Point.Timestamp),

		// points: make(
		// 	chan *client.Point,
		// 	c.Writer.Batch.Size,
		// ),
		done: make(chan struct{}),
	}

	go db.batchWriter()

	return db, nil
}

func (db *Clickhouse) Close() error {
	close(db.done)
	return db.client.Close()
}

func (db *Clickhouse) Create(r record.Record) error {
	tx, _ := db.client.Begin()
	stmt, _ := db.client.Prepare(dbInsert)
	for i := 0; i < 100; i++ {
		if _, err := stmt.Exec([]driver.Value{
			time.Now().Unix(),
			uint64(0),
			uint64(0),
			uint64(0),
			uint64(0),
			uint64(0),
		}); err != nil {
			db.log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		db.log.Fatal(err)
	}

	return nil
}

func (db *Clickhouse) batchWriter() {
	interval := db.config.FlushInterval
	if interval == 0 {
		interval = defaultBatchWriteTime
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-db.done:
			return

		case <-ticker.C:
			db.writeBatch()
		}
	}
}

func (db *Clickhouse) writeBatch() error {
	db.log.Debug("Preparing batch...")

	return nil
}
