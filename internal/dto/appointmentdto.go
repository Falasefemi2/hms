package dto

import (
	"time"
)

type CreateAppointmentRequest struct {
	PatientID       string `json:"patient_id" validate:"required,uuid"`
	DoctorID        string `json:"doctor_id" validate:"required,uuid"`
	AppointmentDate string `json:"appointment_date" validate:"required"`
	DurationMinutes int    `json:"duration_minutes" validate:"required,min=1"`
	Notes           string `json:"notes"`
}

type UpdateAppointmentRequest struct {
	AppointmentDate string `json:"appointment_date" validate:"required"`
	DurationMinutes int    `json:"duration_minutes" validate:"required,min=1"`
	Status          string `json:"status" validate:"required,oneof=PENDING CONFIRMED COMPLETED CANCELLED"`
	Notes           string `json:"notes"`
}

type AppointmentResponse struct {
	AppointmentID   string    `json:"appointment_id"`
	PatientID       string    `json:"patient_id"`
	DoctorID        string    `json:"doctor_id"`
	AppointmentDate time.Time `json:"appointment_date"`
	DurationMinutes int       `json:"duration_minutes"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
