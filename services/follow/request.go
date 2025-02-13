package follow

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
	"github.com/wafi04/chatting-app/services/shared/types"
	"github.com/wafi04/chatting-app/services/shared/utils"
)

func (r *FollowRepository) CreateFollowRequest(ctx context.Context, tx *sqlx.Tx, req *types.CreateFollowRequest) (string, error) {
	var data types.FollowRequest

	query := `
        INSERT INTO follow_request (
            id,
            follower_id,
            following_id,
            status,
            created_at,
            updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, follower_id, following_id, status, created_at
        `
	reqID := utils.GenerateCustomID(utils.IDOptions{
		Prefix:       "REQ",
		NumberLength: 7,
	})
	err := tx.QueryRowContext(ctx,
		query,
		reqID,
		req.FollowerID,
		req.FollowingID,
		"PENDING",
		time.Now(),
		time.Now(),
	).Scan(
		&data.ID,
		&data.FollowerID,
		&data.FollowingID,
		&data.Status,
		&data.CreatedAt,
	)

	if err != nil {
		return "", err
	}

	return data.ID, nil
}

func (r *FollowRepository) RejectFollowRequest(ctx context.Context, requestID string) error {
	query := `
    DELETE FROM follow_request
    WHERE id = $1
    `

	result, err := r.DB.ExecContext(ctx, query, requestID)
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to reject follow request: %v", err)
		return fmt.Errorf("failed to reject follow request: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to check rows affected: %v", err)
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no such follow request exists")
	}

	return nil
}

func (r *FollowRepository) AcceptFollowRequest(ctx context.Context, requestID string) error {
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to begin transaction: %v", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Retrieve the follow request details
	var followRequest types.FollowRequest
	query := `
    SELECT follower_id, following_id
    FROM follow_request
    WHERE id = $1
    `
	err = tx.QueryRowContext(ctx, query, requestID).Scan(&followRequest.FollowerID, &followRequest.FollowingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("follow request not found")
		}
		return fmt.Errorf("failed to retrieve follow request: %w", err)
	}

	// Insert into followers table
	insertQuery := `
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
	_, err = tx.ExecContext(ctx, insertQuery, followerID, followRequest.FollowerID, followRequest.FollowingID, false, false, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert into followers: %w", err)
	}

	// Delete the follow request
	deleteQuery := `
    DELETE FROM follow_request
    WHERE id = $1
    `
	_, err = tx.ExecContext(ctx, deleteQuery, requestID)
	if err != nil {
		return fmt.Errorf("failed to delete follow request: %w", err)
	}

	return nil
}
