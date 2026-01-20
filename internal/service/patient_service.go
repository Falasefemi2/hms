package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
)

type PatientService struct {
	patientRepo *repository.PatientRepository
	userRepo    *repository.UserRepository
}

func NewPatientService(patientRepo *repository.PatientRepository, userRepo *repository.UserRepository) *PatientService {
	return &PatientService{
		patientRepo: patientRepo,
		userRepo:    userRepo,
	}
}

func (p *PatientService) PatientProfile(ctx context.Context, patient *models.Patient) (*models.Patient, error) {
	user, err := p.userRepo.GetByID(ctx, patient.UserID.String())
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Role != "PATIENT" {
		return nil, errors.New("user is not a patient")
	}
	exisitngPatient, err := p.patientRepo.GetByUserID(ctx, user.ID)
	if err == nil && exisitngPatient != nil {
		return nil, errors.New("patient already exists for this user")
	}
	patientProfile, err := p.patientRepo.PatientProfile(ctx, patient)
	if err != nil {
		return nil, fmt.Errorf("failed to create patient profile: %w", err)
	}
	return patientProfile, nil
}
