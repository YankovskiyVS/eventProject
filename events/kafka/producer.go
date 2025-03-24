package kafka

import (
	"github.com/IBM/sarama"
	transactionaloutbox "github.com/YankovskiyVS/eventProject/events/internal/transactional_outbox"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll // Better durability

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
	}, nil
}

func (kp *KafkaProducer) SendMessage(event transactionaloutbox.Message) error {
	headers := make([]sarama.RecordHeader, 0, len(event.Headers))

	for k, v := range event.Headers {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	msg := &sarama.ProducerMessage{
		Topic:   event.Topic,
		Key:     sarama.StringEncoder(event.Key),
		Value:   sarama.StringEncoder(event.Body),
		Headers: headers,
	}

	_, _, err := kp.producer.SendMessage(msg)
	return err
}

func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}
