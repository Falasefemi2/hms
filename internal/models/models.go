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
	DocotrID        uuid.UUID
	UserID          uuid.UUID
	Specialization  string
	LicenseNumber   string
	DepartmentID    uuid.UUID
	ConsultationFee float64
	IsAvailable     bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
