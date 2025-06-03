package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stayfatal/kafka-minio-pipeline/domain"
	myminio "github.com/stayfatal/kafka-minio-pipeline/pkg/minio"
)

func (hm *HandlersManager) Upload(w http.ResponseWriter, r *http.Request) {
	f, header, err := r.FormFile("car")
	if err != nil {
		log.Println(err)
		httpErr(w, err, http.StatusBadRequest)
		return
	}
	defer f.Close()

	contentType, err := detectContentType(f)
	if err != nil {
		log.Println(err)
		httpErr(w, err, http.StatusInternalServerError)
		return
	}

	image := domain.Image{Bucket: "images"}
	image.Hash, err = hashFile(f)
	if err != nil {
		log.Println(err)
		httpErr(w, err, http.StatusInternalServerError)
		return
	}

	obj := myminio.Object{
		Name:        image.Hash,
		ContentType: contentType,
		Size:        header.Size,
		File:        f,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = hm.minioClient.Upload(ctx, "images", obj, minio.PutObjectOptions{})
	if err != nil {
		log.Println(err)
		httpErr(w, err, http.StatusInternalServerError)
		return
	}

	err = hm.kafkaProducer.Send("images", image.Hash, image)
	if err != nil {
		log.Println(err)
		httpErr(w, err, http.StatusInternalServerError)
		return
	}
}

func detectContentType(file io.ReadSeeker) (string, error) {
	buf := make([]byte, 512) // Для детекции нужно первые 512 байт
	_, err := file.Read(buf)
	if err != nil {
		return "", err
	}
	_, err = file.Seek(0, io.SeekStart) // Возвращаем указатель в начало
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buf), nil
}

func hashFile(file io.ReadSeeker) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}

	_, err := file.Seek(0, io.SeekStart) // Возвращаем указатель в начало
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
