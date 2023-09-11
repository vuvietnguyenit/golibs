package kafka

import (
	"github.com/IBM/sarama"
	"golibs/log"
	"os"
	"sync"
)

func (k *Kafka) CreateAsyncProducer() sarama.AsyncProducer {
	k.KafkaSaramaConfig.Producer.Return.Successes = true
	k.KafkaSaramaConfig.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(k.KafkaProperties.Brokers, k.KafkaSaramaConfig)
	if err != nil {
		log.Logger.Error().Msg(err.Error())
		os.Exit(1)
	}
	return producer
}
func (k *Kafka) AsyncProducerObserver(producer sarama.AsyncProducer) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Logger.Debug().Msg("start thread for observer msg produce")
		for {
			select {
			case err := <-producer.Errors():
				log.Logger.Error().Msgf("failed to produce message %v", err)
			case succ := <-producer.Successes():
				log.Logger.Info().Msgf("produce success [topic|key|offset]: [%s|%v|%d]", succ.Topic, succ.Key, succ.Offset)
			}
		}
	}()
}

func (k *Kafka) CreateProducerChannel() chan<- *sarama.ProducerMessage {
	producerChan := k.CreateAsyncProducer()
	return producerChan.Input()
}
