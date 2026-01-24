package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/hms/internal/models"
)

type AppointmentRepository struct {
	pool *pgxpool.Pool
}

func NewAppointmentRepository(pool *pgxpool.Pool) *AppointmentRepository {
	return &AppointmentRepository{
		pool: pool,
	}
}

func (r *AppointmentRepository) Create(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
    INSERT INTO appointments (appointment_id, patient_id, doctor_id, appointment_date, duration_minutes, status, notes)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING created_at, updated_at
`
	err := r.pool.QueryRow(ctx, query,
		appointment.AppointmentID,
		appointment.PatientID,
		appointment.DoctorID,
		appointment.AppointmentDate,
		appointment.DurationMinutes,
		appointment.Status,
		appointment.Notes,
	).Scan(&appointment.CreatedAt, &appointment.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return appointment, nil
}

func (r *AppointmentRepository) GetByID(ctx context.Context, appointmentID uuid.UUID) (*models.Appointment, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT appointment_id, patient_id, doctor_id, appointment_date, duration_minutes, status, notes, created_at, updated_at
		FROM appointments
		WHERE appointment_id = $1
	`

	var appointment models.Appointment
	err := r.pool.QueryRow(ctx, query, appointmentID).Scan(
		&appointment.AppointmentID,
		&appointment.PatientID,
		&appointment.DoctorID,
		&appointment.AppointmentDate,
		&appointment.DurationMinutes,
		&appointment.Status,
		&appointment.Notes,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &appointment, nil
}

func (r *AppointmentRepository) GetByPatientID(ctx context.Context, patientID uuid.UUID) ([]*models.Appointment, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT appointment_id, patient_id, doctor_id, appointment_date, duration_minutes, status, notes, created_at, updated_at
		FROM appointments
		WHERE patient_id = $1
		ORDER BY appointment_date DESC
	`

	rows, err := r.pool.Query(ctx, query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*models.Appointment
	for rows.Next() {
		var appointment models.Appointment
		err := rows.Scan(
			&appointment.AppointmentID,
			&appointment.PatientID,
			&appointment.DoctorID,
			&appointment.AppointmentDate,
			&appointment.DurationMinutes,
			&appointment.Status,
			&appointment.Notes,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}

func (r *AppointmentRepository) GetByDoctorID(ctx context.Context, doctorID uuid.UUID) ([]*models.Appointment, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT appointment_id, patient_id, doctor_id, appointment_date, duration_minutes, status, notes, created_at, updated_at
		FROM appointments
		WHERE doctor_id = $1
		ORDER BY appointment_date DESC
	`

	rows, err := r.pool.Query(ctx, query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*models.Appointment
	for rows.Next() {
		var appointment models.Appointment
		err := rows.Scan(
			&appointment.AppointmentID,
			&appointment.PatientID,
			&appointment.DoctorID,
			&appointment.AppointmentDate,
			&appointment.DurationMinutes,
			&appointment.Status,
			&appointment.Notes,
			&appointment.CreatedAt,
			&appointment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		appointments = append(appointments, &appointment)
	}

	return appointments, nil
}

func (r *AppointmentRepository) Update(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
    UPDATE appointments
    SET appointment_date = $2, duration_minutes = $3, status = $4, notes = $5, updated_at = CURRENT_TIMESTAMP
    WHERE appointment_id = $1
    RETURNING updated_at
`
	err := r.pool.QueryRow(ctx, query,
		appointment.AppointmentID,
		appointment.AppointmentDate,
		appointment.DurationMinutes,
		appointment.Status,
		appointment.Notes,
	).Scan(&appointment.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return appointment, nil
}

func (r *AppointmentRepository) Delete(ctx context.Context, appointmentID uuid.UUID) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `DELETE FROM appointments WHERE appointment_id = $1`

	_, err := r.pool.Exec(ctx, query, appointmentID)
	return err
}
