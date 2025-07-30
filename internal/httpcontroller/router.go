package httpcontroller

import (
	"archiveFiles/config"
	"archiveFiles/doc"

	"archiveFiles/internal/httpcontroller/middleware"
	v1 "archiveFiles/internal/httpcontroller/v1"
	"archiveFiles/internal/usecase"
	"archiveFiles/pkg/logger"
	"net/http"
)

func NewRouter(cfg *config.Config, httpServer *http.Server, u usecase.ILinks) {

	mux := http.NewServeMux()

	// swagger /doc endpoint
	if cfg.Swagger.Enable {
		fsOAPI := http.FileServer(http.FS(v1.OpenApi))
		mux.Handle("/static/", http.StripPrefix("/static", fsOAPI))

		fsSwag := http.FileServer(http.FS(doc.Swagger))
		mux.Handle("/doc/", http.StripPrefix("/doc", fsSwag))
	}

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Routers
	v1.NewRouter(mux, u)

	// Wrap with logger middleware
	httpServer.Handler = middleware.Logger(logger.Log)(mux)
}
