package none

import (
	"github.com/corpix/loggers"
)

const (
	Name = "none"
)

type None struct {
	config Config
	log    loggers.Logger
}

func (l *None) Preprocess(v interface{}) (interface{}, error) {
	return v, nil
}

func (l *None) Close() error {
	l.log.Debug("Closing")

	return nil
}

func New(c Config, l loggers.Logger) (*None, error) {
	return &None{
		config: c,
		log:    l,
	}, nil
}
