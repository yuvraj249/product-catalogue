package functions_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"product-catalogue/config"
	"product-catalogue/routes"
	"product-catalogue/utils"
	"strings"
	"testing"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

var truncateTables = []string{
	"stock_movements",
	"products",
	"categories",
	"users",
	"suppliers",
}

func TruncateAll(t *testing.T) {
	t.Helper()
	if config.DB == nil {
		t.Fatal("DB is nil")
	}
	config.DB.Exec("set foreign_key_checks=0")
	for _, tbl := range truncateTables {
		config.DB.Exec("truncate table " + tbl)
	}
	config.DB.Exec("set foreign_key_checks=1")
}

func FindMigrationFile() (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	root = filepath.Dir(root)
	path := filepath.Join(root, "migration", "init.sql")

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("migration file not found at %s", path)
	}
	return path, nil
}

func SeedAdmin(t *testing.T, email, pass string) {
	t.Helper()
	hash, _ := utils.HashPwd(pass)
	_, err := config.DB.Exec("insert into users(name,email,password_hash,role) values('AdminUser',?,?,?)", email, hash, "system_admin")
	if err != nil {
		t.Fatalf("failed seeding admin: %v", err)
	}
}

func LoginAndGetToken(t *testing.T, email, pass string) string {
	t.Helper()

	r := routes.SetupRouter()

	body := fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, pass)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("login failed: %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	token, ok := resp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("token missing in login response")
	}

	return token
}

func TestMain(m *testing.M) {
	root, _ := os.Getwd()
	root = filepath.Dir(root)
	envPath := filepath.Join(root, ".env")
	_ = godotenv.Load(envPath)
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN must be provided")
	}

	jwt := os.Getenv("JWT_SECRET_KEY")
	if jwt == "" {
		log.Fatal("JWT_SECRET_KEY missing in .env")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	config.DB = db

	mig, err := FindMigrationFile()
	if err != nil {
		log.Fatal(err)
	}

	if err := config.ExecMigration(db, mig); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	os.Exit(code)
}

func CreateSupplierViaFunc(t *testing.T, r http.Handler, adminToken, name, contact, email, company string) int64 {
	t.Helper()

	body := fmt.Sprintf(`{"name":"%s","contact_info":"%s","email":"%s","company":"%s"}`,
		name, contact, email, company)

	req := httptest.NewRequest(http.MethodPost, "/suppliers", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("CreateSupplierViaAPI: expected 200/201 got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("CreateSupplierViaAPI: invalid json response: %v body=%s", err, w.Body.String())
	}
	if v, ok := resp["supplier_id"]; ok {
		return int64(v.(float64))
	}
	if obj, ok := resp["supplier"].(map[string]interface{}); ok {
		if idv, ok := obj["supplier_id"]; ok {
			return int64(idv.(float64))
		}
	}
	for _, v := range resp {
		if n, ok := v.(float64); ok {
			return int64(n)
		}
	}
	t.Fatalf("CreateSupplierViaAPI: supplier_id not found in response: %s", w.Body.String())
	return 0
}

func CreateUserViaFunc(t *testing.T, r http.Handler, adminToken, name, email, password, role string, supplierID int64) {
	t.Helper()

	body := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s","role":"%s","supplier_id":%d}`,
		name, email, password, role, supplierID)

	req := httptest.NewRequest(http.MethodPost, "/users/supplier-admin", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("CreateUserViaAPI: expected 200/201 got %d body=%s", w.Code, w.Body.String())
	}
}

func createCategoryViaFunc(t *testing.T, r http.Handler, adminToken, name, desc string) int64 {
	t.Helper()
	body := fmt.Sprintf(`{"category_name":"%s","category_description":"%s"}`, name, desc)
	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("createCategoryViaAPI: expected 200/201 got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("createCategoryViaAPI: invalid json response: %v body=%s", err, w.Body.String())
	}

	if v, ok := resp["category_id"]; ok {
		return int64(v.(float64))
	}
	if obj, ok := resp["category"].(map[string]interface{}); ok {
		if idv, ok := obj["category_id"]; ok {
			return int64(idv.(float64))
		}
	}
	for _, v := range resp {
		if n, ok := v.(float64); ok {
			return int64(n)
		}
	}
	t.Fatalf("createCategoryViaAPI: category_id not found in response: %s", w.Body.String())
	return 0
}

func CreateProductViaFunc(t *testing.T, r http.Handler, token, name, desc string, cost float64, categoryID int64) int64 {
	t.Helper()
	body := fmt.Sprintf(`{"product_name":"%s","product_description":"%s","product_cost":%f,"product_category_id":%d}`,
		name, desc, cost, categoryID)

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("CreateProductViaAPI: expected 200/201 got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("CreateProductViaAPI: invalid json response: %v body=%s", err, w.Body.String())
	}
	if v, ok := resp["product_id"]; ok {
		return int64(v.(float64))
	}
	if obj, ok := resp["product"].(map[string]interface{}); ok {
		if idv, ok := obj["product_id"]; ok {
			return int64(idv.(float64))
		}
	}
	for _, v := range resp {
		if n, ok := v.(float64); ok {
			return int64(n)
		}
	}
	t.Fatalf("CreateProductViaAPI: product_id not found in response: %s", w.Body.String())
	return 0
}

