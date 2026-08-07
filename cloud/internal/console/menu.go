package console

// MenuItem represents a single navigable item within a workspace's menu structure.
type MenuItem struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	ParentID    string `json:"parent_id,omitempty"` // For creating nested menus
	Name        string `json:"name"`
	Icon        string `json:"icon,omitempty"`
	// Path could be the route to a specific dashboard or page.
	Path        string `json:"path"`
	Order       int    `json:"order"`
}