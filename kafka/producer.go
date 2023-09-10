package kafka

import (
	"github.com/IBM/sarama"
	"golibs/log"
	"os"
)

type ProducerHandler struct {
	SendChann chan *sarama.ProducerMessage
}
type AsyncProducer interface {
	CreateAsyncProducer()
	AsyncProducerObserver()
}

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
	for {
		select {
		case err := <-producer.Errors():
			log.Logger.Error().Msgf("Failed to produce message %v", err)
		case succ := <-producer.Successes():
			log.Logger.Info().Msgf("produce success [topic|key|offset]: [%s|%v|%d]", succ.Topic, succ.Key, succ.Offset)
		}
	}
}

func (k *Kafka) CreateProducerChannel() chan<- *sarama.ProducerMessage {
	producerChan := k.CreateAsyncProducer()
	return producerChan.Input()
}
