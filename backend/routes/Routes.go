package routes

import (
	"os"

	"github.com/arthur-tragante/liven-code-test/controllers"
	"github.com/arthur-tragante/liven-code-test/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, userController *controllers.UserController, addressController *controllers.AddressController) {
	r.POST("/register", userController.RegisterUser)
	r.POST("/login", userController.LoginUser)

	jwtSecret := os.Getenv("JWT_SECRET")
	userGroup := r.Group("/user")
	userGroup.Use(middlewares.AuthMiddleware(jwtSecret))
	{
		userGroup.GET("/", userController.GetUser)
		userGroup.PUT("/", userController.UpdateUser)
		userGroup.DELETE("/", userController.DeleteUser)
		userGroup.POST("/address", addressController.CreateAddress)
		userGroup.GET("/address", addressController.GetAddress)
		userGroup.GET("/address/:id", addressController.GetAddress)
		userGroup.PUT("/address/:id", addressController.UpdateAddress)
		userGroup.DELETE("/address/:id", addressController.DeleteAddress)
	}
}
