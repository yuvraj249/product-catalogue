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

	if err := config.ExecMigration(db, mig); err != nil {
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
	adminEmail, adminPass := "admin_td@test.com", "Admin@123"
	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.POST("/categories", functions.CreateCategory)

	cases := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "no token",
			body:         `{"category_name":"X","category_description":"Y"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "supplier_admin cannot create",
			body:         `{"category_name":"X","category_description":"Y"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "invalid json",
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid category name",
			body:         `{"category_name":"@#$","category_description":"valid"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "valid create",
			body:         `{"category_name":"Electronics","category_description":"Devices"}`,
			expectedCode: http.StatusCreated,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)

			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
				"SupUser", "sup_td@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
			supToken := LoginAndGetToken(t, "sup_td@test.com", "Yuvraj@2411")
			req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(tc.body))

			switch tc.name {
			case "invalid json", "invalid category name", "valid create":
				if adminToken == "" {
					t.Fatalf("token empty for the subset: %s", tc.name)
				} else {

					req.Header.Set("Authorization", "Bearer "+adminToken)
				}
			case "supplier_admin cannot create":
				if supToken == " " {
					t.Fatalf("token empty for the subset: %s", tc.name)
				} else {
					req.Header.Set("Authorization", "Bearer "+supToken)
				}
			case "no token":
				req.Header.Set("Authorization", "Bearer "+"")
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

	adminEmail, adminPass := "admin_get@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup2", "sup2d@test.com", "Comp")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"Sup2User", "sup2d@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supplierToken := LoginAndGetToken(t, "sup2d@test.com", "Yuvraj@2411")

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.GET("/categories", functions.GetCategory)

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

			req := httptest.NewRequest(http.MethodGet, "/categories", nil)
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

	adminEmail, adminPass := "admin_getid@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup3", "sup3d@test.com", "Comp")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"Sup3User", "sup3d@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supplierToken := LoginAndGetToken(t, "sup3d@test.com", "Yuvraj@2411")

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.GET("/categories/:id", functions.GetCategoryByID)

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
			url:          "/categories/abc",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "not found",
			token: adminToken,
			prepare: func() int64 {

				return 0
			},
			url:          "/categories/9999999",
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
			url:          "/categories/1",
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
				url = fmt.Sprintf("/categories/%d", id)
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

	adminEmail, adminPass := "admin_update@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup3", "sup3d@test.com", "Comp")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"Sup3User", "sup3d@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supplierToken := LoginAndGetToken(t, "sup3d@test.com", "Yuvraj@2411")

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.PUT("/categories/:id", functions.UpdateCategory)

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
			url := fmt.Sprintf("/categories/%d", id)
			if id == 0 && tc.name == "invalid id" {
				url = "/categories/abc"
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

	adminEmail, adminPass := "admin_delete@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)

	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SupDel", "sup_del@test.com", "C")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"SupDelUser", "sup_del@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supToken := LoginAndGetToken(t, "sup_del@test.com", "Yuvraj@2411")

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.DELETE("/categories/:id", functions.DeleteCategory)

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
			url := fmt.Sprintf("/categories/%d", id)
			if id == 0 && tc.name == "invalid id bad request" {
				url = "/categories/xyz"
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

func TestCreateSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_create_sup@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SupDel", "sup_del@test.com", "C")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"SupDelUser", "sup_del@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supToken := LoginAndGetToken(t, "sup_del@test.com", "Yuvraj@2411")

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
			name:         "token missing",
			token:        "",
			body:         `{"name":"S1","contact_info":"C","email":"s1p@test.com","company":"Co"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "forbidden for supplier_admin",
			token:        supToken,
			body:         `{"name":"S1","contact_info":"C","email":"s1p@test.com","company":"Co"}`,
			expectedCode: http.StatusForbidden,
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
		TruncateAll(t)
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

func TestUpdateSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_delete@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SupDel", "sup_del@test.com", "C")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"SupDelUser", "sup_del@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supToken := LoginAndGetToken(t, "sup_del@test.com", "Yuvraj@2411")
	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.PUT("/suppliers/:id", functions.UpdateSupplier)
	testcases := []struct {
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
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Tester4", "99999999", "tester@gmail.com", "Apple")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "forbidden for supplier_admin",
			token: supToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Tester4", "99999999", "tester@gmail.com", "Apple")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"supplier_name":"Stubbs"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name:  "invalid id bad request",
			token: adminToken,
			prepare: func() int64 {
				return 0
			},
			body:         `{"supplier_name":"NewName"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "valid update",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Testerrrr", "99999999", "tester@gmail.com", "Apple")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"name":"Testerrrr", "contact_info":"989898989", "email":"tester@gmail.com", "company":"Apple"}`,
			expectedCode: http.StatusOK,
			verify: func(id int64) error {
				var name string
				if err := config.DB.QueryRow("select name from suppliers where supplier_id=?", id).Scan(&name); err != nil {
					return err
				}
				if name != "Tester4" {
					return fmt.Errorf("expected name Hardware got %s", name)
				}
				return nil
			},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			id := tc.prepare()
			url := fmt.Sprintf("/suppliers/%d", id)
			if id == 0 && tc.name == "invalid id bad request" {
				url = "/suppliers/xyz"
			}
			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.body))
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

func TestDeleteSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_create_sup@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	_, _ = config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SupDel", "sup_del@test.com", "C")
	_, _ = config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)",
		"SupDelUser", "sup_del@test.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "supplier_admin")
	supToken := LoginAndGetToken(t, "sup_del@test.com", "Yuvraj@2411")
	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.DELETE("/suppliers/:id", functions.DeleteSupplier)

	testcases := []struct {
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
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Tester4", "99999999", "tester@gmail.com", "Apple")
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
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Tester4", "99999999", "tester@gmail.com", "Apple")
				id, _ := res.LastInsertId()
				return id
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			id := tc.prepare()
			url := fmt.Sprintf("/suppliers/%d", id)
			if id == 0 && tc.name == "invalid id bad request" {
				url = "/suppliers/xyz"
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

func TestGetSupplier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_getid@test.com", "Admin@123"
	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.GET("/suppliers", functions.GetSupplier)

	testcases := []struct {
		name         string
		prepare      func()
		expectedCode int
		expectBody   string
	}{
		{
			name: "unauthorized (no token)",
			prepare: func() {
				_, _ = config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "ABC", "111", "a@test", "CoA")
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "empty list for system_admin",
			prepare: func() {

			},
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin sees all suppliers",
			prepare: func() {
				_, _ = config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "ABC", "111", "a@test", "CoA")
				_, _ = config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "CFD", "222", "b@test", "CoB")
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin sees all suppliers of same company",
			prepare: func() {
				_, _ = config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "ABC", "111", "a@test", "CoA")
				_, _ = config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "CFD", "222", "b@test", "CoA")
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)
			if tc.prepare != nil {
				tc.prepare()
			}
			req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
			if tc.name == "supplier_admin sees all suppliers of same company" {
				var sid int64
				err := config.DB.QueryRow("select supplier_id from suppliers where company= ?", "CoA").Scan(&sid)
				if err != nil {
					t.Fatalf("failed to find supplier row inserted in prepare: %v", err)
				}
				pass := "Yuvraj@2411"
				hash, _ := utils.HashPwd(pass)
				_, err = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)",
					"SupUser", "supget@test.com", hash, "supplier_admin", sid)
				if err != nil {
					t.Fatalf("failed to insert supplier user: %v", err)
				}
				supToken := LoginAndGetToken(t, "supget@test.com", pass)
				if supToken == "" {
					t.Fatalf("supplier login failed for subtest %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+supToken)
			} else if tc.expectedCode != http.StatusUnauthorized {
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("%sexpected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
			if tc.expectedCode == http.StatusOK && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("%s expected body to contain %q; got %s", tc.name, tc.expectBody, w.Body.String())
			}

		})

	}

}

func TestGetSupplierByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_getbyid@test.com", "Admin@123"
	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.GET("/suppliers/:id", functions.GetSupplierByID)

	testcases := []struct {
		name         string
		prepare      func() int64
		expectedCode int
		expectBody   string
	}{
		{
			name: "invalid id param",
			prepare: func() int64 {
				return 0
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "unauthorized (no token)",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"UnauthSup", "111", "unauth@test.com", "CoX")
				id, _ := res.LastInsertId()
				return id
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "not found (system_admin)",
			prepare: func() int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "system_admin found",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"SysFound", "222", "sysfound@test.com", "CoSys")
				id, _ := res.LastInsertId()
				return id
			},
			expectedCode: http.StatusOK,
			expectBody:   "SysFound",
		},
		{
			name: "supplier_admin can get supplier in same company",
			prepare: func() int64 {
				res1, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CoASup1", "333", "coasup1@test.com", "CoA")
				id1, _ := res1.LastInsertId()

				_, _ = config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CoASup2", "444", "coasup2@test.com", "CoA")
				pass := "Supplier@123"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)",
					"CoAAdmin", "coasupadmin@test.com", hash, "supplier_admin", id1)

				res2, _ := config.DB.Exec("select supplier_id from suppliers where email = ?", "coasup2@test.com")
				var targetID int64
				_ = res2
				_ = config.DB.QueryRow("select supplier_id from suppliers where email = ?", "coasup2@test.com").Scan(&targetID)
				return targetID
			},
			expectedCode: http.StatusOK,
			expectBody:   "CoASup2",
		},
		{
			name: "supplier_admin forbidden for different company",
			prepare: func() int64 {
				resA, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CompanyX", "555", "compx@test.com", "CompanyX")
				idA, _ := resA.LastInsertId()

				resB, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CompanyY", "666", "compy@test.com", "CompanyY")
				idB, _ := resB.LastInsertId()
				pass := "Supplier@321"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)",
					"CompXAdmin", "compxadmin@test.com", hash, "supplier_admin", idA)
				return idB
			},
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)

			SeedAdmin(t, adminEmail, adminPass)
			adminToken := LoginAndGetToken(t, adminEmail, adminPass)

			var id int64 = 0
			if tc.prepare != nil {
				id = tc.prepare()
			}

			url := ""
			if tc.name == "invalid id param" {
				url = "/suppliers/abc"
			} else if id == 0 {
				url = fmt.Sprintf("/suppliers/%d", 99999999)
			} else {
				url = fmt.Sprintf("/suppliers/%d", id)
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			switch tc.name {
			case "unauthorized (no token)":
			case "supplier_admin can get supplier in same company":
				supToken := LoginAndGetToken(t, "coasupadmin@test.com", "Supplier@123")
				if supToken == "" {
					t.Fatalf("failed to login supplier admin for subtest %s", tc.name)
				}
				req.Header.Set("Authorization", "Bearer "+supToken)
			case "supplier_admin forbidden for different company":
				supToken := LoginAndGetToken(t, "compxadmin@test.com", "Supplier@321")
				if supToken == "" {
					t.Fatalf("failed to login supplier admin for subtest %s", tc.name)
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
			if tc.expectedCode == http.StatusOK && tc.expectBody != "" && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("%s expected body to contain %q; got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestCreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_prod_create@test.com", "Admin@123"

	r := gin.New()
	r.Use(middleware.AuthMiddleware())
	r.POST("/products", functions.CreateProduct)

	testcases := []struct {
		name         string
		prepare      func() (supplierToken string, catID int)
		body         string
		expectedCode int
	}{
		{
			name: "unauthorized no token",
			prepare: func() (string, int) {
				return "", 0
			},
			body:         `{"product_name":"Product1","product_description":"D","product_cost":100,"product_category_id":1}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for system_admin",
			prepare: func() (string, int) {
				SeedAdmin(t, adminEmail, adminPass)
				adminToken := LoginAndGetToken(t, adminEmail, adminPass)
				res, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatA", "X")
				cid, _ := res.LastInsertId()
				return adminToken, int(cid)
			},
			body:         `{"product_name":"P1","product_description":"D","product_cost":100,"product_category_id":1}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid payload JSON",
			prepare: func() (string, int) {
				SeedAdmin(t, adminEmail, adminPass)
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "SupP", "999", "supprod@test.com", "CoP")
				sid, _ := res.LastInsertId()
				pass := "SupPass@1"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SuppUser", "supprod@test.com", hash, "supplier_admin", sid)
				token := LoginAndGetToken(t, "supprod@test.com", pass)
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatP", "desc")
				cid, _ := rc.LastInsertId()
				return token, int(cid)
			},
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid product name (empty)",
			prepare: func() (string, int) {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "SupP", "999", "supprod@test.com", "CoP")
				sid, _ := res.LastInsertId()
				pass := "SupPass@1"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SuppUser", "supprod@test.com", hash, "supplier_admin", sid)
				token := LoginAndGetToken(t, "supprod@test.com", pass)

				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatP", "desc")
				cid, _ := rc.LastInsertId()
				return token, int(cid)
			},
			body:         `{"product_name":"   ","product_description":"D","product_cost":10,"product_category_id":1}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid category id",
			prepare: func() (string, int) {

				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "SupP", "999", "supprod@test.com", "CoP")
				sid, _ := res.LastInsertId()
				pass := "SupPass@1"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SuppUser", "supprod@test.com", hash, "supplier_admin", sid)
				token := LoginAndGetToken(t, "supprod@test.com", pass)
				return token, 999999
			},
			body:         `{"product_name":"Pvalid","product_description":"D","product_cost":10,"product_category_id":999999}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "valid create by supplier_admin",
			prepare: func() (string, int) {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "SupP", "999", "supprod@test.com", "CoP")
				sid, _ := res.LastInsertId()
				pass := "SupPass@1"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SuppUser", "supprod@test.com", hash, "supplier_admin", sid)
				token := LoginAndGetToken(t, "supprod@test.com", pass)
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatP", "desc")
				cid, _ := rc.LastInsertId()
				return token, int(cid)
			},
			body:         `{"product_name":"Pvalid","product_description":"D","product_cost":10,"product_category_id":1}`,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			TruncateAll(t)
			token, cid := "", 0
			if tc.prepare != nil {
				token, cid = tc.prepare()
			}
			body := tc.body
			if cid > 0 {
				body = strings.Replace(body, `"product_category_id":1`, fmt.Sprintf(`"product_category_id":%d`, cid), 1)
				body = strings.Replace(body, `"product_category_id":999999`, fmt.Sprintf(`"product_category_id":%d`, cid), 1)
			}

			req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedCode {
				t.Fatalf("[%s] expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}
