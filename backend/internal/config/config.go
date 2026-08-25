package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBDriver     string
	DSN          string
	JWTSecret    []byte
	RedisAddr    string
	RedisPwd     string
	StripeSecret string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DSN")
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		if dsn != "" {
			driver = "postgres"
		} else {
			driver = "sqlite"
		}
	}

	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		secret = "super_secret_enterprise_erp_key_2026!@#"
		log.Println("WARNING: Using default JWT secret key. Set JWT_SECRET_KEY in production.")
	}

	return &Config{
		Port:         port,
		DBDriver:     driver,
		DSN:          dsn,
		JWTSecret:    []byte(secret),
		RedisAddr:    os.Getenv("REDIS_ADDR"),
		RedisPwd:     os.Getenv("REDIS_PASSWORD"),
		StripeSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
	}
}
