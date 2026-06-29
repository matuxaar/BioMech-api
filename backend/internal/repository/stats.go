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

	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM devices WHERE user_id = $1`, userID).Scan(&stats.DeviceCount); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM training_jobs WHERE user_id = $1`, userID).Scan(&stats.TotalTrainings); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM training_jobs WHERE user_id = $1 AND status = 'completed'`, userID).Scan(&stats.CompletedTrainings); err != nil {
		return nil, err
	}
	if err := r.db.QueryRow(ctx, `SELECT COALESCE(AVG(accuracy), 0) FROM training_jobs WHERE user_id = $1 AND status = 'completed'`, userID).Scan(&stats.AverageAccuracy); err != nil {
		return nil, err
	}

	topMovements, err := r.db.Query(ctx,
		`SELECT DISTINCT da.name FROM device_actions da
		 JOIN devices d ON d.id = da.device_id
		 WHERE d.user_id = $1 ORDER BY da.name`, userID,
	)
	if err != nil {
		return stats, nil
	}
	defer topMovements.Close()
	for topMovements.Next() {
		var name string
		if err := topMovements.Scan(&name); err != nil {
			return nil, err
		}
		stats.TopMovements = append(stats.TopMovements, name)
	}

	return stats, nil
}
