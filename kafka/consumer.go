package kafka

import (
	"bluebell_backend/dao/mysql"
	"bluebell_backend/models"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

var ReaderForP0 *kafka.Reader
var ReaderForP1 *kafka.Reader

func InitKafkaConsumer() {
	ReaderForP0 = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Topic:    "ccmoblog-user-track",
		GroupID:  "ccmoblog-group-user-track",
		MaxBytes: 10e6, // 10MB
	})
	ReaderForP1 = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		Topic:    "ccmoblog-user-track",
		GroupID:  "ccmoblog-group-user-track",
		MaxBytes: 10e6, // 10MB
	})
}

func Consume(reader *kafka.Reader) {
	for {
		msg, err := reader.FetchMessage(context.Background())
		if err != nil {
			continue
		}
		value := msg.Value
		// 反序列化
		var data models.UserTracking
		json.Unmarshal(value, &data)

		// 数据库同步，同步失败不提交偏移量
		err = mysql.SyncData(data)
		if err != nil {
			continue
		}

		reader.CommitMessages(context.Background(), msg)
	}
}
