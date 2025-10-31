package models

type Product struct {
	ProductID int                 `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductDescription string  `json:"product_description,omitempty"`
	ProductCost int `json:"product_cost"`
	ProductCategoryID *int `json:"product_category_id,omitempty"`
	ProductSupplierID *int `json:"product_supplier_id,omitempty"`
	DiscountType *string `json:"discount_type,omitempty"`
	DiscountValue *float64 `json:"discount_value,omitempty"`

}