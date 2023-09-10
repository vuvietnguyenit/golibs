package kafka

import (
	"github.com/IBM/sarama"
	"golibs/log"
	"testing"
)

func InitKafkaInstance() *Kafka {
	log.InitLogger(&log.Properties{
		Level:          0,
		PrefixFieldLog: "kafka",
	})
	kafka := Kafka{
		KafkaProperties: &Properties{
			Brokers: []string{"kafka.local:9192"},
			Producer: struct {
				Topics         []string
				GroupReBalance string
			}{Topics: []string{"sample"}, GroupReBalance: "roundrobin"},
			Consumer: struct {
				Topics        []string
				StartOffset   int64
				ConsumerGroup string
			}{Topics: []string{"sample"}, StartOffset: sarama.OffsetOldest, ConsumerGroup: "group_test"},
		},
	}
	return &kafka
}

func TestInitConnection(t *testing.T) {
	kafka := InitKafkaInstance()
	kafka.InitConnection()
}
