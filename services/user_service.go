package services

import (
	"cashflow_gin/dto/response"
	"cashflow_gin/repository"

	"github.com/google/uuid"
)

type UserService interface {
	FindAllUser() (*[]response.UserResponse, error)
	GetMyProfile(id uuid.UUID) (*response.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) FindAllUser() (*[]response.UserResponse, error) {
	users, err := s.repo.FindAllUser()
	if err != nil {
		return nil, err
	}
	var userRes []response.UserResponse
	for _, u := range users {
		var WalletRes []response.WalletResponse
		for _, w := range u.Wallets {
			WalletRes = append(WalletRes, response.WalletResponse{
				ID:               w.ID,
				Name:             w.Name,
				Balance:          w.Balance,
				TransactionCount: w.TransactionCount,
			})
		}
		userRes = append(userRes, response.UserResponse{
			ID:       u.ID.String(),
			Username: u.Username,
			Email:    u.Email,
			UserRole: u.UserRole.String(),
			Wallets:  WalletRes,
		})
	}
	return &userRes, nil
}

func (s *userService) GetMyProfile(id uuid.UUID) (*response.UserResponse, error) {
	user, err := s.repo.FindMyProfile(id)
	if err != nil {
		return nil, err
	}
	var UserRes *response.UserResponse
	var WalletRes []response.WalletResponse

	for _, w := range user.Wallets {
		WalletRes = append(WalletRes, response.WalletResponse{
			ID:               w.ID,
			Name:             w.Name,
			Balance:          w.Balance,
			TransactionCount: w.TransactionCount,
		})
	}

	UserRes = &response.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		UserRole: user.UserRole.String(),
		Wallets:  WalletRes,
	}
	return UserRes, nil
}
