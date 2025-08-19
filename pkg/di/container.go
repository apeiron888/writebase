package di

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"write_base/config"
	"write_base/internal/delivery/http/controller"
	"write_base/internal/delivery/http/router"
	"write_base/internal/domain"
	"write_base/internal/infrastructure"
	"write_base/internal/infrastructure/ai"
	"write_base/internal/infrastructure/utils"
	"write_base/internal/policy"
	"write_base/internal/repository"
	"write_base/internal/usecase"

	usecaseai "write_base/internal/usecase/ai"
	usecasecomment "write_base/internal/usecase/comment"
	usecasefollow "write_base/internal/usecase/follow"
	usecasereaction "write_base/internal/usecase/reaction"
	usecasereport "write_base/internal/usecase/report"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Container struct {
	Router      *gin.Engine
	MongoClient *mongo.Client
}

func enableCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
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
		ID:        uuid.New().String(),
		Email:     email,
		Username:  username,
		Password:  hashedPassword,
		Role:      "super_admin",
		Verified:  true,
		IsActive:  true,
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

	// Ensure indexes for performance on articles collection
	if err := ensureArticleIndexes(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	//OAUTH
	//.............
	GoogleOAuthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// Repositories
	articleRepo := repository.NewArticleRepository(db, "articles")
	tagRepo := repository.NewTagRepository(db)
	viewRepo := repository.NewViewRepository(db)
	clapRepo := repository.NewClapRepository(db)

	userRepository := repository.NewUserRepository(db)

	commentRepo := repository.NewMongoCommentRepository(db.Collection("comments"))
	reactionRepo := repository.NewMongoReactionRepository(db.Collection("reactions"))
	followRepo := repository.NewMongoFollowRepository(db.Collection("follows"))
	reportRepo := repository.NewMongoReportRepository(db.Collection("reports"))

	// Utils
	utils := utils.NewUtils()

	startCleanupJob(userRepository, time.Hour, 5*time.Minute)
	// Clean up old revoked tokens (e.g., tokens revoked more than 24 hours ago)
	startRevokedTokenCleanupJob(userRepository, time.Hour, 5*time.Minute)

	// Auth service
	mailtrapService := infrastructure.NewMailtrapService(cfg.MailtrapHost, cfg.MailtrapPort, cfg.MailtrapUsername, cfg.MailtrapPassword, cfg.MailtrapFrom)
	passwordService := infrastructure.NewPasswordService()
	emailService := infrastructure.NewEmailService(mailtrapService, cfg.BackendURL)
	tokenService := infrastructure.NewJWTService([]byte(cfg.JwtSecret))
	authMiddleware := infrastructure.NewMiddleware(tokenService)

	// creating super admin
	err = SeedSuperAdmin(ctx, userRepository, passwordService)
	if err != nil {
		log.Fatal("Failed to seed super admin:", err)
	}

	// AI
	aiClient := ai.NewGeminiClient(cfg.GeminiAPIKey)
	// Policy
	policy := policy.NewArticlePolicy(utils)

	// Usecases
	tagUsecase := usecase.NewTagUsecase(tagRepo, utils)
	viewUsecase := usecase.NewViewUsecase(viewRepo, utils)
	clapUsecase := usecase.NewClapUsecase(clapRepo, utils)
	articleUsecase := usecase.NewArticleUsecase(articleRepo, policy, utils, tagUsecase, viewUsecase, clapUsecase, aiClient)

	userUsecase := usecase.NewUserUsecase(userRepository, passwordService, tokenService, emailService)

	commentUsecase := usecasecomment.NewCommentUsecase(commentRepo)
	reactionUsecase := usecasereaction.NewReactionService(reactionRepo)
	followUsecase := usecasefollow.NewFollowService(followRepo)
	reportUsecase := usecasereport.NewReportService(reportRepo)
	aiGemini := usecaseai.NewGeminiClient(cfg.GeminiAPIKey)

	// Handlers
	tagHandler := controller.NewTagHandler(tagUsecase)
	articleHandler := controller.NewArticleHandler(articleUsecase)

	userController := controller.NewUserController(userUsecase, GoogleOAuthConfig)

	commentController := controller.NewCommentController(commentUsecase)
	reactionController := controller.NewReactionController(reactionUsecase)
	followController := controller.NewFollowController(followUsecase)
	reportController := controller.NewReportController(reportUsecase)
	aiController := controller.NewAIController(aiGemini)

	r := gin.Default()
	r.Use(enableCORS())
	router.RegisterArticleRouter(r, articleHandler)
	router.RegisterTagRouter(r, tagHandler)

	router.UserRouter(r, userController, authMiddleware)
	router.RegisterCommentRoutes(r, commentController)
	router.RegisterReactionRoutes(r, reactionController)
	router.RegisterFollowRoutes(r, followController)
	router.RegisterReportRoutes(r, reportController)
	router.RegisterAIRoutes(r, aiController)

	return &Container{
		Router:      r,
		MongoClient: client,
	}, nil
}

// ensureArticleIndexes creates the indexes used by repository queries
func ensureArticleIndexes(ctx context.Context, db *mongo.Database) error {
	coll := db.Collection("articles")

	models := []mongo.IndexModel{
		// Unique slug for quick retrieval by slug
		{
			Keys:    bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uniq_slug"),
		},
		// Lists by author and status sorted by created_at
		{
			Keys:    bson.D{{Key: "author_id", Value: 1}, {Key: "status", Value: 1}, {Key: "timestamps.created_at", Value: -1}},
			Options: options.Index().SetName("author_status_createdat"),
		},
		// New articles scans by published_at
		{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "timestamps.published_at", Value: -1}},
			Options: options.Index().SetName("status_publishedat"),
		},
		// Popular sort by view_count within published
		{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "stats.view_count", Value: -1}},
			Options: options.Index().SetName("status_viewcount"),
		},
		// Trending: Equality(status), Sort(view_count), Range(published_at)
		{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "stats.view_count", Value: -1}, {Key: "timestamps.published_at", Value: -1}},
			Options: options.Index().SetName("status_viewcount_publishedat"),
		},
		// Tag filtering (multikey)
		// Equality(status), Sort(created_at), with tag filter
		{
			Keys:    bson.D{{Key: "status", Value: 1}, {Key: "tags", Value: 1}, {Key: "timestamps.created_at", Value: -1}},
			Options: options.Index().SetName("status_tags_createdat"),
		},
		// Language filter
		{
			Keys:    bson.D{{Key: "language", Value: 1}},
			Options: options.Index().SetName("language"),
		},
	}

	if _, err := coll.Indexes().CreateMany(ctx, models); err != nil {
		return err
	}

	// Optional text index for future text search on title/excerpt
	// Use default_language and language_override "none" to avoid unsupported language override errors
	textModel := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}, {Key: "title", Value: "text"}, {Key: "excerpt", Value: "text"}},
		Options: options.Index().
			SetDefaultLanguage("none").
			SetLanguageOverride("none").
			SetName("text_title_excerpt"),
	}
	// Ignore error if text index cannot be created due to old server; it's optional
	_, _ = coll.Indexes().CreateOne(ctx, textModel)

	return nil
}
