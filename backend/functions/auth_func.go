package functions

import (
	"database/sql"
	"net/http"
	"os"
	"product-catalogue/config"
	"product-catalogue/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	credentials.Email = strings.TrimSpace(credentials.Email)

	if credentials.Email == "" || credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Enter Password and Email!!"})
		return
	}

	if errEmail := utils.IsValidEmail(credentials.Email); errEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errEmail.Error()})
		return
	}

	if errPwd := utils.IsValidPassword(credentials.Password); errPwd != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errPwd.Error()})
		return
	}

	var id int
	var name, pwd_hash, role string
	var supplierID sql.NullInt64
	err := config.DB.QueryRow("select user_id, name, password_hash, role, supplier_id from users where email = ?", credentials.Email).Scan(&id, &name, &pwd_hash, &role, &supplierID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return

	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database related error"})
		return
	}
	if !utils.CheckPwd(pwd_hash, credentials.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	secret_k := os.Getenv("JWT_SECRET_KEY")
	if secret_k == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT not configured"})
		return
	}

	claims := jwt.MapClaims{
		"user_id": id,
		"name":    name,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	if role == "supplier_admin" && supplierID.Valid {
		claims["supplier_id"] = supplierID.Int64
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed_token, err := token.SignedString([]byte(secret_k))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token Creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "role": role, "token": signed_token})

}
