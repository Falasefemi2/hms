package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/repository"
	"github.com/google/uuid"
)

type AppointmentService struct {
	appointmentRepo *repository.AppointmentRepository
	patientRepo     *repository.PatientRepository
	doctorRepo      *repository.DoctorRepository
}

func NewAppointmentService(appointmentRepo *repository.AppointmentRepository, patientRepo *repository.PatientRepository, doctorRepo *repository.DoctorRepository) *AppointmentService {
	return &AppointmentService{
		appointmentRepo: appointmentRepo,
		patientRepo:     patientRepo,
		doctorRepo:      doctorRepo,
	}
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
	// Validate patient exists
	_, err := s.patientRepo.GetByPatientID(ctx, appointment.PatientID)
	if err != nil {
		return nil, errors.New("patient not found")
	}

	// Validate doctor exists
	_, err = s.doctorRepo.GetDoctorID(ctx, appointment.DoctorID)
	if err != nil {
		return nil, errors.New("doctor not found")
	}

	// Validate appointment date is in the future
	if appointment.AppointmentDate.Before(time.Now()) {
		return nil, errors.New("appointment date must be in the future")
	}

	createdAppointment, err := s.appointmentRepo.Create(ctx, appointment)
	if err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return createdAppointment, nil
}

func (s *AppointmentService) GetAppointmentByID(ctx context.Context, appointmentID uuid.UUID) (*models.Appointment, error) {
	appointment, err := s.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return appointment, nil
}

func (s *AppointmentService) GetAppointmentsByPatientID(ctx context.Context, patientID uuid.UUID) ([]*models.Appointment, error) {
	appointments, err := s.appointmentRepo.GetByPatientID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %w", err)
	}

	return appointments, nil
}

func (s *AppointmentService) GetAppointmentsByDoctorID(ctx context.Context, doctorID uuid.UUID) ([]*models.Appointment, error) {
	appointments, err := s.appointmentRepo.GetByDoctorID(ctx, doctorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %w", err)
	}

	return appointments, nil
}

func (s *AppointmentService) UpdateAppointment(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
	// Check if appointment exists
	existing, err := s.appointmentRepo.GetByID(ctx, appointment.AppointmentID)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	// Validate status transition
	if existing.Status == "COMPLETED" && appointment.Status != "COMPLETED" {
		return nil, errors.New("cannot change status of completed appointment")
	}

	if existing.Status == "CANCELLED" {
		return nil, errors.New("cannot update cancelled appointment")
	}

	updatedAppointment, err := s.appointmentRepo.Update(ctx, appointment)
	if err != nil {
		return nil, fmt.Errorf("failed to update appointment: %w", err)
	}

	return updatedAppointment, nil
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, appointmentID uuid.UUID) error {
	appointment, err := s.appointmentRepo.GetByID(ctx, appointmentID)
	if err != nil {
		return errors.New("appointment not found")
	}

	if appointment.Status == "COMPLETED" {
		return errors.New("cannot delete completed appointment")
	}

	err = s.appointmentRepo.Delete(ctx, appointmentID)
	if err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	return nil
}
