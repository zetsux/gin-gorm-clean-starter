package routes

import (
	"github.com/zetsux/gin-gorm-template-clean/controller"
	"github.com/zetsux/gin-gorm-template-clean/middleware"
	"github.com/zetsux/gin-gorm-template-clean/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userC controller.UserController, jwtS service.JWTService) {
	userRoutes := router.Group("/api/v1/users")
	{
		userRoutes.GET("", middleware.Authenticate(jwtS, "admin"), userC.GetAllUsers)
		userRoutes.GET("/me", middleware.Authenticate(jwtS, "user"), userC.GetMe)
		userRoutes.PUT("/name", middleware.Authenticate(jwtS, "user"), userC.UpdateSelfName)
		userRoutes.DELETE("", middleware.Authenticate(jwtS, "user"), userC.DeleteSelfUser)
		userRoutes.POST("", userC.Register)
		userRoutes.POST("/login", userC.Login)
	}
}
