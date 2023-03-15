package routes

import (
	"fp-rpl/controller"
	"fp-rpl/middleware"
	"fp-rpl/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userC controller.UserController) {
	userRoutes := router.Group("/api/v1/user")
	{
		userRoutes.GET("/", middleware.Authenticate(service.NewJWTService(), "admin"), userC.GetAllUsers)
		userRoutes.GET("/:username", middleware.Authenticate(service.NewJWTService(), "user"), userC.GetUserByUsername)
		userRoutes.PUT("/name", middleware.Authenticate(service.NewJWTService(), "user"), userC.UpdateSelfName)
		userRoutes.DELETE("/", middleware.Authenticate(service.NewJWTService(), "user"), userC.DeleteSelfUser)
		userRoutes.POST("/", userC.Register)
		userRoutes.POST("/login", userC.Login)
	}
}
