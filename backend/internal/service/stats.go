package service

import (
	"context"

	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/repository"
)

type StatsService struct {
	statsRepo *repository.StatsRepository
}

func NewStatsService(statsRepo *repository.StatsRepository) *StatsService {
	return &StatsService{statsRepo: statsRepo}
}

func (s *StatsService) GetDashboardStats(ctx context.Context, userID string) (*model.DashboardStats, error) {
	return s.statsRepo.GetDashboardStats(ctx, userID)
}
