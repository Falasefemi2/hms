package dto

import (
	"time"

	"github.com/google/uuid"
)

type PatientSignUp struct {
	UserID                string    `json:"user_id"`
	DateOfBirth           string    `json:"date_of_birth"`
	Gender                string    `json:"gender"`
	BloodGroup            string    `json:"blood_group"`
	EmergencyContactName  string    `json:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone"`
	MedicalHistory        string    `json:"medical_history"`
}

type PatientResponse struct {
	PatientID             uuid.UUID `json:"patient_id"`
	UserID                uuid.UUID `json:"user_id"`
	DateOfBirth           time.Time `json:"date_of_birth"`
	Gender                string    `json:"gender"`
	BloodGroup            string    `json:"blood_group"`
	EmergencyContactName  string    `json:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone"`
	MedicalHistory        string    `json:"medical_history"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
