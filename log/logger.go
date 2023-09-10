package log

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

func InitLogger(loggerProperties *Properties) {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return fmt.Sprintf("%s", strings.ToUpper(fmt.Sprintf("%s", i)))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s_msg=\"%s\"", loggerProperties.PrefixFieldLog, fmt.Sprintf("%s", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s_%s=", loggerProperties.PrefixFieldLog, i)
	}
	log := zerolog.New(output).
		Level(loggerProperties.Level).
		With().
		Caller().
		Timestamp().
		Logger()
	Logger = &log
}
