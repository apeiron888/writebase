package repository

import (
	"context"
	"time"
	"write_base/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ViewRepositoryImpl struct {
	collection *mongo.Collection
}

func NewViewRepository(db *mongo.Database) domain.ViewRepository {
	collection := db.Collection("views")
	
	// Create TTL index for auto-expiration after 1 day
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"created_at": 1},
		Options: options.Index().SetExpireAfterSeconds(24 * 60 * 60),
	}
	collection.Indexes().CreateOne(context.Background(), indexModel)
	
	return &ViewRepositoryImpl{collection: collection}
}

func (r *ViewRepositoryImpl) Create(ctx context.Context, view *domain.View) error {
	view.CreatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, view)
	return err
}