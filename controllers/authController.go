package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"jobstar.com/api/models"
	"jobstar.com/api/utils"
)

func RegistrationController(c *gin.Context) {
	var user models.User
	err := c.ShouldBindJSON(&user)

	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	// Validate required fields
	if user.FirstName == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide first name", nil)
		return
	}
	if user.LastName == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide last name", nil)
		return
	}
	if user.Email == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide email", nil)
		return
	}
	if user.Password == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide password", nil)
		return
	}
	if user.Location == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide first name", nil)
		return
	}

	err = user.Save()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "Registration Successful", nil)
}
