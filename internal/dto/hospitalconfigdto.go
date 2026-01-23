package dto

import (
	"time"
)

type CreateHospitalConfigRequest struct {
	WorkingHoursStart             string `json:"working_hours_start" validate:"required"`
	WorkingHoursEnd               string `json:"working_hours_end" validate:"required"`
	AppointmentDurationMinutes    int    `json:"appointment_duration_minutes" validate:"required,min=1"`
	MaxSameDayCancellationHours   int    `json:"max_same_day_cancellation_hours" validate:"required,min=0"`
	EnablePatientSelfRegistration *bool  `json:"enable_patient_self_registration"` // pointer to distinguish false from unset
}

type UpdateHospitalConfigRequest struct {
	ConfigID                      string `json:"config_id" validate:"required,uuid"`
	WorkingHoursStart             string `json:"working_hours_start" validate:"required"`
	WorkingHoursEnd               string `json:"working_hours_end" validate:"required"`
	AppointmentDurationMinutes    int    `json:"appointment_duration_minutes" validate:"required,min=1"`
	MaxSameDayCancellationHours   int    `json:"max_same_day_cancellation_hours" validate:"required,min=0"`
	EnablePatientSelfRegistration *bool  `json:"enable_patient_self_registration"`
}

type HospitalConfigResponse struct {
	ConfigID                      string    `json:"config_id"`
	WorkingHoursStart             string    `json:"working_hours_start"`
	WorkingHoursEnd               string    `json:"working_hours_end"`
	AppointmentDurationMinutes    int       `json:"appointment_duration_minutes"`
	MaxSameDayCancellationHours   int       `json:"max_same_day_cancellation_hours"`
	EnablePatientSelfRegistration bool      `json:"enable_patient_self_registration"`
	CreatedAt                     time.Time `json:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at"`
}
