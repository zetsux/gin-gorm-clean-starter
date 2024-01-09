package main

import (
	"log"
	"os"

	"github.com/zetsux/gin-gorm-template-clean/api/v1/controller"
	"github.com/zetsux/gin-gorm-template-clean/api/v1/route"
	"github.com/zetsux/gin-gorm-template-clean/config"
	"github.com/zetsux/gin-gorm-template-clean/internal/middleware"
	"github.com/zetsux/gin-gorm-template-clean/internal/repository"
	"github.com/zetsux/gin-gorm-template-clean/internal/service"
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

		fileC controller.FileController = controller.NewFileController()
	)

	defer config.DBClose(db)

	// Setting Up Server
	server := gin.Default()
	server.Use(
		middleware.CORSMiddleware(),
	)

	// Setting Up Routes
	route.UserRoutes(server, userC, jwtS)
	route.FileRoutes(server, fileC)

	// Running in localhost:8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := server.Run(":" + port)
	if err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
