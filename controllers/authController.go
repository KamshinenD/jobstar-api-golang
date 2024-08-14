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
	user.IsAdmin = false

	err = user.Save()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "Registration Successful", nil)
}

func LoginController(c *gin.Context) {
	var user models.UserLogin
	err := c.ShouldBindJSON(&user)

	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid JSON", err)
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

	err = user.ValidateCredentials()

	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Invalid credentials", err)
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID) // we are able to access user.ID cos it has been binded in ValidateCredentials

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Could not authenticate User", err)
		return
	}
	// context.JSON(http.StatusOK, gin.H{"message":"Login Successful", "token": token})
	utils.RespondJSON(c, http.StatusOK, "Login Successful", gin.H{
		"token": token,
	})
}

func UpdateUser(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "userId does not exist", nil)
		return
	}

	// Type assertion for userId
	userIdStr, ok := userId.(string)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "userId is not a string", nil)
		return
	}

	var user models.UserUpdate
	err := c.ShouldBindJSON(&user)

	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	user.ID = userIdStr

	// Validate required fields
	if user.FirstName == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide first name", nil)
		return
	}
	if user.LastName == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide last name", nil)
		return
	}
	if user.Location == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide first name", nil)
		return
	}

	err = user.Update()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Unable to update user details", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "User details updated successfully", nil)
}
