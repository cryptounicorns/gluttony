package stores

import (
	"github.com/corpix/stores/store/memory"
	"github.com/corpix/stores/store/memoryttl"
)

type Config struct {
	Type      string
	Memory    memory.Config
	MemoryTTL memoryttl.Config
}
