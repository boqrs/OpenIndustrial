package productinstance

type CreateRequest struct {
	SN        string `json:"sn"`
	ProductID string `json:"product_id"`
	OrgID     string `json:"org_id"`
}

type Response struct {
	ID        string `json:"id"`
	SN        string `json:"sn"`
	ProductID string `json:"product_id"`
	State     string `json:"state"`
}