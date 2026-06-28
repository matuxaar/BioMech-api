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

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.ProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		// auto-create user if not found (skip flow, dev mode)
		user, err = s.userRepo.Create(ctx, userID, "dev@biomech.app")
		if err != nil {
			return nil, err
		}
	}

	deviceCount, err := s.userRepo.CountDevices(ctx, userID)
	if err != nil {
		deviceCount = 0
	}

	return &model.ProfileResponse{
		ID:          user.ID,
		Email:       user.Email,
		Nickname:    user.Nickname,
		DisplayName: user.DisplayName,
		PhotoURL:    user.PhotoURL,
		DeviceCount: deviceCount,
	}, nil
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *model.UpdateUserRequest) (*model.User, error) {
	if req.Nickname != nil && *req.Nickname != "" {
		existing, err := s.userRepo.FindByNickname(ctx, *req.Nickname)
		if err == nil && existing.ID != userID {
			return nil, errors.New("nickname already taken")
		}
	}
	return s.userRepo.UpdateProfile(ctx, userID, req)
}
