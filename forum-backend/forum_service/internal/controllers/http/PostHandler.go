package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	utils "github.com/miqxzz/commonmiqx"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/grpc"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/repository"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/usecase"
	"go.uber.org/zap"
)

type PostHandler struct {
	postUsecase usecase.PostUsecase
	postRepo    repository.PostRepository
	jwtUtil     *utils.JWTUtil
	logger      *zap.Logger
	userClient  *grpc.UserClient
}

func NewPostHandler(
	postUsecase usecase.PostUsecase,
	postRepo repository.PostRepository,
	jwtUtil *utils.JWTUtil,
	logger *zap.Logger,
	userClient *grpc.UserClient,
) *PostHandler {
	return &PostHandler{
		postUsecase: postUsecase,
		postRepo:    postRepo,
		jwtUtil:     jwtUtil,
		logger:      logger,
		userClient:  userClient,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
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

	userID, err := h.postRepo.GetUserIDByToken(c.Request.Context(), tokenString)
	if err != nil {
		h.logger.Warn("Invalid token", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var post entity.Post
	if err := c.BindJSON(&post); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.AuthorId = userID

	h.logger.Info("Creating post", zap.Any("post", post))
	createdPost, err := h.postUsecase.CreatePost(c.Request.Context(), post)
	if err != nil {
		h.logger.Error("Failed to create post", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Post created successfully", zap.Any("createdPost", createdPost))
	c.JSON(http.StatusCreated, createdPost)
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Получаем посты с пагинацией
	posts, err := h.postUsecase.GetPosts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем общее количество постов
	total, err := h.postUsecase.GetTotalPostsCount(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Добавляем имена пользователей к постам
	postsWithUsernames := make([]map[string]interface{}, len(posts))
	for i, post := range posts {
		username, err := h.userClient.GetUsername(c.Request.Context(), post.AuthorId)
		if err != nil {
			h.logger.Warn("Failed to get username",
				zap.Int("userID", post.AuthorId),
				zap.Error(err))
			username = ""
		}

		postsWithUsernames[i] = map[string]interface{}{
			"id":        post.ID,
			"title":     post.Title,
			"content":   post.Content,
			"author_id": post.AuthorId,
			"username":  username,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": postsWithUsernames,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
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

	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.logger.Warn("Invalid post ID", zap.String("postID", postIDStr), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
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
		post, err := h.postRepo.GetPostByID(c.Request.Context(), postID)
		if err != nil {
			h.logger.Error("Failed to get post", zap.Int("postID", postID), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get post"})
			return
		}

		if post.AuthorId != userID {
			h.logger.Warn("Unauthorized attempt to delete post",
				zap.Int("userID", userID),
				zap.Int("postAuthorID", post.AuthorId),
				zap.Int("postID", postID))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this post"})
			return
		}
	}

	h.logger.Info("Deleting post", zap.Int("postID", postID))
	err = h.postUsecase.DeletePost(c.Request.Context(), postID)
	if err != nil {
		h.logger.Error("Failed to delete post", zap.Int("postID", postID), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	h.logger.Info("Post deleted successfully", zap.Int("postID", postID))
	c.Status(http.StatusNoContent)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
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

	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		h.logger.Warn("Invalid post ID", zap.String("postID", postIDStr), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
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
	var newpost entity.Post
	if err := c.BindJSON(&newpost); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var updatedpost entity.Post
	if userRole != "admin" {
		post, err := h.postRepo.GetPostByID(c.Request.Context(), postID)
		post.Title = newpost.Title
		post.Content = newpost.Content
		if err != nil {
			h.logger.Error("Failed to get post", zap.Int("postID", postID), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get post"})
			return
		}

		if post.AuthorId != userID {
			h.logger.Warn("Unauthorized attempt to update post",
				zap.Int("userID", userID),
				zap.Int("postAuthorID", post.AuthorId),
				zap.Int("postID", postID))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this post"})
			return
		}
		h.logger.Info("Deleting post", zap.Int("postID", postID))
		updatedpost2, err := h.postUsecase.UpdatePost(c.Request.Context(), *post)
		updatedpost = *updatedpost2
		if err != nil {
			h.logger.Error("Failed to delete post", zap.Int("postID", postID), zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
			return
		}
	}

	h.logger.Info("Post deleted successfully", zap.Int("postID", postID))
	c.JSON(http.StatusOK, updatedpost)
}

func (h *PostHandler) GetTotalPostsCount(c *gin.Context) {
	count, err := h.postUsecase.GetTotalPostsCount(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get total posts count", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
