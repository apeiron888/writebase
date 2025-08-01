package repository

import (
	"context"
	"starter/internal/domain"
	dtodbrep "starter/internal/repository/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoReactionRepository struct {
	collection *mongo.Collection
}

func NewMongoReactionRepository(collection *mongo.Collection) *MongoReactionRepository {
	return &MongoReactionRepository{collection: collection}
}

func (r *MongoReactionRepository) AddReaction(ctx context.Context, reaction *domain.Reaction) error {
	dto := dtodbrep.ReactionResponse{
		ID:        reaction.ID,
		PostID:    reaction.PostID,
		UserID:    reaction.UserID,
		CommentID: reaction.CommentID,
		Type:      string(reaction.Type),
		CreatedAt: reaction.CreatedAt,
	}
	_, err := r.collection.InsertOne(ctx, dto)
	return err
}

func (r *MongoReactionRepository) RemoveReaction(ctx context.Context, reactionID string) error {
	filter := bson.M{"_id": reactionID}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrReactionNotFound
	}
	return nil
}

func (r *MongoReactionRepository) GetReactionsByPost(ctx context.Context, postID string) ([]*domain.Reaction, error) {
	filter := bson.M{"post_id": postID}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []*domain.Reaction
	for cur.Next(ctx) {
		var dto dtodbrep.ReactionResponse
		if err := cur.Decode(&dto); err != nil {
			return nil, err
		}
		results = append(results, &domain.Reaction{
			ID:        dto.ID,
			PostID:    dto.PostID,
			UserID:    dto.UserID,
			CommentID: dto.CommentID,
			Type:      domain.ReactionType(dto.Type),
			CreatedAt: dto.CreatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoReactionRepository) GetReactionsByUser(ctx context.Context, userID string) ([]*domain.Reaction, error) {
	filter := bson.M{"user_id": userID}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []*domain.Reaction
	for cur.Next(ctx) {
		var dto dtodbrep.ReactionResponse
		if err := cur.Decode(&dto); err != nil {
			return nil, err
		}
		results = append(results, &domain.Reaction{
			ID:        dto.ID,
			PostID:    dto.PostID,
			UserID:    dto.UserID,
			CommentID: dto.CommentID,
			Type:      domain.ReactionType(dto.Type),
			CreatedAt: dto.CreatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoReactionRepository) CountReactions(ctx context.Context, postID string, reactionType domain.ReactionType) (int, error) {
	filter := bson.M{"post_id": postID, "type": string(reactionType)}
	count, err := r.collection.CountDocuments(ctx, filter)
	return int(count), err
}
