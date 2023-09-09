package golibs

import (
	"golibs/kafka"
	"golibs/log"
	"testing"
)

func TestEntrypoint(t *testing.T) {
	// init logger first
	log.InitLogger(&log.Properties{Level: 0})
	kafka.CreateConsumer()
}
