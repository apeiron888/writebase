package dto

type FollowRequest struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}

type FollowResponse struct {
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
}
