package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"jobstar.com/api/utils"
)

func Authenticate(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")

	if token == "" {
		utils.RespondError(c, http.StatusUnauthorized, "Not Authorized", nil)
		c.Abort()
		return
	}

	userId, err := utils.VerifyToken(token)

	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "Something went wrong", err)
		c.Abort()
		return
	}

	c.Set("userId", userId)
	c.Next()
}
