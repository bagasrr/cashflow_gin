package services

import (
	"cashflow_gin/models"
	"cashflow_gin/repository"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type UserService interface {
	FindAllUser(ctx context.Context) (*[]models.User, error)
	GetMyProfile(ctx context.Context, id uuid.UUID) (*models.User, error)

	FindUserByID(ctx context.Context, targetID uuid.UUID, requestorRole models.UserRole) (*models.User, error)
	UpdateUser(ctx context.Context, requestorID uuid.UUID, targetUser models.User) (*models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) FindAllUser(ctx context.Context) (*[]models.User, error) {
	users, err := s.repo.FindAllUser(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) GetMyProfile(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindMyProfile(ctx, id)
	fmt.Println(user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil
}

func (s *userService) FindUserByID(ctx context.Context, targetID uuid.UUID, requestorRole models.UserRole) (*models.User, error) {
	if requestorRole != models.RoleAdmin {
		return nil, errors.New("forbidden: acces denied")
	}

	user, err := s.repo.FindUserByID(ctx, targetID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, requestorID uuid.UUID, targetUser models.User) (*models.User, error) {
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

	return updatedUser, nil
}
