package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"backend/internal/domain/entity"
	"backend/internal/domain/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CustomJWTClaims struct {
	UserID   uuid.UUID       `json:"user_id"`
	TenantID uuid.UUID       `json:"tenant_id"`
	Role     entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type AuthUsecase struct {
	authRepo  repository.AuthRepository
	jwtSecret []byte
}

func NewAuthUsecase(authRepo repository.AuthRepository, jwtSecret []byte) *AuthUsecase {
	return &AuthUsecase{
		authRepo:  authRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *AuthUsecase) RegisterTenant(ctx context.Context, tenantName, subdomain, adminEmail, password string) (*entity.Tenant, *entity.User, error) {
	existingTenant, err := u.authRepo.GetTenantBySubdomain(ctx, subdomain)
	if err == nil && existingTenant != nil {
		return nil, nil, errors.New("subdomain already registered")
	}

	tenant := &entity.Tenant{
		ID:        uuid.New(),
		Name:      tenantName,
		Subdomain: subdomain,
	}

	if err := u.authRepo.CreateTenant(ctx, tenant); err != nil {
		return nil, nil, fmt.Errorf("failed to register tenant: %w", err)
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		ID:           uuid.New(),
		TenantID:     tenant.ID,
		Email:        adminEmail,
		PasswordHash: string(hashedPwd),
		Role:         entity.RoleTenantAdmin,
	}

	if err := u.authRepo.CreateUser(ctx, user); err != nil {
		return nil, nil, fmt.Errorf("failed to create tenant admin: %w", err)
	}

	return tenant, user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, subdomain, email, password string) (token string, user *entity.User, err error) {
	tenant, err := u.authRepo.GetTenantBySubdomain(ctx, subdomain)
	if err != nil {
		return "", nil, errors.New("invalid tenant subdomain")
	}

	user, err = u.authRepo.GetUserByEmail(ctx, tenant.ID, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	claims := &CustomJWTClaims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(u.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, user, nil
}
