package dto

import (
	"time"
)

type DoctorSignUpRequest struct {
	UserID         string  `json:"user_id" validate:"required,uuid"`
	Specialization string  `json:"specialization" validate:"required"`
	LicenseNumber  string  `json:"license_number" validate:"required"`
	DepartmentID   string  `json:"department_id" validate:"required,uuid"`
	ConsultationFee float64 `json:"consultation_fee" validate:"required,gt=0"`
}

type DoctorResponse struct {
	DoctorID       string    `json:"doctor_id"`
	UserID         string    `json:"user_id"`
	Specialization string    `json:"specialization"`
	LicenseNumber  string    `json:"license_number"`
	DepartmentID   string    `json:"department_id"`
	ConsultationFee float64   `json:"consultation_fee"`
	IsAvailable    bool      `json:"is_available"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
