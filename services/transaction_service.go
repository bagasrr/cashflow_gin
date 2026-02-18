package services

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/dto/response"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
)

type TransactionService interface {
	Create(userID uuid.UUID, input request.CreateTransactionRequest) (response.TransactionResponse, error)
	GetAll() (*[]response.TransactionResponse, error)
	GetTransactionByID(userID uuid.UUID, transactionID uuid.UUID) (response.TransactionResponse, error)
	UpdateTransaction(userID, transactionID uuid.UUID, input request.UpdateTransactionRequest) (response.TransactionResponse, error)
	SoftDeleteTransaction(userID uuid.UUID, transactionID uuid.UUID, walletID uuid.UUID) error
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	categoryRepo    repository.CategoryRepository
	userRepo        repository.UserRepository
	groupRepo       repository.GroupRepository
	walletRepo      repository.WalletRepository
}

// Constructor minta 2 Repository sekarang
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

func (s *transactionService) Create(userID uuid.UUID, input request.CreateTransactionRequest) (response.TransactionResponse, error) {
	// 1. Parsing UUID
	walletUUID, err := uuid.Parse(input.WalletID)
	if err != nil {
		return response.TransactionResponse{}, errors.New("invalid wallet id")
	}

	wallet, err := s.walletRepo.FindByID(walletUUID)
	if err != nil {
		return response.TransactionResponse{}, errors.New("wallet not found")
	}

	// Personal Wallet
	// Cek Apakah user id yang mengirim = user id yang punya wallet
	isGroupWallet, err := s.groupRepo.IsGroupWallet(walletUUID)
	fmt.Println("Is Group Wallet?", isGroupWallet)
	if err != nil {
		return response.TransactionResponse{}, errors.New("failed to check wallet type")
	}
	if isGroupWallet {
		isGroupMember, err := s.groupRepo.IsGroupMember(*wallet.GroupID, userID)
		fmt.Println("Is Group Member?", isGroupMember)
		if err != nil {
			return response.TransactionResponse{}, errors.New("failed to check group membership")
		}
		if !isGroupMember {
			return response.TransactionResponse{}, errors.New("unauthorized: user is not a member of the group wallet, cannot create personal transaction")
		}
	} else {
		reqUser := s.transactionRepo.IsOwner(userID, input.WalletID)
		if !reqUser {
			return response.TransactionResponse{}, errors.New("unauthorized: wallet does not belong to user")
		}
	}

	// 2. BUSSINESS LOGIC: Cek Category Type (Income/Expense)
	category, err := s.categoryRepo.FindByName(input.CategoryName)
	if err != nil {
		return response.TransactionResponse{}, errors.New("category not found")
	}

	finalAmount := input.Amount

	// Logic Matematika:
	// Jika Category Type == EXPENSE, saldo harus berkurang (negatif)
	// Kita pakai Math.Abs buat mastiin input selalu positif dulu, baru dikali -1
	if category.Type == "EXPENSE" {
		finalAmount = -math.Abs(finalAmount)
	} else {
		// Jika INCOME, pastiin positif
		finalAmount = math.Abs(finalAmount)
	}

	// 3. Construct Object
	transaction := models.Transaction{
		UserID:      userID,
		WalletID:    walletUUID,
		CategoryID:  category.ID,
		Title:       input.Title,
		Amount:      finalAmount, // Nilai sudah otomatis +/- sesuai kategori
		Description: input.Description,
		Date:        input.Date,
	}

	// 4. Save Atomic (Transaction + Wallet Update)
	err = s.transactionRepo.CreateWithWalletUpdate(&transaction)
	if err != nil {
		return response.TransactionResponse{}, err
	}

	res := response.TransactionResponse{
		ID:          transaction.ID.String(),
		Title:       transaction.Title,
		Amount:      transaction.Amount,
		Description: transaction.Description,
		Date:        transaction.Date,
		Category: response.CategoryResponse{
			Name: category.Name,
			Type: category.Type,
		},
	}
	return res, nil
}

