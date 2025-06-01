package main

import (
	"log"
	"net/http"

	"github.com/stayfatal/kafka-minio-pipeline/rest/router"
)

func main() {
	r := router.New()

	log.Println("server is listening on 8080")
	http.ListenAndServe(":8080", r)
}
