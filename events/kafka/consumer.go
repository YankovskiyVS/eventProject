package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	Consumer sarama.Consumer
}

func NewKafkaConsumer(brokers []string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{Consumer: consumer}, nil
}

func (kc *KafkaConsumer) ConsumeMessages(topic string) {
	partitionConsumer, err := kc.Consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received message: %s\n", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v\n", err)
		}
	}
}
