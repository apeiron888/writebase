package di

import (
	"context"
	"time"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gin-gonic/gin"

	"write_base/config"
	"write_base/internal/repository"
	"write_base/internal/usecase"
	"write_base/internal/delivery/http/controller"
	"write_base/internal/delivery/http/router"
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
	//.............
	articleRepo := repository.NewArticleRepository(db, "articles")
	policy := usecase.NewArticlePolicy()
	// Usecases
	articleUsecase := usecase.NewArticleUsecase(articleRepo,policy)

	// Auth service
	//............

	// Handlers
	//.........
	articleHandler := controller.NewArticleHandler(articleUsecase)

	// Router
	// Add all handlers as params
	// Add Auth Sevice as param
	// Add cfg.ServerPort
	r:=gin.Default()
	router.RegisterArticleRouter(r,articleHandler)

	return &Container{
		Router:     r,
		MongoClient: client,
	}, nil
}
