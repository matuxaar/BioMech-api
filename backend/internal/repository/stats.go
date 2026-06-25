package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matuxaar/BioMech-api/internal/model"
)

type StatsRepository struct {
	db *pgxpool.Pool
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) GetDashboardStats(ctx context.Context, userID string) (*model.DashboardStats, error) {
	stats := &model.DashboardStats{}

	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM devices WHERE user_id = $1`, userID).Scan(&stats.DeviceCount)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM training_jobs WHERE user_id = $1`, userID).Scan(&stats.TotalTrainings)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM training_jobs WHERE user_id = $1 AND status = 'completed'`, userID).Scan(&stats.CompletedTrainings)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, `SELECT COALESCE(AVG(accuracy), 0) FROM training_jobs WHERE user_id = $1 AND status = 'completed'`, userID).Scan(&stats.AverageAccuracy)
	if err != nil {
		return nil, err
	}

	stats.TopMovements = []string{"rest", "fist", "open", "pinch", "point"}

	return stats, nil
}
