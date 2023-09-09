package kafka

import (
	"github.com/IBM/sarama"
	"golibs/log"
)

type Kafka struct {
	KafkaProperties   *Properties
	KafkaSaramaConfig *sarama.Config
}

func (k *Kafka) createConfig() *sarama.Config {
	config := sarama.NewConfig()
	// Setup group rebalance
	switch k.KafkaProperties.Options.GroupReBalance {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
		log.Logger.Debug().Msgf("set kafka group rebalance sticky.")
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
		log.Logger.Debug().Msgf("set kafka group rebalance roundrobin.")
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
		log.Logger.Debug().Msgf("set kafka group rebalance range.")
	default:
		log.Logger.Error().Msgf("Unrecognized consumer group partition assignor: %s",
			k.KafkaProperties.Options.GroupReBalance)
	}
	// Setup offset
	config.Consumer.Offsets.Initial = k.KafkaProperties.Options.StartOffset
	return config
}

func (k *Kafka) InitConnection() {
	config := k.createConfig()
	// Set config to class
	k.KafkaSaramaConfig = config
}
