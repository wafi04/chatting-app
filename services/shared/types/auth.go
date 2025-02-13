package types

import "time"

type User struct {
	UserId          string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name            string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email           string   `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	PasswordHash    string   `protobuf:"bytes,4,opt,name=password_hash,json=passwordHash,proto3" json:"password_hash,omitempty"`
	IsEmailVerified bool     `protobuf:"varint,5,opt,name=is_email_verified,json=isEmailVerified,proto3" json:"is_email_verified,omitempty"`
	IsActive        bool     `protobuf:"varint,6,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	CreatedAt       int64    `protobuf:"varint,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt       int64    `protobuf:"varint,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	LastLoginAt     int64    `protobuf:"varint,9,opt,name=last_login_at,json=lastLoginAt,proto3" json:"last_login_at,omitempty"`
	ActiveSessions  []string `protobuf:"bytes,10,rep,name=active_sessions,json=activeSessions,proto3" json:"active_sessions,omitempty"`
	Picture         string   `protobuf:"bytes,11,opt,name=picture,proto3" json:"picture,omitempty"`
}
type UserProfile struct {
	UserID      string     `db:"user_id" json:"userId"`
	Username    *string    `db:"username" json:"username"`
	PlaceBirth  *string    `db:"place_birth" json:"placeBirth"`
	DateBirth   *time.Time `db:"date_birth" json:"dateBirth"`
	Bio         *string    `db:"bio"  json:"bio"`
	IsPrivacy   bool       `db:"is_privacy" json:"isPrivacy"`
	PhoneNumber *string    `db:"phone_number" json:"phoneNumber"`
	Gender      *string    `db:"gender" json:"gender"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAT"`

	// Optional relationship
	User *User `db:"user,omitempty"`
}

