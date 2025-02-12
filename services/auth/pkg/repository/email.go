package authrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/wafi04/chatting-app/services/shared/middleware"
	"github.com/wafi04/chatting-app/services/shared/types"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthRepository) ResendVerification(ctx context.Context, req *types.ResendVerificationRequest) (*types.ResendVerificationResponse, error) {
	// user, err := s.GetUser(ctx, &types.GetUserRequest{
	// 	UserId: req.UserId,
	// })
	// if err != nil {
	// 	s.logger.Log(logger.ErrorLevel, "Unauthorized : %v", err)
	// 	return nil, fmt.Errorf("Unauthorized : %v", err)
	// }

	// verifyCode := logger.GenerateVerificationCode()
	expiresAt := time.Now().Add(1 * time.Hour)

	// appPW := logger.LoadEnv("APP_PASSWORD")
	// cleanPassword := strings.ReplaceAll(appPW, " ", "")
	// emailSender := mailer.NewEmailSender(
	// 	"smtp.gmail.com",
	// 	587,
	// 	"wafiq3040@gmail.com",
	// 	cleanPassword,
	// )

	// toEmail := user.User.Email

	// if err := emailSender.SendVerificationEmail(toEmail, user.User.Name, 0); err != nil {
	// 	return nil, fmt.Errorf("failed to send email : %w", err)
	// }

	// query := `
	//     INSERT INTO verification_tokens (
	//         token,
	//         user_id,
	//         verify_code,
	//         token_type,
	//         expires_at
	//     ) VALUES ($1, $2, $3, $4, $5)
	// `

	// _, err = s.DB.ExecContext(ctx, query,
	// 	req.Token,
	// 	req.UserId,
	// 	verifyCode,
	// 	"EMAIL_VERIFICATION",
	// 	expiresAt,
	// )

	// if err != nil {
	// 	return nil, fmt.Errorf("failed to  verification token: %w", err)
	// }

	return &types.ResendVerificationResponse{
		VerificationToken: req.Token,
		// VerifyCode:        verifyCode,
		Success:   true,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

func (s *AuthRepository) VerifyEmail(ctx context.Context, req *types.VerifyEmailRequest) (*types.VerifyEmailResponse, error) {
	query := `
    SELECT user_id, expires_at  
    FROM verification_tokens 
    WHERE token = $1 
    AND verify_code = $2 
    AND is_used = false 
    AND expires_at > NOW()
`
	var (
		userId    string
		expiresAt time.Time
	)

	err := s.DB.QueryRowContext(ctx, query, req.VerificationToken, req.VerifyCode).Scan(&userId, &expiresAt)
	if err != nil {
		return nil, fmt.Errorf("verification failed: %w", err)
	}

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("transaction error: %w", err)
	}

	updateQuery := `
        UPDATE users 
        SET is_email_verified = true, updated_at = NOW() 
        WHERE user_id = $1
    `

	markUsedQuery := `
        UPDATE verification_tokens 
        SET is_used = true 
        WHERE token = $1
    `

	_, err = tx.ExecContext(ctx, updateQuery, userId)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("update user error: %w", err)
	}

	_, err = tx.ExecContext(ctx, markUsedQuery, req.VerificationToken)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("mark token error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit error: %w", err)
	}

	return &types.VerifyEmailResponse{
		Success: true,
		UserId:  userId,
		Message: "Email verified successfully",
	}, nil
}
func (s *AuthRepository) RequestPasswordReset(ctx context.Context, req *types.RequestPasswordResetRequest) (*types.RequestPasswordResetResponse, error) {
	var user types.User
	checkUserQuery := `SELECT user_id,name,email,is_email_verified,is_active  FROM users WHERE email = $1`
	err := s.DB.QueryRowContext(ctx, checkUserQuery, req.Email).Scan(&user.UserId, user.Name, user.Email, user.IsEmailVerified, user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.RequestPasswordResetResponse{
				Success: false,
			}, nil
		}
		return nil, fmt.Errorf("failed to check user: %v", err)
	}

	token, err := middleware.GenerateToken(&types.UserInfo{
		UserId:          user.UserId,
		Name:            user.Name,
		Email:           user.Email,
		IsEmailVerified: user.IsEmailVerified,
	}, 1)

	if err != nil {
		return nil, err
	}

	insertTokenQuery := `
        INSERT INTO verification_tokens 
        (token, user_id, token_type, expires_at) 
        VALUES ($1, $2, 'PASSWORD', CURRENT_TIMESTAMP + INTERVAL '1 hour')`

	_, err = s.DB.ExecContext(ctx, insertTokenQuery, token, user.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to create reset token: %v", err)
	}

	return &types.RequestPasswordResetResponse{
		Success:    true,
		ResetToken: token,
		ExpiresAt:  time.Now().Add(1 * time.Hour).Unix(),
	}, nil
}

func (s *AuthRepository) ResetPassword(ctx context.Context, req *types.ResetPasswordRequest) (*types.ResetPasswordResponse, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	verifyTokenQuery := `
        SELECT user_id 
        FROM verification_tokens 
        WHERE token = $1 
        AND token_type = 'password_reset'
        AND expires_at > CURRENT_TIMESTAMP 
        AND is_used = false`

	var userID string
	err = tx.QueryRowContext(ctx, verifyTokenQuery, req.ResetToken).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &types.ResetPasswordResponse{
				Success: false,
				Message: "Invalid or expired reset token",
			}, nil
		}
		return nil, fmt.Errorf("failed to verify reset token: %v", err)
	}

	updatePasswordQuery := `
        UPDATE users 
        SET password_hash = $1, 
            updated_at = CURRENT_TIMESTAMP 
        WHERE user_id = $2 
        RETURNING extract(epoch from updated_at)::bigint`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	var updatedAt int64
	err = tx.QueryRowContext(ctx, updatePasswordQuery, hashedPassword, userID).Scan(&updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update password: %v", err)
	}

	markTokenUsedQuery := `
        UPDATE verification_tokens 
        SET is_used = true 
        WHERE token = $1`

	_, err = tx.ExecContext(ctx, markTokenUsedQuery, req.ResetToken)
	if err != nil {
		return nil, fmt.Errorf("failed to mark token as used: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &types.ResetPasswordResponse{
		Success:   true,
		Message:   "Password successfully reset",
		UpdatedAt: updatedAt,
	}, nil
}
