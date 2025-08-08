package di

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"write_base/config"
	"write_base/internal/delivery/http/controller"
	"write_base/internal/delivery/http/router"
	"write_base/internal/domain"
	"write_base/internal/infrastructure"
	"write_base/internal/repository"
	"write_base/internal/usecase"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Container struct {
	Router *gin.Engine
	MongoClient *mongo.Client 
}
func SeedSuperAdmin(ctx context.Context, userRepo domain.IUserRepository, passwordService domain.IPasswordService) error {
	// Define the super admin credentials
	email := "super@admin.com"
	username := "superadmin"
	password := "SuperSecurePass123"

	// Check if already exists
	existingUser, err := userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		fmt.Println("Super admin already exists")
		return nil
	}

	// Hash password
	hashedPassword, err := passwordService.HashPassword(password)
	if err != nil {
		return err
	}

	// Create user model
	superAdmin := &domain.User{
		ID:       uuid.New().String(),
		Email:    email,
		Username: username,
		Password: hashedPassword,
		Role:     "super_admin",
		Verified: true,
		CreatedAt: time.Now(),
	}

	// Create in DB
	err = userRepo.CreateUser(ctx, superAdmin)
	if err != nil {
		return fmt.Errorf("failed to create super admin: %w", err)
	}

	fmt.Println("✅ Super admin created successfully")
	return nil
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
	//OAUTH
	//.............
	GoogleOAuthConfig := &oauth2.Config{
		ClientID:    cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Repositories
	//.............
	userRepository := repository.NewUserRepository(db)
	// clean up job
	startCleanupJob(userRepository, time.Hour, 5*time.Minute)
	// Clean up old revoked tokens (e.g., tokens revoked more than 24 hours ago)
    startRevokedTokenCleanupJob(userRepository, time.Hour, 5*time.Minute)
	// create supper admin
	
	
	// Auth service
	//............
	mailtrapService := infrastructure.NewMailtrapService(cfg.MailtrapHost, cfg.MailtrapPort, cfg.MailtrapUsername, cfg.MailtrapPassword, cfg.MailtrapFrom)
	passwordService := infrastructure.NewPasswordService()
	emailService :=infrastructure.NewEmailService(mailtrapService, cfg.BackendURL)
	tokenService := infrastructure.NewJWTService([]byte(cfg.JwtSecret))
	authMiddleware := infrastructure.NewMiddleware(tokenService)
	
	// creating super admin
	err = SeedSuperAdmin(ctx, userRepository, passwordService)
	if err != nil {
		log.Fatal("Failed to seed super admin:", err)
	}
	// Usecases
	//.........
	userUsecase:= usecase.NewUserUsecase(userRepository, passwordService, tokenService, emailService)
	
	// Handlers
	//.........
	userController := controller.NewUserController(userUsecase, GoogleOAuthConfig)

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
