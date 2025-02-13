package likes

import (
	"context"
	"fmt"
	"time"

	"github.com/wafi04/chatting-app/services/shared/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LikeRepository struct {
	mongoClient *mongo.Client
}

type Repository interface {
	ChangeLikeComment(ctx context.Context, userId string, commentId string) (*types.LikeComment, error)
	ChangeLikePost(ctx context.Context, userId string, postId string) error
	GetCommentLikesCount(ctx context.Context, commentId string) (int64, error)
	GetUserLiked(ctx context.Context, types, commentID, userID string) (*IsLikes, error)
	GetPostLikesCount(ctx context.Context, postId string) (int64, error)
	GetUserCommentLikes(ctx context.Context, userId string) ([]types.LikeComment, error)
	GetUserPostLikes(ctx context.Context, userId string) ([]types.LikePost, error)
}

func NewLikeRepository(mongoClient *mongo.Client) Repository {
	return &LikeRepository{
		mongoClient: mongoClient,
	}
}

func (lr *LikeRepository) ChangeLikeComment(ctx context.Context, userId string, commentId string) (*types.LikeComment, error) {
	likeCollection := lr.mongoClient.Database("chatapp").Collection("likes")

	filter := bson.M{
		"user_id":    userId,
		"comment_id": commentId,
	}

	// Check if like exists
	exists := likeCollection.FindOne(ctx, filter)
	if exists.Err() == nil {
		// Like exists, so remove it
		_, err := likeCollection.DeleteOne(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to remove like: %w", err)
		}
		return nil, nil
	}

	// Like doesn't exist, create it
	like := &types.LikeComment{
		ID:        primitive.NewObjectID(),
		UserID:    userId,
		CommentID: &commentId,
		CreatedAt: time.Now(),
	}

	_, err := likeCollection.InsertOne(ctx, like)
	if err != nil {
		return nil, fmt.Errorf("failed to create like: %w", err)
	}

	return like, nil
}

func (lr *LikeRepository) ChangeLikePost(ctx context.Context, userId string, postId string) error {
	likeCollection := lr.mongoClient.Database("chatapp").Collection("likes")

	filter := bson.M{
		"user_id": userId,
		"post_id": postId,
	}

	// Check if like exists
	exists := likeCollection.FindOne(ctx, filter)
	if exists.Err() == nil {
		// Like exists, so remove it
		_, err := likeCollection.DeleteOne(ctx, filter)
		if err != nil {
			return fmt.Errorf("failed to remove like: %w", err)
		}
		return nil
	}

	// Like doesn't exist, create it
	like := &types.LikePost{
		ID:        primitive.NewObjectID(),
		UserID:    userId,
		PostID:    &postId,
		CreatedAt: time.Now(),
	}

	_, err := likeCollection.InsertOne(ctx, like)
	if err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}

	return nil
}

func (lr *LikeRepository) GetCommentLikesCount(ctx context.Context, commentId string) (int64, error) {
	likeCollection := lr.mongoClient.Database("chatapp").Collection("likes")

	filter := bson.M{
		"comment_id": commentId,
	}

	count, err := likeCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to get likes count: %w", err)
	}

	return count, nil
}

func (lr *LikeRepository) GetPostLikesCount(ctx context.Context, postId string) (int64, error) {
	likeCollection := lr.mongoClient.Database("chatapp").Collection("likes")

	filter := bson.M{
		"post_id": postId,
	}

	count, err := likeCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to get likes count: %w", err)
	}

	return count, nil
}

func (lr *LikeRepository) GetUserCommentLikes(ctx context.Context, userId string) ([]types.LikeComment, error) {
	likeCollection := lr.mongoClient.Database("chatapp").Collection("likes")

	filter := bson.M{
		"user_id":    userId,
		"comment_id": bson.M{"$exists": true},
	}

	cursor, err := likeCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get user comment likes: %w", err)
	}
	defer cursor.Close(ctx)

	var likes []types.LikeComment
	if err = cursor.All(ctx, &likes); err != nil {
		return nil, fmt.Errorf("failed to decode user comment likes: %w", err)
	}

	return likes, nil
}
func (lr *LikeRepository) GetUserLiked(ctx context.Context, types, commentID, userID string) (*IsLikes, error) {
	likeCollection := lr.mongoClient.Database("chatapp").Collection("likes")

	filter := bson.M{
		"user_id": userID,
		types:     commentID,
	}

	err := likeCollection.FindOne(ctx, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &IsLikes{Liked: false}, nil // No like found, return false
		}
		return nil, err // Return error for other cases
	}

	return &IsLikes{Liked: true}, nil // Document exists, return true
}

type IsLikes struct {
	Liked bool `json:"liked"`
}

func (lr *LikeRepository) GetUserPostLikes(ctx context.Context, userId string) ([]types.LikePost, error) {
	likeCollection := lr.mongoClient.Database("chatapp").Collection("likes")

	filter := bson.M{
		"user_id": userId,
		"post_id": bson.M{"$exists": true},
	}

	cursor, err := likeCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get user post likes: %w", err)
	}
	defer cursor.Close(ctx)

	var likes []types.LikePost
	if err = cursor.All(ctx, &likes); err != nil {
		return nil, fmt.Errorf("failed to decode user post likes: %w", err)
	}

	return likes, nil
}
