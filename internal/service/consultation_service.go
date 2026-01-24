package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
	"github.com/google/uuid"
)

type ConsultationService struct {
	consultationRepo *repository.ConsultationRepository
	appointmentRepo  *repository.AppointmentRepository
	patientRepo      *repository.PatientRepository
	doctorRepo       *repository.DoctorRepository
}

func NewConsultationService(consultationRepo *repository.ConsultationRepository, appointmentRepo *repository.AppointmentRepository, patientRepo *repository.PatientRepository, doctorRepo *repository.DoctorRepository) *ConsultationService {
	return &ConsultationService{
		consultationRepo: consultationRepo,
		appointmentRepo:  appointmentRepo,
		patientRepo:      patientRepo,
		doctorRepo:       doctorRepo,
	}
}

func (s *ConsultationService) CreateConsultation(ctx context.Context, consultation *models.Consultation) (*models.Consultation, error) {
	// Validate appointment exists and is completed
	appointment, err := s.appointmentRepo.GetByID(ctx, consultation.AppointmentID)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	if appointment.Status != "COMPLETED" {
		return nil, errors.New("consultation can only be created for completed appointments")
	}

	// Check if consultation already exists for this appointment
	existing, err := s.consultationRepo.GetByAppointmentID(ctx, consultation.AppointmentID)
	if err == nil && existing != nil {
		return nil, errors.New("consultation already exists for this appointment")
	}

	// Validate patient and doctor match the appointment
	if consultation.PatientID != appointment.PatientID || consultation.DoctorID != appointment.DoctorID {
		return nil, errors.New("patient and doctor must match the appointment")
	}

	createdConsultation, err := s.consultationRepo.Create(ctx, consultation)
	if err != nil {
		return nil, fmt.Errorf("failed to create consultation: %w", err)
	}

	return createdConsultation, nil
}

func (s *ConsultationService) GetConsultationByID(ctx context.Context, consultationID uuid.UUID) (*models.Consultation, error) {
	consultation, err := s.consultationRepo.GetByID(ctx, consultationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get consultation: %w", err)
	}

	return consultation, nil
}

func (s *ConsultationService) GetConsultationByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*models.Consultation, error) {
	consultation, err := s.consultationRepo.GetByAppointmentID(ctx, appointmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get consultation: %w", err)
	}

	return consultation, nil
}

func (s *ConsultationService) GetConsultationsByPatientID(ctx context.Context, patientID uuid.UUID) ([]*models.Consultation, error) {
	consultations, err := s.consultationRepo.GetByPatientID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get consultations: %w", err)
	}

	return consultations, nil
}

func (s *ConsultationService) UpdateConsultation(ctx context.Context, consultation *models.Consultation) (*models.Consultation, error) {
	// Check if consultation exists
	existing, err := s.consultationRepo.GetByID(ctx, consultation.ConsultationID)
	if err != nil {
		return nil, errors.New("consultation not found")
	}

	// Check if editable
	if !existing.IsEditable {
		return nil, errors.New("consultation is not editable")
	}

	updatedConsultation, err := s.consultationRepo.Update(ctx, consultation)
	if err != nil {
		return nil, fmt.Errorf("failed to update consultation: %w", err)
	}

	return updatedConsultation, nil
}
