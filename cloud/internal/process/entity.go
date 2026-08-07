package process

// NodeResource defines the binding between a process node and a specific
// resource required to perform the task. This links the logical process
// to the physical world represented in the Resource Graph.
type NodeResource struct {
	// NodeID is the ID of the process node.
	NodeID string `json:"node_id"`
	// ResourceID is the ID of the required resource (e.g., a specific machine, station, or tool).
	ResourceID string `json:"resource_id"`
	// UsageType describes how the resource is used (e.g., "executes_task", "consumes_material").
	UsageType string `json:"usage_type"`
}