package scheduler

import (
	"fmt"
	"github.com/fntkg/container-orchestrator/pkg/models"
)

// Scheduler defines the interface that must implement any scheduling strategy.
type Scheduler interface {
	// Schedule assigns a task to one of the available nodes and returns the selected node.
	Schedule(task models.Task, nodes []models.Node) (*models.Node, error)
}

// DefaultScheduler is a simple implementation of the Scheduler.
type DefaultScheduler struct {
	// You can add additional fields for more advanced strategies.
}

// NewDefaultScheduler returns an instance of DefaultScheduler.
func NewDefaultScheduler() *DefaultScheduler {
	return &DefaultScheduler{}
}

// Schedule assigns the task to the first available node.
// This is a very basic strategy that can be improved in future iterations.
func (s *DefaultScheduler) Schedule(task models.Task, nodes []models.Node) (*models.Node, error) {
	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes available to schedule the task %s", task.ID)
	}
	// Assign the task to the first available node.
	return &nodes[0], nil
}
