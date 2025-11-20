package functions_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-catalogue/functions"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

func InvalidJSONResp(t *testing.T, jsonResp string) map[string]interface{} {
	var mp map[string]interface{}
	if err := json.Unmarshal([]byte(jsonResp), &mp); err != nil {
		t.Errorf("invalid json response : %v", err)
	}
	return mp
}

func MakeTestJWT(role string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
	})
	signed, _ := token.SignedString([]byte("test_secret"))
	return signed
}

func DemoAuthMiddleware(secret_k string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret_k), nil
		})
		if err != nil || !parsed.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims := parsed.Claims.(jwt.MapClaims)
		role := claims["role"].(string)

		c.Set("role", role)
		c.Next()
	}
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testcases := []struct {
		name         string
		reqBody      string
		expectedCode int
	}{

		{
			name:         "Invalid JSON",
			reqBody:      `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid Email",
			reqBody:      `{"email":"abcdefg","password":"Abcd@1234"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Invalid Password",
			reqBody:      `{"email":"unittest@mail.com","password":"passwor"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.POST("/auth/login", functions.Login)
			request, err := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tc.reqBody))
			if err != nil {
				t.Fatalf("Error %v while making new request", err)
			}
			request.Header.Set("Content-Type", "application/json")
			browser := httptest.NewRecorder()
			r.ServeHTTP(browser, request)
			if browser.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d; body=%s", tc.name, tc.expectedCode, browser.Code, browser.Body.String())
			}

		})
	}

}

func TestCreateCatg(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret_k := "test_secret"
	testcases := []struct {
		name         string
		role         string
		reqBody      string
		expectedCode int
	}{
		{
			name:         "System Admin Case (allowed to create category)",
			role:         "system_admin",
			reqBody:      `{}`,
			expectedCode: http.StatusBadRequest,
		}, {
			name:         "Supplier Admin Case (forbiddden)",
			role:         "supplier_admin",
			reqBody:      `{}`,
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			token := MakeTestJWT(tc.role)
			r := gin.New()
			r.Use(DemoAuthMiddleware(secret_k))
			r.POST("/category", functions.CreateCategory)
			request, err := http.NewRequest(http.MethodPost, "/category", strings.NewReader(tc.reqBody))
			if err != nil {
				t.Fatalf("Error %v while making new request", err)
			}
			request.Header.Set("Authorization", "Bearer "+token)
			request.Header.Set("Content-Type", "application/json")
			browser := httptest.NewRecorder()
			r.ServeHTTP(browser, request)
			if browser.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d; body=%s", tc.name, tc.expectedCode, browser.Code, browser.Body.String())
			}

		})
	}
}
