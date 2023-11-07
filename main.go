package main

import (
	"os"

	"github.com/zetsux/gin-gorm-template-clean/config"
	"github.com/zetsux/gin-gorm-template-clean/controller"
	"github.com/zetsux/gin-gorm-template-clean/middleware"
	"github.com/zetsux/gin-gorm-template-clean/repository"
	"github.com/zetsux/gin-gorm-template-clean/routes"
	"github.com/zetsux/gin-gorm-template-clean/service"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func main() {
	var (
		db   *gorm.DB           = config.DBSetup()
		jwtS service.JWTService = service.NewJWTService()

		userR repository.UserRepository = repository.NewUserRepository(db)
		userS service.UserService       = service.NewUserService(userR)
		userC controller.UserController = controller.NewUserController(userS, jwtS)
	)

	defer config.DBClose(db)

	// Setting Up Server
	server := gin.Default()
	server.Use(
		middleware.CORSMiddleware(),
	)

	// Setting Up Routes
	routes.UserRoutes(server, userC, jwtS)

	// Running in localhost:8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}
