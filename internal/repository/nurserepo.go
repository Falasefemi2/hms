package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/hms/internal/models"
)

type NurseRepository struct {
	pool *pgxpool.Pool
}

func NewNurseRepository(pool *pgxpool.Pool) *NurseRepository {
	return &NurseRepository{
		pool: pool,
	}
}

func (n *NurseRepository) Create(ctx context.Context, nurse *models.Nurse) (*models.Nurse, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
	INSERT INTO nurses (nurse_id, user_id, department_id, shift, license_number)
	VALUES ($1,$2,$3,$4,$5)
  RETURNING created_at, updated_at
	`

	err := n.pool.QueryRow(ctx, query,
		nurse.NurseID,
		nurse.UserID,
		nurse.DepartmentID,
		nurse.Shift,
		nurse.LicenseNumber,
	).Scan(&nurse.CreatedAt, &nurse.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return nurse, nil
}

func (n *NurseRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Nurse, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT nurse_id, user_id, department_id, shift, license_number, created_at, updated_at
		FROM nurses
		WHERE user_id = $1
		`

	var nurse models.Nurse
	err := n.pool.QueryRow(ctx, query, userID).Scan(
		&nurse.NurseID,
		&nurse.UserID,
		&nurse.DepartmentID,
		&nurse.Shift,
		&nurse.LicenseNumber,
		&nurse.CreatedAt,
		&nurse.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &nurse, nil
}
