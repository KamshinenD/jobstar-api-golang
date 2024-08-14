package main

import (
	"github.com/gin-gonic/gin"
	"jobstar.com/api/db"
	"jobstar.com/api/routes"
)

func main() {
	db.InitDB()

	server := gin.Default()
	// routes.RegisterRoutes(server)

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

	server.Run(":8080")
}
