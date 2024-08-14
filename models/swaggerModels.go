package models

// UserLoginRequest is used only for Swagger documentation.
type UserLoginRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"password123"`
}

type ErrorData struct{}

// ErrorResponse represents a standard structure for error responses.
type ErrorResponse struct {
	Code    int       `json:"response_code"` // HTTP status code
	Message string    `json:"message"`       // Error message
	Data    ErrorData `json:"data"`
}
type SuccessResponse struct {
	Code    int    `json:"response_code"` // HTTP status code
	Message string `json:"message"`       // Error message
	Data    any    `json:"data"`
}

type Data struct {
	Token string `json:"token"`
}

// AuthResponse represents the structure of the authentication response.
type AuthResponse struct {
	Code    int    `json:"response_code"` // HTTP status code
	Message string `json:"message"`       // Error message
	Data    Data   `json:"data"`
}

type UserRegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Location  string `json:"location"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserUpdateRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Location  string `json:"location"`
}

type JobRequest struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	JobLocation string `json:"jobLocation"`
	Status      string `json:"status"`
	JobType     string `json:"jobType"`
}
