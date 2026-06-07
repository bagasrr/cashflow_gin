package repository

import (
	"cashflow_gin/models"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRepository interface {
	CreateGroupWithWalletAndMembers(ctx context.Context, group *models.Group) (*models.Group, error)
	GetAllGroups(ctx context.Context, limit, offset int) (*[]models.Group, int64, error)
	GetMyGroups(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Group, int64, error)

	IsGroupWallet(ctx context.Context, walletID uuid.UUID) (bool, uuid.UUID, error)
	IsGroupMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error)
	IsGroupAdmin(ctx context.Context, groupID, userID uuid.UUID) (bool, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
	UpdateGroup(ctx context.Context, group *models.Group) (*models.Group, error)
	DeleteGroup(ctx context.Context, groupID uuid.UUID) error

	CreateMembers(ctx context.Context, members []models.GroupMember) error
	RemoveUserFromGroup(ctx context.Context, groupID, userID uuid.UUID) error
}

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) CreateGroupWithWalletAndMembers(ctx context.Context, group *models.Group) (*models.Group, error) {
	// 1. Eksekusi Create (Insert ke DB)
	// Ingat, 'group' adalah pointer. Setelah ini sukses, group.ID akan terisi.
	if err := r.db.WithContext(ctx).Create(group).Error; err != nil {
		return nil, err
	}

	// 2. RELOAD DATA UTUH (The Enterprise Way)
	// Tarik ulang dari database berdasarkan ID yang baru saja terbentuk.
	// Ini menjamin semua relasi dan nilai default DB terbaca sempurna.
	var createdGroup models.Group
	err := r.db.WithContext(ctx).
		Preload("Wallet").
		Preload("Members").      // Pastikan nama relasi sesuai dengan struct GORM lu
		Preload("Members.User"). // Tarik juga data user kalau butuh username di response
		First(&createdGroup, "id = ?", group.ID).Error

	if err != nil {
		return nil, err
	}

	// 3. Kembalikan data yang sudah matang
	return &createdGroup, nil
}

func (r *groupRepository) GetAllGroups(ctx context.Context, limit, offset int) (*[]models.Group, int64, error) {
	var groups []models.Group
	var totalData int64

	// 1. Bangun fondasi query
	query := r.db.WithContext(ctx).Model(&models.Group{})

	// 2. Eksekusi Count
	if err := query.Count(&totalData).Error; err != nil {
		return nil, 0, err
	}

	// 3. Eksekusi pencarian dengan Select, Preload (perbaiki typo), dan Limit
	err := query.
		Select(`groups.*, (SELECT COUNT(*) FROM group_members WHERE group_members.group_id = groups.id) AS member_count`).
		Preload("Wallet").
		Preload("Members"). // UBAH KE "Members"
		Limit(limit).Offset(offset).
		Find(&groups).Error

	return &groups, totalData, err
}

// Pastikan lu me-return int64 untuk totalItems demi paginasi Handler lu
func (r *groupRepository) GetMyGroups(ctx context.Context, userID uuid.UUID, limit, offset int) (*[]models.Group, int64, error) {
	var groups []models.Group
	var totalItems int64

	// 1. BANGUN JEMBATAN QUERY (INNER JOIN)
	// Kita menyuruh Postgres menggabungkan tabel groups dan group_members
	// lalu memfilternya berdasarkan user_id yang ada di group_members.
	query := r.db.WithContext(ctx).
		Model(&models.Group{}).
		Select(`groups.*, (SELECT COUNT(*) FROM group_members WHERE group_members.group_id = groups.id) AS member_count`).
		Joins("JOIN group_members ON group_members.group_id = groups.id").
		Where("group_members.user_id = ?", userID)

	// 2. HITUNG TOTAL DATA (Sebelum kena Limit/Offset)
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// 3. TARIK DATA HALAMAN INI + PRELOAD
	// Preload "Wallet" lu panggil KALAU API List Groups lu butuh nampilin dompet.
	// Jangan Preload "Members" di daftar List Groups, itu bakal bikin query berat
	// (N+1 problem). Detail member cukup dipanggil di API GetGroupByID.
	if err := query.Limit(limit).Offset(offset).
		Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	// 4. KEMBALIKAN ALAMAT MEMORI DAN TOTAL DATA
	return &groups, totalItems, nil
}

func (r *groupRepository) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.WithContext(ctx).Preload("Wallet").Preload("Members.User").First(&group, "id = ?", groupID).Error
	return &group, err
}

func (r *groupRepository) UpdateGroup(ctx context.Context, group *models.Group) (*models.Group, error) {
	// 1. SAVE: Timpa data ke database.
	// GORM otomatis ngebaca ID dari 'group' dan ngelakuin UPDATE query.
	if err := r.db.WithContext(ctx).Save(group).Error; err != nil {
		return nil, err
	}

	// 2. RELOAD: Tarik ulang beserta anak-anaknya (Wallet & Members)
	var updatedGroup models.Group
	err := r.db.WithContext(ctx).
		Preload("Wallet").
		Preload("Members").
		Preload("Members.User").
		First(&updatedGroup, "id = ?", group.ID).Error

	return &updatedGroup, err
}

func (r *groupRepository) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Group{}, "id = ?", groupID).Error
}

// repository/group.go
func (r *groupRepository) CreateMembers(ctx context.Context, members []models.GroupMember) error {
	// Langsung gas simpan.
	// Gak perlu cek GroupID ada atau gak, karena Foreign Key Database bakal nolak otomatis kalau gak ada.
	return r.db.WithContext(ctx).Create(&members).Error
}

func (r *groupRepository) RemoveUserFromGroup(ctx context.Context, groupID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMember{}).Error
}

func (r *groupRepository) IsGroupWallet(ctx context.Context, walletID uuid.UUID) (bool, uuid.UUID, error) {
	var wallet models.Wallet

	// 1. SELECT & FETCH: Suruh Postgres mengambil nilai "group_id" dari baris pertama yang cocok
	err := r.db.WithContext(ctx).
		Model(models.Wallet{}).
		Select("group_id").
		Where("id = ? AND group_id IS NOT NULL", walletID).
		First(&wallet).Error

	if err != nil {
		// 2. GERBANG VALIDASI: Kalau error karena tidak ketemu, berarti bukan dompet grup
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, uuid.Nil, nil
		}
		// 3. ERROR FATAL: Kalau error karena putus koneksi/sistem
		return false, uuid.Nil, err
	}

	// 4. KEMBALIKAN DATA
	// CATATAN MUTLAK: Kalau di models.Wallet lu field GroupID itu bertipe pointer (*uuid.UUID),
	// lu wajib nge-dereference-nya dengan menulis *wallet.GroupID di bawah ini.
	// Tapi kalau tipe datanya murni (uuid.UUID), cukup tulis seperti ini:
	return true, *wallet.GroupID, nil
}

func (r *groupRepository) IsGroupMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count).Error
	fmt.Println("Count : ", count)
	return count > 0, err
}

func (r *groupRepository) IsGroupAdmin(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ? AND members_role = ?", groupID, userID, models.GroupAdmin).Count(&count).Error
	return count > 0, err
}
