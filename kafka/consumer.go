package kafka

import (
	"github.com/IBM/sarama"
	"github.com/vuvietnguyenit/golibs/log"
	"os"
)

func (k *Kafka) CreateConsumerGroup() sarama.ConsumerGroup {
	log.Logger.Debug().Msg("start creating consumer...")
	client, err := sarama.NewConsumerGroup(k.KafkaProperties.Brokers, k.KafkaProperties.Consumer.ConsumerGroup, k.KafkaSaramaConfig)
	if err != nil {
		log.Logger.Error().Msg(err.Error())
		os.Exit(1)
	}
	return client
}
