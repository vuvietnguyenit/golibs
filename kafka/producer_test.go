package kafka

import (
	"testing"
)

func TestProducerHanlder_CreateAsyncProducer(t *testing.T) {
	kafka := InitKafkaInstance()
	kafka.InitConnection()
	// Create producer handle

}
