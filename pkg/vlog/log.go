package vlog

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
)

// This structure is used to declare the information of attributes utilized during logging.
// It is applied immediately when initializing the app's logging instance.
type LoggerConfig struct {
	Level          zerolog.Level
	TimeFormat     string // Time format set for log line
	IncludesCaller bool   //If enabled, the log line output will include the file information where the logging is executed.
}

func createZeroLogInst(c LoggerConfig) *zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Logger()
	logger = logger.With().Timestamp().Logger()
	if c.IncludesCaller {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			return filepath.Base(file) + ":" + strconv.Itoa(line)
		}
		logger = logger.With().Caller().Logger()
	}
	logger = logger.Level(c.Level)
	return &logger
}

func NewJsonLogger(c LoggerConfig) *zerolog.Logger {
	l := createZeroLogInst(c)
	return l
}
