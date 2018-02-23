package input

import (
	"github.com/cryptounicorns/gluttony/consumer"
	"github.com/cryptounicorns/gluttony/databases"
	"github.com/cryptounicorns/gluttony/preprocessors"
)

type Config struct {
	Name         string               `validate:"required"`
	Consumer     consumer.Config      `validate:"required"`
	Preprocessor preprocessors.Config `validate:"required"`
	Database     databases.Config     `validate:"required"`
}
