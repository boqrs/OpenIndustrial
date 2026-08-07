package resource

// Type defines the category of a resource.
type Type string

const (
	TypeFactory  Type = "factory"
	TypeWorkshop Type = "workshop"
	TypeLine     Type = "line"
	TypeStation  Type = "station"
	TypeGateway  Type = "gateway"
	TypeDevice   Type = "device"
	TypeProduct  Type = "product"
	TypePoint    Type = "point"
)

// Status defines the lifecycle status of a resource.
type Status string

const (
	StatusActive      Status = "active"
	StatusInactive    Status = "inactive"
	StatusMaintenance Status = "maintenance"
	StatusRetired     Status = "retired"
)

// RelationType defines the nature of the connection between two resources.
type RelationType string

const (
	// Spatial relationships
	Contains   RelationType = "contains"
	LocatedAt  RelationType = "located_at"
	InstalledIn RelationType = "installed_in"

	// Control relationships
	Controls     RelationType = "controls"
	ControlledBy RelationType = "controlled_by"

	// Data relationships
	Reports RelationType = "reports"
	Measures RelationType = "measures"

	// Product/Process relationships
	ProducedBy  RelationType = "produced_by"
	AssembledAt RelationType = "assembled_at"
	ConsumedBy  RelationType = "consumed_by"
	UsedBy      RelationType = "used_by"
)

// RelationDirection defines the directionality of a relation.
type RelationDirection string

const (
	Directional   RelationDirection = "directional"
	Bidirectional RelationDirection = "bidirectional"
)