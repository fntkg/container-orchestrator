// File: pkg/taskmanager/taskmanager_test.go
package taskmanager_test

import (
	"testing"

	"github.com/fntkg/container-orchestrator/pkg/datastore"
	"github.com/fntkg/container-orchestrator/pkg/models"
	"github.com/fntkg/container-orchestrator/pkg/taskmanager"
)

func TestTaskManager_CreateAndGetTask(t *testing.T) {
	ds := datastore.NewInMemoryDatastore()
	tm := taskmanager.NewTaskManager(ds)

	// Create a new task with a status.
	task := models.Task{ID: "task-1", Status: "pending"}
	err := tm.CreateTask(task)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	// Retrieve the task by its ID.
	retrievedTask, err := tm.GetTask("task-1")
	if err != nil {
		t.Fatalf("failed to retrieve task: %v", err)
	}

	if retrievedTask.ID != task.ID {
		t.Errorf("expected task ID %s, got %s", task.ID, retrievedTask.ID)
	}
	if retrievedTask.Status != "pending" {
		t.Errorf("expected task status 'pending', got '%s'", retrievedTask.Status)
	}
}

func TestTaskManager_UpdateTask(t *testing.T) {
	ds := datastore.NewInMemoryDatastore()
	tm := taskmanager.NewTaskManager(ds)

	// Create an initial task.
	task := models.Task{ID: "task-2", Status: "pending"}
	if err := tm.CreateTask(task); err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	// Update the task's status.
	updatedTask := models.Task{ID: "task-2", Status: "running"}
	if err := tm.UpdateTask(updatedTask); err != nil {
		t.Fatalf("failed to update task: %v", err)
	}

	// Retrieve the task and verify the updated status.
	retrievedTask, err := tm.GetTask("task-2")
	if err != nil {
		t.Fatalf("failed to retrieve task: %v", err)
	}

	if retrievedTask.Status != "running" {
		t.Errorf("expected task status 'running', got '%s'", retrievedTask.Status)
	}
}

func TestTaskManager_GetTasks(t *testing.T) {
	ds := datastore.NewInMemoryDatastore()
	tm := taskmanager.NewTaskManager(ds)

	tasksToCreate := []models.Task{
		{ID: "task-3", Status: "pending"},
		{ID: "task-4", Status: "pending"},
	}

	// Create multiple tasks.
	for _, task := range tasksToCreate {
		if err := tm.CreateTask(task); err != nil {
			t.Fatalf("failed to create task %s: %v", task.ID, err)
		}
	}

	// Retrieve all tasks.
	tasks, err := tm.GetTasks()
	if err != nil {
		t.Fatalf("failed to get tasks: %v", err)
	}

	if len(tasks) != len(tasksToCreate) {
		t.Errorf("expected %d tasks, got %d", len(tasksToCreate), len(tasks))
	}
}
