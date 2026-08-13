package param

// CreateFactoryRequest defines the parameters for creating a new factory.
type CreateFactoryRequest struct {
	Name    string `json:"name" binding:"required"`    // Name of the factory, will be stored in the Resource model.
	Code    string `json:"code" binding:"required"`    // Factory-specific code, e.g., "F001".
	Address string `json:"address" binding:"required"` // Physical address of the factory.
}