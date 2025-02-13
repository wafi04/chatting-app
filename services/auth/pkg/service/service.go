package authservice

import (
	"context"
	"fmt"
	"log"
	"time"

	authrepository "github.com/wafi04/chatting-app/services/auth/pkg/repository"
	"github.com/wafi04/chatting-app/services/shared/types"
)

type AuthService struct {
	authRepo *authrepository.AuthRepository
}

func NewAuthService(authRepo *authrepository.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, req *types.CreateUserRequest) (*types.CreateUserResponse, error) {
	log.Printf("Received CreateUser request for user: %v", req)

	user, err := s.authRepo.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.CreateUserResponse{
		UserId:    user.UserId,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: time.Now().Unix(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *types.LoginRequest) (*types.LoginResponse, error) {
	log.Printf("Received Login request for user: %v", req)

	user, err := s.authRepo.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *AuthService) VerifyEmail(ctx context.Context, req *types.VerifyEmailRequest) (*types.VerifyEmailResponse, error) {
	log.Printf("Received verify email request for user: %v", req)

	user, err := s.authRepo.VerifyEmail(ctx, req)
	if err != nil {
		log.Fatalf("Failed  to verifiy email  :%v ", err)
		return nil, err
	}

	return &types.VerifyEmailResponse{
		Success: true,
		UserId:  user.UserId,
	}, nil
}
func (s *AuthService) ResendVerification(ctx context.Context, req *types.ResendVerificationRequest) (*types.ResendVerificationResponse, error) {
	log.Printf("Received verify email request for user: %v", req)

	user, err := s.authRepo.ResendVerification(ctx, req)
	if err != nil {
		log.Fatalf("Failed  to verifiy email  :%v ", err)
		return nil, err
	}

	return user, nil
}
func (s *AuthService) GetUser(ctx context.Context, req *types.GetUserRequest) (*types.UserInfo, error) {

	user, err := s.authRepo.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get Profile : %s ", err.Error())
	}

	return user, nil
}
func (s *AuthService) Logout(ctx context.Context, req *types.LogoutRequest) (*types.LogoutResponse, error) {

	user, err := s.authRepo.Logout(ctx, req)
	if err != nil {
		return &types.LogoutResponse{}, err
	}

	return user, nil
}
func (s *AuthService) RevokeSession(ctx context.Context, req *types.RevokeSessionRequest) (*types.RevokeSessionResponse, error) {

	user, err := s.authRepo.RevokeSession(ctx, req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *types.RefreshTokenRequest) (*types.RefreshTokenResponse, error) {
	log.Printf("Received verify email request for user: %v", req)

	refresh, err := s.authRepo.RefreshToken(ctx, req)
	if err != nil {
		log.Fatalf("Failed  to refresh token  :%v ", err)
		return nil, err
	}

	return refresh, nil
}

func (s *AuthService) ListSessions(ctx context.Context, req *types.ListSessionsRequest) (*types.ListSessionsResponse, error) {
	log.Printf("Received verify email request for user: %v", req)

	ListSessions, err := s.authRepo.ListSessions(ctx, req)
	if err != nil {
		log.Fatalf("Failed  to ListSessions token  :%v ", err)
		return nil, err
	}

	return ListSessions, nil
}

func (s *AuthService) GetSession(ctx context.Context, req *types.GetSessionRequest) (*types.GetSessionResponse, error) {
	log.Printf("Received Get Session for user: %v", req)

	GetSession, err := s.authRepo.GetSession(ctx, req)
	if err != nil {
		log.Fatalf("Failed  to GetSession token  :%v ", err)
		return nil, err
	}
	return GetSession, nil
}
