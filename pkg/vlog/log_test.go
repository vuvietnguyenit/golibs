package vlog

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestNewJsonLogger(t *testing.T) {
	jsonLogger := NewJsonLogger(LoggerConfig{
		Level:          zerolog.DebugLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	jsonLogger.LegacyHandler.Info().Msg("INFO")
	jsonLogger.LegacyHandler.Debug().Msg("DEBUG")
}
