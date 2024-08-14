package routes

import (
	"github.com/gin-gonic/gin"
	"jobstar.com/api/controllers"
	"jobstar.com/api/middlewares"
)

// func RegisterRoutes(server *gin.Engine) {
// 	server.POST("/register", controllers.RegistrationController)
// 	server.POST("/login", controllers.LoginController)
// 	server.PATCH("/updateUser", middlewares.Authenticate, controllers.UpdateUser)
// }

func RegisterAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", controllers.RegistrationController)
	router.POST("/login", controllers.LoginController)
	router.GET("/verifyAccount", controllers.VerifyAccountController)
	router.PATCH("/updateUser", middlewares.Authenticate, controllers.UpdateUser)
}

func RegisterJobRoutes(router *gin.RouterGroup) {
	router.POST("/", middlewares.Authenticate, controllers.CreateJob)
	router.GET("/", middlewares.Authenticate, controllers.GetJobsByUser)
	router.GET("/stats", middlewares.Authenticate, controllers.ShowStats)
	router.GET("/:id", middlewares.Authenticate, controllers.GetSingleJob)
	router.DELETE("/:id", middlewares.Authenticate, controllers.DeleteJob)
	router.PATCH("/:id", middlewares.Authenticate, controllers.UpdateJob)
}
