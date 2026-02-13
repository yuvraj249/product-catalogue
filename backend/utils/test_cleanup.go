package utils

import (
	"net/http"
	"os"
	"product-catalogue/config"

	"github.com/gin-gonic/gin"
)

func allowOnlyTestEnv(c *gin.Context) bool {
	if os.Getenv("APP_ENV") != "test" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return false
	}
	return true
}

func CleanupUsers(c *gin.Context) {
	if !allowOnlyTestEnv(c) {
		return
	}

	config.DB.Exec("DELETE FROM users WHERE email NOT IN ('yuvrajbisht41@gmail.com')")

	c.JSON(200, gin.H{"message": "Users cleaned"})
}

func CleanupSuppliers(c *gin.Context) {
	if !allowOnlyTestEnv(c) {
		return
	}

	config.DB.Exec("DELETE FROM suppliers")
	c.JSON(200, gin.H{"message": "Suppliers cleaned"})
}

func CleanupCategories(c *gin.Context) {
	if !allowOnlyTestEnv(c) {
		return
	}

	config.DB.Exec("DELETE FROM categories")
	c.JSON(200, gin.H{"message": "Categories cleaned"})
}

func CleanupProducts(c *gin.Context) {
	if !allowOnlyTestEnv(c) {
		return
	}

	config.DB.Exec("DELETE FROM products")
	c.JSON(200, gin.H{"message": "Products cleaned"})
}

func CleanupStock(c *gin.Context) {
	if !allowOnlyTestEnv(c) {
		return
	}

	config.DB.Exec("DELETE FROM stock_movements")
	c.JSON(200, gin.H{"message": "Stocks cleaned"})
}

func CleanupAll(c *gin.Context) {
	if !allowOnlyTestEnv(c) {
		return
	}

	config.DB.Exec("DELETE FROM stock_movements")
	config.DB.Exec("DELETE FROM products")
	config.DB.Exec("DELETE FROM categories")
	config.DB.Exec("DELETE FROM suppliers")
	config.DB.Exec("DELETE FROM users WHERE email NOT IN ('yuvrajbisht41@gmail.com')")

	c.JSON(200, gin.H{"message": "All test data cleaned"})
}
