// File: pkg/scheduler/scheduler_test.go
package scheduler_test

import (
	"testing"

	"github.com/fntkg/container-orchestrator/pkg/models"
	"github.com/fntkg/container-orchestrator/pkg/scheduler"
)

func TestDefaultScheduler_ScheduleSuccess(t *testing.T) {
	// Create a new DefaultScheduler instance.
	sched := scheduler.NewDefaultScheduler()

	// Define a sample task.
	task := models.Task{ID: "task-1"}

	// Define a slice of available nodes.
	nodes := []models.Node{
		{ID: "node-1", Healthy: true},
		{ID: "node-2", Healthy: true},
	}

	// Call Schedule.
	assignedNode, err := sched.Schedule(task, nodes)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify that the first node is returned.
	if assignedNode.ID != "node-1" {
		t.Errorf("expected assigned node ID 'node-1', got '%s'", assignedNode.ID)
	}
}

func TestDefaultScheduler_NoNodes(t *testing.T) {
	// Create a new DefaultScheduler instance.
	sched := scheduler.NewDefaultScheduler()

	// Define a sample task.
	task := models.Task{ID: "task-1"}

	// Define an empty slice of nodes.
	nodes := []models.Node{}

	// Call Schedule expecting an error.
	assignedNode, err := sched.Schedule(task, nodes)
	if err == nil {
		t.Fatalf("expected an error when no nodes are available, got nil")
	}

	// Ensure that no node is returned.
	if assignedNode != nil {
		t.Errorf("expected nil assigned node when no nodes are available, got %+v", assignedNode)
	}
}
