package entity

type RegisterResponse struct {
	Message string `json:"message" example:"User registered successfully"`
}

type LoginResponse struct {
	Token    string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Role     string `json:"role" example:"user"`
	Username string `json:"username" example:"user123"`
	UserID   int    `json:"userID" example:"1"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
