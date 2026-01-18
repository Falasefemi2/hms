package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
	"github.com/falasefemi2/hms/internal/utils"
)

var validAdminRoles = map[string]bool{
	"DOCTOR": true,
	"NURSE":  true,
	"ADMIN":  true,
}

const patientRole = "PATIENT"

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

func (us *UserService) CreatePatientUser(ctx context.Context, username, email, password, firstName, lastName, phone string) (*models.User, error) {
	if err := validatePatientInput(username, email, password, firstName, lastName); err != nil {
		return nil, err
	}
	exists, err := us.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, errors.New("email already registered")
	}
	exists, err = us.repo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, errors.New("username already taken")
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    &firstName,
		LastName:     &lastName,
		Phone:        &phone,
		Role:         patientRole,
		IsActive:     true,
	}

	createdUser, err := us.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create patientuser: %w", err)
	}
	return createdUser, nil
}

func (us *UserService) CreateAdminUser(ctx context.Context, username, email, password, firstName, lastName, phone, role string) (*models.User, error) {
	if err := validateAdminUserInput(username, email, password, firstName, lastName, role); err != nil {
		return nil, err
	}

	exists, err := us.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	exists, err = us.repo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, errors.New("username already taken")
	}

	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    &firstName,
		LastName:     &lastName,
		Phone:        &phone,
		Role:         role,
		IsActive:     true,
	}

	createdUser, err := us.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func (us *UserService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
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

func (us *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, int64, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	users, err := us.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := us.repo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}

func validatePatientInput(username, email, password, firstName, lastName string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}

	if email == "" {
		return errors.New("email is required")
	}
	if !isValidEmail(email) {
		return errors.New("invalid email format")
	}

	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	if firstName == "" {
		return errors.New("first name is required")
	}

	if lastName == "" {
		return errors.New("last name is required")
	}

	return nil
}

func validateAdminUserInput(username, email, password, firstName, lastName, role string) error {
	// First validate basic fields
	if err := validatePatientInput(username, email, password, firstName, lastName); err != nil {
		return err
	}

	// Validate role
	if role == "" {
		return errors.New("role is required")
	}
	if !validAdminRoles[role] {
		return fmt.Errorf("invalid role: %s. must be one of: DOCTOR, NURSE, ADMIN", role)
	}

	// Prevent patient creation through admin endpoint
	if role == patientRole {
		return errors.New("patients must self-register using the patient signup endpoint")
	}

	return nil
}

// validateUserInput validates basic user fields
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

	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
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

// isValidEmail checks if email format is valid
func isValidEmail(email string) bool {
	// Simple email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (us *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.ComparePassword(user.PasswordHash, password) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJwt(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
