package routes

import (
	"os"
	"product-catalogue/functions"
	"product-catalogue/middleware"
	"product-catalogue/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.POST("/auth/login", middleware.RateLimiter(), functions.Login)
	r.POST("/auth/refresh", functions.RefreshToken)
	r.POST("/auth/logout", functions.Logout)

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	authorized.POST("/products", middleware.RoleSupplierAdmin(), functions.CreateProduct)
	authorized.GET("/products", functions.GetProduct)
	authorized.GET("/products/:id", functions.GetProductByID)
	authorized.PUT("/products/:id", middleware.RoleSupplierAdmin(), functions.UpdateProduct)
	authorized.DELETE("/products/:id", middleware.RoleSupplierAdmin(), functions.DeleteProduct)

	authorized.POST("/categories", middleware.RoleSystemAdmin(), functions.CreateCategory)
	authorized.GET("/categories", functions.GetCategory)
	authorized.GET("/categories/:id", functions.GetCategoryByID)
	authorized.PUT("categories/:id", middleware.RoleSystemAdmin(), functions.UpdateCategory)
	authorized.DELETE("/categories/:id", middleware.RoleSystemAdmin(), functions.DeleteCategory)

	authorized.POST("/suppliers", middleware.RoleSystemAdmin(), functions.CreateSupplier)
	authorized.GET("/suppliers", functions.GetSupplier)
	authorized.GET("/suppliers/:id", functions.GetSupplierByID)
	authorized.PUT("/suppliers/:id", middleware.RoleSystemAdmin(), functions.UpdateSupplier)
	authorized.DELETE("/suppliers/:id", middleware.RoleSystemAdmin(), functions.DeleteSupplier)

	authorized.POST("/stock_movements", functions.CreateStockMovement)
	authorized.GET("/stock_movements", functions.GetStockMovements)
	authorized.PUT("/stock_movements/:id", middleware.RoleSupplierAdmin(), functions.UpdateStockMovement)
	authorized.DELETE("/stock_movements/:id", middleware.RoleSupplierAdmin(), functions.DeleteStockMovement)

	authorized.POST("/users/supplier-admin", middleware.RoleSystemAdmin(), functions.CreateSuppAdmin)
	authorized.GET("/users/supplier-admin", middleware.RoleSystemAdmin(), functions.GetsuppAdmin)
	authorized.GET("/users/supplier-admin/:id", middleware.RoleSystemAdmin(), functions.GetsuppAdminByID)
	authorized.DELETE("/users/supplier-admin/:id", middleware.RoleSystemAdmin(), functions.DeleteSuppAdmin)
	authorized.PUT("/users/supplier-admin/:id", middleware.RoleSystemAdmin(), functions.UpdateSuppAdmin)
	authorized.GET("/dashboard", functions.GetDashboard)
	authorized.GET("/audit-logs", middleware.RoleSystemAdmin(), functions.GetAuditLogs)

	if os.Getenv("APP_ENV") == "test" {
		test := r.Group("/test/cleanup")
		test.Use(middleware.AuthMiddleware())

		test.DELETE("/users", utils.CleanupUsers)
		test.DELETE("/suppliers", utils.CleanupSuppliers)
		test.DELETE("/categories", utils.CleanupCategories)
		test.DELETE("/products", utils.CleanupProducts)
		test.DELETE("/stock_movements", utils.CleanupStock)
		test.DELETE("/all", utils.CleanupAll)
	}

	return r

}
