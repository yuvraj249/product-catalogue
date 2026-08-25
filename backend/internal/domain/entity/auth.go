package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleSuperAdmin       UserRole = "SuperAdmin"
	RoleTenantAdmin      UserRole = "TenantAdmin"
	RoleWarehouseManager UserRole = "WarehouseManager"
	RoleCashier          UserRole = "Cashier"
	RoleAuditor          UserRole = "Auditor"
)

type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Subdomain string    `json:"subdomain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
