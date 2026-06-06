package services

import (
	"cashflow_gin/api"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Helper internal untuk menggantikan math.Abs yang cuma buat float64
func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

type TransactionService interface {
	Create(ctx context.Context, userID uuid.UUID, input api.CreateTransactionReq) (*models.Transaction, error)
	GetAll(ctx context.Context, userID uuid.UUID, role models.UserRole) (*[]models.Transaction, error)
	GetTransactionByID(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, userID, transactionID uuid.UUID, input api.UpdateTransactionReq) (*models.Transaction, error)
	SoftDeleteTransaction(ctx context.Context, userID, transactionID uuid.UUID) error
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	categoryRepo    repository.CategoryRepository
	userRepo        repository.UserRepository
	groupRepo       repository.GroupRepository
	walletRepo      repository.WalletRepository
}

func NewTransactionService(
	tRepo repository.TransactionRepository,
	cRepo repository.CategoryRepository,
	uRepo repository.UserRepository,
	gRepo repository.GroupRepository,
	wRepo repository.WalletRepository,
) TransactionService {
	return &transactionService{
		transactionRepo: tRepo,
		categoryRepo:    cRepo,
		userRepo:        uRepo,
		groupRepo:       gRepo,
		walletRepo:      wRepo,
	}
}

func (s *transactionService) Create(ctx context.Context, userID uuid.UUID, input api.CreateTransactionReq) (*models.Transaction, error) {
	// 1. TYPO FIX: Gunakan WalletId, bukan WalletID
	walletUUID, err := uuid.Parse(input.WalletId)
	if err != nil {
		return nil, errors.New("invalid wallet id")
	}

	wallet, err := s.walletRepo.FindByID(ctx, walletUUID)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	isGroupWallet, _, err := s.groupRepo.IsGroupWallet(ctx, walletUUID)
	if err != nil {
		return nil, errors.New("failed to check wallet type")
	}

	if isGroupWallet {
		isGroupMember, err := s.groupRepo.IsGroupMember(ctx, *wallet.GroupID, userID)
		if err != nil {
			return nil, errors.New("failed to check group membership")
		}
		if !isGroupMember {
			return nil, errors.New("unauthorized: user is not a member of the group wallet, cannot create personal transaction")
		}
	} else {
		// TYPO FIX: WalletId
		reqUser := s.transactionRepo.IsOwner(ctx, userID, input.WalletId)
		if !reqUser {
			return nil, errors.New("unauthorized: wallet does not belong to user")
		}
	}

	// 2. TYPO FIX: Gunakan CategoryId, bukan CategoryID
	catId, err := uuid.Parse(input.CategoryId)
	if err != nil {
		return nil, errors.New("invalid category id")
	}

	category, err := s.categoryRepo.FindByID(ctx, catId)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// LOGIKA INT64 MURNI
	finalAmount := input.Amount
	if category.Type == "EXPENSE" {
		finalAmount = -absInt64(finalAmount)
	} else {
		finalAmount = absInt64(finalAmount)
	}

	// 3. PENGAMAN POINTER MUTLAK UNTUK DESCRIPTION
	var safeDescription string
	if input.Description != nil {
		safeDescription = *input.Description
	}

	transaction := models.Transaction{
		UserID:      userID,
		WalletID:    walletUUID,
		CategoryID:  category.ID,
		Title:       input.Title,
		Amount:      finalAmount,
		Description: safeDescription, // Masukkan variabel yang sudah dipastikan aman (string)
		Date:        input.Date,
	}

	newTrx, err := s.transactionRepo.CreateWithWalletUpdate(ctx, &transaction)
	if err != nil {
		return nil, err
	}

	return newTrx, nil
}

func (s *transactionService) GetAll(ctx context.Context, userID uuid.UUID, role models.UserRole) (*[]models.Transaction, error) {
	var transactions *[]models.Transaction
	var err error

	if role == models.RoleAdmin || role == models.RoleModerator {
		transactions, err = s.transactionRepo.FindAll(ctx)
	} else {
		transactions, err = s.transactionRepo.FindAllByUserID(ctx, userID)
	}

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *transactionService) GetTransactionByID(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*models.Transaction, error) {
	user, err := s.userRepo.FindMyProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.ID != userID {
		return nil, errors.New("unauthorized: user not found")
	}

	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	if transaction.UserID != user.ID {
		return nil, errors.New("unauthorized: transaction does not belong to user")
	}

	return transaction, nil
}

func (s *transactionService) UpdateTransaction(ctx context.Context, userID, transactionID uuid.UUID, input api.UpdateTransactionReq) (*models.Transaction, error) {
	reqUser, err := s.userRepo.FindMyProfile(ctx, userID)
	if err != nil {
		return nil, errors.New("unauthorized: user not found")
	}

	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	if transaction.UserID != reqUser.ID {
		return nil, errors.New("unauthorized: transaction does not belong to user")
	}

	oldAmount := transaction.Amount

	if input.Title != "" {
		transaction.Title = input.Title
	}
	if input.Description != "" {
		transaction.Description = input.Description
	}
	if !input.Date.IsZero() {
		transaction.Date = input.Date
	}

	// LOGIKA INT64 MURNI
	var deltaAmount int64
	if input.Amount != 0 {
		var newAmount int64

		if transaction.Category.Type == "EXPENSE" {
			newAmount = -absInt64(input.Amount)
		} else {
			newAmount = absInt64(input.Amount)
		}

		transaction.Amount = newAmount
		deltaAmount = newAmount - oldAmount
	}

	if deltaAmount != 0 {
		fmt.Println("Update Transaction With Wallet Balance")
		// CATATAN: Pastikan UpdateTransactionWithWalletBallance di Repo nerima int64
		err := s.transactionRepo.UpdateTransactionWithWalletBallance(ctx, transaction, deltaAmount)
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Update Transaction ONLY")
		err = s.transactionRepo.UpdateTransaction(ctx, transaction)
		if err != nil {
			return nil, err
		}
	}

	return transaction, nil
}
func (s *transactionService) SoftDeleteTransaction(ctx context.Context, userID, transactionID uuid.UUID) error {
	reqUser, err := s.userRepo.FindMyProfile(ctx, userID)
	if err != nil {
		return errors.New("unauthorized: user not found")
	}

	// a. CARI TRANSAKSI DI DATABASE
	transaction, err := s.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}
	if transaction.UserID != reqUser.ID {
		return errors.New("unauthorized: transaction does not belong to user")
	}

	deltaAmount := -transaction.Amount

	// b. LEMPAR KE REPOSITORY
	// Lihat baris ini! Kita mengambil transaction.WalletID dari objek yang baru aja kita cari.
	// Kita ngirim 4 parameter ke Repo sesuai permintaan Repo lu.
	err = s.transactionRepo.SoftDeleteTransaction(ctx, transactionID, deltaAmount, transaction.WalletID)
	if err != nil {
		return err
	}

	return nil
}
