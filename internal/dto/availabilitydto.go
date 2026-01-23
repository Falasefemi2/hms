package dto

import (
	"time"

	"github.com/google/uuid"
)

type AvailabilityRequest struct {
	DoctorID        uuid.UUID `json:"doctor_id"`
	DayOfWeek       string    `json:"day_of_week"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	MaxAppointments int       `json:"max_appointments"`
}

type AvailabilityResponse struct {
	AvailabilityID  uuid.UUID `json:"availability_id"`
	DoctorID        uuid.UUID `json:"doctor_id"`
	DayOfWeek       string    `json:"day_of_week"`
	StartTime       string    `json:"start_time"`
	EndTime         string    `json:"end_time"`
	MaxAppointments int       `json:"max_appointments"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
