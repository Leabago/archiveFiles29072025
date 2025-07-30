package repo

import (
	"archiveFiles/internal/entity"

	"go.uber.org/zap"
)

type IhttpClient interface {
	DownloadByLink(url string, num int, l *zap.Logger) *entity.DownloadResult
}
