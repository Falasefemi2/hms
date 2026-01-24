package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/hms/internal/models"
)

type ConsultationRepository struct {
	pool *pgxpool.Pool
}

func NewConsultationRepository(pool *pgxpool.Pool) *ConsultationRepository {
	return &ConsultationRepository{
		pool: pool,
	}
}

func (r *ConsultationRepository) Create(ctx context.Context, consultation *models.Consultation) (*models.Consultation, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
    INSERT INTO consultations (consultation_id, appointment_id, patient_id, doctor_id, diagnosis, notes, is_editable)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING created_at
`
	err := r.pool.QueryRow(ctx, query,
		consultation.ConsultationID,
		consultation.AppointmentID,
		consultation.PatientID,
		consultation.DoctorID,
		consultation.Diagnosis,
		consultation.Notes,
		consultation.IsEditable,
	).Scan(&consultation.CreatedAt)

	if err != nil {
		return nil, err
	}

	return consultation, nil
}

func (r *ConsultationRepository) GetByID(ctx context.Context, consultationID uuid.UUID) (*models.Consultation, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT consultation_id, appointment_id, patient_id, doctor_id, diagnosis, notes, created_at, is_editable
		FROM consultations
		WHERE consultation_id = $1
	`

	var consultation models.Consultation
	err := r.pool.QueryRow(ctx, query, consultationID).Scan(
		&consultation.ConsultationID,
		&consultation.AppointmentID,
		&consultation.PatientID,
		&consultation.DoctorID,
		&consultation.Diagnosis,
		&consultation.Notes,
		&consultation.CreatedAt,
		&consultation.IsEditable,
	)
	if err != nil {
		return nil, err
	}

	return &consultation, nil
}

func (r *ConsultationRepository) GetByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*models.Consultation, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT consultation_id, appointment_id, patient_id, doctor_id, diagnosis, notes, created_at, is_editable
		FROM consultations
		WHERE appointment_id = $1
	`

	var consultation models.Consultation
	err := r.pool.QueryRow(ctx, query, appointmentID).Scan(
		&consultation.ConsultationID,
		&consultation.AppointmentID,
		&consultation.PatientID,
		&consultation.DoctorID,
		&consultation.Diagnosis,
		&consultation.Notes,
		&consultation.CreatedAt,
		&consultation.IsEditable,
	)
	if err != nil {
		return nil, err
	}

	return &consultation, nil
}

func (r *ConsultationRepository) GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*models.Consultation, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT consultation_id, appointment_id, patient_id, doctor_id, diagnosis, notes, created_at, is_editable
		FROM consultations
		WHERE patient_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consultations []*models.Consultation
	for rows.Next() {
		var consultation models.Consultation
		err := rows.Scan(
			&consultation.ConsultationID,
			&consultation.AppointmentID,
			&consultation.PatientID,
			&consultation.DoctorID,
			&consultation.Diagnosis,
			&consultation.Notes,
			&consultation.CreatedAt,
			&consultation.IsEditable,
		)
		if err != nil {
			return nil, err
		}
		consultations = append(consultations, &consultation)
	}

	return consultations, nil
}

func (r *ConsultationRepository) Update(ctx context.Context, consultation *models.Consultation) (*models.Consultation, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
    UPDATE consultations
    SET diagnosis = $2, notes = $3
    WHERE consultation_id = $1
`
	_, err := r.pool.Exec(ctx, query,
		consultation.ConsultationID,
		consultation.Diagnosis,
		consultation.Notes,
	)

	if err != nil {
		return nil, err
	}

	return consultation, nil
}
