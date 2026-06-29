package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matuxaar/BioMech-api/internal/model"
)

type TrainingFileRepository struct {
	db *pgxpool.Pool
}

func NewTrainingFileRepository(db *pgxpool.Pool) *TrainingFileRepository {
	return &TrainingFileRepository{db: db}
}

func (r *TrainingFileRepository) Create(ctx context.Context, userID, deviceID, originalName, filePath, label string, fileSize int64) (*model.TrainingFile, error) {
	f := &model.TrainingFile{
		ID:           uuid.New().String(),
		UserID:       userID,
		DeviceID:     deviceID,
		OriginalName: originalName,
		FilePath:     filePath,
		FileSize:     fileSize,
		Label:        label,
		CreatedAt:    time.Now(),
	}

	var devID any = deviceID
	if deviceID == "" {
		devID = nil
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO training_files (id, user_id, device_id, original_name, file_path, file_size, label, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		f.ID, f.UserID, devID, f.OriginalName, f.FilePath, f.FileSize, f.Label, f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *TrainingFileRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM training_files WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func (r *TrainingFileRepository) FindByUserID(ctx context.Context, userID string, page, limit int) ([]model.TrainingFile, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, COALESCE(device_id, ''), original_name, file_path, file_size, COALESCE(label, ''), created_at
		 FROM training_files WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []model.TrainingFile
	for rows.Next() {
		var f model.TrainingFile
		if err := rows.Scan(&f.ID, &f.UserID, &f.DeviceID, &f.OriginalName, &f.FilePath, &f.FileSize, &f.Label, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (r *TrainingFileRepository) FindByID(ctx context.Context, id string) (*model.TrainingFile, error) {
	f := &model.TrainingFile{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, COALESCE(device_id, ''), original_name, file_path, file_size, COALESCE(label, ''), created_at
		 FROM training_files WHERE id = $1`, id,
	).Scan(&f.ID, &f.UserID, &f.DeviceID, &f.OriginalName, &f.FilePath, &f.FileSize, &f.Label, &f.CreatedAt)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *TrainingFileRepository) Delete(ctx context.Context, id, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM training_files WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

func (r *TrainingFileRepository) UpdateFilePath(ctx context.Context, id, filePath string) error {
	_, err := r.db.Exec(ctx, `UPDATE training_files SET file_path = $1 WHERE id = $2`, filePath, id)
	return err
}
