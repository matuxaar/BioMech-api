package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrNicknameTaken = errors.New("nickname already taken")
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
	existing, err := s.userRepo.FindByID(ctx, firebaseUID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("failed to lookup user during sync", "uid", firebaseUID, "error", err)
			return nil, err
		}
		return s.userRepo.Create(ctx, firebaseUID, email)
	}

	if existing.Email != email && email != "" {
		if err := s.userRepo.UpdateEmail(ctx, firebaseUID, email); err != nil {
			slog.Error("failed to update email during sync", "uid", firebaseUID, "error", err)
		}
	}
	return existing, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.ProfileResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
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
			return nil, ErrNicknameTaken
		}
	}
	return s.userRepo.UpdateProfile(ctx, userID, req)
}
