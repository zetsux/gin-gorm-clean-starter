package main

import (
	"fmt"
	"fp-rpl/config"
	"fp-rpl/controller"
	"fp-rpl/middleware"
	"fp-rpl/repository"
	"fp-rpl/routes"
	"fp-rpl/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	// Setting Up Database
	db := config.DBSetup()

	// Setting Up Repositories
	userR := repository.NewUserRepository(db)

	// Setting Up Services
	userS := service.NewUserService(userR)
	jwtS := service.NewJWTService()

	// Setting Up Controllers
	userC := controller.NewUserController(userS, jwtS)

	defer config.DBClose(db)

	// Setting Up Server
	server := gin.Default()
	server.Use(
		middleware.CORSMiddleware(),
	)

	// Setting Up Routes
	routes.UserRoutes(server, userC)

	// Running in localhost:8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}
