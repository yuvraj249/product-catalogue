package models

type User struct {
	UserID       int    `json:"user_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
	SupplierID   *int   `json:"supplier_id,omitempty"`
}
