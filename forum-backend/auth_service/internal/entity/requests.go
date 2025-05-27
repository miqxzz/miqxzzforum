package entity

type RegisterRequest struct {
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"P@ssw0rd"`
	Role     string `json:"role" example:"user"`
}

type LoginRequest struct {
	Username string `json:"username" example:"user123"`
	Password string `json:"password" example:"P@ssw0rd"`
}
