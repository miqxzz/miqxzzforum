package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	commonmiqx "github.com/miqxzz/commonmiqx"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/config"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/grpc"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/http"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/repository"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/usecase"
	"go.uber.org/zap"
)

func main() {
	// Инициализация логгера
	logger, err := commonmiqx.InitLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// Инициализация подключения к базе данных
	db, err := sqlx.Open("sqlite3", cfg.DBPath)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Инициализация репозиториев
	postRepo := repository.NewPostRepository(db, logger)
	commentRepo := repository.NewCommentsRepository(db, logger)

	// Инициализация use cases
	postUsecase := usecase.NewPostUsecase(postRepo, logger)
	commentUsecase := usecase.NewCommentsUsecases(commentRepo, logger)

	// Инициализация gRPC клиента для пользователей
	userClient, err := grpc.NewUserClient(cfg.AuthServiceAddr)
	if err != nil {
		logger.Fatal("Failed to create user client", zap.Error(err))
	}
	defer userClient.Close()

	jwtUtil := commonmiqx.NewJWTUtil("secret")
	// Инициализация HTTP сервера
	router := gin.Default()
	http.NewPostHandler(postUsecase, postRepo, jwtUtil, logger, userClient).Register(router)
	http.NewCommentHandler(commentUsecase, jwtUtil, logger, userClient).Register(router)

	// Запуск HTTP сервера
	go func() {
		if err := router.Run(cfg.HTTPAddr); err != nil {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
}
