package repository

import (
	"context"
	"write_base/internal/domain"
	dtodbrep "write_base/internal/repository/dto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoFollowRepository struct {
	collection *mongo.Collection
}

func NewMongoFollowRepository(collection *mongo.Collection) *MongoFollowRepository {
	return &MongoFollowRepository{collection: collection}
}

func (r *MongoFollowRepository) FollowUser(ctx context.Context, followerID, followeeID string) error {
	doc := dtodbrep.FollowResponse{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}
	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *MongoFollowRepository) UnfollowUser(ctx context.Context, followerID, followeeID string) error {
	filter := bson.M{"follower_id": followerID, "followee_id": followeeID}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return domain.ErrFollowNotFound
	}
	return nil
}

func (r *MongoFollowRepository) GetFollowers(ctx context.Context, userID string) ([]*domain.User, error) {
	filter := bson.M{"followee_id": userID}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var followers []*domain.User
	for cur.Next(ctx) {
		var dto dtodbrep.FollowResponse
		if err := cur.Decode(&dto); err != nil {
			return nil, err
		}
		followers = append(followers, &domain.User{ID: dto.FollowerID})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return followers, nil
}

func (r *MongoFollowRepository) GetFollowing(ctx context.Context, userID string) ([]*domain.User, error) {
	filter := bson.M{"follower_id": userID}
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var following []*domain.User
	for cur.Next(ctx) {
		var dto dtodbrep.FollowResponse
		if err := cur.Decode(&dto); err != nil {
			return nil, err
		}
		following = append(following, &domain.User{ID: dto.FolloweeID})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return following, nil
}

func (r *MongoFollowRepository) IsFollowing(ctx context.Context, followerID, followeeID string) (bool, error) {
	filter := bson.M{"follower_id": followerID, "followee_id": followeeID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
