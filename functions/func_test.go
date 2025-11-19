package functions_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-catalogue/functions"
	"strings"
	"testing"

	"product-catalogue/utils"

	"github.com/gin-gonic/gin"
)

func InvalidJSONResp(t *testing.T, jsonResp string) map[string]interface{} {
	var mp map[string]interface{}
	if err := json.Unmarshal([]byte(jsonResp), &mp); err != nil {
		t.Errorf("invalid json response : %v", err)
	}
	return mp
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.POST("/auth/login", functions.Login)
	request, err := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader("{"))
	if err != nil {
		t.Errorf("Error %v while making new request", err)
	}
	request.Header.Set("Content-Type", "application/json")
	browser := httptest.NewRecorder()
	r.ServeHTTP(browser, request)
	if browser.Code != http.StatusBadRequest {
		t.Errorf("got output %d expected %d", browser.Code, http.StatusBadRequest)
	}

}

func TestValidatorsInvalid(t *testing.T) {
	if err := utils.IsValidEmail("user1"); err == nil {
		t.Fatal("expected error for invalid email, got nil")
	}
	if err := utils.IsValidPassword("user1"); err == nil {
		t.Fatal("expected error for invalid/short password, got nil")
	}
}

func TestValidatorsValid(t *testing.T) {
	if err := utils.IsValidEmail("User@gmail.com"); err != nil {
		t.Fatalf("expected email to be valid but got error : %v", err)
	}
	if err := utils.IsValidPassword("User@1234"); err != nil {
		t.Fatalf("expected password to be valid got error : %v", err)
	}
}
