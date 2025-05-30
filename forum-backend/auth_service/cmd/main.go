package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	user "github.com/miqxzz/miqxzzforum/auth_service/internal/proto"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	commonmiqx "github.com/miqxzz/commonmiqx"
	_ "github.com/miqxzz/miqxzzforum/auth_service/docs"
	"github.com/miqxzz/miqxzzforum/auth_service/internal/config"
	mygrpc "github.com/miqxzz/miqxzzforum/auth_service/internal/delivery/grpc"
	"github.com/miqxzz/miqxzzforum/auth_service/internal/delivery/http"
	"github.com/miqxzz/miqxzzforum/auth_service/internal/repository"
	"github.com/miqxzz/miqxzzforum/auth_service/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// @title Auth Service API
// @version 1.0
// @description This is the API documentation for the Auth Service.
// @host localhost:8080
// @BasePath /

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

	logger.Info("Configuration loaded",
		zap.String("AUTH_SERVICE_PORT", cfg.Port),
		zap.String("AUTH_SERVICE_DB_PATH", cfg.DBPath),
		zap.String("AUTH_SERVICE_MIGRATIONS_PATH", cfg.MigrationsPath),
		zap.String("JWT_SECRET", cfg.JWTSecret),
	)

	// Инициализация подключения к базе данных
	db, err := sqlx.Open("sqlite3", cfg.DBPath)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		logger.Fatal("Failed to create migrate driver", zap.Error(err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+cfg.MigrationsPath,
		"sqlite3", driver)
	if err != nil {
		logger.Fatal("Failed to create migrate instance", zap.Error(err))
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("Failed to apply migrations", zap.Error(err))
	}
	logger.Info("Migrations applied successfully")

	// Инициализация репозитория
	userRepo := repository.NewAuthRepository(db, logger)
	userServer := mygrpc.NewUserServer(userRepo)

	// Инициализация gRPC сервера
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userServer)

	// Запускаем gRPC сервер
	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatal(err)
		}
		log.Println("gRPC user server started on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()


	
	jwtUtil := commonmiqx.NewJWTUtil(cfg.JWTSecret)
	userUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil, logger)
	authHandler := http.NewAuthHandler(userUsecase, jwtUtil, logger)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
		AllowHeaders:     []string{"Content-type", "Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)
	router.POST("/auth/update-role", authHandler.UpdateUserRole)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(cfg.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}

	// Ожидаем сигнал завершения
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down server")
}
