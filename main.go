// @JOBSTAR APP API
// @version 1.0
// @description This is an API for managing and tracking jobs.

// @host localhost:8080
// @BasePath /api/v1

package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files" // This import provides the swaggerFiles variable
	ginSwagger "github.com/swaggo/gin-swagger"
	"jobstar.com/api/db"
	_ "jobstar.com/api/docs" // This import is required to include the generated docs
	"jobstar.com/api/routes"
)

func main() {
	db.InitDB()
	godotenv.Load()

	APIHostURL := os.Getenv("APIHostURL")

	if APIHostURL == "" {
		APIHostURL = "localhost:8080" // Default value if not set
	}

	server := gin.Default()

	// Swagger setup
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Auth Routes
	authRoutes := server.Group("/api/v1/auth")
	{
		routes.RegisterAuthRoutes(authRoutes)
	}

	// Job Routes
	jobRoutes := server.Group("/api/v1/jobs")
	{
		routes.RegisterJobRoutes(jobRoutes)
	}

	// server.Run(":8080")
	server.Run(APIHostURL)
}
