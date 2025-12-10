package main

import (
	"fmt"
	"log"
	"os"
	"product-catalogue/config"
	"product-catalogue/routes"
)

func main() {
	config.JoinDB()

	router := routes.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server running on port %s\n", port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
