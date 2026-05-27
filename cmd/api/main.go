// @title           RealWorld API
// @version         1.0
// @description     Conduit RealWorld API - A social blogging platform built with Go and Gin
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@realworld.io

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Token {your_jwt_token}" to authenticate

package main

import (
	"log"
	"os"

	"realworld-api/config"
	"realworld-api/models"
	"realworld-api/routes"

	_ "realworld-api/docs"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()
	err := config.DB.AutoMigrate(&models.User{}, &models.Follow{}, &models.Article{}, &models.Tag{}, &models.Comment{}, &models.Favorite{})
	if err != nil {
		log.Fatal("Migration error: ", err)
	}
	router := gin.Default()

	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running at http://localhost:%s\n", port)
	router.Run(":" + port)
}
