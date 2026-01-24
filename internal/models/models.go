package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	FirstName    *string
	LastName     *string
	Phone        *string
	Role         string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Department struct {
	ID          uuid.UUID
	Name        string
	Description string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Doctor struct {
	DoctorID        uuid.UUID
	UserID          uuid.UUID
	Specialization  string
	LicenseNumber   string
	DepartmentID    uuid.UUID
	ConsultationFee float64
	IsAvailable     bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Nurse struct {
	NurseID       uuid.UUID
	UserID        uuid.UUID
	DepartmentID  uuid.UUID
	Shift         string
	LicenseNumber string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Patient struct {
	PatientID             uuid.UUID
	UserID                uuid.UUID
	DateOfBirth           time.Time
	Gender                string
	BloodGroup            string
	EmergencyContactName  string
	EmergencyContactPhone string
	MedicalHistory        string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type Availability struct {
	AvailabilityID uuid.UUID
	DoctorID       uuid.UUID
	DayOfWeek      string
	StartTime      string
	EndTime        string
	MaxAppointment int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type HospitalConfig struct {
	ConfigID                      uuid.UUID
	WorkingHoursStart             string
	WorkingHoursEnd               string
	AppointmentDurationMinutes    int
	MaxSameDayCancellationHours   int
	EnablePatientSelfRegistration bool
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
}

type Appointment struct {
	AppointmentID   uuid.UUID
	PatientID       uuid.UUID
	DoctorID        uuid.UUID
	AppointmentDate time.Time
	DurationMinutes int
	Status          string
	Notes           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Consultation struct {
	ConsultationID uuid.UUID
	AppointmentID  uuid.UUID
	PatientID      uuid.UUID
	DoctorID       uuid.UUID
	Diagnosis      string
	Notes          string
	CreatedAt      time.Time
	IsEditable     bool
}
