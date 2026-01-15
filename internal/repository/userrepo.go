package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/falasefemi2/hms/internal/models"
	"github.com/falasefemi2/hms/internal/utils"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (ur *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
	}

	query := `
		INSERT INTO users (
			username,
			email,
			password_hash,
			first_name,
			last_name,
			phone,
			role
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			user_id,
			username,
			email,
			first_name,
			last_name,
			phone,
			role,
			is_active,
			created_at,
			updated_at
	`

	row := ur.pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
	)

	var created models.User

	err := row.Scan(
		&created.ID,
		&created.Username,
		&created.Email,
		&created.FirstName,
		&created.LastName,
		&created.Phone,
		&created.Role,
		&created.IsActive,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (ur *UserRepository) GetByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
		SELECT
			user_id,
			username,
			email,
			password_hash,
			first_name,
			last_name,
			phone,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE user_id = $1
	`

	row := ur.pool.QueryRow(ctx, query, userID)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT
			user_id,
			username,
			email,
			password_hash,
			first_name,
			last_name,
			phone,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	row := ur.pool.QueryRow(ctx, query, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	query := `
		UPDATE users
		SET
			username = $1,
			email = $2,
			first_name = $3,
			last_name = $4,
			phone = $5,
			role = $6,
			is_active = $7,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $8
		RETURNING
			user_id,
			username,
			email,
			first_name,
			last_name,
			phone,
			role,
			is_active,
			created_at,
			updated_at
	`

	row := ur.pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
		user.IsActive,
		user.ID,
	)

	var updated models.User
	err := row.Scan(
		&updated.ID,
		&updated.Username,
		&updated.Email,
		&updated.FirstName,
		&updated.LastName,
		&updated.Phone,
		&updated.Role,
		&updated.IsActive,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (ur *UserRepository) Delete(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE user_id = $1`
	commandTag, err := ur.pool.Exec(ctx, query, userID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return utils.ErrUserNotFound
	}

	return nil
}

func (ur *UserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT
			username,
			email,
			password_hash,
			first_name,
			last_name,
			phone,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := ur.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := ur.pool.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (ur *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := ur.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
