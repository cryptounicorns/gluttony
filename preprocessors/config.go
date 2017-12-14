package preprocessors

import (
	"github.com/cryptounicorns/gluttony/preprocessors/preprocessor/lua"
)

type Config struct {
	Type string     `validator:"required"`
	Lua  lua.Config `validator:"required,dive"`
}
