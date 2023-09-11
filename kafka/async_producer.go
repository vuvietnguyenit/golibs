package kafka

import (
	"github.com/IBM/sarama"
	"golibs/log"
	"os"
	"sync"
)

type AsyncProducer struct {
	KafkaInstance *Kafka
}

func (k *AsyncProducer) Create() sarama.AsyncProducer {
	k.KafkaInstance.KafkaSaramaConfig.Producer.Return.Successes = true
	k.KafkaInstance.KafkaSaramaConfig.Producer.Return.Errors = true
	producer, err := sarama.NewAsyncProducer(k.KafkaInstance.KafkaProperties.Brokers, k.KafkaInstance.KafkaSaramaConfig)
	if err != nil {
		log.Logger.Error().Msg(err.Error())
		os.Exit(1)
	}
	return producer
}
func (k *AsyncProducer) Observer(producer sarama.AsyncProducer) {
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
