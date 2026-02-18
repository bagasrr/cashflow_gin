package controllers

import (
	"cashflow_gin/dto/response"
	"cashflow_gin/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	wallets, err := c.services.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to get wallets",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "Wallets retrieved successfully",
		Errors:  nil,
		Data:    wallets,
	})
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
	reqId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid user ID",
			Errors:  "User ID not found in context",
			Data:    nil,
		})
		return
	}

	userID, err := uuid.Parse(reqId.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid user ID format",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}
	walletID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.BaseResponse{
			Status:  false,
			Message: "Invalid wallet ID format",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}
	groupIDParams, ok := ctx.GetQuery("groupid")
	var groupID uuid.UUID
	if !ok {
		// ctx.JSON(http.StatusBadRequest, response.BaseResponse{
		// 	Status:  false,
		// 	Message: "Cannot find groupid in query params",
		// 	Errors:  nil,
		// 	Data:    nil,
		// })
		// return
		groupID = uuid.Nil
	} else {
		groupID, err = uuid.Parse(groupIDParams)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.BaseResponse{
				Status:  false,
				Message: "Invalid group ID format",
				Errors:  err.Error(),
				Data:    nil,
			})
			return
		}
	}

	wallet, err := c.services.GetWalletByID(userID, walletID, groupID)
	if err != nil {
		// Handle the error
		ctx.JSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  false,
			Message: "Failed to get wallet",
			Errors:  err.Error(),
			Data:    nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, response.BaseResponse{
		Status:  true,
		Message: "Wallet retrieved successfully",
		Errors:  nil,
		Data:    wallet,
	})
}
