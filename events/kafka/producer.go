package kafka

import "github.com/IBM/sarama"

type KafkaProducer struct {
	Producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{Producer: producer}, nil
}

func (kp *KafkaProducer) SendMessage(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := kp.Producer.SendMessage(msg)
	return err
}
