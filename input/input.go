package input

import (
	"context"
	"fmt"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"

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
		log = prefixwrapper.New(
			fmt.Sprintf("Input %s: ", i.config.Name),
			i.log,
		)
		c   databases.Connection
		d   databases.Database
		err error
	)

	c, err = databases.Connect(i.config.Database, log)
	if err != nil {
		return err
	}
	defer c.Close()

	d, err = databases.FromConfig(i.config.Database, c, log)
	if err != nil {
		return err
	}
	defer d.Close()

	err = consumer.PipeConsumerToDatabaseWith(
		i.config.Consumer,
		ctx,
		i.preprocessor.Preprocess,
		d,
		log,
	)
	if err != nil {
		return err
	}

	return nil
}

func (i *Input) Close() error {
	return i.preprocessor.Close()
}

func FromConfig(c Config, l loggers.Logger) (*Input, error) {
	var (
		p   preprocessors.Preprocessor
		err error
	)

	p, err = preprocessors.FromConfig(c.Preprocessor, l)
	if err != nil {
		return nil, err
	}

	return &Input{
		config:       c,
		log:          l,
		preprocessor: p,
	}, nil
}
