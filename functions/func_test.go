package functions_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"product-catalogue/config"
	"product-catalogue/functions"
	"product-catalogue/middleware"
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

func ExecMigration(db *sql.DB, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	queries := strings.Split(string(content), ";")
	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}
	return nil
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

	r := gin.New()
	r.POST("/auth/login", functions.Login)

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
		log.Fatal("DSN not found in .env")
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

	if err := ExecMigration(db, mig); err != nil {
		log.Fatal(err)
	}

	code := m.Run()
	os.Exit(code)
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
			r := gin.New()
			r.POST("/auth/login", functions.Login)
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
	TruncateAll(t)
	adminEmail, adminPass := "admin_td@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup", "sup_td@test.com", "Comp")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"SupUser", "sup_td@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supplierToken := LoginAndGetToken(t, "sup_td@test.com", "Yuvraj@2411")

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.POST("/category", functions.CreateCategory)

	cases := []struct {
		name         string
		token        string
		body         string
		expectedCode int
	}{
		{
			name:         "no token",
			token:        "",
			body:         `{"category_name":"X","category_description":"Y"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "supplier_admin cannot create",
			token:        supplierToken,
			body:         `{"category_name":"X","category_description":"Y"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalid json",
			token:        adminToken,
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid category name",
			token:        adminToken,
			body:         `{"category_name":"@#$","category_description":"valid"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "valid create",
			token:        adminToken,
			body:         `{"category_name":"Electronics","category_description":"Devices"}`,
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)

			req := httptest.NewRequest(http.MethodPost, "/category", strings.NewReader(tc.body))
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectedCode == http.StatusCreated {
				var cnt int
				err := config.DB.QueryRow("select count(*) from categories where category_name = ?", "Electronics").Scan(&cnt)
				if err != nil {
					t.Fatalf("db query error: %v", err)
				}
				if cnt != 1 {
					t.Fatalf("expected 1 category created, got %d", cnt)
				}
			}
		})
	}
}

func TestGetCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	TruncateAll(t)

	adminEmail, adminPass := "admin_get@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup2", "sup2d@test.com", "Comp")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"Sup2User", "sup2d@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supplierToken := LoginAndGetToken(t, "sup2d@test.com", "Yuvraj@2411")

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.GET("/category", functions.GetCategory)

	cases := []struct {
		name         string
		token        string
		seedRows     int
		expectedCode int
		expectBody   string
	}{
		{
			name:         "no token",
			token:        "",
			seedRows:     0,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "system_admin authorized",
			token:        adminToken,
			seedRows:     2,
			expectedCode: http.StatusOK,
			expectBody:   "Electronics",
		},
		{
			name:         "supplier_admin authorized",
			token:        supplierToken,
			seedRows:     2,
			expectedCode: http.StatusOK,
			expectBody:   "Electronics",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			if tc.seedRows > 0 {
				_, _ = config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "Electronics", "Devices")
				if tc.seedRows > 1 {
					_, _ = config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "Home", "Appliances")
				}
			}

			req := httptest.NewRequest(http.MethodGet, "/category", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
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
	TruncateAll(t)

	adminEmail, adminPass := "admin_getid@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup3", "sup3d@test.com", "Comp")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"Sup3User", "sup3d@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supplierToken := LoginAndGetToken(t, "sup3d@test.com", "Yuvraj@2411")

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.GET("/category/:id", functions.GetCategoryByID)

	cases := []struct {
		name         string
		token        string
		prepare      func() int64
		url          string
		expectedCode int
		expectBody   string
	}{
		{
			name:         "invalid id param",
			token:        adminToken,
			prepare:      func() int64 { return 0 },
			url:          "/category/abc",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "not found",
			token: adminToken,
			prepare: func() int64 {

				return 0
			},
			url:          "/category/9999999",
			expectedCode: http.StatusNotFound,
		},
		{
			name:  "found",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "Clothing", "Wearables")
				id, _ := res.LastInsertId()
				return id
			},
			url:          "",
			expectedCode: http.StatusOK,
			expectBody:   "Clothing",
		},
		{
			name:  "system_admin allowed",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "DinnerSet", "Plates")
				id, _ := res.LastInsertId()
				return id
			},
			url:          "",
			expectedCode: http.StatusOK,
			expectBody:   "DinnerSet",
		},
		{
			name:  "supplier_admin allowed",
			token: supplierToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "DinnerSet", "Plates")
				id, _ := res.LastInsertId()
				return id
			},
			url:          "",
			expectedCode: http.StatusOK,
			expectBody:   "DinnerSet",
		},
		{

			name:  "no token",
			token: "",
			prepare: func() int64 {

				return 0
			},
			url:          "/category/1",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			var id = int64(0)
			if tc.prepare != nil {
				id = tc.prepare()
			}
			url := tc.url
			if url == "" && id != 0 {
				url = fmt.Sprintf("/category/%d", id)
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	TruncateAll(t)

	adminEmail, adminPass := "admin_update@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup3", "sup3d@test.com", "Comp")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"Sup3User", "sup3d@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supplierToken := LoginAndGetToken(t, "sup3d@test.com", "Yuvraj@2411")

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.PUT("/category/:id", functions.UpdateCategory)

	cases := []struct {
		name         string
		token        string
		prepare      func() int64
		body         string
		expectedCode int
		verify       func(id int64) error
	}{
		{
			name:  "no fields",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "Sports", "Items")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "forbidden for supplier_admin",
			token: supplierToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "Sports", "Items")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"category_name":"Movies"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name:  "invalid id bad request",
			token: adminToken,
			prepare: func() int64 {
				return 0
			},
			body:         `{"category_name":"NewName"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "valid update",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "Tools", "Hardware")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"category_name":"Hardware"}`,
			expectedCode: http.StatusOK,
			verify: func(id int64) error {
				var name string
				if err := config.DB.QueryRow("select category_name from categories where category_id=?", id).Scan(&name); err != nil {
					return err
				}
				if name != "Hardware" {
					return fmt.Errorf("expected name Hardware got %s", name)
				}
				return nil
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			id := tc.prepare()
			url := fmt.Sprintf("/category/%d", id)
			if id == 0 && tc.name == "invalid id" {
				url = "/category/abc"
			}
			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("[%s] expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.verify != nil {
				if err := tc.verify(id); err != nil {
					t.Fatalf("verification failed: %v", err)
				}
			}
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	TruncateAll(t)

	adminEmail, adminPass := "admin_delete@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)

	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SupDel", "sup_del@test.com", "C")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"SupDelUser", "sup_del@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supToken := LoginAndGetToken(t, "sup_del@test.com", "Yuvraj@2411")

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.DELETE("/category/:id", functions.DeleteCategory)

	cases := []struct {
		name         string
		token        string
		prepare      func() int64
		expectedCode int
	}{
		{
			name:  "invalid id",
			token: adminToken,
			prepare: func() int64 {
				return 0
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "forbidden for supplier_admin",
			token: supToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "ToDelete", "X")
				id, _ := res.LastInsertId()
				return id
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name:  "delete not found",
			token: adminToken,
			prepare: func() int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:  "successful delete",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "Rem", "X")
				id, _ := res.LastInsertId()
				return id
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			id := tc.prepare()
			url := fmt.Sprintf("/category/%d", id)
			if id == 0 && tc.name == "invalid id bad request" {
				url = "/category/xyz"
			}
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func SeedSupplierForTests(t *testing.T, email, pass, company string) int64 {
	t.Helper()
	res, err := config.DB.Exec("insert into suppliers(name, contact_info, email, company) values(?,?,?,?)",
		"TestSupplier", "0000000000", email, company)
	if err != nil {
		t.Fatalf("seed supplier insert error: %v", err)
	}
	sid, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("seed supplier lastid error: %v", err)
	}
	hash, _ := utils.HashPwd(pass)
	_, err = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"SupplierUser", email, hash, "supplier_admin")
	if err != nil {
		t.Fatalf("seed supplier user insert error: %v", err)
	}

	return sid
}

func TestCreateSupplier_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	TruncateAll(t)
	adminEmail, adminPass := "admin_create_sup@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.POST("/supplier", functions.CreateSupplier)

	cases := []struct {
		name         string
		token        string
		body         string
		expectedCode int
	}{
		{
			name:         "forbidden if not system admin",
			token:        "",
			body:         `{"name":"S1","contact_info":"C","email":"s1p@test.com","company":"Co"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid payload",
			token:        adminToken,
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "valid create",
			token:        adminToken,
			body:         `{"name":"SCreate","contact_info":"9999999","email":"screate@test.com","company":"CreateCo"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invvalid contact_info",
			token:        adminToken,
			body:         `{"name":"SCreate","contact_info":"abcdef","email":"screate1@test.com","company":"CreateCo"}`,
			expectedCode: http.StatusBadRequest,
		},

		{
			name:         "duplicate email",
			token:        adminToken,
			body:         `{"name":"SCreate2","contact_info":"CINF","email":"screate@test.com","company":"CreateCo"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			_, _ = config.DB.Exec("delete from suppliers")
			req := httptest.NewRequest(http.MethodPost, "/supplier", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}
