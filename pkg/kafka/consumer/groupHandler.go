package consumer

import (
	"log"

	"github.com/IBM/sarama"
)

type consumerGroupHandler struct {
	handler MessageHandler
}

type MessageHandler func(msg *sarama.ConsumerMessage) error

func newConsumerGroupHandler(handler MessageHandler) *consumerGroupHandler {
	return &consumerGroupHandler{handler: handler}
}

func (cgh *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("Hello im consumer %s", session.MemberID())
	return nil
}

func (cgh *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Printf("Bye im consumer %s", session.MemberID())
	session.Commit()
	return nil
}

func (cgh *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := cgh.handler(msg)
		if err != nil {
			log.Println(err)
			continue
		}
		session.MarkMessage(msg, "")
	}

	return nil
}
