# container-orchestrator
Application that simulates the basic operation of Kubernetes.

## 🗂️ How it works

>  TODO: Revisar si esto es asi de verdad

```mermaid
graph LR
    A[User / CLI] -->|Send requests| B(API Server)
    B -->|Coordinate tasks| C(Scheduler)
    B -->|Manages state| D(Controller Manager)
    C -->|Asign tasks| E(Node Manager)
    D -->|Reconciliation| E
    B -->|Storage| F[(Datastore)]
    E -->|Update state| F
```

## 📂 Project structure

```mermaid
graph TD
    A[cmd/main.go] --> B[pkg/api]
    A --> C[pkg/scheduler]
    A --> D[pkg/controller]
    A --> E[pkg/node]
    A --> F[pkg/datastore]
```

### 🗓️ Scheduler

**Class diagram**
```mermaid
classDiagram
    class DefaultScheduler {
        +availableNodes: Node[]
        +schedule(task: Task, nodes: Node[]): Node
    }
    class Task {
        +id: string
    }
    class Node {
        +id: string
        +healthy: boolean
    }
    
    DefaultScheduler --> "0..*" Node : contains
    DefaultScheduler ..> Task : schedules
    Task --> "1" Node : assignedTo
```

**Sequence diagram**
```mermaid
flowchart LR
    T[Task Request] --> S[DefaultScheduler]
    S --> CN{Are Nodes Available?}
    CN -- Yes --> FN[Select First Node]
    FN --> RN[Return Selected Node]
    CN -- No --> EN[Return Error]
```

### 🎮 Controller

**Class diagram**

```mermaid
classDiagram
    class DefaultScheduler {
        +schedule(task: Task, nodes: Node[]): (Node, error)
    }
    class Task {
        +id: string
    }
    class Node {
        +id: string
    }
    class ControllerManager {
        -scheduler: DefaultScheduler
        -tasks: Task[]
        -nodes: Node[]
        +Run(stopCh)
        -reconcile()
    }
    
    ControllerManager --> DefaultScheduler : uses
    ControllerManager --> "0..*" Task : contains
    ControllerManager --> "0..*" Node : contains
```

**Sequence diagram**

```mermaid
sequenceDiagram
    participant CM as ControllerManager
    participant DS as DefaultScheduler

    Note over CM: Timer Tick (every 5 seconds)
    CM->>CM: Call reconcile()
    loop For each Task in tasks
        CM->>DS: schedule(task, nodes)
        DS-->>CM: returns assigned Node / error
        CM->>CM: Log assignment result
    end
```

### 👮‍♀️ Node manager

**Class diagram**

```mermaid
classDiagram
    class Node {
        +ID: string
        +Healthy: bool
    }
    class Manager {
        -nodes: map[string]Node
        -mu: RWMutex
        +Register(n: Node) error
        +GetNodes() []Node
        +UpdateHealth(nodeID: string, healthy: bool) error
    }

    Manager --> Node : uses
```

### 📀 Datastore

**Class diagram**

```mermaid
classDiagram
    class InMemoryDatastore {
        -nodes: map[string]node.Node
        -tasks: map[string]scheduler.Task
        -mu: RWMutex
        +SaveNode(n: node.Node) error
        +GetNodes() []node.Node
        +SaveTask(t: scheduler.Task) error
        +GetTasks() []scheduler.Task
    }
```