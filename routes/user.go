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
		// admin routes
		userRoutes.GET("", middleware.Authenticate(jwtS, "admin"), userC.GetAllUsers)
		userRoutes.PATCH("/:user_id", middleware.Authenticate(jwtS, "admin"), userC.UpdateUserById)
		userRoutes.DELETE("/:user_id", middleware.Authenticate(jwtS, "admin"), userC.DeleteUserById)

		// user routes
		userRoutes.GET("/me", middleware.Authenticate(jwtS, "user"), userC.GetMe)
		userRoutes.PATCH("/me/name", middleware.Authenticate(jwtS, "user"), userC.UpdateSelfName)
		userRoutes.DELETE("/me", middleware.Authenticate(jwtS, "user"), userC.DeleteSelfUser)
		userRoutes.POST("", userC.Register)
		userRoutes.POST("/login", userC.Login)
	}
}
