package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/hms/internal/models"
)

type DoctorRepository struct {
	pool *pgxpool.Pool
}

func NewDoctorRepository(pool *pgxpool.Pool) *DoctorRepository {
	return &DoctorRepository{
		pool: pool,
	}
}

func (r *DoctorRepository) Create(ctx context.Context, doctor *models.Doctor) (*models.Doctor, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
    INSERT INTO doctors (doctor_id, user_id, department_id, specialization, license_number, consultation_fee)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING created_at, updated_at
`
	err := r.pool.QueryRow(ctx, query,
		doctor.DoctorID,
		doctor.UserID,
		doctor.DepartmentID,
		doctor.Specialization,
		doctor.LicenseNumber,
		doctor.ConsultationFee,
	).Scan(&doctor.CreatedAt, &doctor.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return doctor, nil
}

func (r *DoctorRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Doctor, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT doctor_id, user_id, department_id, specialization, license_number, created_at, updated_at
		FROM doctors
		WHERE user_id = $1
	`

	var doctor models.Doctor
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&doctor.DoctorID,
		&doctor.UserID,
		&doctor.DepartmentID,
		&doctor.Specialization,
		&doctor.LicenseNumber,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (r *DoctorRepository) GetDoctorID(ctx context.Context, doctorID uuid.UUID) (*models.Doctor, error) {

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT doctor_id, user_id, department_id, specialization, license_number, created_at, updated_at
		FROM doctors
		WHERE doctor_id = $1
	`

	var doctor models.Doctor
	err := r.pool.QueryRow(ctx, query, doctorID).Scan(
		&doctor.DoctorID,
		&doctor.UserID,
		&doctor.DepartmentID,
		&doctor.Specialization,
		&doctor.LicenseNumber,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &doctor, nil

}
