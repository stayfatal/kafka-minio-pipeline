package router

import (
	"log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/stayfatal/kafka-minio-pipeline/rest/router/handlers"
)

func New() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	hm, err := handlers.NewManager()
	if err != nil {
		log.Fatal(err)
	}

	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)

	v1ApiRouter := chi.NewRouter()
	apiRouter.Mount("/v1", v1ApiRouter)

	v1ApiRouter.Post("/upload", hm.Upload)

	return r
}
