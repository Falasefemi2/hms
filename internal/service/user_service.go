package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := validateUserInput(user); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}
	createdUser, err := us.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return createdUser, nil
}

func (us *UserService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := us.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	user, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (us *UserService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := validateUserInput(user); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	updatedUser, err := us.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

func (us *UserService) DeleteUser(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return errors.New("invalid user ID")
	}

	err := us.repo.Delete(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (us *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := us.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func validateUserInput(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	if user.Username == "" {
		return errors.New("username is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.PasswordHash == "" {
		return errors.New("password hash is required")
	}

	if user.FirstName == nil || *user.FirstName == "" {
		return errors.New("first name is required")
	}

	if user.LastName == nil || *user.LastName == "" {
		return errors.New("last name is required")
	}

	if user.Role == "" {
		return errors.New("role is required")
	}

	return nil
}
