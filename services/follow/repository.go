package follow

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
	"github.com/wafi04/chatting-app/services/shared/types"
	"github.com/wafi04/chatting-app/services/shared/utils"
	"github.com/wafi04/chatting-app/services/user"
)

type FollowRepository struct {
	DB       *sqlx.DB
	userRepo *user.UserRepository
	log      logger.Logger
}

func NewFollowRepository(db *sqlx.DB, userrepo *user.UserRepository) *FollowRepository {
	return &FollowRepository{
		DB:       db,
		userRepo: userrepo,
	}
}

func (r *FollowRepository) AddFollower(ctx context.Context, req *types.CreateFollowRequest) (*types.RespondFollowRequest, error) {
	// Begin transaction
	tx, err := r.DB.BeginTxx(ctx, nil)
	if err != nil {
		r.log.Log(logger.ErrorLevel, "Failed to begin transaction: %v", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				r.log.Log(logger.ErrorLevel, "Failed to commit transaction: %v", err)
			}
		}
	}()

	// Check if the user being followed has a private account
	isPrivate, err := r.userRepo.CheckIsPrivacy(ctx, tx, req.FollowingID)
	if err != nil {
		return nil, err
	}

	if isPrivate {
		data, err := r.CreateFollowRequest(ctx, tx, req)
		if err != nil {
			return nil, err
		}
		return &types.RespondFollowRequest{
			RequestID: data,
		}, nil
	} else {

		followerID, err := r.CreateFollow(ctx, tx, &types.Follower{
			ID: utils.GenerateCustomID(utils.IDOptions{
				Prefix:       "FLW",
				NumberLength: 7,
			}),
			FollowerID:    req.FollowerID,
			FollowingID:   req.FollowingID,
			IsCloseFriend: isPrivate,
			IsMuted:       false,
			IsBlocked:     false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		})
		if err != nil {
			return nil, err
		}
		return &types.RespondFollowRequest{
			RequestID: followerID,
			Status:    "ACCEPTED",
		}, nil
	}
}
