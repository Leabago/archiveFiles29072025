package v1

import (
	"archive/zip"
	"archiveFiles/internal/httpcontroller/middleware"
	"archiveFiles/internal/httpcontroller/v1/api"
	builderror "archiveFiles/pkg/buildError"
	"archiveFiles/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

	taskLinks := convertLinks(links.Links)

	// обработать ссылки
	errorMess, _, err := v.Usecase.ProcessLinks(zipWriter, taskLinks, l)
	if err != nil {
		l.Debug(fmt.Sprintf("IsLinkLimitError: %v", builderror.IsLinkLimitError(err)))
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

	task, err := v.Usecase.CreateTask(l)
	if err != nil {
		l.Debug(fmt.Sprintf("IsActiveTasksError: %v", builderror.IsActiveTasksLimitError(err)))
		code := http.StatusInternalServerError

		if builderror.IsActiveTasksLimitError(err) {
			code = http.StatusBadRequest
		}

		mess := err.Error()
		respErr := &api.ErrorResponse{
			Code:  &code,
			Error: &mess,
		}
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(respErr)
		return
	}

	apiTask := convertTask(task)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiTask)
}

func (v *V1) AppendLink(w http.ResponseWriter, r *http.Request, taskID int) {
	w.Header().Set("Content-Type", "application/json")
	l := logger.FromContext(r.Context())

	link := &api.AppendLinkJSONRequestBody{}
	err := json.NewDecoder(r.Body).Decode(link)
	if err != nil {
		l.Error("Failed decode request body")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if link.Link == nil {
		http.Error(w, "Link must be present", http.StatusBadRequest)
		return
	}

	//
	task, err := v.Usecase.AppendLink(taskID, *link.Link, l)

	if err != nil {
		code := http.StatusInternalServerError

		if builderror.IsNotFoundError(err) {
			code = http.StatusNotFound
		}

		mess := err.Error()
		respErr := &api.ErrorResponse{
			Code:  &code,
			Error: &mess,
		}
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(respErr)
		return
	}

	apiTask := convertTask(task)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(apiTask)
}

func (v *V1) GetArchive(w http.ResponseWriter, r *http.Request, taskID int) {
	l := logger.FromContext(r.Context())
	content, err := v.Usecase.GetArchive(taskID, l)
	if err != nil {
		code := http.StatusInternalServerError
		if builderror.IsNotFoundError(err) {
			code = http.StatusNotFound
		}

		mess := err.Error()
		respErr := &api.ErrorResponse{
			Code:  &code,
			Error: &mess,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(respErr)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=archive.zip")
	w.Header().Set("Cache-Control", "no-store")

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