func CreateStockViaFunc(t *testing.T, r http.Handler, token string, productID int64, quantity int, movementType, reason string) int64 {
	t.Helper()

	payload := map[string]interface{}{
		"product_id":    productID,
		"quantity":      quantity,
		"movement_type": movementType,
	}
	if reason != "" {
		payload["reason"] = reason
	}

	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("CreateStockViaFunc: json marshal failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/stock_movements", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Fatalf("CreateStockViaFunc: expected 200/201 got %d body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("CreateStockViaFunc: invalid json response: %v body=%s", err, w.Body.String())
	}
	if v, ok := resp["stock_movement_id"]; ok {
		return int64(v.(float64))
	}
	if v, ok := resp["movement_id"]; ok {
		return int64(v.(float64))
	}
	if obj, ok := resp["stock_movement"].(map[string]interface{}); ok {
		if idv, ok := obj["stock_movement_id"]; ok {
			return int64(idv.(float64))
		}
		if idv, ok := obj["movement_id"]; ok {
			return int64(idv.(float64))
		}
		if idv, ok := obj["id"]; ok {
			return int64(idv.(float64))
		}
	}
	for _, v := range resp {
		if n, ok := v.(float64); ok {
			return int64(n)
		}
	}

	t.Fatalf("CreateStockViaFunc: stock movement id not found in response: %s", w.Body.String())
	return 0
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
		{
			name:         "Empty email and Password",
			reqBody:      `{"email":"","password":""}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Valid credentials (system_admin)",
			reqBody:      `{"email":"yuvrajbisht41@gmail.com","password":"Yuvraj@2411"}`,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := routes.SetupRouter()
			request, err := http.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(tc.reqBody))
			if err != nil {
				t.Fatalf("Error %v while making new request", err)
			}
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d; body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}

		})
	}

}

func TestCreateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_td@test.com", "Admin@123"
	const supPass = "Sup@2411"
	const supEmail = "supadmin@sup.com"
	r := routes.SetupRouter()

	cases := []struct {
		name         string
		prepare      func(adminToken string)
		body         string
		expectedCode int
	}{
		{
			name:         "no token",
			prepare:      func(adminToken string) {},
			body:         `{"category_name":"Xsndin","category_description":"Yadaii"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "supplier_admin cannot create",
			prepare: func(adminToken string) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupTD", "111382328828", "sup_td@test.com", "Cfdsdsds")
				CreateUserViaFunc(t, r, adminToken, "SupUser", supEmail, supPass, "supplier_admin", sid)
			},
			body:         `{"category_name":"Xidjais","category_description":"Ysidisd"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid json",
			prepare: func(adminToken string) {

			},
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid category name",
			prepare:      func(adminToken string) {},
			body:         `{"category_name":"@#$","category_description":"valid"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "valid create",
			prepare:      func(adminToken string) {},
			body:         `{"category_name":"Electronics","category_description":"Devices"}`,
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			if tc.prepare != nil {
				tc.prepare(adminToken)
			}

			var supToken string

			req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(tc.body))

			switch tc.name {
			case "invalid json", "invalid category name", "valid create":
				if adminToken == "" {
					t.Fatalf("token empty for the subset: %s", tc.name)
				} else {

					req.Header.Set("Authorization", "Bearer "+adminToken)
				}
			case "supplier_admin cannot create":

				supToken = LoginAndGetToken(t, supEmail, supPass)

				req.Header.Set("Authorization", "Bearer "+supToken)

			case "no token":
				req.Header.Set("Authorization", "Bearer "+"")
			}

			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_td@test.com", "Admin@123"
	const supPass = "Sup@2411"
	const supEmail = "supadmin@sup.com"
	r := routes.SetupRouter()

	cases := []struct {
		name         string
		prepare      func(adminToken string)
		expectedCode int
		expectBody   string
	}{
		{
			name:         "no token",
			prepare:      func(adminToken string) {},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "system_admin authorized",
			prepare: func(adminToken string) {
				_ = createCategoryViaFunc(t, r, adminToken, "Electronics", "Devices")
				_ = createCategoryViaFunc(t, r, adminToken, "Home", "Appliances")
			},
			expectedCode: http.StatusOK,
			expectBody:   "Electronics",
		},
		{
			name: "supplier_admin authorized",
			prepare: func(adminToken string) {
				_ = createCategoryViaFunc(t, r, adminToken, "Electronics", "Devices")
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupTD", "111327727272", "sup_td@test.com", "Cuasaj")
				CreateUserViaFunc(t, r, adminToken, "SupUser", supEmail, supPass, "supplier_admin", sid)
			},
			expectedCode: http.StatusOK,
			expectBody:   "Electronics",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			if tc.prepare != nil {
				tc.prepare(adminToken)
			}

			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
			switch tc.name {
			case "no token":
				req.Header.Set("Authorization", "Bearer "+"")
			case "supplier_admin authorized":
				supToken := LoginAndGetToken(t, supEmail, supPass)
				req.Header.Set("Authorization", "Bearer "+supToken)
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("[%s] expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectedCode == http.StatusOK && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("[%s] expected body to contain %q; got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestGetCategoryByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_td@test.com", "Admin@123"
	const supPass = "Sup@2411"
	const supEmail = "supadmin@sup.com"
	r := routes.SetupRouter()

	cases := []struct {
		name         string
		prepare      func(adminToken string) int64
		url          string
		expectedCode int
		expectBody   string
	}{
		{
			name:         "invalid id param",
			prepare:      func(adminToken string) int64 { return 0 },
			url:          "/categories/abc",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "not found",
			prepare: func(adminToken string) int64 {
				return 0
			},
			url:          "/categories/9999999",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "found",
			prepare: func(adminToken string) int64 {
				return createCategoryViaFunc(t, r, adminToken, "Clothing", "Devices")
			},
			url:          "",
			expectedCode: http.StatusOK,
			expectBody:   "Clothing",
		},
		{
			name: "system_admin allowed",
			prepare: func(adminToken string) int64 {
				return createCategoryViaFunc(t, r, adminToken, "Clothing", "Devices")
			},
			url:          "",
			expectedCode: http.StatusOK,
			expectBody:   "Clothing",
		},
		{
			name: "supplier_admin allowed",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupTD", "111323327272", "sup_td@test.com", "Cdisjdij")
				CreateUserViaFunc(t, r, adminToken, "SupUser", supEmail, supPass, "supplier_admin", sid)
				return createCategoryViaFunc(t, r, adminToken, "Clothing", "Devices")
			},
			url:          "",
			expectedCode: http.StatusOK,
			expectBody:   "Clothing",
		},
		{
			name: "no token",
			prepare: func(adminToken string) int64 {
				return 0
			},
			url:          "/categories/1",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)

			var id int64
			if tc.prepare != nil {
				id = tc.prepare(adminToken)
			}
			url := tc.url
			if url == "" && id != 0 {
				url = fmt.Sprintf("/categories/%d", id)
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			switch tc.name {
			case "supplier_admin allowed":
				supToken := LoginAndGetToken(t, supEmail, supPass)
				req.Header.Set("Authorization", "Bearer "+supToken)
			case "no token":
				req.Header.Set("Authorization", "Bearer "+"")
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectedCode == http.StatusOK && tc.expectBody != "" && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("[%s] expected body to contain %q; got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_td@test.com", "Admin@123"
	const supPass = "Sup@2411"
	const supEmail = "supadmin@sup.com"
	r := routes.SetupRouter()

	cases := []struct {
		name         string
		prepare      func(adminToken string) int64
		body         string
		expectedCode int
	}{
		{
			name: "no fields",
			prepare: func(adminToken string) int64 {
				return createCategoryViaFunc(t, r, adminToken, "Clothing", "Devices")
			},
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupTD", "111328329", "sup_td@test.com", "Cdidind")
				CreateUserViaFunc(t, r, adminToken, "SupUser", supEmail, supPass, "supplier_admin", sid)
				return createCategoryViaFunc(t, r, adminToken, "Clothing", "Devices")
			},
			body:         `{"category_name":"Movies"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid id bad request",
			prepare: func(adminToken string) int64 {
				return 0
			},
			body:         `{"category_name":"NewName"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "valid update",
			prepare: func(adminToken string) int64 {
				return createCategoryViaFunc(t, r, adminToken, "Clothing", "Devices")
			},
			body:         `{"category_name":"Hardware"}`,
			expectedCode: http.StatusOK,
		},
		{
			name: "duplicate category_name",
			prepare: func(adminToken string) int64 {
				return createCategoryViaFunc(t, r, adminToken, "Clothing", "Other")
			},
			body:         `{"category_name":"Clothing"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate category_description",
			prepare: func(adminToken string) int64 {
				return createCategoryViaFunc(t, r, adminToken, "Another", "Devices")
			},
			body:         `{"category_description":"Devices"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)

			id := tc.prepare(adminToken)
			var url string
			if id == 0 && tc.name == "invalid id bad request" {
				url = "/categories/abc"
			} else {
				url = fmt.Sprintf("/categories/%d", id)
			}

			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "forbidden for supplier_admin":
				supToken := LoginAndGetToken(t, supEmail, supPass)
				req.Header.Set("Authorization", "Bearer "+supToken)
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_td@test.com", "Admin@123"
	const supPass = "Sup@2411"
	const supEmail = "supadmin@sup.com"
	r := routes.SetupRouter()

	cases := []struct {
		name         string
		prepare      func(adminToken string) int64
		expectedCode int
	}{
		{
			name: "invalid id",
			prepare: func(adminToken string) int64 {
				return 0
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupTD", "1114234822", "sup_td@test.com", "Cdinasdn")
				CreateUserViaFunc(t, r, adminToken, "SupUser", supEmail, supPass, "supplier_admin", sid)
				return createCategoryViaFunc(t, r, adminToken, "Clothing", "Devices")
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "delete not found",
			prepare: func(adminToken string) int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful delete",
			prepare: func(adminToken string) int64 {
				return createCategoryViaFunc(t, r, adminToken, "Clothing", "Devices")
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)

			id := tc.prepare(adminToken)
			var url string
			if id == 0 && tc.name == "invalid id" {
				url = "/categories/xyz"
			} else {
				url = fmt.Sprintf("/categories/%d", id)
			}

			req := httptest.NewRequest(http.MethodDelete, url, nil)
			switch tc.name {
			case "forbidden for supplier_admin":
				supToken := LoginAndGetToken(t, supEmail, supPass)
				req.Header.Set("Authorization", "Bearer "+supToken)
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestCreateSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_create_sup@test.com", "Admin@123"
	const supPass = "Sup@2411"
	const supEmail = "sup_del@test.com"
	r := routes.SetupRouter()

	cases := []struct {
		name         string
		prepare      func(adminToken string)
		body         string
		expectedCode int
	}{
		{
			name:         "token missing",
			prepare:      func(adminToken string) {},
			body:         `{"name":"Ssadhah","contact_info":"34626323737","email":"s1p@test.com","company":"Co"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func(adminToken string) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupDexsl", "1113773232", supEmail, "Ccjnisnsins")
				CreateUserViaFunc(t, r, adminToken, "SupDelUser", supEmail, supPass, "supplier_admin", sid)
			},
			body:         `{"name":"Sahdshd","contact_info":"57343747347","email":"s1p@test.com","company":"Codks d"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalid payload",
			prepare:      func(adminToken string) {},
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "valid create",
			prepare:      func(adminToken string) {},
			body:         `{"name":"SCreate","contact_info":"9999999","email":"screate@test.com","company":"CreateCo"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid contact_info",
			prepare:      func(adminToken string) {},
			body:         `{"name":"SCreate","contact_info":"abcdef","email":"screate1@test.com","company":"CreateCo"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			prepare: func(adminToken string) {
				CreateSupplierViaFunc(t, r, adminToken, "sdshdsh", "5734737732", "screate@test.com", "CreateCo")
			},
			body:         `{"name":"SCreate","contact_info":"6757475399","email":"screate@test.com","company":"CreateCo"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			if adminToken == "" {
				t.Fatalf("admin login failed")
			}
			if tc.prepare != nil {
				tc.prepare(adminToken)
			}

			req := httptest.NewRequest(http.MethodPost, "/suppliers", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "token missing":
				req.Header.Set("Authorization", "Bearer "+"")
			case "forbidden for supplier_admin":
				supToken := LoginAndGetToken(t, supEmail, supPass)
				if supToken == "" {
					t.Fatalf("supplier login failed for subset: %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+supToken)
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestUpdateSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_update_sup@test.com", "Admin@123"
	const supPass = "Sup@2411"
	const supEmail = "sup_del@test.com"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string) int64
		body         string
		expectedCode int
	}{
		{
			name: "no fields",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Tester", "99999999", "tester@gmail.com", "Apple")
			},
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func(adminToken string) int64 {
				id := CreateSupplierViaFunc(t, r, adminToken, "Tester", "99999999", "tester@gmail.com", "Apple")
				sid := CreateSupplierViaFunc(t, r, adminToken, "OtherSup", "111271282828", supEmail, "OtherCo")
				CreateUserViaFunc(t, r, adminToken, "SupUser", supEmail, supPass, "supplier_admin", sid)
				return id
			},
			body:         `{"supplier_name":"Stubbs"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid id bad request",
			prepare: func(adminToken string) int64 {
				return 0
			},
			body:         `{"supplier_name":"NewName"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "valid update",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Testerrrr", "99999999", "tester@gmail.com", "Apple")
			},
			body:         `{"name":"Testerrrr", "contact_info":"989898989", "email":"tester@gmail.com", "company":"Apple"}`,
			expectedCode: http.StatusOK,
		},
		{
			name: "duplicate name",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Tester", "99999999", "tester@gmail.com", "Apple")
			},
			body:         `{"name":"Tester"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate contact_info",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Tester", "99999999", "tester2@gmail.com", "Apple")
			},
			body:         `{"contact_info":"99999999"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Tester", "88888888", "tester@gmail.com", "Apple")
			},
			body:         `{"email":"tester@gmail.com"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate company",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Tester", "77777777", "tester3@gmail.com", "Apple")
			},
			body:         `{"company":"Apple"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			if adminToken == "" {
				t.Fatalf("admin login failed for subset: %s", tc.name)
			}

			var id int64
			if tc.prepare != nil {
				id = tc.prepare(adminToken)
			}

			url := fmt.Sprintf("/suppliers/%d", id)
			if id == 0 && tc.name == "invalid id bad request" {
				url = "/suppliers/xyz"
			}

			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			if tc.name == "forbidden for supplier_admin" {
				supToken := LoginAndGetToken(t, supEmail, supPass)
				if supToken == "" {
					t.Fatalf("supplier login failed for subset: %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+supToken)
			} else {
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestDeleteSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const supPass = "Sup@2411"
	const supEmail = "sup_del@test.com"
	adminEmail, adminPass := "admin_delete_sup@test.com", "Admin@123"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string) int64
		expectedCode int
	}{
		{
			name: "invalid id",
			prepare: func(adminToken string) int64 {
				return 0
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func(adminToken string) int64 {
				id := CreateSupplierViaFunc(t, r, adminToken, "ToBeDeleted", "99999999", "tbd@test.com", "Apple")
				sid := CreateSupplierViaFunc(t, r, adminToken, "Other", "11111111", supEmail, "OtherCo")
				CreateUserViaFunc(t, r, adminToken, "SupDelUser", supEmail, supPass, "supplier_admin", sid)
				return id
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "delete not found",
			prepare: func(adminToken string) int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful delete",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Removable", "99999999", "rem@test.com", "Apple")
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			if adminToken == "" {
				t.Fatalf("admin login failed for subset: %s", tc.name)
			}

			id := tc.prepare(adminToken)
			url := fmt.Sprintf("/suppliers/%d", id)
			if id == 0 && tc.name == "invalid id" {
				url = "/suppliers/xyz"
			}

			req := httptest.NewRequest(http.MethodDelete, url, nil)

			if tc.name == "forbidden for supplier_admin" {
				supToken := LoginAndGetToken(t, supEmail, supPass)
				if supToken == "" {
					t.Fatalf("supplier login failed for subset: %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+supToken)
			} else {
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_getsup@test.com", "Admin@123"
	const supPass = "Sup@2411"
	const supEmail = "sup_del@test.com"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string)
		expectedCode int
		expectBody   string
	}{
		{
			name: "unauthorized (no token)",
			prepare: func(adminToken string) {
				CreateSupplierViaFunc(t, r, adminToken, "ABC", "11126363", "a@test.com", "CoA")
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "empty list for system_admin",
			prepare:      func(adminToken string) {},
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin sees all suppliers",
			prepare: func(adminToken string) {
				CreateSupplierViaFunc(t, r, adminToken, "ABC", "1111177717", "a@test.com", "CoA")
				CreateSupplierViaFunc(t, r, adminToken, "CFD", "2227171717", "b@test.com", "CoB")
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin sees all suppliers of same company",
			prepare: func(adminToken string) {
				CreateSupplierViaFunc(t, r, adminToken, "ABC", "111772227", "a@test.com", "CoA")
				CreateSupplierViaFunc(t, r, adminToken, "CFD", "222722727", "b@test.com", "CoA")
				var sid int64
				_ = config.DB.QueryRow("select supplier_id from suppliers where company = ? limit 1", "CoA").Scan(&sid)
				CreateUserViaFunc(t, r, adminToken, "SupUser", supEmail, supPass, "supplier_admin", sid)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			if tc.prepare != nil {
				tc.prepare(adminToken)
			}

			req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)

			switch tc.name {
			case "unauthorized (no token)":
			case "supplier_admin sees all suppliers of same company":
				supTok := LoginAndGetToken(t, supEmail, supPass)
				if supTok == "" {
					t.Fatalf("supplier login failed for subtest %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+supTok)
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectedCode == http.StatusOK && tc.expectBody != "" && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("%s expected body to contain %q; got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestGetSupplierByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_getbyid@test.com", "Admin@123"
	const supPass = "Supplier@123"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string) int64
		expectedCode int
		expectBody   string
	}{
		{
			name: "invalid id param",
			prepare: func(adminToken string) int64 {
				return 0
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "unauthorized (no token)",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "UnauthSup", "11128282822828", "unauth@test.com", "CoX")
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "not found (system_admin)",
			prepare: func(adminToken string) int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "system_admin found",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "SysFound", "22228282828", "sysfound@test.com", "CoSys")
			},
			expectedCode: http.StatusOK,
			expectBody:   "SysFound",
		},
		{
			name: "supplier_admin can get supplier of the same company",
			prepare: func(adminToken string) int64 {
				id := CreateSupplierViaFunc(t, r, adminToken, "CoASupA", "333181818", "coasup1@test.com", "CoA")
				CreateSupplierViaFunc(t, r, adminToken, "CoASup", "4448128282822", "coasup2@test.com", "CoB")
				CreateUserViaFunc(t, r, adminToken, "CoAAdmin", "coasupadmin@test.com", supPass, "supplier_admin", id)
				return id
			},
			expectedCode: http.StatusOK,
			expectBody:   "CoASupA",
		},
		{
			name: "supplier_admin forbidden for different company",
			prepare: func(adminToken string) int64 {
				idA := CreateSupplierViaFunc(t, r, adminToken, "CompanyX", "5552828", "compx@test.com", "CompanyX")
				idB := CreateSupplierViaFunc(t, r, adminToken, "CompanyY", "666282288", "compy@test.com", "CompanyY")
				CreateUserViaFunc(t, r, adminToken, "CompXAdmin", "compxadmin@test.com", "Supplier@321", "supplier_admin", idA)
				return idB
			},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			id := tc.prepare(adminToken)

			var url string
			switch tc.name {
			case "invalid id param":
				url = "/suppliers/abc"
			default:
				if id == 0 {
					url = fmt.Sprintf("/suppliers/%d", 99999999)
				} else {
					url = fmt.Sprintf("/suppliers/%d", id)
				}
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)

			switch tc.name {
			case "unauthorized (no token)":
			case "supplier_admin can get supplier of the same company":
				token := LoginAndGetToken(t, "coasupadmin@test.com", supPass)
				if token == "" {
					t.Fatalf("failed to login supplier admin for subtest %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+token)
			case "supplier_admin forbidden for different company":
				token := LoginAndGetToken(t, "compxadmin@test.com", "Supplier@321")
				if token == "" {
					t.Fatalf("failed to login supplier admin for subtest %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+token)
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectedCode == http.StatusOK && tc.expectBody != "" && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("%s expected body to contain %q; got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestCreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_prod_create@test.com", "Admin@123"
	const supPass = "Sup@2411"
	const supEmail = "supprod@test.com"
	r := routes.SetupRouter()

	cases := []struct {
		name         string
		prepare      func(adminToken string) (token string, cid int64)
		body         string
		expectedCode int
	}{
		{
			name: "unauthorized no token",
			prepare: func(adminToken string) (string, int64) {
				cid := createCategoryViaFunc(t, r, adminToken, "CatUnauth", "Desc")
				return "", cid
			},
			body:         `{"product_name":"Product1","product_description":"dajjajs","product_cost":100,"product_category_id":1}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for system_admin",
			prepare: func(adminToken string) (string, int64) {
				cid := createCategoryViaFunc(t, r, adminToken, "CatForbidden", "Desc")
				return adminToken, cid
			},
			body:         `{"product_name":"P1","product_description":"dajdjssjjajs","product_cost":100,"product_category_id":1}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid payload JSON",
			prepare: func(adminToken string) (string, int64) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupP", "9929292929", supEmail, "CoP")
				CreateUserViaFunc(t, r, adminToken, "SuppUser", supEmail, supPass, "supplier_admin", sid)
				cid := createCategoryViaFunc(t, r, adminToken, "CatP", "desc")
				token := LoginAndGetToken(t, supEmail, supPass)
				return token, cid
			},
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid product name (empty)",
			prepare: func(adminToken string) (string, int64) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupPf", "9992929292", "supprod2@test.com", "CoP")
				CreateUserViaFunc(t, r, adminToken, "SuppUser2", "supprod2@test.com", supPass, "supplier_admin", sid)
				cid := createCategoryViaFunc(t, r, adminToken, "CatP2", "desc")
				token := LoginAndGetToken(t, "supprod2@test.com", supPass)
				return token, cid
			},
			body:         `{"product_name":"   ","product_description":"dajjajs","product_cost":10,"product_category_id":1}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid category id",
			prepare: func(adminToken string) (string, int64) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupPdjdj", "9929292999", "supprod3@test.com", "CoP")
				CreateUserViaFunc(t, r, adminToken, "SuppUser3", "supprod3@test.com", supPass, "supplier_admin", sid)
				token := LoginAndGetToken(t, "supprod3@test.com", supPass)
				return token, 999999
			},
			body:         `{"product_name":"Pvalid","product_description":"dajjajs","product_cost":10,"product_category_id":999999}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "valid create by supplier_admin",
			prepare: func(adminToken string) (string, int64) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SupPghg", "999292929", "supprod4@test.com", "CoP")
				CreateUserViaFunc(t, r, adminToken, "SuppUser4", "supprod4@test.com", supPass, "supplier_admin", sid)
				token := LoginAndGetToken(t, "supprod4@test.com", supPass)
				cid := createCategoryViaFunc(t, r, adminToken, "CatP4", "desc")
				return token, cid
			},
			body:         `{"product_name":"Pvalid","product_description":"dajjajs","product_cost":10,"product_category_id":1}`,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			if adminToken == "" {
				t.Fatalf("admin login failed for setup")
			}

			token, cid := "", int64(0)
			if tc.prepare != nil {
				token, cid = tc.prepare(adminToken)
			}

			body := tc.body
			if cid > 0 {
				body = strings.Replace(body, `"product_category_id":1`, fmt.Sprintf(`"product_category_id":%d`, cid), 1)
				body = strings.Replace(body, `"product_category_id":999999`, fmt.Sprintf(`"product_category_id":%d`, cid), 1)
			}

			req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_prod_gets@test.com", "Admin@123"
	const supPass = "SupSee@1"
	const supEmail = "supsee@test.com"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string)
		expectedCode int
		expectBody   string
	}{
		{
			name: "unauthorized",
			prepare: func(adminToken string) {
				cid := createCategoryViaFunc(t, r, adminToken, "CatX", "dajjajs")
				sid := CreateSupplierViaFunc(t, r, adminToken, "Sx", "111882181", "sx_sup@test.com", "CoX")
				CreateUserViaFunc(t, r, adminToken, "SxAdmin", supEmail, supPass, "supplier_admin", sid)
				supToken := LoginAndGetToken(t, supEmail, supPass)
				_ = CreateProductViaFunc(t, r, supToken, "Px", "dajjajdjss", 10.0, cid)
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "system_admin sees all",
			prepare: func(adminToken string) {
				c := createCategoryViaFunc(t, r, adminToken, "Catdss", "dajjajs")
				s1 := CreateSupplierViaFunc(t, r, adminToken, "Supo", "112292991", "s1@s1.com", "CoA")
				CreateUserViaFunc(t, r, adminToken, "S1Admin", "s1@s1.com", "S1Pass@1", "supplier_admin", s1)
				s1Token := LoginAndGetToken(t, "s1@s1.com", "S1Pass@1")
				_ = CreateProductViaFunc(t, r, s1Token, "Prod1", "dajjajs", 11.5, c)
				s2 := CreateSupplierViaFunc(t, r, adminToken, "Srueu", "22282881182", "s2@s2.com", "CoB")
				CreateUserViaFunc(t, r, adminToken, "S2Admin", "s2@s2.com", "S2Pass@1", "supplier_admin", s2)
				s2Token := LoginAndGetToken(t, "s2@s2.com", "S2Pass@1")
				_ = CreateProductViaFunc(t, r, s2Token, "Prod2", "dajjajs", 22.5, c)
			},
			expectedCode: http.StatusOK,
			expectBody:   "Prod1",
		},
		{
			name: "supplier_admin sees same company products only",
			prepare: func(adminToken string) {
				c := createCategoryViaFunc(t, r, adminToken, "CatS", "dajjajs")

				s1 := CreateSupplierViaFunc(t, r, adminToken, "SCsinmsi", "1119292992", "sc1@sameco.com", "SameCo")
				CreateUserViaFunc(t, r, adminToken, "SC1Admin", "sc1@sameco.com", "SC1Pass@1", "supplier_admin", s1)
				sc1Token := LoginAndGetToken(t, "sc1@sameco.com", "SC1Pass@1")
				_ = CreateProductViaFunc(t, r, sc1Token, "SameCoProd", "dajjajs", 10.0, c)
				s2 := CreateSupplierViaFunc(t, r, adminToken, "SCsnsjdd", "2223929292", "sc2@sameco.com", "SameCo")
				CreateUserViaFunc(t, r, adminToken, "SC2Admin", "sc2@sameco.com", "SC2Pass@1", "supplier_admin", s2)
				sc2Token := LoginAndGetToken(t, "sc2@sameco.com", "SC2Pass@1")
				_ = CreateProductViaFunc(t, r, sc2Token, "SameCoProd2", "dajjajs", 12.0, c)
				s3 := CreateSupplierViaFunc(t, r, adminToken, "Other", "329291933", "other@otherco.com", "OtherCo")
				CreateUserViaFunc(t, r, adminToken, "OtherAdmin", "other@otherco.com", "OtherPass@1", "supplier_admin", s3)
				otherToken := LoginAndGetToken(t, "other@otherco.com", "OtherPass@1")
				if otherToken == "" {
					t.Fatalf("failed to login other supplier in setup")
				}
				_ = CreateProductViaFunc(t, r, otherToken, "OtherProd", "dajjajs", 20.0, c)
				CreateUserViaFunc(t, r, adminToken, "SupSee", supEmail, supPass, "supplier_admin", s1)
			},
			expectedCode: http.StatusOK,
			expectBody:   "SameCoProd",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			if tc.prepare != nil {
				tc.prepare(adminToken)
			}

			req := httptest.NewRequest(http.MethodGet, "/products", nil)

			switch tc.name {
			case "system_admin sees all":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			case "supplier_admin sees same company products only":
				supToken := LoginAndGetToken(t, supEmail, supPass)
				if supToken == "" {
					t.Fatalf("supplier login failed for subtest %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+supToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tc.name == "unauthorized" {
				if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
					t.Fatalf("%s expected unauthorized/forbidden got %d body=%s", tc.name, w.Code, w.Body.String())
				}
				return
			}

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectBody != "" && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("%s expected body to contain %q; got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestGetProductByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_prod_getbyid@test.com", "Admin@123"
	const supPass = "Supplier@123"
	const supEmail = "coasupadmin@test.com"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string) int64
		expectedCode int
		expectBody   string
	}{
		{
			name: "invalid id param",
			prepare: func(adminToken string) int64 {
				return 0
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "unauthorized",
			prepare: func(adminToken string) int64 {
				s := CreateSupplierViaFunc(t, r, adminToken, "Sx", "1183282821", "sx@test.com", "CoX")
				c := createCategoryViaFunc(t, r, adminToken, "CatX", "dajjajs")
				CreateUserViaFunc(t, r, adminToken, "SxAdmin", supEmail, supPass, "supplier_admin", s)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				p := CreateProductViaFunc(t, r, supTok, "Pxskakas", "dajjajs", 10.0, c)
				return p
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "not found (system_admin)",
			prepare: func(adminToken string) int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "system_admin can see any product by id",
			prepare: func(adminToken string) int64 {
				s1 := CreateSupplierViaFunc(t, r, adminToken, "Ssnnsns", "11292929211", "s1@tes.com", "CoA")
				CreateUserViaFunc(t, r, adminToken, "S1Admin", supEmail, supPass, "supplier_admin", s1)
				s1Tok := LoginAndGetToken(t, supEmail, supPass)
				c := createCategoryViaFunc(t, r, adminToken, "Cat1", "dajjajs")
				return CreateProductViaFunc(t, r, s1Tok, "Prod1", "dajjajs", 11.5, c)
			},
			expectedCode: http.StatusOK,
			expectBody:   "Prod1",
		},
		{
			name: "supplier_admin can get supplier product from same company",
			prepare: func(adminToken string) int64 {
				id1 := CreateSupplierViaFunc(t, r, adminToken, "CoASupA", "32919133", "coasup1@test.com", "CoA")
				id2 := CreateSupplierViaFunc(t, r, adminToken, "CoASupB", "44191929294", "coasup2@test.com", "CoA")
				c := createCategoryViaFunc(t, r, adminToken, "Cat1", "desc")
				CreateUserViaFunc(t, r, adminToken, "CoASup1Admin", supEmail, supPass, "supplier_admin", id1)
				tok1 := LoginAndGetToken(t, supEmail, supPass)
				_ = CreateProductViaFunc(t, r, tok1, "Prod1", "dish", 11.5, c)

				CreateUserViaFunc(t, r, adminToken, "CoASupAdmin", "hua@hua.com", supPass, "supplier_admin", id2)
				tok2 := LoginAndGetToken(t, "hua@hua.com", supPass)
				pid2 := CreateProductViaFunc(t, r, tok2, "Prod2", "dishtv", 11.5, c)
				CreateUserViaFunc(t, r, adminToken, "CoAAdmin", "hua2@hua.com", supPass, "supplier_admin", id1)
				return pid2
			},
			expectedCode: http.StatusOK,
			expectBody:   "Prod2",
		},
		{
			name: "supplier_admin forbidden for different company",
			prepare: func(adminToken string) int64 {
				idA := CreateSupplierViaFunc(t, r, adminToken, "CompanyX", "5552828", "compx@test.com", "CompanyX")
				idB := CreateSupplierViaFunc(t, r, adminToken, "CompanyY", "666282288", "compy@test.com", "CompanyY")
				c := createCategoryViaFunc(t, r, adminToken, "CatDiff", "desc")
				CreateUserViaFunc(t, r, adminToken, "CompXAdminUser", "compxadmin@test.com", "Supplier@321", "supplier_admin", idA)
				tokA := LoginAndGetToken(t, "compxadmin@test.com", "Supplier@321")
				if tokA == "" {
					t.Fatalf("setup: failed to login compxadmin@test.com")
				}
				_ = CreateProductViaFunc(t, r, tokA, "ProdA", "dish", 11.5, c)

				CreateUserViaFunc(t, r, adminToken, "CompYAdminUser", "compyadmin@test.com", "Supplier@321", "supplier_admin", idB)
				tokB := LoginAndGetToken(t, "compyadmin@test.com", "Supplier@321")
				if tokB == "" {
					t.Fatalf("setup: failed to login compyadmin@test.com")
				}
				pidB := CreateProductViaFunc(t, r, tokB, "ProdB", "dishtv", 11.5, c)
				CreateUserViaFunc(t, r, adminToken, "CompXAdmin", "compxadmin1@test.com", "Supplier@321", "supplier_admin", idA)
				return pidB
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)

			var pID int64
			if tc.prepare != nil {
				pID = tc.prepare(adminToken)
			}

			var url string
			if tc.name == "invalid id param" {
				url = "/products/abc"
			} else if pID == 0 && tc.expectedCode != http.StatusNotFound {
				url = fmt.Sprintf("/products/%d", 99999999)
			} else if pID == 0 {
				url = fmt.Sprintf("/products/%d", 99999999)
			} else {
				url = fmt.Sprintf("/products/%d", pID)
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)

			switch tc.name {
			case "unauthorized":
			case "system_admin can see any product by id", "not found (system_admin)", "invalid id param":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			case "supplier_admin can get supplier product from same company":
				token := LoginAndGetToken(t, supEmail, supPass)
				if token == "" {
					t.Fatalf("failed to login supplier admin for subtest %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+token)
			case "supplier_admin forbidden for different company":
				token := LoginAndGetToken(t, "compxadmin@test.com", "Supplier@321")
				if token == "" {
					t.Fatalf("failed to login supplier admin for subtest %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+token)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tc.name == "unauthorized" {
				if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
					t.Fatalf("%s expected unauthorized/forbidden got %d body=%s", tc.name, w.Code, w.Body.String())
				}
				return
			}

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectBody != "" && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("%s expected body to contain %q; got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_prod_update1@test.com", "Admin@123"
	const supPass = "Supplier@123"
	const supEmail = "Sup@sup.com"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func() int64
		body         string
		expectedCode int
		expectBody   string
	}{
		{
			name: "invalid id param",
			prepare: func() int64 {
				sid := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "UpdSup", "1182282821", "upd_sup@test.com", "UpCo")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "SupUpd", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				if supTok == "" {
					CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "SupUpd", supEmail, supPass, "supplier_admin", sid)
					supTok = LoginAndGetToken(t, supEmail, supPass)
				}
				cid := createCategoryViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatUPD", "dajjajs")
				pid := CreateProductViaFunc(t, r, supTok, "PUpd", "dajjajs", 9.9, cid)
				return pid
			},
			body:         `{"product_name":"Xdfifn"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "unauthorized (no token)",
			prepare: func() int64 {
				sid := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "UnauthSup", "11281818181", "unauthsup@test.com", "XCo")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "UnauthUser", "unauthsup@test.com", supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, "unauthsup@test.com", supPass)
				cid := createCategoryViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatUA", "dcsdd")
				return CreateProductViaFunc(t, r, supTok, "UnauthP", "ddadd", 5.5, cid)
			},
			body:         `{}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for system_admin",
			prepare: func() int64 {
				sid := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "SysOwner", "1128181811", "sysowner@test.com", "SysCo")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "SysSupp", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				cid := createCategoryViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatSys", "dfwsfcw")
				return CreateProductViaFunc(t, r, supTok, "SysProd", "dssd", 12.5, cid)
			},
			body:         `{"product_name":"NewName"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "no fields provided",
			prepare: func() int64 {
				sid := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "OwnerA", "11818228211", "ownera@test.com", "Aco")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "OwnerAUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				cid := createCategoryViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatNF", "dcscsd")
				return CreateProductViaFunc(t, r, supTok, "NoFieldP", "dsdcs", 6.6, cid)
			},
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "cannot change product supplier",
			prepare: func() int64 {
				sid := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "OwnerB", "11218281281", "ownerb@test.com", "BCo")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "OwnerBUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				cid := createCategoryViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatCS", "dsc")
				return CreateProductViaFunc(t, r, supTok, "ChangeSuppP", "descd", 7.7, cid)
			},
			body:         `{"product_supplier_id": 999}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "product not found",
			prepare: func() int64 {
				sid := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "NoProdSup", "121818218211", "nops@test.com", "NPS")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "NoProdUser", supEmail, supPass, "supplier_admin", sid)
				return 99999999
			},
			body:         `{"product_name":"Xsdhjd"}`,
			expectedCode: http.StatusNotFound,
		},
		{
			name: "forbidden update other supplier's product",
			prepare: func() int64 {
				idA := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "Asijs", "11182818381", "a@test.com", "Aco")
				idB := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "Bsihsih", "22821828182", "b@test.com", "Bco")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "BAdmin", "b@test.com", supPass, "supplier_admin", idB)
				bTok := LoginAndGetToken(t, "b@test.com", supPass)
				cid := createCategoryViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatFO", "dajjajs")
				pidB := CreateProductViaFunc(t, r, bTok, "BProd", "dajjajs", 3.3, cid)
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "AAdmin", supEmail, supPass, "supplier_admin", idA)
				return pidB
			},
			body:         `{"product_name":"ILegal"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid category id in payload",
			prepare: func() int64 {
				idS := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatOwner", "182812812811", "catowner@test.com", "Cco")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatOwnerU", supEmail, supPass, "supplier_admin", idS)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				cid := createCategoryViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatValid", "dajjajs")
				return CreateProductViaFunc(t, r, supTok, "CatCheck", "dajjajs", 14.5, cid)
			},
			body:         `{"product_category_id": 9999999}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "successful update",
			prepare: func() int64 {
				sid := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "GoodOwner", "1282182181811", "goodowner@test.com", "Gco")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "GoodSupUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				cid := createCategoryViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatGood", "dajjajs")
				return CreateProductViaFunc(t, r, supTok, "GoodProd", "dajjajs", 20.0, cid)
			},
			body:         `{"product_name":"GoodProdUpdated","product_description":"new","product_cost":21.5}`,
			expectedCode: http.StatusOK,
			expectBody:   "product updated",
		},
		{
			name: "duplicate product_name",
			prepare: func() int64 {
				sid := CreateSupplierViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "GoodOwner", "11131232", "goodowner@test.com", "Gco")
				CreateUserViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "GoodSupUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				cid := createCategoryViaFunc(t, r, LoginAndGetToken(t, "admin_prod_update1@test.com", "Admin@123"), "CatGood", "dajjajs")
				return CreateProductViaFunc(t, r, supTok, "ExistingName", "dajjajs", 10.0, cid)

			},
			body:         `{"product_name":"ExistingName"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)

			var pID int64
			if tc.prepare != nil {
				pID = tc.prepare()
			}

			var url string
			if tc.name == "invalid id param" {
				url = "/products/abc"
			} else if pID == 0 {
				url = fmt.Sprintf("/products/%d", 99999999)
			} else {
				url = fmt.Sprintf("/products/%d", pID)
			}

			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "unauthorized (no token)":
			case "forbidden for system_admin":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			case "invalid id param":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			default:
				supTok := LoginAndGetToken(t, supEmail, supPass)
				if supTok == "" {
					req.Header.Set("Authorization", "Bearer "+adminToken)
				} else {
					req.Header.Set("Authorization", "Bearer "+supTok)
				}
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tc.name == "unauthorized (no token)" {
				if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
					t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
				}
				return
			}
			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectBody != "" && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("%s expected body to contain %q got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_prod_delete@test.com", "Admin@123"
	r := routes.SetupRouter()
	const supPass = "Supplier@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)

	testcases := []struct {
		name         string
		prepare      func() (int64, string)
		expectedCode int
	}{
		{
			name: "invalid id param",
			prepare: func() (int64, string) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "InvSup", "118181282181", "invsup@test.com", "InvCo")
				CreateUserViaFunc(t, r, adminToken, "InvSupUser", "invsupuser@test.com", supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, "invsupuser@test.com", supPass)
				cid := createCategoryViaFunc(t, r, adminToken, "CatInv", "inv")
				_ = CreateProductViaFunc(t, r, supTok, "InvProd", "inv", 1.1, cid)
				return 0, "invsupuser@test.com"
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "unauthorized (no token)",
			prepare: func() (int64, string) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "UnauthSup", "111818228281", "unauthsup@test.com", "XCo")
				CreateUserViaFunc(t, r, adminToken, "UnauthUser", "unauthsup@test.com", supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, "unauthsup@test.com", supPass)
				cid := createCategoryViaFunc(t, r, adminToken, "CatUA", "dcsdd")
				p := CreateProductViaFunc(t, r, supTok, "UnauthP", "ddadd", 5.5, cid)
				return p, ""
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for system_admin",
			prepare: func() (int64, string) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SysDelSup", "111281283238", "sysdelsup@test.com", "SysCo")
				CreateUserViaFunc(t, r, adminToken, "SysDelUser", "sysdeluser@test.com", supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, "sysdeluser@test.com", supPass)
				cid := createCategoryViaFunc(t, r, adminToken, "CatSysDel", "dsjaasj")
				p := CreateProductViaFunc(t, r, supTok, "SysDelProd", "dcsjndsind", 2.2, cid)
				return p, "sysdeluser@test.com"
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "product not found (supplier tries to delete nonexistent id)",
			prepare: func() (int64, string) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "NoProdSup", "11281823831", "nops@test.com", "NPS")
				CreateUserViaFunc(t, r, adminToken, "NoProdUser", "nopsuser@test.com", supPass, "supplier_admin", sid)
				return 99999999, "nopsuser@test.com"
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "forbidden delete other supplier's product",
			prepare: func() (int64, string) {

				idA := CreateSupplierViaFunc(t, r, adminToken, "ASup", "139292911", "asup@test.com", "Aco")
				idB := CreateSupplierViaFunc(t, r, adminToken, "BSup", "23913942922", "bsup@test.com", "Bco")

				CreateUserViaFunc(t, r, adminToken, "BAdmin", "badmin@test.com", supPass, "supplier_admin", idB)
				bTok := LoginAndGetToken(t, "badmin@test.com", supPass)
				if bTok == "" {
					t.Fatalf("setup: failed to login badmin@test.com")
				}
				cid := createCategoryViaFunc(t, r, adminToken, "CatF", "dsanain")
				pid := CreateProductViaFunc(t, r, bTok, "BProd", "ddjid", 3.3, cid)

				CreateUserViaFunc(t, r, adminToken, "AAdmin", "aadmin@test.com", supPass, "supplier_admin", idA)
				return pid, "aadmin@test.com"
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "successful delete",
			prepare: func() (int64, string) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "GoodDelSup", "129139331911", "gooddelsup@test.com", "Gco")
				CreateUserViaFunc(t, r, adminToken, "GoodDelUser", "gooddeluser@test.com", supPass, "supplier_admin", sid)
				gTok := LoginAndGetToken(t, "gooddeluser@test.com", supPass)
				if gTok == "" {
					t.Fatalf("setup: failed to login gooddeluser@test.com")
				}
				cid := createCategoryViaFunc(t, r, adminToken, "CatDel", "ddand")
				pid := CreateProductViaFunc(t, r, gTok, "GoodDeleteProd", "ddunsqun", 10.0, cid)
				return pid, "gooddeluser@test.com"
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })

			var pID int64
			var supplierEmail string
			if tc.prepare != nil {
				pID, supplierEmail = tc.prepare()
			}

			var url string
			if tc.name == "invalid id param" {
				url = "/products/abc"
			} else if pID == 0 {
				url = fmt.Sprintf("/products/%d", 99999999)
			} else {
				url = fmt.Sprintf("/products/%d", pID)
			}

			req := httptest.NewRequest(http.MethodDelete, url, nil)

			switch tc.name {
			case "unauthorized (no token)":
			case "forbidden for system_admin":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			case "invalid id param":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			default:
				if supplierEmail == "" {
					t.Fatalf("test %s expected a supplier user in prepare, but got none", tc.name)
				}
				supTok := LoginAndGetToken(t, supplierEmail, supPass)
				if supTok == "" {
					t.Fatalf("failed to login supplier %s for subtest %s", supplierEmail, tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+supTok)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tc.name == "unauthorized (no token)" {
				if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
					t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
				}
				return
			}

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestCreateSuppAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "sup_admin_create2@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	const supPass = "Supplier@123"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string) int64
		body         string
		expectedCode int
		expectBody   string
	}{
		{
			name: "unauthorized (no token)",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SNoTok", "112811821", "snotok@test.com", "Co")
				return sid
			},
			body:         `{"name":"Ucsic","email":"u@test.com","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "SF", "112338811", "sf@test.com", "CoF")
				CreateUserViaFunc(t, r, adminToken, "SFUser", "sfuser@test.com", supPass, "supplier_admin", sid)
				return sid
			},
			body:         `{"name":"Usicis","email":"u2@test.com","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid json",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Sij", "113292941", "sij@test.com", "Co")
			},
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "missing fields",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Smiss", "1111939391", "smiss@test.com", "Co")
			},
			body:         `{"name":"","email":"", "password":"", "supplier_id":0}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Sinv", "1191291291", "sinv@test.com", "Co")
			},
			body:         `{"name": "Usjsj9oj","email":"notemail","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid password",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Spwd", "1219291311", "spwd@test.com", "Co")
			},
			body:         `{"name":"Udifif","email":"usidjidhn@test.com","password":"not","supplier_id":1}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "successful create",
			prepare: func(adminToken string) int64 {
				return CreateSupplierViaFunc(t, r, adminToken, "Sgood", "1129321931", "sgood@test.com", "GCo")
			},
			body:         `{"name":"TesterAdmin","email":"newadmin@test.com","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusCreated,
			expectBody:   "supplier_admin user created successfully",
		},
		{
			name: "email already registered",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "Sdup", "111291939", "sdup@test.com", "Co")
				CreateUserViaFunc(t, r, adminToken, "Existing", "existing@test.com", supPass, "supplier_admin", sid)
				return sid
			},
			body:         `{"name":"usbdusbd","email":"existing@test.com","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "supplier_id does not exist",
			prepare: func(adminToken string) int64 {
				return 0
			},
			body:         `{"name":"Ucsdd","email":"u4@test.com","password":"Supplier@123","supplier_id":999999}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })

			var sid int64
			if tc.prepare != nil {
				sid = tc.prepare(adminToken)
			}

			body := tc.body
			if sid > 0 {
				body = strings.Replace(body, `"supplier_id":1`, fmt.Sprintf(`"supplier_id":%d`, sid), 1)
			}

			req := httptest.NewRequest(http.MethodPost, "/users/supplier-admin", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")
			case "forbidden for supplier_admin":
				supTok := LoginAndGetToken(t, "sfuser@test.com", supPass)
				if supTok != "" {
					req.Header.Set("Authorization", "Bearer "+supTok)
				}
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectBody != "" && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("%s expected body to contain %q; got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestGetSuppAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "sup_admin_get@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	const supPass = "Supplier@123"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string)
		expectedCode int
	}{
		{
			name: "unauthorized (no token)",
			prepare: func(adminToken string) {

				sid := CreateSupplierViaFunc(t, r, adminToken, "Slusins", "1119291291", "slusdinn@test.com", "Co")
				CreateUserViaFunc(t, r, adminToken, "SLUser", "sluser@test.com", supPass, "supplier_admin", sid)
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func(adminToken string) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "Sup", "112199139391", "sup@test.com", "Co")
				CreateUserViaFunc(t, r, adminToken, "SupUser", "supuser@test.com", supPass, "supplier_admin", sid)
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "system_admin authorized",
			prepare: func(adminToken string) {
				sid := CreateSupplierViaFunc(t, r, adminToken, "Sdajdj", "121939394911", "s1@test.com", "CoA")
				CreateUserViaFunc(t, r, adminToken, "U1", "u1@test.com", supPass, "supplier_admin", sid)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })
			if tc.prepare != nil {
				tc.prepare(adminToken)
			}

			req := httptest.NewRequest(http.MethodGet, "/users/supplier-admin", nil)
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")
			case "forbidden for supplier_admin":
				supTok := LoginAndGetToken(t, "supuser@test.com", supPass)
				if supTok != "" {
					req.Header.Set("Authorization", "Bearer "+supTok)
				}
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetSuppAdminByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "sup_admin_get2@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	const supPass = "Supplier@123"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string) int64
		expectedCode int
	}{
		{
			name: "invalid id param",
			prepare: func(adminToken string) int64 {
				return 0
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "unauthorized (no token)",
			prepare: func(adminToken string) int64 {

				sid := CreateSupplierViaFunc(t, r, adminToken, "NoTokSup", "112329391", "notoksup@test.com", "Co")
				CreateUserViaFunc(t, r, adminToken, "NoTokUser", "notok@test.com", supPass, "supplier_admin", sid)
				return 0
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "not found",
			prepare: func(adminToken string) int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "found",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "UsrS", "1139239491", "usrs@test.com", "Co")
				CreateUserViaFunc(t, r, adminToken, "FoundU", "found@test.com", supPass, "supplier_admin", sid)
				var uid int64
				_ = config.DB.QueryRow("select user_id from users where email = ? limit 1", "found@test.com").Scan(&uid)
				return uid
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "Sup", "112913931", "sup@test.com", "Co")
				CreateUserViaFunc(t, r, adminToken, "SupUser", "supuser@test.com", supPass, "supplier_admin", sid)
				var uid int64
				_ = config.DB.QueryRow("select user_id from users where email = ? limit 1", "supuser@test.com").Scan(&uid)
				return uid
			},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })

			var uid int64
			if tc.prepare != nil {
				uid = tc.prepare(adminToken)
			}

			var url string
			if tc.name == "invalid id param" {
				url = "/users/supplier-admin/abc"
			} else if uid == 0 {
				url = fmt.Sprintf("/users/supplier-admin/%d", 99999999)
			} else {
				url = fmt.Sprintf("/users/supplier-admin/%d", uid)
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")
			case "invalid id param", "not found", "found":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			case "forbidden for supplier_admin":
				supTok := LoginAndGetToken(t, "supuser@test.com", supPass)
				if supTok != "" {
					req.Header.Set("Authorization", "Bearer "+supTok)
				}
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestDeleteSuppAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "sup_admin_get@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	const supPass = "Supplier@123"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func(adminToken string) int64
		expectedCode int
	}{
		{
			name: "unauthorized (no token)",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "NoTokSup", "1291993911", "notok@test.com", "Co")
				CreateUserViaFunc(t, r, adminToken, "NoTokUser", "notok@test.com", supPass, "supplier_admin", sid)
				var uid int64
				_ = config.DB.QueryRow("select user_id from users where email = ? limit 1", "notok@test.com").Scan(&uid)
				return uid
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid id param",
			prepare: func(adminToken string) int64 {
				return 0
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "Sup", "111291393040", "sup@test.com", "Co")
				CreateUserViaFunc(t, r, adminToken, "SupUser", "supuser@test.com", supPass, "supplier_admin", sid)
				var uid int64
				_ = config.DB.QueryRow("select user_id from users where email = ? limit 1", "supuser@test.com").Scan(&uid)
				return uid
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "not found",
			prepare: func(adminToken string) int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful delete",
			prepare: func(adminToken string) int64 {
				sid := CreateSupplierViaFunc(t, r, adminToken, "Sup", "121373711", "sup@test.com", "Co")
				CreateUserViaFunc(t, r, adminToken, "SupUser", "supuser@test.com", supPass, "supplier_admin", sid)
				var uid int64
				_ = config.DB.QueryRow("select user_id from users where email = ? limit 1", "supuser@test.com").Scan(&uid)
				return uid
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })

			var uid int64
			if tc.prepare != nil {
				uid = tc.prepare(adminToken)
			}

			var url string
			if tc.name == "invalid id param" {
				url = "/users/supplier-admin/abc"
			} else if uid == 0 {
				url = fmt.Sprintf("/users/supplier-admin/%d", 99999999)
			} else {
				url = fmt.Sprintf("/users/supplier-admin/%d", uid)
			}

			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")
			case "invalid id param", "not found", "successful delete":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			case "forbidden for supplier_admin":
				supTok := LoginAndGetToken(t, "supuser@test.com", supPass)
				if supTok != "" {
					req.Header.Set("Authorization", "Bearer "+supTok)
				}
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestCreateStockMov(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "sup_admin_get@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	const supPass = "Supplier@123"
	const supEmail = "Supplier@testing.com"
	var fkSupTok string
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func()
		body         string
		expectedCode int
	}{
		{
			name:         "unauthorized (no token)",
			prepare:      func() {},
			body:         `{}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid json",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "MTSup", "19139911", "mtsup@test.com", "MtCo")
				cid := createCategoryViaFunc(t, r, adminToken, "MTCat", "dcdunu")
				CreateUserViaFunc(t, r, adminToken, "MTSupUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				if supTok == "" {
					t.Fatalf("setup: failed to login %s", supEmail)
				}
				_ = CreateProductViaFunc(t, r, supTok, "MTProd", "ddfuhfu", 11.0, cid)
			},
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden for non-supplier (system_admin)",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "FbdSup", "101303011", "fbdsup@test.com", "CoF")
				cid := createCategoryViaFunc(t, r, adminToken, "FbdCat", "dfund")
				supTok := ""
				CreateUserViaFunc(t, r, adminToken, "FbdUser", "fbduser@test.com", supPass, "supplier_admin", sid)
				supTok = LoginAndGetToken(t, "fbduser@test.com", supPass)
				_ = CreateProductViaFunc(t, r, supTok, "FbdProd", "dufucnsn", 10.0, cid)
			},
			body:         `{"product_id":1,"quantity":5,"movement_type":"IN"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid movement type",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "MTSup", "130230411", "mtsup@test.com", "MtCo")
				cid := createCategoryViaFunc(t, r, adminToken, "MTCat", "dcdunu")
				CreateUserViaFunc(t, r, adminToken, "MTSupUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				_ = CreateProductViaFunc(t, r, supTok, "MTProd", "ddfuhfu", 11.0, cid)
			},
			body:         `{"product_id":1,"quantity":5,"movement_type":"INDHDHDH"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid quantity (<=0)",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "MTSup", "132904011", "mtsup@test.com", "MtCo")
				cid := createCategoryViaFunc(t, r, adminToken, "MTCat", "dcdunu")
				CreateUserViaFunc(t, r, adminToken, "MTSupUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				if supTok == "" {
					t.Fatalf("setup: failed to login %s", supEmail)
				}
				_ = CreateProductViaFunc(t, r, supTok, "MTProd", "ddfuhfu", 11.0, cid)
			},
			body:         `{"product_id":1,"quantity":-1,"movement_type":"IN"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "product not found",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "MTSup", "11202031", "mtsup@test.com", "MtCo")
				CreateUserViaFunc(t, r, adminToken, "MTSupUser", supEmail, supPass, "supplier_admin", sid)
			},
			body:         `{"product_id":999999,"quantity":25,"movement_type":"IN"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden supplier (cannot create stock movements for other products)",
			prepare: func() {
				s1 := CreateSupplierViaFunc(t, r, adminToken, "MTSup", "120402011", "mtsup@test.com", "MtCo")
				s2 := CreateSupplierViaFunc(t, r, adminToken, "MTSupdjsjdj", "22298492", "mtsup2@test.com", "MtCodsjsjs")
				cid := createCategoryViaFunc(t, r, adminToken, "MTCat", "dcdunu")
				CreateUserViaFunc(t, r, adminToken, "S1User", "s1@test.com", supPass, "supplier_admin", s1)
				tok1 := LoginAndGetToken(t, "s1@test.com", supPass)
				if tok1 == "" {
					t.Fatalf("setup: failed to login s1@test.com")
				}
				_ = CreateProductViaFunc(t, r, tok1, "MTProd", "ddfuhfu", 11.0, cid)
				CreateUserViaFunc(t, r, adminToken, "TSupdjsjdj", supEmail, supPass, "supplier_admin", s2)
			},
			body:         `{"product_id":1,"quantity":25,"movement_type":"IN"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "OUT causes low stock warning (below threshold)",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "MTSup", "1123924941", "mtsup@test.com", "MtCo")
				cid := createCategoryViaFunc(t, r, adminToken, "MTCat", "dcdunu")
				CreateUserViaFunc(t, r, adminToken, "MTSupUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				pid := CreateProductViaFunc(t, r, supTok, "MTProd", "ddfuhfu", 11.0, cid)
				_ = pid
			},
			body:         `{"product_id":1,"quantity":5,"movement_type":"OUT"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "successful IN by supplier_admin",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "MTSup", "12949411", "mtsup@test.com", "MtCo")
				cid := createCategoryViaFunc(t, r, adminToken, "MTCat", "dcdunu")
				CreateUserViaFunc(t, r, adminToken, "MTSupUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				_ = CreateProductViaFunc(t, r, supTok, "MTProd", "ddfuhfu", 11.0, cid)
			},
			body:         `{"product_id":1,"quantity":25,"movement_type":"IN"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "successful OUT by supplier_admin",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "MTSup", "111393941", "mtsup@test.com", "MtCo")
				cid := createCategoryViaFunc(t, r, adminToken, "MTCat", "dcdunu")
				CreateUserViaFunc(t, r, adminToken, "MTSupUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				pid := CreateProductViaFunc(t, r, supTok, "MTProd", "ddfuhfu", 11.0, cid)
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid, 20, "IN", 1)
			},
			body:         `{"product_id":1,"quantity":5,"movement_type":"OUT","reason":"selling"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "insert FK error (performed_by missing)",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "FKSup", "12943924911", "fksup@test.com", "FKCo")
				cid := createCategoryViaFunc(t, r, adminToken, "FKCat", "fkdesc")
				CreateUserViaFunc(t, r, adminToken, "FKUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				if supTok == "" {
					t.Fatalf("setup: failed to login %s", supEmail)
				}
				fkSupTok = supTok

				pid := CreateProductViaFunc(t, r, supTok, "FKProd", "fkprod", 10.0, cid)
				_ = pid
				_, _ = config.DB.Exec("delete from users where email = ?", supEmail)
			},
			body:         `{"product_id":1,"quantity":50,"movement_type":"IN"}`,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })

			if tc.prepare != nil {
				tc.prepare()
			}
			req := httptest.NewRequest(http.MethodPost, "/stock_movements", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")
			case "forbidden for non-supplier (system_admin)":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			case "invalid json", "invalid movement type", "invalid quantity (<=0)", "product not found", "forbidden supplier (cannot create stock movements for other products)", "OUT causes low stock warning (below threshold)", "successful IN by supplier_admin", "successful OUT by supplier_admin":
				supToken := LoginAndGetToken(t, supEmail, supPass)
				if supToken != "" {
					req.Header.Set("Authorization", "Bearer "+supToken)
				}
			case "insert FK error (performed_by missing)":
				if fkSupTok != "" {
					req.Header.Set("Authorization", "Bearer "+fkSupTok)
				} else {
					supTok := LoginAndGetToken(t, supEmail, supPass)
					req.Header.Set("Authorization", "Bearer "+supTok)
				}
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetStockMov(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "sup_admin_get@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	const supPass = "Supplier@123"
	const supEmail = "Supplier@testing.com"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func()
		query        string
		expectedCode int
	}{
		{
			name: "invalid product_id param",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "tester", "1113939111", "tester@test.com", "TestCo")
				cid := createCategoryViaFunc(t, r, adminToken, "CatSys", "desc")
				CreateUserViaFunc(t, r, adminToken, "SysProdUser", "sysprod@test.com", supPass, "supplier_admin", sid)

				supTok := LoginAndGetToken(t, "sysprod@test.com", supPass)
				if supTok == "" {
					t.Fatalf("setup: failed to login sysprod@test.com")
				}
				pid := CreateProductViaFunc(t, r, supTok, "SysProd", "descj", 10.0, cid)
				if pid != 0 {
					_ = CreateStockViaFunc(t, r, supTok, pid, 50, "IN", "")
					_ = CreateStockViaFunc(t, r, supTok, pid, 2, "OUT", "")
				}
			},
			query:        "?product_id=abc",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "unauthorized (no token)",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "unauthSup", "111139399", "unauthsup@test.com", "UCo")
				cid := createCategoryViaFunc(t, r, adminToken, "UCat", "desc")
				CreateUserViaFunc(t, r, adminToken, "UnauthUser", "unauthuser@test.com", supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, "unauthuser@test.com", supPass)
				if supTok == "" {
					t.Fatalf("setup: failed to login unauthuser@test.com")
				}
				pid := CreateProductViaFunc(t, r, supTok, "UnauthProd", "desc", 10.0, cid)
				if pid != 0 {
					_ = CreateStockViaFunc(t, r, supTok, pid, 10, "IN", "")
				}
			},
			query:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "system_admin sees all movements",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "tester", "1123932912", "tester2@test.com", "TestCo")
				cid := createCategoryViaFunc(t, r, adminToken, "CatSys", "desc")
				CreateUserViaFunc(t, r, adminToken, "SysProdUser", "sysprod@test.com", supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, "sysprod@test.com", supPass)
				if supTok == "" {
					t.Fatalf("setup: failed to login sysprod@test.com")
				}
				pid := CreateProductViaFunc(t, r, supTok, "SysProd", "descj", 10.0, cid)
				if pid != 0 {
					_ = CreateStockViaFunc(t, r, supTok, pid, 50, "IN", "")
					_ = CreateStockViaFunc(t, r, supTok, pid, 2, "OUT", "")
				}
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin filters by product_id",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "tester", "1129323913", "tester3@test.com", "TestCo")
				cid := createCategoryViaFunc(t, r, adminToken, "CatSys", "desc")
				CreateUserViaFunc(t, r, adminToken, "SysProdUser", "sysprod@test.com", supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, "sysprod@test.com", supPass)
				if supTok == "" {
					t.Fatalf("setup: failed to login sysprod@test.com")
				}
				pid1 := CreateProductViaFunc(t, r, supTok, "SysProd", "descj", 10.0, cid)
				_ = CreateProductViaFunc(t, r, supTok, "SysProdsasa", "descjas", 10.0, cid)
				if pid1 != 0 {
					_ = CreateStockViaFunc(t, r, supTok, pid1, 50, "IN", "")
				}
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin no movements returns empty list",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "HUhusup", "111183384", supEmail, "NmCo")
				cid := createCategoryViaFunc(t, r, adminToken, "NmCat", "nm")
				CreateUserViaFunc(t, r, adminToken, "NmUser", supEmail, supPass, "supplier_admin", sid)
				supTok := LoginAndGetToken(t, supEmail, supPass)
				if supTok == "" {
					t.Fatalf("setup: failed to login %s", supEmail)
				}
				_ = CreateProductViaFunc(t, r, supTok, "NmProd", "nmd", 10.0, cid)
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin sees company movements",
			prepare: func() {
				s1 := CreateSupplierViaFunc(t, r, adminToken, "SupCa", "111282835", "supc1@test.com", "SameCo")
				s2 := CreateSupplierViaFunc(t, r, adminToken, "SupCb", "228388225", "supc2@test.com", "SameCo")
				cid := createCategoryViaFunc(t, r, adminToken, "Ccat", "desc")
				CreateUserViaFunc(t, r, adminToken, "CompanyAdmin", supEmail, supPass, "supplier_admin", s1)
				tok1 := LoginAndGetToken(t, supEmail, supPass)
				if tok1 == "" {
					t.Fatalf("setup: failed to login %s", supEmail)
				}
				pid1 := CreateProductViaFunc(t, r, tok1, "Pssind", "dsdudun", 10.0, cid)
				tok2email := "supc2@test.com"
				CreateUserViaFunc(t, r, adminToken, "CompanyAdmin2", tok2email, supPass, "supplier_admin", s2)
				tok2 := LoginAndGetToken(t, tok2email, supPass)
				if tok2 == "" {
					t.Fatalf("setup: failed to login %s", tok2email)
				}
				pid2 := CreateProductViaFunc(t, r, tok2, "Psidjsidj", "dicnicn", 20.0, cid)

				if pid1 != 0 {
					_ = CreateStockViaFunc(t, r, tok1, pid1, 16, "IN", "")
				}
				if pid2 != 0 {
					_ = CreateStockViaFunc(t, r, tok2, pid2, 4, "IN", "")
				}
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin filter by product_id",
			prepare: func() {
				s1 := CreateSupplierViaFunc(t, r, adminToken, "SupCsd", "112939316", "supc1@test.com", "SameCo")
				cid := createCategoryViaFunc(t, r, adminToken, "Ccat", "desc")
				CreateUserViaFunc(t, r, adminToken, "coadmin", supEmail, supPass, "supplier_admin", s1)
				tok := LoginAndGetToken(t, supEmail, supPass)
				if tok == "" {
					t.Fatalf("setup: failed to login %s", supEmail)
				}
				pid1 := CreateProductViaFunc(t, r, tok, "Pssind", "dsdudun", 10.0, cid)
				if pid1 != 0 {
					_ = CreateStockViaFunc(t, r, tok, pid1, 7, "IN", "")
				}
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin no movements returns empty list",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "fakeSup", "11129319317", "fake@test.com", "fakeCo")
				CreateUserViaFunc(t, r, adminToken, "fakeAdmin", supEmail, supPass, "supplier_admin", sid)
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })

			if tc.prepare != nil {
				tc.prepare()
			}

			query := tc.query
			if tc.name == "system_admin filter by product_id" || tc.name == "supplier_admin filter by product_id" {
				var pid int64
				err := config.DB.QueryRow("select product_id from products order by product_id asc limit 1").Scan(&pid)
				if err == nil && pid != 0 {
					query = fmt.Sprintf("?product_id=%d", pid)
				}
			}
			url := "/stock_movements" + query
			req := httptest.NewRequest(http.MethodGet, url, nil)
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")
			case "system_admin sees all movements", "system_admin filters by product_id", "invalid product_id param", "system_admin no movements returns empty list":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			default:
				supToken := LoginAndGetToken(t, supEmail, supPass)
				if supToken != "" {
					req.Header.Set("Authorization", "Bearer "+supToken)
				}
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestDashboard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "sup_admin_get@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	const supPass = "Supplier@123"
	const supEmail = "Supplier@testing.com"
	r := routes.SetupRouter()

	testcases := []struct {
		name         string
		prepare      func()
		expectedCode int
	}{
		{
			name:         "unauthorized (no token)",
			prepare:      func() {},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "system_admin empty counts",
			prepare:      func() {},
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin dashboard (low stock products)",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "LowSup", "128482911", "lowsup@test.com", "LowCo")
				cid := createCategoryViaFunc(t, r, adminToken, "LowCat", "desc")
				CreateUserViaFunc(t, r, adminToken, "LowUser", "lowuser@test.com", supPass, "supplier_admin", sid)

				tok := LoginAndGetToken(t, "lowuser@test.com", supPass)
				pid := CreateProductViaFunc(t, r, tok, "LowProd", "desc", 10.0, cid)
				if pid != 0 {
					_ = CreateStockViaFunc(t, r, tok, pid, 500, "IN", "")
					_ = CreateStockViaFunc(t, r, tok, pid, 100, "OUT", "consume")
				}
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin no low stock products (all above threshold)",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "HighSup", "1128248491", "highsup@test.com", "HighCo")
				cid := createCategoryViaFunc(t, r, adminToken, "HighCat", "desc")
				CreateUserViaFunc(t, r, adminToken, "HighUser", "highuser@test.com", supPass, "supplier_admin", sid)

				tok := LoginAndGetToken(t, "highuser@test.com", supPass)
				if tok == "" {
					t.Fatalf("setup: failed to login %s", "highuser@test.com")
				}

				pid := CreateProductViaFunc(t, r, tok, "HighProd", "desc", 10.0, cid)
				if pid != 0 {
					_ = CreateStockViaFunc(t, r, tok, pid, 100, "IN", "")
				}
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin dashboard (has products)",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "LowSup", "11293010381", "lowsup@test.com", "LowCo")
				cid := createCategoryViaFunc(t, r, adminToken, "LowCat", "desc")
				CreateUserViaFunc(t, r, adminToken, "CompanyAdmin", supEmail, supPass, "supplier_admin", sid)
				tok := LoginAndGetToken(t, supEmail, supPass)
				if tok == "" {
					t.Fatalf("setup: failed to login %s", supEmail)
				}

				_ = CreateProductViaFunc(t, r, tok, "LowProd", "desc", 10.0, cid)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin dashboard (no products)",
			prepare: func() {
				sid := CreateSupplierViaFunc(t, r, adminToken, "LowSup", "111932091", "lowsup@test.com", "LowCo")
				CreateUserViaFunc(t, r, adminToken, "CompanyAdmin", supEmail, supPass, "supplier_admin", sid)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			t.Cleanup(func() { TruncateAll(t) })

			if tc.prepare != nil {
				tc.prepare()
			}

			req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
			req.Header.Set("Content-Type", "application/json")

			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")
			case "system_admin empty counts", "system_admin dashboard (low stock products)", "system_admin no low stock products (all above threshold)":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			case "supplier_admin dashboard (has products)", "supplier_admin dashboard (no products)":
				supToken := LoginAndGetToken(t, supEmail, supPass)
				if supToken != "" {
					req.Header.Set("Authorization", "Bearer "+supToken)
				}
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}
