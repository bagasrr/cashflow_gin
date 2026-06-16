package repository

import (
	"cashflow_gin/models"
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetDashboardSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, walletId uuid.UUID) (*models.CashflowSummary, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetDashboardSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, walletId uuid.UUID) (*models.CashflowSummary, error) {
	var summary models.CashflowSummary

	// Tambahkan kondisi rentang waktu (BETWEEN)
	err := r.db.WithContext(ctx).
		Table("transactions").
		Select(`
            COALESCE(SUM(CASE WHEN categories.type = 'INCOME' THEN transactions.amount ELSE 0 END), 0) as total_inflow,
            COALESCE(SUM(CASE WHEN categories.type = 'EXPENSE' THEN transactions.amount ELSE 0 END), 0) as total_outflow,
            COALESCE(SUM(CASE WHEN categories.type = 'INVESTMENT' THEN transactions.amount ELSE 0 END), 0) as total_investment
        `).
		Joins("JOIN categories ON categories.id = transactions.category_id").
		Where("transactions.user_id = ? AND transactions.deleted_at IS NULL and transactions.wallet_id = ?", userID, walletId).

		// EKSEKUSI MUTLAK: Kunci gembok waktunya di sini
		Where("transactions.date >= ? AND transactions.date <= ?", startDate, endDate).
		Scan(&summary).Error

	return &summary, err
}
