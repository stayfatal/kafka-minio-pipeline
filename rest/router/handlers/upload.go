package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/stayfatal/kafka-minio-pipeline/pkg/detect"
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

	contentType, err := detect.DetectContentType(f)
	if err != nil {
		log.Println(err)
		httpErr(w, err, http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = hm.minioClient.Upload(ctx, "images", myminio.Object{
		Name:        header.Filename,
		ContentType: contentType,
		Size:        header.Size,
		File:        f,
	}, minio.PutObjectOptions{})

	if err != nil {
		log.Println(err)
		httpErr(w, err, http.StatusInternalServerError)
		return
	}
}
