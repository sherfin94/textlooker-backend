package kafka

import (
	"encoding/json"
	"log"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type Text struct {
	ID       string   `json:"id"`
	Content  string   `json:"content"`
	Author   []string `json:"author"`
	Date     string   `json:"date"`
	SourceID int      `json:"source_id"`
}

var TextProcessChannel *chan Text

func InitializeProducer(channel *chan Text) {
	TextProcessChannel = channel
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}

	defer producer.Close()

	topic := "textlooker"
	for text := range *channel {
		if serializedText, err := json.Marshal(text); err == nil {
			producer.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(serializedText),
			}, nil)
		} else {
			log.Println(err)
		}
	}

	producer.Flush(15 * 1000)
}
