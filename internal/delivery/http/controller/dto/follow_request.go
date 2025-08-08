package dto

type FollowRequest struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}
