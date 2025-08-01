package entity

import (
	"sync"
	"sync/atomic"
)

const StatusCreaing = "creating"
const StatusReady = "ready"

var counterID uint64

// сделать ID
func NextID() uint64 {
	return atomic.AddUint64(&counterID, 1)
}

func NewTask() *Task {
	return &Task{
		ID:     NextID(),
		Status: StatusCreaing,
	}
}

type Task struct {
	Status   string
	ID       uint64
	Links    []Link
	FilePath string

	DownloadLink string
	Err          error
}

type Link struct {
	URL      string
	FileName string
	Err      error
}

var (
	tasks           = make(map[uint64]Task)
	activeTaskCount int64
	taskMux         sync.Mutex
)
