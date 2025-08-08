package dto

type FollowResponse struct {
 FollowerID string `json:"follower_id" bson:"follower_id"`
 FolloweeID string `json:"followee_id" bson:"followee_id"`
}
