package models

type Stock_Movement struct {
	StockID      int    `json:"stock_id"`
	ProductID    int    `json:"product_id"`
	ProductName  string `json:"product_name"`
	Quantity     int    `json:"quantity"`
	MovementType string `json:"movement_type"`
	Reason       string `json:"reason,omitempty"`
	PerformedBy  int    `json:"performed_by"`
	Username     string `json:"username"`
}
