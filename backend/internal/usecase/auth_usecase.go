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
	var targetUser *entity.User

	if subdomain != "" {
		tenant, err := u.authRepo.GetTenantBySubdomain(ctx, subdomain)
		if err != nil {
			return "", nil, errors.New("invalid tenant subdomain")
		}
		targetUser, err = u.authRepo.GetUserByEmail(ctx, tenant.ID, email)
		if err != nil {
			return "", nil, errors.New("invalid credentials")
		}
	} else {
		// General email login across tenants
		targetUser, err = u.authRepo.GetUserByEmailAnyTenant(ctx, email)
		if err != nil {
			// Auto-seed default admin user yuvrajbisht41@gmail.com if DB is fresh
			if email == "yuvrajbisht41@gmail.com" || email == "admin@globallogistics.io" {
				if _, user, regErr := u.RegisterTenant(ctx, "Global Logistics Corp", "globallogistics", email, password); regErr == nil {
					targetUser = user
				} else {
					// If globallogistics tenant already exists, try looking up again
					if t, err := u.authRepo.GetTenantBySubdomain(ctx, "globallogistics"); err == nil {
						targetUser, _ = u.authRepo.GetUserByEmail(ctx, t.ID, email)
					}
				}
			}
		}
	}

	if targetUser == nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(targetUser.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	claims := &CustomJWTClaims{
		UserID:   targetUser.ID,
		TenantID: targetUser.TenantID,
		Role:     targetUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   targetUser.ID.String(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString(u.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, targetUser, nil
}
