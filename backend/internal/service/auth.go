package service

import (
	"context"
	"errors"

	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) SyncUser(ctx context.Context, firebaseUID, email string) (*model.User, error) {
	existing, _ := s.userRepo.FindByID(ctx, firebaseUID)
	if existing != nil {
		if existing.Email != email && email != "" {
			_ = s.userRepo.UpdateEmail(ctx, firebaseUID, email)
		}
		return existing, nil
	}

	return s.userRepo.Create(ctx, firebaseUID, email)
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error) {
	return s.userRepo.UpdateProfile(ctx, userID, req)
}
