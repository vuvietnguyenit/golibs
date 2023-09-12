package kafka

import (
	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"os"
)

type SyncProducer struct {
	kafkaInstance *Kafka
}

func (s *SyncProducer) Create() *sarama.SyncProducer {
	producer, err := sarama.NewSyncProducer(s.kafkaInstance.KafkaProperties.Brokers, s.kafkaInstance.KafkaSaramaConfig)
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Logger.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()
	return &producer
}
