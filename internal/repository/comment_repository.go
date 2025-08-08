package repository

import (
	"context"
	"log"
	"write_base/internal/domain"
	dtodbrep "write_base/internal/repository/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCommentRepository struct {
	collection *mongo.Collection
}

func NewMongoCommentRepository(collection *mongo.Collection) *MongoCommentRepository {
	return &MongoCommentRepository{collection: collection}
}

func (r *MongoCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	dto := dtodbrep.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
	_, err := r.collection.InsertOne(ctx, dto)
	return err
}

func (r *MongoCommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	filter := bson.M{"_id": comment.ID}
	update := bson.M{"$set": bson.M{
		"content":    comment.Content,
		"updated_at": comment.UpdatedAt,
	}}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return domain.ErrCommentNotFound
	}
	return nil
}

func (r *MongoCommentRepository) Delete(ctx context.Context, commentID string) error {
	filter := bson.M{"_id": commentID}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrCommentNotFound
	}
	return nil
}

func (r *MongoCommentRepository) GetByID(ctx context.Context, commentID string) (*domain.Comment, error) {
	log.Println("Fetching comment by ID:", commentID)
	objID, err := primitive.ObjectIDFromHex(commentID)
    if err != nil {
        return nil, domain.ErrCommentNotFound // or return a 400 error if you want
    }
	filter := bson.M{"_id": objID}
	var dto dtodbrep.CommentResponse
	err = r.collection.FindOne(ctx, filter).Decode(&dto)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrCommentNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domain.Comment{
		ID:        dto.ID,
		PostID:    dto.PostID,
		UserID:    dto.UserID,
		ParentID:  dto.ParentID,
		Content:   dto.Content,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}, nil
}

func (r *MongoCommentRepository) GetByPostID(ctx context.Context, postID string) ([]*domain.Comment, error) {
	filter := bson.M{"post_id": postID}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []*domain.Comment
	for cur.Next(ctx) {
		var dto dtodbrep.CommentResponse
		if err := cur.Decode(&dto); err != nil {
			return nil, err
		}
		results = append(results, &domain.Comment{
			ID:        dto.ID,
			PostID:    dto.PostID,
			UserID:    dto.UserID,
			ParentID:  dto.ParentID,
			Content:   dto.Content,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoCommentRepository) GetByUserID(ctx context.Context, userID string) ([]*domain.Comment, error) {
	filter := bson.M{"user_id": userID}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []*domain.Comment
	for cur.Next(ctx) {
		var dto dtodbrep.CommentResponse
		if err := cur.Decode(&dto); err != nil {
			return nil, err
		}
		results = append(results, &domain.Comment{
			ID:        dto.ID,
			PostID:    dto.PostID,
			UserID:    dto.UserID,
			ParentID:  dto.ParentID,
			Content:   dto.Content,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *MongoCommentRepository) GetReplies(ctx context.Context, parentID string) ([]*domain.Comment, error) {
	filter := bson.M{"parent_id": parentID}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []*domain.Comment
	for cur.Next(ctx) {
		var dto dtodbrep.CommentResponse
		if err := cur.Decode(&dto); err != nil {
			return nil, err
		}
		results = append(results, &domain.Comment{
			ID:        dto.ID,
			PostID:    dto.PostID,
			UserID:    dto.UserID,
			ParentID:  dto.ParentID,
			Content:   dto.Content,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
