package follow

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
	"github.com/wafi04/chatting-app/services/shared/types"
	"github.com/wafi04/chatting-app/services/shared/utils"
)

func (r *FollowRepository) CreateFollow(ctx context.Context, tx *sqlx.Tx, req *types.Follower) (string, error) {
	query := `
        INSERT INTO followers (
            id,
            follower_id,
            following_id,
            is_close_friend,
            is_blocked,
            created_at,
            updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        `
	followerID := utils.GenerateCustomID(utils.IDOptions{
		Prefix:       "FLW",
		NumberLength: 7,
	})
	_, err := tx.ExecContext(ctx,
		query,
		followerID,
		req.FollowerID,
		req.FollowingID,
		false,
		false,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return "", err
	}

	return followerID, nil
}

func (r *FollowRepository) GetFollowers(ctx context.Context, req *types.GetFollowersRequest) ([]types.Follower, error) {
	query := `
    SELECT 
        id,
        follower_id,
        following_id,
        is_close_friend,
        is_blocked,
        created_at,
        updated_at
    FROM followers
    WHERE following_id = $1
    ORDER BY created_at DESC
    LIMIT $2 OFFSET $3
    `

	var followers []types.Follower
	err := r.DB.SelectContext(ctx, &followers, query, req.UserID, req.Limit, req.Offset)
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to get followers: %v", err)
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	return followers, nil
}

func (r *FollowRepository) GetFollowings(ctx context.Context, req *types.GetFollowingRequest) ([]types.Follower, error) {
	query := `
    SELECT 
        id,
        follower_id,
        following_id,
        is_close_friend,
        is_blocked,
        created_at,
        updated_at
    FROM followers
    WHERE follower_id = $1
    ORDER BY created_at DESC
    LIMIT $2 OFFSET $3
    `

	var followings []types.Follower
	err := r.DB.SelectContext(ctx, &followings, query, req.UserID, req.Limit, req.Offset)
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to get followings: %v", err)
		return nil, fmt.Errorf("failed to get followings: %w", err)
	}

	return followings, nil
}

func (r *FollowRepository) Unfollow(ctx context.Context, followerID, followingID string) error {
	query := `
    DELETE FROM followers
    WHERE follower_id = $1 AND following_id = $2
    `

	result, err := r.DB.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to unfollow: %v", err)
		return fmt.Errorf("failed to unfollow: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to check rows affected: %v", err)
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no such following relationship exists")
	}

	return nil
}
