package links

import (
	"archive/zip"
	"archiveFiles/internal/entity"
	"archiveFiles/internal/repo"
	builderror "archiveFiles/pkg/buildError"
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
)

const temFolderName = "archive-links-*"

type Links struct {
	httpClient  repo.IhttpClient
	maxNumLinks int
}

func New(httpClient repo.IhttpClient, maxNumLinks int) *Links {
	return &Links{
		httpClient:  httpClient,
		maxNumLinks: maxNumLinks,
	}
}

func (p *Links) ProcessLinks(zipWriter *zip.Writer, links *[]string, l *zap.Logger) ([]string, error) {
	errMessages := make([]string, 0, p.maxNumLinks)

	// проверка на кол-во ссылок
	if (len(*links)) > p.maxNumLinks {
		return errMessages, builderror.LinkLimitError{
			Requested:  len(*links),
			MaxAllowed: p.maxNumLinks,
		}
	}

	file := make(chan entity.DownloadResult)
	wg := sync.WaitGroup{}
	wg.Add(len(*links))

	go func() {
		for i, url := range *links {
			// параллельное скачивание файлов

			go func() {
				defer wg.Done()
				l.Debug(fmt.Sprintf("download: file-%d,  url='%s'", i, url))
				// скачать файл по ссылке
				downloadResult := p.httpClient.DownloadByLink(url, i, l)
				file <- *downloadResult
			}()
		}

		wg.Wait()
		close(file)
	}()

	// архивирование body
	for v := range file {

		if v.Error == nil {
			writer, err := zipWriter.Create(v.Filename)
			if err != nil {
				return errMessages, err
			}

			_, err = writer.Write(v.Content)
			if err != nil {
				return errMessages, err
			}
		} else {
			errMessages = append(errMessages, v.Error.Error())
		}
	}

	return errMessages, nil
}

func (u *Links) CreateTask(l *zap.Logger) (*entity.Task, error) {
	// создать папку в /tmp
	dir, err := os.MkdirTemp("", temFolderName)
	if err != nil {
		err := fmt.Errorf("Failed create dir %s, int /tmp folder. Error: %w", temFolderName, err)
		l.Error(err.Error())
		return nil, err
	}

	// взять id папки
	split := strings.Split(dir, "-")
	id := split[len(split)-1]

	task := &entity.Task{
		ID:  id,
		Dir: dir,
	}

	// добавить task в мапу
	entity.CreateTask(*task)

	return task, nil
}
