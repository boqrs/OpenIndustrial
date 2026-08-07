package console

// WorkspaceType defines the category of a workspace, tailored to a specific role or function.
type WorkspaceType string

const (
	WorkspaceManagement  WorkspaceType = "ManagementCenter"
	WorkspaceProduction  WorkspaceType = "ProductionCenter"
	WorkspaceEquipment   WorkspaceType = "EquipmentCenter"
	WorkspaceQuality     WorkspaceType = "QualityCenter"
	WorkspaceWarehouse   WorkspaceType = "WarehouseCenter"
	WorkspaceCustomer    WorkspaceType = "CustomerPortal"
)

// Workspace is a configurable environment that provides a specific set of tools,
// dashboards, and menus for a user role. It is the top-level container for the user interface.
type Workspace struct {
	ID    string        `json:"id"`
	OrgID string        `json:"org_id"`
	Name  string        `json:"name"`
	Type  WorkspaceType `json:"type"`
}