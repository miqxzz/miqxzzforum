package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	utils "github.com/miqxzz/commonmiqx"
	_ "github.com/miqxzz/miqxzzforum/forum_service/docs"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/config"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/chat"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/grpc"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/http"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/repository"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	logger.Info("Configuration loaded",
		zap.String("DB_PATH", cfg.DBPath),
		zap.String("PORT", cfg.Port),
		zap.String("JWT_SECRET", cfg.JWTSecret),
	)

	db, err := sqlx.Open("sqlite3", cfg.DBPath)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	userClient, err := grpc.NewUserClient("localhost:50052") // Или localhost для тестов
	if err != nil {
		log.Fatalf("Failed to create user client: %v", err)
	}
	defer userClient.Close()

	postRepo := repository.NewPostRepository(db, logger)
	commentRepo := repository.NewCommentsRepository(db, logger)
	chatRepo := repository.NewChatRepository(db, logger)
	postUsecase := usecase.NewPostUsecase(postRepo, logger)
	commentUsecase := usecase.NewCommentsUsecases(commentRepo, logger)
	hub := chat.NewHub()
	chatUsecase := usecase.NewChatUsecase(chatRepo, logger)
	jwtUtil := utils.NewJWTUtil(cfg.JWTSecret)

	postHandler := http.NewPostHandler(postUsecase, postRepo, jwtUtil, logger, userClient)
	commentHandler := http.NewCommentHandler(commentUsecase, jwtUtil, logger, userClient)
	chatHandler := http.NewChatHandler(hub, chatUsecase, jwtUtil, logger)

	go hub.Run()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Content-type", "Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ws", chatHandler.ServeWS)
	router.POST("/posts", postHandler.CreatePost)
	router.GET("/posts", postHandler.GetPosts)
	router.DELETE("/posts/:id", postHandler.DeletePost)
	router.POST("/posts/:id/comments", commentHandler.CreateComment)
	router.GET("/posts/:post_id/comments", commentHandler.GetComments)
	router.PUT("/posts/:id", postHandler.UpdatePost)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(cfg.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
