package main

import (
	"context"
	"log"
	"time"
	"write_base/config"
	"write_base/internal/delivery/http/controller"
	"write_base/internal/delivery/http/router"
	"write_base/internal/repository"
	"write_base/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Public routes don't require auth
		if c.Request.URL.Path == "/p/:slug" {
			c.Next()
			return
		}
		
		// Mock token parsing
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}
		
		// Set mock user based on token
		switch token {
		case "admin_token":
			c.Set("user_id", "admin_user")
			c.Set("user_role", "admin")
		case "author_token":
			c.Set("user_id", "author_user")
			c.Set("user_role", "author")
		case "viewer_token":
			c.Set("user_id", "viewer_user")
			c.Set("user_role", "viewer")
		default:
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
func main() {
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal("Not loading env var")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongodbURI))
	if err != nil {
		log.Fatal("DB not connecting")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	db := client.Database(cfg.MongodbName)

	// Repositories
	articleRepo := repository.NewArticleRepository(db, "articles")
	policy := usecase.NewArticlePolicy()
	// Usecases
	articleUsecase := usecase.NewArticleUsecase(articleRepo, policy)

	// Handlers
	articleHandler := controller.NewArticleHandler(articleUsecase)

	// Router setup
	r := gin.Default()

	r.Use(MockAuthMiddleware())
	
	// Only register routes once!
	router.RegisterArticleRouter(r, articleHandler)

	r.Run(":" + cfg.ServerPort) // Use port from config
}