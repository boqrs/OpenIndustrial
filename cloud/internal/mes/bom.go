package mes

// BOM (Bill of Materials) defines the list of materials required to produce one unit of a product model.
type BOM struct {
	ID             string
	ProductModelID string // Foreign key to product.ProductModel
	Version        string
	IsActive       bool
}

// BOMItem is a single line item in a BOM.
type BOMItem struct {
	ID         string
	BOMID      string // Foreign key to BOM
	MaterialID string // Foreign key to wms.Material
	Quantity   float64
	Unit       string
}