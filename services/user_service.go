package services

import (
	"cashflow_gin/dto/response"
	"cashflow_gin/models"
	"cashflow_gin/repository"

	"github.com/google/uuid"
)

type UserService interface {
	FindAllUser() ([]response.UserResponse, error)
	GetMyProfile(id uuid.UUID) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) FindAllUser() ([]response.UserResponse, error) {
	users, err := s.repo.FindAllUser()
	if err != nil {
		return nil, err
	}
	var userRes []response.UserResponse
	for _, u := range users {
		var WalletRes []response.WalletResponse
		for _, w := range u.Wallets {
			var txRes []response.TransactionResponse
			for _, t := range w.Transaction {
				txRes = append(txRes, response.TransactionResponse{
					ID:          t.ID.String(),
					Title:       t.Title,
					Amount:      t.Amount,
					Date:        t.Date,
					Description: t.Description,
					Category: response.CategoryResponse{
						Name: t.Category.Name,
						Type: t.Category.Type,
					},
				})
			}
			WalletRes = append(WalletRes, response.WalletResponse{
				ID:           w.ID,
				Name:         w.Name,
				Balance:      w.Balance,
				Transactions: txRes,
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
	return userRes, nil
}

func (s *userService) GetMyProfile(id uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindMyProfile(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
