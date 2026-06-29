package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
)

type TrainingFileService struct {
	fileRepo  *repository.TrainingFileRepository
	uploadDir string
}

func NewTrainingFileService(fileRepo *repository.TrainingFileRepository, uploadDir string) *TrainingFileService {
	return &TrainingFileService{fileRepo: fileRepo, uploadDir: uploadDir}
}

func (s *TrainingFileService) Upload(ctx context.Context, userID, deviceID, label, originalName string, file io.Reader, fileSize int64) (*model.TrainingFile, error) {
	uploadDir := s.uploadDir
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload dir: %w", err)
	}

	ext := filepath.Ext(originalName)
	uid := userID
	if len(uid) > 8 {
		uid = uid[:8]
	}
	storedName := fmt.Sprintf("%s_%s%s", uid, uuid.New().String(), ext)
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

func (s *TrainingFileService) List(ctx context.Context, userID string, page, limit int) (*model.PaginatedResponse[model.TrainingFile], error) {
	total, err := s.fileRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	files, err := s.fileRepo.FindByUserID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}
	result := model.NewPaginatedResponse(files, total, page, limit)
	return &result, nil
}

func (s *TrainingFileService) Get(ctx context.Context, id, userID string) (*model.TrainingFile, error) {
	f, err := s.fileRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if f.UserID != userID {
		return nil, ErrAccessDenied
	}
	return f, nil
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
