package models

import "time"

type WalletChartPoint struct {
	Date       time.Time `gorm:"column:tx_date"`
	Income     float64   `gorm:"column:total_income"`
	Expense    float64   `gorm:"column:total_expense"`
	Investment float64   `gorm:"column:total_investment"`
}
