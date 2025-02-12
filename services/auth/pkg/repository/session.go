package authrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wafi04/chatting-app/services/shared/middleware"
	"github.com/wafi04/chatting-app/services/shared/pkg/logger"
	"github.com/wafi04/chatting-app/services/shared/types"
)

func (sr *AuthRepository) RevokeSession(ctx context.Context, req *types.RevokeSessionRequest) (*types.RevokeSessionResponse, error) {
	sr.logger.Log(logger.InfoLevel, "Recieved  Session Request ")

	query := `
	DELETE FROM sessions
    WHERE session_id = $1 AND user_id = $2
	`
	_, err := sr.DB.ExecContext(ctx, query, req.SessionId, req.UserId)

	if err != nil {
		sr.logger.Log(logger.ErrorLevel, "Failed to Delete Session : %v", err)
		return nil, nil
	}

	return &types.RevokeSessionResponse{
		Success: true}, nil
}
func (sr *AuthRepository) CreateSession(ctx context.Context, session *types.Session) error {
	// Pisahkan query insert dan update
	insertQuery := `
       INSERT INTO sessions (
           session_id, 
           user_id, 
           access_token, 
           refresh_token, 
           ip_address, 
           device_info, 
           is_active, 
           expires_at, 
           last_activity_at, 
           created_at
       ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
   `

	updateQuery := `
       UPDATE sessions 
       SET 
           access_token = $1, 
           refresh_token = $2, 
           ip_address = $3, 
           last_activity_at = $4
       WHERE user_id = $5 AND device_info = $6
   `

	if session.SessionId == "" {
		session.SessionId = uuid.New().String()
	}

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	// Pertama, coba insert
	_, err := sr.DB.ExecContext(
		ctx,
		insertQuery,
		session.SessionId,
		session.UserId,
		session.AccessToken,
		session.RefreshToken,
		session.IpAddress,
		session.DeviceInfo,
		true,
		expiresAt,
		now,
		now,
	)

	// Jika insert gagal (misal duplicate), lakukan update
	if err != nil {
		_, err = sr.DB.ExecContext(
			ctx,
			updateQuery,
			session.AccessToken,
			session.RefreshToken,
			session.IpAddress,
			now,
			session.UserId,
			session.DeviceInfo,
		)
	}

	if err != nil {
		sr.logger.WithError(err).WithFields(map[string]interface{}{
			"user_id":    session.UserId,
			"session_id": session.SessionId,
		}).Error("Failed to create/update session")
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}
func (sr *AuthRepository) RefreshToken(ctx context.Context, req *types.RefreshTokenRequest) (*types.RefreshTokenResponse, error) {
	sr.logger.Log(logger.InfoLevel, "Refresh Token Incoming")

	query := `
        SELECT u.user_id, u.email, u.role, u.is_email_verified
        FROM sessions s
        JOIN users u ON s.user_id = u.user_id
        WHERE s.session_id = $1
    `
	var user types.UserInfo
	err := sr.DB.QueryRowContext(ctx, query, req.SessionId).Scan(
		&user.UserId,
		&user.Email,
		&user.IsEmailVerified,
	)
	if err != nil {
		sr.logger.Log(logger.ErrorLevel, "Failed to retrieve user: %v", err)
		return nil, err
	}

	updateQuery := `
        UPDATE sessions SET access_token = $1 WHERE session_id = $2
    `
	_, err = sr.DB.ExecContext(ctx, updateQuery, req.RefreshToken, req.SessionId)
	if err != nil {
		sr.logger.Log(logger.ErrorLevel, "Failed to Refresh Token: %v", err)
		return nil, err
	}
	access_token, err := middleware.GenerateToken(&user, 24)
	if err != nil {
		return nil, err
	}
	refresh_token, err := middleware.GenerateToken(&user, 168)
	if err != nil {
		return nil, err
	}

	return &types.RefreshTokenResponse{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(),
	}, nil
}

func (sr *AuthRepository) ListSessions(ctx context.Context, req *types.ListSessionsRequest) (*types.ListSessionsResponse, error) {
	query := `
        SELECT 
            session_id,
            device_info,
            ip_address,
            EXTRACT(EPOCH FROM created_at)::bigint AS created_at,
            EXTRACT(EPOCH FROM last_activity_at)::bigint AS last_activity_at
        FROM sessions
        WHERE user_id = $1
    `

	rows, err := sr.DB.QueryContext(ctx, query, req.UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*types.SessionInfo
	for rows.Next() {
		session := &types.SessionInfo{}
		err := rows.Scan(
			&session.SessionId,
			&session.DeviceInfo,
			&session.IpAddress,
			&session.CreatedAt,
			&session.LastActivityAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &types.ListSessionsResponse{
		Sessions: sessions,
	}, nil
}
func (sr *AuthRepository) GetSession(ctx context.Context, req *types.GetSessionRequest) (*types.GetSessionResponse, error) {
	var created_at, updated_at time.Time
	var session types.Session
	var user types.User
	var picture sql.NullString

	query := `
    SELECT 
        u.user_id,
        u.name,
        u.email,
		u.picture,
        s.session_id,
        s.device_info,
        s.ip_address,
        s.created_at,
        s.updated_at
    FROM 
        sessions s
    JOIN 
        users u
    ON 
        s.user_id = u.user_id
    WHERE 
        s.session_id = $1;
    `

	err := sr.DB.QueryRowContext(ctx, query, req.SessionId).Scan(
		&session.UserId,
		&user.Name,
		&user.Email,
		&picture,
		&session.SessionId,
		&session.DeviceInfo,
		&session.IpAddress,
		&created_at,
		&updated_at,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	session.CreatedAt = created_at.Unix()
	session.LastActivityAt = updated_at.Unix()
	user.Picture = ""
	if picture.Valid {
		user.Picture = picture.String
	}

	return &types.GetSessionResponse{
		SessionInfo: &types.SessionInfo{
			SessionId:      session.SessionId,
			DeviceInfo:     session.DeviceInfo,
			IpAddress:      session.IpAddress,
			CreatedAt:      session.CreatedAt,
			LastActivityAt: session.LastActivityAt,
		},
		UserInfo: &types.User{
			UserId:  session.UserId,
			Name:    user.Name,
			Email:   user.Email,
			Picture: user.Picture,
		},
	}, nil
}
