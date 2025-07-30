package v1

import (
	"archiveFiles/internal/usecase"
	"net/http"
)

func NewRouter(mux *http.ServeMux, usecase usecase.ILinks) {
	v1 := &V1{
		usecase: usecase,
	}

	mux.HandleFunc("POST /api/v1/upload", v1.SetOfLinks)

	mux.HandleFunc("POST /api/v1/task", v1.CreateTask)

}
