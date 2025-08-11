package repository

import (
	"context"
	"time"
	"write_base/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClapRepositoryImpl struct {
	collection *mongo.Collection
}

type ClapDTO struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	ArticleID string    `bson:"article_id"`
	Count     int       `bson:"count"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func NewClapRepository(db *mongo.Database) domain.ClapRepository {
	return &ClapRepositoryImpl{
		collection: db.Collection("claps"),
	}
}

func toClapDTO(clap *domain.Clap) *ClapDTO {
	return &ClapDTO{
		ID:        clap.ID,
		UserID:    clap.UserID,
		ArticleID: clap.ArticleID,
		Count:     clap.Count,
		CreatedAt: clap.CreatedAt,
		UpdatedAt: clap.UpdatedAt,
	}
}

func toDomainClap(dto *ClapDTO) *domain.Clap {
	return &domain.Clap{
		ID:        dto.ID,
		UserID:    dto.UserID,
		ArticleID: dto.ArticleID,
		Count:     dto.Count,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func (r *ClapRepositoryImpl) Create(ctx context.Context, clap *domain.Clap) error {
	clap.CreatedAt = time.Now()
	clap.UpdatedAt = time.Now()
	
	dto := toClapDTO(clap)
	_, err := r.collection.InsertOne(ctx, dto)
	return err
}

func (r *ClapRepositoryImpl) Update(ctx context.Context, clap *domain.Clap) error {
	clap.UpdatedAt = time.Now()
	
	filter := bson.M{"_id": clap.ID}
	update := bson.M{"$set": bson.M{
		"count":      clap.Count,
		"updated_at": clap.UpdatedAt,
	}}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *ClapRepositoryImpl) GetByUserAndArticle(ctx context.Context, userID, articleID string) (*domain.Clap, error) {
	filter := bson.M{
		"user_id":    userID,
		"article_id": articleID,
	}
	
	var dto ClapDTO
	err := r.collection.FindOne(ctx, filter).Decode(&dto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	
	return toDomainClap(&dto), nil
}

func (r *ClapRepositoryImpl) GetArticleClapCount(ctx context.Context, articleID string) (int, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{"article_id": articleID},
		},
		{
			"$group": bson.M{
				"_id":   nil,
				"total": bson.M{"$sum": "$count"},
			},
		},
	}
	
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	
	var result struct {
		Total int `bson:"total"`
	}
	
	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Total, nil
	}
	
	return 0, nil
}