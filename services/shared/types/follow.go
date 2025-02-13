package types

import "time"

type FollowRequest struct {
	ID          string    `db:"id"`
	FollowerID  string    `db:"follower_id"`
	FollowingID string    `db:"following_id"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	// Relationship fields (optional, untuk join)
	Follower  *User `db:"follower,omitempty"`
	Following *User `db:"following,omitempty"`
}

type Follower struct {
	ID            string    `db:"id" json:"id"`
	FollowerID    string    `db:"follower_id" json:"followerId"`
	FollowingID   string    `db:"following_id" json:"followingId"`
	IsCloseFriend bool      `db:"is_close_friend" json:"isCloseFreind"`
	IsMuted       bool      `db:"is_muted" json:"isMuted"`
	IsBlocked     bool      `db:"is_blocked" json:"isBlocked"`
	CreatedAt     time.Time `db:"created_at" json:"createdAT"`
	UpdatedAt     time.Time `db:"updated_at" json:"UpdatedAT"`

	// Relationship fields (optional, untuk join)
	Follower  *User `db:"follower,omitempty" json:"folloers"`
	Following *User `db:"following,omitempty" json:"following"`
}

type CreateFollowRequest struct {
	FollowingID string `json:"following_id" validate:"required"`
	FollowerID  string `db:"follower_id" json:"followerId" validate:"required"`
}

type RespondFollowRequest struct {
	RequestID string `json:"request_id" validate:"required"`
	Status    string `json:"status" validate:"required,oneof=accepted rejected"`
}

type GetFollowersRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Limit  int    `json:"limit" validate:"required,min=1,max=100"`
	Offset int    `json:"offset" validate:"min=0"`
}

type GetFollowingRequest struct {
	UserID string `json:"user_id" validate:"required"`
	Limit  int    `json:"limit" validate:"required,min=1,max=100"`
	Offset int    `json:"offset" validate:"min=0"`
}