type RequestPasswordResetRequest struct {
	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

type RequestPasswordResetResponse struct {
	Success    bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	ResetToken string `protobuf:"bytes,2,opt,name=reset_token,json=resetToken,proto3" json:"reset_token,omitempty"`
	ExpiresAt  int64  `protobuf:"varint,3,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

type ResetPasswordRequest struct {
	ResetToken  string `protobuf:"bytes,1,opt,name=reset_token,json=resetToken,proto3" json:"reset_token,omitempty"`
	OldPassword string `protobuf:"bytes,2,opt,name=old_password,json=oldPassword,proto3" json:"old_password,omitempty"`
	NewPassword string `protobuf:"bytes,3,opt,name=new_password,json=newPassword,proto3" json:"new_password,omitempty"`
}

type ResetPasswordResponse struct {
	Success   bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message   string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	UpdatedAt int64  `protobuf:"varint,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

type UserInfo struct {
	UserId          string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name            string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Email           string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	IsEmailVerified bool   `protobuf:"varint,4,opt,name=is_email_verified,json=isEmailVerified,proto3" json:"is_email_verified,omitempty"`
	CreatedAt       int64  `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt       int64  `protobuf:"varint,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	LastLoginAt     int64  `protobuf:"varint,7,opt,name=last_login_at,json=lastLoginAt,proto3" json:"last_login_at,omitempty"`
	Picture         string `protobuf:"bytes,8,opt,name=picture,proto3" json:"picture,omitempty"`
}

type CreateUserRequest struct {
	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Email      string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Password   string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	Picture    string `protobuf:"bytes,4,opt,name=picture,proto3" json:"picture,omitempty"`
	IpAddress  string `protobuf:"bytes,5,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	DeviceInfo string `protobuf:"bytes,6,opt,name=device_info,json=deviceInfo,proto3" json:"device_info,omitempty"`
}

type CreateUserResponse struct {
	UserId      string   `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AccessToken string   `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	Name        string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Email       string   `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	CreatedAt   int64    `protobuf:"varint,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Picture     string   `protobuf:"bytes,6,opt,name=picture,proto3" json:"picture,omitempty"`
	SessionInfo *Session `protobuf:"bytes,7,opt,name=session_info,json=sessionInfo,proto3,oneof" json:"session_info,omitempty"`
}

type UpdateUserRequest struct {
	UserId   string  `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name     *string `protobuf:"bytes,2,opt,name=name,proto3,oneof" json:"name,omitempty"`
	Email    *string `protobuf:"bytes,3,opt,name=email,proto3,oneof" json:"email,omitempty"`
	Password *string `protobuf:"bytes,4,opt,name=password,proto3,oneof" json:"password,omitempty"`
	Picture  *string `protobuf:"bytes,5,opt,name=picture,proto3,oneof" json:"picture,omitempty"`
}

type GetUserRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type GetUserResponse struct {
	User *UserInfo `protobuf:"bytes,1,opt,name=user,proto3,oneof" json:"user,omitempty"`
}

type Session struct {
	SessionId      string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	UserId         string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AccessToken    string `protobuf:"bytes,3,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken   string `protobuf:"bytes,4,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	DeviceInfo     string `protobuf:"bytes,5,opt,name=device_info,json=deviceInfo,proto3" json:"device_info,omitempty"`
	IpAddress      string `protobuf:"bytes,6,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	CreatedAt      int64  `protobuf:"varint,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	ExpiresAt      int64  `protobuf:"varint,8,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
	LastActivityAt int64  `protobuf:"varint,9,opt,name=last_activity_at,json=lastActivityAt,proto3" json:"last_activity_at,omitempty"`
	IsActive       bool   `protobuf:"varint,10,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
}

type VerificationToken struct {
	Token     string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	UserId    string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Type      string `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	CreatedAt int64  `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	ExpiresAt int64  `protobuf:"varint,5,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
	IsUsed    bool   `protobuf:"varint,6,opt,name=is_used,json=isUsed,proto3" json:"is_used,omitempty"`
}

type LoginRequest struct {
	Name       string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Password   string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	DeviceInfo string `protobuf:"bytes,3,opt,name=device_info,json=deviceInfo,proto3" json:"device_info,omitempty"`
	IpAddress  string `protobuf:"bytes,4,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
}

type LoginResponse struct {
	AccessToken  string       `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken string       `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	UserId       string       `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UserInfo     *UserInfo    `protobuf:"bytes,4,opt,name=user_info,json=userInfo,proto3" json:"user_info,omitempty"`
	SessionInfo  *SessionInfo `protobuf:"bytes,5,opt,name=session_info,json=sessionInfo,proto3" json:"session_info,omitempty"`
	Session      string       `protobuf:"bytes,6,opt,name=session,proto3" json:"session,omitempty"`
	ExpiresAt    int64        `protobuf:"varint,7,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

type SessionInfo struct {
	SessionId      string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	DeviceInfo     string `protobuf:"bytes,2,opt,name=device_info,json=deviceInfo,proto3" json:"device_info,omitempty"`
	IpAddress      string `protobuf:"bytes,3,opt,name=ip_address,json=ipAddress,proto3" json:"ip_address,omitempty"`
	CreatedAt      int64  `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	LastActivityAt int64  `protobuf:"varint,5,opt,name=last_activity_at,json=lastActivityAt,proto3" json:"last_activity_at,omitempty"`
}

type LogoutRequest struct {
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	SessionId   string `protobuf:"bytes,2,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	UserId      string `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type LogoutResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

type ValidateTokenRequest struct {
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

type ValidateTokenResponse struct {
	Valid     bool   `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
	UserId    string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	SessionId string `protobuf:"bytes,3,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	ExpiresAt int64  `protobuf:"varint,4,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	SessionId    string `protobuf:"bytes,2,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

type RefreshTokenResponse struct {
	AccessToken  string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	RefreshToken string `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	ExpiresAt    int64  `protobuf:"varint,3,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

type UpdateUserResponse struct {
	UserId    string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UpdatedAt int64  `protobuf:"varint,2,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

type VerifyEmailRequest struct {
	VerificationToken string `protobuf:"bytes,1,opt,name=verification_token,json=verificationToken,proto3" json:"verification_token,omitempty"`
	VerifyCode        string `protobuf:"bytes,2,opt,name=verify_code,json=verifyCode,proto3" json:"verify_code,omitempty"`
}

type VerifyEmailResponse struct {
	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	UserId  string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

type ResendVerificationRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Type   string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Token  string `protobuf:"bytes,3,opt,name=token,proto3" json:"token,omitempty"`
}

type ResendVerificationResponse struct {
	Success           bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	VerificationToken string `protobuf:"bytes,2,opt,name=verification_token,json=verificationToken,proto3" json:"verification_token,omitempty"`
	ExpiresAt         int64  `protobuf:"varint,3,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
	VerifyCode        string `protobuf:"bytes,4,opt,name=verify_code,json=verifyCode,proto3" json:"verify_code,omitempty"`
}

type GetSessionRequest struct {
	SessionId string `protobuf:"bytes,1,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

type GetSessionResponse struct {
	SessionInfo *SessionInfo `protobuf:"bytes,1,opt,name=session_info,json=sessionInfo,proto3" json:"session_info,omitempty"`
	UserInfo    *User        `protobuf:"bytes,2,opt,name=user_info,json=userInfo,proto3" json:"user_info,omitempty"`
}

type RevokeSessionRequest struct {
	UserId    string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	SessionId string `protobuf:"bytes,2,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
}

type RevokeSessionResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

type ListSessionsRequest struct {
	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

type ListSessionsResponse struct {
	Sessions []*SessionInfo `protobuf:"bytes,1,rep,name=sessions,proto3" json:"sessions,omitempty"`
}
