package console

import "encoding/json"

// Dashboard is a container for a collection of widgets, arranged in a specific layout.
// It is the primary view within a Workspace.
type Dashboard struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	// Layout could be a JSON object defining the grid structure and widget placement.
	Layout      json.RawMessage `json:"layout,omitempty"`
}

// Widget represents a single, configurable component on a dashboard that displays
// a specific piece of information (e.g., a chart, a number, a list).
type Widget struct {
	ID          string `json:"id"`
	DashboardID string `json:"dashboard_id"`
	// Type is the identifier for the widget's frontend component (e.g., "kpi", "timeseries_chart").
	Type        string `json:"type"`
	// Config contains the specific settings for this widget instance, such as the
	// metric to display, the time range, or the target resource.
	Config      json.RawMessage `json:"config"`
}