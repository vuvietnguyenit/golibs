package vlog

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
)

type LoggerConfig struct {
	Level              zerolog.Level
	TimeFormat         string
	IncludesCaller     bool
	EnableHTTPTraceLog bool
}

type JsonLogger struct {
	LegacyHandler *zerolog.Logger
}

func createZeroLogInst(c LoggerConfig) zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Logger()
	logger = logger.With().Timestamp().Logger()
	if c.IncludesCaller {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return filepath.Base(file) + ":" + strconv.Itoa(line)
		}
		logger = logger.With().Caller().Logger()
	}
	logger = logger.Level(c.Level)
	return logger
}

func NewJsonLogger(c LoggerConfig) JsonLogger {
	l := createZeroLogInst(c)
	// Create HTTP handlers

	return JsonLogger{
		LegacyHandler: &l,
	}
}
