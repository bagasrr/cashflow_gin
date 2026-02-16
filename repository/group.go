package repository

import (
	"cashflow_gin/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRepository interface {
	CreateGroupWithWalletAndMembers(group *models.Group, wallet *models.Wallet, members *[]models.GroupMember) error

	GetGroupByID(groupID uuid.UUID) (*models.Group, error)
	UpdateGroup(group *models.Group) error
	DeleteGroup(groupID uuid.UUID) error

	AddUserToGroup(groupID, userID uuid.UUID) error
	RemoveUserFromGroup(groupID, userID uuid.UUID) error
}

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) CreateGroupWithWalletAndMembers(group *models.Group, wallet *models.Wallet, members *[]models.GroupMember) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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

func (r *groupRepository) GetGroupByID(groupID uuid.UUID) (*models.Group, error) {
	var group models.Group
	err := r.db.Preload("Members").First(&group, "id = ?", groupID).Error
	return &group, err
}

func (r *groupRepository) UpdateGroup(group *models.Group) error {
	return r.db.Save(group).Error
}

func (r *groupRepository) DeleteGroup(groupID uuid.UUID) error {
	return r.db.Delete(&models.Group{}, "id = ?", groupID).Error
}

func (r *groupRepository) AddUserToGroup(groupID, userID uuid.UUID) error {
	var group models.Group
	if err := r.db.First(&group, "id = ?", groupID).Error; err != nil {
		return err
	}
	group.Members = append(group.Members, models.GroupMember{
		GroupID: groupID,
		UserID:  userID,
	})
	return r.db.Save(&group).Error
}

func (r *groupRepository) RemoveUserFromGroup(groupID, userID uuid.UUID) error {
	return r.db.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMember{}).Error
}
