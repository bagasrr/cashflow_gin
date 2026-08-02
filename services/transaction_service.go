package services

import (
	"cashflow_gin/api"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"errors"
	"time"

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
	GetTransactionsByWallet(ctx context.Context, userID, walletID uuid.UUID, params api.GetWalletTransactionsParams, page, limit, offset int) ([]models.Transaction, int64, error)
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
	// 1. VALIDASI WALLET MUTLAK
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
			return nil, errors.New("unauthorized: user is not a member of the group wallet")
		}
	} else {
		reqUser := s.transactionRepo.IsOwner(ctx, userID, input.WalletId)
		if !reqUser {
			return nil, errors.New("unauthorized: wallet does not belong to user")
		}
	}

	// 2. VALIDASI CATEGORY MUTLAK
	catId, err := uuid.Parse(input.CategoryId)
	if err != nil {
		return nil, errors.New("invalid category id")
	}

	category, err := s.categoryRepo.FindByID(ctx, catId)
	if err != nil {
		return nil, errors.New("category not found")
	}

	// 3. LOGIKA AKUNTANSI MUTLAK (Pemisahan Fakta dan Dampak)
	absoluteAmount := absInt64(input.Amount)

	impactAmount := absoluteAmount
	if category.Type == "EXPENSE" || category.Type == "INVESTMENT" {
		impactAmount = -absoluteAmount // Dampak minus ke dompet
	}

	// 4. PENGAMAN POINTER
	var safeDescription string
	if input.Description != nil {
		safeDescription = *input.Description
	}

	// Struct mentah sebelum masuk database
	transaction := models.Transaction{
		UserID:      userID,
		WalletID:    walletUUID,
		CategoryID:  category.ID,
		Title:       input.Title,
		Amount:      absoluteAmount, // Fakta nominal positif
		Description: safeDescription,
		Date:        input.Date,
	}

	// 5. EKSEKUSI DAN TANGKAP HASIL (Integrasi dengan Repo baru lu)
	// newTrx di sini adalah data yang sudah di-reload (memiliki data Category utuh)
	newTrx, err := s.transactionRepo.CreateWithWalletUpdate(ctx, &transaction, impactAmount)
	if err != nil {
		return nil, err
	}

	// Kembalikan newTrx, bukan &transaction mentah
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

	// 1. TENTUKAN DAMPAK TRANSAKSI LAMA
	oldImpact := transaction.Amount
	if transaction.Category.Type == "EXPENSE" || transaction.Category.Type == "INVESTMENT" {
		oldImpact = -transaction.Amount
	}

	// 2. TERAPKAN PERUBAHAN DATA TEKS & TANGGAL
	if input.Title != "" {
		transaction.Title = input.Title
	}
	if input.Description != "" {
		transaction.Description = input.Description
	}
	if !input.Date.IsZero() {
		transaction.Date = input.Date
	}

	// 3. TERAPKAN PERUBAHAN KATEGORI (MUTLAK)
	if input.CategoryId != "" && input.CategoryId != transaction.CategoryID.String() {
		newCatID, err := uuid.Parse(input.CategoryId)
		if err != nil {
			return nil, errors.New("bad request: format category_id tidak valid")
		}

		// Panggil Repository Kategori lu. (Pastikan lu punya fungsi FindByID ini di categoryRepo lu)
		// Kalau balikan repo lu itu struct utuh (bukan pointer), hilangkan bintang/ampersand sesuai kebutuhan.
		newCategory, err := s.categoryRepo.FindByID(ctx, newCatID)
		if err != nil {
			return nil, errors.New("bad request: kategori baru tidak ditemukan")
		}

		// Update memori sementara agar perhitungan delta di bawah tidak cacat
		transaction.CategoryID = newCategory.ID
		transaction.Category = *newCategory
	}

	// 4. TERAPKAN PERUBAHAN NOMINAL
	if input.Amount != 0 {
		transaction.Amount = absInt64(input.Amount)
	}

	// 5. TENTUKAN DAMPAK TRANSAKSI BARU
	// Karena transaction.Category udah di-update di atas (jika ada perubahan),
	// Type yang dibaca di sini adalah Type yang 100% akurat.
	newImpact := transaction.Amount
	if transaction.Category.Type == "EXPENSE" || transaction.Category.Type == "INVESTMENT" {
		newImpact = -transaction.Amount
	}

	// 6. RUMUS DELTA MUTLAK
	deltaAmount := newImpact - oldImpact

	// 7. EKSEKUSI KE REPO GABUNGAN
	err = s.transactionRepo.UpdateTransaction(ctx, transaction, deltaAmount)
	if err != nil {
		return nil, err
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

func (s *transactionService) GetTransactionsByWallet(ctx context.Context, userID, walletID uuid.UUID, params api.GetWalletTransactionsParams, page, limit, offset int) ([]models.Transaction, int64, error) {

	// 1. Logika Jatuh Bebas (Fallback) Tanggal
	now := time.Now()
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0) // Default 1 bulan terakhir

	if params.EndDate != nil {
		parsedEnd := params.EndDate.Time
		endDate = time.Date(parsedEnd.Year(), parsedEnd.Month(), parsedEnd.Day(), 23, 59, 59, 0, parsedEnd.Location())
	}

	if params.StartDate != nil {
		parsedStart := params.StartDate.Time
		startDate = time.Date(parsedStart.Year(), parsedStart.Month(), parsedStart.Day(), 0, 0, 0, 0, parsedStart.Location())
	}

	if startDate.After(endDate) {
		return nil, 0, errors.New("bad request: start_date tidak boleh melebihi end_date")
	}

	// 3. Lempar ke DB
	var searchStr string
	if params.Search != nil {
		searchStr = *params.Search
	}

	validOrder := "desc"

	if params.SortOrder != nil {
		orderStr := string(*params.SortOrder)

		if orderStr == "asc" || orderStr == "desc" {
			validOrder = orderStr
		}
	}

	// Validasi Sort By (WHITELIST MUTLAK MENCEGAH SQL INJECTION)
	validSortBy := "date" // Default
	if params.SortBy != nil {
		switch *params.SortBy {
		case "amount", "title", "date", "created_at": // Hanya izinkan kolom ini
			validSortBy = *params.SortBy
		default:
			// Tendang user kalau dia nyoba masukin nama kolom aneh
			return nil, 0, errors.New("bad request: parameter sort_by tidak valid")
		}
	}

	return s.transactionRepo.GetTransactionsByWallet(ctx, userID, walletID, startDate, endDate, limit, offset, searchStr, validSortBy, validOrder)
}
