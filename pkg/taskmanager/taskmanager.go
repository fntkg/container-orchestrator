// File: pkg/taskmanager/taskmanager.go
package taskmanager

import (
	"errors"

	"github.com/fntkg/container-orchestrator/pkg/datastore"
	"github.com/fntkg/container-orchestrator/pkg/models"
)

type TaskManagement interface {
	CreateTask(task models.Task) error
	GetTask(taskID string) (*models.Task, error)
	GetTasks() ([]models.Task, error)
	UpdateTask(task models.Task) error
}

// TaskManager handles the lifecycle of tasks.
type TaskManager struct {
	ds datastore.Datastore
}

// NewTaskManager returns a new instance of TaskManager.
func NewTaskManager(ds datastore.Datastore) *TaskManager {
	return &TaskManager{
		ds: ds,
	}
}

// CreateTask stores a new task in the datastore.
func (tm *TaskManager) CreateTask(task models.Task) error {
	return tm.ds.SaveTask(task)
}

// GetTask retrieves a task by ID.
func (tm *TaskManager) GetTask(taskID string) (*models.Task, error) {
	tasks, err := tm.ds.GetTasks()
	if err != nil {
		return nil, err
	}
	for _, t := range tasks {
		if t.ID == taskID {
			return &t, nil
		}
	}
	return nil, errors.New("task not found")
}

// GetTasks retrieves all tasks.
func (tm *TaskManager) GetTasks() ([]models.Task, error) {
	return tm.ds.GetTasks()
}

// UpdateTask updates an existing task.
func (tm *TaskManager) UpdateTask(task models.Task) error {
	// In our simple datastore, SaveTask will overwrite any existing task with the same ID.
	return tm.ds.SaveTask(task)
}
