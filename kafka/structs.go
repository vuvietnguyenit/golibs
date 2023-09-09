package kafka

type Properties struct {
	Brokers []string
	Options struct {
		GroupReBalance string
		StartOffset    int64
	}
	ConsumerGroup string
}
