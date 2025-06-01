package minio

import "time"

type Opt func(*Client)

func SetRetryAmount(retryAmount int) Opt {
	return func(c *Client) {
		c.cfg.retryAmount = retryAmount
	}
}

func SetRetryInterval(retryInterval time.Duration) Opt {
	return func(c *Client) {
		c.cfg.retryInterval = retryInterval
	}
}
