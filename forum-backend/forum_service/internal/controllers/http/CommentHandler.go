package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	utils "github.com/miqxzz/commonmiqx"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/grpc"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/usecase"
	"go.uber.org/zap"
)

type CommentHandler struct {
	commentUsecase usecase.CommentsUsecases
	jwtUtil        *utils.JWTUtil
	logger         *zap.Logger
	userClient     *grpc.UserClient
}

func NewCommentHandler(commentUsecase usecase.CommentsUsecases, jwtUtil *utils.JWTUtil, logger *zap.Logger, userClient *grpc.UserClient) *CommentHandler {
	return &CommentHandler{commentUsecase: commentUsecase, jwtUtil: jwtUtil, logger: logger, userClient: userClient}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.logger.Error("Invalid post ID", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var comment entity.Comment
	if err := c.BindJSON(&comment); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.logger.Error("Authorization header required")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == authHeader {
		h.logger.Error("Invalid Authorization header format")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	userID, err := h.jwtUtil.GetUserIDFromToken(tokenString)
	if err != nil {
		h.logger.Error("Invalid token or user ID", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or user ID"})
		return
	}

	comment.PostId = postID
	comment.AuthorId = userID

	createdComment, err := h.commentUsecase.CreateComment(c.Request.Context(), comment)
	if err != nil {
		h.logger.Error("Failed to create comment", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Comment created successfully", zap.Int("postID", postID), zap.Int("userID", userID))
	c.JSON(http.StatusCreated, createdComment)
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	comments, err := h.commentUsecase.GetComments(c.Request.Context(), postID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := h.commentUsecase.GetTotalCommentsCount(c.Request.Context(), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	commentsWithUsernames := make([]map[string]interface{}, len(comments))
	for i, comment := range comments {
		username, err := h.userClient.GetUsername(c.Request.Context(), comment.AuthorId)
		if err != nil {
			h.logger.Warn("Failed to get username",
				zap.Int("userID", comment.AuthorId),
				zap.Error(err))
			username = ""
		}

		commentsWithUsernames[i] = map[string]interface{}{
			"id":        comment.ID,
			"author_id": comment.AuthorId,
			"post_id":   comment.PostId,
			"content":   comment.Content,
			"username":  username,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": commentsWithUsernames,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *CommentHandler) GetTotalCommentsCount(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.logger.Error("Invalid post ID", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	count, err := h.commentUsecase.GetTotalCommentsCount(c.Request.Context(), postID)
	if err != nil {
		h.logger.Error("Failed to get total comments count", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
