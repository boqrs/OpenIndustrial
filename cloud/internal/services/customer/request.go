package customer

type CreateRequest struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	ContactName  string `json:"contactName"`
	ContactEmail string `json:"contactEmail"`
	ContactPhone string `json:"contactPhone"`
	Address      string `json:"address"`
}

type UpdateRequest struct {
	Name         *string `json:"name,omitempty"`
	Type         *string `json:"type,omitempty"`
	ContactName  *string `json:"contactName,omitempty"`
	ContactEmail *string `json:"contactEmail,omitempty"`
	ContactPhone *string `json:"contactPhone,omitempty"`
	Address      *string `json:"address,omitempty"`
}

type ListRequest struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}
