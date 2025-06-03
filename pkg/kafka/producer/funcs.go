package producer

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

func (p *Producer) Send(topic string, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	_, _, err = p.conn.SendMessage(
		&sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.ByteEncoder(data),
		},
	)

	if err != nil {
		return fmt.Errorf("p.conn.SendMessage: %w", err)
	}

	return nil
}
