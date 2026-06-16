package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/services"
	"cashflow_gin/utils"
	"context"
)

type DashboardAPI struct {
	Service services.DashboardService
}

func (c *DashboardAPI) GetDashboardSummary(ctx context.Context, request api.GetDashboardSummaryRequestObject) (api.GetDashboardSummaryResponseObject, error) {
	// 1. Ambil UserID dari KTP (Middleware)
	userId, _, err := utils.GetUserInfo(ctx)
	if err != nil {
		return api.GetDashboardSummary200JSONResponse{
			Status:  utils.BoolPtr(false),
			Data:    nil,
			Message: utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	// 2. Panggil Service (Kirim parameter request untuk tanggal)
	summary, err := c.Service.GetSummary(ctx, userId, request.Params)
	if err != nil {
		return api.GetDashboardSummary200JSONResponse{
			Status:  utils.BoolPtr(false),
			Data:    nil,
			Message: utils.StringPtr("ERR : " + err.Error()),
		}, nil
	}

	// 3. KALKULASI MUTLAK DI MEMORI (Bukan di DB)
	totalCashflow := summary.TotalInflow - summary.TotalOutflow - summary.TotalInvestment

	// 4. Return ke OpenAPI Response
	return api.GetDashboardSummary200JSONResponse{
		Status:  utils.BoolPtr(true),
		Message: utils.StringPtr("Success get dashboard summary"),
		Data: &api.SummaryRes{
			TotalCashflow:   totalCashflow,
			TotalInflow:     summary.TotalInflow,
			TotalOutflow:    summary.TotalOutflow,
			TotalInvestment: summary.TotalInvestment,
		},
	}, nil
}
