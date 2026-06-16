package models

type CashflowSummary struct {
	TotalInflow     float64 `gorm:"column:total_inflow"`
	TotalOutflow    float64 `gorm:"column:total_outflow"`
	TotalInvestment float64 `gorm:"column:total_investment"`
}
