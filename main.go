package main

import (
	"github.com/gin-gonic/gin"
	"jobstar.com/api/db"
	"jobstar.com/api/routes"
)

func main() {
	db.InitDB()

	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
