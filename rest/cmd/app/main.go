package main

import (
	"net/http"

	"github.com/stayfatal/kafka-minio-pipeline/rest/router"
)

func main() {
	r := router.New()

	http.ListenAndServe(":8080", r)
}
