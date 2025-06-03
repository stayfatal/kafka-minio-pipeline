package producer

import (
	"time"
)

type Opt func(*Producer)

func SetRetryAmount(retryAmount int) Opt {
	return func(p *Producer) {
		p.cfg.retryAmount = retryAmount
	}
}

func SetRetryInterval(interval time.Duration) Opt {
	return func(p *Producer) {
		p.cfg.retryInterval = interval
	}
}
