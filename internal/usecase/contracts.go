package usecase

import (
	"archive/zip"
	"archiveFiles/internal/entity"

	"go.uber.org/zap"
)

type ILinks interface {
	ProcessLinks(zipWriter *zip.Writer, links *[]string, l *zap.Logger) ([]string, error)
	CreateTask(l *zap.Logger) (*entity.Task, error)
}
