package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
)

type NurseSerivce struct {
	nurseRepo *repository.NurseRepository
	userRepo  *repository.UserRepository
}

func NewNurseService(nurseRepo *repository.NurseRepository, userRepo *repository.UserRepository) *NurseSerivce {
	return &NurseSerivce{
		nurseRepo: nurseRepo,
		userRepo:  userRepo,
	}
}

func (n *NurseSerivce) CreateNurse(ctx context.Context, nurse *models.Nurse) (*models.Nurse, error) {

	user, err := n.userRepo.GetByID(ctx, nurse.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role != "NURSE" {
		return nil, errors.New("user is not a nurse")
	}

	existingDoctor, err := n.nurseRepo.GetByUserID(ctx, user.ID)
	if err == nil && existingDoctor != nil {
		return nil, errors.New("nurse already exists for this user")
	}

	createdNurse, err := n.nurseRepo.Create(ctx, nurse)
	if err != nil {
		return nil, fmt.Errorf("failed to create nurse: %w", err)
	}

	return createdNurse, nil
}
