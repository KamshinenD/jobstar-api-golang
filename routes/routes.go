package routes

import (
	"github.com/gin-gonic/gin"
	"jobstar.com/api/controllers"
)

func RegisterRoutes(server *gin.Engine) {
	server.POST("/register", controllers.RegistrationController)
}
