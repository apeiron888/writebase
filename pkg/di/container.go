package di

import (
	"context"
	"fmt"
	"log"
	"time"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"write_base/config"
	"write_base/internal/delivery/http/controller"
	"write_base/internal/delivery/http/router"
	"write_base/internal/domain"
	"write_base/internal/repository"
	usecaseai "write_base/internal/usecase/ai"
	usecasecomment "write_base/internal/usecase/comment"
	usecasefollow "write_base/internal/usecase/follow"
	usecasereaction "write_base/internal/usecase/reaction"
	usecasereport "write_base/internal/usecase/report"

	"github.com/gin-gonic/gin"
	"write_base/internal/usecase"
	"write_base/internal/infrastructure"

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

commentRepo := repository.NewMongoCommentRepository(db.Collection("comments"))
reactionRepo := repository.NewMongoReactionRepository(db.Collection("reactions"))
followRepo := repository.NewMongoFollowRepository(db.Collection("follows"))
reportRepo := repository.NewMongoReportRepository(db.Collection("reports"))

// Usecases
commentUsecase := usecasecomment.NewCommentUsecase(commentRepo)
reactionUsecase := usecasereaction.NewReactionService(reactionRepo)
followUsecase := usecasefollow.NewFollowService(followRepo)
reportUsecase := usecasereport.NewReportService(reportRepo)
aiGemini := usecaseai.NewGeminiClient(cfg.GeminiAPIKey)

// Controllers
commentController := controller.NewCommentController(commentUsecase)
reactionController := controller.NewReactionController(reactionUsecase)
followController := controller.NewFollowController(followUsecase)
reportController := controller.NewReportController(reportUsecase)
aiController := controller.NewAIController(aiGemini)


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
r := gin.Default()
// router.RegisterArticleRouter(r,articleHandler)
router.UserRouter(r, userController, authMiddleware)
router.RegisterCommentRoutes(r, commentController)
router.RegisterReactionRoutes(r, reactionController)
router.RegisterFollowRoutes(r, followController)
router.RegisterReportRoutes(r, reportController)
router.RegisterAIRoutes(r, aiController)

return &Container{
	   Router:     r,
	   MongoClient: client,
}, nil
}
