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
	"write_base/internal/usecase"
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
commentUsecase := usecase.NewCommentUsecase(commentRepo)
reactionUsecase := usecase.NewReactionUsecase(reactionRepo)
followUsecase := usecase.NewFollowUsecase(followRepo)
reportUsecase := usecase.NewReportUsecase(reportRepo)

// Controllers
commentController := controller.NewCommentController(commentUsecase)
reactionController := controller.NewReactionController(reactionUsecase)
followController := controller.NewFollowController(followUsecase)
reportController := controller.NewReportController(reportUsecase)

// Router
r := gin.Default()
router.RegisterCommentRoutes(r, commentController)
router.RegisterReactionRoutes(r, reactionController)
router.RegisterFollowRoutes(r, followController)
router.RegisterReportRoutes(r, reportController)

return &Container{
	   Router:     r,
	   MongoClient: client,
}, nil
}
