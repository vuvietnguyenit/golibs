package vlog_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/vuvietnguyenit/golibs/pkg/vlog"
)

func ExampleNewJsonLogger() {
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.InfoLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	jsonLogger.Info().Msg("This is INFO message")
	jsonLogger.Debug().Msg("This is DEBUG message")

	// {"level":"info","time":"2023-06-16T07:46:47Z07:00","message":"This is INFO message"}

	// Change Level logger to lower level if you want to show DEBUG level
	jsonLogger = vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.DebugLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})

	// {"level":"info","time":"2025-02-06T18:55:31+07:00","caller":"log_test.go:34","message":"This is INFO message"}
	// {"level":"debug","time":"2025-02-06T18:55:31+07:00","caller":"log_test.go:35","message":"This is DEBUG message"}
}

func TestNewJsonLogger(t *testing.T) {
	jsonLogger := vlog.NewJsonLogger(vlog.LoggerConfig{
		Level:          zerolog.DebugLevel,
		TimeFormat:     "2006-01-02T15:04:05Z07:00",
		IncludesCaller: true,
	})
	jsonLogger.Info().Msg("This is INFO message")
	jsonLogger.Debug().Msg("This is DEBUG message")

	//
}
