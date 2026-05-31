package services

import (
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type WalletService interface {
	CreatePersonalWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error)
	CreateGroupWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error)
	GetAll(ctx context.Context) (*[]models.Wallet, error)
	GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (*models.Wallet, error)
	GetMine(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Wallet, error)

	DeleteWallet(ctx context.Context, walletId, userId uuid.UUID) error
}

type walletService struct {
	// Kita butuh WalletRepository untuk akses data wallet
	walletRepo repository.WalletRepository
	groupRepo  repository.GroupRepository
}

func NewWalletService(wRepo repository.WalletRepository, gRepo repository.GroupRepository) WalletService {
	return &walletService{walletRepo: wRepo, groupRepo: gRepo}
}

func (s *walletService) CreatePersonalWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error) {

	createdWallet, err := s.walletRepo.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return createdWallet, nil

}

func (s *walletService) CreateGroupWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error) {

	createdWallet, err := s.walletRepo.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}
	return createdWallet, nil

}

func (s *walletService) GetAll(ctx context.Context) (*[]models.Wallet, error) {
	limit := 10
	offset := 0
	// nanti ganti jadi dinamis
	wallet, err := s.walletRepo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func (s *walletService) GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (*models.Wallet, error) { // groupID DIBUANG

	// 1. FETCH DULUAN. Biarkan database yang berbicara ini dompet apa.
	wallet, err := s.walletRepo.FindByID(ctx, walletID)
	if err != nil {
		return &models.Wallet{}, err
	}

	// 2. INTEROGASI WUJUD DOMPETNYA
	if wallet.GroupID != nil {
		// SKENARIO A: INI DOMPET GRUP
		isMember, err := s.groupRepo.IsGroupMember(ctx, *wallet.GroupID, userID)
		if err != nil || !isMember {
			return &models.Wallet{}, errors.New("unauthorized: you are not a member of this group")
		}
	} else if wallet.UserID != nil {
		// SKENARIO B: INI DOMPET PERSONAL
		// Perhatikan tanda '*'. Kita membandingkan VALUE-nya, bukan alamat memorinya.
		if *wallet.UserID != userID {
			return &models.Wallet{}, errors.New("unauthorized: wallet does not belong to you")
		}
	} else {
		// Data cacat di database (tidak punya pemilik)
		return &models.Wallet{}, errors.New("internal error: wallet has no owner")
	}

	return wallet, nil
}

// services/wallet_service.go
func (s *walletService) GetMine(ctx context.Context, userID uuid.UUID, page, limit int) (*[]models.Wallet, error) {
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

	return wallets, nil
}

func (s *walletService) DeleteWallet(ctx context.Context, walletId, userId uuid.UUID) error {
	isGroupWallet, groupId, err := s.groupRepo.IsGroupWallet(ctx, walletId)
	if err != nil {
		return err
	}
	if isGroupWallet {
		isGroupAdmin, err := s.groupRepo.IsGroupAdmin(ctx, groupId, userId)
		if err != nil || !isGroupAdmin {
			return err
		}
		err = s.walletRepo.SoftDeleteWallet(ctx, walletId)
		if err != nil {
			return err
		}
	}

	err = s.walletRepo.SoftDeleteWallet(ctx, walletId)
	if err != nil {
		return err
	}
	return nil
}
