package main

import (
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"golibs/kafka/example"
	"time"
)

func main() {
	kafka := example.InitKafkaInstance()
	kafka.InitConnection()
	producer := kafka.CreateAsyncProducer()
	kafka.AsyncProducerObserver(producer)
	for {
		uuidIn := uuid.New()
		messageProduce := sarama.ProducerMessage{
			Topic: kafka.KafkaProperties.Producer.Topics[0],
			Value: sarama.StringEncoder(uuidIn.String()),
		}
		producer.Input() <- &messageProduce
		time.Sleep(100 * time.Millisecond)
	}
}
