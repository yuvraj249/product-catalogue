package models

type Stock_Movement struct {
	StockID      int    `json:"stock_id"`
	ProductID    int    `json:"product_id"`
	Quantity     int    `json:"quantity"`
	MovementType string `json:"movement_type"`
	Reason       string `json:"reason,omitempty"`
	PerformedBy  int    `json:"performed_by"`
}
