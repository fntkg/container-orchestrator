package controller

import (
	"testing"

	"github.com/fntkg/container-orchestrator/pkg/node"
	"github.com/fntkg/container-orchestrator/pkg/scheduler"
)

// FakeScheduler is a mock implementation of scheduler.Scheduler for testing.
type FakeScheduler struct {
	scheduledTasks []scheduler.Task
	nodeToReturn   node.Node
	errToReturn    error
}

// Schedule records the task and returns a predefined node or error.
func (fs *FakeScheduler) Schedule(task scheduler.Task, nodes []node.Node) (*node.Node, error) {
	fs.scheduledTasks = append(fs.scheduledTasks, task)
	if fs.errToReturn != nil {
		return nil, fs.errToReturn
	}
	return &fs.nodeToReturn, nil
}

// FakeNodeManager is a fake implementation of a node manager for testing.
type FakeNodeManager struct {
	nodes []node.Node
}

// GetNodes returns the fake list of nodes.
func (fnm *FakeNodeManager) GetNodes() []node.Node {
	return fnm.nodes
}

// Register register a fake node
func Register(n node.Node) error {
	return nil
}

// UpdateHealth updates the health of a Node
func UpdateHealth(nodeID string, healthy bool) error {
	return nil
}

func TestControllerManager_Reconcile(t *testing.T) {
	// Create sample tasks.
	tasks := []scheduler.Task{
		{ID: "task-1"},
		{ID: "task-2"},
	}

	// Create sample nodes.
	nodesList := []node.Node{
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

	// Initialize the ControllerManager with the fake scheduler, tasks, and fake Node Manager.
	cm := NewControllerManager(fakeScheduler, tasks, fakeNodeManager)

	// Call the reconcile method directly.
	cm.reconcile()

	// Verify that the scheduler was called for each task.
	if len(fakeScheduler.scheduledTasks) != len(tasks) {
		t.Errorf("Expected %d scheduled tasks, got %d", len(tasks), len(fakeScheduler.scheduledTasks))
	}

	// Check that the recorded tasks match the expected tasks.
	for i, task := range tasks {
		if fakeScheduler.scheduledTasks[i].ID != task.ID {
			t.Errorf("Expected task ID %s, got %s", task.ID, fakeScheduler.scheduledTasks[i].ID)
		}
	}
}
