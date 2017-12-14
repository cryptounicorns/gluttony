package preprocessors

import (
	"fmt"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"

	"github.com/cryptounicorns/gluttony/preprocessors/preprocessor/lua"
)

type Preprocessor interface {
	Preprocess(interface{}) (interface{}, error)
	Close() error
}

func New(c Config, l loggers.Logger) (Preprocessor, error) {
	var (
		t   = strings.ToLower(c.Type)
		log = prefixwrapper.New(
			fmt.Sprintf("Preprocessor(%s): ", t),
			l,
		)
	)

	switch t {
	case lua.Name:
		return lua.New(c.Lua, log)
	default:
		return nil, NewErrUnknownPreprocessorType(c.Type)
	}
}
