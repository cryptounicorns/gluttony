package stores

import (
	"fmt"
	"strings"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"

	"github.com/corpix/stores/store/memory"
	"github.com/corpix/stores/store/memoryttl"
)

func New(c Config, l loggers.Logger) (Store, error) {
	var (
		t   = strings.ToLower(c.Type)
		log = prefixwrapper.New(
			fmt.Sprintf(
				"Store(%s): ",
				t,
			),
			l,
		)
	)

	switch t {
	case memory.Name:
		return memory.New(
			c.Memory,
			log,
		)
	case memoryttl.Name:
		return memoryttl.New(
			c.MemoryTTL,
			log,
		)
	default:
		return nil, NewErrUnknownStoreType(c.Type)
	}
}
