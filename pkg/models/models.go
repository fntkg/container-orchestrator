// pkg/models/models.go
package models

// Node represents a cluster node.
type Node struct {
	ID      string
	Healthy bool
}

// Task represents a task that needs scheduling.
type Task struct {
	ID string
}
