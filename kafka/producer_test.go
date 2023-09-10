package kafka

import (
	"github.com/IBM/sarama"
	"testing"
)

func TestProducerHanlder_CreateAsyncProducer(t *testing.T) {
	kafka := InitKafkaInstance()
	kafka.InitConnection()
	// Create producer handle
	producerHandler := ProducerHandler{
		SendChann: make(chan *sarama.ProducerMessage),
	}

	kafka.CreateAsyncProducer(&producerHandler)

	msg := &sarama.ProducerMessage{Topic: kafka.KafkaProperties.Producer.Topics[0], Key: nil,
		Value: sarama.StringEncoder("testing 123")}
	producerHandler.SendChann <- msg

}
