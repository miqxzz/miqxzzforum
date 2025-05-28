package http

import (
	"net/http"

	utils "github.com/miqxzz/commonmiqx"
	entity "github.com/miqxzz/miqxzzforum/auth_service/internal/entity"
	usecase "github.com/miqxzz/miqxzzforum/auth_service/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	jwtUtil     *utils.JWTUtil
	logger      *zap.Logger
}

func NewAuthHandler(authUsecase usecase.AuthUsecase, jwtUtil *utils.JWTUtil, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, jwtUtil: jwtUtil, logger: logger}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body entity.RegisterRequest true "Данные для регистрации"
// @Success 200 {object} entity.RegisterResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req entity.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON for registration", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.authUsecase.Register(req.Username, req.Password, req.Role); err != nil {
		h.logger.Error("Failed to register user", zap.Error(err), zap.String("username", req.Username))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("User registered successfully", zap.String("username", req.Username))
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Вход пользователя в систему и получение токена
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Учетные данные пользователя"
// @Success 200 {object} entity.LoginResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON for login", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.authUsecase.Login(req.Username, req.Password)
	if err != nil {
		h.logger.Error("Failed to login user", zap.Error(err), zap.String("username", req.Username))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	role, err := h.authUsecase.GetUserRole(req.Username)
	if err != nil {
		h.logger.Error("Failed to get user role", zap.Error(err), zap.String("username", req.Username))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userId, err := h.jwtUtil.GetUserIDFromToken(token)
	if err != nil {
		h.logger.Error("Failed to get user ID from token", zap.Error(err), zap.String("username", req.Username))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("User logged in successfully", zap.String("username", req.Username))
	c.JSON(http.StatusOK, gin.H{"token": token, "role": role, "username": req.Username, "userID": userId})
}

// UpdateUserRole godoc
// @Summary Изменение роли пользователя
// @Description Изменяет роль пользователя (только для администраторов)
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param request body entity.UpdateRoleRequest true "Данные для изменения роли"
// @Success 200 {object} entity.RegisterResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 403 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/update-role [post]
func (h *AuthHandler) UpdateUserRole(c *gin.Context) {
	// Проверяем роль текущего пользователя
	token := c.GetHeader("Authorization")
	if token == "" {
		h.logger.Error("No authorization token provided")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
		return
	}

	role, err := h.jwtUtil.GetRoleFromToken(token)
	if err != nil {
		h.logger.Error("Failed to get role from token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен"})
		return
	}

	if role != "admin" {
		h.logger.Error("Unauthorized role update attempt", zap.String("role", role))
		c.JSON(http.StatusForbidden, gin.H{"error": "недостаточно прав"})
		return
	}

	var req entity.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON for role update", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authUsecase.UpdateUserRole(req.UserID, req.NewRole); err != nil {
		h.logger.Error("Failed to update user role", zap.Error(err), zap.Int("userID", req.UserID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("User role updated successfully", zap.Int("userID", req.UserID), zap.String("newRole", req.NewRole))
	c.JSON(http.StatusOK, gin.H{"message": "Роль пользователя успешно обновлена"})
}
