package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/hms/internal/models"
)

type HospitalConfigRepository struct {
	pool *pgxpool.Pool
}

func NewHospitalConfigRepository(pool *pgxpool.Pool) *HospitalConfigRepository {
	return &HospitalConfigRepository{
		pool: pool,
	}
}

func (r *HospitalConfigRepository) Create(ctx context.Context, config *models.HospitalConfig) (*models.HospitalConfig, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
    INSERT INTO hospital_config (config_id, working_hours_start, working_hours_end, appointment_duration_minutes, max_same_day_cancellation_hours, enable_patient_self_registration)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING created_at, updated_at
`
	err := r.pool.QueryRow(ctx, query,
		config.ConfigID,
		config.WorkingHoursStart,
		config.WorkingHoursEnd,
		config.AppointmentDurationMinutes,
		config.MaxSameDayCancellationHours,
		config.EnablePatientSelfRegistration,
	).Scan(&config.CreatedAt, &config.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (r *HospitalConfigRepository) GetByID(ctx context.Context, configID uuid.UUID) (*models.HospitalConfig, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT config_id, working_hours_start, working_hours_end, appointment_duration_minutes, max_same_day_cancellation_hours, enable_patient_self_registration, created_at, updated_at
		FROM hospital_config
		WHERE config_id = $1
	`

	var config models.HospitalConfig
	err := r.pool.QueryRow(ctx, query, configID).Scan(
		&config.ConfigID,
		&config.WorkingHoursStart,
		&config.WorkingHoursEnd,
		&config.AppointmentDurationMinutes,
		&config.MaxSameDayCancellationHours,
		&config.EnablePatientSelfRegistration,
		&config.CreatedAt,
		&config.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (r *HospitalConfigRepository) GetAll(ctx context.Context) ([]*models.HospitalConfig, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		SELECT config_id, working_hours_start, working_hours_end, appointment_duration_minutes, max_same_day_cancellation_hours, enable_patient_self_registration, created_at, updated_at
		FROM hospital_config
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*models.HospitalConfig
	for rows.Next() {
		var config models.HospitalConfig
		err := rows.Scan(
			&config.ConfigID,
			&config.WorkingHoursStart,
			&config.WorkingHoursEnd,
			&config.AppointmentDurationMinutes,
			&config.MaxSameDayCancellationHours,
			&config.EnablePatientSelfRegistration,
			&config.CreatedAt,
			&config.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &config)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

func (r *HospitalConfigRepository) Update(ctx context.Context, config *models.HospitalConfig) (*models.HospitalConfig, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
    UPDATE hospital_config
    SET working_hours_start = $2, working_hours_end = $3, appointment_duration_minutes = $4, max_same_day_cancellation_hours = $5, enable_patient_self_registration = $6, updated_at = CURRENT_TIMESTAMP
    WHERE config_id = $1
    RETURNING updated_at
`
	err := r.pool.QueryRow(ctx, query,
		config.ConfigID,
		config.WorkingHoursStart,
		config.WorkingHoursEnd,
		config.AppointmentDurationMinutes,
		config.MaxSameDayCancellationHours,
		config.EnablePatientSelfRegistration,
	).Scan(&config.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (r *HospitalConfigRepository) Delete(ctx context.Context, configID uuid.UUID) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
    DELETE FROM hospital_config
    WHERE config_id = $1
`
	_, err := r.pool.Exec(ctx, query, configID)
	return err
}
