package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mikemeh/ecommerce-api/config"
	"github.com/mikemeh/ecommerce-api/controllers"
	docs "github.com/mikemeh/ecommerce-api/docs"
	"github.com/mikemeh/ecommerce-api/middleware"
	"github.com/mikemeh/ecommerce-api/models"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title E-commerce API
// @version 1.0
// @description This is a sample e-commerce API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Initialize database
	db, err := config.SetupDB()
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{})

	// Initialize router
	r := gin.Default()

	// Initialize swagger docs
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Initialize controllers
	userController := controllers.NewUserController(db)
	productController := controllers.NewProductController(db)
	orderController := controllers.NewOrderController(db)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// User routes
		v1.POST("/register", userController.Register)
		v1.POST("/login", userController.Login)

		// Product routes
		products := v1.Group("/products")
		products.Use(middleware.AuthMiddleware())
		{
			products.POST("/", productController.CreateProduct)
			products.GET("/", productController.GetProducts)
			products.GET("/:id", productController.GetProduct)
			products.PUT("/:id", productController.UpdateProduct)
			products.DELETE("/:id", productController.DeleteProduct)
		}

		// Order routes
		orders := v1.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.POST("/", orderController.CreateOrder)
			orders.GET("/", orderController.GetOrders)
			orders.GET("/:id", orderController.GetOrder)
			orders.PUT("/:id/cancel", orderController.CancelOrder)
			orders.PUT("/:id/status", middleware.AuthMiddleware(), orderController.UpdateOrderStatus)
		}
	}

	// Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
