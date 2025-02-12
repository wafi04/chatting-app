package authrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/wafi04/chatting-app/services/shared/middleware"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
	"github.com/wafi04/chatting-app/services/shared/types"
	"github.com/wafi04/chatting-app/services/shared/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	DB     *sqlx.DB
	logger logger.Logger
}

func NewUserRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, req *types.CreateUserRequest) (types.CreateUserResponse, error) {
	userID := uuid.New().String()
	now := time.Now()

	hashPw, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		r.logger.Log(logger.ErrorLevel, "Failes Password : %v", err)
	}

	query := `
        INSERT INTO users (
            user_id, name, email, password_hash, 
            is_active, is_email_verified, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err = r.DB.ExecContext(
		ctx, query,
		userID, req.Name, req.Email, hashPw,
		true, false, now, now,
	)

	if err != nil {
		return types.CreateUserResponse{}, fmt.Errorf("failed to create verification token: %w", err)
	}

	access_token, err := middleware.GenerateToken(&types.UserInfo{
		UserId:          userID,
		Name:            req.Name,
		Email:           req.Email,
		IsEmailVerified: false,
	}, 24)
	if err != nil {
		return types.CreateUserResponse{}, fmt.Errorf("failed to generate tokens: %w", err)
	}
	refresh_token, err := middleware.GenerateToken(&types.UserInfo{
		UserId:          userID,
		Name:            req.Name,
		Email:           req.Email,
		IsEmailVerified: false,
	}, 168)
	if err != nil {
		return types.CreateUserResponse{}, fmt.Errorf("failed to generate tokens: %w", err)
	}

	session := types.Session{
		SessionId: utils.GenerateCustomID(utils.IDOptions{
			NumberLength: 10,
		}),
		UserId:         userID,
		AccessToken:    access_token,
		RefreshToken:   refresh_token,
		IpAddress:      req.IpAddress,
		DeviceInfo:     req.DeviceInfo,
		CreatedAt:      time.Now().Unix(),
		LastActivityAt: time.Now().Unix(),
		IsActive:       true,
		ExpiresAt:      time.Now().Unix(),
	}

	err = r.CreateSession(ctx, &session)
	if err != nil {
		return types.CreateUserResponse{}, fmt.Errorf("failed to create session: %w", err)
	}

	return types.CreateUserResponse{
		UserId:      userID,
		Name:        req.Name,
		Email:       req.Email,
		Picture:     req.Picture,
		AccessToken: "",
		SessionInfo: &types.Session{
			SessionId:  session.SessionId,
			DeviceInfo: session.DeviceInfo,
			IpAddress:  session.IpAddress,
		},
	}, nil

}

type dbUser struct {
	UserID          string
	Name            string
	Email           string
	Password        string
	Picture         string
	IsEmailVerified bool
	CreatedAt       int64
	UpdatedAt       int64
	LastLoginAt     int64
	IsActive        bool
}

func (r *AuthRepository) Login(ctx context.Context, login *types.LoginRequest) (*types.LoginResponse, error) {

	query := `
    SELECT
        user_id,
        name,
        email,
        password_hash,
        COALESCE(picture, ''),
        COALESCE(is_email_verified, false)::boolean,  
        EXTRACT(EPOCH FROM created_at)::bigint,
        EXTRACT(EPOCH FROM updated_at)::bigint,
        EXTRACT(EPOCH FROM COALESCE(last_login_at, created_at))::bigint,
        is_active::boolean
    FROM users
    WHERE name = $1
