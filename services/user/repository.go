package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/chatting-app/services/shared/types"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// UserID      string     `db:"user_id" json:"userId"`
// 	Username    *string    `db:"username" json:"username"`
// 	PlaceBirth  *string    `db:"place_birth" json:"placeBirth"`
// 	DateBirth   *time.Time `db:"date_birth" json:"dateBirth"`
// 	Bio         *string    `db:"bio"  json:"bio"`
// 	IsPrivacy   bool       `db:"is_privacy" json:"isPrivacy"`
// 	PhoneNumber *string    `db:"phone_number" json:"phoneNumber"`
// 	Gender      *string    `db:"gender" json:"gender"`
// 	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAT"`

func (r *UserRepository) CreateUserProfile(ctx context.Context, req *types.UserProfile) (*types.UserProfile, error) {
	var profile types.UserProfile
	query := `
	INSERT INTO user_profile  
	(user_id,username,place_birth,date_birth,bio,is_privacy,phone_number,gender,updated_at)
	VALUES ($1,$2,$3,$4,$5,$6, $7, $8, $9)
	RETURNING
		user_id,username,place_birth,date_birth,bio,is_privacy,phone_number,gender,updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		req.UserID,
		req.Username,
		req.PlaceBirth,
		req.DateBirth,
		req.IsPrivacy,
		req.PhoneNumber,
		req.Gender,
		time.Now(),
	).Scan(
		&profile.UserID,
		&profile.Username,
		&profile.PlaceBirth,
		&profile.PlaceBirth,
		&profile.IsPrivacy,
		&profile.PhoneNumber,
		&profile.Gender,
		&profile.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create profile : %s", err)
	}

	return &profile, nil
}

func (r *UserRepository) CheckIsPrivacy(ctx context.Context, tx *sqlx.Tx, userID string) (bool, error) {
	var isPrivacy bool
	query := `
    SELECT is_privacy 
    FROM user_profile 
    WHERE user_id = $1
    `
	err := tx.QueryRowContext(ctx, query, userID).Scan(&isPrivacy)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows are found, assume the user profile does not exist or is public by default
			return false, fmt.Errorf("user profile not found for user_id: %s", userID)
		}
		// For other errors, return the error directly
		return false, fmt.Errorf("failed to check privacy setting: %w", err)
	}
	return isPrivacy, nil
}
