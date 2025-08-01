package usecase

import (
	"archive/zip"
	"archiveFiles/internal/entity"

	"go.uber.org/zap"
)

type ILinks interface {
	ProcessLinks(zipWriter *zip.Writer, links []entity.Link, l *zap.Logger) ([]string, []entity.DownloadResult, error)
	CreateTask(l *zap.Logger) (*entity.Task, error)
	AppendLink(taskID int, link string, l *zap.Logger) (*entity.Task, error)
	GetArchive(taskID int, l *zap.Logger) ([]byte, error)
}
