package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/hms/internal/models"
)

type AvailabilityRepository struct {
	pool *pgxpool.Pool
}

func NewAvailabilityRepository(pool *pgxpool.Pool) *AvailabilityRepository {
	return &AvailabilityRepository{
		pool: pool,
	}
}

func (a *AvailabilityRepository) CreateAvailability(ctx context.Context, availability *models.Availability) (*models.Availability, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}
	query := `
	INSERT INTO doctor_availability (availability_id, doctor_id, day_of_week, start_time, end_time, max_appointments)
	VALUES ($1, $2, $3, $4::TIME, $5::TIME, $6)
	RETURNING created_at, updated_at
	`
	err := a.pool.QueryRow(ctx, query,
		availability.AvailabilityID,
		availability.DoctorID,
		availability.DayOfWeek,
		availability.StartTime,
		availability.EndTime,
		availability.MaxAppointment,
	).Scan(&availability.CreatedAt, &availability.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return availability, nil
}
