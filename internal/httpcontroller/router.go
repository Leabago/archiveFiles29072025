package httpcontroller

import (
	"archiveFiles/config"
	"archiveFiles/doc"

	"archiveFiles/internal/httpcontroller/middleware"
	v1 "archiveFiles/internal/httpcontroller/v1"
	"archiveFiles/internal/httpcontroller/v1/api"
	"archiveFiles/internal/usecase"
	"archiveFiles/pkg/logger"
	"net/http"
)

func NewRouter(cfg *config.Config, httpServer *http.Server, u usecase.ILinks) {

	baseRouter := http.NewServeMux()

	// swagger /doc endpoint
	if cfg.Swagger.Enable {
		fsOAPI := http.FileServer(http.FS(v1.OpenApi))
		baseRouter.Handle("/static/", http.StripPrefix("/static", fsOAPI))

		fsSwag := http.FileServer(http.FS(doc.Swagger))
		baseRouter.Handle("/doc/", http.StripPrefix("/doc", fsSwag))
	}

	baseRouter.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Routers
	v1 := &v1.V1{
		Usecase: u,
	}
	handler := api.HandlerFromMux(v1, baseRouter)

	// Wrap handler with logger middleware
	handler = middleware.Logger(logger.Log)(handler)

	httpServer.Handler = handler
}
