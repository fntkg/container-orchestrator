package controller

import (
	"fmt"
	"testing"

	"github.com/fntkg/container-orchestrator/pkg/models"
)

// FakeScheduler implements the scheduler.Scheduler interface for testing.
type FakeScheduler struct {
	scheduledTasks []models.Task
	nodeToReturn   models.Node
	errToReturn    error
}

// Schedule records the task and returns a predefined node (or error).
func (fs *FakeScheduler) Schedule(task models.Task, nodes []models.Node) (*models.Node, error) {
	fs.scheduledTasks = append(fs.scheduledTasks, task)
	if fs.errToReturn != nil {
		return nil, fs.errToReturn
	}
	return &fs.nodeToReturn, nil
}

// FakeNodeManager implements the node.NodeManager interface for testing.
type FakeNodeManager struct {
	nodes []models.Node
}

// Register adds a node to the fake manager.
func (fnm *FakeNodeManager) Register(n models.Node) error {
	fnm.nodes = append(fnm.nodes, n)
	return nil
}

// GetNodes returns the list of nodes.
func (fnm *FakeNodeManager) GetNodes() []models.Node {
	return fnm.nodes
}

// UpdateHealth updates the health status of a node.
func (fnm *FakeNodeManager) UpdateHealth(id string, healthy bool) error {
	for i, n := range fnm.nodes {
		if n.ID == id {
			fnm.nodes[i].Healthy = healthy
			return nil
		}
	}
	return fmt.Errorf("node not found")
}

// TestControllerManager_Reconcile verifies that the reconcile method schedules each task.
func TestControllerManager_Reconcile(t *testing.T) {
	// Create sample tasks.
	tasks := []models.Task{
		{ID: "task-1"},
		{ID: "task-2"},
	}

	// Create sample nodes.
	nodesList := []models.Node{
		{ID: "node-1", Healthy: true},
		{ID: "node-2", Healthy: true},
	}

	// Create a FakeNodeManager with the sample nodes.
	fakeNodeManager := &FakeNodeManager{
		nodes: nodesList,
	}

	// Create a FakeScheduler that always returns the first node.
	fakeScheduler := &FakeScheduler{
		nodeToReturn: nodesList[0],
	}

	// Create the Controller Manager using the fake dependencies.
	cm := NewControllerManager(fakeScheduler, tasks, fakeNodeManager)

	// Invoke the reconcile logic.
	cm.reconcile()

	// Verify that the scheduler was called for each task.
	if len(fakeScheduler.scheduledTasks) != len(tasks) {
		t.Errorf("Expected %d scheduled tasks, got %d", len(tasks), len(fakeScheduler.scheduledTasks))
	}

	// Check that the recorded tasks match the expected ones.
	for i, task := range tasks {
		if fakeScheduler.scheduledTasks[i].ID != task.ID {
			t.Errorf("Expected task ID %s, got %s", task.ID, fakeScheduler.scheduledTasks[i].ID)
		}
	}
}
