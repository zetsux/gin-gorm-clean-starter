package route

import (
	"github.com/zetsux/gin-gorm-template-clean/api/v1/controller"
	"github.com/zetsux/gin-gorm-template-clean/common/middleware"
	"github.com/zetsux/gin-gorm-template-clean/common/standard"
	"github.com/zetsux/gin-gorm-template-clean/internal/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, userC controller.UserController, jwtS service.JWTService) {
	userRoutes := router.Group("/api/v1/users")
	{
		// admin routes
		userRoutes.GET("", middleware.Authenticate(jwtS, standard.EnumRoleAdmin), userC.GetAllUsers)
		userRoutes.PATCH("/:user_id", middleware.Authenticate(jwtS, standard.EnumRoleAdmin), userC.UpdateUserByID)
		userRoutes.DELETE("/:user_id", middleware.Authenticate(jwtS, standard.EnumRoleAdmin), userC.DeleteUserByID)

		// user routes
		userRoutes.GET("/me", middleware.Authenticate(jwtS, standard.EnumRoleUser), userC.GetMe)
		userRoutes.PATCH("/me/name", middleware.Authenticate(jwtS, standard.EnumRoleUser), userC.UpdateSelfName)
		userRoutes.DELETE("/me", middleware.Authenticate(jwtS, standard.EnumRoleUser), userC.DeleteSelfUser)
		userRoutes.POST("", userC.Register)
		userRoutes.POST("/login", userC.Login)
		userRoutes.PATCH("/picture", middleware.Authenticate(jwtS, standard.EnumRoleUser), userC.ChangePicture)
		userRoutes.DELETE("/picture/:user_id", middleware.Authenticate(jwtS, standard.EnumRoleUser), userC.DeletePicture)
	}
}