func (s *transactionService) GetAll() (*[]response.TransactionResponse, error) {
	// 1. Panggil Repository (Filter by UserID biar gak bocor data orang lain)
	transactions, err := s.transactionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 2. Mapping dari []models.Transaction ke []response.TransactionResponse
	// Kita inisialisasi slice kosong biar kalau data 0, returnnya "[]" bukan "null"
	transactionResponses := []response.TransactionResponse{}

	for _, t := range transactions {
		res := response.TransactionResponse{
			ID:          t.ID.String(),
			Title:       t.Title,
			Amount:      t.Amount,
			Description: t.Description,
			Date:        t.Date,
			Category: response.CategoryResponse{
				Name: t.Category.Name,
				Type: t.Category.Type,
			},
		}
		transactionResponses = append(transactionResponses, res)
	}

	return &transactionResponses, nil
}

func (s *transactionService) GetTransactionByID(userID uuid.UUID, transactionID uuid.UUID) (response.TransactionResponse, error) {
	user, err := s.userRepo.FindMyProfile(userID)
	if err != nil {
		return response.TransactionResponse{}, err
	}

	if user.ID != userID {
		return response.TransactionResponse{}, errors.New("unauthorized: user not found")
	}

	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return response.TransactionResponse{}, err
	}

	if transaction.UserID != user.ID {
		return response.TransactionResponse{}, errors.New("unauthorized: transaction does not belong to user")
	}
	res := response.TransactionResponse{
		ID:          transaction.ID.String(),
		Title:       transaction.Title,
		Amount:      transaction.Amount,
		Description: transaction.Description,
		Date:        transaction.Date,
		Category: response.CategoryResponse{
			Name: transaction.Category.Name,
			Type: transaction.Category.Type,
		},
	}
	return res, nil
}

func (s *transactionService) UpdateTransaction(userID, transactionID uuid.UUID, input request.UpdateTransactionRequest) (response.TransactionResponse, error) {
	reqUser, err := s.userRepo.FindMyProfile(userID)
	if err != nil {
		return response.TransactionResponse{}, errors.New("unauthorized: user not found")
	}
	// Cari Transaksi dan validasi
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return response.TransactionResponse{}, errors.New("transaction not found")
	}
	if transaction.UserID != reqUser.ID {
		return response.TransactionResponse{}, errors.New("unauthorized: transaction does not belong to user")
	}

	oldAmount := transaction.Amount

	// Update fields
	if input.Title != "" {
		transaction.Title = input.Title
	}
	if input.Description != "" {
		transaction.Description = input.Description
	}
	if !input.Date.IsZero() {
		transaction.Date = input.Date
	}

	var deltaAmount float64
	if input.Amount != 0 {
		var newAmount float64

		if transaction.Category.Type == "EXPENSE" {
			newAmount = -math.Abs(input.Amount)
		} else {
			newAmount = math.Abs(input.Amount)
		}

		transaction.Amount = newAmount
		deltaAmount = newAmount - oldAmount

	}

	if deltaAmount != 0 {
		fmt.Println("Update Transaction With Wallet Ballance")
		err := s.transactionRepo.UpdateTransactionWithWalletBallance(transaction, deltaAmount)
		if err != nil {
			return response.TransactionResponse{}, err
		}
	} else {
		fmt.Println("Update Transaction ONLY")

		// Simpan perubahan
		err = s.transactionRepo.UpdateTransaction(transaction)
		if err != nil {
			return response.TransactionResponse{}, err
		}
	}

	res := response.TransactionResponse{
		ID:          transaction.ID.String(),
		Title:       transaction.Title,
		Amount:      transaction.Amount,
		Description: transaction.Description,
		Date:        transaction.Date,
		Category: response.CategoryResponse{
			Name: transaction.Category.Name,
			Type: transaction.Category.Type,
		},
	}

	return res, nil
}

func (s *transactionService) SoftDeleteTransaction(userID, transactionID, walletID uuid.UUID) error {
	reqUser, err := s.userRepo.FindMyProfile(userID)
	if err != nil {
		return errors.New("unauthorized: user not found")
	}
	// Cari Transaksi dan validasi
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return errors.New("transaction not found")
	}
	if transaction.UserID != reqUser.ID {
		return errors.New("unauthorized: transaction does not belong to user")
	}

	// Logic Matematika:
	// Untuk Soft Delete, kita harus ngurangin balance wallet dengan amount transaksi yang mau dihapus
	deltaAmount := -transaction.Amount
	fmt.Println("Delta Amount", deltaAmount)

	err = s.transactionRepo.SoftDeleteTransaction(transactionID, deltaAmount, walletID)
	if err != nil {
		return err
	}

	return nil
}
