package services

import (
	"cashflow_gin/api"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"errors"

	"github.com/google/uuid"
)

type GroupService interface {
	CreateGroup(ctx context.Context, ownerID uuid.UUID, input *api.CreateGroupReq) (*models.Group, error)
	GetAllGroups(ctx context.Context, limit, offset int) (*[]models.Group, int64, error)

	GetMyGroups(ctx context.Context, limit, offset int, userId uuid.UUID) (*[]models.Group, int64, error)
	GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
	UpdateGroup(ctx context.Context, groupID uuid.UUID, name string) (*models.Group, error)
	DeleteGroup(ctx context.Context, userId, groupID uuid.UUID) error

	AddUserToGroup(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error
	RemoveUserFromGroup(ctx context.Context, groupID, userID uuid.UUID) error
}

type groupService struct {
	repo repository.GroupRepository
}

func NewGroupService(r repository.GroupRepository) GroupService {
	return &groupService{repo: r}
}

// UBAH MUTLAK: Parameter input jadi pointer ke Body aslinya
func (s *groupService) CreateGroup(ctx context.Context, ownerID uuid.UUID, input *api.CreateGroupReq) (*models.Group, error) {
	uniqMemberID := make(map[uuid.UUID]bool)
	uniqMemberID[ownerID] = true

	// Pastikan input.Memberids tidak nil sebelum di-loop (bawaan oapi-codegen)
	if input.Memberids != nil {
		for _, idStr := range input.Memberids {
			id, err := uuid.Parse(idStr)
			if err != nil {
				// STOP! Jangan di-continue. Tolak requestnya kalau datanya cacat.
				return nil, errors.New("terdapat format member ID yang tidak valid")
			}
			uniqMemberID[id] = true
		}
	}

	var members []models.GroupMember
	for userID := range uniqMemberID {
		role := models.GroupParticipant
		if userID == ownerID {
			role = models.GroupAdmin
		}
		members = append(members, models.GroupMember{
			UserID:      userID,
			MembersRole: role,
		})
	}

	//var wallets []models.Wallet
	//for

	// 1. RAKIT SEMUANYA DALAM SATU WADAH (GORM Association)
	// Deskripsi harus di-dereference dengan aman karena tipenya *string di OpenAPI
	var desc string
	if input.Description != nil {
		desc = *input.Description
	}

	memberCount := len(members)

	newGroup := models.Group{
		Name:        input.Name,
		Description: desc,
		OwnerID:     ownerID,
		MemberCount: memberCount,
		Wallet: []models.Wallet{
			{
				Name:     "Wallet " + input.Name,
				Balance:  0,
				Currency: "IDR",
			},
		},
		Members: members,
	}

	// 2. SAVE KE DB (Satu pemanggilan saja)
	createdGroup, err := s.repo.CreateGroupWithWalletAndMembers(ctx, &newGroup)
	if err != nil {
		return nil, err
	}

	// 3. KEMBALIKAN MODEL ASLINYA
	// Jangan mapping DTO di sini! Biarkan Handler yang mengurusnya.
	return createdGroup, nil
}

// UBAH TANDA TANGAN: kembalikan int64 murni untuk diserahkan ke Handler
func (s *groupService) GetAllGroups(ctx context.Context, limit, offset int) (*[]models.Group, int64, error) {
	// Tarik data dan total baris dari Repo
	groups, totalData, err := s.repo.GetAllGroups(ctx, limit, offset)
	if err != nil {
		// KEMBALIKAN 3 NILAI: nil untuk struct, 0 untuk angka, dan err
		return nil, 0, err
	}

	// KEMBALIKAN LANGSUNG. Biarkan Handler yang pusing ngurusin kalkulasi halaman UI.
	return groups, totalData, nil
}

func (s *groupService) GetMyGroups(ctx context.Context, limit, offset int, userId uuid.UUID) (*[]models.Group, int64, error) {
	groups, totalData, err := s.repo.GetMyGroups(ctx, userId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return groups, totalData, nil
}

func (s *groupService) GetGroupByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error) {
	// Implementasi logika untuk mendapatkan grup berdasarkan ID
	group, err := s.repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *groupService) UpdateGroup(ctx context.Context, groupID uuid.UUID, name string) (*models.Group, error) {
	// Implementasi logika untuk memperbarui nama grup
	return &models.Group{}, nil
}

func (s *groupService) DeleteGroup(ctx context.Context, userId, groupID uuid.UUID) error {
	isGroupAdmin, err := s.repo.IsGroupAdmin(ctx, userId, groupID)
	if err != nil {
		return err
	}
	if !isGroupAdmin {
		return errors.New("Cant delete this, youre not group admin")
	}
	return s.repo.DeleteGroup(ctx, groupID)
}

// services/group_service.go
func (s *groupService) AddUserToGroup(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error {
	// 1. (Opsional) Cek dulu Group-nya ada gak?
	// _, err := s.repo.GetGroupByID(groupID)
	// if err != nil { return errors.New("group not found") }

	// 2. Mapping Logic (Business Logic)
	var members []models.GroupMember
	for _, uid := range userIDs {
		members = append(members, models.GroupMember{
			GroupID:     groupID,
			UserID:      uid,
			MembersRole: models.GroupParticipant, // Enaknya di Service: Bisa atur default role disini
			// JoinedAt:    time.Now(),              // Atau set waktu join custom
		})
	}

	// 3. Panggil Repo buat nyimpen
	return s.repo.CreateMembers(ctx, members)
}

func (s *groupService) RemoveUserFromGroup(ctx context.Context, groupID, userID uuid.UUID) error {
	return s.repo.RemoveUserFromGroup(ctx, groupID, userID)
}
