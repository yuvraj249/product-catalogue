package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"backend/internal/domain/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserContextKey   contextKey = "auth.user.claims"
	TenantContextKey contextKey = "auth.tenant.id"
)

type CustomJWTClaims struct {
	UserID   uuid.UUID       `json:"user_id"`
	TenantID uuid.UUID       `json:"tenant_id"`
	Role     entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"unauthorized: missing bearer token"}`, http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &CustomJWTClaims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing algorithm")
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error":"unauthorized: token invalid or expired"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			ctx = context.WithValue(ctx, TenantContextKey, claims.TenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
