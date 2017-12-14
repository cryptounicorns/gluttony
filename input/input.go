package input

import (
	"context"

	"github.com/corpix/loggers"

	"github.com/cryptounicorns/gluttony/consumer"
	"github.com/cryptounicorns/gluttony/databases"
	"github.com/cryptounicorns/gluttony/preprocessors"
)

type Input struct {
	config       Config
	log          loggers.Logger
	preprocessor preprocessors.Preprocessor
}

func (i *Input) Run(ctx context.Context) error {
	var (
		c   databases.Connection
		d   databases.Database
		err error
	)

	c, err = databases.Connect(i.config.Database, i.log)
	if err != nil {
		return err
	}
	defer c.Close()

	d, err = databases.New(i.config.Database, c, i.log)
	if err != nil {
		return err
	}
	defer d.Close()

	return consumer.PipeConsumerToDatabaseWith(
		i.config.Consumer,
		ctx,
		i.preprocessor.Preprocess,
		d,
		i.log,
	)
}

func (i *Input) Close() error {
	return i.preprocessor.Close()
}

func New(c Config, l loggers.Logger) (*Input, error) {
	var (
		p   preprocessors.Preprocessor
		err error
	)

	p, err = preprocessors.New(c.Preprocessor, l)
	if err != nil {
		return nil, err
	}

	return &Input{
		config:       c,
		log:          l,
		preprocessor: p,
	}, nil
}
