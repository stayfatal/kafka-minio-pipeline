package minio

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	URL       string
	AccessKey string
	SecretKey string
	UseSSL    bool

	retryAmount   int
	retryInterval time.Duration
}

type Client struct {
	conn *minio.Client

	cfg Config
}

var (
	defaultRetryAmount   = 10
	defaultRetryInterval = time.Second * 3
	ErrConnectionFailed  = errors.New("connection to minio failed")
)

func New(cfg Config, opts ...Opt) (*Client, error) {
	client := &Client{cfg: cfg}
	client.cfg.retryAmount = defaultRetryAmount
	client.cfg.retryInterval = defaultRetryInterval

	for _, opt := range opts {
		opt(client)
	}

	err := client.connect()
	if err != nil {
		return nil, fmt.Errorf("client.connect: %w", err)
	}

	log.Println("подключение установлено")

	return client, nil
}

func (c *Client) connect() error {
	var err error
	c.conn, err = minio.New(c.cfg.URL, &minio.Options{
		Creds:  credentials.NewStaticV4(c.cfg.AccessKey, c.cfg.SecretKey, ""),
		Secure: c.cfg.UseSSL,
	})

	if err != nil {
		return fmt.Errorf("minio.New:%w", err)
	}

	cancel, err := c.conn.HealthCheck(c.cfg.retryInterval)
	defer cancel()
	if err != nil {
		return fmt.Errorf("c.conn.HealthCheck:%w", err)
	}

	var ok bool
	for i := c.cfg.retryAmount; i > 0; i-- {
		ok = c.conn.IsOnline()
		if ok {
			break
		}

		time.Sleep(c.cfg.retryInterval)
	}

	if !ok {
		return fmt.Errorf("c.conn.IsOnline: %w", ErrConnectionFailed)
	}

	return nil
}
