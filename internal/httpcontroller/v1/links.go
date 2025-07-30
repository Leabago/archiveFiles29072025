package v1

import (
	"archive/zip"
	"archiveFiles/internal/httpcontroller/middleware"
	"archiveFiles/internal/httpcontroller/v1/api"
	builderror "archiveFiles/pkg/buildError"
	"archiveFiles/pkg/logger"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

func (v *V1) SetOfLinks(w http.ResponseWriter, r *http.Request) {

	// логгер с id запроса "X-Request-ID"
	l := logger.FromContext(r.Context())
	defer r.Body.Close()

	//
	links := &api.SetOfLinksJSONRequestBody{}
	err := json.NewDecoder(r.Body).Decode(links)

	if err != nil {
		l.Error("Failed decode request body")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// количество ссылок должно быть больше 0
	if links.Links == nil || len(*links.Links) == 0 {
		http.Error(w, "Links count must be > 0", http.StatusBadRequest)
		return
	}

	// zip для архивирования файлов
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	// обработать ссылки
	errorMess, err := v.usecase.ProcessLinks(zipWriter, links.Links, l)
	if err != nil {
		if builderror.IsLinkLimitError(err) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// для скачивания файла archive.zip
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
	w.Header().Set("Cache-Control", "no-store")

	// файл/файлы по ссылке не скачался или был отфильтрован
	if len(errorMess) > 0 {
		// хедер для проверки общего статуса скачанных по url файлов
		w.Header().Set(middleware.XError, "true")
		w.Header().Set(middleware.XErrorMessage, strings.Join(errorMess, " | "))
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// закрыть zip
	if err := zipWriter.Close(); err != nil {
		http.Error(w, "Failed to create ZIP archive", http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, buf); err != nil {
		l.Error("Failed to send ZIP file", zap.Error(err))
	}
}

func (v *V1) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	l := logger.FromContext(r.Context())

	task, err := v.usecase.CreateTask(l)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	now := time.Now()
	apiTask := &api.Task{
		Id:      &task.ID,
		Created: &now,
	}

	json.NewEncoder(w).Encode(apiTask)

	w.WriteHeader(http.StatusCreated)

}
