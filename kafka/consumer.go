package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"golibs/log"
	"os"
	"sync"
)

type ConsumerGroupHandler struct {
	Kafka
	ready       chan bool
	messageChan chan *sarama.ConsumerMessage
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	log.Logger.Debug().Msg("setup consumer is ready.")

	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Logger.Debug().Msgf("Message topic:%q partition:%d offset:%d", msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
		consumer.messageChan <- msg
	}
	return nil
}

func (k *Kafka) CreateConsumer(handler *ConsumerGroupHandler) {
	log.Logger.Debug().Msg("start creating consumer...")
	ctx := context.Background()
	client, err := sarama.NewConsumerGroup(k.KafkaProperties.Brokers, k.KafkaProperties.Consumer.ConsumerGroup, k.KafkaSaramaConfig)
	if err != nil {
		panic(err)
	}
	topics := k.KafkaProperties.Consumer.Topics

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, topics, handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Logger.Error().Msg(err.Error())
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			handler.ready = make(chan bool)
		}
	}()
	<-handler.ready
	log.Logger.Info().Msg("sarama consumer up and running...")
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Logger.Error().Msgf("Error closing client: %v", err)
		os.Exit(1)
	}
}
