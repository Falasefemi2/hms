package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
)

type AvailabilityService struct {
	availabilityRepo *repository.AvailabilityRepository
	doctorRepo       *repository.DoctorRepository
}

func NewAvailabilityService(availabilityRepo *repository.AvailabilityRepository, doctorRepo *repository.DoctorRepository) *AvailabilityService {
	return &AvailabilityService{
		availabilityRepo: availabilityRepo,
		doctorRepo:       doctorRepo,
	}
}

func (a *AvailabilityService) CreateDoctorAvailability(ctx context.Context, availability *models.Availability) (*models.Availability, error) {
	doctor, err := a.doctorRepo.GetDoctorID(ctx, availability.DoctorID)
	if err != nil {
		return nil, errors.New("doctor not found")
	}
	if doctor == nil {
		return nil, errors.New("doctor does not exist")
	}

	doctorAvailability, err := a.availabilityRepo.CreateAvailability(ctx, availability)
	if err != nil {
		return nil, fmt.Errorf("failed to create availability: %w", err)
	}
	return doctorAvailability, nil
}
