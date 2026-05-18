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
	r.GET("/home", controllers.GetAllProduct)
	r.GET("/detailProduct/:id", controllers.DetailProduct)
	r.GET("/search", controllers.SearchProduct)
	r.POST("/api/notification/payment", user.HandleMidtransNotification)

	// AUTH user
	userRoutes := r.Group("/user")

	userRoutes.Use(auth.AuthMiddleware())
	userRoutes.Use(auth.RoleMiddleware("customer"))
	{
		userRoutes.POST("/transaction", user.CreateOrder)
		userRoutes.PUT("/toSeller", user.UpdateToSeller)

	}

	// AUTH SELLER
	sellerRoutes := r.Group("/seller")

	sellerRoutes.Use(auth.AuthMiddleware())
	sellerRoutes.Use(auth.RoleMiddleware("seller"))
	{
		sellerRoutes.POST("/product/createProduct", seller.MakeProduct)
		sellerRoutes.GET("/product", seller.GetProduct)
		sellerRoutes.POST("/product/image", seller.UploadImage)
		sellerRoutes.DELETE("/product/:id", seller.DeleteProduct)
		sellerRoutes.PUT("/product/:id", seller.EditProduct)
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
