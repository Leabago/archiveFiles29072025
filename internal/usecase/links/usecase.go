package links

import (
	"archive/zip"
	"archiveFiles/internal/entity"
	"archiveFiles/internal/repo"
	builderror "archiveFiles/pkg/buildError"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"go.uber.org/zap"
)

const temFolderName = "archive-links-*"
const archiveName = "archive-*.zip"
const downloadLink = "/api/v1/arch/%d"

type Links struct {
	httpClient repo.IhttpClient

	// масимальное число ссылок
	maxNumLinks int
	// максимальное число активных тасок
	maxTaskCount int

	tasks           map[uint64]entity.Task
	activeTaskCount int64
	taskMux         sync.Mutex
}

func New(httpClient repo.IhttpClient, maxNumLinks, maxTaskCount int) *Links {
	return &Links{
		httpClient: httpClient,
		// tempFileManager: tempFileManager,
		maxNumLinks:  maxNumLinks,
		maxTaskCount: maxTaskCount,
		tasks:        make(map[uint64]entity.Task),
	}
}

func (p *Links) ProcessLinks(zipWriter *zip.Writer, links []entity.Link, l *zap.Logger) ([]string, []entity.DownloadResult, error) {
	errMessages := make([]string, 0, p.maxNumLinks)
	downloadResults := make([]entity.DownloadResult, 0, p.maxNumLinks)

	// проверка на кол-во ссылок
	if (len(links)) > p.maxNumLinks {
		return errMessages, downloadResults, builderror.LinkLimitError{
			Requested:  len(links),
			MaxAllowed: p.maxNumLinks,
		}
	}

	file := make(chan entity.DownloadResult)
	wg := sync.WaitGroup{}
	wg.Add(len(links))

	go func() {
		for i, url := range links {
			// параллельное скачивание файлов

			go func() {
				defer wg.Done()
				l.Debug(fmt.Sprintf("download: file-%d,  url='%s'", i, url))
				// скачать файл по ссылке
				downloadResult := p.httpClient.DownloadByLink(url.URL, i, l)
				file <- *downloadResult
			}()
		}

		wg.Wait()
		close(file)
	}()

	// архивирование body
	for v := range file {
		downloadResults = append(downloadResults, v)
		if v.Error == nil {
			writer, err := zipWriter.Create(v.Filename)
			if err != nil {
				return errMessages, downloadResults, err
			}

			_, err = writer.Write(v.Content)
			if err != nil {
				return errMessages, downloadResults, err
			}
		} else {
			errMessages = append(errMessages, v.Error.Error())
		}
	}

	if len(errMessages) == len(links) {
		return errMessages, downloadResults, fmt.Errorf("Can't download at least one file")
	}

	return errMessages, downloadResults, nil
}

func (u *Links) CreateTask(l *zap.Logger) (*entity.Task, error) {

	taskCount := u.GetActiveTaskCount()
	if taskCount >= u.maxTaskCount {
		return nil, builderror.ActiveTasksLimitError{
			MaxAllowed: u.maxTaskCount,
		}
	}

	task := entity.NewTask()
	u.AddTask(*task)

	return task, nil
}

func (u *Links) AppendLink(taskID int, link string, l *zap.Logger) (*entity.Task, error) {
	taskIDint64 := uint64(taskID)
	linkEnt := entity.Link{
		URL: link,
	}
	exist, task := u.CheckTaskExists(taskIDint64)

	if !exist {
		return nil, builderror.NotFoundError{
			Entity: "Task",
			ID:     strconv.Itoa(taskID),
		}
	} else {
		l.Debug(task.Status)
		if task.Status == entity.StatusReady {
			return task, nil
		}
		_, task = u.AddLinkToTask(taskIDint64, linkEnt)
	}

	if len(task.Links) >= u.maxNumLinks {
		// сделать архив при достижении максимума ссылок
		// zip для архивирования файлов
		buf := new(bytes.Buffer)
		zipWriter := zip.NewWriter(buf)
		defer zipWriter.Close()

		// обработать ссылки
		_, downloadResult, err := u.ProcessLinks(zipWriter, task.Links, l)
		if err != nil {
			l.Debug(fmt.Sprintf("IsLinkLimitError: %v", builderror.IsLinkLimitError(err)))
			return nil, err
		}

		// set filename and fill link errors if exist and
		for _, v := range downloadResult {
			if v.Error != nil {
				task.Links[v.FileNum].Err = v.Error
			}
			task.Links[v.FileNum].FileName = v.Filename
		}

		// create zip file
		err = zipWriter.Close()
		if err != nil {
			l.Error(fmt.Sprintf("Failed to close ZIP writer: %v", err))
			return nil, err
		}

		arch, err := os.CreateTemp("", archiveName)
		if err != nil {
			l.Error(fmt.Sprintf("Failed to create ZIP file: %v", err))
			return nil, err
		}
		defer arch.Close()

		_, err = arch.Write(buf.Bytes())
		if err != nil {
			l.Error(fmt.Sprintf("Failed write to ZIP file: %v", err))
			return nil, err
		}

		task.Status = entity.StatusReady
		task.FilePath = arch.Name()
		task.DownloadLink = fmt.Sprintf(downloadLink, taskID)
		u.UpdateTask(task.ID, *task)
	}

	exist, task = u.CheckTaskExists(taskIDint64)
	return task, nil
}

func (u *Links) GetArchive(taskID int, l *zap.Logger) ([]byte, error) {
	taskIDint64 := uint64(taskID)

	exist, task := u.CheckTaskExists(taskIDint64)

	if !exist || task.Status != entity.StatusReady {
		return nil, builderror.NotFoundError{
			Entity: "Archive",
			ID:     strconv.Itoa(taskID),
		}
	}

	content, err := os.ReadFile(task.FilePath)
	if err != nil {
		l.Error(fmt.Sprintf("Read file error: %w", err), zap.String(
			"filePath", task.FilePath,
		))
		return nil, err
	}

	return content, nil
}

func (u *Links) AddTask(task entity.Task) {
	u.taskMux.Lock()
	defer u.taskMux.Unlock()
	u.tasks[task.ID] = task
	u.activeTaskCount++
}

func (u *Links) AddLinkToTask(taskID uint64, link entity.Link) (bool, *entity.Task) {
	u.taskMux.Lock()
	defer u.taskMux.Unlock()
	if task, exists := u.tasks[taskID]; exists {
		task.Links = append(task.Links, link)
		u.tasks[taskID] = task
		return true, &task
	}
	return false, nil
}

func (u *Links) CheckTaskExists(taskID uint64) (bool, *entity.Task) {
	u.taskMux.Lock()
	defer u.taskMux.Unlock()
	if task, exists := u.tasks[taskID]; exists {
		return true, &task
	}
	return false, nil
}

func (u *Links) UpdateTask(taskID uint64, task entity.Task) bool {
	u.taskMux.Lock()
	defer u.taskMux.Unlock()
	if _, exists := u.tasks[taskID]; exists {
		u.tasks[taskID] = task
		return true
	}
	return false
}

func (u *Links) GetActiveTaskCount() int {
	return int(atomic.LoadInt64(&u.activeTaskCount))
}
