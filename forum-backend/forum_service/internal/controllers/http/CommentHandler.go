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
	userClient     grpc.UserClientInterface
}

func NewCommentHandler(commentUsecase usecase.CommentsUsecases, jwtUtil *utils.JWTUtil, logger *zap.Logger, userClient grpc.UserClientInterface) *CommentHandler {
	return &CommentHandler{commentUsecase: commentUsecase, jwtUtil: jwtUtil, logger: logger, userClient: userClient}
}

func (h *CommentHandler) Register(router *gin.Engine) {
	router.POST("/posts/:id/comments", h.CreateComment)
	router.GET("/posts/:id/comments", h.GetComments)
	router.DELETE("/comments/:id", h.DeleteComment)
}

// CreateComment godoc
// @Summary Создать новый комментарий
// @Description Создает новый комментарий к указанному посту
// @Tags Комментарии
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID поста"
// @Param comment body entity.Comment true "Данные комментария"
// @Success 201 {object} entity.Comment
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /posts/{id}/comments [post]
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

// GetComments returns paginated comments for a post
// @Summary Получить комментарии
// @Description Получить комментарии
// @Tags Комментарии
// @Accept json
// @Produce json
// @Param post_id path int true "Post ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} map[string]interface{} "comments and pagination info"
// @Router /posts/{post_id}/comments [get]
func (h *CommentHandler) GetComments(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
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
			username = "" // Используем пустое имя, если не удалось получить
		}

		commentsWithUsernames[i] = map[string]interface{}{
			"id":        comment.ID,
			"author_id": comment.AuthorId,
			"post_id":   comment.PostId,
			"content":   comment.Content,
			"username":  username, // Добавляем имя пользователя
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

// DeleteComment godoc
// @Summary Удалить комментарий
// @Description Удаляет комментарий по ID (доступно автору или администратору)
// @Tags Комментарии
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID комментария"
// @Success 204 "No Content"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 401 {object} entity.ErrorResponse
// @Failure 403 {object} entity.ErrorResponse
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		h.logger.Warn("Authorization header required")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == authHeader {
		h.logger.Warn("Invalid Authorization header format", zap.String("header", authHeader))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		h.logger.Warn("Invalid comment ID", zap.String("commentID", commentIDStr), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID, err := h.jwtUtil.GetUserIDFromToken(tokenString)
	if err != nil {
		h.logger.Warn("Invalid token or user ID", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or user ID"})
		return
	}

	userRole, err := h.jwtUtil.GetRoleFromToken(tokenString)
	if err != nil {
		h.logger.Warn("Invalid token or user role", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or user role"})
		return
	}

	if userRole != "admin" {
		comment, err := h.commentUsecase.GetCommentByID(c.Request.Context(), commentID)
		if err != nil {
			h.logger.Error("Failed to get comment", zap.Int("commentID", commentID), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comment"})
			return
		}

		if comment.AuthorId != userID {
			h.logger.Warn("Unauthorized attempt to delete comment",
				zap.Int("userID", userID),
				zap.Int("commentAuthorID", comment.AuthorId),
				zap.Int("commentID", commentID))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this comment"})
			return
		}
	}

	h.logger.Info("Deleting comment", zap.Int("commentID", commentID))
	err = h.commentUsecase.DeleteComment(c.Request.Context(), commentID)
	if err != nil {
		h.logger.Error("Failed to delete comment", zap.Int("commentID", commentID), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	h.logger.Info("Comment deleted successfully", zap.Int("commentID", commentID))
	c.Status(http.StatusNoContent)
}
