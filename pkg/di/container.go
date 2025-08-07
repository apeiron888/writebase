package di

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"write_base/config"
	"write_base/internal/delivery/http/controller"
	"write_base/internal/delivery/http/router"
	"write_base/internal/domain"
	"write_base/internal/infrastructure"
	"write_base/internal/repository"
	"write_base/internal/usecase"
)

type Container struct {
	Router *gin.Engine
	MongoClient *mongo.Client 
}
func startCleanupJob(userRepo domain.IUserRepository, interval, expiration time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            <-ticker.C
            ctx := context.Background()
            if err := userRepo.DeleteUnverifiedExpiredUsers(ctx, expiration); err != nil {
                fmt.Println("Cleanup job error:", err)
            }
        }
    }()
}
func startRevokedTokenCleanupJob(userRepo domain.IUserRepository, interval, olderThan time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            <-ticker.C
            ctx := context.Background()
            if err := userRepo.DeleteOldRevokedTokens(ctx, olderThan); err != nil {
                fmt.Println("Revoked token cleanup job error:", err)
            }
        }
    }()
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongodbURI))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB connection failed: %w", err)
	}
	db := client.Database(cfg.MongodbName)

	// Repositories
	//.............
	userRepository := repository.NewUserRepository(db)
	// clean up job
	startCleanupJob(userRepository, time.Hour, 5*time.Minute)
	// Clean up old revoked tokens (e.g., tokens revoked more than 24 hours ago)
    startRevokedTokenCleanupJob(userRepository, time.Hour, 5*time.Minute)


	
	// Auth service
	//............
	mailtrapService := infrastructure.NewMailtrapService(cfg.MailtrapHost, cfg.MailtrapPort, cfg.MailtrapUsername, cfg.MailtrapPassword, cfg.MailtrapFrom)
	passwordService := infrastructure.NewPasswordService()
	emailService :=infrastructure.NewEmailService(mailtrapService, cfg.BackendURL)
	tokenService := infrastructure.NewJWTService([]byte(cfg.JwtSecret))
	authMiddleware := infrastructure.NewMiddleware(tokenService)
	
	// Usecases
	//.........
	userUsecase:= usecase.NewUserUsecase(userRepository, passwordService, tokenService, emailService)
	
	// Handlers
	//.........
	userController := controller.NewUserController(userUsecase)

	// Router
	// Add all handlers as params
	// Add Auth Sevice as param
	// Add cfg.ServerPort
	r := gin.Default()
	router.UserRouter(r, userController, authMiddleware)

	return &Container{
		Router:     r,
		MongoClient: client,
	}, nil
}
