package preprocessors

import (
	"github.com/cryptounicorns/gluttony/preprocessors/preprocessor/lua"
	"github.com/cryptounicorns/gluttony/preprocessors/preprocessor/none"
)

type Config struct {
	Type string `validate:"required"`
	Lua  *lua.Config
	None *none.Config
}
