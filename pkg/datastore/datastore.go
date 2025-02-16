package datastore

import (
	"sync"

	"github.com/fntkg/container-orchestrator/pkg/models"
)

// Datastore defines the methods to store and retrieve the cluster state.
type Datastore interface {
	SaveNode(n models.Node) error
	GetNodes() ([]models.Node, error)
	SaveTask(t models.Task) error
	GetTasks() ([]models.Task, error)
}

// InMemoryDatastore is a simple in-memory implementation of Datastore.
type InMemoryDatastore struct {
	nodes map[string]models.Node
	tasks map[string]models.Task
	mu    sync.RWMutex
}

// NewInMemoryDatastore creates a new instance of InMemoryDatastore.
func NewInMemoryDatastore() *InMemoryDatastore {
	return &InMemoryDatastore{
		nodes: make(map[string]models.Node),
		tasks: make(map[string]models.Task),
	}
}

// SaveNode stores a node in the datastore.
func (ds *InMemoryDatastore) SaveNode(n models.Node) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.nodes[n.ID] = n
	return nil
}

// GetNodes retrieves all nodes from the datastore.
func (ds *InMemoryDatastore) GetNodes() ([]models.Node, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	nodes := make([]models.Node, 0, len(ds.nodes))
	for _, n := range ds.nodes {
		nodes = append(nodes, n)
	}
	return nodes, nil
}

// SaveTask stores a task in the datastore.
func (ds *InMemoryDatastore) SaveTask(t models.Task) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.tasks[t.ID] = t
	return nil
}

// GetTasks retrieves all tasks from the datastore.
func (ds *InMemoryDatastore) GetTasks() ([]models.Task, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	tasks := make([]models.Task, 0, len(ds.tasks))
	for _, t := range ds.tasks {
		tasks = append(tasks, t)
	}
	return tasks, nil
}
