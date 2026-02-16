package repositories

import (
	"context"
	"database/sql"
	"ecom-appz/internal/models"
	"errors"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository{
	return  &UserRepository{
		db: db,
	}
}

// Create inserts a new user
func(r *UserRepository) Create(ctx context.Context, user *models.User) error{
	query := `
	INSERT INTO users (email, password, role, is_active)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`
	return  r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Role,
		user.IsActive,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

// GetByEmail fetches user by email
func (r *UserRepository)GetByEmail(ctx context.Context, email string)(*models.User, error){
	user := &models.User{}

	query :=`
			SELECT id, email, password, role, is_active, created_at, updated_at
			FROM users
			WHERE email = $1
	`
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows{
		return nil, errors.New("user not found")
	}

	return  user, err
}

// GetByID fetches user by ID

func (r * UserRepository) GetByID(ctx context.Context, id string) (*models.User, error){
	user := &models.User{}

	query := `
			SELECT id, email, role, is_active, created_at, updated_at
			FROM users
			WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows{
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *UserRepository) UpdateProfile(user *models.User) error{
	query :=`UPDATE users
	SET full_name = $1,
	email = $2,
	phone = $3
	Updated_at = NOW()
	WHERE id = $4`
	_, err := r.db.Exec(query, user.Fullname, user.Phone, user.Email, user.ID)
	return err
}

