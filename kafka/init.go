package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"golibs/log"
	"os"
)

type Kafka struct {
	KafkaProperties   *Properties
	KafkaSaramaConfig *sarama.Config
}

func (k *Kafka) createClient() {
	broker := sarama.NewBroker(k.KafkaProperties.Brokers[0])
	err := broker.Open(k.KafkaSaramaConfig)
	if err != nil {
		log.Logger.Error().Msg(err.Error())
		os.Exit(1)
	}

	request := sarama.MetadataRequest{}
	response, err := broker.GetMetadata(&request)
	if err != nil {
		_ = broker.Close()
		log.Logger.Error().Msg(err.Error())
		os.Exit(1)
	}
	log.Logger.Debug().Msgf("there are %d topics active in the cluster.", len(response.Topics))
	// Get list topic in this cluster
	var topicsActive []string
	for _, topic := range response.Topics {
		topicDetail := fmt.Sprintf("%s:%d", topic.Name, len(topic.Partitions))
		topicsActive = append(topicsActive, topicDetail)
	}
	log.Logger.Debug().Msgf("list topic %v", topicsActive)
	if err = broker.Close(); err != nil {
		log.Logger.Error().Msg(err.Error())
		os.Exit(1)
	}
}

func (k *Kafka) createConfig() {
	config := sarama.NewConfig()
	// Setup group rebalance
	switch k.KafkaProperties.Producer.GroupReBalance {
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
			k.KafkaProperties.Producer.GroupReBalance)
		os.Exit(1)
	}
	// Setup offset
	config.Consumer.Offsets.Initial = k.KafkaProperties.Consumer.StartOffset
	config.Consumer.Return.Errors = true
	k.KafkaSaramaConfig = config
}

func (k *Kafka) InitConnection() {
	log.Logger.Debug().Msgf("start create connection to kafka")
	k.createConfig()
	k.createClient()
	// Set config to class
}
