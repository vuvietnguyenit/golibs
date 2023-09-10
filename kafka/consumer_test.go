package kafka

import (
	"github.com/IBM/sarama"
	"golibs/log"
	"sync"
	"testing"
)

func TestKafka_CreateConsumer(t *testing.T) {
	kafka := InitKafkaInstance()
	kafka.InitConnection()

	// Create consumer group handle
	consumerGroupHandler := ConsumerGroupHandler{
		ready:       make(chan bool),
		messageChan: make(chan *sarama.ConsumerMessage),
	}
	// Create async producer chan
	// channel here is received message for produce topic
	producer := kafka.CreateAsyncProducer()

	// Get message consumed
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// go routine for init consumer
		kafka.CreateConsumer(&consumerGroupHandler)
	}()
	wg.Add(1)
	go func() {
		// go routine for get message from channel
		log.Logger.Debug().Msg("read message from consumer channel")
		for {
			select {
			case msg := <-consumerGroupHandler.messageChan:
				// create message for producer
				messageProduce := sarama.ProducerMessage{
					Topic: kafka.KafkaProperties.Producer.Topics[0],
					Value: sarama.StringEncoder(msg.Value),
				}
				// step produce message to topic
				producer.Input() <- &messageProduce
			}
		}
	}()
	// create goroutine observer message producer
	wg.Add(1)
	go func() {
		kafka.AsyncProducerObserver(producer)
	}()
	wg.Wait()

}
