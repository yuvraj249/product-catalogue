package models

type Supplier struct {
	SupplierID  int    `json:"supplier_id"`
	Name        string `json:"name"`
	ContactInfo int    `json:"contact_info,omitempty"`
	Email       string `json:"email,omitempty"`
	Comapany    string `json:"company,omitempty"`
}
