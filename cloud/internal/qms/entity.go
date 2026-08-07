package qms

import (
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/trace"
	"github.com/OpenGongChang/OpenIndustrial/cloud/internal/wms"
)

// ProductPassport is a comprehensive digital record of a single product instance's entire lifecycle.
// It aggregates data from manufacturing, quality, material, and usage history.
type ProductPassport struct {
	ProductID string `json:"product_id"`

	// ManufactureHistory contains the sequence of production steps from the trace module.
	ManufactureHistory []trace.TraceRecord `json:"manufacture_history"`

	// MaterialGenealogy lists all material batches consumed during production.
	MaterialGenealogy []wms.MaterialConsumption `json:"material_genealogy"`

	// QualityHistory contains all inspection records and their results.
	QualityHistory []InspectionRecord `json:"quality_history"`

	// UsageHistory could contain a summary of telemetry data from the product's operational life.
	// (This would be populated by the analytics platform in the next stage).
	UsageHistory interface{} `json:"usage_history,omitempty"`
}