package http

import (
	"net/http"
	"strconv"

	utils "github.com/miqxzz/commonmiqx"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/chat"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/usecase"
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
