package di

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"write_base/config"
	"write_base/internal/infrastructure/ai"
	"write_base/internal/infrastructure/utils"
	"write_base/internal/delivery/http/controller"
	"write_base/internal/delivery/http/router"
	"write_base/internal/policy"
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
	articleRepo := repository.NewArticleRepository(db, "articles")
	tagRepo := repository.NewTagRepository(db)
	viewRepo := repository.NewViewRepository(db)
	clapRepo := repository.NewClapRepository(db)

	// Utils
	utils := utils.NewUtils()

	// AI
	aiClient := ai.NewGeminiClient(cfg.GeminiAPIKey)
	// Policy
	policy := policy.NewArticlePolicy(utils)

	// Usecases
	tagUsecase := usecase.NewTagUsecase(tagRepo,utils)
	viewUsecase := usecase.NewViewUsecase(viewRepo,utils)
	clapUsecase := usecase.NewClapUsecase(clapRepo,utils)
	articleUsecase := usecase.NewArticleUsecase(articleRepo,policy,utils,tagUsecase,viewUsecase,clapUsecase,aiClient)

	// Handlers
	tagHandler := controller.NewTagHandler(tagUsecase)
	articleHandler := controller.NewArticleHandler(articleUsecase)
	


	r:=gin.Default()
	router.RegisterArticleRouter(r,articleHandler)
	router.RegisterTagRouter(r, tagHandler)

	return &Container{
		Router:     r,
		MongoClient: client,
	}, nil
}
