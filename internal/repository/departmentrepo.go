package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/hms/internal/models"
)

type DepartmentRepository struct {
	pool *pgxpool.Pool
}

type UpdateDepartmentRequest struct {
	Name        *string
	Description *string
	IsActive    *bool
}

type PaginationParams struct {
	Limit  int
	Offset int
}

type PaginatedResponse struct {
	Data       []*models.Department
	TotalCount int
}

func NewDepartmentRepository(pool *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{
		pool: pool,
	}
}

func (dept *DepartmentRepository) CreateDepartment(ctx context.Context, department *models.Department) (*models.Department, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `

	INSERT INTO departments (
		name,
		description
	)
	VALUES ($1, $2)
	RETURNING 
	department_id,
	name,
	description,
	is_active,
	created_at,
	updated_at
	`
	row := dept.pool.QueryRow(
		ctx,
		query,
		department.Name,
		department.Description,
	)

	var created models.Department

	err := row.Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.IsActive,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (dept *DepartmentRepository) GetByID(ctx context.Context, deptID string) (*models.Department, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
	SELECT
	department_id,
	name,
	description,
	is_active,
	created_at,
	updated_at
	FROM departments 
	WHERE department_id = $1
	`
	row := dept.pool.QueryRow(ctx, query, deptID)

	var department models.Department
	err := row.Scan(
		&department.ID,
		&department.Name,
		&department.Description,
		&department.IsActive,
		&department.CreatedAt,
		&department.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("department not found")
		}
		return nil, err
	}
	return &department, nil
}

func (dept *DepartmentRepository) GetAll(ctx context.Context, pagination PaginationParams) (*PaginatedResponse, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	if pagination.Limit <= 0 {
		pagination.Limit = 10
	}
	if pagination.Offset < 0 {
		pagination.Offset = 0
	}

	countQuery := `
	SELECT COUNT(*)
	FROM departments
	WHERE is_active = true
	`

	var totalCount int
	err := dept.pool.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	query := `
	SELECT
	department_id,
	name,
	description,
	is_active,
	created_at,
	updated_at
	FROM departments 
	WHERE is_active = true
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`
	rows, err := dept.pool.Query(ctx, query, pagination.Limit, pagination.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	departments := make([]*models.Department, 0)

	for rows.Next() {
		var department models.Department
		err := rows.Scan(
			&department.ID,
			&department.Name,
			&department.Description,
			&department.IsActive,
			&department.CreatedAt,
			&department.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		departments = append(departments, &department)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &PaginatedResponse{
		Data:       departments,
		TotalCount: totalCount,
	}, nil
}

func (dept *DepartmentRepository) UpdateDepartment(ctx context.Context, deptID string, request *UpdateDepartmentRequest) (*models.Department, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	existing, err := dept.GetByID(ctx, deptID)
	if err != nil {
		return nil, err
	}
	query := `UPDATE departments SET `
	args := []interface{}{}
	argCounter := 1

	hasUpdates := false

	if request.Name != nil {
		if hasUpdates {
			query += ", "
		}
		query += `name = $` + fmt.Sprintf("%d", argCounter)
		args = append(args, *request.Name)
		argCounter++
		hasUpdates = true
	}

	if request.Description != nil {
		if hasUpdates {
			query += ", "
		}
		query += `description = $` + fmt.Sprintf("%d", argCounter)
		args = append(args, *request.Description)
		argCounter++
		hasUpdates = true
	}
	if request.IsActive != nil {
		if hasUpdates {
			query += ", "
		}
		query += `is_active = $` + fmt.Sprintf("%d", argCounter)
		args = append(args, *request.IsActive)
		argCounter++
		hasUpdates = true
	}
	if !hasUpdates {
		return existing, nil
	}
	query += `, updated_at = CURRENT_TIMESTAMP `
	query += `WHERE department_id = $` + fmt.Sprintf("%d", argCounter)
	args = append(args, deptID)

	query += `
	RETURNING 
	department_id,
	name,
	description,
	is_active,
	created_at,
	updated_at
	`

	row := dept.pool.QueryRow(ctx, query, args...)

	var updated models.Department
	err = row.Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.IsActive,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (dept *DepartmentRepository) DeleteDepartment(ctx context.Context, deptID string) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	_, err := dept.GetByID(ctx, deptID)
	if err != nil {
		return err
	}

	query := `
	UPDATE departments 
	SET is_active = false, updated_at = CURRENT_TIMESTAMP
	WHERE department_id = $1
	`

	result, err := dept.pool.Exec(ctx, query, deptID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("department not found")
	}

	return nil
}
