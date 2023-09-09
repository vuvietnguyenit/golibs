package kafka

import (
	"testing"
)

func TestKafka_CreateConsumer(t *testing.T) {
	kafka := InitKafkaInstance()
	kafka.InitConnection()
	kafka.CreateConsumer()
}
