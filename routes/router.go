package routes

import (
	"github.com/KishiEdward/back-skenaid/handlers"
	"github.com/KishiEdward/back-skenaid/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authHandler := handlers.NewAuthHandler()
	productHandler := handlers.NewProductHandler()

	v1 := r.Group("/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "skenaid-backend",
		})
	})

	auth := v1.Group("/auth")
	auth.POST("/verify-token", authHandler.VerifyToken)

	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware())

	products := protected.Group("/products")
	
	products.GET("", productHandler.GetAll)
	products.GET("/:id", productHandler.GetByID)

	adminProducts := products.Group("")
	adminProducts.Use(middleware.AdminOnly())

	adminProducts.POST("", productHandler.Create)
	adminProducts.PUT("/:id", productHandler.Update)
	adminProducts.DELETE("/:id", productHandler.Delete)

    cartHandler := handlers.NewCartHandler()
    orderHandler := handlers.NewOrderHandler()

    cart := protected.Group("/cart")
    cart.GET("", cartHandler.GetCart)
    cart.POST("", cartHandler.AddToCart)
    cart.PUT("/:id", cartHandler.UpdateCartItem)
    cart.DELETE("/:id", cartHandler.RemoveCartItem)
    cart.DELETE("", cartHandler.ClearCart)

    orders := protected.Group("/orders")
    orders.GET("", orderHandler.GetMyOrders)
    orders.GET("/:id", orderHandler.GetOrderDetail)
    orders.POST("/checkout", orderHandler.Checkout)

	return r
}