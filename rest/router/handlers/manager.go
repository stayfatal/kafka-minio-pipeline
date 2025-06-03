package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stayfatal/kafka-minio-pipeline/pkg/kafka/producer"
	myminio "github.com/stayfatal/kafka-minio-pipeline/pkg/minio"
)

type HandlersManager struct {
	minioClient   *myminio.Client
	kafkaProducer *producer.Producer
}

func NewManager() (*HandlersManager, error) {
	minioClient, err := myminio.New(myminio.Config{
		URL:       "minio:9000",
		AccessKey: "admin",
		SecretKey: "password",
		UseSSL:    false,
	})
	if err != nil {
		return nil, fmt.Errorf("minio.New:%w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err = minioClient.SetupBucket(ctx, "images", minio.MakeBucketOptions{})
	if err != nil {
		return nil, fmt.Errorf("minioClient.SetupBucket: %w", err)
	}

	producer, err := producer.New(&producer.Config{
		Brokers: []string{"kafka-1:19092", "kafka-2:19092", "kafka-3:19092"},
	})
	if err != nil {
		return nil, fmt.Errorf("producer.New: %w", err)
	}

	return &HandlersManager{minioClient: minioClient, kafkaProducer: producer}, nil
}

func httpErr(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	resp, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	w.Write(resp)
}
