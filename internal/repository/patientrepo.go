package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/hms/internal/models"
)

type PatientRepository struct {
	pool *pgxpool.Pool
}

func NewPatientRepository(pool *pgxpool.Pool) *PatientRepository {
	return &PatientRepository{
		pool: pool,
	}
}

func (p *PatientRepository) PatientProfile(ctx context.Context, patient *models.Patient) (*models.Patient, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `

	INSERT INTO patients (patient_id, user_id, date_of_birth, gender, blood_group, emergency_contact_name, emergency_contact_phone, medical_history)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING created_at, updated_at
	`
	err := p.pool.QueryRow(ctx, query,
		patient.PatientID,
		patient.UserID,
		patient.DateOfBirth,
		patient.Gender,
		patient.BloodGroup,
		patient.EmergencyContactName,
		patient.EmergencyContactPhone,
		patient.MedicalHistory,
	).Scan(&patient.CreatedAt, &patient.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return patient, nil
}

func (p *PatientRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Patient, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT patient_id, user_id, date_of_birth, gender, blood_group, emergency_contact_name, emergency_contact_phone, medical_history, created_at, updated_at
		FROM patients
		WHERE user_id = $1
	`

	var patient models.Patient
	err := p.pool.QueryRow(ctx, query, userID).Scan(
		&patient.PatientID,
		&patient.UserID,
		&patient.DateOfBirth,
		&patient.Gender,
		&patient.BloodGroup,
		&patient.EmergencyContactName,
		&patient.EmergencyContactPhone,
		&patient.MedicalHistory,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &patient, nil
}
