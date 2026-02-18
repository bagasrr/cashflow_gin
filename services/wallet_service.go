package services

import (
	"cashflow_gin/dto/response"
	"cashflow_gin/repository"
	"errors"

	"github.com/google/uuid"
)

type WalletService interface {
	GetAll() ([]response.WalletResponse, error)
	GetWalletByID(userID, walletID, groupID uuid.UUID) (response.WalletResponse, error)
}

type walletService struct {
	// Kita butuh WalletRepository untuk akses data wallet
	walletRepo repository.WalletRepository
	groupRepo  repository.GroupRepository
}

func NewWalletService(wRepo repository.WalletRepository, gRepo repository.GroupRepository) WalletService {
	return &walletService{walletRepo: wRepo, groupRepo: gRepo}
}

func (s *walletService) GetAll() ([]response.WalletResponse, error) {
	// Implementasi untuk mendapatkan semua wallet
	wallet, err := s.walletRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []response.WalletResponse
	for _, w := range *wallet {
		var transactions []response.TransactionResponse
		for _, t := range w.Transactions {
			transactions = append(transactions, response.TransactionResponse{
				ID:          t.ID.String(),
				Title:       t.Title,
				Amount:      t.Amount,
				Date:        t.Date,
				Description: t.Description,
				User: response.UserResponse{
					// ID:       t.User.ID.String(),
					Username: t.User.Username,
					Email:    t.User.Email,
					// UserRole: t.User.UserRole.String(),
				},
				Category: response.CategoryResponse{
					Name: t.Category.Name,
					Type: t.Category.Type,
				},
			})
		}
		responses = append(responses, response.WalletResponse{
			ID:           w.ID,
			Name:         w.Name,
			Balance:      w.Balance,
			Transactions: transactions,
		})
	}
	return responses, nil
}

func (s *walletService) GetWalletByID(userID, walletID, groupID uuid.UUID) (response.WalletResponse, error) {
	// Cek Apakah Group id nya Nil atau bukan
	if groupID != uuid.Nil {
		isMember, err := s.groupRepo.IsGroupMember(groupID, userID)
		if err != nil || !isMember {
			return response.WalletResponse{}, err
		}
	}

	wallet, err := s.walletRepo.FindByID(walletID)
	if err != nil {
		return response.WalletResponse{}, err
	}

	// cek apakah si wallet itu milik user atau grup yang dia ikuti
	if wallet.UserID != &userID || (groupID != uuid.Nil && wallet.GroupID != &groupID) {
		return response.WalletResponse{}, errors.New("unauthorized: wallet does not belong to user")
	}

	var transactions []response.TransactionResponse
	for _, t := range wallet.Transactions {
		transactions = append(transactions, response.TransactionResponse{
			ID:          t.ID.String(),
			Title:       t.Title,
			Amount:      t.Amount,
			Date:        t.Date,
			Description: t.Description,
			User: response.UserResponse{
				ID:       t.User.ID.String(),
				Username: t.User.Username,
				Email:    t.User.Email,
				UserRole: t.User.UserRole.String(),
			},
			Category: response.CategoryResponse{
				Name: t.Category.Name,
				Type: t.Category.Type,
			},
		})
	}

	response := response.WalletResponse{
		ID:           wallet.ID,
		Name:         wallet.Name,
		Balance:      wallet.Balance,
		Transactions: transactions,
	}
	return response, nil
}
