package scheduler

import (
	"github.com/fntkg/container-orchestrator/pkg/node"
	"testing"
)

func TestDefaultScheduler_Schedule(t *testing.T) {
	// Creates an instance of the scheduler
	scheduler := NewDefaultScheduler()

	// Define a sample task
	task := Task{ID: "task-1"}

	// Defines a list of available nodes
	nodes := []node.Node{
		{ID: "node-1", Healthy: true},
		{ID: "node-2", Healthy: true},
	}

	// Try to assign the task
	selectedNode, err := scheduler.Schedule(task, nodes)
	if err != nil {
		t.Fatalf("Error when scheduling the task: %v", err)
	}

	// Verify that the first node has been selected.
	if selectedNode.ID != "node-1" {
		t.Errorf("The task was expected to be assigned to node-1, but was assigned to %s", selectedNode.ID)
	}
}

func TestDefaultScheduler_NoNodes(t *testing.T) {
	scheduler := NewDefaultScheduler()
	task := Task{ID: "task-2"}
	var nodes []node.Node // No nodes available

	_, err := scheduler.Schedule(task, nodes)
	if err == nil {
		t.Fatal("An error was expected as no nodes were available, but none were obtained.")
	}
}
