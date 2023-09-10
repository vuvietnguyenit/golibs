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
	// Get message consumed
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		// go routine for init consumer
		kafka.CreateConsumer(&consumerGroupHandler)
		wg.Done()
	}()
	go func() {
		// go routine for get message from channel
		log.Logger.Debug().Msg("read message from consumer channel")
		for {
			select {
			case msg := <-consumerGroupHandler.messageChan:
				log.Logger.Info().Msg(string(msg.Value))
			}
		}
	}()
	wg.Wait()

}
