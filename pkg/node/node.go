package node

import (
	"errors"
	"sync"
)

// Node represents a cluster node with an ID and a health status.
type Node struct {
	ID      string
	Healthy bool
}

// NodeManager is an interface that defines the behavior of a node manager.
type NodeManager interface {
	GetNodes() []Node
}

// Manager manages the nodes in the cluster.
type Manager struct {
	nodes map[string]Node
	mu    sync.RWMutex
}

// NewManager initializes and returns a new Node Manager.
func NewManager() *Manager {
	return &Manager{
		nodes: make(map[string]Node),
	}
}

// Register adds a new node to the manager.
func (m *Manager) Register(n Node) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.nodes[n.ID]; exists {
		return errors.New("node already registered")
	}
	m.nodes[n.ID] = n
	return nil
}

// GetNodes returns a slice of all registered nodes.
func (m *Manager) GetNodes() []Node {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []Node
	for _, node := range m.nodes {
		result = append(result, node)
	}
	return result
}

// UpdateHealth updates the health status of a node.
func (m *Manager) UpdateHealth(nodeID string, healthy bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	n, exists := m.nodes[nodeID]
	if !exists {
		return errors.New("node not found")
	}
	n.Healthy = healthy
	m.nodes[nodeID] = n
	return nil
}
