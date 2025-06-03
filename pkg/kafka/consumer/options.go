package consumer

import (
	"time"
)

type Opt func(*Consumer)

func SetRetryAmount(retryAmount int) Opt {
	return func(cg *Consumer) {
		cg.cfg.retryAmount = retryAmount
	}
}

func SetRetryInterval(backoff time.Duration) Opt {
	return func(cg *Consumer) {
		cg.cfg.retryInterval = backoff
	}
}
