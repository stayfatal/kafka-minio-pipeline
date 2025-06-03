package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/minio/minio-go/v7"
	"github.com/stayfatal/kafka-minio-pipeline/domain"
	"github.com/stayfatal/kafka-minio-pipeline/pkg/kafka/consumer"
	myminio "github.com/stayfatal/kafka-minio-pipeline/pkg/minio"
)

func main() {
	minioClient, err := myminio.New(myminio.Config{
		URL:       "minio:9000",
		AccessKey: "admin",
		SecretKey: "password",
		UseSSL:    false,
	})
	if err != nil {
		log.Fatal(err)
	}

	consumerGroup, err := consumer.New(&consumer.Config{Brokers: []string{"kafka-1:19092", "kafka-2:19092", "kafka-3:19092"}})
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	consumerGroup.Consume(
		ctx,
		[]string{"images"},
		func(msg *sarama.ConsumerMessage) error {
			log.Println(string(msg.Value))

			image := domain.Image{}
			err := json.Unmarshal(msg.Value, &image)
			if err != nil {
				log.Println(err)
				return err
			}
			log.Println(image)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			obj, err := minioClient.Download(ctx, image.Bucket, image.Hash, minio.GetObjectOptions{})
			if err != nil {
				log.Println(err)
				return err
			}

			hasher := sha256.New()
			_, err = io.Copy(hasher, obj.File)
			if err != nil {
				log.Println(err)
				return err
			}
			hash := hex.EncodeToString(hasher.Sum(nil))
			if hash != image.Hash {
				err = errors.New("wrong hash")
				log.Println(err)
				return err
			}

			log.Println("Success!!!")

			return nil
		},
	)

}
