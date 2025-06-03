package consumer

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

var (
	defaultRetryAmount   = 10
	defaultRetryInterval = time.Second * 3
)

// If KafkaCfg nil, using default sarama config
type Config struct {
	Brokers  []string
	KafkaCfg *sarama.Config

	retryAmount   int
	retryInterval time.Duration
}

type Consumer struct {
	group_id string
	conn     sarama.ConsumerGroup
	cfg      *Config
}

func New(cfg *Config, opts ...Opt) (*Consumer, error) {
	cfg.retryAmount = defaultRetryAmount
	cfg.retryInterval = defaultRetryInterval

	if cfg.KafkaCfg == nil {
		cfg.KafkaCfg = sarama.NewConfig()
	}

	consumerGroup := &Consumer{group_id: uuid.NewString(), cfg: cfg}

	for _, opt := range opts {
		opt(consumerGroup)
	}

	err := consumerGroup.connect()
	if err != nil {
		return nil, fmt.Errorf("consumerGroup.connect: %w", err)
	}

	return consumerGroup, nil
}

func (cg *Consumer) connect() error {
	var err error
	for i := cg.cfg.retryAmount; i > 0; i-- {
		cg.conn, err = sarama.NewConsumerGroup(cg.cfg.Brokers, cg.group_id, cg.cfg.KafkaCfg)
		if err == nil {
			return nil
		}

		time.Sleep(cg.cfg.retryInterval)
	}

	return fmt.Errorf("sarama.NewConsumerGroup: %w", err)
}
