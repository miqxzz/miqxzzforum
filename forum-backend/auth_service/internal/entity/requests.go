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

type UpdateRoleRequest struct {
	UserID  int    `json:"user_id" example:"1"`
	NewRole string `json:"new_role" example:"admin"`
}
