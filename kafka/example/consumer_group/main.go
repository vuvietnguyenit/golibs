package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"golibs/kafka/example"
	"golibs/log"
	"sync"
)

type ConsumerGroupHandler struct {
	ready    chan bool
	totalMsg int
}

// Setup Implement all function for ConsumerGroupHandler
func (consumer *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Logger.Warn().Msg("message channel was closed")
				return nil
			}
			log.Logger.Info().Msgf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")
			consumer.totalMsg++
			fmt.Println("Total msg: ", consumer.totalMsg)
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func main() {
	kafkaInstance := example.InitKafkaInstance()
	kafkaInstance.InitConnection()
	handler := kafkaInstance.CreateConsumerGroup()
	consumer := ConsumerGroupHandler{
		ready: make(chan bool),
	}
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := handler.Consume(ctx, kafkaInstance.KafkaProperties.Consumer.Topics, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Logger.Error().Msgf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Logger.Info().Msg("Sarama consumer up and running!...")
	wg.Wait()
	if err := handler.Close(); err != nil {
		log.Logger.Error().Msgf("Error closing client: %v", err)
	}
}
