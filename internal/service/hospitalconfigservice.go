package service

import (
	"context"
	"fmt"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
	"github.com/google/uuid"
)

type HospitalConfigService struct {
	hospitalConfigRepo *repository.HospitalConfigRepository
}

func NewHospitalConfigService(hospitalConfigRepo *repository.HospitalConfigRepository) *HospitalConfigService {
	return &HospitalConfigService{
		hospitalConfigRepo: hospitalConfigRepo,
	}
}

func (s *HospitalConfigService) CreateHospitalConfig(ctx context.Context, config *models.HospitalConfig) (*models.HospitalConfig, error) {
	// Generate new UUID if not provided
	if config.ConfigID == uuid.Nil {
		config.ConfigID = uuid.New()
	}

	// Set defaults
	if config.AppointmentDurationMinutes == 0 {
		config.AppointmentDurationMinutes = 30
	}
	if config.MaxSameDayCancellationHours == 0 {
		config.MaxSameDayCancellationHours = 24
	}
	// EnablePatientSelfRegistration defaults to true in schema

	createdConfig, err := s.hospitalConfigRepo.Create(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create hospital config: %w", err)
	}

	return createdConfig, nil
}

func (s *HospitalConfigService) GetHospitalConfigByID(ctx context.Context, configID uuid.UUID) (*models.HospitalConfig, error) {
	config, err := s.hospitalConfigRepo.GetByID(ctx, configID)
	if err != nil {
		return nil, fmt.Errorf("failed to get hospital config: %w", err)
	}

	return config, nil
}

func (s *HospitalConfigService) GetAllHospitalConfigs(ctx context.Context) ([]*models.HospitalConfig, error) {
	configs, err := s.hospitalConfigRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get hospital configs: %w", err)
	}

	return configs, nil
}

func (s *HospitalConfigService) UpdateHospitalConfig(ctx context.Context, config *models.HospitalConfig) (*models.HospitalConfig, error) {
	updatedConfig, err := s.hospitalConfigRepo.Update(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to update hospital config: %w", err)
	}

	return updatedConfig, nil
}

func (s *HospitalConfigService) DeleteHospitalConfig(ctx context.Context, configID uuid.UUID) error {
	err := s.hospitalConfigRepo.Delete(ctx, configID)
	if err != nil {
		return fmt.Errorf("failed to delete hospital config: %w", err)
	}

	return nil
}
