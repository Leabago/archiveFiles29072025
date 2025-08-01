package tempfilemanager

import (
	"archiveFiles/internal/entity"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const dirName = "/tmp/archive-file-links"

type TempFileManager struct {
	dirPath string
}

func New() (*TempFileManager, error) {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &TempFileManager{
		dirPath: dirName,
	}, nil
}

// сделать файловую структуру:
// создать папку task-id,
// внутри task-id папки сделать файл id.json для отслеживания статуса задачиt
// внутри task-id папки сделать папку files для складвания загружаемых файлов
func (t *TempFileManager) CreateTask() (*entity.Task, error) {

	task, dir, err := t.CreateFileStruct()
	if err != nil {
		if dir != "" {
			DeleteFolderIfExists(dir)
		}
		return nil, err
	}

	return task, nil
}

func (t *TempFileManager) CreateFileStruct() (*entity.Task, string, error) {

	// папка task-id
	dir, err := os.MkdirTemp(t.dirPath, "task-*")
	if err != nil {
		return nil, dir, err
	}

	// взять id папки
	split := strings.Split(dir, "-")
	id := split[len(split)-1]

	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, dir, err
	}

	task := &entity.Task{
		Status:   entity.StatusCreaing,
		ID:       idInt,
		FilePath: dir,
	}

	// файл id.json для статуса задачи
	file, err := os.Create(filepath.Join(dir, id+".json"))
	if err != nil {
		return nil, dir, err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(task)
	if err != nil {
		return nil, dir, err
	}

	// папка files для храненя загружаемых файлов
	err = os.MkdirAll(filepath.Join(dir, "files"), os.ModePerm)
	if err != nil {
		return nil, dir, err
	}

	return task, dir, nil
}

func DeleteFolderIfExists(path string) error {

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed to remove folder: %w", err)
	}

	return nil
}
