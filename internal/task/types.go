package task

import (
	"sync"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusPending TaskStatus = "pending"
	StatusPacking TaskStatus = "packing"
	StatusReady   TaskStatus = "ready"
	StatusFailed  TaskStatus = "failed"
)

type Task struct {
	ID         string
	Status     TaskStatus
	Files      []FileLink
	ArchiveURL string
	Errors     []string

	Mu sync.Mutex
}

type FileLink struct {
	URL      string
	Filename string
	Status   string
	Error    string
}

func NewTask() *Task {
	return &Task{
		ID:     uuid.New().String(),
		Status: StatusPending,
		Files:  make([]FileLink, 0),
		Errors: make([]string, 0),
	}
}
