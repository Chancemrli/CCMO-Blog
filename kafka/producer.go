package kafka

import (
	"github.com/segmentio/kafka-go"
)

var KWriter *kafka.Writer

func InitKafkaProducer() {
	KWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Topic:        "ccmoblog-user-track",
		Balancer:     &kafka.LeastBytes{},
		Async:        true,
		RequiredAcks: int(kafka.RequireOne),
	})
}

func Close() {
	defer KWriter.Close()
	defer ReaderForP0.Close()
	defer ReaderForP1.Close()
}
