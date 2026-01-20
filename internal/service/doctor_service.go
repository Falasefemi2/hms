
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
)

type DoctorService struct {
	doctorRepo *repository.DoctorRepository
	userRepo   *repository.UserRepository
}

func NewDoctorService(doctorRepo *repository.DoctorRepository, userRepo *repository.UserRepository) *DoctorService {
	return &DoctorService{
		doctorRepo: doctorRepo,
		userRepo:   userRepo,
	}
}

func (s *DoctorService) CreateDoctor(ctx context.Context, doctor *models.Doctor) (*models.Doctor, error) {
	user, err := s.userRepo.GetByID(ctx, doctor.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Role != "DOCTOR" {
		return nil, errors.New("user is not a doctor")
	}

	existingDoctor, err := s.doctorRepo.GetByUserID(ctx, user.ID)
	if err == nil && existingDoctor != nil {
		return nil, errors.New("doctor already exists for this user")
	}

	createdDoctor, err := s.doctorRepo.Create(ctx, doctor)
	if err != nil {
		return nil, fmt.Errorf("failed to create doctor: %w", err)
	}

	return createdDoctor, nil
}
