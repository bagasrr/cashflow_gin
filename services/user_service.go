package services

import (
	"cashflow_gin/dto/response"
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type UserService interface {
	FindAllUser(ctx context.Context) (*[]response.UserResponse, error)
	GetMyProfile(ctx context.Context, id uuid.UUID) (*response.UserResponse, error)

	FindUserByID(ctx context.Context, targetID uuid.UUID, requestorRole models.UserRole) (*response.UserResponse, error)
	UpdateUser(ctx context.Context, requestorID uuid.UUID, targetUser models.User) (*response.UserResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) FindAllUser(ctx context.Context) (*[]response.UserResponse, error) {
	users, err := s.repo.FindAllUser(ctx)
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

func (s *userService) GetMyProfile(ctx context.Context, id uuid.UUID) (*response.UserResponse, error) {
	user, err := s.repo.FindMyProfile(ctx, id)
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

func (s *userService) FindUserByID(ctx context.Context, targetID uuid.UUID, requestorRole models.UserRole) (*response.UserResponse, error) {
	if requestorRole != models.RoleAdmin {
		return nil, errors.New("forbidden: acces denied")
	}

	user, err := s.repo.FindUserByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		UserRole: user.UserRole.String(),
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, requestorID uuid.UUID, targetUser models.User) (*response.UserResponse, error) {
	requestor, err := s.repo.FindUserByID(ctx, requestorID)
	if err != nil {
		return nil, err
	}

	if requestor.UserRole != models.RoleAdmin && requestor.ID != targetUser.ID {
		return nil, errors.New("forbidden: access denied")
	}

	fmt.Println("TargetUser : ", targetUser)
	fmt.Println("Requestor : ", requestor)

	updatedUser, err := s.repo.UpdateUser(ctx, targetUser)
	if err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:       updatedUser.ID.String(),
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
		UserRole: updatedUser.UserRole.String(),
	}, nil
}
