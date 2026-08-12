package param

import (
	"github.com/OpenIndustrial/cloud/internal/persistence/model"
)

// ResourceResponse is the standard API response for a resource, combining the
// core resource model with its dynamic attributes.
type ResourceResponse struct {
	*model.Resource
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}