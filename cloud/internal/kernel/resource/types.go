package resource

// ResourceType is a string representing the type of a resource.
type ResourceType string

const (
	// Factory-related resource types
	ResourceTypeFactory        ResourceType = "FACTORY"
	ResourceTypeWorkshop       ResourceType = "WORKSHOP"
	ResourceTypeProductionLine ResourceType = "PRODUCTION_LINE"
	ResourceTypeProductionCell ResourceType = "PRODUCTION_CELL"
	ResourceTypeWorkCenter     ResourceType = "WORK_CENTER"

	// Device-related resource types
	ResourceTypeDevice     ResourceType = "DEVICE"
	ResourceTypeDeviceType ResourceType = "DEVICE_TYPE"

	// Product-related resource types
	ResourceTypeProductModel    ResourceType = "PRODUCT_MODEL"
	//ResourceTypeProductInstance ResourceType = "PRODUCT_INSTANCE"
)