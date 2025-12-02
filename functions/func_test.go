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
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	r := routes.SetupRouter()

	cases := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "no token",
			body:         `{"category_name":"Xsndin","category_description":"Yadaii"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "supplier_admin cannot create",
			body:         `{"category_name":"Xidjais","category_description":"Ysidisd"}`,
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
			t.Cleanup(func() {
				TruncateAll(t)
			})
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

	r := routes.SetupRouter()

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
			t.Cleanup(func() {
				TruncateAll(t)
			})
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
	r := routes.SetupRouter()

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
			t.Cleanup(func() {
				TruncateAll(t)
			})
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

	r := routes.SetupRouter()

	cases := []struct {
		name         string
		token        string
		prepare      func() int64
		body         string
		expectedCode int
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
		},
		{
			name:  "duplicate category_name",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "Category", "Hardware")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"category_name":"Category"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "duplicate category_description",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into categories(category_name, category_description) values(?,?)", "Category", "Hardware")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"category_description":"Hardware"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
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
				t.Fatalf("%sexpected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
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

	r := routes.SetupRouter()

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
			t.Cleanup(func() {
				TruncateAll(t)
			})
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

	r := routes.SetupRouter()

	cases := []struct {
		name         string
		token        string
		prepare      func()
		body         string
		expectedCode int
	}{
		{
			name:         "token missing",
			token:        "",
			body:         `{"name":"Ssadhah","contact_info":"34626323737","email":"s1p@test.com","company":"Co"}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "forbidden for supplier_admin",
			token:        supToken,
			body:         `{"name":"Sahdshd","contact_info":"57343747347","email":"s1p@test.com","company":"Co"}`,
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
			name:  "duplicate email",
			token: adminToken,
			prepare: func() {
				_, _ = config.DB.Exec("insert into suppliers(name, contact_info,email, company) values(?,?,?,?)", "sdshdsh", "5734737732", "screate@test.com", "CreateCo")
			},
			body:         `{"name":"SCreate","contact_info":"6757475399","email":"screate@test.com","company":"CreateCo"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
			_, _ = config.DB.Exec("delete from suppliers")
			if tc.prepare != nil {
				tc.prepare()
			}
			req := httptest.NewRequest(http.MethodPost, "/suppliers", strings.NewReader(tc.body))
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
	r := routes.SetupRouter()
	testcases := []struct {
		name         string
		token        string
		prepare      func() int64
		body         string
		expectedCode int
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
		},
		{
			name:  "duplicate name",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Tester", "99999999", "tester@gmail.com", "Apple")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"name":"Tester"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "duplicate contact_info",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Tester", "99999999", "tester@gmail.com", "Apple")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"contact_info":"99999999"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "duplicate email",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Tester", "99999999", "tester@gmail.com", "Apple")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"email":"tester@gmail.com"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:  "duplicate company",
			token: adminToken,
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Tester", "99999999", "tester@gmail.com", "Apple")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"company":"Apple"}`,
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
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
	r := routes.SetupRouter()
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
			t.Cleanup(func() {
				TruncateAll(t)
			})
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
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	r := routes.SetupRouter()

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
			t.Cleanup(func() {
				TruncateAll(t)
			})
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
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	r := routes.SetupRouter()

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
			name: "supplier_admin can get supplier of the same company",
			prepare: func() int64 {
				res1, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CoASup1", "333", "coasup1@test.com", "CoA")
				id1, _ := res1.LastInsertId()

				_, _ = config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CoASup2", "444", "coasup2@test.com", "CoB")
				pass := "Supplier@123"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)",
					"CoAAdmin", "coasupadmin@test.com", hash, "supplier_admin", id1)

				res2 := config.DB.QueryRow("select supplier_id from suppliers where company = ?", "CoA")
				var targetID int64
				_ = res2
				_ = config.DB.QueryRow("select supplier_id from suppliers where company = ?", "CoA").Scan(&targetID)
				return targetID
			},
			expectedCode: http.StatusOK,
			expectBody:   "CoASup1",
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
			t.Cleanup(func() {
				TruncateAll(t)
			})
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
				req.Header.Set("Authorization", "Bearer "+"")
			case "supplier_admin can get supplier of the same company":
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
	r := routes.SetupRouter()
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
			body:         `{"product_name":"Product1","product_description":"dajjajs","product_cost":100,"product_category_id":1}`,
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
			body:         `{"product_name":"P1","product_description":"dajjajs","product_cost":100,"product_category_id":1}`,
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
			body:         `{"product_name":"   ","product_description":"dajjajs","product_cost":10,"product_category_id":1}`,
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
			body:         `{"product_name":"Pvalid","product_description":"dajjajs","product_cost":10,"product_category_id":999999}`,
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
			body:         `{"product_name":"Pvalid","product_description":"dajjajs","product_cost":10,"product_category_id":1}`,
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
			t.Cleanup(func() {
				t.Cleanup(func() {
					TruncateAll(t)
				})
			})
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
				t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

