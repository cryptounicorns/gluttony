package input

import (
	"github.com/cryptounicorns/gluttony/consumer"
	"github.com/cryptounicorns/gluttony/databases"
	"github.com/cryptounicorns/gluttony/preprocessors"
)

type Config struct {
	Name         string               `validate:"required"`
	Consumer     consumer.Config      `validate:"required,dive"`
	Preprocessor preprocessors.Config `validate:"required,dive"`
	Database     databases.Config     `validate:"required,dive"`
}
