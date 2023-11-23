package route

import (
	"github.com/zetsux/gin-gorm-template-clean/api/v1/controller"
	"github.com/zetsux/gin-gorm-template-clean/common"
	"github.com/zetsux/gin-gorm-template-clean/internal/middleware"
	"github.com/zetsux/gin-gorm-template-clean/internal/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userC controller.UserController, jwtS service.JWTService) {
	userRoutes := router.Group("/api/v1/users")
	{
		// admin routes
		userRoutes.GET("", middleware.Authenticate(jwtS, common.ENUM_ROLE_ADMIN), userC.GetAllUsers)
		userRoutes.PATCH("/:user_id", middleware.Authenticate(jwtS, common.ENUM_ROLE_ADMIN), userC.UpdateUserById)
		userRoutes.DELETE("/:user_id", middleware.Authenticate(jwtS, common.ENUM_ROLE_ADMIN), userC.DeleteUserById)

		// user routes
		userRoutes.GET("/me", middleware.Authenticate(jwtS, common.ENUM_ROLE_USER), userC.GetMe)
		userRoutes.PATCH("/me/name", middleware.Authenticate(jwtS, common.ENUM_ROLE_USER), userC.UpdateSelfName)
		userRoutes.DELETE("/me", middleware.Authenticate(jwtS, common.ENUM_ROLE_USER), userC.DeleteSelfUser)
		userRoutes.POST("", userC.Register)
		userRoutes.POST("/login", userC.Login)
	}
}
