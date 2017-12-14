package input

import (
	"github.com/cryptounicorns/gluttony/consumer"
	"github.com/cryptounicorns/gluttony/databases"
	"github.com/cryptounicorns/gluttony/preprocessors"
)

type Config struct {
	Consumer     consumer.Config      `validator:"required,dive"`
	Preprocessor preprocessors.Config `validator:"required,dive"`
	Database     databases.Config     `validator:"required,dive"`
}
