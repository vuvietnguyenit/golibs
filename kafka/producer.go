package kafka

import (
	"github.com/IBM/sarama"
	"golibs/log"
	"os"
	"os/signal"
	"sync"
)

type ProducerHandler struct {
	SendChann chan *sarama.ProducerMessage
}

var (
	wg                                  sync.WaitGroup
	enqueued, successes, producerErrors int
)

func (k *Kafka) CreateAsyncProducer(producerHandle *ProducerHandler) {
	producer, err := sarama.NewAsyncProducer(k.KafkaProperties.Brokers, k.KafkaSaramaConfig)
	if err != nil {
		log.Logger.Error().Msg(err.Error())
		os.Exit(1)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	wg.Add(1)
	go func() {
		log.Logger.Debug().Msg("read success goroutine")
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()
	wg.Add(1)
	go func() {
		log.Logger.Debug().Msg("read error goroutine")
		defer wg.Done()
		for err := range producer.Errors() {
			log.Logger.Error().Msg(err.Error())
			producerErrors++
		}
	}()
ProducerLoop:
	for {
		message := &sarama.ProducerMessage{Topic: k.KafkaProperties.Producer.Topics[0], Value: sarama.StringEncoder("testing 123")}
		select {
		case producer.Input() <- message:
			enqueued++

		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			break ProducerLoop
		}
	}

	wg.Wait()

	log.Logger.Info().Msgf("Successfully produced: %d; errors: %d", successes, producerErrors)

}

func (k *Kafka) CreateProducerChannel() {

}
