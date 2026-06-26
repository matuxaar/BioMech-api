package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
)

type TrainingFileService struct {
	fileRepo *repository.TrainingFileRepository
}

func NewTrainingFileService(fileRepo *repository.TrainingFileRepository) *TrainingFileService {
	return &TrainingFileService{fileRepo: fileRepo}
}

func (s *TrainingFileService) Upload(ctx context.Context, userID, deviceID, label, originalName string, file io.Reader, fileSize int64) (*model.TrainingFile, error) {
	uploadDir := "uploads/training"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload dir: %w", err)
	}

	ext := filepath.Ext(originalName)
	storedName := fmt.Sprintf("%s_%s%s", userID[:8], uuid.New().String(), ext)
	filePath := filepath.Join(uploadDir, storedName)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return s.fileRepo.Create(ctx, userID, deviceID, originalName, filePath, label, written)
}

func (s *TrainingFileService) List(ctx context.Context, userID string) ([]model.TrainingFile, error) {
	return s.fileRepo.FindByUserID(ctx, userID)
}

func (s *TrainingFileService) Get(ctx context.Context, id string) (*model.TrainingFile, error) {
	return s.fileRepo.FindByID(ctx, id)
}

func (s *TrainingFileService) Delete(ctx context.Context, id, userID string) error {
	f, err := s.fileRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if f.UserID != userID {
		return ErrAccessDenied
	}
	os.Remove(f.FilePath)
	return s.fileRepo.Delete(ctx, id, userID)
}

func (s *TrainingFileService) GetFilePath(ctx context.Context, id string) (string, error) {
	f, err := s.fileRepo.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	return f.FilePath, nil
}

var _ = time.Now
