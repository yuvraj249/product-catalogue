package functions

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"product-catalogue/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	refreshTokenExpiry = 7 * 24 * time.Hour // 7 days
	refreshCookieName  = "refresh_token"
)

func generateRandomToken() (string, error) {
	b := make([]byte, 64)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func generateFamilyID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func issueAccessToken(userID int, name, role string, supplierID sql.NullInt64) (string, error) {
	secretK := os.Getenv("JWT_SECRET_KEY")
	if secretK == "" {
		return "", fmt.Errorf("JWT not configured")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"name":    name,
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	if role == "supplier_admin" && supplierID.Valid {
		claims["supplier_id"] = supplierID.Int64
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretK))
}

func createRefreshToken(ctx context.Context, userID int, familyID string) (string, error) {
	rawToken, err := generateRandomToken()
	if err != nil {
		return "", err
	}

	tokenHash := hashToken(rawToken)
	expiresAt := time.Now().Add(refreshTokenExpiry)

	_, err = config.DB.ExecContext(ctx,
		"INSERT INTO refresh_tokens (token_hash, user_id, family_id, expires_at) VALUES (?, ?, ?, ?)",
		tokenHash, userID, familyID, expiresAt,
	)
	if err != nil {
		return "", err
	}

	return rawToken, nil
}

func setRefreshCookie(c *gin.Context, token string) {
	secure := os.Getenv("APP_ENV") != "test"
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		refreshCookieName,
		token,
		int(refreshTokenExpiry.Seconds()),
		"/",
		"",
		secure,
		true, // HttpOnly
	)
}

func clearRefreshCookie(c *gin.Context) {
	secure := os.Getenv("APP_ENV") != "test"
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(refreshCookieName, "", -1, "/", "", secure, true)
}

func RefreshToken(c *gin.Context) {
	rawToken, err := c.Cookie(refreshCookieName)
	if err != nil || rawToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	tokenHash := hashToken(rawToken)

	ctx, cancel := CtxTimeout(c)
	defer cancel()

	var id, userID int
	var familyID string
	var isUsed bool
	var expiresAt time.Time

	err = config.DB.QueryRowContext(ctx,
		"SELECT id, user_id, family_id, is_used, expires_at FROM refresh_tokens WHERE token_hash = ?",
		tokenHash,
	).Scan(&id, &userID, &familyID, &isUsed, &expiresAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Token reuse detected — revoke entire family
	if isUsed {
		config.DB.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE family_id = ?", familyID)
		clearRefreshCookie(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token reuse detected, all sessions revoked"})
		return
	}

	// Token expired
	if time.Now().After(expiresAt) {
		config.DB.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE id = ?", id)
		clearRefreshCookie(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
		return
	}

	// Mark current token as used
	config.DB.ExecContext(ctx, "UPDATE refresh_tokens SET is_used = TRUE WHERE id = ?", id)

	// Fetch user info for new access token
	var name, role string
	var supplierID sql.NullInt64
	err = config.DB.QueryRowContext(ctx,
		"SELECT name, role, supplier_id FROM users WHERE user_id = ?", userID,
	).Scan(&name, &role, &supplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	// Issue new access token
	accessToken, err := issueAccessToken(userID, name, role, supplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create access token"})
		return
	}

	// Issue new refresh token in same family
	newRefresh, err := createRefreshToken(ctx, userID, familyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create refresh token"})
		return
	}

	setRefreshCookie(c, newRefresh)
	c.JSON(http.StatusOK, gin.H{
		"message": "token refreshed",
		"token":   accessToken,
		"role":    role,
	})
}

func Logout(c *gin.Context) {
	ctx, cancel := CtxTimeout(c)
	defer cancel()

	rawToken, _ := c.Cookie(refreshCookieName)
	if rawToken != "" {
		tokenHash := hashToken(rawToken)

		var familyID string
		err := config.DB.QueryRowContext(ctx,
			"SELECT family_id FROM refresh_tokens WHERE token_hash = ?", tokenHash,
		).Scan(&familyID)
		if err == nil {
			config.DB.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE family_id = ?", familyID)
		}
	}

	clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
