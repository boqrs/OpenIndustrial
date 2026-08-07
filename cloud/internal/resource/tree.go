package resource

// Node represents a node in the resource tree.
// It contains a resource and its children, allowing for hierarchical representation.
type Node struct {
	*Resource
	Children []*Node
}