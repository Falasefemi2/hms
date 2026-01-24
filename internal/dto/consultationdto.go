package dto

import (
	"time"
)

type CreateConsultationRequest struct {
	AppointmentID string `json:"appointment_id" validate:"required,uuid"`
	PatientID     string `json:"patient_id" validate:"required,uuid"`
	DoctorID      string `json:"doctor_id" validate:"required,uuid"`
	Diagnosis     string `json:"diagnosis" validate:"required"`
	Notes         string `json:"notes"`
}

type UpdateConsultationRequest struct {
	Diagnosis string `json:"diagnosis" validate:"required"`
	Notes     string `json:"notes"`
}

type ConsultationResponse struct {
	ConsultationID string    `json:"consultation_id"`
	AppointmentID  string    `json:"appointment_id"`
	PatientID      string    `json:"patient_id"`
	DoctorID       string    `json:"doctor_id"`
	Diagnosis      string    `json:"diagnosis"`
	Notes          string    `json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	IsEditable     bool      `json:"is_editable"`
}
