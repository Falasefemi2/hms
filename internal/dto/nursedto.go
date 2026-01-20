package dto

import "time"

type NurseSignupRequest struct {
	UserID        string `json:"user_id" validate:"required,uuid"`
	Shift         string `json:"shift" validate:"required"`
	LicenseNumber string `json:"license_number" validate:"required"`
	DepartmentID  string `json:"department_id" validate:"required,uuid"`
}

type NurseResponse struct {
	NurseID       string    `json:"nurse_id"`
	UserID        string    `json:"user_id"`
	LicenseNumber string    `json:"license_number"`
	DepartmentID  string    `json:"department_id"`
	Shift         string    `json:"shift"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
