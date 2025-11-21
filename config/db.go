package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func JoinDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found!!")
	}
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN not found")
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to db: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not reach to db: %v", err)
	}
	DB = db
	err = ExecMigration(DB, "migration/init.sql")
	if err != nil {
		log.Fatalf("Error while performing migration: %v", err)
	}
	fmt.Println("Database connected as well as migration executed!!")
}

func ExecMigration(db *sql.DB, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read path: %w", err)
	}
	sqlcommands := strings.Split(string(b), ";")
	for _, s := range sqlcommands {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, err := db.Exec(s); err != nil {
			return fmt.Errorf("migration failed: %w due to %s", err, s)
		}
	}
	return nil

}
