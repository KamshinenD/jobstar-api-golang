// utils/response.go
package utils

import (
	"github.com/gin-gonic/gin"
)

// RespondJSON formats and sends a JSON response.
func RespondJSON(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"message":       message,
		"response_code": status,
		"data":          data,
	})
}

// RespondError sends a JSON response for errors.
func RespondError(c *gin.Context, status int, message string, err error) {

	var errorMessage string

	if err != nil {
		errorMessage = err.Error()
	} else {
		errorMessage = ""
	}

	c.JSON(status, gin.H{
		"message":       message,
		"response_code": status,
		"error":         errorMessage,
	})
}

// usage
// func RegistrationController(c *gin.Context) {
// 	// Assuming some data is retrieved or processed here
// 	data := map[string]string{
// 		"username": "example_user",
// 	}

// func ExampleErrorController(c *gin.Context) {
// 	// Simulate an error
// 	err := errors.New("something went wrong")

// 	// Use the RespondError function to return the error
// 	utils.RespondError(c, http.StatusBadRequest, "Failed to process request", err)
// }
