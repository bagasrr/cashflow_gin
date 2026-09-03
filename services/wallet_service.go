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

type WalletService interface {
	CreatePersonalWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error)
	CreateGroupWallet(ctx context.Context, wallet models.Wallet) (*models.Wallet, error)
	GetAll(ctx context.Context) (*[]models.Wallet, error)
	GetWalletByID(ctx context.Context, userID, walletID uuid.UUID) (*models.Wallet, error)
	GetMine(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Wallet, int64, error)

	UpdateWalletName(ctx context.Context, userID, walletID uuid.UUID, newName string) (*models.Wallet, error)
	DeleteWallet(ctx context.Context, walletId, userId uuid.UUID) error
	TransferBalance(ctx context.Context, userID, fromWalletID, toWalletID uuid.UUID, amount int64, notes string) error

	GetWalletChartData(ctx context.Context, userID, walletID uuid.UUID, params api.GetWalletChartDataParams) ([]models.WalletChartPoint, error)

}

type walletService struct {
	// Kita butuh WalletRepository untuk akses data wallet
	walletRepo repository.WalletRepository
	groupRepo  repository.GroupRepository
	categoryRepo repository.CategoryRepository

}

func NewWalletService(wRepo repository.WalletRepository, gRepo repository.GroupRepository, cRepo repository.CategoryRepository) WalletService {
	return &walletService{walletRepo: wRepo, groupRepo: gRepo, categoryRepo: cRepo}
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

func (s *walletService) GetMine(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Wallet, int64, error) {
	// 1. Validasi cegah Hacker/Bug Frontend (Opsional karena udah di-cover Helper Handler, tapi bagus untuk Double Security)
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// 2. Langsung Panggil Repository (HAPUS PERHITUNGAN OFFSET DI SINI)
	wallets, totalItems, err := s.walletRepo.FindAllMine(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return wallets, totalItems, nil
}
func (s *walletService) DeleteWallet(ctx context.Context, walletId, userId uuid.UUID) error {
	// 1. FETCH: Tarik wujud asli dompetnya dari database (Tambahkan fungsi GetWalletByID di repo lu kalau belum ada)
	wallet, err := s.walletRepo.GetWalletByID(ctx, walletId)
	if err != nil {
		return errors.New("wallet not found")
	}

	// 2. OTORISASI MUTLAK (Aman dari Nil Pointer)
	if wallet.GroupID != nil {
		// SKENARIO A: Dompet Grup
		// Karena kita udah ngecek != nil, melakukan *wallet.GroupID di bawah ini 100% aman dari Panic
		isAdmin, err := s.groupRepo.IsGroupAdmin(ctx, *wallet.GroupID, userId)
		if err != nil || !isAdmin {
			return errors.New("unauthorized: only group admin can delete this group wallet")
		}
	} else {
		// SKENARIO B: Dompet Personal
		// CEK MUTLAK: Pastikan user yang login adalah pemiliknya
		if wallet.UserID == nil || *wallet.UserID != userId {
			return errors.New("unauthorized: you do not own this personal wallet")
		}
	}

	// 3. EKSEKUSI DELETE TUNGGAL
	// Pastikan di dalam SoftDeleteWallet ini lu juga me-soft delete semua transaksi
	// yang terkait dengan walletId tersebut menggunakan db.Transaction()
	err = s.walletRepo.SoftDeleteWallet(ctx, walletId)
	if err != nil {
		return err
	}

	return nil
}

// Tanda tangan wajib menerima userID pembawa request
func (s *walletService) UpdateWalletName(ctx context.Context, userID, walletID uuid.UUID, newName string) (*models.Wallet, error) {
	// 1. FETCH: Tarik dompet aslinya
	wallet, err := s.walletRepo.GetWalletByID(ctx, walletID)
	if err != nil {
		return nil, errors.New("wallet not found")
	}

	// 2. OTORISASI MUTLAK
	if wallet.GroupID != nil {
		// SKENARIO A: Dompet Grup
		// Panggil fungsi IsGroupAdmin dari GroupRepo lu
		isAdmin, err := s.groupRepo.IsGroupAdmin(ctx, *wallet.GroupID, userID)
		if err != nil || !isAdmin {
			return nil, errors.New("forbidden: only group admin can update this wallet")
		}
	} else {
		// SKENARIO B: Dompet Personal
		if wallet.UserID == nil || *wallet.UserID != userID {
			return nil, errors.New("forbidden: you do not own this personal wallet")
		}
	}

	// 3. MODIFY: Ubah namanya
	wallet.Name = newName

	// 4. SAVE: Lempar ke Repo
	return s.walletRepo.UpdateWallet(ctx, wallet)
}

func (s *walletService) GetWalletChartData(ctx context.Context, userID, walletID uuid.UUID, params api.GetWalletChartDataParams) ([]models.WalletChartPoint, error) {


	// 1. DEFINISIKAN WAKTU DEFAULT FALLBACK (30 Hari Terakhir)
	now := time.Now()

	// Set Jam ke 23:59:59 untuk batas akhir hari ini
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	// Tarik mundur 30 hari, set jam ke 00:00:00 untuk awal hari
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -30)

	// 2. TIMPA DENGAN INPUT FRONTEND JIKA ADA (Menggunakan objek .Time dari openapi_types.Date)
	if params.EndDate != nil {
		parsedEnd := params.EndDate.Time
		endDate = time.Date(parsedEnd.Year(), parsedEnd.Month(), parsedEnd.Day(), 23, 59, 59, 0, parsedEnd.Location())
	}

	if params.StartDate != nil {
		parsedStart := params.StartDate.Time
		startDate = time.Date(parsedStart.Year(), parsedStart.Month(), parsedStart.Day(), 0, 0, 0, 0, parsedStart.Location())
	}

	// 3. JARING PENGAMAN (DEFENSIVE CODING) - BLIND SPOT VALIDASI WAKTU
	// Jangan pernah percaya 100% sama Frontend. Ada kalanya mereka salah kirim rentang waktu.
	// Jika start_date lebih maju daripada end_date (misal: start=Desember, end=Januari), query DB akan hancur atau kosong.
	if startDate.After(endDate) {
		return nil, errors.New("bad request: tanggal mulai (start_date) tidak boleh melewati tanggal akhir (end_date)")
	}

	// Maksimal rentang waktu (Opsional, tapi krusial untuk kestabilan server)
	// Batasi rentang pencarian grafik maksimal 1 tahun (365 hari) agar user gak iseng narik data grafik 10 tahun sekaligus yang bikin DB hang.
	if endDate.Sub(startDate).Hours() > 24*365 {
		return nil, errors.New("bad request: rentang waktu grafik maksimal adalah 1 tahun")
	}

	// 4. LEMPAR KE REPOSITORY
	// Menggunakan fungsi Repo yang udah kita bedah di obrolan sebelumnya
	chartPoints, err := s.walletRepo.GetWalletChartData(ctx, userID, walletID, startDate, endDate)

	if err != nil {
		return nil, err
	}

	return chartPoints, nil
}

func (s *walletService) getOrCreateTransferCategory(ctx context.Context, userID uuid.UUID, name string, catType models.CategoryType) (uuid.UUID, error) {
	category, err := s.categoryRepo.FindTransferCategory(ctx, name, catType, userID)
	if err == nil {
		return category.ID, nil
	}

	newCategory := models.Category{
		UserID: &userID,
		Name:   name,
		Type:   catType,
	}
	
	created, err := s.categoryRepo.Create(ctx, &newCategory)
	if err != nil {
		return uuid.Nil, err
	}
	
	return created.ID, nil
}

func (s *walletService) TransferBalance(ctx context.Context, userID, fromWalletID, toWalletID uuid.UUID, amount int64, notes string) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if fromWalletID == toWalletID {
		return errors.New("cannot transfer to the same wallet")
	}

	fromWallet, err := s.walletRepo.FindByID(ctx, fromWalletID)
	if err != nil {
		return errors.New("from_wallet not found")
	}

	if fromWallet.GroupID != nil {
		isGroupMember, err := s.groupRepo.IsGroupMember(ctx, *fromWallet.GroupID, userID)
		if err != nil || !isGroupMember {
			return errors.New("unauthorized: from_wallet")
		}
	} else if *fromWallet.UserID != userID {
		return errors.New("unauthorized: from_wallet")
	}

	if fromWallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	toWallet, err := s.walletRepo.FindByID(ctx, toWalletID)
	if err != nil {
		return errors.New("to_wallet not found")
	}

	if toWallet.GroupID != nil {
		isGroupMember, err := s.groupRepo.IsGroupMember(ctx, *toWallet.GroupID, userID)
		if err != nil || !isGroupMember {
			return errors.New("unauthorized: to_wallet")
		}
	} else if *toWallet.UserID != userID {
		return errors.New("unauthorized: to_wallet")
	}

	fromCatID, err := s.getOrCreateTransferCategory(ctx, userID, "Transfer Out", models.CategoryExpense)
	if err != nil {
		return err
	}
	
	toCatID, err := s.getOrCreateTransferCategory(ctx, userID, "Transfer In", models.CategoryIncome)
	if err != nil {
		return err
	}

	expenseTrx := &models.Transaction{
		UserID:      userID,
		WalletID:    fromWalletID,
		CategoryID:  fromCatID,
		Title:       "Transfer Out",
		Amount:      amount,
		Description: notes,
		Date:        time.Now(),
	}
	
	incomeTrx := &models.Transaction{
		UserID:      userID,
		WalletID:    toWalletID,
		CategoryID:  toCatID,
		Title:       "Transfer In",
		Amount:      amount,
		Description: notes,
		Date:        time.Now(),
	}

	return s.walletRepo.TransferBalance(ctx, expenseTrx, incomeTrx)
}
