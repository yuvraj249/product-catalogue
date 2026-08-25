package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"backend/internal/domain/entity"
	"github.com/google/uuid"
)

type AuthRepository struct {
	DB *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (r *AuthRepository) CreateTenant(ctx context.Context, tenant *entity.Tenant) error {
	if tenant.ID == uuid.Nil {
		tenant.ID = uuid.New()
	}
	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	query := `INSERT INTO tenants (id, name, subdomain, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.ExecContext(ctx, query, tenant.ID, tenant.Name, tenant.Subdomain, tenant.CreatedAt, tenant.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetTenantByID(ctx context.Context, tenantID uuid.UUID) (*entity.Tenant, error) {
	query := `SELECT id, name, subdomain, created_at, updated_at FROM tenants WHERE id = $1`
	var t entity.Tenant
	err := r.DB.QueryRowContext(ctx, query, tenantID).Scan(&t.ID, &t.Name, &t.Subdomain, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tenant by ID: %w", err)
	}
	return &t, nil
}

func (r *AuthRepository) GetTenantBySubdomain(ctx context.Context, subdomain string) (*entity.Tenant, error) {
	query := `SELECT id, name, subdomain, created_at, updated_at FROM tenants WHERE subdomain = $1`
	var t entity.Tenant
	err := r.DB.QueryRowContext(ctx, query, subdomain).Scan(&t.ID, &t.Name, &t.Subdomain, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tenant by subdomain: %w", err)
	}
	return &t, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `INSERT INTO users (id, tenant_id, email, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.DB.ExecContext(ctx, query, user.ID, user.TenantID, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*entity.User, error) {
	query := `SELECT id, tenant_id, email, password_hash, role, created_at, updated_at FROM users WHERE tenant_id = $1 AND email = $2`
	var u entity.User
	err := r.DB.QueryRowContext(ctx, query, tenantID, email).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}
	return &u, nil
}

func (r *AuthRepository) GetUserByEmailAnyTenant(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, tenant_id, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1 LIMIT 1`
	var u entity.User
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}
	return &u, nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, tenantID, userID uuid.UUID) (*entity.User, error) {
	query := `SELECT id, tenant_id, email, password_hash, role, created_at, updated_at FROM users WHERE tenant_id = $1 AND id = $2`
	var u entity.User
	err := r.DB.QueryRowContext(ctx, query, tenantID, userID).Scan(&u.ID, &u.TenantID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by ID: %w", err)
	}
	return &u, nil
}
