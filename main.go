package main

import (
	"github.com/IBM/sarama"
	"golibs/kafka"
	"golibs/log"
)

func main() {
	log.InitLogger(&log.Properties{Level: 0})
	kafka := kafka.Kafka{
		KafkaProperties: &kafka.Properties{
			Brokers: []string{"kafka.local:9192"},
			Producer: struct {
				Topics         []string
				GroupReBalance string
			}{Topics: []string{"producer-topic"}, GroupReBalance: "roundrobin"},
			Consumer: struct {
				Topics        []string
				StartOffset   int64
				ConsumerGroup string
			}{Topics: []string{"sample"}, StartOffset: sarama.OffsetOldest, ConsumerGroup: "group_test"},
		},
	}
	kafka.InitConnection()
}
