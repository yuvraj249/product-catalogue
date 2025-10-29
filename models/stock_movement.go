package models

type Stock_Movement struct {
	StockID   int    `json:"stock_id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Reason    string `json:"reason,omitempty"`
	ChangedBy *int   `json:"changed_by,omitempty"`
}
