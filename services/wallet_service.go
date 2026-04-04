package services

import (
	"cashflow_gin/dto/response"
	"cashflow_gin/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type WalletService interface {
	GetAll(ctx context.Context) (*[]response.WalletResponse, error)
	GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (*response.WalletResponse, error)
	GetMine(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]response.WalletResponse, error)
}

type walletService struct {
	// Kita butuh WalletRepository untuk akses data wallet
	walletRepo repository.WalletRepository
	groupRepo  repository.GroupRepository
}

func NewWalletService(wRepo repository.WalletRepository, gRepo repository.GroupRepository) WalletService {
	return &walletService{walletRepo: wRepo, groupRepo: gRepo}
}

func (s *walletService) GetAll(ctx context.Context) (*[]response.WalletResponse, error) {
	limit := 10
	offset := 0
	// nanti ganti jadi dinamis
	wallet, err := s.walletRepo.FindAll(ctx, limit, offset)
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
			ID:               w.ID,
			Name:             w.Name,
			Balance:          w.Balance,
			Transactions:     transactions,
			TransactionCount: w.TransactionCount,
		})
	}
	return &responses, nil
}

func (s *walletService) GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (*response.WalletResponse, error) { // groupID DIBUANG

	// 1. FETCH DULUAN. Biarkan database yang berbicara ini dompet apa.
	wallet, err := s.walletRepo.FindByID(ctx, walletID)
	if err != nil {
		return &response.WalletResponse{}, err
	}

	// 2. INTEROGASI WUJUD DOMPETNYA
	if wallet.GroupID != nil {
		// SKENARIO A: INI DOMPET GRUP
		isMember, err := s.groupRepo.IsGroupMember(ctx, *wallet.GroupID, userID)
		if err != nil || !isMember {
			return &response.WalletResponse{}, errors.New("unauthorized: you are not a member of this group")
		}
	} else if wallet.UserID != nil {
		// SKENARIO B: INI DOMPET PERSONAL
		// Perhatikan tanda '*'. Kita membandingkan VALUE-nya, bukan alamat memorinya.
		if *wallet.UserID != userID {
			return &response.WalletResponse{}, errors.New("unauthorized: wallet does not belong to you")
		}
	} else {
		// Data cacat di database (tidak punya pemilik)
		return &response.WalletResponse{}, errors.New("internal error: wallet has no owner")
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
		ID:               wallet.ID,
		Name:             wallet.Name,
		Balance:          wallet.Balance,
		Transactions:     transactions,
		TransactionCount: wallet.TransactionCount,
	}
	return &response, nil
}

// services/wallet_service.go
func (s *walletService) GetMine(ctx context.Context, userID uuid.UUID, page, limit int) (*[]response.WalletResponse, error) {
	// 1. Validasi cegah angka minus atau nol (Hacker/Bug Frontend)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 { // Jangan biarkan frontend request 1 juta data sekaligus
		limit = 10
	}

	// 2. Hitung Offset
	offset := (page - 1) * limit

	// 3. Panggil Repository dengan parameter yang udah matang
	wallets, err := s.walletRepo.FindAllMine(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var res []response.WalletResponse

	for _, w := range *wallets {
		res = append(res, response.WalletResponse{
			ID:               w.ID,
			Name:             w.Name,
			Balance:          w.Balance,
			TransactionCount: w.TransactionCount,
		})
	}

	return &res, nil
}
