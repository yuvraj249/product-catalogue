package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func JoinDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found!!")
	}

	dsn := os.Getenv("DSN")
	dbType := os.Getenv("DB_TYPE")

	var db *sql.DB
	var dbErr error

	if dbType != "sqlite" && dsn != "" {
		db, dbErr = sql.Open("mysql", dsn)
		if dbErr == nil {
			dbErr = db.Ping()
		}
	} else {
		dbErr = fmt.Errorf("forced sqlite or missing dsn")
	}

	if dbErr != nil {
		log.Printf("MySQL not reachable (%v). Falling back to embedded SQLite database...\n", dbErr)
		db, dbErr = sql.Open("sqlite", "product_catalogue.db?_pragma=foreign_keys(1)")
		if dbErr != nil {
			log.Fatalf("Could not open embedded SQLite database: %v", dbErr)
		}
		log.Println("Successfully connected to embedded SQLite database (product_catalogue.db)")
	} else {
		log.Println("Successfully connected to MySQL database")
	}

	DB = db

	// Convert MySQL DDL dialect to SQLite compatible if running on SQLite
	err = ExecMigration(DB, "migration/init.sql")
	if err != nil {
		log.Printf("Warning on init migration: %v\n", err)
	}

	err = ExecMigration(DB, "migration/002_security.sql")
	if err != nil {
		log.Printf("Warning on security migration: %v\n", err)
	}

	// Ensure default system_admin user exists
	ensureDefaultAdminUser(DB)

	fmt.Println("Database connected as well as migrations executed!!")
}

func ensureDefaultAdminUser(db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = 'yuvrajbisht41@gmail.com'").Scan(&count)
	if err == nil && count == 0 {
		_, _ = db.Exec(
			"INSERT INTO users (name, email, password_hash, role) VALUES (?, ?, ?, ?)",
			"Owner", "yuvrajbisht41@gmail.com", "$2a$10$CCw/Xx/.lW1BCcc0MYIH5.xh2QJq7pBqMrWeE.WPxgRI4F8Af12s2", "system_admin",
		)
		log.Println("Seeded default system_admin user: yuvrajbisht41@gmail.com")
	}
}

func ExecMigration(db *sql.DB, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read path: %w", err)
	}

	sqlStr := string(b)
	// Replace MySQL specific keywords for SQLite compatibility
	sqlStr = strings.ReplaceAll(sqlStr, "auto_increment", "AUTOINCREMENT")
	sqlStr = strings.ReplaceAll(sqlStr, "AUTO_INCREMENT", "AUTOINCREMENT")
	sqlStr = strings.ReplaceAll(sqlStr, "enum('system_admin', 'supplier_admin')", "TEXT")
	sqlStr = strings.ReplaceAll(sqlStr, "enum('percent', 'flat')", "TEXT")
	sqlStr = strings.ReplaceAll(sqlStr, "enum('admin','manager','staff')", "TEXT")
	sqlStr = strings.ReplaceAll(sqlStr, "enum('IN', 'OUT')", "TEXT")
	sqlStr = strings.ReplaceAll(sqlStr, "int auto_increment primary key", "INTEGER PRIMARY KEY AUTOINCREMENT")
	sqlStr = strings.ReplaceAll(sqlStr, "int AUTOINCREMENT primary key", "INTEGER PRIMARY KEY AUTOINCREMENT")
	sqlStr = strings.ReplaceAll(sqlStr, "INT AUTOINCREMENT PRIMARY KEY", "INTEGER PRIMARY KEY AUTOINCREMENT")
	sqlStr = strings.ReplaceAll(sqlStr, "DATETIME DEFAULT CURRENT_TIMESTAMP", "TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
	sqlStr = strings.ReplaceAll(sqlStr, "on delete set null on update cascade", "")
	sqlStr = strings.ReplaceAll(sqlStr, "insert ignore into", "INSERT OR IGNORE INTO")

	sqlcommands := strings.Split(sqlStr, ";")
	for _, s := range sqlcommands {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, err := db.Exec(s); err != nil {
			log.Printf("Migration notice on statement [%s]: %v\n", s, err)
		}
	}
	return nil
}
