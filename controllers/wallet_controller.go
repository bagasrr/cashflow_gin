package controllers

import (
	"cashflow_gin/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WalletController struct {
	services services.WalletService
}

func NewWalletController(s services.WalletService) *WalletController {
	return &WalletController{services: s}
}

// GetAllWallets godoc
// @Summary      Get All Wallets
// @Description  Mendapatkan semua dompet yang dimiliki pengguna, termasuk dompet pribadi dan grup.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Success      200 {object} response.BaseResponse{data=[]response.WalletResponse}
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /wallets [get]
func (c *WalletController) GetAllWallets(ctx *gin.Context) {
	// Implementasi untuk mendapatkan semua wallet
	wallets, err := c.services.GetAll(ctx.Request.Context())
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to get wallets", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Wallets retrieved successfully", wallets)
}

// GetWalletByID godoc
// @Summary      Get Wallet By ID
// @Description  Mendapatkan detail dompet berdasarkan ID, termasuk transaksi terkait.
// @Tags         Wallets
// @Accept       json
// @Produce      json
// @Param        id path string true "Wallet ID"
// @Param        groupid query string true "Group ID"
// @Success      200 {object} response.BaseResponse{data=response.WalletResponse}
// @Failure      400 {object} response.BaseResponse
// @Failure      500 {object} response.BaseResponse
// @Security 	 BearerAuth
// @Router       /wallets/{id} [get]
func (c *WalletController) GetWalletByID(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, "Invalid user ID", err)
		return
	}

	walletID, err := GetParamID(ctx, "id")
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid wallet ID format", err)
		return
	}

	wallet, err := c.services.GetWalletByID(ctx.Request.Context(), userID, walletID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to get wallet", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Wallet retrieved successfully", wallet)
}

func (c *WalletController) GetMine(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		SendError(ctx, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	wallets, err := c.services.GetMine(ctx.Request.Context(), userID, page, limit)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to retrieve wallets", err)
		return
	}

	SendSuccess(ctx, http.StatusOK, "Wallets retrieved successfully", wallets)
}
