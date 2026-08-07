package resource

// ResourceType defines the type of a resource.
// In the future, this will be loaded from a database table for dynamic extension.
type ResourceType string

const (
	TypeOrg       ResourceType = "org"
	TypeFactory   ResourceType = "factory"
	TypeWorkshop  ResourceType = "workshop"
	TypeLine      ResourceType = "line"
	TypeStation   ResourceType = "station"
	TypeGateway   ResourceType = "gateway"
	TypeDevice    ResourceType = "device"
	TypeProduct   ResourceType = "product"
	TypeCustomer  ResourceType = "customer"
	TypeWarehouse ResourceType = "warehouse"
	TypeMaterial  ResourceType = "material"
)

// RelationType defines the nature of a relationship between two resources.
type RelationType string

const (
	RelationContains    RelationType = "contains"
	RelationConnects    RelationType = "connects"
	RelationProduces    RelationType = "produces"
	RelationInstalledOn RelationType = "installed_on"
	RelationUses        RelationType = "uses"
	RelationBelongs     RelationType = "belongs"
	RelationControls    RelationType = "controls"
)