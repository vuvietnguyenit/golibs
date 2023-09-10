package kafka

type Properties struct {
	Brokers  []string
	Producer struct {
		Topics         []string
		GroupReBalance string
	}
	Consumer struct {
		Topics        []string
		StartOffset   int64
		ConsumerGroup string
	}
}
