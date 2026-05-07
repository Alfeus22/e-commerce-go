package routes

import (
	"github.com/Alfeus22/ecommerce-go/controllers"
	admin "github.com/Alfeus22/ecommerce-go/controllers/admin"
	"github.com/Alfeus22/ecommerce-go/controllers/seller"
	auth "github.com/Alfeus22/ecommerce-go/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/home")

	// AUTH user
	userRoutes := r.Group("/user")

	userRoutes.Use(auth.AuthMiddleware())
	{
		userRoutes.POST("/transaction")
		userRoutes.GET("/product/:id")
		userRoutes.GET("/order")
	}

	// AUTH SELLER
	sellerRoutes := r.Group("/seller")

	sellerRoutes.Use(auth.AuthMiddleware())
	sellerRoutes.Use(auth.RoleMiddleware("SELLER"))
	{
		sellerRoutes.POST("/product", seller.MakeProduct)
		sellerRoutes.DELETE("/product/:id")
		sellerRoutes.PUT("/product/:id")
	}

	// AUTH ADMIN
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(auth.AuthMiddleware())
	adminRoutes.Use(auth.RoleMiddleware("ADMIN"))
	{
		adminRoutes.GET("/pending-seller", admin.GetPendingSeller)
		adminRoutes.POST("/deleteSeller/:id")
		adminRoutes.PUT("/updateUser/:id", admin.UpdateUser)
	}

}
