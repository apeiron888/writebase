package di

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gin-gonic/gin"

	"write_base/config"
	"write_base/internal/repository"
	"write_base/internal/usecase"
	"write_base/internal/delivery/http/controller"
	"write_base/internal/delivery/http/router"
	"write_base/internal/infrastructure"
)

type Container struct {
	Router *gin.Engine
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// MongoDB client
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongodbURI))
	if err != nil {
		return nil, err
	}
	db := client.Database(cfg.MongodbName)

	// Repositories
	//.............
	articleRepo := repository.NewArticleRepository(db, "articles")
	clapRepo := repository.NewArticleClapRepository(db, "article_claps")
	viewRepo := repository.NewArticleViewRepository(db, "article_views")

	// Usecases
	//.........
	clapUsecase := usecase.NewClapUsecase(clapRepo)
	viewUsecase := usecase.NewViewUsecase(viewRepo)
	articleUsecase := usecase.NewArticleUsecase(articleRepo, clapUsecase, viewUsecase)

	// Auth service
	//............

	// Handlers
	//.........
	articleHandler := controller.NewArticleHandler(articleUsecase)

	// Router
	// Add all handlers as params
	// Add Auth Sevice as param
	// Add cfg.ServerPort
	router := router.NewRouter(articleHandler)

	return &Container{
		Router: router,
	}, nil
}
