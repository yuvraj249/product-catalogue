package entity

import (
	"time"

	"github.com/google/uuid"
)

type LocationType string

const (
	LocationInternal      LocationType = "Internal"
	LocationCustomer      LocationType = "Customer"
	LocationVendor        LocationType = "Vendor"
	LocationInventoryLoss LocationType = "InventoryLoss"
	LocationProduction    LocationType = "Production"
)

type Location struct {
	ID        uuid.UUID    `json:"id"`
	TenantID  uuid.UUID    `json:"tenant_id"`
	Name      string       `json:"name"`
	Type      LocationType `json:"type"`
	IsScrap   bool         `json:"is_scrap"`
	CreatedAt time.Time    `json:"created_at"`
}
