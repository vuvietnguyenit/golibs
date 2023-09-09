package kafka

import (
	"github.com/IBM/sarama"
	"golibs/log"
	"testing"
)

func InitKafkaInstance() *Kafka {
	log.InitLogger(&log.Properties{Level: 0})
	kafka := Kafka{
		KafkaProperties: &Properties{
			Brokers: []string{"localhost:9092"},
			Options: struct {
				GroupReBalance string
				StartOffset    int64
			}{GroupReBalance: "roundrobin", StartOffset: sarama.OffsetOldest},
			ConsumerGroup: "group-test",
		},
	}
	return &kafka
}

func TestInitConnection(t *testing.T) {
	kafka := InitKafkaInstance()
	kafka.InitConnection()

}
