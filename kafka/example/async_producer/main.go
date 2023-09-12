package main

import (
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/vuvietnguyenit/golibs/kafka"
	"github.com/vuvietnguyenit/golibs/kafka/example"
	"time"
)

func main() {
	kafkaInstance := example.InitKafkaInstance()
	kafkaInstance.InitConnection()
	newProducer := kafka.AsyncProducer{
		KafkaInstance: kafkaInstance,
	}
	asyncProducer := newProducer.Create()
	newProducer.Observer(asyncProducer)
	for {
		uuidIn := uuid.New()
		messageProduce := sarama.ProducerMessage{
			Topic: kafkaInstance.KafkaProperties.Producer.Topics[0],
			Value: sarama.StringEncoder(uuidIn.String()),
		}
		asyncProducer.Input() <- &messageProduce
		time.Sleep(100 * time.Millisecond)
	}
}
