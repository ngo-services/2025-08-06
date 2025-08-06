package task

import (
	"errors"
	"sync"
)

type Manager struct {
	tasks     map[string]*Task
	mu        sync.Mutex
	maxActive int
}

func NewManager(maxActive int) *Manager {
	return &Manager{
		tasks:     make(map[string]*Task),
		maxActive: maxActive,
	}
}

func (m *Manager) activeCountUnlocked() int {
	count := 0
	for _, t := range m.tasks {
		if t.Status == StatusPending || t.Status == StatusPacking {
			count++
		}
	}
	return count
}

func (m *Manager) NewTask() (*Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.activeCountUnlocked() >= m.maxActive {
		return nil, errors.New("maximum active tasks reached")
	}
	task := NewTask()
	m.tasks[task.ID] = task
	return task, nil
}

func (m *Manager) GetTask(id string) (*Task, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tasks[id]
	return t, ok
}

func (m *Manager) AddFileToTask(id string, file FileLink, maxFiles int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tasks[id]
	if !ok {
		return errors.New("task not found")
	}
	if len(t.Files) >= maxFiles {
		return errors.New("max files reached for this task")
	}
	t.Files = append(t.Files, file)
	return nil
}

func (m *Manager) SetTaskStatus(id string, status TaskStatus, archiveURL string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tasks[id]
	if ok {
		t.Status = status
		if status == StatusReady {
			t.ArchiveURL = archiveURL
		}
	}
}
