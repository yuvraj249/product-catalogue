package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/interface/http/delivery"
	"backend/internal/interface/http/middleware"
	"backend/internal/interface/storage/postgres"
	"backend/internal/interface/storage/redis"
	"backend/internal/usecase"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Idempotency-Key")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.LoadConfig()

	// 1. Initialize Storage Layer (Postgres / Embedded SQLite + Redis)
	db, err := postgres.InitDB(cfg.DBDriver, cfg.DSN)
	if err != nil {
		log.Fatalf("Fatal: Database initialization error: %v", err)
	}
	defer db.Close()

	redisClient := redis.NewRedisClient(cfg.RedisAddr, cfg.RedisPwd)

	// 2. Repositories
	authRepo := postgres.NewAuthRepository(db)
	inventoryRepo := postgres.NewInventoryRepository(db)
	invoiceRepo := postgres.NewInvoiceRepository(db)

	// 3. Usecases
	authUsecase := usecase.NewAuthUsecase(authRepo, cfg.JWTSecret)
	inventoryUsecase := usecase.NewInventoryUsecase(inventoryRepo)
	paymentUsecase := usecase.NewPaymentUsecase(invoiceRepo, inventoryRepo)

	// 4. Handlers & Middleware
	authHandler := delivery.NewAuthHandler(authUsecase)
	inventoryHandler := delivery.NewInventoryHandler(inventoryUsecase)
	invoiceHandler := delivery.NewInvoiceHandler(invoiceRepo)
	webhookHandler := delivery.NewWebhookHandler(paymentUsecase)

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	idempotencyMiddleware := middleware.IdempotencyMiddleware(redisClient)

	// 5. Router Setup
	mux := http.NewServeMux()

	// Public Auth Endpoints
	mux.HandleFunc("POST /api/v1/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/v1/webhooks/stripe", webhookHandler.HandleStripeWebhook)

	// Protected ERP Routes
	apiMux := http.NewServeMux()

	// Inventory
	apiMux.Handle("POST /api/v1/inventory/items", middleware.RequireRole("SuperAdmin", "TenantAdmin", "WarehouseManager")(http.HandlerFunc(inventoryHandler.CreateItem)))
	apiMux.Handle("GET /api/v1/inventory/items", http.HandlerFunc(inventoryHandler.ListItems))
	apiMux.Handle("GET /api/v1/inventory/barcode", http.HandlerFunc(inventoryHandler.GetItemByBarcode))
	apiMux.Handle("POST /api/v1/inventory/transfer", middleware.RequireRole("SuperAdmin", "TenantAdmin", "WarehouseManager")(http.HandlerFunc(inventoryHandler.TransferStock)))
	apiMux.Handle("POST /api/v1/inventory/scrap", middleware.RequireRole("SuperAdmin", "TenantAdmin", "WarehouseManager")(http.HandlerFunc(inventoryHandler.ScrapStock)))

	// Invoices & Billing
	apiMux.Handle("POST /api/v1/invoices", idempotencyMiddleware(middleware.RequireRole("SuperAdmin", "TenantAdmin", "Cashier")(http.HandlerFunc(invoiceHandler.CreateInvoice))))

	// Attach Auth Middleware to API routes
	protectedHandler := authMiddleware(apiMux)

	mux.Handle("/api/v1/inventory/", protectedHandler)
	mux.Handle("/api/v1/invoices", protectedHandler)

	// Apply CORS
	finalHandler := corsMiddleware(mux)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      finalHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("🏛️ Enterprise Multi-Tenant ERP Backend running on port %s...\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server ListenAndServe error: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down ERP Backend server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown error: %v", err)
	}

	fmt.Println("ERP Backend server stopped cleanly.")
}
