package example

import (
	"github.com/IBM/sarama"
	"golibs/kafka"
	"golibs/log"
)

func InitKafkaInstance() *kafka.Kafka {
	log.InitLogger(&log.Properties{
		Level:          0,
		PrefixFieldLog: "kafka",
	})
	instance := kafka.Kafka{
		KafkaProperties: &kafka.Properties{
			Brokers: []string{"kafka.local:9192"},
			Producer: struct {
				Topics         []string
				GroupReBalance string
			}{Topics: []string{"topic-input"}, GroupReBalance: "roundrobin"},
			Consumer: struct {
				Topics        []string
				StartOffset   int64
				ConsumerGroup string
			}{Topics: []string{"topic-input"}, StartOffset: sarama.OffsetOldest, ConsumerGroup: "group_test"},
		},
	}
	return &instance
}
