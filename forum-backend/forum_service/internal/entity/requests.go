package entity

type WSAuthRequest struct {
	Token  string `form:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	UserID string `form:"userID" binding:"required" example:"123"`
	Auth   string `form:"auth" binding:"required" example:"true"`
}
