package entity

import "sync"

type Task struct {
	ID    string
	Links []Link
	Dir   string
}

type Link struct {
	URL string
}

var (
	tasks   = make(map[string]Task)
	taskMux sync.RWMutex
)

func CreateTask(task Task) {
	taskMux.Lock()
	defer taskMux.Unlock()
	tasks[task.ID] = task
}

func AddLink(taskID string, link Link) bool {
	taskMux.Lock()
	defer taskMux.Unlock()
	if task, exists := tasks[taskID]; exists {
		task.Links = append(task.Links, link)
		tasks[taskID] = task
		return true
	}
	return false
}
