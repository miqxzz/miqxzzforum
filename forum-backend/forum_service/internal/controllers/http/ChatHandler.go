package http

import (
	utils "github.com/Engls/EnglsJwt"
	"net/http"
	"strconv"

	"github.com/Engls/forum-project2/forum_service/internal/controllers/chat"
	"github.com/Engls/forum-project2/forum_service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatHandler struct {
	hub     *chat.Hub
	chatUC  usecase.ChatUsecase
	jwtUtil *utils.JWTUtil
	logger  *zap.Logger
}

func NewChatHandler(hub *chat.Hub, chatUC usecase.ChatUsecase, jwtUtil *utils.JWTUtil, logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		hub:     hub,
		chatUC:  chatUC,
		jwtUtil: jwtUtil,
		logger:  logger,
	}
}

// ServeWS godoc
// @Summary Установить WebSocket соединение для чата
// @Description Обновляет HTTP соединение до WebSocket для обмена сообщениями в реальном времени
// @Tags Чат
// @Accept json
// @Produce json
// @Param token query string true "JWT токен авторизации"
// @Param userID query int true "ID пользователя"
// @Param username query string false "Имя пользователя"
// @Param auth query bool true "Флаг аутентификации"
// @Success 101 "Switching Protocols" {object} nil
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /ws/chat [get]
func (h *ChatHandler) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade error", zap.Error(err))
		return
	}

	realUserID, err := h.jwtUtil.GetUserIDFromToken(c.Query("token"))
	if err != nil {
		h.logger.Error("Failed to get user ID from token", zap.Error(err))
	}
	userID, err := strconv.Atoi(c.Query("userID"))
	if err != nil {
		h.logger.Error("Failed to convert user ID to integer", zap.Error(err))
	}
	h.logger.Info("Real userID", zap.Int("realUserID", realUserID))
	username := c.Query("username")
	isAuthenticated := c.Query("auth") == "true"

	client := &chat.Client{
		Hub:             h.hub,
		Conn:            conn,
		Send:            make(chan []byte, 256),
		UserID:          userID,
		Username:        username,
		IsAuthenticated: isAuthenticated,
		ChatUC:          h.chatUC,
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
