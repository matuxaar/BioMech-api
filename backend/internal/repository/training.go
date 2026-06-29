package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matuxaar/BioMech-api/internal/model"
)

var ErrTrainingJobNotFound = errors.New("training job not found")

type TrainingRepository struct {
	db *pgxpool.Pool
}

func NewTrainingRepository(db *pgxpool.Pool) *TrainingRepository {
	return &TrainingRepository{db: db}
}

func (r *TrainingRepository) Create(ctx context.Context, userID string, req *model.CreateTrainingJobRequest) (*model.TrainingJob, error) {
	job := &model.TrainingJob{
		ID:         uuid.New().String(),
		UserID:     userID,
		SessionIDs: req.SessionIDs,
		Status:     model.TrainingStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	_, err := r.db.Exec(ctx,
		`INSERT INTO training_jobs (id, user_id, session_ids, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		job.ID, job.UserID, job.SessionIDs, job.Status, job.CreatedAt, job.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (r *TrainingRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM training_jobs WHERE user_id = $1`, userID).Scan(&count)
	return count, err
}

func (r *TrainingRepository) FindByUserID(ctx context.Context, userID string, page, limit int) ([]model.TrainingJob, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, session_ids, status, COALESCE(model_path, ''), COALESCE(accuracy, 0), COALESCE(error_message, ''), created_at, updated_at
		 FROM training_jobs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []model.TrainingJob
	for rows.Next() {
		var j model.TrainingJob
		if err := rows.Scan(&j.ID, &j.UserID, &j.SessionIDs, &j.Status,
			&j.ModelPath, &j.Accuracy, &j.ErrorMessage, &j.CreatedAt, &j.UpdatedAt,
		); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func (r *TrainingRepository) FindByID(ctx context.Context, id string) (*model.TrainingJob, error) {
	j := &model.TrainingJob{}
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, session_ids, status, COALESCE(model_path, ''), COALESCE(accuracy, 0), COALESCE(error_message, ''), created_at, updated_at
		 FROM training_jobs WHERE id = $1`, id,
	).Scan(&j.ID, &j.UserID, &j.SessionIDs, &j.Status,
		&j.ModelPath, &j.Accuracy, &j.ErrorMessage, &j.CreatedAt, &j.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (r *TrainingRepository) UpdateStatus(ctx context.Context, id string, status model.TrainingStatus, modelPath string, accuracy float64, errMsg string) error {
	ct, err := r.db.Exec(ctx,
		`UPDATE training_jobs SET status = $1, model_path = $2, accuracy = $3, error_message = $4, updated_at = $5
		 WHERE id = $6`,
		status, modelPath, accuracy, errMsg, time.Now(), id,
	)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
