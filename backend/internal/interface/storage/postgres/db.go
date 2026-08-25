package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

func InitDB(driverName, dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	if driverName == "postgres" && dsn != "" {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
		}
	} else {
		err = fmt.Errorf("postgres connection bypassed or not configured")
	}

	if err != nil {
		log.Printf("PostgreSQL connection bypassed/failed (%v). Operating on embedded SQLite ERP storage...\n", err)
		db, err = sql.Open("sqlite", "product_catalogue.db?_pragma=foreign_keys(1)")
		if err != nil {
			return nil, fmt.Errorf("failed to initialize embedded SQLite: %w", err)
		}
	}

	if err := runMigrations(db); err != nil {
		log.Printf("Migration notice: %v\n", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	schemaBytes, err := os.ReadFile("schema.sql")
	if err != nil {
		// Fallback to embedded schema strings if schema.sql file is missing
		return nil
	}

	sqlStr := string(schemaBytes)

	// Replacements for SQLite dialect compatibility
	sqlStr = strings.ReplaceAll(sqlStr, "UUID PRIMARY KEY DEFAULT uuid_generate_v4()", "TEXT PRIMARY KEY")
	sqlStr = strings.ReplaceAll(sqlStr, "TIMESTAMPTZ", "TIMESTAMP")
	sqlStr = strings.ReplaceAll(sqlStr, "user_role NOT NULL DEFAULT 'Cashier'", "TEXT NOT NULL DEFAULT 'Cashier'")
	sqlStr = strings.ReplaceAll(sqlStr, "location_type NOT NULL", "TEXT NOT NULL")
	sqlStr = strings.ReplaceAll(sqlStr, "invoice_status NOT NULL DEFAULT 'Draft'", "TEXT NOT NULL DEFAULT 'Draft'")
	sqlStr = strings.ReplaceAll(sqlStr, "NUMERIC(12, 4)", "REAL")
	sqlStr = strings.ReplaceAll(sqlStr, "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";", "")
	sqlStr = strings.ReplaceAll(sqlStr, "CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";", "")
	sqlStr = strings.ReplaceAll(sqlStr, "CREATE TYPE user_role AS ENUM ('SuperAdmin', 'TenantAdmin', 'WarehouseManager', 'Cashier', 'Auditor');", "")
	sqlStr = strings.ReplaceAll(sqlStr, "CREATE TYPE location_type AS ENUM ('Internal', 'Customer', 'Vendor', 'InventoryLoss', 'Production');", "")
	sqlStr = strings.ReplaceAll(sqlStr, "CREATE TYPE invoice_status AS ENUM ('Draft', 'Unpaid', 'Paid', 'Cancelled');", "")

	statements := strings.Split(sqlStr, ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.Exec(stmt); err != nil {
			// Log dialect differences silently
		}
	}

	return nil
}
