package services

import (
	"cashflow_gin/api"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"time"

	"github.com/google/uuid"
)

type DashboardService interface {
	GetSummary(ctx context.Context, userID uuid.UUID, params api.GetDashboardSummaryParams) (*models.CashflowSummary, error)
}
type dashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(r repository.DashboardRepository) DashboardService {
	return &dashboardService{repo: r}
}
func (s *dashboardService) GetSummary(ctx context.Context, userID uuid.UUID, params api.GetDashboardSummaryParams) (*models.CashflowSummary, error) {
	walletId, err := uuid.Parse(params.WalletId)
	if err != nil {
		return nil, err
	}
	//now := time.Now()
	//endDate := now
	//
	//startDate := now.AddDate(0, 0, -30)

	parsedEnd := params.EndDate.Time
	endDate := time.Date(parsedEnd.Year(), parsedEnd.Month(), parsedEnd.Day(), 23, 59, 59, 0, parsedEnd.Location())

	parsedStart := params.StartDate.Time
	// Set ke jam 00:00:00
	startDate := time.Date(parsedStart.Year(), parsedStart.Month(), parsedStart.Day(), 0, 0, 0, 0, parsedStart.Location())

	// 4. Lempar rentang waktu yang udah absolut ke Repository
	return s.repo.GetDashboardSummary(ctx, userID, startDate, endDate, walletId)
}
