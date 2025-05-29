package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	utils "github.com/miqxzz/commonmiqx"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/chat"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/usecase"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	hub         *chat.Hub
	chatUsecase usecase.ChatUsecase
	jwtUtil     *utils.JWTUtil
	logger      *zap.Logger
}

func NewChatHandler(hub *chat.Hub, chatUsecase usecase.ChatUsecase, jwtUtil *utils.JWTUtil, logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		hub:         hub,
		chatUsecase: chatUsecase,
		jwtUtil:     jwtUtil,
		logger:      logger,
	}
}

func (h *ChatHandler) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", zap.Error(err))
		return
	}

	userID, _ := strconv.Atoi(c.Query("userID"))
	username := c.Query("username")
	isAuth := c.Query("auth") == "true"

	client := &chat.Client{
		Hub:             h.hub,
		Conn:            conn,
		Send:            make(chan []byte, 256),
		UserID:          userID,
		Username:        username,
		IsAuthenticated: isAuth,
		ChatUC:          h.chatUsecase,
	}

	h.hub.Register <- client
	go client.WritePump()
	client.ReadPump()
}
