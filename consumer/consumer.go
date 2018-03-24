package consumer

import (
	"context"

	"github.com/corpix/effects/closer"
	"github.com/corpix/formats"
	"github.com/corpix/loggers"

	"github.com/cryptounicorns/gluttony/databases"
	"github.com/cryptounicorns/queues"
	"github.com/cryptounicorns/queues/consumer"
	"github.com/cryptounicorns/queues/result"
)

type PrepareForDatabaseFn = func(v interface{}) (interface{}, error)

func PipeConsumerToDatabaseWith(
	c queues.GenericConfig,
	ctx context.Context,
	fn PrepareForDatabaseFn,
	d databases.Database,
	l loggers.Logger,
) error {
	var (
		closers = closer.Closers{}
		f       formats.Format
		q       queues.Queue
		cr      consumer.Consumer
		mcr     consumer.Generic
		err     error
	)

	go func() {
		<-ctx.Done()
		if err := closers.Close(); err != nil {
			l.Error(err)
		}
		return
	}()

	f, err = formats.New(c.Format)
	if err != nil {
		return err
	}

	q, err = queues.New(c.Queue, l)
	if err != nil {
		return err
	}
	closers = append(closers, q)

	cr, err = q.Consumer()
	if err != nil {
		return err
	}
	closers = append(closers, cr)

	mcr = consumer.NewUnmarshal(
		cr,
		new(interface{}),
		f,
	)
	closers = append(closers, mcr)

	var (
		stream <-chan result.Generic
		v      interface{}
	)

	stream, err = mcr.Consume()
	if err != nil {
		return err
	}

	for r := range stream {
		if r.Err != nil {
			return r.Err
		}

		// FIXME: pool?
		// FIXME: more flexible error handling?
		v, err = fn(r.Value)
		if err != nil {
			return err
		}

		err = d.Create(v)
		if err != nil {
			return err
		}
	}

	return nil
}