func TestGetProducts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_prod_gets@test.com", "Admin@123"
	r := routes.SetupRouter()
	testcases := []struct {
		name         string
		prepare      func()
		expectedCode int
		expectBody   string
	}{
		{
			name: "unauthorized",
			prepare: func() {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Sx", "111", "sx@test.com", "CoX")
				sid, _ := res.LastInsertId()
				_, _ = config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatX", "dajjajs")
				cidRes, _ := config.DB.Exec("select category_id from categories where category_name = ?", "CatX")
				_ = cidRes

				_, _ = config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "Px", "dajjajs", 10, 1, sid)
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "system_admin sees all",
			prepare: func() {
				SeedAdmin(t, adminEmail, adminPass)
				res1, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "S1", "111", "s1@test", "CoA")
				sid1, _ := res1.LastInsertId()
				res2, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "S2", "222", "s2@test", "CoB")
				sid2, _ := res2.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "Cat1", "dajjajs")
				cid, _ := rc.LastInsertId()
				_, _ = config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "Prod1", "dajjajs", 11.5, cid, sid1)
				_, _ = config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "Prod2", "dajjajs", 22.5, cid, sid2)
			},
			expectedCode: http.StatusOK,
			expectBody:   "Prod1",
		},
		{
			name: "supplier_admin sees same company products only",
			prepare: func() {
				SeedAdmin(t, adminEmail, adminPass)
				res1, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "SC1", "111", "sc1@test.com", "SameCo")
				sid1, _ := res1.LastInsertId()
				_, _ = config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "SC2", "222", "sc2@test.com", "SameCo")
				res3, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "diddy", "333", "dx@test.com", "OtherCo")
				sid3, _ := res3.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatS", "dajjajs")
				cid, _ := rc.LastInsertId()
				_, _ = config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "SameCoProd", "dajjajs", 10, cid, sid1)
				_, _ = config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "OtherProd", "dajjajs", 20, cid, sid3)
				pass := "SupSee@1"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SupSee", "supsee@test.com", hash, "supplier_admin", sid1)
			},
			expectedCode: http.StatusOK,
			expectBody:   "SameCoProd",
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
			if tc.prepare != nil {
				tc.prepare()
			}

			req := httptest.NewRequest(http.MethodGet, "/products", nil)

			if tc.name == "system_admin sees all" || tc.name == "supplier_admin sees same company products only" {
				adminToken := LoginAndGetToken(t, adminEmail, adminPass)
				if tc.name == "system_admin sees all" {
					req.Header.Set("Authorization", "Bearer "+adminToken)
				} else {

					supToken := LoginAndGetToken(t, "supsee@test.com", "SupSee@1")
					if supToken == "" {
						t.Fatalf("supplier login failed for subtest %s", tc.name)
					}
					req.Header.Set("Authorization", "Bearer "+supToken)
				}
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tc.name == "unauthorized" {
				if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
					t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
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
	adminEmail, adminPass := "admin_prod_getssss@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)

	r := routes.SetupRouter()
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
			name: "unauthorized",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Sx", "111", "sx@test.com", "CoX")
				sid, _ := res.LastInsertId()
				_, _ = config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatX", "dajjajs")
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "Px", "dajjajs", 10, 1, sid)
				pid, _ := prod.LastInsertId()
				return pid
			},

			expectedCode: http.StatusForbidden,
		},
		{
			name: "not found (system_admin)",
			prepare: func() int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "system_admin can see any product by id",
			prepare: func() int64 {
				res1, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "S1", "111", "s1@test", "CoA")
				sid1, _ := res1.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "Cat1", "dajjajs")
				cid, _ := rc.LastInsertId()
				prod1, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "Prod1", "dajjajs", 11.5, cid, sid1)
				pid1, _ := prod1.LastInsertId()
				return pid1

			},
			expectedCode: http.StatusOK,
			expectBody:   "Prod1",
		},
		{
			name: "supplier_admin can get products made by supplier admin of the same company",
			prepare: func() int64 {
				res1, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CoASup1", "333", "coasup1@test.com", "CoA")
				id1, _ := res1.LastInsertId()

				res2, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CoASup2", "444", "coasup2@test.com", "CoA")
				id2, _ := res2.LastInsertId()
				cid, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "Cat1", "desc")
				cID, _ := cid.LastInsertId()
				prod1, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "Prod1", "dish", 11.5, cID, id1)
				_, _ = prod1.LastInsertId()
				prod2, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "Prod2", "dishtv", 11.5, cID, id2)
				pid2, _ := prod2.LastInsertId()
				pass := "Supplier@123"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)",
					"CoAAdmin", "coasupadmin@test.com", hash, "supplier_admin", id1)
				return pid2
			},
			expectedCode: http.StatusOK,
			expectBody:   "Prod2",
		},
		{
			name: "supplier_admin cannot see products made by supplier admin of the different company",
			prepare: func() int64 {
				res1, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CoASup1", "333", "coasup1@test.com", "CoA")
				id1, _ := res1.LastInsertId()

				res2, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)",
					"CoASup2", "444", "coasup2@test.com", "CoB")
				id2, _ := res2.LastInsertId()
				cid, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "Cat1", "desc")
				cID, _ := cid.LastInsertId()
				prod1, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "Prod1", "dish", 11.5, cID, id1)
				_, _ = prod1.LastInsertId()
				prod2, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "Prod2", "dishtv", 11.5, cID, id2)
				pid2, _ := prod2.LastInsertId()
				pass := "Supplier@123"
				hash, _ := utils.HashPwd(pass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)",
					"CoAAdmin", "coasupadmin@test.com", hash, "supplier_admin", id1)
				return pid2
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
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
			req := httptest.NewRequest(http.MethodGet, url, nil)
			switch tc.name {
			case "system_admin can see any product by id", "invalid id param", "not found (system_admin)":
				req.Header.Set("Authorization", "Bearer "+adminToken)
			case "supplier_admin can get products made by supplier admin of the same company", "supplier_admin cannot see products made by supplier admin of the different company":
				supToken := LoginAndGetToken(t, "coasupadmin@test.com", "Supplier@123")
				req.Header.Set("Authorization", "Bearer "+supToken)
			case "unauthorized":
				req.Header.Set("Authorization", "Bearer "+"")

			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if tc.name == "unauthorized" {
				if w.Code != http.StatusUnauthorized && w.Code != http.StatusForbidden {
					t.Fatalf("%s expected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
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
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	const supPass = "Supplier@123"

	r := routes.SetupRouter()
	testcases := []struct {
		name         string
		prepare      func() (int64, string)
		body         string
		expectedCode int
		expectBody   string
	}{
		{
			name: "invalid id param",
			prepare: func() (int64, string) {
				resSup, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "UpdSup", "111", "upd_sup@test.com", "UpCo")
				sid, _ := resSup.LastInsertId()
				rCatg, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatUPD", "dajjajs")
				cid, _ := rCatg.LastInsertId()
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "PUpd", "dajjajs", 9.9, cid, sid)
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SupUpd", "supupd@test.com", hash, "supplier_admin", sid)
				return pid, "supupd@test.com"
			},
			body:         `{"product_name":"Xdfifn"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "unauthorized (no token)",
			prepare: func() (int64, string) {
				resS, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "UnauthSup", "111", "unauthsup@test.com", "XCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatUA", "dcsdd")
				cid, _ := rc.LastInsertId()
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "UnauthP", "ddadd", 5.5, cid, sid)
				pid, _ := prod.LastInsertId()
				return pid, ""
			},
			body:         `{}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for system_admin",
			prepare: func() (int64, string) {
				resS, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "SysOwner", "111", "sysowner@test.com", "SysCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatSys", "dfwsfcw")
				cid, _ := rc.LastInsertId()
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "SysProd", "dssd", 12.5, cid, sid)
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SysSupp", "syssupp@test.com", hash, "supplier_admin", sid)
				return pid, "syssupp@test.com"
			},
			body:         `{"product_name":"NewName"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "no fields provided",
			prepare: func() (int64, string) {
				resS, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "OwnerA", "111", "ownera@test.com", "Aco")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatNF", "dcscsd")
				cid, _ := rc.LastInsertId()
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "NoFieldP", "dsdcs", 6.6, cid, sid)
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "OwnerAUser", "ownerauser@test.com", hash, "supplier_admin", sid)
				return pid, "ownerauser@test.com"
			},
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "cannot change product supplier",
			prepare: func() (int64, string) {
				resS, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "OwnerB", "111", "ownerb@test.com", "BCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatCS", "dsc")
				cid, _ := rc.LastInsertId()
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "ChangeSuppP", "descd", 7.7, cid, sid)
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "OwnerBUser", "ownerbuser@test.com", hash, "supplier_admin", sid)
				return pid, "ownerbuser@test.com"
			},
			body:         `{"product_supplier_id": 999}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "product not found",
			prepare: func() (int64, string) {
				resS, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "NoProdSup", "111", "nops@test.com", "NPS")
				sid, _ := resS.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "NoProdUser", "nopsuser@test.com", hash, "supplier_admin", sid)
				return 99999999, "nopsuser@test.com"
			},
			body:         `{"product_name":"Xsdhjd"}`,
			expectedCode: http.StatusNotFound,
		},
		{
			name: "forbidden update other supplier's product",
			prepare: func() (int64, string) {
				resA, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Asijs", "111", "a@test.com", "Aco")
				idA, _ := resA.LastInsertId()
				resB, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "Bsihsih", "222", "b@test.com", "Bco")
				idB, _ := resB.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatFO", "dajjajs")
				cid, _ := rc.LastInsertId()
				prodB, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "BProd", "dajjajs", 3.3, cid, idB)
				pidB, _ := prodB.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "AAdmin", "aadmin@test.com", hash, "supplier_admin", idA)
				return pidB, "aadmin@test.com"
			},
			body:         `{"product_name":"ILegal"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid category id in payload",
			prepare: func() (int64, string) {
				resS, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "CatOwner", "111", "catowner@test.com", "Cco")
				idS, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatValid", "dajjajs")
				cid, _ := rc.LastInsertId()
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "CatCheck", "dajjajs", 14.5, cid, idS)
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "CatOwnerU", "catowneru@test.com", hash, "supplier_admin", idS)
				return pid, "catowneru@test.com"
			},
			body:         `{"product_category_id": 9999999}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "successful update",
			prepare: func() (int64, string) {
				resS, err := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "GoodOwner", "111", "goodowner@test.com", "Gco")
				if err != nil {
					t.Fatalf("prepare supplier failed: %v", err)
				}
				sid, _ := resS.LastInsertId()
				rc, err := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatGood", "dajjajs")
				if err != nil {
					t.Fatalf("prepare category failed: %v", err)
				}
				cid, _ := rc.LastInsertId()
				prod, err := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "GoodProd", "dajjajs", 20.0, cid, sid)
				if err != nil {
					t.Fatalf("prepare product failed: %v", err)
				}
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "GoodSupUser", "goodsup@test.com", hash, "supplier_admin", sid)
				return pid, "goodsup@test.com"
			},
			body:         `{"product_name":"GoodProdUpdated","product_description":"new","product_cost":21.5}`,
			expectedCode: http.StatusOK,
			expectBody:   "product updated",
		},
		{
			name: "duplicate product_name",
			prepare: func() (int64, string) {
				resS, err := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "GoodOwner", "111", "goodowner@test.com", "Gco")
				if err != nil {
					t.Fatalf("prepare supplier failed: %v", err)
				}
				sid, _ := resS.LastInsertId()
				rc, err := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatGood", "dajjajs")
				if err != nil {
					t.Fatalf("prepare category failed: %v", err)
				}
				cid, _ := rc.LastInsertId()
				prod, err := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "GoodProd", "dajjajs", 20.0, cid, sid)
				if err != nil {
					t.Fatalf("prepare product failed: %v", err)
				}
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "GoodSupUser", "goodsup@test.com", hash, "supplier_admin", sid)
				return pid, "goodsup@test.com"
			},
			body:         `{"product_name":"GoodProd"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate product_description",
			prepare: func() (int64, string) {
				resS, err := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "GoodOwner", "111", "goodowner@test.com", "Gco")
				if err != nil {
					t.Fatalf("prepare supplier failed: %v", err)
				}
				sid, _ := resS.LastInsertId()
				rc, err := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatGood", "dajjajs")
				if err != nil {
					t.Fatalf("prepare category failed: %v", err)
				}
				cid, _ := rc.LastInsertId()
				prod, err := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "GoodProd", "dajjajs", 20.0, cid, sid)
				if err != nil {
					t.Fatalf("prepare product failed: %v", err)
				}
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "GoodSupUser", "goodsup@test.com", hash, "supplier_admin", sid)
				return pid, "goodsup@test.com"
			},
			body:         `{"product_description":"dajjajs"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate product_cost",
			prepare: func() (int64, string) {
				resS, err := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "GoodOwner", "111", "goodowner@test.com", "Gco")
				if err != nil {
					t.Fatalf("prepare supplier failed: %v", err)
				}
				sid, _ := resS.LastInsertId()
				rc, err := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatGood", "dajjajs")
				if err != nil {
					t.Fatalf("prepare category failed: %v", err)
				}
				cid, _ := rc.LastInsertId()
				prod, err := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "GoodProd", "dajjajs", 20.0, cid, sid)
				if err != nil {
					t.Fatalf("prepare product failed: %v", err)
				}
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "GoodSupUser", "goodsup@test.com", hash, "supplier_admin", sid)
				return pid, "goodsup@test.com"
			},
			body:         `{"product_cost": 20.0}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate product_category_id",
			prepare: func() (int64, string) {
				resS, err := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "GoodOwner", "111", "goodowner@test.com", "Gco")
				if err != nil {
					t.Fatalf("prepare supplier failed: %v", err)
				}
				sid, _ := resS.LastInsertId()
				rc, err := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatGood", "dajjajs")
				if err != nil {
					t.Fatalf("prepare category failed: %v", err)
				}
				cid, _ := rc.LastInsertId()
				prod, err := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "GoodProd", "dajjajs", 20.0, cid, sid)
				if err != nil {
					t.Fatalf("prepare product failed: %v", err)
				}
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "GoodSupUser", "goodsup@test.com", hash, "supplier_admin", sid)
				return pid, "goodsup@test.com"
			},
			body:         `{"product_category_id":1}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "duplicate discount_type and discount_value",
			prepare: func() (int64, string) {
				resS, err := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "GoodOwner", "111", "goodowner@test.com", "Gco")
				if err != nil {
					t.Fatalf("prepare supplier failed: %v", err)
				}
				sid, _ := resS.LastInsertId()
				rc, err := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatGood", "dajjajs")
				if err != nil {
					t.Fatalf("prepare category failed: %v", err)
				}
				cid, _ := rc.LastInsertId()
				prod, err := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id, discount_type, discount_value) values(?,?,?,?,?,?,?)", "GoodProd", "dajjajs", 20.0, cid, sid, "percent", 50)
				if err != nil {
					t.Fatalf("prepare product failed: %v", err)
				}
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "GoodSupUser", "goodsup@test.com", hash, "supplier_admin", sid)
				return pid, "goodsup@test.com"
			},
			body:         `{"discount_type":"percent","discount_value":50}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
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

			req := httptest.NewRequest(http.MethodPut, url, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
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
			if tc.expectBody != "" && !strings.Contains(w.Body.String(), tc.expectBody) {
				t.Fatalf("%s expected body to contain %q got %s", tc.name, tc.expectBody, w.Body.String())
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	adminEmail, adminPass := "admin_prod_delete@test.com", "Admin@123"
	SeedAdmin(t, adminEmail, adminPass)
	adminToken := LoginAndGetToken(t, adminEmail, adminPass)
	const supPass = "Supplier@123"
	r := routes.SetupRouter()
	testcases := []struct {
		name         string
		prepare      func() (int64, string)
		expectedCode int
	}{
		{
			name: "invalid id param",
			prepare: func() (int64, string) {
				resSup, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "InvSup", "111", "invsup@test.com", "InvCo")
				sid, _ := resSup.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatInv", "inv")
				cid, _ := rc.LastInsertId()
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "InvProd", "inv", 1.1, cid, sid)
				_, _ = prod.LastInsertId()

				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "InvSupUser", "invsupuser@test.com", hash, "supplier_admin", sid)
				return 0, "invsupuser@test.com"
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "unauthorized (no token)",
			prepare: func() (int64, string) {
				resSup, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "UnauthSup", "111", "unauthsup@test.com", "XCo")
				sid, _ := resSup.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatUA", "dcsdd")
				cid, _ := rc.LastInsertId()
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "UnauthP", "ddadd", 5.5, cid, sid)
				pid, _ := prod.LastInsertId()
				return pid, ""
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for system_admin",
			prepare: func() (int64, string) {
				resSup, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "SysDelSup", "111", "sysdelsup@test.com", "SysCo")
				sid, _ := resSup.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatSysDel", "d")
				cid, _ := rc.LastInsertId()
				prod, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "SysDelProd", "d", 2.2, cid, sid)
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SysDelUser", "sysdeluser@test.com", hash, "supplier_admin", sid)
				return pid, "sysdeluser@test.com"
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "product not found (supplier tries to delete nonexistent id)",
			prepare: func() (int64, string) {
				resSup, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "NoProdSup", "111", "nops@test.com", "NPS")
				sid, _ := resSup.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "NoProdUser", "nopsuser@test.com", hash, "supplier_admin", sid)
				return 99999999, "nopsuser@test.com"
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "forbidden delete other supplier's product",
			prepare: func() (int64, string) {
				resA, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "ASup", "111", "asup@test.com", "Aco")
				idA, _ := resA.LastInsertId()
				resB, _ := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "BSup", "222", "bsup@test.com", "Bco")
				idB, _ := resB.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatF", "d")
				cid, _ := rc.LastInsertId()
				prodB, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "BProd", "d", 3.3, cid, idB)
				pidB, _ := prodB.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "AAdmin", "aadmin@test.com", hash, "supplier_admin", idA)
				return pidB, "aadmin@test.com"
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "successful delete",
			prepare: func() (int64, string) {
				resS, err := config.DB.Exec("insert into suppliers(name,contact_info,email,company) values(?,?,?,?)", "GoodDelSup", "111", "gooddelsup@test.com", "Gco")
				if err != nil {
					t.Fatalf("prepare supplier failed: %v", err)
				}
				sid, _ := resS.LastInsertId()
				rc, err := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatDel", "d")
				if err != nil {
					t.Fatalf("prepare category failed: %v", err)
				}
				cid, _ := rc.LastInsertId()
				prod, err := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "GoodDeleteProd", "d", 10.0, cid, sid)
				if err != nil {
					t.Fatalf("prepare product failed: %v", err)
				}
				pid, _ := prod.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "GoodDelUser", "gooddeluser@test.com", hash, "supplier_admin", sid)
				return pid, "gooddeluser@test.com"
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
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
			if tc.name == "successful delete" && w.Code == http.StatusOK {
				var cnt int
				err := config.DB.QueryRow("select count(*) from products where product_id = ?", pID).Scan(&cnt)
				if err != nil {
					t.Fatalf("db check failed: %v", err)
				}
				if cnt != 0 {
					t.Fatalf("expected product to be deleted but found %d", cnt)
				}
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
		prepare      func() int64
		body         string
		expectedCode int
		expectBody   string
	}{
		{
			name: "unauthorized (no token)",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SNoTok", "snotok@test.com", "Co")
				id, _ := res.LastInsertId()
				return id
			},

			body:         `{name":"Ucsic","email":"u@test.com","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SF", "sf@test.com", "CoF")
				sid, _ := res.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SFUser", "sfuser@test.com", hash, "supplier_admin", sid)
				return sid
			},
			body:         `{name":"Usicis","email":"u2@test.com","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid json",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sij", "sij@test.com", "Co")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "missing fields",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Smiss", "smiss@test.com", "Co")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"name":"","email":"", "password":"", "supplier_id":0}`,
			expectedCode: http.StatusBadRequest,
		},
		{

			name: "invalid email",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sinv", "sinv@test.com", "Co")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"name": "Usjsj9oj","email":"notemail","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid password",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Spwd", "spwd@test.com", "Co")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"name":"Udifif","email":"usidjidhn@test.com","password":"not","supplier_id":1}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "successful create",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sgood", "sgood@test.com", "GCo")
				id, _ := res.LastInsertId()
				return id
			},
			body:         `{"name":"TesterAdmin","email":"newadmin@test.com","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusCreated,
			expectBody:   "supplier_admin user created successfully",
		},
		{
			name: "email already registered",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sdup", "sdup@test.com", "Co")
				sid, _ := res.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "Existing", "existing@test.com", hash, "supplier_admin", sid)
				return sid
			},
			body:         `{"name":"usbdusbd","email":"existing@test.com","password":"Supplier@123","supplier_id":1}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "supplier_id does not exist",
			prepare: func() int64 {
				return 0
			},
			body:         `{"name":"Ucsdd","email":"u4@test.com","password":"Supplier@123","supplier_id":999999}`,
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				t.Cleanup(func() {
					TruncateAll(t)
				})
			})
			if tc.prepare != nil {
				tc.prepare()
			}
			req := httptest.NewRequest(http.MethodPost, "/users/supplier-admin", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")

			case "forbidden for supplier_admin":

				supToken := LoginAndGetToken(t, "sfuser@test.com", supPass)
				if supToken != "" {
					req.Header.Set("Authorization", "Bearer "+supToken)
				}
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.expectedCode {
				t.Fatalf("%sexpected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
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
		prepare      func()
		expectedCode int
	}{
		{
			name: "unauthorized (no token)",
			prepare: func() {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Slusins", "slusdinn@test.com", "Co")
				sid, _ := res.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SLUser", "sluser@test.com", hash, "supplier_admin", sid)
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func() {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup", "sup@test.com", "Co")
				sid, _ := res.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SupUser", "supuser@test.com", hash, "supplier_admin", sid)
			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "system_admin authorized",
			prepare: func() {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "S1", "s1@test.com", "CoA")
				sid, _ := res.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "U1", "u1@test.com", hash, "supplier_admin", sid)
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			t.Cleanup(func() {
				t.Cleanup(func() {
					TruncateAll(t)
				})
			})
			if tc.prepare != nil {
				tc.prepare()
			}
			req := httptest.NewRequest(http.MethodGet, "/users/supplier-admin", nil)
			req.Header.Set("Content-Type", "application/json")
			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")
			case "forbidden for supplier_admin":
				supToken := LoginAndGetToken(t, "supuser@test.com", supPass)
				if supToken != "" {
					req.Header.Set("Authorization", "Bearer "+supToken)
				}
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.expectedCode {
				t.Fatalf("%sexpected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
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
		prepare      func() int64
		expectedCode int
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
				res, _ := config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)", "NoTok", "notok@test.com", "hash", "system_admin")
				id, _ := res.LastInsertId()
				return id
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "not found",
			prepare: func() int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "found",
			prepare: func() int64 {
				resSup, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "UsrS", "usrs@test.com", "Co")
				sid, _ := resSup.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				resU, _ := config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "FoundU", "found@test.com", hash, "supplier_admin", sid)
				uid, _ := resU.LastInsertId()
				return uid
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup", "sup@test.com", "Co")
				sid, _ := res.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				prod, _ := config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SupUser", "supuser@test.com", hash, "supplier_admin", sid)
				prodID, _ := prod.LastInsertId()
				return prodID

			},
			expectedCode: http.StatusForbidden,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			t.Cleanup(func() {
				t.Cleanup(func() {
					TruncateAll(t)
				})
			})
			var uid int64
			if tc.prepare != nil {
				uid = tc.prepare()
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
				supToken := LoginAndGetToken(t, "supuser@test.com", supPass)
				if supToken != "" {
					req.Header.Set("Authorization", "Bearer "+supToken)
				}
			default:
				req.Header.Set("Authorization", "Bearer "+adminToken)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.expectedCode {
				t.Fatalf("%sexpected %d got %d body=%s", tc.name, tc.expectedCode, w.Code, w.Body.String())
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
		prepare      func() int64
		expectedCode int
	}{
		{
			name: "unauthorized (no token)",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into users(name,email,password_hash,role) values(?,?,?,?)", "NoTok", "notok@test.com", "hash", "system_admin")
				id, _ := res.LastInsertId()
				return id
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid id param",
			prepare: func() int64 {
				return 0
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden for supplier_admin",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup", "sup@test.com", "Co")
				sid, _ := res.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				prod, _ := config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SupUser", "supuser@test.com", hash, "supplier_admin", sid)
				prodID, _ := prod.LastInsertId()
				return prodID

			},
			expectedCode: http.StatusForbidden,
		},
		{
			name: "not found",
			prepare: func() int64 {
				return 99999999
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "successful delete",
			prepare: func() int64 {
				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "Sup", "sup@test.com", "Co")
				sid, _ := res.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				prod, _ := config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "SupUser", "supuser@test.com", hash, "supplier_admin", sid)
				prodID, _ := prod.LastInsertId()
				return prodID
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				t.Cleanup(func() {
					TruncateAll(t)
				})
			})
			var uid int64
			if tc.prepare != nil {
				uid = tc.prepare()
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
				supToken := LoginAndGetToken(t, "supuser@test.com", supPass)
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

func TestCreateStockMov(t *testing.T) {
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
		body         string
		expectedCode int
	}{
		{
			name: "unauthorized (no token)",
			prepare: func() {

			},
			body:         `{}`,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "invalid json",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "MTSup", "mtsup@test", "MtCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "MTCat", "dcdunu")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "MTProd", "ddfuhfu", 11, cid, sid)
				pid, _ := pr.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "MTSupUser", supEmail, hash, "supplier_admin", sid)

				_ = pid

			},
			body:         `{`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden for non-supplier (system_admin)",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "FbdSup", "fbdsup@test", "CoF")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "FbdCat", "dfund")
				cid, _ := rc.LastInsertId()
				_, _ = config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "FbdProd", "dufucnsn", 10, cid, sid)
			},
			body:         `{"product_id":1,"quantity":5,"movement_type":"IN"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "invalid movement type",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "MTSup", "mtsup@test", "MtCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "MTCat", "dcdunu")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "MTProd", "ddfuhfu", 11, cid, sid)
				pid, _ := pr.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "MTSupUser", supEmail, hash, "supplier_admin", sid)

				_ = pid
			},
			body:         `{"product_id":1,"quantity":5,"movement_type":"INDHDHDH"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid quantity (<=0)",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "MTSup", "mtsup@test", "MtCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "MTCat", "dcdunu")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "MTProd", "ddfuhfu", 11, cid, sid)
				pid, _ := pr.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "MTSupUser", supEmail, hash, "supplier_admin", sid)

				_ = pid
			},
			body:         `{"product_id":1,"quantity":-1,"movement_type":"IN"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "product not found",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "MTSup", "mtsup@test", "MtCo")
				sid, _ := resS.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "MTSupUser", supEmail, hash, "supplier_admin", sid)
			},
			body:         `{"product_id":999999,"quantity":25,"movement_type":"IN"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "forbidden supplier (cannot create stock movements for other products)",
			prepare: func() {
				res1, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "MTSup", "mtsup@test", "MtCo")
				s1, _ := res1.LastInsertId()
				res2, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "MTSupdjsjdj", "mtsup@test.com", "MtCodsjsjs")
				s2, _ := res2.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "MTCat", "dcdunu")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "MTProd", "ddfuhfu", 11, cid, s1)
				pid, _ := pr.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "TSupdjsjdj", supEmail, hash, "supplier_admin", s2)
				_ = pid

			},
			body:         `{"product_id":1,"quantity":25,"movement_type":"IN"}`,
			expectedCode: http.StatusForbidden,
		},
		{
			name: "OUT causes low stock warning (below threshold)",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "MTSup", "mtsup@test", "MtCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "MTCat", "dcdunu")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "MTProd", "ddfuhfu", 11, cid, sid)
				pid, _ := pr.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "MTSupUser", supEmail, hash, "supplier_admin", sid)

				_ = pid
			},
			body:         `{"product_id":1,"quantity":5,"movement_type":"OUT"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "successful IN by supplier_admin",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "MTSup", "mtsup@test", "MtCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "MTCat", "dcdunu")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "MTProd", "ddfuhfu", 11, cid, sid)
				pid, _ := pr.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "MTSupUser", supEmail, hash, "supplier_admin", sid)

				_ = pid
			},
			body:         `{"product_id":1,"quantity":25,"movement_type":"IN"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "successful OUT by supplier_admin",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "MTSup", "mtsup@test", "MtCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "MTCat", "dcdunu")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "MTProd", "ddfuhfu", 11, cid, sid)
				pid, _ := pr.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "MTSupUser", supEmail, hash, "supplier_admin", sid)
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid, 20, "IN", 1)

			},
			body:         `{"product_id":1,"quantity":5,"movement_type":"OUT","reason":"selling"}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "insert FK error (performed_by missing)",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "FKSup", "fksup@test", "FKCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "FKCat", "fkdesc")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "FKProd", "fkprod", 10, cid, sid)
				_, _ = pr.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				resU, _ := config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "FKUser", supEmail, hash, "supplier_admin", sid)
				_, _ = resU.LastInsertId()

			},
			body:         `{"product_id":1,"quantity":50,"movement_type":"IN"}`,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
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
				supToken := LoginAndGetToken(t, supEmail, supPass)
				if supToken != "" {
					req.Header.Set("Authorization", "Bearer "+supToken)
				}
				_, _ = config.DB.Exec("delete from users where email = ?", supEmail)

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
		body         string
		expectedCode int
	}{
		{
			name: "invalid product_id param",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "tester", "tester@test.com", "TestCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatSys", "desc")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)",
					"SysProd", "descj", 10, cid, sid)
				pid, _ := pr.LastInsertId()
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid, 5, "IN", 1)
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid, 2, "OUT", 1)

			},
			query:        "?product_id=abc",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "unauthorized (no token)",
			prepare: func() {
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", 1, 10, "IN", 1)
			},
			query:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "system_admin sees all movements",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "tester", "tester@test.com", "TestCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatSys", "desc")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)",
					"SysProd", "descj", 10, cid, sid)
				pid, _ := pr.LastInsertId()
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid, 5, "IN", 1)
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid, 2, "OUT", 1)
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin filters by product_id",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "tester", "tester@test.com", "TestCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "CatSys", "desc")
				cid, _ := rc.LastInsertId()
				prd1, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)",
					"SysProd", "descj", 10, cid, sid)
				pid1, _ := prd1.LastInsertId()
				prd2, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)",
					"SysProdsasa", "descjas", 10, cid, sid)
				pid2, _ := prd2.LastInsertId()
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid1, 5, "IN", 1)
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid2, 2, "OUT", 1)
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin no movements returns empty list",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "HUhusup", "nmsup@test", "NmCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "NmCat", "nm")
				cid, _ := rc.LastInsertId()
				_, _ = config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "NmProd", "nmd", 10, cid, sid)
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin sees company movements",
			prepare: func() {
				res1, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SupC1", "supc1@test.com", "SameCo")
				s1, _ := res1.LastInsertId()
				res2, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SupC2", "supc2@test.com", "SameCo")
				s2, _ := res2.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "Ccat", "desc")
				cid, _ := rc.LastInsertId()
				pr1, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)",
					"Pssind", "dsdudun", 10, cid, s1)
				pid1, _ := pr1.LastInsertId()
				pr2, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)",
					"Psidjsidj", "dicnicn", 20, cid, s2)
				pid2, _ := pr2.LastInsertId()
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid1, 7, "IN", 1)
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid2, 4, "IN", 1)
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)",
					"CompanyAdmin", supEmail, hash, "supplier_admin", s1)
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin filter by product_id",
			prepare: func() {
				res1, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SupC1", "supc1@test.com", "SameCo")
				s1, _ := res1.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "Ccat", "desc")
				cid, _ := rc.LastInsertId()
				pr1, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)",
					"Pssind", "dsdudun", 10, cid, s1)
				pid1, _ := pr1.LastInsertId()
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid1, 7, "IN", 1)
				res2, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "SupC2", "supc2@test.com", "SameCo")
				s2, _ := res2.LastInsertId()
				pr2, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)",
					"Psidjsidj", "dicnicn", 20, cid, s2)
				pid2, _ := pr2.LastInsertId()
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid2, 4, "IN", 1)
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)",
					"coadmin", supEmail, hash, "supplier_admin", s1)

			},
			query:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin no movements returns empty list",
			prepare: func() {

				res, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "fakeSup", "fake@test.com", "fakeCo")
				sid, _ := res.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)",
					"fakeAdmin", supEmail, hash, "supplier_admin", sid)
			},
			query:        "",
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			t.Cleanup(func() {
				TruncateAll(t)
			})
			if tc.prepare != nil {
				tc.prepare()
			}
			query := tc.query
			if tc.name == "system_admin filter by product_id" {
				var pid int64
				err := config.DB.QueryRow("select product_id from products order by product_id asc limit 1").Scan(&pid)
				if err != nil {
					t.Fatalf("failed to fetch product id for test %s: %v", tc.name, err)
				}
				query = fmt.Sprintf("?product_id=%d", pid)
			}
			if tc.name == "supplier_admin filter by product_id" {
				var pid int64
				err := config.DB.QueryRow("select p.product_id from products p join suppliers s on p.product_supplier_id = s.supplier_id where s.company = (select company from suppliers where supplier_id = (select supplier_id from users where email = ? limit 1)) order by p.product_id asc limit 1", supEmail).Scan(&pid)
				if err != nil {
					_ = config.DB.QueryRow("select product_id from products order by product_id asc limit 1").Scan(&pid)
				}
				query = fmt.Sprintf("?product_id=%d", pid)
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
			name: "unauthorized (no token)",
			prepare: func() {

			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "system_admin empty counts",
			prepare: func() {

			},
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin dashboard (low stock products)",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "LowSup", "lowsup@test.com", "LowCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "LowCat", "desc")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "LowProd", "desc", 10, cid, sid)
				pid, _ := pr.LastInsertId()
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid, 5, "IN", 1)
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid, 100, "OUT", 1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "system_admin no low stock products (all above threshold)",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "HighSup", "highsup@test.com", "HighCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "HighCat", "desc")
				cid, _ := rc.LastInsertId()
				pr, _ := config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "HighProd", "desc", 10, cid, sid)
				pid, _ := pr.LastInsertId()
				_, _ = config.DB.Exec("insert into stock_movements(product_id,quantity,movement_type,performed_by) values(?,?,?,?)", pid, 100, "IN", 1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin dashboard (has products)",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "LowSup", "lowsup@test.com", "LowCo")
				sid, _ := resS.LastInsertId()
				rc, _ := config.DB.Exec("insert into categories(category_name,category_description) values(?,?)", "LowCat", "desc")
				cid, _ := rc.LastInsertId()
				_, _ = config.DB.Exec("insert into products(product_name,product_description,product_cost,product_category_id,product_supplier_id) values(?,?,?,?,?)", "LowProd", "desc", 10, cid, sid)
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "CompanyAdmin", supEmail, hash, "supplier_admin", sid)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "supplier_admin dashboard (no products)",
			prepare: func() {
				resS, _ := config.DB.Exec("insert into suppliers(name,email,company) values(?,?,?)", "LowSup", "lowsup@test.com", "LowCo")
				sid, _ := resS.LastInsertId()
				hash, _ := utils.HashPwd(supPass)
				_, _ = config.DB.Exec("insert into users(name,email,password_hash,role,supplier_id) values(?,?,?,?,?)", "CompanyAdmin", supEmail, hash, "supplier_admin", sid)
			},
			expectedCode: http.StatusOK,
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				TruncateAll(t)
			})
			if tc.prepare != nil {
				tc.prepare()
			}
			req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
			req.Header.Set("Content-Type", "application/json")
			switch tc.name {
			case "unauthorized (no token)":
				req.Header.Set("Authorization", "Bearer "+"")
			case "system_admin empty counts", "system_admin dashboard (low stock products)":
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
