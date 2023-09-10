package kafka

import (
	"github.com/IBM/sarama"
	uuid2 "github.com/google/uuid"
	"sync"
	"testing"
	"time"
)

func TestProducerHanlder_CreateAsyncProducer(t *testing.T) {
	kafka := InitKafkaInstance()
	kafka.InitConnection()
	// Create producer handle

	// Create async producer
	producer := kafka.CreateAsyncProducer()
	// create goroutine observer message producer
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		kafka.AsyncProducerObserver(producer)
	}()
	go func() {
		for {
			uuid := uuid2.New()
			messageProduce := sarama.ProducerMessage{
				Topic: kafka.KafkaProperties.Producer.Topics[0],
				Value: sarama.StringEncoder(uuid.String()),
			}
			producer.Input() <- &messageProduce
			time.Sleep(100 * time.Millisecond)

		}
	}()
	wg.Wait()

}
