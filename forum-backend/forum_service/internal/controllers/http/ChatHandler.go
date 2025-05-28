package http

import (
	"github.com/gin-gonic/gin"
	utils "github.com/miqxzz/commonmiqx"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/chat"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/usecase"
	"go.uber.org/zap"
)

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
	// TODO: Implement WebSocket handler
}
