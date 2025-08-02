package httpclient

import (
	"archiveFiles/internal/entity"
	"archiveFiles/internal/repo"
	builderror "archiveFiles/pkg/buildError"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
	"go.uber.org/zap"
)

var _ repo.IhttpClient = (*Httpclient)(nil)

var Extensions = map[string]string{
	"application/json": ".json",
	"application/xml":  ".xml",
	"text/xml":         ".xml",
	"text/csv":         ".csv",
	"application/zip":  ".zip",
	"application/pdf":  ".pdf",
	// images
	"image/jpeg":               ".jpg",
	"image/jpg":                ".jpg",
	"image/png":                ".png",
	"image/gif":                ".gif",
	"image/webp":               ".webp",
	"image/svg+xml":            ".svg",
	"image/tiff":               ".tiff",
	"image/bmp":                ".bmp",
	"image/x-ms-bmp":           ".bmp",
	"image/vnd.microsoft.icon": ".ico",
}

type Httpclient struct {
	client             *http.Client
	allowedContentType []string
	maxBytesResp       int64
}

func New(contentType string, maxBytesResp int64) *Httpclient {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 5 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
		},
	}

	ht := &Httpclient{
		client:       client,
		maxBytesResp: maxBytesResp,
	}

	if contentType != "" {
		ht.allowedContentType = strings.Split(contentType, ",")
	}

	return ht
}

// DownloadByLink сделать запорс и провалидировать ответ
func (h *Httpclient) DownloadByLink(url string, num int, l *zap.Logger) *entity.DownloadResult {

	result := &entity.DownloadResult{
		Filename: "file" + strconv.Itoa(num),
		FileNum:  num,
		Content:  []byte{},
	}

	// запорс на скачивание
	resp, err := h.client.Get(url)
	if err != nil {
		result.Error = builderror.ErrUrl(fmt.Sprintf("Failed download file %d. %s", num, err), url)
		l.Error(result.Error.Error())
		return result
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Error = builderror.ErrUrl(fmt.Sprintf("Status not OK, status == %d", resp.StatusCode), url)
		l.Warn(result.Error.Error())
		return result
	}

	// отфильтровать по "Content-Type"
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		result.Error = builderror.ErrUrl("'Content-Type' header not defined", url)
		l.Warn(result.Error.Error())
		return result
	}

	if !lo.Contains(h.allowedContentType, contentType) {
		result.Error = builderror.ErrUrl(fmt.Sprintf("'Content-Type' header did not pass the filter. 'Content-Type' == %s, allowed Content-Type == %s", contentType, h.allowedContentType), url)
		l.Warn(result.Error.Error())
		return result
	}

	// отфильтровать по максимальному размеру ответа
	limitedReader := io.LimitReader(resp.Body, h.maxBytesResp)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		result.Error = builderror.ErrUrl(fmt.Sprintf("Error reading response body. %s", err), url)
		l.Warn(result.Error.Error())
		return result
	}

	if int64(len(body)) >= h.maxBytesResp {
		result.Error = builderror.ErrUrl(fmt.Sprintf("Response body exceeded limit. Read body %d bytes, max %d bytes", int64(len(body)), h.maxBytesResp), url)
		l.Warn(result.Error.Error())
		return result
	}

	// добавить расширение если есть
	result.Filename += Extensions[contentType]

	// ок
	result.Content = body
	return result
}
