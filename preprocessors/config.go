package preprocessors

import (
	"github.com/cryptounicorns/gluttony/preprocessors/preprocessor/lua"
	"github.com/cryptounicorns/gluttony/preprocessors/preprocessor/none"
)

type Config struct {
	Type string      `validator:"required"`
	Lua  lua.Config  `validator:"dive"`
	None none.Config `validator:"dive"`
}
