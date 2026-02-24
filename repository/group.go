package repository

import (
	"cashflow_gin/models"
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRepository interface {
	CreateGroupWithWalletAndMembers(ctx context.Context, group *models.Group, wallet *models.Wallet, members *[]models.GroupMember) error
	GetAllGroups(ctx context.Context) (*[]models.Group, error)

	IsGroupWallet(ctx context.Context, walletID uuid.UUID) (bool, error)
	IsGroupMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
	UpdateGroup(ctx context.Context, group *models.Group) error
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

func (r *groupRepository) CreateGroupWithWalletAndMembers(ctx context.Context, group *models.Group, wallet *models.Wallet, members *[]models.GroupMember) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// A. Create Group dulu (biar dapet ID Group)
		if err := tx.Create(group).Error; err != nil {
			return err
		}

		// B. Assign GroupID ke Wallet & Create Wallet
		wallet.GroupID = &group.ID // Asumsi di model Wallet pake pointer *uuid.UUID
		if err := tx.Create(wallet).Error; err != nil {
			return err
		}

		// C. Assign GroupID ke semua Member & Create Members
		for i := range *members {
			(*members)[i].GroupID = group.ID
		}

		// Batch Insert Members
		if err := tx.Create(members).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *groupRepository) GetAllGroups(ctx context.Context) (*[]models.Group, error) {
	var groups []models.Group

	err := r.db.WithContext(ctx).
		Table("groups").
		Select(`
			groups.*,(
				SELECT COUNT(*)
				FROM group_members
				WHERE group_members.group_id = groups.id
			) AS member_count
		`).Preload("Wallet").Find(&groups).Error

	return &groups, err
}

func (r *groupRepository) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.WithContext(ctx).Preload("Wallet").Preload("Members.User").First(&group, "id = ?", groupID).Error
	return &group, err
}

func (r *groupRepository) UpdateGroup(ctx context.Context, group *models.Group) error {
	return r.db.WithContext(ctx).Save(group).Error
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

func (r *groupRepository) IsGroupWallet(ctx context.Context, walletID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Wallet{}).Where("id = ? AND group_id IS NOT NULL", walletID).Count(&count).Error
	return count > 0, err
}

func (r *groupRepository) IsGroupMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&count).Error
	fmt.Println("Count : ", count)
	return count > 0, err
}
