package routes

import (
	"github.com/Alfeus22/ecommerce-go/controllers"
	admin "github.com/Alfeus22/ecommerce-go/controllers/admin"
	"github.com/Alfeus22/ecommerce-go/controllers/seller"
	user "github.com/Alfeus22/ecommerce-go/controllers/user"
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
	userRoutes.Use(auth.RoleMiddleware("customer"))
	{
		userRoutes.POST("/transaction")
		userRoutes.GET("/product/:id")
		userRoutes.GET("/order")
		userRoutes.PUT("/toSeller", user.UpdateToSeller)

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
	adminRoutes.Use(auth.RoleMiddleware("admin"))
	{
		adminRoutes.GET("/showSeller", admin.GetSeller)
		adminRoutes.GET("/pending-seller", admin.GetPendingSeller)
		adminRoutes.DELETE("/deleteSeller/:id", admin.DeleteSeller)
		adminRoutes.PUT("/updateUser/:id", admin.UpdateUser)
	}

}
