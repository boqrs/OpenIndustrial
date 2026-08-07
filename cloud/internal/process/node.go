package process

// ProcessNodeType defines the category of action a node represents.
type ProcessNodeType string

const (
	// NodeTypeStation represents a physical or logical work station.
	NodeTypeStation ProcessNodeType = "station"
	// NodeTypeInspection represents a quality control or inspection point.
	NodeTypeInspection ProcessNodeType = "inspection"
	// NodeTypeManual represents a task performed by a human operator.
	NodeTypeManual ProcessNodeType = "manual"
	// NodeTypeAutomatic represents an automated task performed by a machine or system.
	NodeTypeAutomatic ProcessNodeType = "automatic"
	// NodeTypeApproval represents a decision or approval point in the process.
	NodeTypeApproval ProcessNodeType = "approval"
)

// ProcessNode represents a single step or stage within a ProcessDefinition.
// It's a vertex in the process graph.
type ProcessNode struct {
	ID                string          `json:"id"`
	// ProcessDefinitionID links the node to its parent process definition.
	ProcessDefinitionID string          `json:"process_definition_id"`
	Name              string          `json:"name"`
	Type              ProcessNodeType `json:"type"`
	// Order provides a general sequence but the graph is defined by Edges.
	Order             int             `json:"order"`
}

// ProcessEdge represents a directed connection between two ProcessNodes.
// It defines the flow of the process and can have conditions.
type ProcessEdge struct {
	ID                  string `json:"id"`
	ProcessDefinitionID string `json:"process_definition_id"`
	FromNodeID          string `json:"from_node_id"`
	ToNodeID            string `json:"to_node_id"`
	// Condition is an expression that must evaluate to true for this edge to be taken.
	// e.g., "result == 'PASS'" or "data.temperature < 60".
	Condition           string `json:"condition,omitempty"`
}