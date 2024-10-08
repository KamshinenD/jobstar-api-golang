package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"jobstar.com/api/email"
	"jobstar.com/api/models"
	"jobstar.com/api/utils"
)

func generateVerificationToken() (string, error) {
	// Create a slice to hold 40 random bytes
	bytes := make([]byte, 40)

	// Read random bytes into the slice
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("error generating random bytes: %w", err)
	}

	// Convert the byte slice to a hexadecimal string
	token := hex.EncodeToString(bytes)

	return token, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// @Summary Register a new user
// @Description Registers a new user
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param user body models.UserRegisterRequest true "User Data"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [POST]
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
	token, err := generateVerificationToken()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Please provide first name", err)
	}

	user.VerificationToken = token

	err = user.Save()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	baseURL := os.Getenv("APIHostURL")

	// Create the verification link
	verifyLink := fmt.Sprintf("%s/api/v1/auth/verifyAccount?e=%s&t=%s", baseURL, user.Email, token)

	subject := "Welcome to JobStar!"
	name := user.FirstName
	body := fmt.Sprintf(`
		<p>Thank you for registering with JobStar! We're excited to have you on board.</p>
		<p>Please verify your account by clicking the following link:</p>
		<p><a href="%s">Verify your account</a></p>
	`, verifyLink)

	err = email.SendEmail(user.Email, subject, name, body)

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Unable to send email", err)
		// You can choose to handle the error or ignore it
	}

	utils.RespondJSON(c, http.StatusOK, "Registration Successful", nil)
}

// @Summary Login a user
// @Description Authenticates a user and returns a JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param credentials body models.UserLoginRequest true "User Login Data"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/login [POST]
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

// UpdateUser updates a user's details
// @Summary Update user details
// @Description Updates the authenticated user's details
// @Tags Auth
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param user body models.UserUpdateRequest true "User Update Data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /auth/updateUser [PATCH]
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

// VerifyAccountController verifies a user's email account
// @Summary Verify user's email account
// @Description Verifies a user's email account using the link sent to your email
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param e query string true "User Email"
// @Param t query string true "Verification Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /auth/verifyAccount [GET]
func VerifyAccountController(c *gin.Context) {
	email := c.Query("e")
	token := c.Query("t")

	// Check if the parameters are provided
	if email == "" || token == "" {
		utils.RespondJSON(c, http.StatusOK, "Email and Token are required", nil)
		return
	}

	err := models.Verify(email, token)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Unable to update user details", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "User details verified successfully", nil)
}
