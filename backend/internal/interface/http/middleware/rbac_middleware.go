package middleware

import (
	"net/http"

	"backend/internal/domain/entity"
)

func RequireRole(allowedRoles ...entity.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*CustomJWTClaims)
			if !ok || claims == nil {
				http.Error(w, `{"error":"forbidden: security context missing"}`, http.StatusForbidden)
				return
			}

			permitted := false
			for _, role := range allowedRoles {
				if claims.Role == role || claims.Role == entity.RoleSuperAdmin {
					permitted = true
					break
				}
			}

			if !permitted {
				http.Error(w, `{"error":"forbidden: insufficient role authorization"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
