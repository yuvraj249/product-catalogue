package repository

import (
	"context"

	"backend/internal/domain/entity"
	"github.com/google/uuid"
)

type AuthRepository interface {
	CreateTenant(ctx context.Context, tenant *entity.Tenant) error
	GetTenantByID(ctx context.Context, tenantID uuid.UUID) (*entity.Tenant, error)
	GetTenantBySubdomain(ctx context.Context, subdomain string) (*entity.Tenant, error)
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*entity.User, error)
}
