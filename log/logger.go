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
		return fmt.Sprintf("vti_msg=\"%s\"", fmt.Sprintf("%s", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("vti_%s=", i)
	}
	log := zerolog.New(output).
		Level(loggerProperties.Level).
		With().
		Caller().
		Timestamp().
		Logger()
	Logger = &log
}
