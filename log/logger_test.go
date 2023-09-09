package log

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func TestLogger(t *testing.T) {
	InitLogger(&Properties{
		Level: 0,
	})
	log.Logger.Debug().Str("connnection", "conn ok").Int("duration", 1000).Msg("Hello world")
	log.Logger.Info().Int("duration", 1000).Msg("Done")
	log.Logger.Error().Int("error_code", 500).Msg("Connect to service failed")
}
