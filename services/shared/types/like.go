package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Like represents a like for either a post or a comment
type Like struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id"  json:"userId"`
	PostID    *string            `bson:"post_id,omitempty"  json:"postId"`
	CommentID *string            `bson:"comment_id,omitempty" json:"commentId"`
	CreatedAt time.Time          `bson:"created_at"  json:"createdAt"`
}

type LikeComment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id"  json:"userId"`
	CommentID *string            `bson:"comment_id,omitempty" json:"commentId"`
	CreatedAt time.Time          `bson:"created_at"  json:"createdAt"`
}

type LikePost struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id"  json:"userId"`
	PostID    *string            `bson:"post_id,omitempty"  json:"postId"`
	CreatedAt time.Time          `bson:"created_at"  json:"createdAt"`
}