`

	var dbuser dbUser

	err := r.DB.QueryRowContext(ctx, query, login.Name).Scan(
		&dbuser.UserID,
		&dbuser.Name,
		&dbuser.Email,
		&dbuser.Password,
		&dbuser.Picture,
		&dbuser.IsEmailVerified,
		&dbuser.CreatedAt,
		&dbuser.UpdatedAt,
		&dbuser.LastLoginAt,
		&dbuser.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	userInfo := &types.UserInfo{
		UserId:          dbuser.UserID,
		Name:            dbuser.Name,
		Email:           dbuser.Email,
		IsEmailVerified: dbuser.IsEmailVerified,
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(login.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	access_token, err := middleware.GenerateToken(&types.UserInfo{
		UserId:          userInfo.UserId,
		Name:            userInfo.Name,
		Email:           userInfo.Email,
		IsEmailVerified: userInfo.IsEmailVerified,
	}, 24)
	if err != nil {
		return &types.LoginResponse{}, fmt.Errorf("failed to generate tokens: %w", err)
	}
	refresh_token, err := middleware.GenerateToken(&types.UserInfo{
		UserId:          userInfo.UserId,
		Name:            userInfo.Name,
		Email:           userInfo.Email,
		IsEmailVerified: userInfo.IsEmailVerified,
	}, 168)
	if err != nil {
		return &types.LoginResponse{}, fmt.Errorf("failed to generate tokens: %w", err)
	}
	query = `
        SELECT 
            session_id, 
            ip_address,
            device_info, 
            EXTRACT(EPOCH FROM created_at)::bigint, 
          	EXTRACT(EPOCH FROM last_activity_at)::bigint
        FROM sessions 
        WHERE user_id = $1 AND is_active = true AND device_info = $2
    `

	var existingSession types.Session
	err = r.DB.QueryRowContext(ctx, query, userInfo.UserId, login.DeviceInfo).Scan(
		&existingSession.SessionId,
		&existingSession.IpAddress,
		&existingSession.DeviceInfo,
		&existingSession.CreatedAt,
		&existingSession.LastActivityAt,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error checking existing session: %w", err)
	}

	if err == sql.ErrNoRows {
		existingSession = types.Session{
			SessionId: utils.GenerateCustomID(utils.IDOptions{
				NumberLength: 10,
			}),
			UserId:         userInfo.UserId,
			AccessToken:    access_token,
			RefreshToken:   refresh_token,
			IpAddress:      login.IpAddress,
			DeviceInfo:     login.DeviceInfo,
			CreatedAt:      time.Now().Unix(),
			LastActivityAt: time.Now().Unix(),
			IsActive:       true,
			ExpiresAt:      time.Now().Unix(),
		}

		err = r.CreateSession(ctx, &existingSession)
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %w", err)
		}
	}

	_, err = r.DB.ExecContext(
		ctx,
		"UPDATE users SET last_login_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE user_id = $1",
		userInfo.UserId,
	)
	if err != nil {
		r.logger.Log(logger.ErrorLevel, "Failed to update last login: %v", err)
	}

	return &types.LoginResponse{
		AccessToken:  access_token,
		UserId:       userInfo.UserId,
		RefreshToken: refresh_token,
		Session:      existingSession.SessionId,
		SessionInfo: &types.SessionInfo{
			SessionId:      existingSession.SessionId,
			DeviceInfo:     existingSession.DeviceInfo,
			IpAddress:      existingSession.IpAddress,
			CreatedAt:      existingSession.CreatedAt,
			LastActivityAt: existingSession.LastActivityAt,
		},
	}, nil
}

func (sr *AuthRepository) GetUser(ctx context.Context, req *types.GetUserRequest) (*types.GetUserResponse, error) {
	query := `
        SELECT 
            user_id, 
            name, 
            email,
            picture, 
            is_active, 
            is_email_verified,
            created_at, 
            updated_at, 
            last_login_at
        FROM users
        WHERE user_id = $1
    `
	sr.logger.Log(logger.InfoLevel, "data")

	user := &types.GetUserResponse{
		User: &types.UserInfo{},
	}

	var (
		isActive                          bool
		createdAt, updatedAt, lastLoginAt time.Time
		picture                           sql.NullString
	)
	err := sr.DB.QueryRowContext(ctx, query, req.UserId).Scan(
		&user.User.UserId,
		&user.User.Name,
		&user.User.Email,
		&picture,
		&isActive,
		&user.User.IsEmailVerified,
		&createdAt,
		&updatedAt,
		&lastLoginAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed : %s", err.Error())
	}

	if picture.Valid {
		user.User.Picture = picture.String
	}

	user.User.CreatedAt = createdAt.Unix()
	user.User.UpdatedAt = updatedAt.Unix()
	user.User.LastLoginAt = lastLoginAt.Unix()
	return user, nil
}

func (sr *AuthRepository) Logout(ctx context.Context, req *types.LogoutRequest) (*types.LogoutResponse, error) {
	query := `
	DELETE FROM sessions
    WHERE access_token = $1 AND user_id = $2
	`
	_, err := sr.DB.ExecContext(ctx, query, req.AccessToken, req.UserId)

	if err != nil {
		return nil, fmt.Errorf("failed : %s", err.Error())

	}

	return &types.LogoutResponse{
		Success: true,
	}, nil
}
