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
	"write_base/internal/repository"
	usecaseai "write_base/internal/usecase/ai"
	usecasecomment "write_base/internal/usecase/comment"
	usecasefollow "write_base/internal/usecase/follow"
	usecasereaction "write_base/internal/usecase/reaction"
	usecasereport "write_base/internal/usecase/report"
)

type Container struct {
	Router *gin.Engine
	MongoClient *mongo.Client 
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


// Router
r := gin.Default()
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
