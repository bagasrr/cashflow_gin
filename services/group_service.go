package services

import (
	"cashflow_gin/dto/request"
	"cashflow_gin/dto/response"
	"cashflow_gin/models"
	"cashflow_gin/repository"

	"github.com/google/uuid"
)

type GroupService interface {
	CreateGroup(ownerID uuid.UUID, input request.CreateGroupRequest) (*response.GroupResponse, error)
	GetAllGroups() (*[]response.GroupResponse, error)

	GetGroupByID(groupID uuid.UUID) (*response.GroupResponse, error)
	UpdateGroup(groupID uuid.UUID, name string) (*response.GroupResponse, error)
	DeleteGroup(groupID uuid.UUID) error

	AddUserToGroup(groupID uuid.UUID, userIDs []uuid.UUID) error
	RemoveUserFromGroup(groupID, userID uuid.UUID) error
}

type groupService struct {
	repo repository.GroupRepository
}

func NewGroupService(r repository.GroupRepository) GroupService {
	return &groupService{repo: r}
}

func (s *groupService) CreateGroup(ownerID uuid.UUID, input request.CreateGroupRequest) (*response.GroupResponse, error) {
	uniqMemberID := make(map[uuid.UUID]bool)
	uniqMemberID[ownerID] = true

	for _, idStr := range input.MemberIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		uniqMemberID[id] = true
	}

	// 1. Inisialisasi Group
	newGroup := models.Group{
		Name:        input.Name,
		Description: input.Description,
		OwnerID:     ownerID, // Set Owner
	}

	// 2. Inisialisasi Wallet untuk Group
	// Ingat model Wallet kita sebelumnya (UserID null, GroupID terisi)
	newWallet := models.Wallet{
		Name:    "Wallet " + input.Name,
		Balance: 0,
		// GroupID akan diisi otomatis oleh GORM lewat relasi, atau bisa manual nanti
	}

	// 3. Siapkan List Member
	// Member PERTAMA wajib si OWNER itu sendiri (Role: ADMIN)
	// members := []models.GroupMember{
	// 	{
	// 		UserID:      ownerID,
	// 		MembersRole: models.GroupAdmin, // Owner otomatis jadi Admin
	// 	},
	// }

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

	// (Opsional) Tambahin member lain dari input jika ada
	// for _, invitedIDStr := range input.MemberIDs {
	// 	invitedID, _ := uuid.Parse(invitedIDStr)
	// 	members = append(members, models.GroupMember{
	// 		UserID:      invitedID,
	// 		MembersRole: 2,
	// 	})
	// }

	// 4. SAVE KE DB (Panggil Repo yang Transactional)
	// Kita kirim pointer biar ID-nya ke-generate dan balik ke variable ini
	err := s.repo.CreateGroupWithWalletAndMembers(&newGroup, &newWallet, &members)
	if err != nil {
		return &response.GroupResponse{}, err
	}

	// 5. MAPPING KE RESPONSE (Manual Mapping biar Rapi)
	// Ambil data member yang baru disimpan buat ditampilkan
	var memberResponses []response.GroupMemberResponse
	for _, m := range members {
		memberResponses = append(memberResponses, response.GroupMemberResponse{
			ID:     m.ID.String(),
			UserID: m.UserID.String(),
			Role:   m.MembersRole.String(),
			// Username idealnya di-preload di repo atau fetch ulang,
			// disini kita skip dulu atau set kosong
			Username: m.User.Username,
		})
	}

	res := response.GroupResponse{
		ID:          newGroup.ID.String(),
		Name:        newGroup.Name,
		Description: newGroup.Description,
		Wallet: response.WalletResponse{
			ID:      newWallet.ID,
			Name:    newWallet.Name,
			Balance: newWallet.Balance,
		},
		Members: memberResponses,
	}

	return &res, nil
}

func (s *groupService) GetAllGroups() (*[]response.GroupResponse, error) {
	// Implementasi logika untuk mendapatkan semua grup
	groups, err := s.repo.GetAllGroups()
	if err != nil {
		return nil, err
	}

	var groupResponses []response.GroupResponse
	for _, group := range *groups {

		var walletRes response.WalletResponse

		for _, w := range group.Wallet {
			walletRes = response.WalletResponse{
				ID:      w.ID,
				Name:    w.Name,
				Balance: w.Balance,
			}
			break
		}
		groupResponses = append(groupResponses, response.GroupResponse{
			ID:           group.ID.String(),
			Name:         group.Name,
			Description:  group.Description,
			Wallet:       walletRes,
			TotalMembers: group.MemberCount,
		})
	}

	return &groupResponses, nil
}

func (s *groupService) GetGroupByID(groupID uuid.UUID) (*response.GroupResponse, error) {
	// Implementasi logika untuk mendapatkan grup berdasarkan ID

	group, err := s.repo.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}

	// Mapping ke response
	var memberResponses []response.GroupMemberResponse
	for _, m := range group.Members {
		memberResponses = append(memberResponses, response.GroupMemberResponse{
			ID:       m.ID.String(),
			UserID:   m.UserID.String(),
			Role:     m.MembersRole.String(),
			Username: m.User.Username, // Idealnya di-preload di repo atau fetch ulang
		})
	}

	var walletRes response.WalletResponse
	for _, w := range group.Wallet {
		walletRes = response.WalletResponse{
			ID:      w.ID,
			Name:    w.Name,
			Balance: w.Balance,
			// GroupID: w.GroupID,
		}
		break // Asumsi cuma 1 wallet per group, keluar setelah dapat yang pertama
	}

	res := response.GroupResponse{
		ID:           group.ID.String(),
		Name:         group.Name,
		Description:  group.Description,
		Wallet:       walletRes,
		Members:      memberResponses,
		TotalMembers: group.MemberCount,
	}

	return &res, nil
}

func (s *groupService) UpdateGroup(groupID uuid.UUID, name string) (*response.GroupResponse, error) {
	// Implementasi logika untuk memperbarui nama grup
	return &response.GroupResponse{}, nil
}

func (s *groupService) DeleteGroup(groupID uuid.UUID) error {
	return s.repo.DeleteGroup(groupID)
}

// services/group_service.go
func (s *groupService) AddUserToGroup(groupID uuid.UUID, userIDs []uuid.UUID) error {
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
	return s.repo.CreateMembers(members)
}

func (s *groupService) RemoveUserFromGroup(groupID, userID uuid.UUID) error {
	return s.repo.RemoveUserFromGroup(groupID, userID)
}
