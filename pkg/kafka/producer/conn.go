package producer

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
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

type Producer struct {
	conn sarama.SyncProducer
	cfg  *Config
}

func New(cfg *Config, opts ...Opt) (*Producer, error) {
	cfg.retryAmount = defaultRetryAmount
	cfg.retryInterval = defaultRetryInterval

	if cfg.KafkaCfg != nil {
		cfg.KafkaCfg.Producer.Return.Successes = true
	}

	producer := &Producer{cfg: cfg}

	for _, opt := range opts {
		opt(producer)
	}

	err := producer.connect()
	if err != nil {
		return nil, fmt.Errorf("producer.connect: %w", err)
	}

	return producer, nil
}

func (p *Producer) connect() error {
	var err error
	for i := p.cfg.retryAmount; i > 0; i-- {
		p.conn, err = sarama.NewSyncProducer(p.cfg.Brokers, p.cfg.KafkaCfg)
		if _, ok := err.(sarama.ConfigurationError); ok {
			return fmt.Errorf("sarama.NewSyncProducer: %w", err)
		}
		if err == nil {
			return nil
		}

		time.Sleep(p.cfg.retryInterval)
	}

	return fmt.Errorf("sarama.NewSyncProducer: %w", err)
}
