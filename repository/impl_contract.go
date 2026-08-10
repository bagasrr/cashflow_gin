package repository

import (
	"cashflow_gin/utils"
	"context"

	"gorm.io/gorm"
)

type gormTxManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) *gormTxManager {
	return &gormTxManager{db: db}
}

func (m *gormTxManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Bungkus 'tx' ke dalam context baru
		txCtx := utils.InjectTx(ctx, tx)
		// Jalankan fungsi bisnis lu dengan context yang sudah berisi transaksi
		return fn(txCtx)
	})
}
